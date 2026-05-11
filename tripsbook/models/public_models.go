package models

import (
	"fmt"
	"math"
	"time"
)

// Public API Models for Mobile Frontend

// BaseService represents the base structure for all services
type BaseService struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Rating      float64   `json:"rating"`
	Price       string    `json:"price"`
	Location    Location  `json:"location"`
	Contact     Contact   `json:"contact"`
	Features    []string  `json:"features,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Location represents geographical coordinates
type Location struct {
	Address     string      `json:"address"`
	City        string      `json:"city"`
	State       string      `json:"state"`
	Coordinates Coordinates `json:"coordinates"`
}

// Coordinates represents lat/lng coordinates
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Contact represents contact information
type Contact struct {
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`
}

// OperatingHours represents service operating hours
type OperatingHours struct {
	Open    string   `json:"open"`
	Close   string   `json:"close"`
	Days    []string `json:"days"`
	OpenNow bool     `json:"open_now"`
}

// Hotel represents hotel-specific service data
type Hotel struct {
	BaseService
	Type           string          `json:"type"` // "hotel"
	Rooms          []Room          `json:"rooms"`
	Amenities      []string        `json:"amenities"`
	StarRating     int             `json:"star_rating"`
	CheckIn        string          `json:"check_in"`
	CheckOut       string          `json:"check_out"`
	OperatingHours *OperatingHours `json:"operating_hours,omitempty"`
}

// Room represents hotel room information
type Room struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Price     float64  `json:"price"`
	Capacity  int      `json:"capacity"`
	Available bool     `json:"available"`
	Images    []string `json:"images"`
}

// Restaurant represents restaurant-specific service data
type Restaurant struct {
	BaseService
	Type              string          `json:"type"` // "restaurant"
	Cuisine           []string        `json:"cuisine"`
	PriceRange        string          `json:"price_range"` // "$" | "$$" | "$$$" | "$$$$"
	MenuItems         []MenuItem      `json:"menu_items,omitempty"`
	DeliveryAvailable bool            `json:"delivery_available"`
	DineInAvailable   bool            `json:"dine_in_available"`
	OperatingHours    *OperatingHours `json:"operating_hours,omitempty"`
}

// MenuItem represents restaurant menu item
type MenuItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

// Transport represents transport-specific service data
type Transport struct {
	BaseService
	Type           string          `json:"type"` // "transport"
	VehicleTypes   []string        `json:"vehicle_types"`
	PricingModel   string          `json:"pricing_model"` // "per_km" | "per_hour" | "fixed"`
	BasePrice      float64         `json:"base_price"`
	Availability   bool            `json:"availability"`
	DriverInfo     *DriverInfo     `json:"driver_info,omitempty"`
	OperatingHours *OperatingHours `json:"operating_hours,omitempty"`
}

// DriverInfo represents driver information for transport services
type DriverInfo struct {
	Name    string  `json:"name"`
	Rating  float64 `json:"rating"`
	Trips   int     `json:"trips"`
	Photo   string  `json:"photo"`
	Vehicle string  `json:"vehicle"`
}

// Shopping represents shopping-specific service data
type Shopping struct {
	BaseService
	Type              string          `json:"type"` // "shopping"
	Categories        []string        `json:"categories"`
	StoreType         string          `json:"store_type"` // "mall", "supermarket", "boutique", "market"
	OpeningHours      *OperatingHours `json:"operating_hours,omitempty"`
	DeliveryAvailable bool            `json:"delivery_available"`
	OnlineShopping    bool            `json:"online_shopping"`
}

// Category represents service category
type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Count int    `json:"count"`
}

// FeaturedService represents featured service for explore tab
type FeaturedService struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Count       int    `json:"count"`
	Image       string `json:"image"`
}

// PopularDestination represents popular destination
type PopularDestination struct {
	Name         string  `json:"name"`
	Country      string  `json:"country"`
	Rating       float64 `json:"rating"`
	Distance     string  `json:"distance"`
	Image        string  `json:"image"`
	ServiceCount int     `json:"service_count"`
}

// NearbyService represents nearby service with distance
type NearbyService struct {
	BaseService
	Distance     string `json:"distance"`
	OpenNow      bool   `json:"open_now"`
	AvailableNow bool   `json:"available_now"`
}

// TrendingService represents trending service
type TrendingService struct {
	BaseService
	Distance      string  `json:"distance"`
	Trending      bool    `json:"trending"`
	Image         string  `json:"image"`
	WeeklyChange  float64 `json:"weekly_change"`
	TrendingBadge string  `json:"trending_badge"`
}

// TrendingCategory represents trending category
type TrendingCategory struct {
	Name   string `json:"name"`
	Change string `json:"change"`
	Color  string `json:"color"`
	Icon   string `json:"icon"`
}

// SearchResponse represents search API response
type SearchResponse struct {
	Services    []BaseService   `json:"services"`
	Suggestions []string        `json:"suggestions"`
	Categories  []CategoryCount `json:"categories"`
}

// CategoryCount represents category with count
type CategoryCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// SearchSuggestion represents search suggestion
type SearchSuggestion struct {
	Text     string `json:"text"`
	Type     string `json:"type"` // "service" | "category"
	Category string `json:"category,omitempty"`
}

// UserLocation represents user location
type UserLocation struct {
	City        string      `json:"city"`
	State       string      `json:"state"`
	Country     string      `json:"country"`
	Coordinates Coordinates `json:"coordinates"`
	Timezone    string      `json:"timezone"`
}

// PopularLocation represents popular location
type PopularLocation struct {
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	ServiceCount int    `json:"service_count"`
	Image        string `json:"image"`
}

// MapMarker represents map marker
type MapMarker struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Location Coordinates `json:"location"`
	Price    string      `json:"price"`
	Rating   float64     `json:"rating"`
	Icon     string      `json:"icon"`
}

// MapCluster represents map cluster
type MapCluster struct {
	ID       string      `json:"id"`
	Location Coordinates `json:"location"`
	Count    int         `json:"count"`
	Type     string      `json:"type"`
	Icon     string      `json:"icon"`
}

// MapViewResponse represents map view response
type MapViewResponse struct {
	Markers  []MapMarker  `json:"markers"`
	Clusters []MapCluster `json:"clusters"`
}

// QuickBookingRequest represents quick booking request
type QuickBookingRequest struct {
	ServiceID      string         `json:"service_id"`
	CustomerInfo   CustomerInfo   `json:"customer_info"`
	BookingDetails BookingDetails `json:"booking_details"`
	Location       Location       `json:"location"`
}

// CustomerInfo represents customer information for quick booking
type CustomerInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// BookingDetails represents booking details
type BookingDetails struct {
	Date     string `json:"date"`
	Time     string `json:"time"`
	Duration int    `json:"duration"` // in minutes
	Notes    string `json:"notes"`
}

// QuickBookingResponse represents quick booking response
type QuickBookingResponse struct {
	BookingID      string   `json:"booking_id"`
	Status         string   `json:"status"`
	EstimatedPrice float64  `json:"estimated_price"`
	NextSteps      []string `json:"next_steps"`
}

// BookingStatus represents booking status
type BookingStatus struct {
	BookingID       string     `json:"booking_id"`
	Status          string     `json:"status"`
	ProviderName    string     `json:"provider_name"`
	ConfirmedAt     *time.Time `json:"confirmed_at,omitempty"`
	TotalAmount     float64    `json:"total_amount"`
	PaymentRequired bool       `json:"payment_required"`
	PaymentLink     string     `json:"payment_link,omitempty"`
}

// DistanceFilter represents distance filter option
type DistanceFilter struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Count int    `json:"count"`
}

// AnalyticsEvent represents analytics tracking event
type AnalyticsEvent struct {
	Event           string                 `json:"event"`
	Data            map[string]interface{} `json:"data"`
	UserFingerprint string                 `json:"user_fingerprint"`
	Timestamp       time.Time              `json:"timestamp"`
}

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Meta    APIMeta     `json:"meta"`
}

// APIMetadata represents API metadata
type APIMeta struct {
	Timestamp  time.Time       `json:"timestamp"`
	Location   string          `json:"location,omitempty"`
	RequestID  string          `json:"request_id,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	Total   int64 `json:"total"`
	Pages   int   `json:"pages"`
	HasNext bool  `json:"has_next"`
	HasPrev bool  `json:"has_prev"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
	Meta    APIMeta  `json:"meta"`
}

// APIError represents API error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Constants for error codes
const (
	ErrInvalidLocation   = "INVALID_LOCATION"
	ErrServiceNotFound   = "SERVICE_NOT_FOUND"
	ErrSearchTooBroad    = "SEARCH_TOO_BROAD"
	ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrInvalidParameters = "INVALID_PARAMETERS"
	ErrInternalServer    = "INTERNAL_SERVER_ERROR"
)

// Constants for service types
const (
	ServiceTypeHotel      = "hotel"
	ServiceTypeRestaurant = "restaurant"
	ServiceTypeTransport  = "transport"
	ServiceTypeShopping   = "shopping"
	ServiceTypeEvent      = "event"
)

// Constants for booking status
const (
	BookingStatusPending    = "pending"
	BookingStatusConfirmed  = "confirmed"
	BookingStatusInProgress = "in_progress"
	BookingStatusCompleted  = "completed"
	BookingStatusCancelled  = "cancelled"
	BookingStatusNoShow     = "no_show"
)

// Constants for payment status
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
	PaymentStatusFailed   = "failed"
)

// Helper functions

// NewAPIResponse creates a new API response
func NewAPIResponse(data interface{}, message string) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta: APIMeta{
			Timestamp: time.Now(),
		},
	}
}

// NewAPIResponseWithLocation creates a new API response with location
func NewAPIResponseWithLocation(data interface{}, message, location string) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta: APIMeta{
			Timestamp: time.Now(),
			Location:  location,
		},
	}
}

// NewAPIResponseWithPagination creates a new API response with pagination
func NewAPIResponseWithPagination(data interface{}, message string, page, limit int, total int64) *APIResponse {
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}

	return &APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta: APIMeta{
			Timestamp: time.Now(),
			Pagination: &PaginationMeta{
				Page:    page,
				Limit:   limit,
				Total:   total,
				Pages:   pages,
				HasNext: page < pages,
				HasPrev: page > 1,
			},
		},
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message, details string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error: APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: APIMeta{
			Timestamp: time.Now(),
		},
	}
}

// IsValidServiceType checks if service type is valid
func IsValidServiceType(serviceType string) bool {
	validTypes := []string{
		ServiceTypeHotel,
		ServiceTypeRestaurant,
		ServiceTypeTransport,
		ServiceTypeShopping,
		ServiceTypeEvent,
	}

	for _, validType := range validTypes {
		if serviceType == validType {
			return true
		}
	}
	return false
}

// IsValidBookingStatus checks if booking status is valid
func IsValidBookingStatus(status string) bool {
	validStatuses := []string{
		BookingStatusPending,
		BookingStatusConfirmed,
		BookingStatusInProgress,
		BookingStatusCompleted,
		BookingStatusCancelled,
		BookingStatusNoShow,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// FormatPrice formats price for display
func FormatPrice(price float64) string {
	if price < 1000 {
		return "₦" + fmt.Sprintf("%.0f", price)
	} else if price < 100000 {
		return "₦" + fmt.Sprintf("%.1fK", price/1000)
	} else {
		return "₦" + fmt.Sprintf("%.1fM", price/1000000)
	}
}

// GetPriceRange returns price range symbol
func GetPriceRange(price float64) string {
	if price < 1000 {
		return "$"
	} else if price < 5000 {
		return "$$"
	} else if price < 20000 {
		return "$$$"
	} else {
		return "$$$$"
	}
}

// CalculateDistance calculates distance between two coordinates
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert to radians
	lat1Rad := lat1 * 3.14159 / 180
	lng1Rad := lng1 * 3.14159 / 180
	lat2Rad := lat2 * 3.14159 / 180
	lng2Rad := lng2 * 3.14159 / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlng := lng2Rad - lng1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// FormatDistance formats distance for display
func FormatDistance(distance float64) string {
	if distance < 1 {
		return fmt.Sprintf("%.0f m", distance*1000)
	}
	return fmt.Sprintf("%.1f km", distance)
}
