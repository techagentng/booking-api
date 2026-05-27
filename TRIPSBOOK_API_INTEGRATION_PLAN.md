# TripsBook API Integration Plan

## Current Implementation Status

### ✅ **COMPLETED MODULES**

#### 1. **Authentication & Authorization**
- ✅ User registration (signup)
- ✅ User login (JWT tokens)
- ✅ Google OAuth integration
- ✅ Email verification (for providers)
- ✅ Role-based access (Admin, Provider, Customer, Moderator)
- ✅ Authorization middleware

**Endpoints:**
- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `POST /api/v1/google/user/login`
- `POST /api/v1/auth/logout`
- `POST /api/v1/public/onboarding/register` (provider)
- `POST /api/v1/public/onboarding/verify-email`

#### 2. **Customer Identity System**
- ✅ Customer profile management
- ✅ Saved services/favorites
- ✅ Customer booking history
- ✅ Customer preferences
- ✅ Auto-create customer on login

**Endpoints:**
- `GET /api/v1/customer/me`
- `PATCH /api/v1/customer/me`
- `GET /api/v1/customer/saved`
- `POST /api/v1/customer/saved`
- `DELETE /api/v1/customer/saved/:id`
- `GET /api/v1/customer/bookings`
- `GET /api/v1/customer/preferences`
- `PATCH /api/v1/customer/preferences`

#### 3. **Provider Onboarding System**
- ✅ Provider registration
- ✅ Email verification
- ✅ Business information management
- ✅ Service creation
- ✅ Verification status tracking
- ✅ Training completion
- ✅ Account activation

**Endpoints:**
- `POST /api/v1/public/onboarding/register`
- `POST /api/v1/public/onboarding/verify-email`
- `PUT /api/v1/provider/onboarding/business-info`
- `POST /api/v1/provider/onboarding/services`
- `GET /api/v1/provider/onboarding/verification`
- `POST /api/v1/provider/onboarding/training/:module_id`
- `POST /api/v1/provider/onboarding/activate`
- `GET /api/v1/provider/onboarding/status`

#### 4. **Hotel/Hall Booking System**
- ✅ Room management
- ✅ Reservation management
- ✅ Hall booking
- ✅ Check-in/check-out
- ✅ Calendar availability
- ✅ Guest management
- ✅ Staff management
- ✅ Room service menu
- ✅ Service requests

**Endpoints:**
- `GET/POST /api/v1/rooms`
- `GET/POST /api/v1/reservations`
- `GET/POST /api/v1/hall-bookings`
- `GET/POST /api/v1/guests`
- `GET/POST /api/v1/staff`
- `GET /api/v1/calendar/availability`
- `GET/POST /api/v1/room-service/orders`
- `GET/POST /api/v1/service-requests`

#### 5. **Payment Integration**
- ✅ Stripe integration
- ✅ Payment processing
- ✅ Payment details retrieval

**Endpoints:**
- `GET /api/v1/reservations/:id/payment`

#### 6. **Notifications**
- ✅ Server-Sent Events (SSE)
- ✅ Notification hub
- ✅ Test notifications
- ✅ Notification stats

**Endpoints:**
- `GET /api/v1/notifications/stream`
- `GET /api/v1/notifications/stats`
- `POST /api/v1/notifications/test`

---

### ❌ **MISSING MODULES**

#### 1. **Service Discovery & Search**
**Priority: HIGH**

**Required Features:**
- Search services by category, location, name
- Filter by price, rating, availability
- Nearby businesses discovery
- Service categories (hotels, restaurants, transport, experiences)
- Search without login (public access)

**Database Models Needed:**
```go
type ServiceCategory struct {
    ID          uuid.UUID
    Name        string
    Description string
    Icon        string
    ParentID    uuid.UUID
}

type ServiceSearchResult struct {
    ServiceID   uuid.UUID
    ProviderID  uuid.UUID
    Title       string
    Category    string
    Location    string
    Price       float64
    Rating      float64
    ImageURL    string
    Distance    float64 // for nearby search
}
```

**API Endpoints Needed:**
- `GET /api/v1/public/services/search?q=query&category=cat&location=loc`
- `GET /api/v1/public/services/categories`
- `GET /api/v1/public/services/nearby?lat=lat&lng=lng&radius=km`
- `GET /api/v1/public/services/:id`
- `GET /api/v1/public/providers/:id/services`

#### 2. **Provider Service Management**
**Priority: HIGH**

**Required Features:**
- CRUD operations for provider services
- Service availability management
- Service pricing management
- Service images
- Service features
- Service reviews/ratings

**Database Models Needed:**
```go
type ProviderService struct {
    ID          uuid.UUID
    ProviderID  uuid.UUID
    Title       string
    Description string
    CategoryID  string
    PriceType   string // fixed, hourly, per_item
    BasePrice   float64
    Duration    int
    Features    []string
    Images      []string
    IsActive    bool
    IsAvailable bool
}

type ServiceReview struct {
    ID          uuid.UUID
    ServiceID   uuid.UUID
    CustomerID  uuid.UUID
    Rating      int
    Comment     string
    CreatedAt   time.Time
}
```

**API Endpoints Needed:**
- `GET /api/v1/provider/services` (provider's own services)
- `POST /api/v1/provider/services`
- `PUT /api/v1/provider/services/:id`
- `DELETE /api/v1/provider/services/:id`
- `PATCH /api/v1/provider/services/:id/availability`
- `GET /api/v1/public/services/:id/reviews`
- `POST /api/v1/customer/services/:id/reviews`

#### 3. **Service Booking System**
**Priority: HIGH**

**Required Features:**
- Book services (not just hotels)
- Booking confirmation
- Booking status tracking
- Booking cancellation
- Booking modification
- Provider acceptance/rejection

**Database Models Needed:**
```go
type ServiceBooking struct {
    ID          uuid.UUID
    CustomerID  uuid.UUID
    ProviderID  uuid.UUID
    ServiceID   uuid.UUID
    BookingDate time.Time
    Status      string // pending, confirmed, cancelled, completed
    Price       float64
    Notes       string
    CreatedAt   time.Time
}
```

**API Endpoints Needed:**
- `POST /api/v1/customer/bookings/service`
- `GET /api/v1/customer/bookings/:id`
- `PUT /api/v1/customer/bookings/:id/cancel`
- `PUT /api/v1/provider/bookings/:id/accept`
- `PUT /api/v1/provider/bookings/:id/reject`
- `PUT /api/v1/provider/bookings/:id/complete`

#### 4. **Chat/Messaging System**
**Priority: MEDIUM**

**Required Features:**
- Customer-provider messaging
- Real-time chat (WebSocket/SSE)
- Message history
- Read receipts
- File attachments

**Database Models Needed:**
```go
type Conversation struct {
    ID         uuid.UUID
    CustomerID uuid.UUID
    ProviderID uuid.UUID
    CreatedAt  time.Time
}

type Message struct {
    ID             uuid.UUID
    ConversationID uuid.UUID
    SenderID       uint
    Content        string
    AttachmentURL  string
    IsRead         bool
    CreatedAt      time.Time
}
```

**API Endpoints Needed:**
- `GET /api/v1/customer/conversations`
- `GET /api/v1/customer/conversations/:id/messages`
- `POST /api/v1/customer/conversations/:id/messages`
- `GET /api/v1/provider/conversations`
- `GET /api/v1/provider/conversations/:id/messages`
- `POST /api/v1/provider/conversations/:id/messages`
- `WebSocket /api/v1/chat/:conversation_id`

#### 5. **Callback Request System**
**Priority: MEDIUM**

**Required Features:**
- Request callback from provider
- Schedule callback time
- Callback status tracking
- Provider callback management

**Database Models Needed:**
```go
type CallbackRequest struct {
    ID          uuid.UUID
    CustomerID  uuid.UUID
    ProviderID  uuid.UUID
    ServiceID   uuid.UUID
    RequestedAt time.Time
    ScheduledAt time.Time
    Status      string // pending, scheduled, completed, cancelled
    Notes       string
}
```

**API Endpoints Needed:**
- `POST /api/v1/customer/callbacks`
- `GET /api/v1/customer/callbacks`
- `PUT /api/v1/customer/callbacks/:id/cancel`
- `GET /api/v1/provider/callbacks`
- `PUT /api/v1/provider/callbacks/:id/complete`

#### 6. **Provider Earnings & Analytics**
**Priority: MEDIUM**

**Required Features:**
- Earnings tracking
- Revenue breakdown
- Booking statistics
- Performance metrics
- Payout management

**Database Models Needed:**
```go
type ProviderEarnings struct {
    ID          uuid.UUID
    ProviderID  uuid.UUID
    Amount      float64
    BookingID   uuid.UUID
    Commission  float64
    NetAmount   float64
    Status      string // pending, paid
    CreatedAt   time.Time
}

type ProviderStats struct {
    ProviderID      uuid.UUID
    TotalBookings   int
    CompletedBookings int
    TotalRevenue    float64
    AverageRating   float64
    ResponseTime    float64
}
```

**API Endpoints Needed:**
- `GET /api/v1/provider/earnings`
- `GET /api/v1/provider/earnings/:id`
- `GET /api/v1/provider/analytics`
- `GET /api/v1/provider/analytics/dashboard`
- `POST /api/v1/provider/payouts/request`

#### 7. **Provider Availability Management**
**Priority: MEDIUM**

**Required Features:**
- Set availability hours
- Block/unblock dates
- Real-time availability status
- Calendar integration

**Database Models Needed:**
```go
type ProviderAvailability struct {
    ID         uuid.UUID
    ProviderID uuid.UUID
    DayOfWeek  string
    OpenTime   string
    CloseTime  string
    IsAvailable bool
}

type ProviderBlockedDate struct {
    ID         uuid.UUID
    ProviderID uuid.UUID
    Date       time.Time
    Reason     string
}
```

**API Endpoints Needed:**
- `GET /api/v1/provider/availability`
- `PUT /api/v1/provider/availability`
- `POST /api/v1/provider/blocked-dates`
- `DELETE /api/v1/provider/blocked-dates/:id`
- `GET /api/v1/public/providers/:id/availability`

#### 8. **Admin Verification Workflow**
**Priority: MEDIUM**

**Required Features:**
- Verify provider documents
- Approve/reject providers
- Verification status management
- Admin notes on verification

**Database Models Needed:**
```go
type AdminVerification struct {
    ID          uuid.UUID
    ProviderID  uuid.UUID
    AdminID     uint
    Status      string // pending, approved, rejected
    Notes       string
    VerifiedAt  time.Time
}
```

**API Endpoints Needed:**
- `GET /api/v1/admin/providers/pending`
- `GET /api/v1/admin/providers/:id/verification`
- `PUT /api/v1/admin/providers/:id/verify`
- `PUT /api/v1/admin/providers/:id/reject`
- `GET /api/v1/admin/verifications`

#### 9. **Admin Analytics Dashboard**
**Priority: LOW**

**Required Features:**
- Platform statistics
- User growth metrics
- Booking analytics
- Revenue tracking
- Provider performance

**API Endpoints Needed:**
- `GET /api/v1/admin/analytics/overview`
- `GET /api/v1/admin/analytics/users`
- `GET /api/v1/admin/analytics/bookings`
- `GET /api/v1/admin/analytics/revenue`
- `GET /api/v1/admin/analytics/providers`

#### 10. **Social Features**
**Priority: LOW**

**Required Features:**
- Service reviews
- Provider ratings
- Customer testimonials
- Like/favorite services

**Database Models Needed:**
```go
type ServiceReview struct {
    ID          uuid.UUID
    ServiceID   uuid.UUID
    CustomerID  uuid.UUID
    Rating      int
    Comment     string
    Images      []string
    CreatedAt   time.Time
}

type ProviderRating struct {
    ID          uuid.UUID
    ProviderID  uuid.UUID
    CustomerID  uuid.UUID
    Rating      int
    Categories  map[string]int // service_quality, communication, value
    Comment     string
    CreatedAt   time.Time
}
```

**API Endpoints Needed:**
- `GET /api/v1/public/services/:id/reviews`
- `POST /api/v1/customer/services/:id/reviews`
- `GET /api/v1/public/providers/:id/reviews`
- `POST /api/v1/customer/providers/:id/reviews`
- `PUT /api/v1/customer/reviews/:id`

---

## Implementation Priority

### **Phase 1: Core Service Discovery (HIGH)**
1. Service search & discovery endpoints
2. Service categories management
3. Nearby services search
4. Provider service catalog (CRUD)

### **Phase 2: Service Booking System (HIGH)**
1. Service booking creation
2. Booking status management
3. Provider booking acceptance/rejection
4. Booking cancellation

### **Phase 3: Provider Features (MEDIUM)**
1. Provider availability management
2. Provider earnings tracking
3. Provider analytics dashboard
4. Service reviews & ratings

### **Phase 4: Communication (MEDIUM)**
1. Chat/messaging system
2. Callback request system
3. Real-time notifications for bookings

### **Phase 5: Admin Features (MEDIUM)**
1. Admin verification workflow
2. Admin analytics dashboard
3. Platform moderation tools

### **Phase 6: Social Features (LOW)**
1. Advanced reviews system
2. Social sharing
3. Recommendations

---

## Database Schema Updates Needed

### **New Tables Required:**
1. `service_categories`
2. `provider_services` (extend existing)
3. `service_reviews`
4. `provider_ratings`
5. `service_bookings`
6. `conversations`
7. `messages`
8. `callback_requests`
9. `provider_earnings`
10. `provider_availability`
11. `provider_blocked_dates`
12. `admin_verifications`

### **Existing Tables to Extend:**
1. `customers` - add more profile fields
2. `service_providers` - add analytics fields
3. `reservations` - add service booking fields

---

## Frontend Integration Points

### **Mobile App Routes:**
- `/search` - Service discovery
- `/services/:id` - Service details
- `/providers/:id` - Provider profile
- `/bookings` - Customer bookings
- `/messages` - Chat
- `/profile` - Customer profile
- `/saved` - Saved services

### **Provider Dashboard:**
- `/provider/services` - Service management
- `/provider/bookings` - Booking management
- `/provider/earnings` - Earnings
- `/provider/availability` - Availability
- `/provider/analytics` - Analytics

### **Admin Panel:**
- `/admin/verifications` - Provider verification
- `/admin/analytics` - Platform analytics
- `/admin/providers` - Provider management
- `/admin/bookings` - Booking oversight

---

## Summary

**Completed:** 6 major modules (Auth, Customer Identity, Provider Onboarding, Hotel Booking, Payments, Notifications)

**Remaining:** 10 major modules (Service Discovery, Provider Services, Service Booking, Chat, Callbacks, Earnings, Availability, Admin Verification, Admin Analytics, Social Features)

**Estimated Completion:** 
- Phase 1-2 (Core): 2-3 weeks
- Phase 3-4 (Enhanced): 2-3 weeks
- Phase 5-6 (Advanced): 1-2 weeks

**Total Estimated Time:** 5-8 weeks for full API integration
