# Room API Integration Guide

## Overview

This guide helps frontend developers integrate with the Room Management API endpoints.

**Base URL:** `http://localhost:8080/api/v1`

**Authentication:** All endpoints require `Authorization: Bearer <token>` header.

---

## Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/rooms` | Get all rooms (paginated) |
| GET | `/rooms/:id` | Get room by ID |
| POST | `/rooms` | Create a new room |
| PUT | `/rooms/:id` | Update a room |
| DELETE | `/rooms/:id` | Delete a room |
| GET | `/rooms/available` | Get available rooms for date range |
| POST | `/rooms/availability` | Check room availability (alternative) |
| GET | `/rooms/type/:type` | Get rooms by type |
| GET | `/rooms/stats` | Get room statistics |
| PUT | `/rooms/:id/status` | Update room status |

---

## 1. Get All Rooms

**Endpoint:** `GET /rooms`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | number | 1 | Page number |
| `page_size` | number | 10 | Items per page |
| `search` | string | - | Search by room number or type |
| `status` | string | - | Filter by status (available, occupied, maintenance, cleaning) |

**Example Request:**
```typescript
const response = await axios.get('/rooms', {
  params: {
    page: 1,
    page_size: 10,
    status: 'available'
  }
});
```

**Response:**
```json
{
  "message": "Rooms retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "room_number": "101",
        "room_type": "Standard",
        "floor": 1,
        "capacity": 2,
        "price_per_night": 100.00,
        "status": "available",
        "description": "Cozy standard room with city view",
        "amenities": "",
        "bed_type": "Queen",
        "created_at": "2024-12-05T00:00:00Z",
        "updated_at": "2024-12-05T00:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total": 11,
      "total_pages": 2
    }
  },
  "errors": "",
  "status": "OK"
}
```

**Frontend Usage:**
```typescript
// hooks/useRooms.ts
import { useState, useEffect } from 'react';
import axios from '../lib/axios';

interface Room {
  id: number;
  room_number: string;
  room_type: string;
  floor: number;
  capacity: number;
  price_per_night: number;
  status: string;
  description: string;
  bed_type: string;
}

interface RoomListResponse {
  data: Room[];
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

export function useRooms(page = 1, pageSize = 10, status?: string) {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [meta, setMeta] = useState<RoomListResponse['meta'] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        setLoading(true);
        const { data } = await axios.get('/rooms', {
          params: { page, page_size: pageSize, status }
        });
        setRooms(data.data.data);
        setMeta(data.data.meta);
      } catch (err: any) {
        setError(err.response?.data?.errors || 'Failed to fetch rooms');
      } finally {
        setLoading(false);
      }
    };

    fetchRooms();
  }, [page, pageSize, status]);

  return { rooms, meta, loading, error };
}
```

---

## 2. Get Room by ID

**Endpoint:** `GET /rooms/:id`

**Example Request:**
```typescript
const response = await axios.get('/rooms/1');
```

**Response:**
```json
{
  "message": "Room retrieved successfully",
  "data": {
    "id": 1,
    "room_number": "101",
    "room_type": "Standard",
    "floor": 1,
    "capacity": 2,
    "price_per_night": 100.00,
    "status": "available",
    "description": "Cozy standard room with city view",
    "amenities": "",
    "bed_type": "Queen",
    "created_at": "2024-12-05T00:00:00Z",
    "updated_at": "2024-12-05T00:00:00Z"
  },
  "errors": "",
  "status": "OK"
}
```

---

## 3. Create Room

**Endpoint:** `POST /rooms`

**Request Body:**
```json
{
  "room_number": "501",
  "room_type": "Suite",
  "floor": 5,
  "capacity": 4,
  "price_per_night": 250.00,
  "description": "Luxury suite with ocean view",
  "amenities": "WiFi, Mini Bar, Jacuzzi",
  "bed_type": "King"
}
```

**Required Fields:**
- `room_number` (string) - Unique room identifier
- `room_type` (string) - Must be: "Standard", "Deluxe", or "Suite"
- `floor` (number) - Must be positive integer
- `capacity` (number) - Must be positive integer
- `price_per_night` (number) - Must be greater than 0

**Example Request:**
```typescript
const response = await axios.post('/rooms', {
  room_number: '501',
  room_type: 'Suite',
  floor: 5,
  capacity: 4,
  price_per_night: 250.00,
  description: 'Luxury suite with ocean view',
  bed_type: 'King'
});
```

**Response (201 Created):**
```json
{
  "message": "Room created successfully",
  "data": {
    "id": 12,
    "room_number": "501",
    "room_type": "Suite",
    "floor": 5,
    "capacity": 4,
    "price_per_night": 250.00,
    "status": "available",
    "description": "Luxury suite with ocean view",
    "bed_type": "King",
    "created_at": "2024-12-05T10:00:00Z",
    "updated_at": "2024-12-05T10:00:00Z"
  },
  "errors": "",
  "status": "Created"
}
```

**Error Response (409 Conflict):**
```json
{
  "message": "Room creation failed",
  "data": null,
  "errors": "room number already exists",
  "status": "Conflict"
}
```

---

## 4. Update Room

**Endpoint:** `PUT /rooms/:id`

**Request Body (all fields optional):**
```json
{
  "price_per_night": 275.00,
  "capacity": 5,
  "description": "Updated description"
}
```

**Example Request:**
```typescript
const response = await axios.put('/rooms/12', {
  price_per_night: 275.00,
  capacity: 5
});
```

**Response:**
```json
{
  "message": "Room updated successfully",
  "data": {
    "id": 12,
    "room_number": "501",
    "room_type": "Suite",
    "floor": 5,
    "capacity": 5,
    "price_per_night": 275.00,
    "status": "available",
    ...
  },
  "errors": "",
  "status": "OK"
}
```

---

## 5. Delete Room

**Endpoint:** `DELETE /rooms/:id`

**Example Request:**
```typescript
const response = await axios.delete('/rooms/12');
```

**Response:**
```json
{
  "message": "Room deleted successfully",
  "data": null,
  "errors": "",
  "status": "OK"
}
```

---

## 6. Get Available Rooms (GET - Recommended)

**Endpoint:** `GET /rooms/available`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `check_in` | string | Yes | Check-in date (YYYY-MM-DD) |
| `check_out` | string | Yes | Check-out date (YYYY-MM-DD) |
| `room_type` | string | No | Filter by room type |

**Example Request:**
```typescript
const response = await axios.get('/rooms/available', {
  params: {
    check_in: '2024-12-15',
    check_out: '2024-12-18',
    room_type: 'Deluxe'
  }
});
```

**Response:**
```json
{
  "message": "Available rooms retrieved successfully",
  "data": [
    {
      "id": 5,
      "room_number": "201",
      "room_type": "Deluxe",
      "floor": 2,
      "capacity": 2,
      "price_per_night": 150.00,
      "status": "available",
      "description": "Spacious deluxe room with balcony",
      "bed_type": "King"
    },
    {
      "id": 6,
      "room_number": "202",
      "room_type": "Deluxe",
      "floor": 2,
      "capacity": 3,
      "price_per_night": 150.00,
      "status": "available",
      "description": "Deluxe room with ocean view",
      "bed_type": "King"
    }
  ],
  "errors": "",
  "status": "OK"
}
```

**Frontend Usage for Booking Form:**
```typescript
// In your booking form component
const checkAvailability = async (roomType: string, checkIn: string, checkOut: string) => {
  try {
    setLoading(true);
    const { data } = await axios.get('/rooms/available', {
      params: {
        check_in: checkIn,
        check_out: checkOut,
        room_type: roomType
      }
    });
    
    // Transform for dropdown
    const availableRooms = data.data.map((room: any) => ({
      id: room.id.toString(),
      number: room.room_number,
      type: room.room_type,
      price: room.price_per_night
    }));
    
    setAvailableRooms(availableRooms);
    
    // Auto-select first room
    if (availableRooms.length > 0) {
      setSelectedRoom(availableRooms[0].id);
    }
  } catch (err: any) {
    setError(err.response?.data?.errors || 'Failed to check availability');
  } finally {
    setLoading(false);
  }
};
```

---

## 7. Check Room Availability (POST - Alternative)

**Endpoint:** `POST /rooms/availability`

**Request Body:**
```json
{
  "room_type": "Standard",
  "check_in_date": "2024-12-15",
  "check_out_date": "2024-12-18"
}
```

**Example Request:**
```typescript
const response = await axios.post('/rooms/availability', {
  room_type: 'Standard',
  check_in_date: '2024-12-15',
  check_out_date: '2024-12-18'
});
```

**Response:** Same as GET `/rooms/available`

---

## 8. Update Room Status

**Endpoint:** `PUT /rooms/:id/status`

**Request Body:**
```json
{
  "status": "maintenance"
}
```

**Valid Status Values:**
- `available` - Room is ready for guests
- `occupied` - Guest is currently in the room
- `maintenance` - Room requires maintenance
- `cleaning` - Room is being cleaned

**Example Request:**
```typescript
const response = await axios.put('/rooms/1/status', {
  status: 'maintenance'
});
```

**Response:**
```json
{
  "message": "Room status updated successfully",
  "data": null,
  "errors": "",
  "status": "OK"
}
```

---

## 9. Get Room Statistics

**Endpoint:** `GET /rooms/stats`

**Example Request:**
```typescript
const response = await axios.get('/rooms/stats');
```

**Response:**
```json
{
  "message": "Room statistics retrieved successfully",
  "data": [
    {
      "room_type": "Standard",
      "total_rooms": 4,
      "available": 3,
      "occupied": 1,
      "maintenance": 0,
      "price_per_night": 100.00
    },
    {
      "room_type": "Deluxe",
      "total_rooms": 4,
      "available": 4,
      "occupied": 0,
      "maintenance": 0,
      "price_per_night": 150.00
    },
    {
      "room_type": "Suite",
      "total_rooms": 3,
      "available": 3,
      "occupied": 0,
      "maintenance": 0,
      "price_per_night": 250.00
    }
  ],
  "errors": "",
  "status": "OK"
}
```

---

## 10. Get Rooms by Type

**Endpoint:** `GET /rooms/type/:type`

**Example Request:**
```typescript
const response = await axios.get('/rooms/type/Deluxe');
```

**Response:**
```json
{
  "message": "Rooms retrieved successfully",
  "data": [
    {
      "id": 5,
      "room_number": "201",
      "room_type": "Deluxe",
      ...
    },
    {
      "id": 6,
      "room_number": "202",
      "room_type": "Deluxe",
      ...
    }
  ],
  "errors": "",
  "status": "OK"
}
```

---

## TypeScript Interfaces

```typescript
// types/room.ts

export interface Room {
  id: number;
  room_number: string;
  room_type: 'Standard' | 'Deluxe' | 'Suite';
  floor: number;
  capacity: number;
  price_per_night: number;
  status: 'available' | 'occupied' | 'maintenance' | 'cleaning';
  description: string;
  amenities: string;
  bed_type: string;
  created_at: string;
  updated_at: string;
}

export interface CreateRoomRequest {
  room_number: string;
  room_type: 'Standard' | 'Deluxe' | 'Suite';
  floor: number;
  capacity: number;
  price_per_night: number;
  description?: string;
  amenities?: string;
  bed_type?: string;
}

export interface UpdateRoomRequest {
  room_number?: string;
  room_type?: 'Standard' | 'Deluxe' | 'Suite';
  floor?: number;
  capacity?: number;
  price_per_night?: number;
  status?: 'available' | 'occupied' | 'maintenance' | 'cleaning';
  description?: string;
  amenities?: string;
  bed_type?: string;
}

export interface RoomListResponse {
  data: Room[];
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

export interface RoomTypeStats {
  room_type: string;
  total_rooms: number;
  available: number;
  occupied: number;
  maintenance: number;
  price_per_night: number;
}

export interface AvailabilityRequest {
  check_in: string;  // YYYY-MM-DD
  check_out: string; // YYYY-MM-DD
  room_type?: string;
}
```

---

## API Service Example

```typescript
// services/roomService.ts
import axios from '../lib/axios';
import { Room, CreateRoomRequest, UpdateRoomRequest, RoomListResponse, RoomTypeStats } from '../types/room';

export const roomService = {
  // Get all rooms with pagination
  async getAll(page = 1, pageSize = 10, status?: string): Promise<RoomListResponse> {
    const { data } = await axios.get('/rooms', {
      params: { page, page_size: pageSize, status }
    });
    return data.data;
  },

  // Get room by ID
  async getById(id: number): Promise<Room> {
    const { data } = await axios.get(`/rooms/${id}`);
    return data.data;
  },

  // Create room
  async create(room: CreateRoomRequest): Promise<Room> {
    const { data } = await axios.post('/rooms', room);
    return data.data;
  },

  // Update room
  async update(id: number, room: UpdateRoomRequest): Promise<Room> {
    const { data } = await axios.put(`/rooms/${id}`, room);
    return data.data;
  },

  // Delete room
  async delete(id: number): Promise<void> {
    await axios.delete(`/rooms/${id}`);
  },

  // Check availability
  async checkAvailability(checkIn: string, checkOut: string, roomType?: string): Promise<Room[]> {
    const { data } = await axios.get('/rooms/available', {
      params: { check_in: checkIn, check_out: checkOut, room_type: roomType }
    });
    return data.data;
  },

  // Update room status
  async updateStatus(id: number, status: string): Promise<void> {
    await axios.put(`/rooms/${id}/status`, { status });
  },

  // Get room statistics
  async getStats(): Promise<RoomTypeStats[]> {
    const { data } = await axios.get('/rooms/stats');
    return data.data;
  }
};
```

---

## Error Handling

All endpoints return errors in this format:

```json
{
  "message": "Error description",
  "data": null,
  "errors": "Detailed error message",
  "status": "BadRequest"
}
```

**Common HTTP Status Codes:**
| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 404 | Not Found |
| 409 | Conflict (duplicate room number) |
| 500 | Internal Server Error |

**Frontend Error Handling:**
```typescript
try {
  const room = await roomService.create(roomData);
  toast.success('Room created successfully');
} catch (err: any) {
  const errorMessage = err.response?.data?.errors || 'Failed to create room';
  toast.error(errorMessage);
}
```

---

## Seeded Test Data

The backend automatically seeds these rooms on startup:

| Room # | Type | Floor | Capacity | Price/Night |
|--------|------|-------|----------|-------------|
| 101-104 | Standard | 1 | 2 | $100 |
| 201-204 | Deluxe | 2 | 2-3 | $150 |
| 301-303 | Suite | 3 | 4 | $250 |

All rooms start with status `available`.

---

## Quick Start Checklist

1. ✅ Set up axios instance with base URL and auth interceptor
2. ✅ Create TypeScript interfaces for Room types
3. ✅ Implement roomService with all API calls
4. ✅ Create useRooms hook for room list with pagination
5. ✅ Implement room availability check in booking form
6. ✅ Add room management CRUD pages (admin)
7. ✅ Display room statistics on dashboard
