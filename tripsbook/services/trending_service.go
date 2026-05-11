package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TrendingService handles trending algorithm calculations
type TrendingService struct {
	db *gorm.DB
}

// NewTrendingService creates a new trending service instance
func NewTrendingService(db *gorm.DB) *TrendingService {
	return &TrendingService{db: db}
}

// TrendingServiceResponse represents a trending service
type TrendingServiceResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Distance      string   `json:"distance"`
	Rating        float64  `json:"rating"`
	Price         string   `json:"price"`
	Trending      bool     `json:"trending"`
	Image         string   `json:"image"`
	WeeklyChange  float64  `json:"weeklyChange"`
	TrendingBadge string   `json:"trendingBadge"`
	Location      Location `json:"location"`
}

// TrendingCategory represents a trending category
type TrendingCategory struct {
	Name   string `json:"name"`
	Change string `json:"change"`
	Color  string `json:"color"`
	Icon   string `json:"icon"`
}

// TrendingRequest represents a trending search request
type TrendingRequest struct {
	Location string `json:"location"` // city name
	Period   string `json:"period"`   // "week" or "month"
	Category string `json:"category"` // optional filter
	Limit    int    `json:"limit"`    // pagination
}

// GetTrendingServices gets trending services based on various metrics
func (s *TrendingService) GetTrendingServices(req TrendingRequest) ([]TrendingServiceResponse, error) {
	var services []TrendingServiceResponse

	// Calculate time range
	endDate := time.Now()
	var startDate time.Time

	if req.Period == "month" {
		startDate = endDate.AddDate(0, -1, 0)
	} else {
		startDate = endDate.AddDate(0, 0, -7)
	}

	// Get trending services with booking growth and engagement metrics
	query := s.db.Raw(`
		WITH current_period AS (
			SELECT 
				sp.id,
				sp.business_name as name,
				sc.name as type,
				sp.average_rating as rating,
				sp.logo as image,
				sp.latitude,
				sp.longitude,
				sp.address,
				COUNT(DISTINCT b.id) as current_bookings,
				COUNT(DISTINCT b.customer_id) as unique_customers,
				AVG(r.rating) as avg_review_score,
				COUNT(DISTINCT r.id) as review_count
			FROM service_providers sp
			JOIN service_sub_categories sc ON sp.sub_category_id = sc.id
			LEFT JOIN bookings b ON sp.id = b.provider_id 
				AND b.created_at BETWEEN ? AND ? 
				AND b.status = 'completed'
			LEFT JOIN reviews r ON sp.id = r.provider_id 
				AND r.created_at BETWEEN ? AND ?
			WHERE sp.is_verified = true AND sp.is_available = true
			GROUP BY sp.id, sp.business_name, sc.name, sp.average_rating, sp.logo, sp.latitude, sp.longitude, sp.address
		),
		previous_period AS (
			SELECT 
				sp.id,
				COUNT(DISTINCT b.id) as previous_bookings,
				COUNT(DISTINCT b.customer_id) as previous_customers,
				AVG(r.rating) as prev_review_score,
				COUNT(DISTINCT r.id) as prev_review_count
			FROM service_providers sp
			LEFT JOIN bookings b ON sp.id = b.provider_id 
				AND b.created_at BETWEEN ? AND ? 
				AND b.status = 'completed'
			LEFT JOIN reviews r ON sp.id = r.provider_id 
				AND r.created_at BETWEEN ? AND ?
			WHERE sp.is_verified = true AND sp.is_available = true
			GROUP BY sp.id
		),
		search_metrics AS (
			SELECT 
				sp.id,
				COUNT(*) as search_count,
				COUNT(DISTINCT user_fingerprint) as unique_searchers
			FROM service_providers sp
			LEFT JOIN search_analytics sa ON sp.id = sa.service_id 
				AND sa.created_at BETWEEN ? AND ?
			WHERE sp.is_verified = true AND sp.is_available = true
			GROUP BY sp.id
		)
		SELECT 
			cp.*,
			COALESCE(pp.previous_bookings, 0) as previous_bookings,
			COALESCE(pp.previous_customers, 0) as previous_customers,
			COALESCE(sm.search_count, 0) as search_count,
			COALESCE(sm.unique_searchers, 0) as unique_searchers,
			-- Calculate trending score based on multiple factors
			CASE 
				WHEN COALESCE(pp.previous_bookings, 0) = 0 THEN 
					(cp.current_bookings * 10) + (cp.unique_customers * 5) + (COALESCE(sm.search_count, 0) * 2)
				ELSE 
					((cp.current_bookings - pp.previous_bookings) / NULLIF(pp.previous_bookings, 0)) * 100 +
					((cp.unique_customers - pp.previous_customers) / NULLIF(pp.previous_customers, 0)) * 50 +
					(COALESCE(sm.search_count, 0) * 2)
			END as trending_score,
			-- Calculate weekly change percentage
			CASE 
				WHEN COALESCE(pp.previous_bookings, 0) = 0 THEN 
					CASE WHEN cp.current_bookings > 0 THEN 100 ELSE 0 END
				ELSE 
					((cp.current_bookings - pp.previous_bookings) / NULLIF(pp.previous_bookings, 0)) * 100
			END as weekly_change
		FROM current_period cp
		LEFT JOIN previous_period pp ON cp.id = pp.id
		LEFT JOIN search_metrics sm ON cp.id = sm.id
		WHERE cp.current_bookings > 0 OR COALESCE(sm.search_count, 0) > 0
		ORDER BY trending_score DESC
		LIMIT ?
	`, startDate, endDate, startDate, endDate,
		startDate.AddDate(0, 0, -7), startDate, startDate.AddDate(0, 0, -7), startDate,
		startDate, endDate, req.Limit)

	var results []struct {
		ID                string  `json:"id"`
		Name              string  `json:"name"`
		Type              string  `json:"type"`
		Rating            float64 `json:"rating"`
		Image             string  `json:"image"`
		Latitude          float64 `json:"latitude"`
		Longitude         float64 `json:"longitude"`
		Address           string  `json:"address"`
		CurrentBookings   int64   `json:"current_bookings"`
		UniqueCustomers   int64   `json:"unique_customers"`
		AvgReviewScore    float64 `json:"avg_review_score"`
		ReviewCount       int64   `json:"review_count"`
		PreviousBookings  int64   `json:"previous_bookings"`
		PreviousCustomers int64   `json:"previous_customers"`
		SearchCount       int64   `json:"search_count"`
		UniqueSearchers   int64   `json:"unique_searchers"`
		TrendingScore     float64 `json:"trending_score"`
		WeeklyChange      float64 `json:"weekly_change"`
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get trending services: %w", err)
	}

	// Process results and format response
	for _, result := range results {
		service := TrendingServiceResponse{
			ID:           result.ID,
			Name:         result.Name,
			Type:         result.Type,
			Rating:       result.Rating,
			Image:        result.Image,
			WeeklyChange: result.WeeklyChange,
			Trending:     true,
			Location: Location{
				Latitude:  result.Latitude,
				Longitude: result.Longitude,
			},
			Price:         s.getServicePrice(result.ID),
			Distance:      "2.3 km", // Would calculate based on user location
			TrendingBadge: s.getTrendingBadge(result.WeeklyChange),
		}

		services = append(services, service)
	}

	return services, nil
}

// GetTrendingCategories gets trending categories with growth percentages
func (s *TrendingService) GetTrendingCategories(req TrendingRequest) ([]TrendingCategory, error) {
	var categories []TrendingCategory

	// Calculate time range
	endDate := time.Now()
	var startDate time.Time

	if req.Period == "month" {
		startDate = endDate.AddDate(0, -1, 0)
	} else {
		startDate = endDate.AddDate(0, 0, -7)
	}

	query := s.db.Raw(`
		WITH current_period AS (
			SELECT 
				c.name,
				c.color,
				c.icon,
				COUNT(DISTINCT b.id) as current_bookings,
				COUNT(DISTINCT sp.id) as active_providers
			FROM service_categories c
			JOIN service_sub_categories sc ON c.id = sc.category_id
			JOIN service_providers sp ON sc.id = sp.sub_category_id
			LEFT JOIN bookings b ON sp.id = b.provider_id 
				AND b.created_at BETWEEN ? AND ? 
				AND b.status = 'completed'
			WHERE sp.is_verified = true AND sp.is_available = true
			GROUP BY c.id, c.name, c.color, c.icon
		),
		previous_period AS (
			SELECT 
				c.name,
				COUNT(DISTINCT b.id) as previous_bookings
			FROM service_categories c
			JOIN service_sub_categories sc ON c.id = sc.category_id
			JOIN service_providers sp ON sc.id = sp.sub_category_id
			LEFT JOIN bookings b ON sp.id = b.provider_id 
				AND b.created_at BETWEEN ? AND ? 
				AND b.status = 'completed'
			WHERE sp.is_verified = true AND sp.is_available = true
			GROUP BY c.id, c.name
		)
		SELECT 
			cp.name,
			cp.color,
			cp.icon,
			COALESCE(pp.previous_bookings, 0) as previous_bookings,
			cp.current_bookings,
			cp.active_providers,
			-- Calculate growth percentage
			CASE 
				WHEN COALESCE(pp.previous_bookings, 0) = 0 THEN 
					CASE WHEN cp.current_bookings > 0 THEN 100 ELSE 0 END
				ELSE 
					((cp.current_bookings - pp.previous_bookings) / NULLIF(pp.previous_bookings, 0)) * 100
			END as growth_percentage
		FROM current_period cp
		LEFT JOIN previous_period pp ON cp.name = pp.name
		WHERE cp.current_bookings > 0
		ORDER BY growth_percentage DESC
	`, startDate, endDate, startDate.AddDate(0, 0, -7), startDate)

	var results []struct {
		Name             string  `json:"name"`
		Color            string  `json:"color"`
		Icon             string  `json:"icon"`
		PreviousBookings int64   `json:"previous_bookings"`
		CurrentBookings  int64   `json:"current_bookings"`
		ActiveProviders  int64   `json:"active_providers"`
		GrowthPercentage float64 `json:"growth_percentage"`
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get trending categories: %w", err)
	}

	for _, result := range results {
		category := TrendingCategory{
			Name:   result.Name,
			Color:  result.Color,
			Icon:   result.Icon,
			Change: fmt.Sprintf("%+.0f%%", result.GrowthPercentage),
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// TrackSearchAnalytics tracks search interactions for trending calculations
func (s *TrendingService) TrackSearchAnalytics(serviceID string, userFingerprint string, query string) error {
	analytics := SearchAnalytics{
		ServiceID:       serviceID,
		UserFingerprint: userFingerprint,
		Query:           query,
		CreatedAt:       time.Now(),
	}

	return s.db.Create(&analytics).Error
}

// TrackServiceView tracks when a user views a service
func (s *TrendingService) TrackServiceView(serviceID string, userFingerprint string) error {
	view := ServiceView{
		ServiceID:       serviceID,
		UserFingerprint: userFingerprint,
		CreatedAt:       time.Now(),
	}

	return s.db.Create(&view).Error
}

// GetTrendingScore calculates a comprehensive trending score for a service
func (s *TrendingService) GetTrendingScore(serviceID string, period string) (float64, error) {
	// Calculate time range
	endDate := time.Now()
	var startDate time.Time

	if period == "month" {
		startDate = endDate.AddDate(0, -1, 0)
	} else {
		startDate = endDate.AddDate(0, 0, -7)
	}

	var score float64

	// Get booking growth
	var bookingGrowth float64
	s.db.Raw(`
		WITH current_bookings AS (
			SELECT COUNT(*) as count
			FROM bookings 
			WHERE provider_id = ? AND created_at BETWEEN ? AND ? AND status = 'completed'
		),
		previous_bookings AS (
			SELECT COUNT(*) as count
			FROM bookings 
			WHERE provider_id = ? AND created_at BETWEEN ? AND ? AND status = 'completed'
		)
		SELECT 
			CASE 
				WHEN (SELECT count FROM previous_bookings) = 0 THEN 
					CASE WHEN (SELECT count FROM current_bookings) > 0 THEN 100 ELSE 0 END
				ELSE 
					((SELECT count FROM current_bookings) - (SELECT count FROM previous_bookings)) / 
					NULLIF((SELECT count FROM previous_bookings), 0) * 100
			END
	`, serviceID, startDate, endDate, serviceID, startDate.AddDate(0, 0, -7), startDate).Scan(&bookingGrowth)

	// Get search volume
	var searchCount int64
	s.db.Model(&SearchAnalytics{}).
		Where("service_id = ? AND created_at BETWEEN ? AND ?", serviceID, startDate, endDate).
		Count(&searchCount)

	// Get review activity
	var reviewCount int64
	s.db.Table("reviews").
		Where("provider_id = ? AND created_at BETWEEN ? AND ?", serviceID, startDate, endDate).
		Count(&reviewCount)

	// Calculate composite score
	score = (bookingGrowth * 0.6) + (float64(searchCount) * 0.3) + (float64(reviewCount) * 0.1)

	return score, nil
}

// Helper functions

// getTrendingBadge determines the trending badge based on growth percentage
func (s *TrendingService) getTrendingBadge(weeklyChange float64) string {
	if weeklyChange >= 50 {
		return "🔥 Hot"
	} else if weeklyChange >= 20 {
		return "📈 Rising"
	} else if weeklyChange >= 10 {
		return "⭐ Popular"
	} else if weeklyChange >= 0 {
		return "🔸 Trending"
	} else {
		return "📉 Declining"
	}
}

// getServicePrice gets the price range for a service
func (s *TrendingService) getServicePrice(serviceID string) string {
	var price string
	s.db.Table("services").
		Select("CASE WHEN base_price < 1000 THEN '$' WHEN base_price < 5000 THEN '$$' WHEN base_price < 20000 THEN '$$$' ELSE '$$$$' END").
		Where("provider_id = ? AND is_active = true", serviceID).
		Order("base_price ASC").
		Limit(1).
		Scan(&price)

	if price == "" {
		return "$$"
	}
	return price
}

// Database models for analytics

// SearchAnalytics tracks search interactions
type SearchAnalytics struct {
	ID              uint      `gorm:"primaryKey"`
	ServiceID       string    `json:"service_id"`
	UserFingerprint string    `json:"user_fingerprint"` // Anonymous user identifier
	Query           string    `json:"query"`
	CreatedAt       time.Time `json:"created_at"`
}

// ServiceView tracks service views
type ServiceView struct {
	ID              uint      `gorm:"primaryKey"`
	ServiceID       string    `json:"service_id"`
	UserFingerprint string    `json:"user_fingerprint"`
	CreatedAt       time.Time `json:"created_at"`
}

// UpdateTrendingData runs scheduled updates for trending calculations
func (s *TrendingService) UpdateTrendingData() error {
	// This would be called by a cron job to pre-calculate trending scores
	// and cache them for better performance

	// Update service trending scores
	services, err := s.GetTrendingServices(TrendingRequest{
		Period: "week",
		Limit:  100,
	})
	if err != nil {
		return fmt.Errorf("failed to update trending services: %w", err)
	}

	// Cache the results in Redis or similar
	// This is a placeholder for caching logic
	_ = services

	// Update category trending
	categories, err := s.GetTrendingCategories(TrendingRequest{
		Period: "week",
	})
	if err != nil {
		return fmt.Errorf("failed to update trending categories: %w", err)
	}

	// Cache category results
	_ = categories

	return nil
}

// GetPersonalizedTrending gets trending services personalized for a user
func (s *TrendingService) GetPersonalizedTrending(userFingerprint string, location string, limit int) ([]TrendingServiceResponse, error) {
	// Get user's search history
	var searchHistory []string
	s.db.Model(&SearchAnalytics{}).
		Where("user_fingerprint = ?", userFingerprint).
		Order("created_at DESC").
		Limit(10).
		Pluck("query", &searchHistory)

	// Get user's viewed services
	var viewedServices []string
	s.db.Model(&ServiceView{}).
		Where("user_fingerprint = ?", userFingerprint).
		Order("created_at DESC").
		Limit(20).
		Pluck("service_id", &viewedServices)

	// Get trending services and filter based on user preferences
	trendingReq := TrendingRequest{
		Location: location,
		Period:   "week",
		Limit:    limit * 2, // Get more to allow for filtering
	}

	trending, err := s.GetTrendingServices(trendingReq)
	if err != nil {
		return nil, err
	}

	// Filter and personalize based on user history
	var personalized []TrendingServiceResponse
	for _, service := range trending {
		// Simple personalization: prioritize services in categories user has searched for
		shouldInclude := true
		if len(searchHistory) > 0 {
			// Check if service type matches user's search history
			for _, query := range searchHistory {
				if contains(service.Type, query) || contains(query, service.Type) {
					shouldInclude = true
					break
				}
			}
		}

		if shouldInclude && len(personalized) < limit {
			personalized = append(personalized, service)
		}
	}

	return personalized, nil
}

// contains checks if a string contains another string (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

// containsSubstring checks if substring exists in string (simple implementation)
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
