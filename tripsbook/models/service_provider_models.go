package models

import (
	"time"
)

// ServiceProvider represents a service provider in the system
type ServiceProvider struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	BusinessName    string    `json:"business_name"`
	DisplayName     string    `json:"display_name"`
	Description     string    `json:"description"`
	CategoryID      string    `json:"category_id"`
	SubCategories   []string  `json:"sub_categories"`
	
	// Contact Information
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Website         string    `json:"website"`
	Address         string    `json:"address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	
	// Business Details
	BusinessType    string    `json:"business_type"`
	EstablishedYear int       `json:"established_year"`
	EmployeesCount  int       `json:"employees_count"`
	
	// Service Areas
	ServiceAreas    []string  `json:"service_areas"`
	ServiceRadius   float64   `json:"service_radius"`
	
	// Media
	Logo            string    `json:"logo"`
	BannerImage     string    `json:"banner_image"`
	Gallery         []string  `json:"gallery"`
	
	// Verification & Status
	IsVerified      bool      `json:"is_verified"`
	VerificationStatus string  `json:"verification_status"`
	IsActive        bool      `json:"is_active"`
	IsFeatured      bool      `json:"is_featured"`
	
	// Admin Positioning
	AdminPosition   int       `json:"admin_position"`
	PositionCategory string   `json:"position_category"`
	FeaturedUntil   *time.Time `json:"featured_until"`
	
	// Ratings & Reviews
	AverageRating   float64   `json:"average_rating"`
	TotalReviews    int       `json:"total_reviews"`
	RatingBreakdown map[int]int `json:"rating_breakdown"`
	
	// Service Stats
	TotalServices   int       `json:"total_services"`
	ActiveServices  int       `json:"active_services"`
	CompletedBookings int     `json:"completed_bookings"`
	
	// Timestamps
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	LastActiveAt    time.Time `json:"last_active_at"`
}

// ProviderService represents a service offered by a provider
type ProviderService struct {
	ID              string    `json:"id"`
	ProviderID      string    `json:"provider_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	CategoryID      string    `json:"category_id"`
	
	// Pricing
	PriceType       string    `json:"price_type"`
	BasePrice       float64   `json:"base_price"`
	PriceRange      string    `json:"price_range"`
	
	// Availability
	IsActive        bool      `json:"is_active"`
	IsAvailable     bool      `json:"is_available"`
	AvailableHours  []TimeSlot `json:"available_hours"`
	
	// Service Details
	Duration        int       `json:"duration"`
	AdvanceNotice   int       `json:"advance_notice"`
	MaxBookingsPerDay int     `json:"max_bookings_per_day"`
	
	// Media
	Images          []string  `json:"images"`
	Video           string    `json:"video"`
	
	// Requirements
	Requirements    []string  `json:"requirements"`
	
	// Tags & Features
	Tags            []string  `json:"tags"`
	Features        []string  `json:"features"`
	
	// Location Service
	IsMobileService bool      `json:"is_mobile_service"`
	ServiceRadius   float64   `json:"service_radius"`
	
	// Stats
	BookingCount    int       `json:"booking_count"`
	CompletionRate  float64   `json:"completion_rate"`
	
	// Timestamps
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
}

// AdminPositioning represents admin-controlled provider positioning
type AdminPositioning struct {
	ID              string    `json:"id"`
	CategoryID      string    `json:"category_id"`
	ProviderID      string    `json:"provider_id"`
	Position        int       `json:"position"`
	IsFeatured      bool      `json:"is_featured"`
	FeaturedReason  string    `json:"featured_reason"`
	PriorityScore   float64   `json:"priority_score"`
	AdminNotes      string    `json:"admin_notes"`
	
	// Positioning Rules
	AlwaysOnTop     bool      `json:"always_on_top"`
	BoostFactor     float64   `json:"boost_factor"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      *time.Time `json:"valid_until"`
	
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Review represents a customer review for a provider
type Review struct {
	ID              string    `json:"id"`
	BookingID       string    `json:"booking_id"`
	ProviderID      string    `json:"provider_id"`
	CustomerID      string    `json:"customer_id"`
	Rating          int       `json:"rating"`
	Comment         string    `json:"comment"`
	
	// Review Categories
	QualityRating   int       `json:"quality_rating"`
	ValueRating     int       `json:"value_rating"`
	PunctualityRating int     `json:"punctuality_rating"`
	ProfessionalRating int    `json:"professional_rating"`
	
	// Images
	ReviewImages    []string  `json:"review_images"`
	
	// Moderation
	IsVerified      bool      `json:"is_verified"`
	IsApproved      bool      `json:"is_approved"`
	ReportedCount   int       `json:"reported_count"`
	
	// Response
	ProviderResponse string   `json:"provider_response"`
	RespondedAt     *time.Time `json:"responded_at"`
	
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TimeSlot represents a time slot for availability
type TimeSlot struct {
	DayOfWeek string `json:"day_of_week"` // "monday", "tuesday", etc.
	OpenTime  string `json:"open_time"`   // "09:00"
	CloseTime string `json:"close_time"`  // "17:00"
}

// ProviderListResponse represents the response for provider listings
type ProviderListResponse struct {
	Providers      []ServiceProvider `json:"providers"`
	Total          int               `json:"total"`
	Page           int               `json:"page"`
	Limit          int               `json:"limit"`
	HasNext        bool              `json:"has_next"`
	HasPrev        bool              `json:"has_prev"`
	Category       string            `json:"category"`
	FeaturedCount  int               `json:"featured_count"`
}

// ProviderSearchRequest represents a provider search request
type ProviderSearchRequest struct {
	Query          string   `json:"query"`
	Category       string   `json:"category"`
	SubCategories  []string `json:"sub_categories"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	Lat            float64  `json:"lat"`
	Lng            float64  `json:"lng"`
	Radius         float64  `json:"radius"`
	MinRating      float64  `json:"min_rating"`
	VerifiedOnly   bool     `json:"verified_only"`
	FeaturedOnly   bool     `json:"featured_only"`
	PriceRange     string   `json:"price_range"`
	SortBy         string   `json:"sort_by"` // "position", "rating", "reviews", "price"
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
}

// Constants for provider status
const (
	ProviderStatusDraft      = "draft"
	ProviderStatusPending    = "pending"
	ProviderStatusActive     = "active"
	ProviderStatusPaused     = "paused"
	ProviderStatusExpired    = "expired"
	ProviderStatusSuspended  = "suspended"
	
	VerificationStatusPending = "pending"
	VerificationStatusVerified = "verified"
	VerificationStatusRejected = "rejected"
	
	BusinessTypeIndividual   = "individual"
	BusinessTypeCompany      = "company"
	BusinessTypeFranchise    = "franchise"
	
	PriceTypeFixed    = "fixed"
	PriceTypeHourly   = "hourly"
	PriceTypePerItem  = "per_item"
	PriceTypeCustom   = "custom"
	
	SortByPosition    = "position"
	SortByRating      = "rating"
	SortByReviews     = "reviews"
	SortByPrice       = "price"
	SortByDistance    = "distance"
	SortByNewest      = "newest"
)
