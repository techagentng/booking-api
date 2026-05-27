# Frontend-Backend API Integration Plan

## Frontend Status Overview

**Fully Implemented (3/10):**
- ✅ Service Discovery & Search
- ✅ Provider Service Management  
- ✅ Service Booking System

**Partially Implemented (4/10):**
- ⚠️ Provider Earnings & Analytics
- ⚠️ Admin Verification Workflow
- ⚠️ Admin Analytics Dashboard
- ⚠️ Social Features

**Missing (3/10):**
- ❌ Chat/Messaging System
- ❌ Callback Request System
- ❌ Provider Availability Management

---

## Integration Plan by Module

### 1. **Authentication & Identity Integration**

#### **Frontend Pages:**
- Login pages
- Signup pages
- Google OAuth flow

#### **Backend API Endpoints:**
- `POST /api/v1/auth/signup` - Manual signup
- `POST /api/v1/auth/login` - Manual login
- `POST /api/v1/google/user/login` - Google OAuth
- `POST /api/v1/auth/logout` - Logout

#### **Integration Tasks:**
- [ ] Create API client for auth endpoints
- [ ] Store JWT tokens (access_token, refresh_token)
- [ ] Implement token refresh logic
- [ ] Handle auth cookies (iwe_access_token)
- [ ] Create auth context/provider for frontend
- [ ] Handle role-based redirects after login:
  - Customer → `/customer/saved` or `/customer/bookings`
  - Provider → `/provider/dashboard` or `/provider/onboarding`
  - Admin → `/admin/dashboard`

#### **Customer Identity Integration:**
- [ ] After login, call `GET /api/v1/customer/me` to get customer profile
- [ ] Auto-create customer profile if not exists (handled by backend)
- [ ] Store customer_id in frontend state

---

### 2. **Service Discovery & Search** ✅ (FULLY IMPLEMENTED)

#### **Frontend Pages:**
- `pages/index.tsx` - Main landing with search
- `[id].tsx` - Provider detail pages
- `index.tsx` - Destinations page

#### **Backend API Endpoints Needed:**
- `GET /api/v1/public/services/search?q=query&category=cat&location=loc`
- `GET /api/v1/public/services/categories`
- `GET /api/v1/public/services/nearby?lat=lat&lng=lng&radius=km`
- `GET /api/v1/public/services/:id`
- `GET /api/v1/public/providers/:id/services`

#### **Integration Tasks:**
- [ ] Connect existing search hooks to backend API
- [ ] Map category filters to API parameters
- [ ] Implement location-based search (geolocation API)
- [ ] Cache search results for performance
- [ ] Handle empty states and loading states

---

### 3. **Provider Service Management** ✅ (FULLY IMPLEMENTED)

#### **Frontend Pages:**
- `index.tsx` - Provider services management
- `pages/onboarding/services.tsx` - Service setup

#### **Backend API Endpoints:**
- `POST /api/v1/provider/onboarding/services` - Create services (during onboarding)
- `GET /api/v1/provider/services` - Get provider's services
- `POST /api/v1/provider/services` - Add new service
- `PUT /api/v1/provider/services/:id` - Update service
- `DELETE /api/v1/provider/services/:id` - Delete service
- `PATCH /api/v1/provider/services/:id/availability` - Toggle availability

#### **Integration Tasks:**
- [ ] Connect onboarding service creation to `POST /api/v1/provider/onboarding/services`
- [ ] Connect service management page to CRUD endpoints
- [ ] Implement image upload for service images
- [ ] Handle service features/tags
- [ ] Add validation for service data

---

### 4. **Service Booking System** ✅ (FULLY IMPLEMENTED)

#### **Frontend Pages:**
- `new.tsx` - New booking flow
- `pages/customer/bookings.tsx` - Customer bookings
- `bookings.tsx` - Provider bookings
- `bookings.tsx` - Admin bookings
- `pages/hotelbooking.tsx` - Hotel bookings
- `pages/hall-bookings.tsx` - Hall bookings

#### **Backend API Endpoints:**
- `GET /api/v1/customer/bookings` - Get customer bookings
- `GET /api/v1/customer/bookings?status=upcoming|completed|cancelled`
- `POST /api/v1/customer/bookings/service` - Book a service (NEW)
- `PUT /api/v1/customer/bookings/:id/cancel` - Cancel booking
- `PUT /api/v1/provider/bookings/:id/accept` - Accept booking (NEW)
- `PUT /api/v1/provider/bookings/:id/reject` - Reject booking (NEW)
- `PUT /api/v1/provider/bookings/:id/complete` - Complete booking (NEW)

#### **Integration Tasks:**
- [ ] Connect customer bookings page to `GET /api/v1/customer/bookings`
- [ ] Implement status filtering (upcoming, completed, cancelled)
- [ ] Create service booking endpoint in backend
- [ ] Connect new booking flow to service booking API
- [ ] Implement provider booking management (accept/reject/complete)
- [ ] Add booking status updates via SSE notifications
- [ ] Handle booking cancellation logic

---

### 5. **Provider Onboarding** ✅ (FULLY IMPLEMENTED)

#### **Frontend Pages:**
- `pages/onboarding/services.tsx` - Service setup
- `pages/onboarding/verification.tsx` - Verification status

#### **Backend API Endpoints:**
- `POST /api/v1/public/onboarding/register` - Provider registration
- `POST /api/v1/public/onboarding/verify-email` - Email verification
- `PUT /api/v1/provider/onboarding/business-info` - Business info
- `POST /api/v1/provider/onboarding/services` - Service creation
- `GET /api/v1/provider/onboarding/verification` - Verification status
- `POST /api/v1/provider/onboarding/training/:module_id` - Training completion
- `POST /api/v1/provider/onboarding/activate` - Account activation
- `GET /api/v1/provider/onboarding/status` - Onboarding status

#### **Integration Tasks:**
- [ ] Create provider registration page
- [ ] Implement email verification flow
- [ ] Connect business info form to API
- [ ] Connect service creation to onboarding
- [ ] Show verification status from API
- [ ] Implement training module completion
- [ ] Handle account activation
- [ ] Redirect to dashboard after activation

---

### 6. **Customer Identity Features** ✅ (FULLY IMPLEMENTED)

#### **Frontend Pages:**
- `pages/customer/saved.tsx` - Saved services
- `pages/customer/bookings.tsx` - Customer bookings
- `pages/customer/profile.tsx` - Customer profile (needs creation)

#### **Backend API Endpoints:**
- `GET /api/v1/customer/me` - Get customer profile
- `PATCH /api/v1/customer/me` - Update customer profile
- `GET /api/v1/customer/saved` - Get saved services
- `POST /api/v1/customer/saved` - Save a service
- `DELETE /api/v1/customer/saved/:id` - Remove saved service
- `GET /api/v1/customer/bookings` - Get customer bookings
- `GET /api/v1/customer/preferences` - Get customer preferences
- `PATCH /api/v1/customer/preferences` - Update customer preferences

#### **Integration Tasks:**
- [ ] Create customer profile page
- [ ] Connect saved services page to API
- [ ] Implement save/unsave functionality on service pages
- [ ] Connect customer bookings to API
- [ ] Implement customer preferences page
- [ ] Add notification preferences
- [ ] Implement auth redirect for customer routes:
  - `/customer/saved` → redirect to `/login?redirectTo=/customer/saved` if not logged in
  - `/customer/bookings` → redirect to `/login?redirectTo=/customer/bookings` if not logged in
  - `/customer/profile` → redirect to `/login?redirectTo=/customer/profile` if not logged in

---

### 7. **Provider Earnings & Analytics** ⚠️ (PARTIALLY IMPLEMENTED)

#### **Frontend Pages:**
- `pages/dashboard.tsx` - General dashboard
- `index.tsx` - Provider dashboard

#### **Backend API Endpoints Needed:**
- `GET /api/v1/provider/earnings` - Get earnings history
- `GET /api/v1/provider/earnings/:id` - Get specific earning
- `GET /api/v1/provider/analytics` - Get analytics data
- `GET /api/v1/provider/analytics/dashboard` - Dashboard analytics
- `POST /api/v1/provider/payouts/request` - Request payout

#### **Integration Tasks:**
- [ ] Create earnings API endpoints in backend
- [ ] Create analytics API endpoints in backend
- [ ] Connect provider dashboard to analytics API
- [ ] Implement earnings breakdown by service
- [ ] Add payout request functionality
- [ ] Create earnings history page
- [ ] Add revenue charts and graphs

---

### 8. **Admin Verification Workflow** ⚠️ (PARTIALLY IMPLEMENTED)

#### **Frontend Pages:**
- `[id].tsx` - Admin provider details
- `index.tsx` - Admin providers list
- `pages/onboarding/verification.tsx` - Verification status

#### **Backend API Endpoints Needed:**
- `GET /api/v1/admin/providers/pending` - Get pending verifications
- `GET /api/v1/admin/providers/:id/verification` - Get verification details
- `PUT /api/v1/admin/providers/:id/verify` - Approve provider
- `PUT /api/v1/admin/providers/:id/reject` - Reject provider
- `GET /api/v1/admin/verifications` - Get all verifications

#### **Integration Tasks:**
- [ ] Create admin verification API endpoints in backend
- [ ] Connect admin providers list to pending verifications
- [ ] Implement verification detail view
- [ ] Add approve/reject functionality
- [ ] Show verification status on provider profile
- [ ] Add admin notes to verification

---

### 9. **Admin Analytics Dashboard** ⚠️ (PARTIALLY IMPLEMENTED)

#### **Frontend Pages:**
- `pages/dashboard.tsx` - Basic dashboard
- `bookings.tsx` - Admin bookings view

#### **Backend API Endpoints Needed:**
- `GET /api/v1/admin/analytics/overview` - Platform overview
- `GET /api/v1/admin/analytics/users` - User analytics
- `GET /api/v1/admin/analytics/bookings` - Booking analytics
- `GET /api/v1/admin/analytics/revenue` - Revenue analytics
- `GET /api/v1/admin/analytics/providers` - Provider analytics

#### **Integration Tasks:**
- [ ] Create admin analytics API endpoints in backend
- [ ] Connect admin dashboard to analytics API
- [ ] Implement user growth charts
- [ ] Add booking statistics
- [ ] Create revenue tracking
- [ ] Add provider performance metrics

---

### 10. **Social Features** ⚠️ (PARTIALLY IMPLEMENTED)

#### **Frontend Pages:**
- `pages/customer/saved.tsx` - Saved/favorite services

#### **Backend API Endpoints Needed:**
- `GET /api/v1/public/services/:id/reviews` - Get service reviews
- `POST /api/v1/customer/services/:id/reviews` - Add service review
- `GET /api/v1/public/providers/:id/reviews` - Get provider reviews
- `POST /api/v1/customer/providers/:id/reviews` - Add provider review
- `PUT /api/v1/customer/reviews/:id` - Update review

#### **Integration Tasks:**
- [ ] Create reviews API endpoints in backend
- [ ] Create reviews database models
- [ ] Add review components to service pages
- [ ] Implement rating system (stars)
- [ ] Add review filtering/sorting
- [ ] Connect saved services to API (already done)
- [ ] Add social sharing buttons

---

### 11. **Chat/Messaging System** ❌ (MISSING)

#### **Frontend Pages:**
- `pages/example-chat.tsx` - Demo/template only

#### **Backend API Endpoints Needed:**
- `GET /api/v1/customer/conversations` - Get customer conversations
- `GET /api/v1/customer/conversations/:id/messages` - Get messages
- `POST /api/v1/customer/conversations/:id/messages` - Send message
- `GET /api/v1/provider/conversations` - Get provider conversations
- `GET /api/v1/provider/conversations/:id/messages` - Get messages
- `POST /api/v1/provider/conversations/:id/messages` - Send message
- `WebSocket /api/v1/chat/:conversation_id` - Real-time chat

#### **Integration Tasks:**
- [ ] Create chat API endpoints in backend
- [ ] Create conversations/messages database models
- [ ] Implement WebSocket server for real-time chat
- [ ] Build chat UI components
- [ ] Implement message history
- [ ] Add read receipts
- [ ] Handle file attachments
- [ ] Add typing indicators
- [ ] Implement push notifications for new messages

---

### 12. **Callback Request System** ❌ (MISSING)

#### **Frontend Pages:**
- None (needs creation)

#### **Backend API Endpoints Needed:**
- `POST /api/v1/customer/callbacks` - Request callback
- `GET /api/v1/customer/callbacks` - Get customer callbacks
- `PUT /api/v1/customer/callbacks/:id/cancel` - Cancel callback
- `GET /api/v1/provider/callbacks` - Get provider callbacks
- `PUT /api/v1/provider/callbacks/:id/complete` - Complete callback

#### **Integration Tasks:**
- [ ] Create callback API endpoints in backend
- [ ] Create callback request database model
- [ ] Build callback request UI on service pages
- [ ] Implement callback scheduling
- [ ] Add callback status tracking
- [ ] Create provider callback management page
- [ ] Add callback notifications

---

### 13. **Provider Availability Management** ❌ (MISSING)

#### **Frontend Pages:**
- None (needs creation)

#### **Backend API Endpoints Needed:**
- `GET /api/v1/provider/availability` - Get availability
- `PUT /api/v1/provider/availability` - Update availability
- `POST /api/v1/provider/blocked-dates` - Block date
- `DELETE /api/v1/provider/blocked-dates/:id` - Unblock date
- `GET /api/v1/public/providers/:id/availability` - Get provider availability

#### **Integration Tasks:**
- [ ] Create availability API endpoints in backend
- [ ] Create availability database models
- [ ] Build availability management UI
- [ ] Implement weekly schedule editor
- [ ] Add date blocking functionality
- [ ] Show availability on service pages
- [ ] Integrate with booking system

---

## Implementation Priority

### **Phase 1: Core Integration (Week 1-2)**
1. **Authentication Integration**
   - Connect login/signup to backend
   - Implement auth context
   - Handle role-based redirects

2. **Customer Identity**
   - Connect customer profile
   - Connect saved services
   - Connect customer bookings
   - Implement auth redirects

3. **Service Discovery**
   - Connect search to backend API
   - Implement category filters
   - Add location-based search

### **Phase 2: Provider Features (Week 2-3)**
1. **Provider Onboarding**
   - Connect registration flow
   - Connect verification
   - Connect service creation

2. **Service Management**
   - Connect CRUD operations
   - Add image upload
   - Implement availability toggle

3. **Service Booking**
   - Create service booking API
   - Connect booking flow
   - Implement status management

### **Phase 3: Enhanced Features (Week 3-4)**
1. **Provider Earnings**
   - Create earnings API
   - Connect dashboard
   - Add payout requests

2. **Admin Verification**
   - Create verification API
   - Connect admin panel
   - Implement approval workflow

3. **Social Features**
   - Create reviews API
   - Add review components
   - Implement ratings

### **Phase 4: Advanced Features (Week 4-5)**
1. **Chat System**
   - Create chat API
   - Implement WebSocket
   - Build chat UI

2. **Callback System**
   - Create callback API
   - Build callback UI
   - Add notifications

3. **Availability Management**
   - Create availability API
   - Build availability UI
   - Integrate with bookings

### **Phase 5: Admin Analytics (Week 5-6)**
1. **Admin Dashboard**
   - Create analytics API
   - Connect admin dashboard
   - Add charts and metrics

---

## API Client Structure

### **Recommended API Client Setup:**

```typescript
// lib/api/client.ts
const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
  const token = getAuthToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Refresh token logic
    }
    return Promise.reject(error);
  }
);
```

### **API Service Modules:**

```typescript
// lib/api/auth.ts
export const authAPI = {
  signup: (data) => apiClient.post('/auth/signup', data),
  login: (data) => apiClient.post('/auth/login', data),
  googleLogin: (data) => apiClient.post('/google/user/login', data),
  logout: () => apiClient.post('/auth/logout'),
};

// lib/api/customer.ts
export const customerAPI = {
  getProfile: () => apiClient.get('/customer/me'),
  updateProfile: (data) => apiClient.patch('/customer/me', data),
  getSaved: () => apiClient.get('/customer/saved'),
  saveService: (data) => apiClient.post('/customer/saved', data),
  removeSaved: (id) => apiClient.delete(`/customer/saved/${id}`),
  getBookings: (status) => apiClient.get('/customer/bookings', { params: { status } }),
  getPreferences: () => apiClient.get('/customer/preferences'),
  updatePreferences: (data) => apiClient.patch('/customer/preferences', data),
};

// lib/api/provider.ts
export const providerAPI = {
  register: (data) => apiClient.post('/public/onboarding/register', data),
  verifyEmail: (data) => apiClient.post('/public/onboarding/verify-email', data),
  updateBusinessInfo: (data) => apiClient.put('/provider/onboarding/business-info', data),
  createServices: (data) => apiClient.post('/provider/onboarding/services', data),
  getVerificationStatus: () => apiClient.get('/provider/onboarding/verification'),
  completeTraining: (moduleId) => apiClient.post(`/provider/onboarding/training/${moduleId}`),
  activateAccount: () => apiClient.post('/provider/onboarding/activate'),
  getOnboardingStatus: () => apiClient.get('/provider/onboarding/status'),
};

// lib/api/services.ts
export const servicesAPI = {
  search: (params) => apiClient.get('/public/services/search', { params }),
  getCategories: () => apiClient.get('/public/services/categories'),
  getNearby: (params) => apiClient.get('/public/services/nearby', { params }),
  getService: (id) => apiClient.get(`/public/services/${id}`),
  getProviderServices: (id) => apiClient.get(`/public/providers/${id}/services`),
};
```

---

## State Management

### **Auth Context:**
```typescript
// context/AuthContext.tsx
interface AuthContextType {
  user: User | null;
  customer: Customer | null;
  provider: Provider | null;
  login: (credentials) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}
```

### **Customer Context:**
```typescript
// context/CustomerContext.tsx
interface CustomerContextType {
  savedServices: SavedService[];
  bookings: Booking[];
  preferences: Preferences;
  saveService: (service) => Promise<void>;
  removeSaved: (id) => Promise<void>;
  updatePreferences: (prefs) => Promise<void>;
}
```

---

## Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## Summary

**Total Tasks:** 80+ integration tasks across 13 modules

**Estimated Timeline:** 5-6 weeks for full integration

**Immediate Priorities:**
1. Authentication integration
2. Customer identity features
3. Service discovery connection
4. Provider onboarding connection

**Key Decisions Needed:**
1. State management approach (Context API vs Redux vs Zustand)
2. WebSocket implementation for chat
3. File upload service for images/documents
4. Notification system integration
