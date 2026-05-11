# TripsBook Frontend API Consumption Guide

## 📱 Overview
This guide provides comprehensive documentation for frontend developers to consume the TripsBook mobile-first public APIs. All endpoints are **public (no authentication required)** and optimized for mobile performance.

## 🌐 Base URL
```
Production: https://api.tripsbook.com/api/v1/public/
Development: http://localhost:8080/api/v1/public/
```

## 📋 Standard Response Format

### Success Response
```json
{
  "success": true,
  "data": {},
  "message": "Operation successful",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria",
    "request_id": "req_123456789"
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "INVALID_LOCATION",
    "message": "Invalid location provided",
    "details": "Coordinates are out of service area"
  },
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "req_123456789"
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [],
  "message": "Data retrieved successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "pages": 8,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

## 🏷️ Categories & Discovery

### 1. Get All Categories
```http
GET /api/v1/public/categories
```

**Response Example:**
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
  ],
  "message": "Categories retrieved successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria"
  }
}
```

**Frontend Implementation:**
```typescript
interface Category {
  id: string;
  name: string;
  icon: string;
  color: string;
  count: number;
}

async function getCategories(): Promise<Category[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/categories`);
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to fetch categories:', error);
    throw error;
  }
}
```

### 2. Get Services by Category
```http
GET /api/v1/public/services/category/{category}?location={city}&page={page}&limit={limit}
```

**Example Request:**
```
GET /api/v1/public/services/category/hotels?location=Lagos&page=1&limit=20
```

**Response Example:**
```json
{
  "success": true,
  "data": [
    {
      "id": "svc_123",
      "name": "Transcorp Hilton Hotel",
      "type": "hotel",
      "description": "Luxury hotel in the heart of Abuja with world-class amenities",
      "image": "https://cdn.tripsbook.com/images/hotels/transcorp.jpg",
      "rating": 4.5,
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
      "features": ["WiFi", "Pool", "Gym", "Spa", "Restaurant"]
    }
  ],
  "message": "Services retrieved successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria",
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 450,
      "pages": 23,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

**Frontend Implementation:**
```typescript
interface Service {
  id: string;
  name: string;
  type: string;
  description: string;
  image: string;
  rating: number;
  price: string;
  location: {
    address: string;
    city: string;
    state: string;
    coordinates: { lat: number; lng: number };
  };
  contact: {
    phone?: string;
    email?: string;
    website?: string;
  };
  features: string[];
}

async function getServicesByCategory(
  category: string,
  location: string = 'Lagos',
  page: number = 1,
  limit: number = 20
): Promise<{ services: Service[]; pagination: any }> {
  try {
    const params = new URLSearchParams({
      location,
      page: page.toString(),
      limit: limit.toString()
    });
    
    const response = await fetch(`${API_BASE_URL}/services/category/${category}?${params}`);
    const result = await response.json();
    
    if (result.success) {
      return {
        services: result.data,
        pagination: result.meta.pagination
      };
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to fetch services:', error);
    throw error;
  }
}
```

## 🗺️ Explore Tab

### 1. Get Featured Services
```http
GET /api/v1/public/explore/featured?location={city}
```

**Response Example:**
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
    },
    {
      "title": "Transport",
      "description": "Ride services & rentals",
      "icon": "car",
      "color": "bg-green-500",
      "count": 67,
      "image": "https://cdn.tripsbook.com/images/categories/transport.jpg"
    },
    {
      "title": "Food & Dining",
      "description": "Restaurants & catering",
      "icon": "utensils",
      "color": "bg-orange-500",
      "count": 89,
      "image": "https://cdn.tripsbook.com/images/categories/food.jpg"
    }
  ],
  "message": "Featured services retrieved successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria"
  }
}
```

### 2. Get Popular Destinations
```http
GET /api/v1/public/explore/destinations
```

**Response Example:**
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
      "service_count": 1250
    },
    {
      "name": "Abuja",
      "country": "Nigeria",
      "rating": 4.6,
      "distance": "500 km",
      "image": "https://cdn.tripsbook.com/images/destinations/abuja.jpg",
      "service_count": 850
    },
    {
      "name": "Port Harcourt",
      "country": "Nigeria",
      "rating": 4.5,
      "distance": "650 km",
      "image": "https://cdn.tripsbook.com/images/destinations/portharcourt.jpg",
      "service_count": 420
    }
  ],
  "message": "Popular destinations retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface FeaturedService {
  title: string;
  description: string;
  icon: string;
  color: string;
  count: number;
  image: string;
}

interface PopularDestination {
  name: string;
  country: string;
  rating: number;
  distance: string;
  image: string;
  service_count: number;
}

async function getExploreData(location: string = 'Lagos') {
  try {
    const [featuredResponse, destinationsResponse] = await Promise.all([
      fetch(`${API_BASE_URL}/explore/featured?location=${location}`),
      fetch(`${API_BASE_URL}/explore/destinations`)
    ]);

    const featured = await featuredResponse.json();
    const destinations = await destinationsResponse.json();

    if (featured.success && destinations.success) {
      return {
        featured: featured.data,
        destinations: destinations.data
      };
    }
  } catch (error) {
    console.error('Failed to fetch explore data:', error);
    throw error;
  }
}
```

## 📍 Nearby Tab

### 1. Get Nearby Services
```http
GET /api/v1/public/services/nearby?lat={lat}&lng={lng}&radius={km}&category={type}
```

**Example Request:**
```
GET /api/v1/public/services/nearby?lat=6.4474&lng=3.3903&radius=5&category=hotels
```

**Response Example:**
```json
{
  "success": true,
  "data": [
    {
      "id": "svc_456",
      "name": "Eko Hotels & Suites",
      "type": "Hotel",
      "description": "Luxury beachfront hotel with stunning ocean views",
      "image": "https://cdn.tripsbook.com/images/hotels/eko.jpg",
      "rating": 4.4,
      "price": "$$$",
      "location": {
        "address": "1415 Adetokunbo Ademola Street",
        "city": "Lagos",
        "state": "Lagos State",
        "coordinates": {
          "lat": 6.4474,
          "lng": 3.3903
        }
      },
      "contact": {
        "phone": "+234-1-2778000",
        "email": "eko.suites@ekohotels.com",
        "website": "https://ekohotels.com"
      },
      "features": ["WiFi", "Pool", "Beach Access", "Spa", "Restaurant"],
      "distance": "1.2 km",
      "open_now": true,
      "available_now": true
    },
    {
      "id": "svc_789",
      "name": "Bolt Driver - John",
      "type": "Transport",
      "description": "Professional driver with 5+ years experience",
      "image": "https://cdn.tripsbook.com/images/drivers/john.jpg",
      "rating": 4.8,
      "price": "₦500/km",
      "location": {
        "address": "Victoria Island",
        "city": "Lagos",
        "state": "Lagos State",
        "coordinates": {
          "lat": 6.4474,
          "lng": 3.3903
        }
      },
      "contact": {
        "phone": "+234-800-000-0001"
      },
      "features": ["Air Conditioning", "GPS", "Experienced Driver"],
      "distance": "0.8 km",
      "open_now": true,
      "available_now": true
    }
  ],
  "message": "Nearby services retrieved successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria"
  }
}
```

### 2. Get Distance Filters
```http
GET /api/v1/public/services/nearby/filters?lat={lat}&lng={lng}&category={type}
```

**Response Example:**
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
      "value": 0,
      "count": 1250
    }
  ],
  "message": "Distance filters retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface NearbyService extends Service {
  distance: string;
  open_now: boolean;
  available_now: boolean;
}

interface DistanceFilter {
  label: string;
  value: number;
  count: number;
}

async function getNearbyServices(
  lat: number,
  lng: number,
  radius: number = 5,
  category?: string
): Promise<NearbyService[]> {
  try {
    const params = new URLSearchParams({
      lat: lat.toString(),
      lng: lng.toString(),
      radius: radius.toString()
    });
    
    if (category) {
      params.append('category', category);
    }
    
    const response = await fetch(`${API_BASE_URL}/services/nearby?${params}`);
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to fetch nearby services:', error);
    throw error;
  }
}

// Get user's current location
async function getCurrentLocation(): Promise<{ lat: number; lng: number }> {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('Geolocation is not supported'));
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        resolve({
          lat: position.coords.latitude,
          lng: position.coords.longitude
        });
      },
      (error) => {
        reject(error);
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000 // 5 minutes
      }
    );
  });
}

// Usage example
async function loadNearbyServices() {
  try {
    const location = await getCurrentLocation();
    const services = await getNearbyServices(location.lat, location.lng);
    // Update UI with services
  } catch (error) {
    console.error('Error loading nearby services:', error);
    // Show error message to user
  }
}
```

## 🔥 Trending Tab

### 1. Get Trending Services
```http
GET /api/v1/public/services/trending?period={week|month}&location={city}
```

**Example Request:**
```
GET /api/v1/public/services/trending?period=week&location=Lagos
```

**Response Example:**
```json
{
  "success": true,
  "data": [
    {
      "id": "svc_111",
      "name": "Federal Palace Hotel",
      "type": "Hotel",
      "description": "Historic luxury hotel with modern amenities",
      "image": "https://cdn.tripsbook.com/images/hotels/federal.jpg",
      "rating": 4.6,
      "price": "$$$$",
      "location": {
        "coordinates": {
          "lat": 6.4474,
          "lng": 3.3903
        }
      },
      "distance": "3.1 km",
      "trending": true,
      "weekly_change": 23.5,
      "trending_badge": "🔥 Hot"
    },
    {
      "id": "svc_222",
      "name": "Uber Premium",
      "type": "Transport",
      "description": "Premium ride service with professional drivers",
      "image": "https://cdn.tripsbook.com/images/transport/uber.jpg",
      "rating": 4.7,
      "price": "₦800/km",
      "location": {
        "coordinates": {
          "lat": 6.4474,
          "lng": 3.3903
        }
      },
      "distance": "1.5 km",
      "trending": true,
      "weekly_change": 18.2,
      "trending_badge": "📈 Rising"
    },
    {
      "id": "svc_333",
      "name": "Terra Kulture",
      "type": "Restaurant",
      "description": "Contemporary Nigerian cuisine with cultural experience",
      "image": "https://cdn.tripsbook.com/images/restaurants/terra.jpg",
      "rating": 4.5,
      "price": "$$$",
      "location": {
        "coordinates": {
          "lat": 6.4474,
          "lng": 3.3903
        }
      },
      "distance": "2.8 km",
      "trending": true,
      "weekly_change": 12.7,
      "trending_badge": "⭐ Popular"
    }
  ],
  "message": "Trending services retrieved successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria"
  }
}
```

### 2. Get Trending Categories
```http
GET /api/v1/public/categories/trending?period={week|month}
```

**Response Example:**
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
    },
    {
      "name": "Shopping",
      "change": "+8%",
      "color": "bg-pink-500",
      "icon": "shopping-bag"
    }
  ],
  "message": "Trending categories retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface TrendingService extends Service {
  distance: string;
  trending: boolean;
  image: string;
  weekly_change: number;
  trending_badge: string;
}

interface TrendingCategory {
  name: string;
  change: string;
  color: string;
  icon: string;
}

async function getTrendingServices(
  period: 'week' | 'month' = 'week',
  location: string = 'Lagos'
): Promise<TrendingService[]> {
  try {
    const params = new URLSearchParams({ period, location });
    const response = await fetch(`${API_BASE_URL}/services/trending?${params}`);
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to fetch trending services:', error);
    throw error;
  }
}

async function getTrendingCategories(period: 'week' | 'month' = 'week'): Promise<TrendingCategory[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/categories/trending?period=${period}`);
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to fetch trending categories:', error);
    throw error;
  }
}
```

## 🔍 Search Functionality

### 1. Global Search
```http
GET /api/v1/public/search?q={query}&location={city}&category={type}&lat={lat}&lng={lng}
```

**Example Request:**
```
GET /api/v1/public/search?q=hotel&location=Lagos&category=hotels
```

**Response Example:**
```json
{
  "success": true,
  "data": {
    "services": [
      {
        "id": "svc_123",
        "name": "Transcorp Hilton",
        "type": "Hotel",
        "distance": "2.3 km",
        "rating": 4.5,
        "price": "$$$$",
        "image": "https://cdn.tripsbook.com/images/hotels/transcorp.jpg",
        "match_score": 95
      },
      {
        "id": "svc_456",
        "name": "Eko Hotels & Suites",
        "type": "Hotel",
        "distance": "1.2 km",
        "rating": 4.4,
        "price": "$$$",
        "image": "https://cdn.tripsbook.com/images/hotels/eko.jpg",
        "match_score": 88
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
  },
  "message": "Search completed successfully",
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "location": "Lagos, Nigeria"
  }
}
```

### 2. Search Suggestions (Autocomplete)
```http
GET /api/v1/public/search/suggestions?q={partial}&location={city}
```

**Example Request:**
```
GET /api/v1/public/search/suggestions?q=trans&location=Lagos
```

**Response Example:**
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
      "text": "Transcorp Hotels Lagos",
      "type": "service",
      "category": "Hotel"
    },
    {
      "text": "Transport Services",
      "type": "category",
      "category": "Transport"
    }
  ],
  "message": "Search suggestions retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface SearchResult {
  services: Service[];
  suggestions: string[];
  categories: { name: string; count: number }[];
}

interface SearchSuggestion {
  text: string;
  type: 'service' | 'category';
  category?: string;
}

// Search with debouncing
function useSearch() {
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<SearchResult | null>(null);
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);

  const searchServices = useCallback(
    debounce(async (query: string, location: string, category?: string) => {
      if (!query.trim()) {
        setResults(null);
        setSuggestions([]);
        return;
      }

      setLoading(true);
      try {
        const params = new URLSearchParams({ q: query, location });
        if (category && category !== 'all') {
          params.append('category', category);
        }

        const response = await fetch(`${API_BASE_URL}/search?${params}`);
        const result = await response.json();

        if (result.success) {
          setResults(result.data);
        }
      } catch (error) {
        console.error('Search failed:', error);
      } finally {
        setLoading(false);
      }
    }, 300),
    []
  );

  const getSuggestions = useCallback(
    debounce(async (query: string, location: string) => {
      if (!query.trim()) {
        setSuggestions([]);
        return;
      }

      try {
        const response = await fetch(`${API_BASE_URL}/search/suggestions?q=${query}&location=${location}`);
        const result = await response.json();

        if (result.success) {
          setSuggestions(result.data);
        }
      } catch (error) {
        console.error('Failed to get suggestions:', error);
      }
    }, 200),
    []
  );

  return { searchServices, getSuggestions, loading, results, suggestions };
}

// Debounce utility
function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}
```

## 📍 Location Services

### 1. Get Current Location
```http
GET /api/v1/public/location/current
```

**Response Example:**
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
  },
  "message": "Current location retrieved successfully"
}
```

### 2. Update Location
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

### 3. Get Popular Locations
```http
GET /api/v1/public/locations/popular
```

**Response Example:**
```json
{
  "success": true,
  "data": [
    {
      "city": "Lagos",
      "state": "Lagos State",
      "country": "Nigeria",
      "service_count": 1250,
      "image": "https://cdn.tripsbook.com/images/cities/lagos.jpg"
    },
    {
      "city": "Abuja",
      "state": "FCT",
      "country": "Nigeria",
      "service_count": 850,
      "image": "https://cdn.tripsbook.com/images/cities/abuja.jpg"
    }
  ],
  "message": "Popular locations retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface UserLocation {
  city: string;
  state: string;
  country: string;
  coordinates: { lat: number; lng: number };
  timezone: string;
}

interface PopularLocation {
  city: string;
  state: string;
  country: string;
  service_count: number;
  image: string;
}

class LocationService {
  private currentLocation: UserLocation | null = null;

  async getCurrentLocation(): Promise<UserLocation> {
    try {
      const response = await fetch(`${API_BASE_URL}/location/current`);
      const result = await response.json();
      
      if (result.success) {
        this.currentLocation = result.data;
        return result.data;
      } else {
        throw new Error(result.error.message);
      }
    } catch (error) {
      console.error('Failed to get current location:', error);
      throw error;
    }
  }

  async updateLocation(city: string, coordinates: { lat: number; lng: number }): Promise<void> {
    try {
      const response = await fetch(`${API_BASE_URL}/location/update`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ city, coordinates }),
      });

      const result = await response.json();
      
      if (result.success) {
        this.currentLocation = { city, coordinates } as UserLocation;
      } else {
        throw new Error(result.error.message);
      }
    } catch (error) {
      console.error('Failed to update location:', error);
      throw error;
    }
  }

  async getPopularLocations(): Promise<PopularLocation[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/locations/popular`);
      const result = await response.json();
      
      if (result.success) {
        return result.data;
      } else {
        throw new Error(result.error.message);
      }
    } catch (error) {
      console.error('Failed to get popular locations:', error);
      throw error;
    }
  }

  async detectUserLocation(): Promise<UserLocation> {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error('Geolocation not supported'));
        return;
      }

      navigator.geolocation.getCurrentPosition(
        async (position) => {
          const { latitude, longitude } = position.coords;
          
          // Reverse geocode to get city name (you'd use a geocoding service)
          const city = await this.reverseGeocode(latitude, longitude);
          
          const location: UserLocation = {
            city,
            state: '', // Would come from geocoding
            country: 'Nigeria',
            coordinates: { lat: latitude, lng: longitude },
            timezone: 'Africa/Lagos'
          };

          this.currentLocation = location;
          resolve(location);
        },
        (error) => {
          reject(error);
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 300000
        }
      );
    });
  }

  private async reverseGeocode(lat: number, lng: number): Promise<string> {
    // Implement reverse geocoding using a service like Google Maps API
    // For now, return a default
    return 'Lagos';
  }
}
```

## 🗺️ Map Integration

### Get Map View Services
```http
GET /api/v1/public/services/map?sw_lat={lat}&sw_lng={lng}&ne_lat={lat}&ne_lng={lng}&category={type}
```

**Example Request:**
```
GET /api/v1/public/services/map?sw_lat=6.4&sw_lng=3.3&ne_lat=6.5&ne_lng=3.4&category=hotels
```

**Response Example:**
```json
{
  "success": true,
  "data": {
    "markers": [
      {
        "id": "svc_123",
        "name": "Transcorp Hilton",
        "type": "hotel",
        "location": {
          "lat": 9.0579,
          "lng": 7.4951
        },
        "price": "$$$$",
        "rating": 4.5,
        "icon": "hotel"
      },
      {
        "id": "svc_456",
        "name": "Eko Hotels",
        "type": "hotel",
        "location": {
          "lat": 6.4474,
          "lng": 3.3903
        },
        "price": "$$$",
        "rating": 4.4,
        "icon": "hotel"
      }
    ],
    "clusters": [
      {
        "id": "cluster_1",
        "location": {
          "lat": 6.4480,
          "lng": 3.3905
        },
        "count": 15,
        "type": "hotel",
        "icon": "hotel"
      }
    ]
  },
  "message": "Map services retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface MapMarker {
  id: string;
  name: string;
  type: string;
  location: { lat: number; lng: number };
  price: string;
  rating: number;
  icon: string;
}

interface MapCluster {
  id: string;
  location: { lat: number; lng: number };
  count: number;
  type: string;
  icon: string;
}

interface MapViewData {
  markers: MapMarker[];
  clusters: MapCluster[];
}

async function getMapViewServices(
  bounds: {
    southwest: { lat: number; lng: number };
    northeast: { lat: number; lng: number };
  },
  category?: string
): Promise<MapViewData> {
  try {
    const params = new URLSearchParams({
      sw_lat: bounds.southwest.lat.toString(),
      sw_lng: bounds.southwest.lng.toString(),
      ne_lat: bounds.northeast.lat.toString(),
      ne_lng: bounds.northeast.lng.toString()
    });

    if (category && category !== 'all') {
      params.append('category', category);
    }

    const response = await fetch(`${API_BASE_URL}/services/map?${params}`);
    const result = await response.json();

    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to get map services:', error);
    throw error;
  }
}

// React Hook for map integration
function useMapServices() {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<MapViewData | null>(null);

  const loadMapServices = useCallback(async (
    bounds: any,
    category?: string
  ) => {
    setLoading(true);
    try {
      const mapData = await getMapViewServices(bounds, category);
      setData(mapData);
    } catch (error) {
      console.error('Failed to load map services:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  return { loadMapServices, loading, data };
}
```

## 📊 Service Details

### Get Service Details
```http
GET /api/v1/public/services/{id}
```

**Example Request:**
```
GET /api/v1/public/services/svc_123
```

**Response Example:**
```json
{
  "success": true,
  "data": {
    "id": "svc_123",
    "name": "Transcorp Hilton Hotel",
    "type": "hotel",
    "description": "Luxury hotel in the heart of Abuja with world-class amenities and exceptional service",
    "images": [
      "https://cdn.tripsbook.com/images/hotels/transcorp1.jpg",
      "https://cdn.tripsbook.com/images/hotels/transcorp2.jpg",
      "https://cdn.tripsbook.com/images/hotels/transcorp3.jpg"
    ],
    "rating": 4.5,
    "review_count": 342,
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
    "operating_hours": {
      "open": "00:00",
      "close": "23:59",
      "days": ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"],
      "open_now": true
    },
    "features": ["WiFi", "Pool", "Gym", "Spa", "Restaurant", "Bar", "Business Center"],
    "specific_data": {
      "star_rating": 5,
      "check_in": "15:00",
      "check_out": "11:00",
      "rooms": [
        {
          "id": "room_1",
          "type": "Deluxe Room",
          "price": 45000,
          "capacity": 2,
          "available": true,
          "images": ["https://cdn.tripsbook.com/images/rooms/deluxe1.jpg"]
        },
        {
          "id": "room_2",
          "type": "Executive Suite",
          "price": 75000,
          "capacity": 3,
          "available": true,
          "images": ["https://cdn.tripsbook.com/images/rooms/suite1.jpg"]
        }
      ],
      "amenities": ["24/7 Front Desk", "Room Service", "Laundry", "Concierge"]
    },
    "reviews": [
      {
        "id": "rev_1",
        "user_name": "John Doe",
        "rating": 5,
        "comment": "Excellent service and beautiful rooms!",
        "date": "2024-01-15",
        "helpful": 23
      },
      {
        "id": "rev_2",
        "user_name": "Jane Smith",
        "rating": 4,
        "comment": "Great location, friendly staff",
        "date": "2024-01-10",
        "helpful": 15
      }
    ],
    "distance": "2.3 km",
    "nearby_services": [
      {
        "id": "svc_456",
        "name": "Federal Secretariat",
        "type": "Government",
        "distance": "0.5 km"
      }
    ]
  },
  "message": "Service details retrieved successfully"
}
```

**Frontend Implementation:**
```typescript
interface ServiceDetails {
  id: string;
  name: string;
  type: string;
  description: string;
  images: string[];
  rating: number;
  review_count: number;
  price: string;
  location: {
    address: string;
    city: string;
    state: string;
    coordinates: { lat: number; lng: number };
  };
  contact: {
    phone?: string;
    email?: string;
    website?: string;
  };
  operating_hours?: {
    open: string;
    close: string;
    days: string[];
    open_now: boolean;
  };
  features: string[];
  specific_data?: any; // Type-specific data
  reviews: Review[];
  distance: string;
  nearby_services: NearbyService[];
}

interface Review {
  id: string;
  user_name: string;
  rating: number;
  comment: string;
  date: string;
  helpful: number;
}

async function getServiceDetails(serviceId: string): Promise<ServiceDetails> {
  try {
    const response = await fetch(`${API_BASE_URL}/services/${serviceId}`);
    const result = await response.json();

    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error.message);
    }
  } catch (error) {
    console.error('Failed to get service details:', error);
    throw error;
  }
}

// React Hook for service details
function useServiceDetails(serviceId: string) {
  const [loading, setLoading] = useState(false);
  const [service, setService] = useState<ServiceDetails | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadServiceDetails = async () => {
      if (!serviceId) return;

      setLoading(true);
      setError(null);

      try {
        const details = await getServiceDetails(serviceId);
        setService(details);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load service details');
      } finally {
        setLoading(false);
      }
    };

    loadServiceDetails();
  }, [serviceId]);

  return { loading, service, error };
}
```

## 📈 Analytics Tracking

### Track User Interactions
```http
POST /api/v1/public/analytics/track
```

**Request Body:**
```json
{
  "event": "service_view",
  "data": {
    "service_id": "svc_123",
    "category": "hotel",
    "source": "search",
    "location": "Lagos",
    "query": "luxury hotel"
  },
  "user_fingerprint": "fp_123456789",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Frontend Implementation:**
```typescript
class AnalyticsService {
  private userFingerprint: string;

  constructor() {
    this.userFingerprint = this.generateFingerprint();
  }

  private generateFingerprint(): string {
    // Generate a unique fingerprint for anonymous users
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    if (ctx) {
      ctx.textBaseline = 'top';
      ctx.font = '14px Arial';
      ctx.fillText('User fingerprint', 2, 2);
    }
    
    const fingerprint = [
      navigator.userAgent,
      navigator.language,
      screen.width + 'x' + screen.height,
      new Date().getTimezoneOffset(),
      canvas?.toDataURL() || ''
    ].join('|');
    
    return btoa(fingerprint).substring(0, 16);
  }

  async track(event: string, data: Record<string, any>): Promise<void> {
    try {
      await fetch(`${API_BASE_URL}/analytics/track`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          event,
          data,
          user_fingerprint: this.userFingerprint,
          timestamp: new Date().toISOString()
        })
      });
    } catch (error) {
      console.error('Failed to track analytics:', error);
    }
  }

  trackServiceView(serviceId: string, category: string, source: string) {
    this.track('service_view', {
      service_id: serviceId,
      category,
      source
    });
  }

  trackSearch(query: string, resultsCount: number, category?: string) {
    this.track('search', {
      query,
      results_count: resultsCount,
      category
    });
  }

  trackCategorySelect(category: string) {
    this.track('category_select', {
      category
    });
  }

  trackLocationChange(location: string) {
    this.track('location_change', {
      location
    });
  }
}

// Usage
const analytics = new AnalyticsService();

// Track service view
analytics.trackServiceView('svc_123', 'hotel', 'search');

// Track search
analytics.trackSearch('luxury hotel', 45, 'hotels');

// Track category selection
analytics.trackCategorySelect('hotels');
```

## 🛠️ Complete Frontend Integration Example

```typescript
// API Service Class
class TripsBookAPI {
  private baseURL: string;
  private locationService: LocationService;
  private analytics: AnalyticsService;

  constructor(baseURL: string = 'https://api.tripsbook.com/api/v1/public') {
    this.baseURL = baseURL;
    this.locationService = new LocationService();
    this.analytics = new AnalyticsService();
  }

  // Categories
  async getCategories(): Promise<Category[]> {
    const response = await fetch(`${this.baseURL}/categories`);
    const result = await response.json();
    
    if (result.success) {
      this.analytics.track('categories_viewed', { count: result.data.length });
      return result.data;
    }
    throw new Error(result.error.message);
  }

  // Services
  async getServicesByCategory(category: string, location: string, page = 1, limit = 20) {
    const params = new URLSearchParams({ location, page: page.toString(), limit: limit.toString() });
    const response = await fetch(`${this.baseURL}/services/category/${category}?${params}`);
    const result = await response.json();
    
    if (result.success) {
      return {
        services: result.data,
        pagination: result.meta.pagination
      };
    }
    throw new Error(result.error.message);
  }

  // Nearby
  async getNearbyServices(lat: number, lng: number, radius = 5, category?: string) {
    const params = new URLSearchParams({
      lat: lat.toString(),
      lng: lng.toString(),
      radius: radius.toString()
    });
    
    if (category) params.append('category', category);
    
    const response = await fetch(`${this.baseURL}/services/nearby?${params}`);
    const result = await response.json();
    
    if (result.success) {
      this.analytics.track('nearby_search', { 
        lat, 
        lng, 
        radius, 
        category,
        results_count: result.data.length 
      });
      return result.data;
    }
    throw new Error(result.error.message);
  }

  // Search
  async search(query: string, location: string, category?: string) {
    const params = new URLSearchParams({ q: query, location });
    if (category && category !== 'all') params.append('category', category);
    
    const response = await fetch(`${this.baseURL}/search?${params}`);
    const result = await response.json();
    
    if (result.success) {
      this.analytics.trackSearch(query, result.data.services.length, category);
      return result.data;
    }
    throw new Error(result.error.message);
  }

  // Service Details
  async getServiceDetails(serviceId: string) {
    const response = await fetch(`${this.baseURL}/services/${serviceId}`);
    const result = await response.json();
    
    if (result.success) {
      this.analytics.trackServiceView(serviceId, result.data.type, 'direct');
      return result.data;
    }
    throw new Error(result.error.message);
  }

  // Location
  async getCurrentLocation() {
    return this.locationService.getCurrentLocation();
  }

  async updateLocation(city: string, coordinates: { lat: number; lng: number }) {
    const result = await this.locationService.updateLocation(city, coordinates);
    this.analytics.trackLocationChange(city);
    return result;
  }

  // Trending
  async getTrendingServices(period: 'week' | 'month' = 'week', location = 'Lagos') {
    const response = await fetch(`${this.baseURL}/services/trending?period=${period}&location=${location}`);
    const result = await response.json();
    
    if (result.success) return result.data;
    throw new Error(result.error.message);
  }
}

// React Context for Global State
interface TripsBookContextType {
  api: TripsBookAPI;
  currentLocation: UserLocation | null;
  categories: Category[];
  loading: boolean;
  error: string | null;
  updateLocation: (location: UserLocation) => void;
}

const TripsBookContext = createContext<TripspsBookContextType | null>(null);

export function TripsBookProvider({ children }: { children: React.ReactNode }) {
  const [api] = useState(() => new TripsBookAPI());
  const [currentLocation, setCurrentLocation] = useState<UserLocation | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const initializeApp = async () => {
      try {
        setLoading(true);
        
        // Load current location
        const location = await api.getCurrentLocation();
        setCurrentLocation(location);
        
        // Load categories
        const categoriesData = await api.getCategories();
        setCategories(categoriesData);
        
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to initialize app');
      } finally {
        setLoading(false);
      }
    };

    initializeApp();
  }, [api]);

  const updateLocation = async (location: UserLocation) => {
    try {
      await api.updateLocation(location.city, location.coordinates);
      setCurrentLocation(location);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update location');
    }
  };

  return (
    <TripsBookContext.Provider value={{
      api,
      currentLocation,
      categories,
      loading,
      error,
      updateLocation
    }}>
      {children}
    </TripsBookContext.Provider>
  );
}

export function useTripsBook() {
  const context = useContext(TripsBookContext);
  if (!context) {
    throw new Error('useTripsBook must be used within TripsBookProvider');
  }
  return context;
}
```

## 🚀 Quick Start Guide

### 1. Installation
```bash
npm install @tripsbook/api-client
# or
yarn add @tripsbook/api-client
```

### 2. Basic Usage
```typescript
import { TripsBookAPI, TripsBookProvider, useTripsBook } from '@tripsbook/api-client';

function App() {
  return (
    <TripsBookProvider>
      <MyApp />
    </TripsBookProvider>
  );
}

function MyApp() {
  const { api, currentLocation, loading } = useTripsBook();

  if (loading) return <div>Loading...</div>;

  const handleSearch = async (query: string) => {
    try {
      const results = await api.search(query, currentLocation?.city || 'Lagos');
      console.log('Search results:', results);
    } catch (error) {
      console.error('Search failed:', error);
    }
  };

  return (
    <div>
      <SearchBar onSearch={handleSearch} />
      {/* Rest of your app */}
    </div>
  );
}
```

### 3. Error Handling
```typescript
try {
  const services = await api.getServicesByCategory('hotels', 'Lagos');
  // Handle success
} catch (error) {
  if (error.message.includes('INVALID_LOCATION')) {
    // Handle location error
  } else if (error.message.includes('SERVICE_NOT_FOUND')) {
    // Handle not found error
  } else {
    // Handle other errors
  }
}
```

## 📱 Mobile Optimization Tips

### 1. Image Loading
```typescript
const ImageWithFallback = ({ src, alt, ...props }) => {
  const [imageSrc, setImageSrc] = useState(src);
  const [loading, setLoading] = useState(true);

  return (
    <div className="relative">
      {loading && <div className="animate-pulse bg-gray-200" />}
      <img
        src={imageSrc}
        alt={alt}
        onLoad={() => setLoading(false)}
        onError={() => setImageSrc('/fallback-image.jpg')}
        loading="lazy"
        {...props}
      />
    </div>
  );
};
```

### 2. Infinite Scroll
```typescript
function useInfiniteScroll(fetchMore: () => Promise<void>) {
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const loadMore = useCallback(async () => {
    if (loading || !hasMore) return;
    
    setLoading(true);
    try {
      await fetchMore();
    } catch (error) {
      console.error('Failed to load more:', error);
    } finally {
      setLoading(false);
    }
  }, [fetchMore, loading, hasMore]);

  return { loadMore, loading, hasMore };
}
```

### 3. Offline Support
```typescript
// Cache API responses for offline use
const cache = new Map();

async function cachedFetch(url: string, options?: RequestInit) {
  const cacheKey = `${url}_${JSON.stringify(options)}`;
  
  if (cache.has(cacheKey)) {
    return cache.get(cacheKey);
  }

  const response = await fetch(url, options);
  const data = await response.clone().json();
  
  cache.set(cacheKey, data);
  return data;
}
```

## 🎯 Best Practices

### 1. Performance
- Use debouncing for search inputs
- Implement lazy loading for images
- Cache API responses appropriately
- Use pagination for large datasets

### 2. User Experience
- Show loading states during API calls
- Provide meaningful error messages
- Implement retry mechanisms for failed requests
- Track user interactions for analytics

### 3. Error Handling
- Always check `success` field in responses
- Handle network errors gracefully
- Show user-friendly error messages
- Implement fallback data when needed

### 4. Security
- Never expose API keys in frontend
- Validate user inputs before sending
- Use HTTPS for all API calls
- Implement rate limiting on the client side

This comprehensive guide provides everything your frontend team needs to successfully integrate with the TripsBook mobile-first APIs, including dummy data, TypeScript interfaces, React hooks, and best practices for mobile optimization.
