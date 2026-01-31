# 🎯 Hall Booking Statuses - Frontend Integration Guide

## 📋 Overview

This guide explains how to retrieve and manage hall booking statuses (pending, confirmed, cancelled, completed) from the frontend. The backend provides comprehensive endpoints for filtering, updating, and tracking booking statuses.

---

## 🚀 Available Status Values

Hall bookings can have the following statuses:

| Status | Description | Typical Use |
|--------|-------------|-------------|
| `pending` | Awaiting admin confirmation | New bookings that need review |
| `confirmed` | Booking approved and scheduled | Ready for event execution |
| `cancelled` | Booking was cancelled | By admin or customer request |
| `completed` | Event has taken place | Post-event processing |

---

## 🛠 API Endpoints

### 1. Get All Bookings with Status Filtering

**Endpoint**: `GET /api/v1/admin/bookings`  
**Authentication**: Required (Bearer token)  
**Purpose**: Retrieve bookings filtered by status with pagination

#### Query Parameters

| Parameter | Type | Required | Example | Description |
|-----------|------|----------|---------|-------------|
| `status` | string | No | `pending` | Filter by specific status |
| `status` | array | No | `pending&status=confirmed` | Filter by multiple statuses |
| `page` | number | No | `1` | Page number (default: 1) |
| `limit` | number | No | `20` | Items per page (default: 20) |
| `date_from` | string | No | `2024-01-01` | Filter by start date |
| `date_to` | string | No | `2024-12-31` | Filter by end date |
| `search` | string | No | `john` | Search by organizer name/email |
| `sort_by` | string | No | `created_at` | Sort field |
| `sort_order` | string | No | `desc` | Sort direction (asc/desc) |

#### Request Examples

```typescript
// Get all pending bookings
GET /api/v1/admin/bookings?status=pending

// Get multiple statuses
GET /api/v1/admin/bookings?status=pending&status=confirmed

// Get all bookings with pagination
GET /api/v1/admin/bookings?page=1&limit=10

// Get bookings by date range
GET /api/v1/admin/bookings?date_from=2024-01-01&date_to=2024-01-31

// Search and filter
GET /api/v1/admin/bookings?search=john&status=confirmed
```

#### Response Format

```json
{
  "message": "Bookings retrieved successfully",
  "data": [
    {
      "id": 1,
      "booking_id": "HB-2024-001",
      "organizer_name": "John Doe",
      "organizer_email": "john@example.com",
      "organizer_phone": "+1234567890",
      "event_type": "wedding",
      "guest_count": 50,
      "special_requests": "Floral arrangements",
      "booking_date": "2024-06-15",
      "start_time": "14:00",
      "end_time": "18:00",
      "total_price": 2500.00,
      "deposit_required": 500.00,
      "payment_method": "online",
      "status": "pending",
      "confirmed_by": null,
      "confirmed_at": null,
      "cancelled_by": null,
      "cancelled_at": null,
      "updated_by": null,
      "status_history": [],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "total": 25,
    "page": 1,
    "page_size": 20,
    "total_pages": 2
  }
}
```

---

### 2. Get Single Booking by ID

**Endpoint**: `GET /api/v1/admin/bookings/:id`  
**Authentication**: Required  
**Purpose**: Retrieve detailed information for a specific booking

```typescript
// Example request
GET /api/v1/admin/bookings/123
```

---

### 3. Update Booking Status

**Endpoint**: `PUT /api/v1/admin/bookings/:id/status`  
**Authentication**: Required  
**Purpose**: Update the status of a booking

#### Request Body

```json
{
  "status": "confirmed",  // Required: pending, confirmed, completed, cancelled
  "notes": "Customer called to confirm the event details"
}
```

#### Example Request

```typescript
// Confirm a booking
PUT /api/v1/admin/bookings/123/status
{
  "status": "confirmed",
  "notes": "Deposit received, booking confirmed"
}

// Cancel a booking
PUT /api/v1/admin/bookings/123/status
{
  "status": "cancelled",
  "notes": "Customer requested cancellation due to scheduling conflict"
}

// Mark as completed
PUT /api/v1/admin/bookings/123/status
{
  "status": "completed",
  "notes": "Event successfully completed"
}
```

---

### 4. Get Booking Status History

**Endpoint**: `GET /api/v1/admin/bookings/:id/history`  
**Authentication**: Required  
**Purpose**: Track all status changes for a specific booking

#### Response Format

```json
{
  "message": "Booking history retrieved successfully",
  "data": [
    {
      "id": 1,
      "booking_id": 123,
      "old_status": null,
      "new_status": "pending",
      "changed_by": 1,
      "changed_at": "2024-01-15T10:30:00Z",
      "notes": "Initial booking created"
    },
    {
      "id": 2,
      "booking_id": 123,
      "old_status": "pending",
      "new_status": "confirmed",
      "changed_by": 1,
      "changed_at": "2024-01-16T09:15:00Z",
      "notes": "Deposit received, booking confirmed"
    }
  ]
}
```

---

### 5. Get Booking Statistics by Status

**Endpoint**: `GET /api/v1/admin/bookings/stats`  
**Authentication**: Required  
**Purpose**: Get statistics grouped by booking status

#### Query Parameters

| Parameter | Type | Required | Example | Description |
|-----------|------|----------|---------|-------------|
| `period` | string | No | `month` | today, week, month, year |
| `date_from` | string | No | `2024-01-01` | Custom date range start |
| `date_to` | string | No | `2024-12-31` | Custom date range end |

#### Example Requests

```typescript
// Get monthly stats
GET /api/v1/admin/bookings/stats?period=month

// Get custom date range stats
GET /api/v1/admin/bookings/stats?date_from=2024-01-01&date_to=2024-01-31
```

#### Response Format

```json
{
  "data": {
    "data": {
      "total_bookings": 18,
      "pending_bookings": 3,
      "confirmed_bookings": 0,
      "completed_bookings": 0,
      "cancelled_bookings": 0,
      "total_revenue": 6483,
      "revenue_by_status": {
        "confirmed": 6483,
        "completed": 0
      },
      "average_guests": 22.72,
      "occupancy_rate": 22.72,
      "popular_event_types": [
        {
          "event_type": "party",
          "count": 9
        },
        {
          "event_type": "wedding",
          "count": 5
        }
      ],
      "monthly_bookings": [
        {
          "month": "2026-01-01 00:00:00+01",
          "bookings": 2
        }
      ],
      "monthly_revenue": [
        {
          "month": "2026-01-01 00:00:00+01",
          "revenue": 323
        }
      ]
    }
  },
  "errors": "",
  "message": "Booking statistics retrieved successfully",
  "status": "OK"
}
```

#### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `total_bookings` | number | Total number of bookings |
| `pending_bookings` | number | Count of pending bookings |
| `confirmed_bookings` | number | Count of confirmed bookings |
| `completed_bookings` | number | Count of completed bookings |
| `cancelled_bookings` | number | Count of cancelled bookings |
| `total_revenue` | number | Total revenue from confirmed + completed bookings |
| `revenue_by_status` | object | Revenue breakdown by status |
| `average_guests` | number | Average guest count per booking |
| `occupancy_rate` | number | Occupancy rate percentage |
| `popular_event_types` | array | Top 5 most popular event types |
| `monthly_bookings` | array | Monthly booking volume data |
| `monthly_revenue` | array | Monthly revenue data |

---

## 💻 Frontend Implementation Examples

### React Hook Example

```typescript
// hooks/useHallBookings.ts
import { useState, useEffect } from 'react';
import axios from 'axios';

interface BookingFilters {
  status?: string | string[];
  page?: number;
  limit?: number;
  dateFrom?: string;
  dateTo?: string;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
}

export const useHallBookings = (filters: BookingFilters = {}) => {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [meta, setMeta] = useState(null);

  const fetchBookings = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const params = new URLSearchParams();
      
      // Handle single status or array of statuses
      if (filters.status) {
        if (Array.isArray(filters.status)) {
          filters.status.forEach(status => params.append('status', status));
        } else {
          params.append('status', filters.status);
        }
      }
      
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.limit) params.append('limit', filters.limit.toString());
      if (filters.dateFrom) params.append('date_from', filters.dateFrom);
      if (filters.dateTo) params.append('date_to', filters.dateTo);
      if (filters.search) params.append('search', filters.search);
      if (filters.sortBy) params.append('sort_by', filters.sortBy);
      if (filters.sortOrder) params.append('sort_order', filters.sortOrder);

      const response = await axios.get(`/api/v1/admin/bookings?${params}`);
      setBookings(response.data.data);
      setMeta(response.data.meta);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch bookings');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBookings();
  }, [JSON.stringify(filters)]);

  return { bookings, loading, error, meta, refetch: fetchBookings };
};
```

### Status Update Function

```typescript
// services/bookingService.ts
export const updateBookingStatus = async (
  bookingId: number,
  status: 'pending' | 'confirmed' | 'completed' | 'cancelled',
  notes?: string
) => {
  try {
    const response = await axios.put(`/api/v1/admin/bookings/${bookingId}/status`, {
      status,
      notes: notes || ''
    });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Failed to update status');
  }
};
```

### React Component Example

```typescript
// components/BookingStatusManager.tsx
import React from 'react';
import { useHallBookings } from '../hooks/useHallBookings';
import { updateBookingStatus } from '../services/bookingService';

const BookingStatusManager: React.FC = () => {
  const [selectedStatus, setSelectedStatus] = React.useState<string>('all');
  const { bookings, loading, error, meta, refetch } = useHallBookings({
    status: selectedStatus === 'all' ? undefined : selectedStatus,
    page: 1,
    limit: 20
  });

  const handleStatusChange = async (bookingId: number, newStatus: string, notes?: string) => {
    try {
      await updateBookingStatus(bookingId, newStatus as any, notes);
      refetch(); // Refresh the list
    } catch (error) {
      console.error('Failed to update status:', error);
    }
  };

  const statusFilters = [
    { value: 'all', label: 'All Bookings', color: 'gray' },
    { value: 'pending', label: 'Pending', color: 'yellow' },
    { value: 'confirmed', label: 'Confirmed', color: 'green' },
    { value: 'completed', label: 'Completed', color: 'blue' },
    { value: 'cancelled', label: 'Cancelled', color: 'red' }
  ];

  if (loading) return <div>Loading bookings...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      {/* Status Filter Tabs */}
      <div className="flex space-x-4 mb-6">
        {statusFilters.map(filter => (
          <button
            key={filter.value}
            onClick={() => setSelectedStatus(filter.value)}
            className={`px-4 py-2 rounded ${
              selectedStatus === filter.value
                ? `bg-${filter.color}-500 text-white`
                : 'bg-gray-200 text-gray-700'
            }`}
          >
            {filter.label}
          </button>
        ))}
      </div>

      {/* Bookings List */}
      <div className="space-y-4">
        {bookings.map(booking => (
          <div key={booking.id} className="border rounded p-4">
            <div className="flex justify-between items-start">
              <div>
                <h3 className="font-semibold">{booking.booking_id}</h3>
                <p>{booking.organizer_name} - {booking.event_type}</p>
                <p className="text-sm text-gray-600">
                  {booking.booking_date} at {booking.start_time}-{booking.end_time}
                </p>
              </div>
              
              <div className="flex items-center space-x-2">
                <span className={`px-2 py-1 rounded text-sm ${
                  booking.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                  booking.status === 'confirmed' ? 'bg-green-100 text-green-800' :
                  booking.status === 'completed' ? 'bg-blue-100 text-blue-800' :
                  'bg-red-100 text-red-800'
                }`}>
                  {booking.status}
                </span>
                
                <select
                  value={booking.status}
                  onChange={(e) => handleStatusChange(booking.id, e.target.value)}
                  className="border rounded px-2 py-1 text-sm"
                >
                  <option value="pending">Pending</option>
                  <option value="confirmed">Confirmed</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
```

---

## 🔄 Status Workflow

### Typical Status Flow

1. **pending** → **confirmed**
   - When admin reviews and approves the booking
   - Deposit payment received
   - Venue availability confirmed

2. **pending** → **cancelled**
   - Customer requests cancellation
   - Admin rejects booking
   - Payment issues

3. **confirmed** → **completed**
   - Event successfully conducted
   - Post-event cleanup completed
   - Final payment processed

4. **confirmed** → **cancelled**
   - Customer cancels before event
   - Force majeure situations
   - Venue issues

### Status Validation Rules

The backend enforces these status transitions:
- `pending` → `confirmed`, `cancelled`
- `confirmed` → `completed`, `cancelled`
- `completed` → (no further changes)
- `cancelled` → (no further changes)

---

## 📊 Real-time Updates

For real-time status updates, you can implement polling or use Server-Sent Events (SSE):

```typescript
// Polling example
useEffect(() => {
  const interval = setInterval(() => {
    refetch();
  }, 30000); // Poll every 30 seconds

  return () => clearInterval(interval);
}, [refetch]);
```

---

## 🎨 UI/UX Best Practices

1. **Color Coding**: Use consistent colors for each status
   - Pending: Yellow/Orange
   - Confirmed: Green
   - Completed: Blue
   - Cancelled: Red

2. **Status Badges**: Make status easily scannable with badges

3. **Bulk Actions**: Allow updating multiple bookings at once

4. **Status History**: Show timeline of status changes for transparency

5. **Filters**: Provide quick filters for each status

6. **Search**: Combine status filters with search functionality

---

## 🚨 Error Handling

Common errors and how to handle them:

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Missing/invalid token | Redirect to login |
| 403 Forbidden | Insufficient permissions | Show permission error |
| 404 Not Found | Booking doesn't exist | Show booking not found |
| 400 Bad Request | Invalid status transition | Show validation error |

---

## 📱 Mobile Considerations

- Use dropdown selects for mobile status updates
- Implement swipe actions for quick status changes
- Show status prominently in mobile list view
- Use touch-friendly buttons for status actions

---

## 🔧 Testing

Test these scenarios:
1. Filter by each status individually
2. Filter by multiple statuses
3. Status transitions work correctly
4. Invalid status transitions are rejected
5. Pagination with status filters
6. Search combined with status filters
7. Empty states for each status filter

---

## 📞 Support

For any issues with the booking status endpoints:
1. Check the API response for detailed error messages
2. Verify authentication tokens are valid
3. Ensure status values match exactly (case-sensitive)
4. Check network connectivity for API calls

---

**🎉 Your frontend is now ready to handle hall booking statuses efficiently!**
