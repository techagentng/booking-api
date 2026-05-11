package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hotel/tripsbook/models"
	"hotel/tripsbook/services"
)

// PublicAPIHandler handles public API endpoints
type PublicAPIHandler struct {
	db              *gorm.DB
	locationService *services.LocationService
	trendingService *services.TrendingService
	searchService   *services.SearchService
}

// NewPublicAPIHandler creates a new public API handler
func NewPublicAPIHandler(db *gorm.DB) *PublicAPIHandler {
	return &PublicAPIHandler{
		db:              db,
		locationService: services.NewLocationService(db),
		trendingService: services.NewTrendingService(db),
		searchService:   services.NewSearchService(db),
	}
}

// GetCategories handles GET /api/v1/public/categories
func (h *PublicAPIHandler) GetCategories(c *gin.Context) {
	var categories []models.Category

	// Query categories with service counts
	query := h.db.Table("service_categories c").
		Select(`
			c.id,
			c.name,
			c.icon,
			c.color,
			COALESCE(COUNT(DISTINCT sp.id), 0) as count
		`).
		Joins("LEFT JOIN service_sub_categories sc ON c.id = sc.category_id").
		Joins("LEFT JOIN service_providers sp ON sc.id = sp.sub_category_id").
		Where("sp.is_verified = true AND sp.is_available = true").
		Group("c.id, c.name, c.icon, c.color").
		Order("c.sort_order ASC")

	if err := query.Scan(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch categories",
			err.Error(),
		))
		return
	}

	// Add "All" category at the beginning
	allCategory := models.Category{
		ID:    "all",
		Name:  "All",
		Icon:  "home",
		Color: "bg-blue-500",
		Count: h.getTotalServiceCount(),
	}

	categories = append([]models.Category{allCategory}, categories...)

	response := models.NewAPIResponseWithLocation(categories, "Categories retrieved successfully", h.getCurrentLocation(c))
	c.JSON(http.StatusOK, response)
}

// GetServicesByCategory handles GET /api/v1/public/services/category/{category}
func (h *PublicAPIHandler) GetServicesByCategory(c *gin.Context) {
	category := c.Param("category")
	location := c.Query("location")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validate inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get services
	services, total, err := h.getServicesByCategory(category, location, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch services",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponseWithPagination(services, "Services retrieved successfully", page, limit, total)
	response.Meta.Location = h.getCurrentLocation(c)
	c.JSON(http.StatusOK, response)
}

// GetExploreFeatured handles GET /api/v1/public/explore/featured
func (h *PublicAPIHandler) GetExploreFeatured(c *gin.Context) {
	location := c.Query("location")

	featured := []models.FeaturedService{
		{
			Title:       "Hotels",
			Description: "Find perfect accommodation",
			Icon:        "building",
			Color:       "bg-purple-500",
			Count:       h.getCategoryServiceCount("hotels"),
			Image:       "https://cdn.tripsbook.com/images/categories/hotels.jpg",
		},
		{
			Title:       "Airport Transfers",
			Description: "Reliable airport pickup",
			Icon:        "plane",
			Color:       "bg-blue-500",
			Count:       h.getCategoryServiceCount("airport pickup"),
			Image:       "https://cdn.tripsbook.com/images/categories/airport.jpg",
		},
		{
			Title:       "Transport",
			Description: "Ride services & rentals",
			Icon:        "car",
			Color:       "bg-green-500",
			Count:       h.getCategoryServiceCount("transport"),
			Image:       "https://cdn.tripsbook.com/images/categories/transport.jpg",
		},
		{
			Title:       "Food & Dining",
			Description: "Restaurants & catering",
			Icon:        "utensils",
			Color:       "bg-orange-500",
			Count:       h.getCategoryServiceCount("restaurant"),
			Image:       "https://cdn.tripsbook.com/images/categories/food.jpg",
		},
	}

	response := models.NewAPIResponseWithLocation(featured, "Featured services retrieved successfully", location)
	c.JSON(http.StatusOK, response)
}

// GetExploreDestinations handles GET /api/v1/public/explore/destinations
func (h *PublicAPIHandler) GetExploreDestinations(c *gin.Context) {
	destinations := []models.PopularDestination{
		{
			Name:         "Lagos",
			Country:      "Nigeria",
			Rating:       4.7,
			Distance:     "0 km",
			Image:        "https://cdn.tripsbook.com/images/destinations/lagos.jpg",
			ServiceCount: h.getLocationServiceCount("Lagos"),
		},
		{
			Name:         "Abuja",
			Country:      "Nigeria",
			Rating:       4.6,
			Distance:     "500 km",
			Image:        "https://cdn.tripsbook.com/images/destinations/abuja.jpg",
			ServiceCount: h.getLocationServiceCount("Abuja"),
		},
		{
			Name:         "Port Harcourt",
			Country:      "Nigeria",
			Rating:       4.5,
			Distance:     "650 km",
			Image:        "https://cdn.tripsbook.com/images/destinations/portharcourt.jpg",
			ServiceCount: h.getLocationServiceCount("Port Harcourt"),
		},
	}

	response := models.NewAPIResponse(destinations, "Popular destinations retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetNearbyServices handles GET /api/v1/public/services/nearby
func (h *PublicAPIHandler) GetNearbyServices(c *gin.Context) {
	// Parse parameters
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radius, _ := strconv.ParseFloat(c.Query("radius"), 64)
	category := c.Query("category")
	location := c.Query("location")

	// Validate coordinates
	if lat == 0 || lng == 0 {
		if location != "" {
			// Get coordinates for location
			coords, err := h.getLocationCoordinates(location)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.NewErrorResponse(
					models.ErrInvalidLocation,
					"Invalid location",
					err.Error(),
				))
				return
			}
			lat, lng = coords.Lat, coords.Lng
		} else {
			// Default to Lagos coordinates
			lat, lng = 6.4474, 3.3903
		}
	}

	// Set default radius
	if radius == 0 {
		radius = 5.0 // 5km
	}

	// Create request
	req := services.NearbyServiceRequest{
		Location: services.Location{
			Latitude:  lat,
			Longitude: lng,
		},
		RadiusKm: radius,
		Category: category,
		Limit:    50,
	}

	// Get nearby services
	services, err := h.locationService.GetNearbyServices(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch nearby services",
			err.Error(),
		))
		return
	}

	// Convert to public models
	var nearbyServices []models.NearbyService
	for _, service := range services {
		nearbyService := models.NearbyService{
			BaseService: models.BaseService{
				ID:          service.ID,
				Name:        service.Name,
				Type:        service.Type,
				Description: service.Description,
				Image:       service.Image,
				Rating:      service.Rating,
				Price:       service.Price,
				Location: models.Location{
					Address: service.Address,
					City:    h.getCityFromCoordinates(service.Location),
					State:   h.getStateFromCoordinates(service.Location),
					Coordinates: models.Coordinates{
						Lat: service.Location.Latitude,
						Lng: service.Location.Longitude,
					},
				},
				Contact: models.Contact{
					Phone:   service.Contact.Phone,
					Email:   service.Contact.Email,
					Website: service.Contact.Website,
				},
				Features: service.Features,
			},
			Distance:     service.DistanceText,
			OpenNow:      service.OpenNow,
			AvailableNow: service.AvailableNow,
		}

		nearbyServices = append(nearbyServices, nearbyService)
	}

	response := models.NewAPIResponseWithLocation(nearbyServices, "Nearby services retrieved successfully", location)
	c.JSON(http.StatusOK, response)
}

// GetDistanceFilters handles GET /api/v1/public/services/nearby/filters
func (h *PublicAPIHandler) GetDistanceFilters(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	category := c.Query("category")

	// Default to Lagos if no coordinates provided
	if lat == 0 || lng == 0 {
		lat, lng = 6.4474, 3.3903
	}

	location := services.Location{Latitude: lat, Longitude: lng}
	filters, err := h.locationService.GetDistanceFilters(location, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch distance filters",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponse(filters, "Distance filters retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetTrendingServices handles GET /api/v1/public/services/trending
func (h *PublicAPIHandler) GetTrendingServices(c *gin.Context) {
	period := c.DefaultQuery("period", "week")
	location := c.Query("location")

	req := services.TrendingRequest{
		Location: location,
		Period:   period,
		Limit:    20,
	}

	trending, err := h.trendingService.GetTrendingServices(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch trending services",
			err.Error(),
		))
		return
	}

	// Convert to public models
	var trendingServices []models.TrendingService
	for _, service := range trending {
		trendingService := models.TrendingService{
			BaseService: models.BaseService{
				ID:     service.ID,
				Name:   service.Name,
				Type:   service.Type,
				Image:  service.Image,
				Rating: service.Rating,
				Price:  service.Price,
				Location: models.Location{
					Coordinates: models.Coordinates{
						Lat: service.Location.Latitude,
						Lng: service.Location.Longitude,
					},
				},
			},
			Distance:      service.Distance,
			Trending:      service.Trending,
			Image:         service.Image,
			WeeklyChange:  service.WeeklyChange,
			TrendingBadge: service.TrendingBadge,
		}

		trendingServices = append(trendingServices, trendingService)
	}

	response := models.NewAPIResponseWithLocation(trendingServices, "Trending services retrieved successfully", location)
	c.JSON(http.StatusOK, response)
}

// GetTrendingCategories handles GET /api/v1/public/categories/trending
func (h *PublicAPIHandler) GetTrendingCategories(c *gin.Context) {
	period := c.DefaultQuery("period", "week")

	req := services.TrendingRequest{
		Period: period,
	}

	categories, err := h.trendingService.GetTrendingCategories(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch trending categories",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponse(categories, "Trending categories retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// SearchServices handles GET /api/v1/public/search
func (h *PublicAPIHandler) SearchServices(c *gin.Context) {
	query := c.Query("q")
	location := c.Query("location")
	category := c.Query("category")
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)

	if query == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.ErrInvalidParameters,
			"Search query is required",
			"",
		))
		return
	}

	// Create search request
	req := services.SearchRequest{
		Query:    query,
		Location: location,
		Category: category,
		Lat:      lat,
		Lng:      lng,
		Limit:    20,
	}

	// Perform search
	results, err := h.searchService.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Search failed",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponseWithLocation(results, "Search completed successfully", location)
	c.JSON(http.StatusOK, response)
}

// GetSearchSuggestions handles GET /api/v1/public/search/suggestions
func (h *PublicAPIHandler) GetSearchSuggestions(c *gin.Context) {
	query := c.Query("q")
	location := c.Query("location")

	if query == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.ErrInvalidParameters,
			"Search query is required",
			"",
		))
		return
	}

	suggestions, err := h.searchService.GetSuggestions(query, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to get suggestions",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponse(suggestions, "Search suggestions retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetCurrentLocation handles GET /api/v1/public/location/current
func (h *PublicAPIHandler) GetCurrentLocation(c *gin.Context) {
	location := models.UserLocation{
		City:    "Lagos",
		State:   "Lagos State",
		Country: "Nigeria",
		Coordinates: models.Coordinates{
			Lat: 6.4474,
			Lng: 3.3903,
		},
		Timezone: "Africa/Lagos",
	}

	response := models.NewAPIResponse(location, "Current location retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// UpdateLocation handles POST /api/v1/public/location/update
func (h *PublicAPIHandler) UpdateLocation(c *gin.Context) {
	var req struct {
		City        string             `json:"city"`
		Coordinates models.Coordinates `json:"coordinates"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.ErrInvalidParameters,
			"Invalid location data",
			err.Error(),
		))
		return
	}

	// In a real implementation, you would store this in session or user profile
	// For now, just return success

	response := models.NewAPIResponse(nil, "Location updated successfully")
	c.JSON(http.StatusOK, response)
}

// GetPopularLocations handles GET /api/v1/public/locations/popular
func (h *PublicAPIHandler) GetPopularLocations(c *gin.Context) {
	locations := []models.PopularLocation{
		{
			City:         "Lagos",
			State:        "Lagos State",
			Country:      "Nigeria",
			ServiceCount: h.getLocationServiceCount("Lagos"),
			Image:        "https://cdn.tripsbook.com/images/cities/lagos.jpg",
		},
		{
			City:         "Abuja",
			State:        "FCT",
			Country:      "Nigeria",
			ServiceCount: h.getLocationServiceCount("Abuja"),
			Image:        "https://cdn.tripsbook.com/images/cities/abuja.jpg",
		},
		{
			City:         "Port Harcourt",
			State:        "Rivers State",
			Country:      "Nigeria",
			ServiceCount: h.getLocationServiceCount("Port Harcourt"),
			Image:        "https://cdn.tripsbook.com/images/cities/portharcourt.jpg",
		},
	}

	response := models.NewAPIResponse(locations, "Popular locations retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetServiceDetails handles GET /api/v1/public/services/{id}
func (h *PublicAPIHandler) GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("id")

	// Get service details
	service, err := h.getServiceDetails(serviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				models.ErrServiceNotFound,
				"Service not found",
				"",
			))
			return
		}

		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch service details",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponse(service, "Service details retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetMapViewServices handles GET /api/v1/public/services/map
func (h *PublicAPIHandler) GetMapViewServices(c *gin.Context) {
	swLat, _ := strconv.ParseFloat(c.Query("sw_lat"), 64)
	swLng, _ := strconv.ParseFloat(c.Query("sw_lng"), 64)
	neLat, _ := strconv.ParseFloat(c.Query("ne_lat"), 64)
	neLng, _ := strconv.ParseFloat(c.Query("ne_lng"), 64)
	category := c.Query("category")

	bounds := services.LocationBounds{
		Southwest: services.Location{Latitude: swLat, Longitude: swLng},
		Northeast: services.Location{Latitude: neLat, Longitude: neLng},
	}

	result, err := h.locationService.GetMapViewServices(bounds, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to fetch map services",
			err.Error(),
		))
		return
	}

	// Convert to public models
	var markers []models.MapMarker
	for _, marker := range result.Markers {
		publicMarker := models.MapMarker{
			ID:   marker.ID,
			Name: marker.Name,
			Type: marker.Type,
			Location: models.Coordinates{
				Lat: marker.Location.Latitude,
				Lng: marker.Location.Longitude,
			},
			Price:  marker.Price,
			Rating: marker.Rating,
			Icon:   marker.Icon,
		}
		markers = append(markers, publicMarker)
	}

	var clusters []models.MapCluster
	for _, cluster := range result.Clusters {
		publicCluster := models.MapCluster{
			ID: cluster.ID,
			Location: models.Coordinates{
				Lat: cluster.Location.Latitude,
				Lng: cluster.Location.Longitude,
			},
			Count: cluster.Count,
			Type:  cluster.Type,
			Icon:  cluster.Icon,
		}
		clusters = append(clusters, publicCluster)
	}

	mapViewResponse := models.MapViewResponse{
		Markers:  markers,
		Clusters: clusters,
	}

	response := models.NewAPIResponse(mapViewResponse, "Map services retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// TrackAnalytics handles POST /api/v1/public/analytics/track
func (h *PublicAPIHandler) TrackAnalytics(c *gin.Context) {
	var req models.AnalyticsEvent

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.ErrInvalidParameters,
			"Invalid analytics data",
			err.Error(),
		))
		return
	}

	// Track analytics event
	err := h.trackAnalyticsEvent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			models.ErrInternalServer,
			"Failed to track analytics",
			err.Error(),
		))
		return
	}

	response := models.NewAPIResponse(nil, "Analytics tracked successfully")
	c.JSON(http.StatusOK, response)
}

// Helper methods

func (h *PublicAPIHandler) getCurrentLocation(c *gin.Context) string {
	if location := c.Query("location"); location != "" {
		return location
	}
	return "Lagos, Nigeria"
}

func (h *PublicAPIHandler) getTotalServiceCount() int {
	var count int64
	h.db.Model(&struct{ ID uint }{}).
		Table("service_providers").
		Where("is_verified = true AND is_available = true").
		Count(&count)
	return int(count)
}

func (h *PublicAPIHandler) getCategoryServiceCount(category string) int {
	var count int64
	h.db.Model(&struct{ ID uint }{}).
		Table("service_providers sp").
		Joins("JOIN service_sub_categories sc ON sp.sub_category_id = sc.id").
		Joins("JOIN service_categories c ON sc.category_id = c.id").
		Where("c.name = ? AND sp.is_verified = true AND sp.is_available = true", category).
		Count(&count)
	return int(count)
}

func (h *PublicAPIHandler) getLocationServiceCount(location string) int {
	var count int64
	h.db.Model(&struct{ ID uint }{}).
		Table("service_providers").
		Where("city = ? AND is_verified = true AND is_available = true", location).
		Count(&count)
	return int(count)
}

func (h *PublicAPIHandler) getServicesByCategory(category, location string, page, limit int) ([]models.BaseService, int64, error) {
	// Implementation for getting services by category
	// This would involve complex queries and pagination
	return []models.BaseService{}, 0, nil
}

func (h *PublicAPIHandler) getLocationCoordinates(location string) (models.Coordinates, error) {
	// Implementation for geocoding location to coordinates
	// This would use a geocoding service
	return models.Coordinates{Lat: 6.4474, Lng: 3.3903}, nil
}

func (h *PublicAPIHandler) getCityFromCoordinates(location services.Location) string {
	// Implementation for reverse geocoding
	return "Lagos"
}

func (h *PublicAPIHandler) getStateFromCoordinates(location services.Location) string {
	// Implementation for reverse geocoding
	return "Lagos State"
}

func (h *PublicAPIHandler) getServiceDetails(serviceID string) (interface{}, error) {
	// Implementation for getting detailed service information
	return nil, nil
}

func (h *PublicAPIHandler) trackAnalyticsEvent(event models.AnalyticsEvent) error {
	// Implementation for tracking analytics events
	return nil
}
