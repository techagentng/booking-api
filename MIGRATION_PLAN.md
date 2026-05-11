# TripsBook Migration & Implementation Plan

## Phase 1: Module Restructuring (Week 1-2)

### 1.1 Directory Structure Reorganization
```
Hotel/
├── main.go                    # Entry point - handles both platforms
├── shared/                    # Common components
│   ├── auth/                  # Authentication system
│   ├── payments/              # Stripe integration
│   ├── notifications/         # Email/SMS services
│   ├── middleware/            # Common middleware
│   ├── models/                # Base models (User, Role, etc.)
│   └── utils/                 # Helper functions
├── hotel/                     # Existing hall booking
│   ├── models/                # Hall-specific models
│   ├── handlers/              # Hall booking handlers
│   ├── services/              # Hall business logic
│   └── routes/                # Hall routes
├── tripsbook/                 # New service platform
│   ├── models/                # Service provider models
│   ├── handlers/              # TripsBook handlers
│   ├── services/              # TripsBook business logic
│   └── routes/                # TripsBook routes
└── config/                    # Configuration management
```

### 1.2 Code Migration Steps
1. **Extract Shared Components**
   - Move `models/user.go`, `models/role.go` to `shared/models/`
   - Move auth handlers to `shared/auth/`
   - Move payment logic to `shared/payments/`
   - Move notification system to `shared/notifications/`

2. **Isolate Hotel Module**
   - Move hall booking specific files to `hotel/`
   - Update imports to use shared components
   - Create hotel-specific router

3. **Create TripsBook Module**
   - Set up initial module structure
   - Create base models and handlers
   - Set up TripsBook router

## Phase 2: Database Schema Implementation (Week 2-3)

### 2.1 Migration Strategy
```sql
-- Step 1: Extend existing User table
ALTER TABLE users ADD COLUMN user_type VARCHAR(20) DEFAULT 'customer';
ALTER TABLE users ADD COLUMN profile_image TEXT;
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN location VARCHAR(255);
ALTER TABLE users ADD COLUMN latitude DECIMAL(10, 8);
ALTER TABLE users ADD COLUMN longitude DECIMAL(11, 8);
ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN average_rating DECIMAL(3, 2) DEFAULT 0;
ALTER TABLE users ADD COLUMN total_reviews INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN business_name VARCHAR(255);
ALTER TABLE users ADD COLUMN business_reg_no VARCHAR(100);
ALTER TABLE users ADD COLUMN service_radius INTEGER;

-- Step 2: Create new TripsBook tables
CREATE TABLE service_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE service_sub_categories (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES service_categories(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 2.2 GORM Model Implementation
- Create all models defined in schema
- Implement relationships and constraints
- Add database indexes for performance
- Set up soft deletes

## Phase 3: API Development (Week 3-4)

### 3.1 Router Configuration
```go
// main.go - Updated router setup
func main() {
    // ... existing setup ...
    
    // Shared router
    sharedRouter := gin.Group("/api/v1")
    sharedRouter.Use(middleware.AuthMiddleware())
    
    // Hotel routes
    hotel.SetupRoutes(sharedRouter.Group("/hotel"))
    
    // TripsBook routes
    tripsbook.SetupRoutes(sharedRouter.Group("/tripsbook"))
    
    // Public routes (no auth required)
    publicRouter := gin.Group("/api/v1/public")
    tripsbook.SetupPublicRoutes(publicRouter)
}
```

### 3.2 API Endpoints Priority
1. **Week 3 - Core Provider APIs**
   - `POST /tripsbook/providers/register` - Provider registration
   - `GET /tripsbook/providers` - List providers (with filters)
   - `GET /tripsbook/providers/:id` - Provider profile
   - `PUT /tripsbook/providers/:id` - Update provider profile

2. **Week 4 - Booking & Service APIs**
   - `GET /tripsbook/categories` - Service categories
   - `GET /tripsbook/services` - Service catalog
   - `POST /tripsbook/bookings` - Create booking
   - `GET /tripsbook/bookings` - User bookings
   - `POST /tripsbook/reviews` - Submit review

## Phase 4: Provider Onboarding System (Week 4-5)

### 4.1 Bulk Provider Import
```go
// tripsbook/services/provider_onboarding.go
type ProviderImport struct {
    BusinessName    string `json:"business_name"`
    ContactName     string `json:"contact_name"`
    Email          string `json:"email"`
    Phone          string `json:"phone"`
    Category       string `json:"category"`
    SubCategory    string `json:"sub_category"`
    Address        string `json:"address"`
    City           string `json:"city"`
    State          string `json:"state"`
    Description    string `json:"description"`
}

func (s *ProviderService) BulkImportProviders(providers []ProviderImport) error {
    // Implementation for bulk import
    // 1. Validate data
    // 2. Create user accounts
    // 3. Create provider profiles
    // 4. Send welcome emails
}
```

### 4.2 Provider Categories Setup
```go
// Initial seed data for 100 target providers
var ServiceCategories = []ServiceCategory{
    {
        Name: "Transport & Mobility",
        SubCategories: []ServiceSubCategory{
            {Name: "Bolt Drivers"},
            {Name: "Uber Drivers"},
            {Name: "InDrive Drivers"},
            {Name: "Private Drivers"},
            {Name: "Chauffeur Services"},
            {Name: "Car Hire Companies"},
            {Name: "Self-drive Car Rental"},
        },
    },
    {
        Name: "Airport & Travel Services",
        SubCategories: []ServiceSubCategory{
            {Name: "Airport Pickup Operators"},
            {Name: "Travel Assistants"},
            {Name: "Tour Guides"},
        },
    },
    // ... other 8 categories
}
```

## Phase 5: Commission & Payment System (Week 5-6)

### 5.1 Commission Calculation Service
```go
type CommissionService struct {
    repo *CommissionRepository
    stripeService *StripeService
}

func (s *CommissionService) CalculateCommission(booking *Booking) *Commission {
    commissionRate := 0.05 // 5%
    commissionAmount := booking.TotalAmount * commissionRate
    
    return &Commission{
        BookingID:        booking.ID,
        ProviderID:       booking.ProviderID,
        CommissionRate:   commissionRate,
        BaseAmount:       booking.TotalAmount,
        CommissionAmount: commissionAmount,
        Status:          "pending",
    }
}

func (s *CommissionService) ProcessPayout(providerID uint) error {
    // Weekly payout processing
    // 1. Calculate total commission for provider
    // 2. Transfer provider's share (95%)
    // 3. Keep commission (5%)
    // 4. Update commission records
}
```

### 5.2 Payment Flow Integration
```go
// Enhanced payment processing for TripsBook
func (s *PaymentService) ProcessServiceBookingPayment(booking *Booking) error {
    // 1. Create payment intent for full amount
    // 2. Hold payment in escrow
    // 3. Release to provider after completion
    // 4. Deduct commission automatically
}
```

## Phase 6: Testing & Deployment (Week 6-7)

### 6.1 Testing Strategy
1. **Unit Tests**
   - Model validation
   - Business logic
   - Commission calculations

2. **Integration Tests**
   - API endpoints
   - Payment processing
   - Database operations

3. **Load Testing**
   - Provider search performance
   - Booking system under load
   - Payment processing speed

### 6.2 Deployment Checklist
- [ ] Database migrations tested
- [ ] API endpoints documented
- [ ] Payment system verified
- [ ] Provider onboarding tested
- [ ] Monitoring and logging setup
- [ ] Backup procedures verified

## Risk Mitigation

### Technical Risks
1. **Database Performance**
   - Add indexes for location-based queries
   - Implement caching for provider searches
   - Monitor query performance

2. **Payment Processing**
   - Test refund flows thoroughly
   - Implement retry mechanisms
   - Set up fraud detection

3. **Scalability**
   - Design for 1000+ providers
   - Implement rate limiting
   - Plan for horizontal scaling

### Business Risks
1. **Provider Adoption**
   - Simple onboarding process
   - Clear value proposition
   - Support system for providers

2. **Customer Trust**
   - Verification system
   - Review authenticity
   - Dispute resolution

## Success Metrics

### Technical KPIs
- API response time < 200ms
- 99.9% uptime
- Zero payment failures
- < 1% booking errors

### Business KPIs
- 100 providers onboarded in Phase 1
- 5% commission rate achieved
- 4.0+ average provider rating
- 50+ bookings per week by Month 2

## Next Steps

1. **Immediate (This Week)**
   - Start module restructuring
   - Set up new directory structure
   - Begin shared component extraction

2. **Short Term (Week 2-3)**
   - Implement database schema
   - Create core models
   - Set up basic API structure

3. **Medium Term (Week 4-6)**
   - Complete API development
   - Implement provider onboarding
   - Set up commission system

4. **Long Term (Week 7+)**
   - Testing and optimization
   - Provider recruitment
   - Marketing and launch

This plan allows you to build TripsBook while keeping your hotel booking system operational, sharing infrastructure where beneficial while maintaining clear separation between the two platforms.
