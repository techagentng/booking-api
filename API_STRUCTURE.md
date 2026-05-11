# TripsBook API Structure Design

## API Architecture Overview

### Base URL Structure
```
https://api.tripsbook.com/v1/
```

### Module Separation
```
/hotel/           - Existing hall booking system
/tripsbook/       - New service provider platform
/shared/          - Common services (auth, payments, notifications)
/public/          - Public endpoints (no authentication)
/admin/           - Admin panel endpoints
```

## Authentication & Authorization

### JWT Token Structure
```json
{
  "user_id": 123,
  "user_type": "customer|provider|admin",
  "permissions": ["book_service", "manage_profile", "view_analytics"],
  "exp": 1640995200
}
```

### Middleware Layers
1. **Auth Middleware** - JWT validation
2. **Role Middleware** - Permission checking
3. **Rate Limiting** - API abuse prevention
4. **CORS** - Cross-origin requests

## Core API Endpoints

### 1. Authentication (`/auth/`)
```go
// Public endpoints
POST   /auth/register              // Customer registration
POST   /auth/login                 // User login
POST   /auth/provider/register     // Provider registration
POST   /auth/social/google         // Google OAuth
POST   /auth/social/facebook       // Facebook OAuth
POST   /auth/forgot-password       // Password reset
POST   /auth/verify-email/:token   // Email verification

// Protected endpoints
POST   /auth/refresh-token         // Token refresh
POST   /auth/logout                // Logout
PUT    /auth/change-password       // Change password
PUT    /auth/update-profile        // Update profile
```

### 2. Hotel Module (`/hotel/`)
```go
// Public endpoints
GET    /hotel/availability         // Hall availability
POST   /hotel/bookings/public      // Public booking request

// Protected endpoints
GET    /hotel/bookings             // User bookings
POST   /hotel/bookings             // Create booking
PUT    /hotel/bookings/:id         // Update booking
DELETE /hotel/bookings/:id         // Cancel booking
GET    /hotel/calendar             // Calendar view
```

### 3. TripsBook Core (`/tripsbook/`)

#### 3.1 Service Categories
```go
// Public endpoints
GET    /tripsbook/categories                    // All categories
GET    /tripsbook/categories/:id                // Category details
GET    /tripsbook/categories/:id/subcategories  // Subcategories
GET    /tripsbook/subcategories/:id/services     // Services in subcategory

// Admin endpoints
POST   /tripsbook/categories                    // Create category
PUT    /tripsbook/categories/:id                // Update category
DELETE /tripsbook/categories/:id                // Delete category
```

#### 3.2 Service Providers
```go
// Public endpoints
GET    /tripsbook/providers                     // List providers (with filters)
GET    /tripsbook/providers/:id                 // Provider profile
GET    /tripsbook/providers/:id/services       // Provider services
GET    /tripsbook/providers/:id/reviews        // Provider reviews
GET    /tripsbook/providers/search              // Search providers
GET    /tripsbook/providers/nearby              // Location-based search

// Provider endpoints
POST   /tripsbook/providers/profile             // Create/complete profile
PUT    /tripsbook/providers/profile             // Update profile
POST   /tripsbook/providers/documents           // Upload verification docs
GET    /tripsbook/providers/availability        // My availability
PUT    /tripsbook/providers/availability        // Update availability
GET    /tripsbook/providers/bookings           // My bookings
GET    /tripsbook/providers/analytics           // Performance analytics
GET    /tripsbook/providers/earnings            // Earnings report

// Admin endpoints
GET    /tripsbook/providers/admin/list          // All providers (admin)
PUT    /tripsbook/providers/:id/verify          // Verify provider
PUT    /tripsbook/providers/:id/feature         // Feature provider
POST   /tripsbook/providers/:id/suspend         // Suspend provider
```

#### 3.3 Services
```go
// Public endpoints
GET    /tripsbook/services                      // All services (with filters)
GET    /tripsbook/services/:id                  // Service details
GET    /tripsbook/services/search               // Search services

// Provider endpoints
POST   /tripsbook/services                      // Create service
PUT    /tripsbook/services/:id                  // Update service
DELETE /tripsbook/services/:id                  // Deactivate service
POST   /tripsbook/services/:id/images           // Upload service images
```

#### 3.4 Bookings
```go
// Customer endpoints
GET    /tripsbook/bookings                      // My bookings
POST   /tripsbook/bookings                      // Create booking
PUT    /tripsbook/bookings/:id                  // Update booking
DELETE /tripsbook/bookings/:id                  // Cancel booking
POST   /tripsbook/bookings/:id/review           // Review booking

// Provider endpoints
GET    /tripsbook/bookings/provider             // My bookings (provider view)
PUT    /tripsbook/bookings/:id/confirm          // Confirm booking
PUT    /tripsbook/bookings/:id/start            // Start service
PUT    /tripsbook/bookings/:id/complete         // Complete service
PUT    /tripsbook/bookings/:id/noshow           // Mark no-show

// Common endpoints
GET    /tripsbook/bookings/:id                  // Booking details
POST   /tripsbook/bookings/:id/messages         // Send message
GET    /tripsbook/bookings/:id/messages         // Booking messages
```

#### 3.5 Reviews & Ratings
```go
// Public endpoints
GET    /tripsbook/reviews/provider/:id          // Provider reviews
GET    /tripsbook/reviews/service/:id           // Service reviews

// Customer endpoints
POST   /tripsbook/reviews                       // Create review
PUT    /tripsbook/reviews/:id                   // Update review
DELETE /tripsbook/reviews/:id                   // Delete review
POST   /tripsbook/reviews/:id/helpful          // Mark as helpful

// Provider endpoints
GET    /tripsbook/reviews/my                    // My reviews
POST   /tripsbook/reviews/:id/response          // Respond to review
```

#### 3.6 Messaging
```go
GET    /tripsbook/conversations                 // My conversations
GET    /tripsbook/conversations/:id              // Conversation details
POST   /tripsbook/conversations                 // Start conversation
POST   /tripsbook/conversations/:id/messages     // Send message
PUT    /tripsbook/messages/:id/read             // Mark messages as read
```

### 4. Payments (`/payments/`)
```go
// Customer endpoints
POST   /payments/intent/booking                  // Create payment intent
POST   /payments/confirm                        // Confirm payment
GET    /payments/history                        // Payment history

// Provider endpoints
GET    /payments/earnings                       // My earnings
GET    /payments/payouts                        // Payout history
POST   /payments/payout/request                 // Request payout

// Webhooks
POST   /payments/webhooks/stripe                // Stripe webhook
```

### 5. Admin Panel (`/admin/`)
```go
GET    /admin/dashboard                        // Admin dashboard
GET    /admin/providers                         // Provider management
GET    /admin/customers                         // Customer management
GET    /admin/bookings                          // Booking overview
GET    /admin/payments                          // Payment overview
GET    /admin/reports                           // Analytics reports
POST   /admin/providers/bulk-import             // Bulk provider import
```

## Request/Response Formats

### Standard Response Structure
```json
{
  "success": true,
  "data": {},
  "message": "Operation successful",
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req_123456789"
  }
}
```

### Error Response Structure
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      }
    ]
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req_123456789"
  }
}
```

### Paginated Response Structure
```json
{
  "success": true,
  "data": [],
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5,
      "has_next": true,
      "has_prev": false
    },
    "filters": {
      "category": "transport",
      "location": "Abuja",
      "rating_min": 4.0
    }
  }
}
```

## API Filtering & Search

### Provider Search Parameters
```go
GET /tripsbook/providers?
    category=transport&
    subcategory=bolt_drivers&
    location=Abuja&
    lat=9.0579&
    lng=7.4951&
    radius=10&
    rating_min=4.0&
    price_max=5000&
    available=true&
    featured=true&
    page=1&
    limit=20
```

### Service Search Parameters
```go
GET /tripsbook/services?
    query=airport%20pickup&
    category=airport_services&
    min_price=1000&
    max_price=10000&
    duration_min=30&
    duration_max=180&
    sort=price_asc|rating_desc|distance_asc
```

## Rate Limiting Strategy

### Rate Limits by Endpoint Type
```go
// Authentication endpoints
auth endpoints:        5 requests/minute/IP

// Public search endpoints
public search:        100 requests/minute/IP

// Customer endpoints
customer endpoints:    1000 requests/hour/user

// Provider endpoints  
provider endpoints:    500 requests/hour/user

// Admin endpoints
admin endpoints:       2000 requests/hour/user
```

## API Versioning Strategy

### URL Versioning
```
/v1/ - Current stable version
/v2/ - Future major updates
```

### Backward Compatibility
- Maintain v1 endpoints for 12 months after v2 release
- Use feature flags for gradual rollout
- Provide migration guides for API consumers

## Security Implementation

### 1. Authentication Security
```go
// JWT Configuration
- Access token expiry: 15 minutes
- Refresh token expiry: 7 days
- Token rotation on refresh
- Blacklist compromised tokens
```

### 2. Data Validation
```go
// Input validation rules
- Email format validation
- Phone number validation (NG format)
- Location coordinate validation
- Price range validation
- File upload restrictions
```

### 3. SQL Injection Prevention
```go
// GORM parameterized queries
- Use GORM's built-in protection
- Raw queries with proper escaping
- Input sanitization
```

### 4. CORS Configuration
```go
// CORS settings
- Allowed origins: whitelist
- Allowed methods: GET, POST, PUT, DELETE
- Allowed headers: Authorization, Content-Type
- Max age: 86400 seconds
```

## Performance Optimization

### 1. Database Optimization
```go
// Indexes for common queries
- Provider location (latitude, longitude)
- Provider category and subcategory
- Booking dates and status
- User email and phone
- Service price ranges
```

### 2. Caching Strategy
```go
// Redis caching
- Provider profiles: 1 hour
- Service categories: 24 hours
- Popular searches: 30 minutes
- User sessions: 15 minutes
```

### 3. Pagination
```go
// Cursor-based pagination for large datasets
- Efficient for mobile apps
- Stable during data changes
- Better performance than offset
```

## Monitoring & Logging

### 1. API Metrics
```go
// Key metrics to track
- Response times by endpoint
- Error rates by type
- Request volume by endpoint
- User authentication success/failure
- Payment success/failure rates
```

### 2. Logging Strategy
```go
// Structured logging
- Request ID for traceability
- User ID for audit trail
- Action performed
- Timestamp and duration
- Error details with stack traces
```

### 3. Health Checks
```go
// Health check endpoints
GET /health                    // Basic health check
GET /health/db                 // Database connectivity
GET /health/redis              // Cache connectivity
GET /health/stripe             // Payment service
GET /health/email              // Email service
```

## API Documentation

### 1. OpenAPI Specification
```yaml
# Comprehensive API documentation
- Endpoint descriptions
- Request/response schemas
- Authentication requirements
- Error codes and meanings
- Example requests/responses
```

### 2. SDK Development
```go
// Client SDKs planned
- JavaScript/TypeScript
- Go (for server-to-server)
- Python (for data analysis)
- Swift (iOS app)
- Kotlin (Android app)
```

This API structure provides:
1. **Clear separation** between hotel and TripsBook functionality
2. **Scalable architecture** for 100+ providers
3. **Comprehensive feature set** for all user types
4. **Security-first approach** with proper authentication and validation
5. **Performance optimization** with caching and indexing
6. **Monitoring capabilities** for operational excellence
7. **Future-proof design** with versioning and extensibility
