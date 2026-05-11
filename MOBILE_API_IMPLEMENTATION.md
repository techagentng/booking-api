# Mobile-First Public API Implementation

## Overview
Public APIs for the TripsBook mobile-first customer interface. No authentication required - focus on discovery, search, and basic booking functionality.

## 🏗️ API Structure

### Base URL
```
https://api.tripsbook.com/api/v1/public/
```

### Response Format
```json
{
  "success": true,
  "data": {},
  "message": "Operation successful",
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "location": "Lagos, Nigeria"
  }
}
```

## 📱 Core Public Endpoints

### 1. Categories & Services

#### Get All Categories
```http
GET /api/v1/public/categories
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "all",
      "name": "All",
      "icon": "home",
      "color": "bg-blue-500",
      "count": 1250
    },
    {
      "id": "hotels",
      "name": "Hotels",
      "icon": "building",
      "color": "bg-purple-500",
      "count": 450
    },
    {
      "id": "transport",
      "name": "Transport",
      "icon": "car",
      "color": "bg-green-500",
      "count": 320
    },
    {
      "id": "food",
      "name": "Food",
      "icon": "utensils",
      "color": "bg-orange-500",
      "count": 280
    },
    {
      "id": "shopping",
      "name": "Shopping",
      "icon": "shopping-bag",
      "color": "bg-pink-500",
      "count": 200
    }
  ]
}
```

#### Get Services by Category
```http
GET /api/v1/public/services/category/{category}?location={city}&page={page}&limit={limit}
```

**Parameters:**
- `category`: all, hotels, transport, food, shopping
- `location`: Lagos, Abuja, etc.
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "svc_123",
      "name": "Transcorp Hilton Hotel",
      "type": "hotel",
      "description": "Luxury hotel in the heart of Abuja",
      "image": "https://cdn.tripsbook.com/images/hotels/transcorp.jpg",
      "rating": 4.5,
      "price": "$$$",
      "location": {
        "address": "1 Constitution Avenue",
        "city": "Abuja",
        "state": "FCT",
        "coordinates": {
          "lat": 9.0579,
          "lng": 7.4951
        }
      },
      "contact": {
        "phone": "+234-800-000-0000",
        "email": "abuja@transcorphilton.com"
      },
      "operatingHours": {
        "open": "00:00",
        "close": "23:59",
        "days": ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
      },
      "features": ["WiFi", "Pool", "Gym", "Spa", "Restaurant"],
      "distance": "2.3 km"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 450,
      "pages": 23
    }
  }
}
```

### 2. Explore Tab Endpoints

#### Featured Services
```http
GET /api/v1/public/explore/featured?location={city}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "title": "Hotels",
      "description": "Find perfect accommodation",
      "icon": "building",
      "color": "bg-purple-500",
      "count": 45,
      "image": "https://cdn.tripsbook.com/images/categories/hotels.jpg"
    },
    {
      "title": "Airport Transfers",
      "description": "Reliable airport pickup",
      "icon": "plane",
      "color": "bg-blue-500",
      "count": 23,
      "image": "https://cdn.tripsbook.com/images/categories/airport.jpg"
    }
  ]
}
```

#### Popular Destinations
```http
GET /api/v1/public/explore/destinations
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "name": "Lagos",
      "country": "Nigeria",
      "rating": 4.7,
      "distance": "0 km",
      "image": "https://cdn.tripsbook.com/images/destinations/lagos.jpg",
      "serviceCount": 1250
    },
    {
      "name": "Abuja",
      "country": "Nigeria",
      "rating": 4.6,
      "distance": "500 km",
      "image": "https://cdn.tripsbook.com/images/destinations/abuja.jpg",
      "serviceCount": 850
    }
  ]
}
```

### 3. Nearby Tab Endpoints

#### Nearby Services
```http
GET /api/v1/public/services/nearby?lat={lat}&lng={lng}&radius={km}&category={type}
```

**Parameters:**
- `lat`: Latitude (required if no city)
- `lng`: Longitude (required if no city)
- `radius`: Radius in km (1, 5, 10, default: 5)
- `category`: Filter by category (optional)
- `location`: City name (alternative to coordinates)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "svc_456",
      "name": "Eko Hotels & Suites",
      "type": "Hotel",
      "distance": "1.2 km",
      "rating": 4.4,
      "price": "$$$",
      "image": "https://cdn.tripsbook.com/images/hotels/eko.jpg",
      "coordinates": {
        "lat": 6.4474,
        "lng": 3.3903
      },
      "openNow": true,
      "featured": true
    },
    {
      "id": "svc_789",
      "name": "Bolt Driver - John",
      "type": "Transport",
      "distance": "0.8 km",
      "rating": 4.8,
      "price": "₦500/km",
      "image": "https://cdn.tripsbook.com/images/drivers/john.jpg",
      "coordinates": {
        "lat": 6.4474,
        "lng": 3.3903
      },
      "availableNow": true
    }
  ]
}
```

#### Distance Filter Options
```http
GET /api/v1/public/services/nearby/filters
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "label": "< 1 km",
      "value": 1,
      "count": 45
    },
    {
      "label": "< 5 km",
      "value": 5,
      "count": 128
    },
    {
      "label": "< 10 km",
      "value": 10,
      "count": 234
    },
    {
      "label": "Any distance",
      "value": null,
      "count": 1250
    }
  ]
}
```

### 4. Trending Tab Endpoints

#### Trending Services
```http
GET /api/v1/public/services/trending?period={period}&location={city}
```

**Parameters:**
- `period`: week, month (default: week)
- `location`: City name (optional)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "svc_111",
      "name": "Federal Palace Hotel",
      "type": "Hotel",
      "distance": "3.1 km",
      "rating": 4.6,
      "price": "$$$",
      "trending": true,
      "image": "https://cdn.tripsbook.com/images/hotels/federal.jpg",
      "weeklyChange": 23.5,
      "trendingBadge": "🔥 Hot"
    },
    {
      "id": "svc_222",
      "name": "Uber Premium",
      "type": "Transport",
      "distance": "1.5 km",
      "rating": 4.7,
      "price": "₦800/km",
      "trending": true,
      "image": "https://cdn.tripsbook.com/images/transport/uber.jpg",
      "weeklyChange": 18.2,
      "trendingBadge": "📈 Rising"
    }
  ]
}
```

#### Trending Categories
```http
GET /api/v1/public/categories/trending?period={week|month}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "name": "Hotels",
      "change": "+23%",
      "color": "bg-purple-500",
      "icon": "building"
    },
    {
      "name": "Transport",
      "change": "+18%",
      "color": "bg-green-500",
      "icon": "car"
    },
    {
      "name": "Food",
      "change": "+12%",
      "color": "bg-orange-500",
      "icon": "utensils"
    }
  ]
}
```

### 5. Search Functionality

#### Global Search
```http
GET /api/v1/public/search?q={query}&location={city}&category={type}&lat={lat}&lng={lng}
```

**Parameters:**
- `q`: Search query (required)
- `location`: City name (optional)
- `category`: Filter by category (optional)
- `lat`: Latitude for proximity (optional)
- `lng`: Longitude for proximity (optional)

**Response:**
```json
{
  "success": true,
  "data": {
    "services": [
      {
        "id": "svc_333",
        "name": "Transcorp Hilton",
        "type": "Hotel",
        "distance": "2.3 km",
        "rating": 4.5,
        "price": "$$$$",
        "image": "https://cdn.tripsbook.com/images/hotels/transcorp.jpg",
        "matchScore": 95
      }
    ],
    "suggestions": [
      "Transcorp Hilton Abuja",
      "Transcorp Hotels Lagos",
      "Hilton Hotels Nigeria"
    ],
    "categories": [
      {
        "name": "Hotels",
        "count": 45
      }
    ]
  }
}
```

#### Search Suggestions (Autocomplete)
```http
GET /api/v1/public/search/suggestions?q={partial}&location={city}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "text": "Transcorp Hilton Abuja",
      "type": "service",
      "category": "Hotel"
    },
    {
      "text": "Hotels in Abuja",
      "type": "category",
      "category": "Hotel"
    }
  ]
}
```

### 6. Location Services

#### Current Location
```http
GET /api/v1/public/location/current
```

**Response:**
```json
{
  "success": true,
  "data": {
    "city": "Lagos",
    "state": "Lagos State",
    "country": "Nigeria",
    "coordinates": {
      "lat": 6.4474,
      "lng": 3.3903
    },
    "timezone": "Africa/Lagos"
  }
}
```

#### Update Location
```http
POST /api/v1/public/location/update
```

**Request Body:**
```json
{
  "city": "Abuja",
  "coordinates": {
    "lat": 9.0579,
    "lng": 7.4951
  }
}
```

#### Popular Locations
```http
GET /api/v1/public/locations/popular
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "city": "Lagos",
      "state": "Lagos State",
      "country": "Nigeria",
      "serviceCount": 1250,
      "image": "https://cdn.tripsbook.com/images/cities/lagos.jpg"
    },
    {
      "city": "Abuja",
      "state": "FCT",
      "country": "Nigeria",
      "serviceCount": 850,
      "image": "https://cdn.tripsbook.com/images/cities/abuja.jpg"
    }
  ]
}
```

### 7. Service Details

#### Get Service Details
```http
GET /api/v1/public/services/{id}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "svc_123",
    "name": "Transcorp Hilton Hotel",
    "type": "hotel",
    "description": "Luxury hotel in the heart of Abuja with world-class amenities",
    "images": [
      "https://cdn.tripsbook.com/images/hotels/transcorp1.jpg",
      "https://cdn.tripsbook.com/images/hotels/transcorp2.jpg",
      "https://cdn.tripsbook.com/images/hotels/transcorp3.jpg"
    ],
    "rating": 4.5,
    "reviewCount": 342,
    "price": "$$$$",
    "location": {
      "address": "1 Constitution Avenue",
      "city": "Abuja",
      "state": "FCT",
      "coordinates": {
        "lat": 9.0579,
        "lng": 7.4951
      }
    },
    "contact": {
      "phone": "+234-800-000-0000",
      "email": "abuja@transcorphilton.com",
      "website": "https://transcorphilton.com"
    },
    "operatingHours": {
      "open": "00:00",
      "close": "23:59",
      "days": ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"],
      "openNow": true
    },
    "features": ["WiFi", "Pool", "Gym", "Spa", "Restaurant", "Bar", "Business Center"],
    "specificData": {
      "starRating": 5,
      "checkIn": "15:00",
      "checkOut": "11:00",
      "rooms": [
        {
          "id": "room_1",
          "type": "Deluxe Room",
          "price": 45000,
          "capacity": 2,
          "available": true,
          "images": ["https://cdn.tripsbook.com/images/rooms/deluxe1.jpg"]
        }
      ],
      "amenities": ["24/7 Front Desk", "Room Service", "Laundry", "Concierge"]
    },
    "reviews": [
      {
        "id": "rev_1",
        "userName": "John Doe",
        "rating": 5,
        "comment": "Excellent service and beautiful rooms!",
        "date": "2024-01-15",
        "helpful": 23
      }
    ],
    "distance": "2.3 km",
    "nearbyServices": [
      {
        "id": "svc_456",
        "name": "Federal Secretariat",
        "type": "Government",
        "distance": "0.5 km"
      }
    ]
  }
}
```

### 8. Public Booking (Optional - No Auth)

#### Create Quick Booking Request
```http
POST /api/v1/public/bookings/quick
```

**Request Body:**
```json
{
  "serviceId": "svc_123",
  "customerInfo": {
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+234-800-000-0000"
  },
  "bookingDetails": {
    "date": "2024-01-20",
    "time": "15:00",
    "duration": 120,
    "notes": "Business meeting"
  },
  "location": {
    "address": "Customer pickup location",
    "coordinates": {
      "lat": 9.0579,
      "lng": 7.4951
    }
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "bookingId": "bk_123456",
    "status": "pending",
    "estimatedPrice": 45000,
    "nextSteps": [
      "Wait for provider confirmation",
      "You'll receive a confirmation call/email",
      "Payment will be required upon confirmation"
    ]
  }
}
```

#### Check Booking Status
```http
GET /api/v1/public/bookings/{bookingId}/status
```

**Response:**
```json
{
  "success": true,
  "data": {
    "bookingId": "bk_123456",
    "status": "confirmed",
    "providerName": "Transcorp Hilton Hotel",
    "confirmedAt": "2024-01-15T10:30:00Z",
    "totalAmount": 45000,
    "paymentRequired": true,
    "paymentLink": "https://pay.tripsbook.com/bk_123456"
  }
}
```

## 🗺️ Map Integration Endpoints

#### Map View Services
```http
GET /api/v1/public/services/map?bounds={sw_lat,sw_lng,ne_lat,ne_lng}&category={type}
```

**Parameters:**
- `bounds`: Southwest and northeast coordinates
- `category`: Filter by category (optional)

**Response:**
```json
{
  "success": true,
  "data": {
    "markers": [
      {
        "id": "svc_123",
        "name": "Transcorp Hilton",
        "type": "hotel",
        "coordinates": {
          "lat": 9.0579,
          "lng": 7.4951
        },
        "price": "$$$$",
        "rating": 4.5,
        "icon": "hotel"
      }
    ],
    "clusters": [
      {
        "id": "cluster_1",
        "coordinates": {
          "lat": 9.0580,
          "lng": 7.4952
        },
        "count": 15,
        "type": "hotel"
      }
    ]
  }
}
```

## 📊 Analytics & Tracking

#### Track User Interaction
```http
POST /api/v1/public/analytics/track
```

**Request Body:**
```json
{
  "event": "service_view",
  "data": {
    "serviceId": "svc_123",
    "category": "hotel",
    "source": "search",
    "location": "Lagos",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "userFingerprint": "fp_123456789"
}
```

## 🚀 Implementation Plan

### Phase 1: Core Discovery APIs (Week 1)
- [ ] Categories endpoint
- [ ] Services by category
- [ ] Basic search functionality
- [ ] Location services

### Phase 2: Location-Based Features (Week 2)
- [ ] Nearby services with geolocation
- [ ] Distance filtering
- [ ] Map integration endpoints
- [ ] Popular destinations

### Phase 3: Advanced Features (Week 3)
- [ ] Trending algorithms
- [ ] Search suggestions
- [ ] Service details
- [ ] Analytics tracking

### Phase 4: Booking Integration (Week 4)
- [ ] Quick booking requests
- [ ] Booking status tracking
- [ ] Payment integration
- [ ] Provider notifications

## 🔧 Technical Implementation

### Go Models
```go
// tripsbook/models/public_api.go
type PublicService struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Description string    `json:"description"`
    Image       string    `json:"image"`
    Rating      float64   `json:"rating"`
    Price       string    `json:"price"`
    Distance    string    `json:"distance,omitempty"`
    Location    Location  `json:"location"`
    Contact     Contact   `json:"contact"`
    OperatingHours *OperatingHours `json:"operatingHours,omitempty"`
    Features    []string  `json:"features,omitempty"`
    OpenNow     bool      `json:"openNow,omitempty"`
    AvailableNow bool     `json:"availableNow,omitempty"`
    Featured    bool      `json:"featured,omitempty"`
}

type Category struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Icon  string `json:"icon"`
    Color string `json:"color"`
    Count int    `json:"count"`
}

type TrendingService struct {
    PublicService
    Trending      bool    `json:"trending"`
    WeeklyChange  float64 `json:"weeklyChange"`
    TrendingBadge string  `json:"trendingBadge"`
}

type SearchResponse struct {
    Services    []PublicService `json:"services"`
    Suggestions []string        `json:"suggestions"`
    Categories  []CategoryCount `json:"categories"`
}
```

### Performance Optimizations
- **Redis caching** for popular queries
- **Database indexes** on location and category
- **Image CDN** for fast loading
- **Response compression** for mobile networks
- **Rate limiting** to prevent abuse

### Error Handling
```go
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// Common error codes
const (
    ErrInvalidLocation   = "INVALID_LOCATION"
    ErrServiceNotFound   = "SERVICE_NOT_FOUND"
    ErrSearchTooBroad    = "SEARCH_TOO_BROAD"
    ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
)
```

This mobile-first public API design provides:
1. **No authentication barrier** for easy customer access
2. **Location-based discovery** for nearby services
3. **Trending algorithms** for popular recommendations
4. **Fast, optimized responses** for mobile performance
5. **Comprehensive search** with autocomplete
6. **Map integration** for visual discovery
7. **Optional booking** without complex auth flows
