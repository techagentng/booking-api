# TripsBook Database Schema Design

## Core Architecture: Multi-Service Platform

### 1. Enhanced User Model (Extension of Existing)
```go
type User struct {
    Model
    Fullname       string    `json:"fullname"`
    Username       string    `json:"username" gorm:"unique"`
    Telephone      string    `json:"telephone"`
    Email          string    `json:"email" gorm:"unique"`
    HashedPassword string    `json:"-"`
    IsSocial       bool      `json:"is_social"`
    AdminStatus    bool      `json:"is_admin"`
    RoleID         uuid.UUID `json:"role_id"`
    
    // New fields for TripsBook
    UserType       string    `json:"user_type" gorm:"default:'customer'"` // customer, provider, both
    ProfileImage   string    `json:"profile_image"`
    Bio            string    `json:"bio"`
    Location       string    `json:"location"`
    Latitude       *float64  `json:"latitude"`
    Longitude      *float64  `json:"longitude"`
    IsVerified     bool      `json:"is_verified" gorm:"default:false"`
    AverageRating  float64   `json:"average_rating" gorm:"default:0"`
    TotalReviews   int       `json:"total_reviews" gorm:"default:0"`
    
    // Provider specific fields
    BusinessName   string    `json:"business_name"`
    BusinessRegNo  string    `json:"business_reg_no"`
    ServiceRadius  int       `json:"service_radius"` // in kilometers
    
    Role           Role      `json:"role"`
}
```

### 2. Service Categories
```go
type ServiceCategory struct {
    Model
    Name         string `json:"name" gorm:"unique"`
    Description  string `json:"description"`
    Icon         string `json:"icon"`
    IsActive     bool   `json:"is_active" gorm:"default:true"`
    SortOrder    int    `json:"sort_order"`
    
    SubCategories []ServiceSubCategory `json:"sub_categories"`
}

type ServiceSubCategory struct {
    Model
    CategoryID   uint   `json:"category_id"`
    Name         string `json:"name"`
    Description  string `json:"description"`
    IsActive     bool   `json:"is_active" gorm:"default:true"`
    
    Category     ServiceCategory `json:"category"`
    Providers    []ServiceProvider `json:"providers"`
}
```

### 3. Service Providers
```go
type ServiceProvider struct {
    Model
    UserID           uint      `json:"user_id"`
    SubCategoryID    uint      `json:"sub_category_id"`
    BusinessName     string    `json:"business_name"`
    Description      string    `json:"description"`
    Address          string    `json:"address"`
    City             string    `json:"city"`
    State            string    `json:"state"`
    Latitude         *float64  `json:"latitude"`
    Longitude        *float64  `json:"longitude"`
    Phone            string    `json:"phone"`
    Email            string    `json:"email"`
    Website          string    `json:"website"`
    Logo             string    `json:"logo"`
    CoverImage       string    `json:"cover_image"`
    IsAvailable      bool      `json:"is_available" gorm:"default:true"`
    CommissionRate   float64   `json:"commission_rate" gorm:"default:0.05"` // 5%
    AverageRating    float64   `json:"average_rating" gorm:"default:0"`
    TotalReviews     int       `json:"total_reviews" gorm:"default:0"`
    TotalBookings    int       `json:"total_bookings" gorm:"default:0"`
    VerifiedAt       *time.Time `json:"verified_at"`
    Featured         bool      `json:"featured" gorm:"default:false"`
    
    User             User                  `json:"user"`
    SubCategory      ServiceSubCategory    `json:"sub_category"`
    Services         []Service             `json:"services"`
    Reviews          []Review              `json:"reviews"`
    Bookings         []Booking             `json:"bookings"`
    Availability     []ProviderAvailability `json:"availability"`
}
```

### 4. Services
```go
type Service struct {
    Model
    ProviderID       uint      `json:"provider_id"`
    Name             string    `json:"name"`
    Description      string    `json:"description"`
    BasePrice        float64   `json:"base_price"`
    PriceType        string    `json:"price_type" gorm:"default:'fixed'"` // fixed, hourly, per_person, custom
    Duration         int       `json:"duration"` // in minutes
    MinAdvanceNotice int       `json:"min_advance_notice"` // in hours
    CancellationPolicy string  `json:"cancellation_policy"`
    IsActive         bool      `json:"is_active" gorm:"default:true"`
    Featured         bool      `json:"featured" gorm:"default:false"`
    
    Provider         ServiceProvider `json:"provider"`
    Bookings         []Booking       `json:"bookings"`
    ServiceImages    []ServiceImage  `json:"images"`
}
```

### 5. Bookings (Universal)
```go
type Booking struct {
    Model
    BookingID        string    `json:"booking_id" gorm:"uniqueIndex;not null"`
    CustomerID       uint      `json:"customer_id"`
    ProviderID       uint      `json:"provider_id"`
    ServiceID        uint      `json:"service_id"`
    
    // Booking details
    ScheduledDate    time.Time `json:"scheduled_date"`
    StartTime        time.Time `json:"start_time"`
    EndTime          time.Time `json:"end_time"`
    Duration         int       `json:"duration"` // in minutes
    Location         string    `json:"location"`
    Latitude         *float64  `json:"latitude"`
    Longitude        *float64  `json:"longitude"`
    
    // Pricing
    BasePrice        float64   `json:"base_price"`
    AdditionalFees   float64   `json:"additional_fees"`
    CommissionAmount float64   `json:"commission_amount"`
    TotalAmount      float64   `json:"total_amount"`
    
    // Status tracking
    Status           string    `json:"status" gorm:"default:'pending'"` // pending, confirmed, in_progress, completed, cancelled, no_show
    PaymentStatus    string    `json:"payment_status" gorm:"default:'pending'"` // pending, paid, refunded
    CustomerNotes    string    `json:"customer_notes"`
    ProviderNotes    string    `json:"provider_notes"`
    
    // Timestamps
    ConfirmedAt      *time.Time `json:"confirmed_at"`
    StartedAt        *time.Time `json:"started_at"`
    CompletedAt      *time.Time `json:"completed_at"`
    CancelledAt      *time.Time `json:"cancelled_at"`
    CancelledBy      string     `json:"cancelled_by"` // customer, provider, system
    
    Customer         User            `json:"customer"`
    Provider         ServiceProvider `json:"provider"`
    Service          Service         `json:"service"`
    Reviews          []Review        `json:"reviews"`
    Payments         []Payment       `json:"payments"`
}
```

### 6. Reviews & Ratings
```go
type Review struct {
    Model
    BookingID        uint      `json:"booking_id"`
    CustomerID       uint      `json:"customer_id"`
    ProviderID       uint      `json:"provider_id"`
    Rating           float64   `json:"rating" gorm:"not null"` // 1.0 - 5.0
    Comment          string    `json:"comment"`
    
    // Service quality metrics
    Punctuality      float64   `json:"punctuality"`      // 1-5
    Professionalism  float64   `json:"professionalism"`  // 1-5
    Communication    float64   `json:"communication"`    // 1-5
    ValueForMoney    float64   `json:"value_for_money"`  // 1-5
    
    IsVerified       bool      `json:"is_verified" gorm:"default:false"`
    HelpfulCount     int       `json:"helpful_count" gorm:"default:0"`
    
    Booking          Booking         `json:"booking"`
    Customer         User            `json:"customer"`
    Provider         ServiceProvider `json:"provider"`
}
```

### 7. Provider Availability
```go
type ProviderAvailability struct {
    Model
    ProviderID       uint      `json:"provider_id"`
    DayOfWeek        int       `json:"day_of_week"` // 0-6 (Sunday-Saturday)
    StartTime        string    `json:"start_time"`   // "09:00"
    EndTime          string    `json:"end_time"`     // "17:00"
    IsAvailable      bool      `json:"is_available" gorm:"default:true"`
    
    Provider         ServiceProvider `json:"provider"`
}

type AvailabilityException struct {
    Model
    ProviderID       uint      `json:"provider_id"`
    ExceptionDate    time.Time `json:"exception_date"`
    StartTime        string    `json:"start_time"`
    EndTime          string    `json:"end_time"`
    IsUnavailable    bool      `json:"is_unavailable" gorm:"default:false"`
    Reason           string    `json:"reason"`
    
    Provider         ServiceProvider `json:"provider"`
}
```

### 8. Notifications & Messaging
```go
type Conversation struct {
    Model
    CustomerID       uint      `json:"customer_id"`
    ProviderID       uint      `json:"provider_id"`
    BookingID        *uint     `json:"booking_id"`
    LastMessageAt    time.Time `json:"last_message_at"`
    
    Customer         User            `json:"customer"`
    Provider         ServiceProvider `json:"provider"`
    Booking          *Booking         `json:"booking"`
    Messages         []Message       `json:"messages"`
}

type Message struct {
    Model
    ConversationID   uint      `json:"conversation_id"`
    SenderID         uint      `json:"sender_id"`
    SenderType       string    `json:"sender_type"` // customer, provider
    Content          string    `json:"content"`
    MessageType      string    `json:"message_type" gorm:"default:'text'"` // text, image, location
    IsRead           bool      `json:"is_read" gorm:"default:false"`
    ReadAt           *time.Time `json:"read_at"`
    
    Conversation     Conversation `json:"conversation"`
}
```

## Migration Strategy

### Phase 1: Isolate Hotel Booking
1. Create separate module structure:
   - `/hotel/` - Existing hall booking system
   - `/tripsbook/` - New service platform
   - `/shared/` - Common components (auth, payments, notifications)

### Phase 2: Extend Existing Models
1. Enhance User model with provider fields
2. Create new service-specific models
3. Add role-based access control

### Phase 3: API Layer
1. Separate routers for hotel vs tripsbook
2. Shared middleware for auth/payments
3. Provider onboarding endpoints

### Phase 4: Provider Onboarding
1. Bulk import system for 100 providers
2. Category management
3. Verification workflow

## Commission System Design

### Commission Calculation
```go
type Commission struct {
    Model
    BookingID        uint      `json:"booking_id"`
    ProviderID       uint      `json:"provider_id"`
    CommissionRate   float64   `json:"commission_rate"`
    BaseAmount       float64   `json:"base_amount"`
    CommissionAmount float64   `json:"commission_amount"`
    Status           string    `json:"status" gorm:"default:'pending'"` // pending, collected, paid
    CollectedAt      *time.Time `json:"collected_at"`
    PaidToProviderAt *time.Time `json:"paid_to_provider_at"`
    
    Booking          Booking         `json:"booking"`
    Provider         ServiceProvider `json:"provider"`
}
```

### Payment Flow
1. Customer pays full amount to TripsBook
2. TripsBook holds payment in escrow
3. After service completion:
   - 5% commission goes to TripsBook
   - 95% transferred to provider
4. Weekly payout system for providers

## API Structure

### Core Endpoints
```
/auth/*          - Shared authentication
/hotel/*         - Existing hall booking
/tripsbook/      - New service platform
  /providers/*   - Provider management
  /services/*    - Service catalog
  /bookings/*    - Booking management
  /reviews/*     - Rating system
  /messages/*    - Messaging
  /categories/*  - Service categories
```

This architecture allows you to:
1. Keep existing hotel booking functional
2. Build TripsBook alongside it
3. Share common infrastructure
4. Scale to 100+ providers
5. Implement 5% commission model
6. Support all 10 service categories
