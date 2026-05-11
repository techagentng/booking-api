package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Import types from models
type (
	APIResponse struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Meta    APIMeta     `json:"meta"`
	}

	APIResponseWithPagination struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Meta    APIMeta     `json:"meta"`
	}

	APIMeta struct {
		Timestamp  string          `json:"timestamp"`
		Location   string          `json:"location,omitempty"`
		RequestID  string          `json:"request_id,omitempty"`
		Pagination *PaginationMeta `json:"pagination,omitempty"`
	}

	PaginationMeta struct {
		Page    int   `json:"page"`
		Limit   int   `json:"limit"`
		Total   int64 `json:"total"`
		Pages   int   `json:"pages"`
		HasNext bool  `json:"has_next"`
		HasPrev bool  `json:"has_prev"`
	}

	DistanceFilter struct {
		Label string `json:"label"`
		Value int    `json:"value"`
		Count int    `json:"count"`
	}

	SearchResponse struct {
		Services    []BaseServiceResult `json:"services"`
		Suggestions []string            `json:"suggestions"`
		Categories  []CategoryResult    `json:"categories"`
	}

	BaseServiceResult struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		Type       string  `json:"type"`
		Distance   string  `json:"distance"`
		Rating     float64 `json:"rating"`
		Price      string  `json:"price"`
		Image      string  `json:"image"`
		MatchScore int     `json:"match_score"`
	}

	CategoryResult struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	SearchSuggestion struct {
		Text     string `json:"text"`
		Type     string `json:"type"`
		Category string `json:"category,omitempty"`
	}

	MapViewResponse struct {
		Markers  []MapMarker  `json:"markers"`
		Clusters []MapCluster `json:"clusters"`
	}

	MapMarker struct {
		ID       string      `json:"id"`
		Name     string      `json:"name"`
		Type     string      `json:"type"`
		Location Coordinates `json:"location"`
		Price    string      `json:"price"`
		Rating   float64     `json:"rating"`
		Icon     string      `json:"icon"`
	}

	MapCluster struct {
		ID       string      `json:"id"`
		Location Coordinates `json:"location"`
		Count    int         `json:"count"`
		Type     string      `json:"type"`
		Icon     string      `json:"icon"`
	}

	Coordinates struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	BaseService struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		Image       string   `json:"image"`
		Rating      float64  `json:"rating"`
		Price       string   `json:"price"`
		Location    Location `json:"location"`
		Contact     Contact  `json:"contact"`
		Features    []string `json:"features"`
	}

	Location struct {
		Address     string      `json:"address"`
		City        string      `json:"city"`
		State       string      `json:"state"`
		Coordinates Coordinates `json:"coordinates"`
	}

	Contact struct {
		Phone   string `json:"phone,omitempty"`
		Email   string `json:"email,omitempty"`
		Website string `json:"website,omitempty"`
	}
)

// Mock data variables
var (
	MockCategories = []Category{
		{ID: "all", Name: "All", Icon: "home", Color: "bg-blue-500", Count: 1250},
		{ID: "hotels", Name: "Hotels", Icon: "building", Color: "bg-purple-500", Count: 450},
		{ID: "transport", Name: "Transport", Icon: "car", Color: "bg-green-500", Count: 320},
		{ID: "food", Name: "Food", Icon: "utensils", Color: "bg-orange-500", Count: 280},
		{ID: "shopping", Name: "Shopping", Icon: "shopping-bag", Color: "bg-pink-500", Count: 200},
	}

	MockFeaturedServices = []FeaturedService{
		{Title: "Hotels", Description: "Find perfect accommodation", Icon: "building", Color: "bg-purple-500", Count: 45, Image: "https://cdn.tripsbook.com/images/categories/hotels.jpg"},
		{Title: "Airport Transfers", Description: "Reliable airport pickup", Icon: "plane", Color: "bg-blue-500", Count: 23, Image: "https://cdn.tripsbook.com/images/categories/airport.jpg"},
		{Title: "Transport", Description: "Ride services & rentals", Icon: "car", Color: "bg-green-500", Count: 67, Image: "https://cdn.tripsbook.com/images/categories/transport.jpg"},
		{Title: "Food & Dining", Description: "Restaurants & catering", Icon: "utensils", Color: "bg-orange-500", Count: 89, Image: "https://cdn.tripsbook.com/images/categories/food.jpg"},
	}

	MockPopularDestinations = []PopularDestination{
		{Name: "Lagos", Country: "Nigeria", Rating: 4.7, Distance: "0 km", Image: "https://cdn.tripsbook.com/images/destinations/lagos.jpg", ServiceCount: 1250},
		{Name: "Abuja", Country: "Nigeria", Rating: 4.6, Distance: "500 km", Image: "https://cdn.tripsbook.com/images/destinations/abuja.jpg", ServiceCount: 850},
		{Name: "Port Harcourt", Country: "Nigeria", Rating: 4.5, Distance: "650 km", Image: "https://cdn.tripsbook.com/images/destinations/portharcourt.jpg", ServiceCount: 420},
	}

	MockTrendingServices = []TrendingService{
		{ID: "svc_111", Name: "Federal Palace Hotel", Type: "Hotel", Description: "Historic luxury hotel with modern amenities", Image: "https://cdn.tripsbook.com/images/hotels/federal.jpg", Rating: 4.6, Price: "$$$$", Location: Location{Coordinates: Coordinates{Lat: 9.0579, Lng: 7.4951}}, Distance: "3.1 km", Trending: true, WeeklyChange: 23.5, TrendingBadge: "🔥 Hot"},
		{ID: "svc_222", Name: "Uber Premium", Type: "Transport", Description: "Premium ride service with professional drivers", Image: "https://cdn.tripsbook.com/images/transport/uber.jpg", Rating: 4.7, Price: "₦800/km", Location: Location{Coordinates: Coordinates{Lat: 6.4474, Lng: 3.3903}}, Distance: "1.5 km", Trending: true, WeeklyChange: 18.2, TrendingBadge: "📈 Rising"},
		{ID: "svc_333", Name: "Terra Kulture", Type: "Restaurant", Description: "Contemporary Nigerian cuisine with cultural experience", Image: "https://cdn.tripsbook.com/images/restaurants/terra.jpg", Rating: 4.5, Price: "$$$", Location: Location{Coordinates: Coordinates{Lat: 6.4474, Lng: 3.3903}}, Distance: "2.8 km", Trending: true, WeeklyChange: 12.7, TrendingBadge: "⭐ Popular"},
	}

	MockTrendingCategories = []TrendingCategory{
		{Name: "Hotels", Change: "+23%", Color: "bg-purple-500", Icon: "building"},
		{Name: "Transport", Change: "+18%", Color: "bg-green-500", Icon: "car"},
		{Name: "Food", Change: "+12%", Color: "bg-orange-500", Icon: "utensils"},
		{Name: "Shopping", Change: "+8%", Color: "bg-pink-500", Icon: "shopping-bag"},
	}

	MockNearbyServices = []NearbyService{
		{
			BaseService: BaseService{
				ID: "svc_456", Name: "Eko Hotels & Suites", Type: "Hotel", Description: "Luxury beachfront hotel with stunning ocean views", Image: "https://cdn.tripsbook.com/images/hotels/eko.jpg", Rating: 4.4, Price: "$$$",
				Location: Location{Address: "1415 Adetokunbo Ademola Street", City: "Lagos", State: "Lagos State", Coordinates: Coordinates{Lat: 6.4474, Lng: 3.3903}},
				Contact:  Contact{Phone: "+234-1-2778000", Email: "eko.suites@ekohotels.com", Website: "https://ekohotels.com"},
				Features: []string{"WiFi", "Pool", "Beach Access", "Spa", "Restaurant"},
			},
			Distance: "1.2 km", OpenNow: true, AvailableNow: true,
		},
		{
			BaseService: BaseService{
				ID: "svc_789", Name: "Bolt Driver - John", Type: "Transport", Description: "Professional driver with 5+ years experience", Image: "https://cdn.tripsbook.com/images/drivers/john.jpg", Rating: 4.8, Price: "₦500/km",
				Location: Location{Address: "Victoria Island", City: "Lagos", State: "Lagos State", Coordinates: Coordinates{Lat: 6.4474, Lng: 3.3903}},
				Contact:  Contact{Phone: "+234-800-000-0001"},
				Features: []string{"Air Conditioning", "GPS", "Experienced Driver"},
			},
			Distance: "0.8 km", OpenNow: true, AvailableNow: true,
		},
	}

	MockUserLocation = UserLocation{
		City: "Lagos", State: "Lagos State", Country: "Nigeria",
		Coordinates: Coordinates{Lat: 6.4474, Lng: 3.3903}, Timezone: "Africa/Lagos",
	}

	MockPopularLocations = []PopularLocation{
		{City: "Lagos", State: "Lagos State", Country: "Nigeria", ServiceCount: 1250, Image: "https://cdn.tripsbook.com/images/cities/lagos.jpg"},
		{City: "Abuja", State: "FCT", Country: "Nigeria", ServiceCount: 850, Image: "https://cdn.tripsbook.com/images/cities/abuja.jpg"},
		{City: "Port Harcourt", State: "Rivers State", Country: "Nigeria", ServiceCount: 420, Image: "https://cdn.tripsbook.com/images/cities/portharcourt.jpg"},
	}
)

// Additional types needed
type (
	Category struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
		Count int    `json:"count"`
	}

	FeaturedService struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		Count       int    `json:"count"`
		Image       string `json:"image"`
	}

	PopularDestination struct {
		Name         string  `json:"name"`
		Country      string  `json:"country"`
		Rating       float64 `json:"rating"`
		Distance     string  `json:"distance"`
		Image        string  `json:"image"`
		ServiceCount int     `json:"service_count"`
	}

	TrendingService struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		Type          string   `json:"type"`
		Description   string   `json:"description"`
		Image         string   `json:"image"`
		Rating        float64  `json:"rating"`
		Price         string   `json:"price"`
		Location      Location `json:"location"`
		Distance      string   `json:"distance"`
		Trending      bool     `json:"trending"`
		WeeklyChange  float64  `json:"weekly_change"`
		TrendingBadge string   `json:"trending_badge"`
	}

	TrendingCategory struct {
		Name   string `json:"name"`
		Change string `json:"change"`
		Color  string `json:"color"`
		Icon   string `json:"icon"`
	}

	NearbyService struct {
		BaseService
		Distance     string `json:"distance"`
		OpenNow      bool   `json:"open_now"`
		AvailableNow bool   `json:"available_now"`
	}

	UserLocation struct {
		City        string      `json:"city"`
		State       string      `json:"state"`
		Country     string      `json:"country"`
		Coordinates Coordinates `json:"coordinates"`
		Timezone    string      `json:"timezone"`
	}

	PopularLocation struct {
		City         string `json:"city"`
		State        string `json:"state"`
		Country      string `json:"country"`
		ServiceCount int    `json:"service_count"`
		Image        string `json:"image"`
	}
)

// MockAPIHandler provides mock data for development
type MockAPIHandler struct{}

// NewMockAPIHandler creates a new mock API handler
func NewMockAPIHandler() *MockAPIHandler {
	return &MockAPIHandler{}
}

// GetCategories returns mock categories
func (h *MockAPIHandler) GetCategories(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    MockCategories,
		Message: "Categories retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
			Location:  "Lagos, Nigeria",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetExploreFeatured returns mock featured services
func (h *MockAPIHandler) GetExploreFeatured(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    MockFeaturedServices,
		Message: "Featured services retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
			Location:  "Lagos, Nigeria",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetExploreDestinations returns mock popular destinations
func (h *MockAPIHandler) GetExploreDestinations(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    MockPopularDestinations,
		Message: "Popular destinations retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetTrendingServices returns mock trending services
func (h *MockAPIHandler) GetTrendingServices(c *gin.Context) {
	_ = c.DefaultQuery("period", "week") // period parameter
	location := c.DefaultQuery("location", "Lagos")

	// Filter by location if needed
	services := MockTrendingServices

	response := APIResponse{
		Success: true,
		Data:    services,
		Message: "Trending services retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
			Location:  location,
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetTrendingCategories returns mock trending categories
func (h *MockAPIHandler) GetTrendingCategories(c *gin.Context) {
	_ = c.DefaultQuery("period", "week") // period parameter

	categories := MockTrendingCategories

	response := APIResponse{
		Success: true,
		Data:    categories,
		Message: "Trending categories retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetNearbyServices returns mock nearby services
func (h *MockAPIHandler) GetNearbyServices(c *gin.Context) {
	_ = c.Query("lat")    // lat parameter
	_ = c.Query("lng")    // lng parameter
	_ = c.Query("radius") // radius parameter
	category := c.Query("category")

	// Default coordinates if not provided - using Lagos coordinates
	// lat, lng = 6.4474, 3.3903

	// Filter by category if specified
	services := MockNearbyServices
	if category != "" && category != "all" {
		// Simple filtering logic
		filtered := []NearbyService{}
		for _, service := range services {
			if service.Type == category {
				filtered = append(filtered, service)
			}
		}
		services = filtered
	}

	response := APIResponse{
		Success: true,
		Data:    services,
		Message: "Nearby services retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
			Location:  "Lagos, Nigeria",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetDistanceFilters returns mock distance filters
func (h *MockAPIHandler) GetDistanceFilters(c *gin.Context) {
	filters := []DistanceFilter{
		{Label: "< 1 km", Value: 1, Count: 45},
		{Label: "< 5 km", Value: 5, Count: 128},
		{Label: "< 10 km", Value: 10, Count: 234},
		{Label: "Any distance", Value: 0, Count: 1250},
	}

	response := APIResponse{
		Success: true,
		Data:    filters,
		Message: "Distance filters retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetCurrentLocation returns mock current location
func (h *MockAPIHandler) GetCurrentLocation(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    MockUserLocation,
		Message: "Current location retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetPopularLocations returns mock popular locations
func (h *MockAPIHandler) GetPopularLocations(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    MockPopularLocations,
		Message: "Popular locations retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// SearchServices returns mock search results
func (h *MockAPIHandler) SearchServices(c *gin.Context) {
	_ = c.Query("q")        // query parameter
	_ = c.Query("location") // location parameter
	_ = c.Query("category") // category parameter

	// Mock search response
	searchResponse := SearchResponse{
		Services: []BaseServiceResult{
			{
				ID:         "svc_123",
				Name:       "Transcorp Hilton",
				Type:       "Hotel",
				Distance:   "2.3 km",
				Rating:     4.5,
				Price:      "$$$$",
				Image:      "https://cdn.tripsbook.com/images/hotels/transcorp.jpg",
				MatchScore: 95,
			},
			{
				ID:         "svc_456",
				Name:       "Eko Hotels & Suites",
				Type:       "Hotel",
				Distance:   "1.2 km",
				Rating:     4.4,
				Price:      "$$$",
				Image:      "https://cdn.tripsbook.com/images/hotels/eko.jpg",
				MatchScore: 88,
			},
		},
		Suggestions: []string{
			"Transcorp Hilton Abuja",
			"Transcorp Hotels Lagos",
			"Hilton Hotels Nigeria",
		},
		Categories: []CategoryResult{
			{
				Name:  "Hotels",
				Count: 45,
			},
		},
	}

	response := APIResponse{
		Success: true,
		Data:    searchResponse,
		Message: "Search completed successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
			Location:  "Lagos",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetSearchSuggestions returns mock search suggestions
func (h *MockAPIHandler) GetSearchSuggestions(c *gin.Context) {
	_ = c.DefaultQuery("period", "week")    // period parameter
	_ = c.DefaultQuery("location", "Lagos") // location parameter

	suggestions := []SearchSuggestion{
		{
			Text:     "Transcorp Hilton Abuja",
			Type:     "service",
			Category: "Hotel",
		},
		{
			Text:     "Transcorp Hotels Lagos",
			Type:     "service",
			Category: "Hotel",
		},
		{
			Text:     "Transport Services",
			Type:     "category",
			Category: "Transport",
		},
	}

	response := APIResponse{
		Success: true,
		Data:    suggestions,
		Message: "Search suggestions retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetServiceDetails returns mock service details
func (h *MockAPIHandler) GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("id")

	// Mock service details
	serviceDetails := map[string]interface{}{
		"id":           serviceID,
		"name":         "Transcorp Hilton Hotel",
		"type":         "hotel",
		"description":  "Luxury hotel in the heart of Abuja with world-class amenities and exceptional service",
		"images":       []string{"https://cdn.tripsbook.com/images/hotels/transcorp1.jpg", "https://cdn.tripsbook.com/images/hotels/transcorp2.jpg"},
		"rating":       4.5,
		"review_count": 342,
		"price":        "$$$$",
		"location": map[string]interface{}{
			"address": "1 Constitution Avenue",
			"city":    "Abuja",
			"state":   "FCT",
			"coordinates": map[string]float64{
				"lat": 9.0579,
				"lng": 7.4951,
			},
		},
		"contact": map[string]string{
			"phone":   "+234-800-000-0000",
			"email":   "abuja@transcorphilton.com",
			"website": "https://transcorphilton.com",
		},
		"operating_hours": map[string]interface{}{
			"open":     "00:00",
			"close":    "23:59",
			"days":     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			"open_now": true,
		},
		"features": []string{"WiFi", "Pool", "Gym", "Spa", "Restaurant", "Bar", "Business Center"},
		"distance": "2.3 km",
		"nearby_services": []map[string]interface{}{
			{
				"id":       "svc_456",
				"name":     "Federal Secretariat",
				"type":     "Government",
				"distance": "0.5 km",
			},
		},
	}

	response := APIResponse{
		Success: true,
		Data:    serviceDetails,
		Message: "Service details retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetMapViewServices returns mock map view data
func (h *MockAPIHandler) GetMapViewServices(c *gin.Context) {
	mapViewResponse := MapViewResponse{
		Markers: []MapMarker{
			{
				ID:   "svc_123",
				Name: "Transcorp Hilton",
				Type: "hotel",
				Location: Coordinates{
					Lat: 9.0579,
					Lng: 7.4951,
				},
				Price:  "$$$$",
				Rating: 4.5,
				Icon:   "hotel",
			},
			{
				ID:   "svc_456",
				Name: "Eko Hotels",
				Type: "hotel",
				Location: Coordinates{
					Lat: 6.4474,
					Lng: 3.3903,
				},
				Price:  "$$$",
				Rating: 4.4,
				Icon:   "hotel",
			},
		},
		Clusters: []MapCluster{
			{
				ID: "cluster_1",
				Location: Coordinates{
					Lat: 6.4480,
					Lng: 3.3905,
				},
				Count: 15,
				Type:  "hotel",
				Icon:  "hotel",
			},
		},
	}

	response := APIResponse{
		Success: true,
		Data:    mapViewResponse,
		Message: "Map services retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// UpdateLocation handles location update (mock)
func (h *MockAPIHandler) UpdateLocation(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    nil,
		Message: "Location updated successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// TrackAnalytics handles analytics tracking (mock)
func (h *MockAPIHandler) TrackAnalytics(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Data:    nil,
		Message: "Analytics tracked successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetServicesByCategory returns mock services by category
func (h *MockAPIHandler) GetServicesByCategory(c *gin.Context) {
	category := c.Param("category")
	location := c.DefaultQuery("location", "Lagos")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Mock services based on category
	var services []BaseService

	switch category {
	case "hotels":
		services = []BaseService{
			{
				ID:          "svc_123",
				Name:        "Transcorp Hilton Hotel",
				Type:        "hotel",
				Description: "Luxury hotel in the heart of Abuja",
				Image:       "https://cdn.tripsbook.com/images/hotels/transcorp.jpg",
				Rating:      4.5,
				Price:       "$$$$",
				Location: Location{
					Address: "1 Constitution Avenue",
					City:    "Abuja",
					State:   "FCT",
					Coordinates: Coordinates{
						Lat: 9.0579,
						Lng: 7.4951,
					},
				},
				Contact: Contact{
					Phone:   "+234-800-000-0000",
					Email:   "abuja@transcorphilton.com",
					Website: "https://transcorphilton.com",
				},
				Features: []string{"WiFi", "Pool", "Gym", "Spa"},
			},
		}
	default:
		services = []BaseService{}
	}

	response := APIResponseWithPagination{
		Success: true,
		Data:    services,
		Message: "Services retrieved successfully",
		Meta: APIMeta{
			Timestamp: "2024-01-15T10:30:00Z",
			Location:  location,
			Pagination: &PaginationMeta{
				Page:    page,
				Limit:   limit,
				Total:   450,
				Pages:   23,
				HasNext: page < 23,
				HasPrev: page > 1,
			},
		},
	}
	c.JSON(http.StatusOK, response)
}
