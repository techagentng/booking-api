package services

import (
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

// LocationService handles location-based service discovery
type LocationService struct {
	db *gorm.DB
}

// NewLocationService creates a new location service instance
func NewLocationService(db *gorm.DB) *LocationService {
	return &LocationService{db: db}
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// LocationBounds represents a rectangular area for map queries
type LocationBounds struct {
	Southwest Location `json:"southwest"`
	Northeast Location `json:"northeast"`
}

// NearbyServiceRequest represents a nearby service search request
type NearbyServiceRequest struct {
	Location  Location `json:"location" binding:"required"`
	RadiusKm  float64  `json:"radius" binding:"min=0.1,max=50"` // 100m to 50km
	Category  string   `json:"category"`                       // optional filter
	Limit     int      `json:"limit" binding:"min=1,max=100"`   // pagination
	Offset    int      `json:"offset" binding:"min=0"`
}

// NearbyServiceResponse represents a nearby service result
type NearbyServiceResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	Rating         float64   `json:"rating"`
	Price          string    `json:"price"`
	Distance       float64   `json:"distance"`        // in km
	DistanceText   string    `json:"distance_text"`   // "1.2 km"
	Location       Location  `json:"location"`
	Address        string    `json:"address"`
	OpenNow        bool      `json:"open_now"`
	AvailableNow   bool      `json:"available_now"`
	Features       []string  `json:"features"`
	Contact        Contact   `json:"contact"`
	OperatingHours *OperatingHours `json:"operating_hours,omitempty"`
}

// Contact represents service contact information
type Contact struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

// OperatingHours represents service operating hours
type OperatingHours struct {
	Open     string   `json:"open"`
	Close    string   `json:"close"`
	Days     []string `json:"days"`
	OpenNow  bool     `json:"open_now"`
}

// MapMarker represents a service for map display
type MapMarker struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Location Location `json:"location"`
	Price    string  `json:"price"`
	Rating   float64 `json:"rating"`
	Icon     string  `json:"icon"`
}

// MapCluster represents a cluster of nearby services
type MapCluster struct {
	ID         string  `json:"id"`
	Location   Location `json:"location"`
	Count      int     `json:"count"`
	Type       string  `json:"type"`
	Icon       string  `json:"icon"`
}

// MapViewResponse represents map view results
type MapViewResponse struct {
	Markers  []MapMarker  `json:"markers"`
	Clusters []MapCluster `json:"clusters"`
}

// GetNearbyServices finds services within a specified radius
func (s *LocationService) GetNearbyServices(req NearbyServiceRequest) ([]NearbyServiceResponse, error) {
	var services []NearbyServiceResponse

	// Calculate bounding box for the radius
	bounds := s.calculateBoundingBox(req.Location, req.RadiusKm)

	// Build query
	query := s.db.Table("service_providers sp").
		Select(`
			sp.id,
			sp.business_name as name,
			sc.name as type,
			sp.description,
			sp.logo as image,
			sp.average_rating as rating,
			sp.address,
			sp.latitude,
			sp.longitude,
			sp.phone,
			sp.email,
			sp.website,
			sp.is_available as available_now
		`).
		Joins("JOIN service_sub_categories sc ON sp.sub_category_id = sc.id").
		Where(`
			sp.latitude BETWEEN ? AND ? AND 
			sp.longitude BETWEEN ? AND ? AND 
			sp.is_verified = true AND 
			sp.is_available = true
		`, bounds.Southwest.Latitude, bounds.Northeast.Latitude,
			bounds.Southwest.Longitude, bounds.Northeast.Longitude)

	// Add category filter if specified
	if req.Category != "" && req.Category != "all" {
		query = query.Joins("JOIN service_categories c ON sc.category_id = c.id").
			Where("c.name = ?", req.Category)
	}

	// Add pagination
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}

	var results []struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Type         string  `json:"type"`
		Description  string  `json:"description"`
		Image        string  `json:"image"`
		Rating       float64 `json:"rating"`
		Address      string  `json:"address"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Phone        string  `json:"phone"`
		Email        string  `json:"email"`
		Website      string  `json:"website"`
		AvailableNow bool    `json:"available_now"`
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to query nearby services: %w", err)
	}

	// Process results and calculate exact distances
	for _, result := range results {
		distance := s.calculateDistance(req.Location, Location{
			Latitude:  result.Latitude,
			Longitude: result.Longitude,
		})

		// Only include services within the exact radius
		if distance <= req.RadiusKm {
			service := NearbyServiceResponse{
				ID:           result.ID,
				Name:         result.Name,
				Type:         result.Type,
				Description:  result.Description,
				Image:        result.Image,
				Rating:       result.Rating,
				Price:        s.getServicePrice(result.ID),
				Distance:     distance,
				DistanceText: s.formatDistance(distance),
				Location: Location{
					Latitude:  result.Latitude,
					Longitude: result.Longitude,
				},
				Address:      result.Address,
				AvailableNow: result.AvailableNow,
				Contact: Contact{
					Phone:   result.Phone,
					Email:   result.Email,
					Website: result.Website,
				},
				OpenNow: s.isServiceOpenNow(result.ID),
			}

			// Get operating hours
			if hours := s.getOperatingHours(result.ID); hours != nil {
				service.OperatingHours = hours
			}

			// Get features
			service.Features = s.getServiceFeatures(result.ID)

			services = append(services, service)
		}
	}

	return services, nil
}

// GetMapViewServices gets services for map display with clustering
func (s *LocationService) GetMapViewServices(bounds LocationBounds, category string) (*MapViewResponse, error) {
	response := &MapViewResponse{
		Markers:  []MapMarker{},
		Clusters: []MapCluster{},
	}

	// Query services within bounds
	query := s.db.Table("service_providers sp").
		Select(`
			sp.id,
			sp.business_name as name,
			sc.name as type,
			sp.latitude,
			sp.longitude,
			sp.average_rating as rating,
			sp.logo as image
		`).
		Joins("JOIN service_sub_categories sc ON sp.sub_category_id = sc.id").
		Where(`
			sp.latitude BETWEEN ? AND ? AND 
			sp.longitude BETWEEN ? AND ? AND 
			sp.is_verified = true AND 
			sp.is_available = true
		`, bounds.Southwest.Latitude, bounds.Northeast.Latitude,
			bounds.Southwest.Longitude, bounds.Northeast.Longitude)

	// Add category filter if specified
	if category != "" && category != "all" {
		query = query.Joins("JOIN service_categories c ON sc.category_id = c.id").
			Where("c.name = ?", category)
	}

	var results []struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		Type      string  `json:"type"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Rating    float64 `json:"rating"`
		Image     string  `json:"image"`
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to query map services: %w", err)
	}

	// Create markers and apply clustering
	clusterThreshold := 0.01 // ~1km clustering threshold

	for _, result := range results {
		// Check if this service should be clustered
		clustered := false
		for _, cluster := range response.Clusters {
			distance := s.calculateDistance(Location{
				Latitude:  result.Latitude,
				Longitude: result.Longitude,
			}, cluster.Location)

			if distance < clusterThreshold {
				cluster.Count++
				clustered = true
				break
			}
		}

		if !clustered {
			// Check if close to existing markers
			closeToMarker := false
			for _, marker := range response.Markers {
				distance := s.calculateDistance(Location{
					Latitude:  result.Latitude,
					Longitude: result.Longitude,
				}, marker.Location)

				if distance < clusterThreshold {
					closeToMarker = true
					break
				}
			}

			if !closeToMarker {
				marker := MapMarker{
					ID:   result.ID,
					Name: result.Name,
					Type: result.Type,
					Location: Location{
						Latitude:  result.Latitude,
						Longitude: result.Longitude,
					},
					Price:  s.getServicePrice(result.ID),
					Rating: result.Rating,
					Icon:   s.getServiceIcon(result.Type),
				}
				response.Markers = append(response.Markers, marker)
			}
		}
	}

	return response, nil
}

// GetDistanceFilters returns available distance filters with counts
func (s *LocationService) GetDistanceFilters(location Location, category string) ([]DistanceFilter, error) {
	filters := []DistanceFilter{
		{Label: "< 1 km", Value: 1},
		{Label: "< 5 km", Value: 5},
		{Label: "< 10 km", Value: 10},
		{Label: "Any distance", Value: 0},
	}

	// Get service counts for each radius
	for i := range filters {
		if filters[i].Value > 0 {
			req := NearbyServiceRequest{
				Location: location,
				RadiusKm: float64(filters[i].Value),
				Category: category,
				Limit:    1000, // Get count only
			}

			services, err := s.GetNearbyServices(req)
			if err == nil {
				filters[i].Count = len(services)
			}
		} else {
			// Count all services in the area
			count, err := s.countAllServicesInArea(location, category)
			if err == nil {
				filters[i].Count = count
			}
		}
	}

	return filters, nil
}

// DistanceFilter represents a distance filter option
type DistanceFilter struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Count int    `json:"count"`
}

// Helper functions

// calculateDistance calculates the distance between two points using Haversine formula
func (s *LocationService) calculateDistance(loc1, loc2 Location) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert to radians
	lat1 := loc1.Latitude * math.Pi / 180
	lon1 := loc1.Longitude * math.Pi / 180
	lat2 := loc2.Latitude * math.Pi / 180
	lon2 := loc2.Longitude * math.Pi / 180

	// Haversine formula
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// calculateBoundingBox calculates a bounding box for a given location and radius
func (s *LocationService) calculateBoundingBox(location Location, radiusKm float64) LocationBounds {
	// Approximate conversion (1 degree ≈ 111 km)
	deltaLat := radiusKm / 111.0
	deltaLng := radiusKm / (111.0 * math.Cos(location.Latitude*math.Pi/180))

	return LocationBounds{
		Southwest: Location{
			Latitude:  location.Latitude - deltaLat,
			Longitude: location.Longitude - deltaLng,
		},
		Northeast: Location{
			Latitude:  location.Latitude + deltaLat,
			Longitude: location.Longitude + deltaLng,
		},
	}
}

// formatDistance formats distance for display
func (s *LocationService) formatDistance(distance float64) string {
	if distance < 1 {
		return fmt.Sprintf("%.0f m", distance*1000)
	}
	return fmt.Sprintf("%.1f km", distance)
}

// getServicePrice gets the price range for a service
func (s *LocationService) getServicePrice(serviceID string) string {
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

// getServiceIcon gets the icon for a service type
func (s *LocationService) getServiceIcon(serviceType string) string {
	iconMap := map[string]string{
		"hotel":           "building",
		"restaurant":      "utensils",
		"transport":       "car",
		"shopping":        "shopping-bag",
		"airport pickup":  "plane",
		"tour guide":      "map",
		"courier":         "truck",
		"cleaning":        "broom",
		"event planning":  "calendar",
		"graphics":        "palette",
	}
	
	if icon, exists := iconMap[serviceType]; exists {
		return icon
	}
	return "star"
}

// isServiceOpenNow checks if a service is currently open
func (s *LocationService) isServiceOpenNow(serviceID string) bool {
	now := time.Now()
	dayOfWeek := int(now.Weekday())
	currentTime := now.Format("15:04")

	var availability ProviderAvailability
	err := s.db.Where("provider_id = ? AND day_of_week = ?", serviceID, dayOfWeek).
		First(&availability).Error
	
	if err != nil {
		return false // No availability set
	}

	if !availability.IsAvailable {
		return false
	}

	// Check if current time is within operating hours
	return currentTime >= availability.StartTime && currentTime <= availability.EndTime
}

// getOperatingHours gets operating hours for a service
func (s *LocationService) getOperatingHours(serviceID string) *OperatingHours {
	var availabilities []ProviderAvailability
	err := s.db.Where("provider_id = ? AND is_available = ?", serviceID, true).
		Order("day_of_week").
		Find(&availabilities).Error
	
	if err != nil || len(availabilities) == 0 {
		return nil
	}

	// Find common operating hours (simplified)
	hours := &OperatingHours{
		Open:    availabilities[0].StartTime,
		Close:   availabilities[0].EndTime,
		OpenNow: s.isServiceOpenNow(serviceID),
	}

	// Add days
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for _, availability := range availabilities {
		if availability.IsAvailable {
			hours.Days = append(hours.Days, days[availability.DayOfWeek])
		}
	}

	return hours
}

// getServiceFeatures gets features for a service
func (s *LocationService) getServiceFeatures(serviceID string) []string {
	var features []string
	s.db.Table("services").
		Select("name").
		Where("provider_id = ? AND is_active = true", serviceID).
		Limit(5).
		Pluck("name", &features)
	
	return features
}

// countAllServicesInArea counts all services in a larger area
func (s *LocationService) countAllServicesInArea(location Location, category string) (int, error) {
	bounds := s.calculateBoundingBox(location, 50) // 50km radius

	query := s.db.Table("service_providers sp").
		Joins("JOIN service_sub_categories sc ON sp.sub_category_id = sc.id").
		Where(`
			sp.latitude BETWEEN ? AND ? AND 
			sp.longitude BETWEEN ? AND ? AND 
			sp.is_verified = true AND 
			sp.is_available = true
		`, bounds.Southwest.Latitude, bounds.Northeast.Latitude,
			bounds.Southwest.Longitude, bounds.Northeast.Longitude)

	if category != "" && category != "all" {
		query = query.Joins("JOIN service_categories c ON sc.category_id = c.id").
			Where("c.name = ?", category)
	}

	var count int64
	err := query.Count(&count).Error
	return int(count), err
}

// ProviderAvailability represents provider availability
type ProviderAvailability struct {
	ID           uint      `gorm:"primaryKey"`
	ProviderID   uint      `json:"provider_id"`
	DayOfWeek    int       `json:"day_of_week"` // 0-6 (Sunday-Saturday)
	StartTime    string    `json:"start_time"`   // "09:00"
	EndTime      string    `json:"end_time"`     // "17:00"
	IsAvailable  bool      `json:"is_available"`
}
