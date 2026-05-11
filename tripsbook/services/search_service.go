package services

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

// SearchService handles search functionality
type SearchService struct {
	db *gorm.DB
}

// NewSearchService creates a new search service instance
func NewSearchService(db *gorm.DB) *SearchService {
	return &SearchService{db: db}
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query    string  `json:"query"`
	Location string  `json:"location"`
	Category string  `json:"category"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Limit    int     `json:"limit"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Services    []BaseServiceResult `json:"services"`
	Suggestions []string            `json:"suggestions"`
	Categories  []CategoryResult    `json:"categories"`
}

// BaseServiceResult represents a base service result
type BaseServiceResult struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Rating      float64 `json:"rating"`
	Price       string  `json:"price"`
	Distance    string  `json:"distance,omitempty"`
	MatchScore  int     `json:"match_score"`
}

// CategoryResult represents a category result
type CategoryResult struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Search performs a comprehensive search
func (s *SearchService) Search(req SearchRequest) (*SearchResponse, error) {
	response := &SearchResponse{
		Services:    []BaseServiceResult{},
		Suggestions: []string{},
		Categories:  []CategoryResult{},
	}

	// Build search query
	searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(req.Query))

	// Search services
	services, err := s.searchServices(searchQuery, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search services: %w", err)
	}
	response.Services = services

	// Get suggestions
	suggestions, err := s.GetSuggestions(req.Query, req.Location)
	if err == nil {
		response.Suggestions = suggestions
	}

	// Get categories
	categories, err := s.getCategoryCounts(searchQuery)
	if err == nil {
		response.Categories = categories
	}

	return response, nil
}

// searchServices performs service search
func (s *SearchService) searchServices(query string, req SearchRequest) ([]BaseServiceResult, error) {
	var services []BaseServiceResult

	baseQuery := s.db.Table("service_providers sp").
		Select(`
			sp.id,
			sp.business_name as name,
			sc.name as type,
			sp.description,
			sp.logo as image,
			sp.average_rating as rating,
			sp.address,
			sp.city,
			sp.state,
			sp.latitude,
			sp.longitude,
			sp.phone,
			sp.email,
			sp.website
		`).
		Joins("JOIN service_sub_categories sc ON sp.sub_category_id = sc.id").
		Where(`
			sp.is_verified = true AND 
			sp.is_available = true AND 
			(
				LOWER(sp.business_name) LIKE ? OR
				LOWER(sp.description) LIKE ? OR
				LOWER(sc.name) LIKE ? OR
				LOWER(sp.address) LIKE ? OR
				LOWER(sp.city) LIKE ?
			)
		`, query, query, query, query, query)

	// Add category filter
	if req.Category != "" && req.Category != "all" {
		baseQuery = baseQuery.Joins("JOIN service_categories c ON sc.category_id = c.id").
			Where("c.name = ?", req.Category)
	}

	// Add location filter
	if req.Location != "" {
		baseQuery = baseQuery.Where("sp.city = ?", req.Location)
	}

	// Add proximity filter if coordinates provided
	if req.Lat != 0 && req.Lng != 0 {
		// This would involve distance calculations
		// For simplicity, just filter by city for now
	}

	// Add limit
	if req.Limit > 0 {
		baseQuery = baseQuery.Limit(req.Limit)
	}

	var results []struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		Image       string  `json:"image"`
		Rating      float64 `json:"rating"`
		Address     string  `json:"address"`
		City        string  `json:"city"`
		State       string  `json:"state"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Phone       string  `json:"phone"`
		Email       string  `json:"email"`
		Website     string  `json:"website"`
	}

	if err := baseQuery.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}

	// Process results
	for _, result := range results {
		service := BaseServiceResult{
			ID:          result.ID,
			Name:        result.Name,
			Type:        result.Type,
			Description: result.Description,
			Image:       result.Image,
			Rating:      result.Rating,
			Price:       s.getServicePrice(result.ID),
			MatchScore:  s.calculateMatchScore(req.Query, result),
		}

		// Calculate distance if coordinates available
		if req.Lat != 0 && req.Lng != 0 {
			distance := s.calculateDistance(req.Lat, req.Lng, result.Latitude, result.Longitude)
			service.Distance = s.formatDistance(distance)
		}

		services = append(services, service)
	}

	return services, nil
}

// GetSuggestions gets search suggestions
func (s *SearchService) GetSuggestions(query, location string) ([]string, error) {
	var suggestions []string
	searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(query))

	// Get service name suggestions
	var serviceNames []string
	s.db.Table("service_providers").
		Select("DISTINCT business_name").
		Where("is_verified = true AND is_available = true AND LOWER(business_name) LIKE ?", searchQuery).
		Order("business_name").
		Limit(5).
		Pluck("business_name", &serviceNames)

	// Get category suggestions
	var categories []string
	s.db.Table("service_categories").
		Select("name").
		Where("LOWER(name) LIKE ?", searchQuery).
		Order("name").
		Limit(3).
		Pluck("name", &categories)

	// Combine suggestions
	suggestions = append(suggestions, serviceNames...)
	suggestions = append(suggestions, categories...)

	return suggestions, nil
}

// getCategoryCounts gets category counts for search
func (s *SearchService) getCategoryCounts(query string) ([]CategoryResult, error) {
	var categories []CategoryResult

	err := s.db.Raw(`
		SELECT 
			c.name,
			COUNT(DISTINCT sp.id) as count
		FROM service_categories c
		JOIN service_sub_categories sc ON c.id = sc.category_id
		JOIN service_providers sp ON sc.id = sp.sub_category_id
		WHERE 
			sp.is_verified = true AND 
			sp.is_available = true AND
			(
				LOWER(sp.business_name) LIKE ? OR
				LOWER(sp.description) LIKE ? OR
				LOWER(c.name) LIKE ? OR
				LOWER(sc.name) LIKE ?
			)
		GROUP BY c.id, c.name
		ORDER BY count DESC
	`, query, query, query, query).Scan(&categories).Error

	return categories, err
}

// Helper functions

// calculateMatchScore calculates how well a result matches the query
func (s *SearchService) calculateMatchScore(query string, result interface{}) int {
	// Simple match scoring based on exact name match
	if strings.Contains(strings.ToLower(result.(struct {
		ID          string
		Name        string
		Type        string
		Description string
		Image       string
		Rating      float64
		Address     string
		City        string
		State       string
		Latitude    float64
		Longitude   float64
		Phone       string
		Email       string
		Website     string
	}).Name), strings.ToLower(query)) {
		return 100
	}
	return 50
}

// getServicePrice gets price range for a service
func (s *SearchService) getServicePrice(serviceID string) string {
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

// calculateDistance calculates distance between two points
func (s *SearchService) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert to radians
	lat1Rad := lat1 * 3.14159 / 180
	lng1Rad := lng1 * 3.14159 / 180
	lat2Rad := lat2 * 3.14159 / 180
	lng2Rad := lng2 * 3.14159 / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlng := lng2Rad - lng1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// formatDistance formats distance for display
func (s *SearchService) formatDistance(distance float64) string {
	if distance < 1 {
		return fmt.Sprintf("%.0f m", distance*1000)
	}
	return fmt.Sprintf("%.1f km", distance)
}
