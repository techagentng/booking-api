# Hotel Management System - API Documentation

**Base URL:** `http://localhost:8080/api/v1`

**Authentication:** All endpoints (except `/auth/*` and `/health`) require a Bearer token in the Authorization header.

```
Authorization: Bearer <token>
```

---

## Table of Contents

1. [Authentication](#authentication)
2. [Guests](#guests)
3. [Rooms](#rooms)
4. [Reservations](#reservations)
5. [Room Service](#room-service)
6. [Guest Services (Tablet)](#guest-services-tablet)
7. [Service Requests (Admin)](#service-requests-admin)
8. [Hotel Info](#hotel-info)

---

## Authentication

### Login
```http
POST /auth/login
```
**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
**Response:**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "role": "User"
    }
  }
}
```

### Register
```http
POST /auth/register
```
**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

### Logout
```http
POST /auth/logout
```

---

## Guests

### List Guests
```http
GET /guests?page=1&page_size=20
```
**Response:**
```json
{
  "message": "Guests retrieved",
  "data": {
    "guests": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100
    }
  }
}
```

### Get Guest by ID
```http
GET /guests/:id
```

### Create Guest
```http
POST /guests
```
**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "nationality": "USA",
  "id_type": "passport",
  "id_number": "AB123456"
}
```

### Update Guest
```http
PUT /guests/:id
```

### Delete Guest
```http
DELETE /guests/:id
```

### Get Guest History
```http
GET /guests/:id/history
```
**Response:**
```json
{
  "data": {
    "total_stays": 5,
    "total_spent": 2500.00,
    "average_stay": 3,
    "last_visit": "2024-12-01"
  }
}
```

### Get Guest Preferences
```http
GET /guests/:id/preferences
```

### Get Guest AI Insights
```http
GET /guests/:id/ai-insights
```

---

## Rooms

### List Rooms
```http
GET /rooms?page=1&page_size=20&status=available&type=Deluxe
```
**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number (default: 1) |
| page_size | int | Items per page (default: 20) |
| status | string | Filter by status: `available`, `occupied`, `maintenance`, `cleaning` |
| type | string | Filter by room type |

### Get Available Rooms
```http
GET /rooms/available
```

### Get Room Statistics
```http
GET /rooms/stats
```
**Response:**
```json
{
  "data": [
    { "room_type": "Standard", "total": 20, "available": 15, "occupied": 5 },
    { "room_type": "Deluxe", "total": 15, "available": 10, "occupied": 5 }
  ]
}
```

### Get Rooms by Type
```http
GET /rooms/type/:type
```

### Get Room by ID
```http
GET /rooms/:id
```

### Create Room
```http
POST /rooms
```
**Request:**
```json
{
  "room_number": "301",
  "room_type": "Suite",
  "floor": 3,
  "price_per_night": 250.00,
  "capacity": 4,
  "amenities": ["wifi", "minibar", "jacuzzi"],
  "description": "Luxury suite with city view"
}
```

### Update Room
```http
PUT /rooms/:id
```

### Update Room Status
```http
PUT /rooms/:id/status
```
**Request:**
```json
{
  "status": "maintenance"
}
```

### Check Room Availability
```http
POST /rooms/availability
```
**Request:**
```json
{
  "room_id": 5,
  "check_in_date": "2024-12-10",
  "check_out_date": "2024-12-15"
}
```

### Get Guest by Room Number (Tablet)
```http
GET /rooms/guest/:room_number
```
**Response:**
```json
{
  "data": {
    "guest": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "phone": "+1234567890"
    },
    "reservation": {
      "id": 1,
      "confirmation_number": "BK-2024-001",
      "check_in_date": "2024-12-05",
      "check_out_date": "2024-12-08",
      "number_of_nights": 3
    },
    "room": {
      "id": 5,
      "room_number": "201",
      "room_type": "Deluxe",
      "floor": 2
    },
    "is_checked_in": true
  }
}
```

---

## Reservations

### List Reservations
```http
GET /reservations?page=1&page_size=20&status=confirmed
```
**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number |
| page_size | int | Items per page |
| status | string | `pending`, `confirmed`, `checked-in`, `checked-out`, `cancelled` |
| guest_id | int | Filter by guest |
| room_id | int | Filter by room |
| check_in_date | string | Filter by check-in date (YYYY-MM-DD) |

### Get Reservation by ID
```http
GET /reservations/:id
```

### Create Reservation
```http
POST /reservations
```
**Request:**
```json
{
  "guest_id": 1,
  "room_id": 5,
  "check_in_date": "2024-12-10",
  "check_out_date": "2024-12-15",
  "number_of_guests": 2,
  "total_amount": 750.00,
  "special_requests": "Late check-in",
  "payment_method": "credit_card"
}
```

### Update Reservation
```http
PUT /reservations/:id
```

### Update Reservation Status
```http
PUT /reservations/:id/status
```
**Request:**
```json
{
  "status": "confirmed"
}
```

### Check In
```http
POST /reservations/:id/checkin
```
**Request:**
```json
{
  "actual_check_in_time": "2024-12-10T14:30:00Z",
  "deposit_collected": 100.00,
  "notes": "Guest requested extra pillows"
}
```

### Check Out
```http
POST /reservations/:id/checkout
```

### Express Checkout (Tablet)
```http
POST /reservations/:id/checkout/express
```
**Request:**
```json
{
  "guest_id": 1,
  "email_receipt": true,
  "feedback": "Great stay!"
}
```
**Response:**
```json
{
  "data": {
    "reservation_id": 1,
    "checkout_status": "pending_review",
    "total_charges": 450.00,
    "additional_charges": 75.00,
    "final_total": 525.00,
    "receipt_sent": true,
    "estimated_processing_time": "15 minutes"
  }
}
```

### Cancel Reservation
```http
PUT /reservations/:id/cancel
```

### Get Check-In Statistics
```http
GET /reservations/checkin/stats?date=2024-12-05
```

### Get Payment Details
```http
GET /reservations/:id/payment
```

---

## Room Service

### Menu

#### List Menu Items
```http
GET /room-service/menu?category=Main Course&available=true
```
**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| category | string | Filter by category |
| available | bool | Filter by availability |

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Club Sandwich",
      "description": "Triple-decker with turkey, bacon, lettuce",
      "category": "Main Course",
      "price": 15.00,
      "preparation_time": 15,
      "available": true,
      "image_url": "/images/club-sandwich.jpg",
      "allergens": ["gluten", "dairy"],
      "dietary_info": ["contains-meat"]
    }
  ]
}
```

#### Get Menu Categories
```http
GET /room-service/menu/categories
```
**Response:**
```json
{
  "data": [
    { "name": "Breakfast", "count": 5, "icon": "coffee" },
    { "name": "Main Course", "count": 8, "icon": "utensils" },
    { "name": "Desserts", "count": 4, "icon": "cake" },
    { "name": "Beverages", "count": 6, "icon": "glass-water" }
  ]
}
```

#### Get Menu Item by ID
```http
GET /room-service/menu/:id
```

#### Create Menu Item
```http
POST /room-service/menu
```
**Request:**
```json
{
  "name": "Grilled Salmon",
  "description": "Atlantic salmon with seasonal vegetables",
  "category": "Main Course",
  "price": 32.00,
  "preparation_time": 25,
  "available": true,
  "allergens": ["fish"],
  "dietary_info": ["gluten-free"]
}
```

#### Update Menu Item
```http
PUT /room-service/menu/:id
```

#### Delete Menu Item
```http
DELETE /room-service/menu/:id
```

### Orders

#### List Orders
```http
GET /room-service/orders?page=1&page_size=20&status=pending
```
**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number |
| page_size | int | Items per page |
| status | string | `pending`, `preparing`, `ready`, `delivering`, `delivered`, `cancelled` |
| room_id | int | Filter by room |

#### Get Active Orders
```http
GET /room-service/orders/active
```

#### Get Guest Active Orders (Tablet)
```http
GET /room-service/orders/guest/:guest_id/active
```
**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "order_number": "RS-2024-001",
      "status": "preparing",
      "items": [
        { "name": "Club Sandwich", "quantity": 2 },
        { "name": "Coffee", "quantity": 2 }
      ],
      "total_amount": 38.00,
      "ordered_at": "2024-12-05T14:30:00Z",
      "estimated_delivery_time": "2024-12-05T15:00:00Z",
      "time_elapsed": 15
    }
  ]
}
```

#### Get Order by ID
```http
GET /room-service/orders/:id
```

#### Get Order Status (Tablet)
```http
GET /room-service/orders/:id/status
```
**Response:**
```json
{
  "data": {
    "id": 1,
    "order_number": "RS-2024-001",
    "status": "preparing",
    "status_history": [
      { "status": "pending", "timestamp": "2024-12-05T14:30:00Z" },
      { "status": "preparing", "timestamp": "2024-12-05T14:35:00Z", "updated_by": "Chef John" }
    ],
    "estimated_delivery_time": "2024-12-05T15:00:00Z",
    "current_step": 2,
    "total_steps": 5
  }
}
```

#### Create Order
```http
POST /room-service/orders
```
**Request:**
```json
{
  "room_id": 5,
  "guest_id": 1,
  "items": [
    { "menu_item_id": 1, "quantity": 2, "special_instructions": "No onions" },
    { "menu_item_id": 5, "quantity": 1 }
  ],
  "special_instructions": "Please deliver to room door",
  "delivery_time": "2024-12-05T15:00:00Z"
}
```
**Response:**
```json
{
  "data": {
    "id": 1,
    "order_number": "RS-2024-001",
    "status": "pending",
    "items": [...],
    "subtotal": 47.00,
    "tax": 4.70,
    "service_charge": 7.05,
    "total_amount": 58.75,
    "estimated_delivery_time": "2024-12-05T15:00:00Z"
  }
}
```

#### Update Order Status
```http
PUT /room-service/orders/:id/status
```
**Request:**
```json
{
  "status": "preparing",
  "prepared_by": "Chef John"
}
```

#### Deliver Order
```http
POST /room-service/orders/:id/deliver
```
**Request:**
```json
{
  "delivered_by": "Staff Member"
}
```

#### Cancel Order
```http
PUT /room-service/orders/:id/cancel
```
**Request:**
```json
{
  "cancellation_reason": "Guest changed their mind"
}
```

#### Get Room Service Statistics
```http
GET /room-service/stats
```

---

## Guest Services (Tablet)

### Request Housekeeping
```http
POST /services/housekeeping
```
**Request:**
```json
{
  "room_id": 5,
  "guest_id": 1,
  "service_type": "cleaning",
  "priority": "normal",
  "notes": "Please clean bathroom thoroughly",
  "preferred_time": "2024-12-05T16:00:00Z"
}
```
**Service Types:** `cleaning`, `towels`, `amenities`, `bedding`, `turndown`

**Response:**
```json
{
  "data": {
    "id": 1,
    "request_number": "HK-2024-001",
    "type": "housekeeping",
    "service_type": "cleaning",
    "status": "pending",
    "estimated_time": "2024-12-05T16:30:00Z"
  }
}
```

### Request Maintenance
```http
POST /services/maintenance
```
**Request:**
```json
{
  "room_id": 5,
  "guest_id": 1,
  "issue_type": "air_conditioning",
  "priority": "high",
  "description": "AC is not cooling properly",
  "preferred_time": "2024-12-05T17:00:00Z"
}
```
**Issue Types:** `air_conditioning`, `plumbing`, `electrical`, `appliances`, `furniture`, `other`

**Response:**
```json
{
  "data": {
    "id": 2,
    "request_number": "MN-2024-001",
    "type": "maintenance",
    "issue_type": "air_conditioning",
    "status": "pending",
    "priority": "high",
    "estimated_response_time": "2024-12-05T16:00:00Z"
  }
}
```

### Get Guest Service Requests
```http
GET /services/guest/:guest_id?status=active
```
**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | `active`, `completed`, `cancelled`, or specific status |

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "request_number": "HK-2024-001",
      "type": "housekeeping",
      "service_type": "cleaning",
      "status": "in_progress",
      "priority": "normal",
      "requested_at": "2024-12-05T14:00:00Z",
      "estimated_time": "2024-12-05T14:30:00Z"
    }
  ]
}
```

---

## Service Requests (Admin)

### List All Service Requests
```http
GET /service-requests?page=1&page_size=20&status=pending&type=housekeeping
```
**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number (default: 1) |
| page_size | int | Items per page (default: 20, max: 100) |
| status | string | `pending`, `in_progress`, `completed`, `cancelled` |
| type | string | `housekeeping`, `maintenance` |

**Response:**
```json
{
  "data": {
    "data": [
      {
        "id": 1,
        "request_number": "HK-2024-001",
        "room_id": 5,
        "guest_id": 1,
        "type": "housekeeping",
        "service_type": "cleaning",
        "status": "pending",
        "priority": "normal",
        "notes": "Please clean bathroom",
        "requested_at": "2024-12-05T14:00:00Z",
        "room": { "room_number": "201", "floor": 2 },
        "guest": { "name": "John Doe" }
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

---

## Hotel Info

### Get Hotel Information
```http
GET /hotel/info
```
**Response:**
```json
{
  "data": {
    "name": "Grand Hotel",
    "description": "Luxury accommodation in the heart of the city",
    "facilities": [
      { "name": "Swimming Pool", "hours": "6:00 AM - 10:00 PM", "location": "Rooftop", "icon": "pool" },
      { "name": "Fitness Center", "hours": "24/7", "location": "Ground Floor", "icon": "dumbbell" },
      { "name": "Spa", "hours": "9:00 AM - 9:00 PM", "location": "2nd Floor", "icon": "spa" }
    ],
    "dining": [
      { "name": "Main Restaurant", "cuisine": "International", "hours": "7:00 AM - 11:00 PM", "location": "Ground Floor" },
      { "name": "Rooftop Bar", "type": "Bar & Lounge", "hours": "5:00 PM - 1:00 AM", "location": "Rooftop" }
    ],
    "contact": {
      "reception": "ext. 0",
      "room_service": "ext. 1",
      "concierge": "ext. 2",
      "emergency": "ext. 911"
    },
    "wifi": {
      "network": "GrandHotel-Guest",
      "password": "welcome2024"
    }
  }
}
```

---

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## Response Format

All responses follow this structure:

```json
{
  "message": "Description of the result",
  "data": { ... },
  "error": null
}
```

On error:
```json
{
  "message": "Error description",
  "data": null,
  "error": "Detailed error message"
}
```

---

## Pagination

Paginated endpoints return:

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## Frontend Integration Tips

### React Query Example

```typescript
// hooks/useRooms.ts
import { useQuery } from '@tanstack/react-query';
import { api } from '../lib/api';

export const useRooms = (params?: { status?: string; type?: string }) => {
  return useQuery({
    queryKey: ['rooms', params],
    queryFn: () => api.get('/rooms', { params }),
  });
};

export const useRoomGuest = (roomNumber: string) => {
  return useQuery({
    queryKey: ['room-guest', roomNumber],
    queryFn: () => api.get(`/rooms/guest/${roomNumber}`),
    enabled: !!roomNumber,
  });
};
```

### API Client Setup

```typescript
// lib/api.ts
import axios from 'axios';

export const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

## Test Data

The server seeds test data on startup:

| Guest | Room | Status |
|-------|------|--------|
| John Doe (john.doe@example.com) | First room | Confirmed |
| Jane Smith (jane.smith@example.com) | Second room | Checked-In |

**Menu Items:** 21 items across 5 categories (Breakfast, Main Course, Salads, Desserts, Beverages)
