package tripsbook

// Mock data for development and testing

// Mock Categories
var MockCategories = []Category{
	{
		ID:    "all",
		Name:  "All",
		Icon:  "home",
		Color: "bg-blue-500",
		Count: 1250,
	},
	{
		ID:    "hotels",
		Name:  "Hotels",
		Icon:  "building",
		Color: "bg-purple-500",
		Count: 450,
	},
	{
		ID:    "transport",
		Name:  "Transport",
		Icon:  "car",
		Color: "bg-green-500",
		Count: 320,
	},
	{
		ID:    "food",
		Name:  "Food",
		Icon:  "utensils",
		Color: "bg-orange-500",
		Count: 280,
	},
	{
		ID:    "shopping",
		Name:  "Shopping",
		Icon:  "shopping-bag",
		Color: "bg-pink-500",
		Count: 200,
	},
}

// Mock Featured Services
var MockFeaturedServices = []FeaturedService{
	{
		Title:       "Hotels",
		Description: "Find perfect accommodation",
		Icon:        "building",
		Color:       "bg-purple-500",
		Count:       45,
		Image:       "https://cdn.tripsbook.com/images/categories/hotels.jpg",
	},
	{
		Title:       "Airport Transfers",
		Description: "Reliable airport pickup",
		Icon:        "plane",
		Color:       "bg-blue-500",
		Count:       23,
		Image:       "https://cdn.tripsbook.com/images/categories/airport.jpg",
	},
	{
		Title:       "Transport",
		Description: "Ride services & rentals",
		Icon:        "car",
		Color:       "bg-green-500",
		Count:       67,
		Image:       "https://cdn.tripsbook.com/images/categories/transport.jpg",
	},
	{
		Title:       "Food & Dining",
		Description: "Restaurants & catering",
		Icon:        "utensils",
		Color:       "bg-orange-500",
		Count:       89,
		Image:       "https://cdn.tripsbook.com/images/categories/food.jpg",
	},
}

// Mock Popular Destinations
var MockPopularDestinations = []PopularDestination{
	{
		Name:         "Lagos",
		Country:      "Nigeria",
		Rating:       4.7,
		Distance:     "0 km",
		Image:        "https://cdn.tripsbook.com/images/destinations/lagos.jpg",
		ServiceCount: 1250,
	},
	{
		Name:         "Abuja",
		Country:      "Nigeria",
		Rating:       4.6,
		Distance:     "500 km",
		Image:        "https://cdn.tripsbook.com/images/destinations/abuja.jpg",
		ServiceCount: 850,
	},
	{
		Name:         "Port Harcourt",
		Country:      "Nigeria",
		Rating:       4.5,
		Distance:     "650 km",
		Image:        "https://cdn.tripsbook.com/images/destinations/portharcourt.jpg",
		ServiceCount: 420,
	},
}

// Mock Trending Services
var MockTrendingServices = []TrendingService{
	{
		ID:          "svc_111",
		Name:        "Federal Palace Hotel",
		Type:        "Hotel",
		Description: "Historic luxury hotel with modern amenities",
		Image:       "https://cdn.tripsbook.com/images/hotels/federal.jpg",
		Rating:      4.6,
		Price:       "$$$$",
		Location: Location{
			Coordinates: Coordinates{
				Lat: 9.0579,
				Lng: 7.4951,
			},
		},
		Distance:      "3.1 km",
		Trending:      true,
		WeeklyChange:  23.5,
		TrendingBadge: "🔥 Hot",
	},
	{
		ID:          "svc_222",
		Name:        "Uber Premium",
		Type:        "Transport",
		Description: "Premium ride service with professional drivers",
		Image:       "https://cdn.tripsbook.com/images/transport/uber.jpg",
		Rating:      4.7,
		Price:       "₦800/km",
		Location: Location{
			Coordinates: Coordinates{
				Lat: 6.4474,
				Lng: 3.3903,
			},
		},
		Distance:      "1.5 km",
		Trending:      true,
		WeeklyChange:  18.2,
		TrendingBadge: "📈 Rising",
	},
	{
		ID:          "svc_333",
		Name:        "Terra Kulture",
		Type:        "Restaurant",
		Description: "Contemporary Nigerian cuisine with cultural experience",
		Image:       "https://cdn.tripsbook.com/images/restaurants/terra.jpg",
		Rating:      4.5,
		Price:       "$$$",
		Location: Location{
			Coordinates: Coordinates{
				Lat: 6.4474,
				Lng: 3.3903,
			},
		},
		Distance:      "2.8 km",
		Trending:      true,
		WeeklyChange:  12.7,
		TrendingBadge: "⭐ Popular",
	},
}

// Mock Trending Categories
var MockTrendingCategories = []TrendingCategory{
	{
		Name:   "Hotels",
		Change: "+23%",
		Color:  "bg-purple-500",
		Icon:   "building",
	},
	{
		Name:   "Transport",
		Change: "+18%",
		Color:  "bg-green-500",
		Icon:   "car",
	},
	{
		Name:   "Food",
		Change: "+12%",
		Color:  "bg-orange-500",
		Icon:   "utensils",
	},
	{
		Name:   "Shopping",
		Change: "+8%",
		Color:  "bg-pink-500",
		Icon:   "shopping-bag",
	},
}

// Mock Nearby Services
var MockNearbyServices = []NearbyService{
	{
		BaseService: BaseService{
			ID:          "svc_456",
			Name:        "Eko Hotels & Suites",
			Type:        "Hotel",
			Description: "Luxury beachfront hotel with stunning ocean views",
			Image:       "https://cdn.tripsbook.com/images/hotels/eko.jpg",
			Rating:      4.4,
			Price:       "$$$",
			Location: Location{
				Address: "1415 Adetokunbo Ademola Street",
				City:    "Lagos",
				State:   "Lagos State",
				Coordinates: Coordinates{
					Lat: 6.4474,
					Lng: 3.3903,
				},
			},
			Contact: Contact{
				Phone:   "+234-1-2778000",
				Email:   "eko.suites@ekohotels.com",
				Website: "https://ekohotels.com",
			},
			Features: []string{"WiFi", "Pool", "Beach Access", "Spa", "Restaurant"},
		},
		Distance:     "1.2 km",
		OpenNow:      true,
		AvailableNow: true,
	},
	{
		BaseService: BaseService{
			ID:          "svc_789",
			Name:        "Bolt Driver - John",
			Type:        "Transport",
			Description: "Professional driver with 5+ years experience",
			Image:       "https://cdn.tripsbook.com/images/drivers/john.jpg",
			Rating:      4.8,
			Price:       "₦500/km",
			Location: Location{
				Address: "Victoria Island",
				City:    "Lagos",
				State:   "Lagos State",
				Coordinates: Coordinates{
					Lat: 6.4474,
					Lng: 3.3903,
				},
			},
			Contact: Contact{
				Phone: "+234-800-000-0001",
			},
			Features: []string{"Air Conditioning", "GPS", "Experienced Driver"},
		},
		Distance:     "0.8 km",
		OpenNow:      true,
		AvailableNow: true,
	},
}

// Mock User Location
var MockUserLocation = UserLocation{
	City:    "Lagos",
	State:   "Lagos State",
	Country: "Nigeria",
	Coordinates: Coordinates{
		Lat: 6.4474,
		Lng: 3.3903,
	},
	Timezone: "Africa/Lagos",
}

// Mock Popular Locations
var MockPopularLocations = []PopularLocation{
	{
		City:         "Lagos",
		State:        "Lagos State",
		Country:      "Nigeria",
		ServiceCount: 1250,
		Image:        "https://cdn.tripsbook.com/images/cities/lagos.jpg",
	},
	{
		City:         "Abuja",
		State:        "FCT",
		Country:      "Nigeria",
		ServiceCount: 850,
		Image:        "https://cdn.tripsbook.com/images/cities/abuja.jpg",
	},
	{
		City:         "Port Harcourt",
		State:        "Rivers State",
		Country:      "Nigeria",
		ServiceCount: 420,
		Image:        "https://cdn.tripsbook.com/images/cities/portharcourt.jpg",
	},
}

// Type definitions for mock data
type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Count int    `json:"count"`
}

type FeaturedService struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Count       int    `json:"count"`
	Image       string `json:"image"`
}

type PopularDestination struct {
	Name         string  `json:"name"`
	Country      string  `json:"country"`
	Rating       float64 `json:"rating"`
	Distance     string  `json:"distance"`
	Image        string  `json:"image"`
	ServiceCount int     `json:"service_count"`
}

type TrendingService struct {
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

type TrendingCategory struct {
	Name   string `json:"name"`
	Change string `json:"change"`
	Color  string `json:"color"`
	Icon   string `json:"icon"`
}

type NearbyService struct {
	BaseService
	Distance     string `json:"distance"`
	OpenNow      bool   `json:"open_now"`
	AvailableNow bool   `json:"available_now"`
}

type BaseService struct {
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

type Location struct {
	Address     string      `json:"address"`
	City        string      `json:"city"`
	State       string      `json:"state"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Contact struct {
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`
}

type UserLocation struct {
	City        string      `json:"city"`
	State       string      `json:"state"`
	Country     string      `json:"country"`
	Coordinates Coordinates `json:"coordinates"`
	Timezone    string      `json:"timezone"`
}

type PopularLocation struct {
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	ServiceCount int    `json:"service_count"`
	Image        string `json:"image"`
}
