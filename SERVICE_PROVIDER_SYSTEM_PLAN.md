# Service Provider System Design Plan

## Overview
Design a comprehensive service provider management system with admin-controlled positioning, ratings, and service lifecycle management.

## 1. Service Provider Data Structure

### Core ServiceProvider Model
```go
type ServiceProvider struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`          // Link to user account
    BusinessName    string    `json:"business_name"`
    DisplayName     string    `json:"display_name"`     // Name shown to customers
    Description     string    `json:"description"`
    CategoryID      string    `json:"category_id"`      // Primary category
    SubCategories   []string  `json:"sub_categories"`  // Multiple sub-categories
    
    // Contact Information
    Phone           string    `json:"phone"`
    Email           string    `json:"email"`
    Website         string    `json:"website"`
    Address         string    `json:"address"`
    City            string    `json:"city"`
    State           string    `json:"state"`
    
    // Business Details
    BusinessType    string    `json:"business_type"`    // "individual" | "company" | "franchise"
    EstablishedYear int       `json:"established_year"`
    EmployeesCount  int       `json:"employees_count"`
    
    // Service Areas
    ServiceAreas    []string  `json:"service_areas"`    // Cities/areas they serve
    ServiceRadius   float64   `json:"service_radius"`   // KM radius from base location
    
    // Media
    Logo            string    `json:"logo"`
    BannerImage     string    `json:"banner_image"`
    Gallery         []string  `json:"gallery"`
    
    // Verification & Status
    IsVerified      bool      `json:"is_verified"`
    VerificationStatus string  `json:"verification_status"` // "pending" | "verified" | "rejected"
    IsActive        bool      `json:"is_active"`        // Provider can receive bookings
    IsFeatured      bool      `json:"is_featured"`      // Admin-selected featured provider
    
    // Admin Positioning
    AdminPosition   int       `json:"admin_position"`   // 1-100, lower = higher priority
    PositionCategory string   `json:"position_category"` // Category for this positioning
    FeaturedUntil   *time.Time `json:"featured_until"`  // Temporary featured status
    
    // Ratings & Reviews
    AverageRating   float64   `json:"average_rating"`
    TotalReviews    int       `json:"total_reviews"`
    RatingBreakdown map[int]int `json:"rating_breakdown"` // 5-star: 45, 4-star: 23, etc.
    
    // Service Stats
    TotalServices   int       `json:"total_services"`
    ActiveServices  int       `json:"active_services"`
    CompletedBookings int     `json:"completed_bookings"`
    
    // Timestamps
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    LastActiveAt    time.Time `json:"last_active_at"`
}
```

### Service Model (Provider's Services)
```go
type ProviderService struct {
    ID              string    `json:"id"`
    ProviderID      string    `json:"provider_id"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    CategoryID      string    `json:"category_id"`
    
    // Pricing
    PriceType       string    `json:"price_type"`        // "fixed" | "hourly" | "per_item" | "custom"
    BasePrice       float64   `json:"base_price"`
    PriceRange      string    `json:"price_range"`       // "₦5000-₦15000"
    
    // Availability
    IsActive        bool      `json:"is_active"`          // Service is visible to customers
    IsAvailable     bool      `json:"is_available"`       // Currently accepting bookings
    AvailableHours  []TimeSlot `json:"available_hours"`   // Weekly schedule
    
    // Service Details
    Duration        int       `json:"duration"`           // Minutes
    AdvanceNotice   int       `json:"advance_notice"`     // Hours notice required
    MaxBookingsPerDay int     `json:"max_bookings_per_day"`
    
    // Media
    Images          []string  `json:"images"`
    Video           string    `json:"video"`
    
    // Requirements
    Requirements    []string  `json:"requirements"`      // What customer needs to provide
    
    // Tags & Features
    Tags            []string  `json:"tags"`
    Features        []string  `json:"features"`
    
    // Location Service
    IsMobileService bool      `json:"is_mobile_service"`  // Provider travels to customer
    ServiceRadius   float64   `json:"service_radius"`     // KM for mobile service
    
    // Stats
    BookingCount    int       `json:"booking_count"`
    CompletionRate  float64   `json:"completion_rate"`
    
    // Timestamps
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    ExpiresAt       *time.Time `json:"expires_at"`       // Service expiration
}
```

## 2. Admin Positioning System

### AdminPositioning Model
```go
type AdminPositioning struct {
    ID              string    `json:"id"`
    CategoryID      string    `json:"category_id"`
    ProviderID      string    `json:"provider_id"`
    Position        int       `json:"position"`          // 1-100 (1 = highest priority)
    IsFeatured      bool      `json:"is_featured"`       // Special featured status
    FeaturedReason   string    `json:"featured_reason"`   // Why admin featured them
    PriorityScore   float64   `json:"priority_score"`    // Algorithm-calculated priority
    AdminNotes      string    `json:"admin_notes"`       // Internal admin notes
    
    // Positioning Rules
    AlwaysOnTop     bool      `json:"always_on_top"`     // Force to top regardless of algorithm
    BoostFactor     float64   `json:"boost_factor"`      // 1.0-10.0 multiplier
    ValidFrom       time.Time `json:"valid_from"`
    ValidUntil      *time.Time `json:"valid_until"`      // Null = permanent
    
    CreatedBy       string    `json:"created_by"`        // Admin user ID
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### Positioning Algorithm
```go
type PositioningCriteria struct {
    RatingWeight    float64 `json:"rating_weight"`       // 0.0-1.0
    ReviewsWeight   float64 `json:"reviews_weight"`      // 0.0-1.0
    BookingsWeight  float64 `json:"bookings_weight"`     // 0.0-1.0
    ResponseWeight  float64 `json:"response_weight"`     // 0.0-1.0
    VerificationWeight float64 `json:"verification_weight"` // 0.0-1.0
    FeaturedWeight  float64 `json:"featured_weight"`     // 0.0-1.0
}
```

## 3. Rating System

### Review Model
```go
type Review struct {
    ID              string    `json:"id"`
    BookingID       string    `json:"booking_id"`        // Link to completed booking
    ProviderID      string    `json:"provider_id"`
    CustomerID      string    `json:"customer_id"`
    Rating          int       `json:"rating"`            // 1-5 stars
    Comment         string    `json:"comment"`
    
    // Review Categories
    QualityRating   int       `json:"quality_rating"`    // Service quality
    ValueRating     int       `json:"value_rating"`      // Value for money
    PunctualityRating int     `json:"punctuality_rating"` // On-time performance
    ProfessionalRating int    `json:"professional_rating"` // Professionalism
    
    // Images
    ReviewImages    []string  `json:"review_images"`
    
    // Moderation
    IsVerified      bool      `json:"is_verified"`       // Verified purchase
    IsApproved      bool      `json:"is_approved"`       // Admin approved
    ReportedCount   int       `json:"reported_count"`
    
    // Response
    ProviderResponse string   `json:"provider_response"` // Provider's reply
    RespondedAt     *time.Time `json:"responded_at"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

## 4. Service Lifecycle Management

### Service Status Flow
1. **Draft** → Provider creates service, not visible
2. **Pending** → Submitted for admin approval
3. **Active** → Approved and visible to customers
4. **Paused** → Temporarily hidden by provider
5. **Expired** → Automatically or manually expired
6. **Suspended** → Admin suspended for policy violation

### Service Expiration Rules
- **Auto-expire**: Services older than X days without activity
- **Manual expire**: Provider chooses to expire
- **Admin expire**: Admin forces expiration
- **Seasonal expire**: Time-limited seasonal services

## 5. Mock Data Strategy

### Mock Service Providers (100 providers across categories)

#### Hotels Category (Top 20)
1. **Eko Hotels & Suites** - Position 1, Featured, 4.6★, 1200+ reviews
2. **Federal Palace Hotel** - Position 2, Featured, 4.5★, 980+ reviews  
3. **Lagos Continental Hotel** - Position 3, 4.4★, 750+ reviews
4. **Southern Sun Ikoyi** - Position 4, Featured, 4.7★, 620+ reviews
5. **The Wheatbaker** - Position 5, 4.8★, 450+ reviews
... (15 more)

#### Restaurants Category (Top 20)
1. **Terra Kulture** - Position 1, Featured, 4.5★, 890+ reviews
2. **Nigerian Kitchen** - Position 2, 4.6★, 720+ reviews
3. **Bukka Hut** - Position 3, Featured, 4.3★, 1100+ reviews
4. **Yellow Chilli** - Position 4, 4.7★, 580+ reviews
5. **Ocean Basket** - Position 5, 4.4★, 420+ reviews
... (15 more)

#### Transport Category (Top 20)
1. **Uber Premium** - Position 1, Featured, 4.7★, 2500+ reviews
2. **Bolt Professional** - Position 2, 4.5★, 1800+ reviews
3. **Lagos Taxi** - Position 3, Featured, 4.2★, 950+ reviews
4. **Airport Shuttle Pro** - Position 4, 4.6★, 680+ reviews
5. **Corporate Car Service** - Position 5, 4.8★, 320+ reviews
... (15 more)

#### Shopping Category (Top 20)
1. **Palms Shopping Mall** - Position 1, Featured, 4.4★, 1500+ reviews
2. **Ikeja City Mall** - Position 2, 4.3★, 1200+ reviews
3. **Eko Atlantic Mall** - Position 3, Featured, 4.6★, 890+ reviews
4. **Adeniran Ogunsanya Mall** - Position 4, 4.2★, 780+ reviews
5. **Circle Mall** - Position 5, 4.5★, 450+ reviews
... (15 more)

#### Events Category (Top 20)
1. **Eko Convention Center** - Position 1, Featured, 4.8★, 320+ reviews
2. **Landmark Event Center** - Position 2, 4.6★, 280+ reviews
3. **Terra Events** - Position 3, Featured, 4.5★, 190+ reviews
4. **Federal Palace Events** - Position 4, 4.7★, 150+ reviews
5. **Muson Center** - Position 5, 4.4★, 120+ reviews
... (15 more)

### Provider Positioning Examples

#### Admin Control Panel Features
- **Drag & Drop Reordering**: Visual positioning interface
- **Bulk Positioning**: Set multiple providers at once
- **Featured Slots**: Limited featured positions per category
- **Priority Boost**: Temporary visibility boost
- **Category Transfer**: Move providers between categories
- **Position Analytics**: See impact of positioning changes

#### Positioning Rules
1. **Positions 1-5**: Premium featured spots (paid/admin selected)
2. **Positions 6-20**: High-performing organic providers
3. **Positions 21-50**: Established providers with good ratings
4. **Positions 51-100**: New or lower-rated providers

## 6. API Endpoints

### Provider Management
```
POST   /api/v1/admin/providers                    // Create provider
GET    /api/v1/admin/providers                    // List all providers
GET    /api/v1/admin/providers/{id}               // Get provider details
PUT    /api/v1/admin/providers/{id}               // Update provider
DELETE /api/v1/admin/providers/{id}               // Delete provider
POST   /api/v1/admin/providers/{id}/verify        // Verify provider
POST   /api/v1/admin/providers/{id}/suspend       // Suspend provider
```

### Position Management
```
POST   /api/v1/admin/positioning                   // Set provider positions
GET    /api/v1/admin/positioning/{category}       // Get category positioning
PUT    /api/v1/admin/positioning/{provider}       // Update provider position
DELETE /api/v1/admin/positioning/{provider}       // Remove positioning
POST   /api/v1/admin/positioning/bulk             // Bulk positioning
POST   /api/v1/admin/positioning/featured          // Set featured providers
```

### Service Management
```
POST   /api/v1/providers/services                  // Create service
GET    /api/v1/providers/services                  // List provider services
PUT    /api/v1/providers/services/{id}             // Update service
DELETE /api/v1/providers/services/{id}             // Delete service
POST   /api/v1/providers/services/{id}/expire      // Expire service
POST   /api/v1/providers/services/{id}/reactivate  // Reactivate service
```

### Customer-Facing APIs
```
GET    /api/v1/public/providers/{category}        // Get providers by category
GET    /api/v1/public/providers/{category}/featured // Get featured providers
GET    /api/v1/public/providers/search             // Search providers
GET    /api/v1/public/providers/{id}               // Get provider details
GET    /api/v1/public/providers/{id}/services      // Get provider services
GET    /api/v1/public/providers/{id}/reviews        // Get provider reviews
```

## 7. Implementation Priority

### Phase 1: Core Models & APIs
- ServiceProvider and ProviderService models
- Basic CRUD operations
- Mock data generation

### Phase 2: Admin Positioning
- Admin positioning interface
- Positioning algorithm
- Featured provider system

### Phase 3: Rating System
- Review submission and management
- Rating calculation
- Review moderation

### Phase 4: Service Lifecycle
- Service creation workflow
- Expiration system
- Status management

### Phase 5: Advanced Features
- Analytics dashboard
- Automated positioning suggestions
- Performance tracking

## 8. Mock Data Implementation

The mock data will include:
- 100 service providers across 5 categories
- 500+ services with varying details
- 2000+ reviews with realistic ratings
- Admin positioning for top 100 positions per category
- Realistic business information and contact details

This system provides complete admin control over provider visibility while maintaining fairness through the positioning algorithm and rating system.
