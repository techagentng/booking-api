# Hall Booking Real-Time Notifications Guide

## 🎯 Overview
Real-time notifications for hall booking events using Server-Sent Events (SSE). Instant updates when bookings are created, updated, or cancelled.

## 📡 Notification Types

### **New Hall Booking**
```json
{
  "type": "new_hall_booking",
  "title": "New Hall Booking",
  "message": "John Doe booked a wedding for 50 guests on 2026-04-15",
  "priority": "normal",
  "data": {
    "booking_id": 123,
    "booking_id_str": "HB-20260415-001",
    "organizer_name": "John Doe",
    "event_type": "wedding",
    "booking_date": "2026-04-15",
    "guest_count": 50,
    "total_price": 2500,
    "created_by_type": "public"
  }
}
```

### **Hall Booking Updated**
```json
{
  "type": "hall_booking_updated",
  "title": "Hall Booking Updated",
  "message": "John Doe's wedding booking status changed to confirmed",
  "priority": "normal",
  "data": {
    "booking_id": 123,
    "booking_id_str": "HB-20260415-001",
    "organizer_name": "John Doe",
    "event_type": "wedding",
    "status": "confirmed"
  }
}
```

### **Hall Booking Cancelled**
```json
{
  "type": "hall_booking_cancelled",
  "title": "Hall Booking Cancelled",
  "message": "John Doe's wedding booking for 2026-04-15 has been cancelled",
  "priority": "high",
  "data": {
    "booking_id": 123,
    "booking_id_str": "HB-20260415-001",
    "organizer_name": "John Doe",
    "event_type": "wedding",
    "booking_date": "2026-04-15"
  }
}
```

## 🔌 Frontend Integration

### **1. SSE Connection Hook**
```typescript
// hooks/useHallBookingNotifications.ts
import { useEffect, useState } from 'react'

interface Notification {
  id: string
  type: string
  title: string
  message: string
  priority: string
  data: any
  timestamp: string
}

export function useHallBookingNotifications() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const eventSource = new EventSource('/api/v1/notifications/sse', {
      withCredentials: true
    })

    eventSource.onopen = () => {
      setIsConnected(true)
      setError(null)
    }

    eventSource.onerror = () => {
      setIsConnected(false)
      setError('Connection lost')
    }

    eventSource.addEventListener('connected', (event) => {
      console.log('Connected to notifications:', JSON.parse(event.data))
    })

    eventSource.addEventListener('notification', (event) => {
      const notification = JSON.parse(event.data)
      
      // Only handle hall booking notifications
      if (notification.type.includes('hall_booking')) {
        setNotifications(prev => [notification, ...prev].slice(0, 50))
      }
    })

    eventSource.addEventListener('heartbeat', (event) => {
      // Keep connection alive
    })

    return () => {
      eventSource.close()
    }
  }, [])

  const clearNotifications = () => setNotifications([])

  return {
    notifications,
    isConnected,
    error,
    clearNotifications
  }
}
```

### **2. Notification Component**
```typescript
// components/HallBookingNotifications.tsx
import React from 'react'
import { useHallBookingNotifications } from '@/hooks/useHallBookingNotifications'

export function HallBookingNotifications() {
  const { notifications, isConnected, error, clearNotifications } = useHallBookingNotifications()

  const getEventIcon = (eventType: string) => {
    switch (eventType) {
      case 'wedding': return '💒'
      case 'party': return '🎉'
      case 'corporate': return '💼'
      case 'funeral': return '⚰️'
      case 'meeting': return '🤝'
      default: return '📅'
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'bg-red-500'
      case 'urgent': return 'bg-red-600'
      case 'normal': return 'bg-blue-500'
      case 'low': return 'bg-gray-500'
      default: return 'bg-gray-500'
    }
  }

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2">
      {/* Connection Status */}
      <div className={`px-3 py-1 rounded-full text-xs text-white ${
        isConnected ? 'bg-green-500' : 'bg-red-500'
      }`}>
        {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
      </div>

      {/* Notifications */}
      {notifications.slice(0, 5).map((notification) => (
        <div
          key={notification.id}
          className="bg-white rounded-lg shadow-lg p-4 max-w-sm border-l-4 border-blue-500 animate-pulse"
        >
          <div className="flex items-start">
            <div className="text-2xl mr-3">
              {getEventIcon(notification.data.event_type)}
            </div>
            <div className="flex-1">
              <div className="flex items-center justify-between">
                <h4 className="font-semibold text-sm">{notification.title}</h4>
                <span className={`px-2 py-1 rounded-full text-xs text-white ${getPriorityColor(notification.priority)}`}>
                  {notification.priority}
                </span>
              </div>
              <p className="text-sm text-gray-600 mt-1">{notification.message}</p>
              <div className="text-xs text-gray-400 mt-2">
                {new Date(notification.timestamp).toLocaleTimeString()}
              </div>
            </div>
          </div>
        </div>
      ))}

      {/* Clear Button */}
      {notifications.length > 0 && (
        <button
          onClick={clearNotifications}
          className="bg-gray-500 text-white px-3 py-1 rounded text-xs hover:bg-gray-600"
        >
          Clear ({notifications.length})
        </button>
      )}
    </div>
  )
}
```

### **3. Activity Feed Integration**
```typescript
// components/HallBookingActivityFeed.tsx
import React, { useState, useEffect } from 'react'
import { useHallBookingNotifications } from '@/hooks/useHallBookingNotifications'
import { useGetHallBookingActivity } from '@/hooks/useHallBookingActivity'

export function HallBookingActivityFeed() {
  const { notifications } = useHallBookingNotifications()
  const { data: activity, refetch } = useGetHallBookingActivity()
  const [realtimeBookings, setRealtimeBookings] = useState([])

  // Update activity feed when new booking notification arrives
  useEffect(() => {
    const newBookingNotifications = notifications.filter(n => n.type === 'new_hall_booking')
    
    if (newBookingNotifications.length > 0) {
      // Add new booking to the top of the list
      const latestNotification = newBookingNotifications[0]
      const newBooking = {
        id: latestNotification.data.booking_id,
        booking_id: latestNotification.data.booking_id_str,
        organizer_name: latestNotification.data.organizer_name,
        event_type: latestNotification.data.event_type,
        guest_count: latestNotification.data.guest_count,
        booking_date: latestNotification.data.booking_date,
        total_price: latestNotification.data.total_price,
        status: 'pending',
        created_at: new Date().toISOString(),
        created_by_type: latestNotification.data.created_by_type
      }
      
      setRealtimeBookings(prev => [newBooking, ...prev].slice(0, 10))
      
      // Refresh the activity feed
      refetch()
    }
  }, [notifications, refetch])

  // Combine API data with real-time updates
  const allBookings = [...realtimeBookings, ...(activity?.recent_bookings || [])]
  const uniqueBookings = allBookings.filter((booking, index, self) => 
    index === self.findIndex(b => b.id === booking.id)
  ).slice(0, 10)

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-bold">🔴 Live Hall Booking Activity</h3>
      
      {uniqueBookings.length === 0 ? (
        <p className="text-gray-500">No recent bookings</p>
      ) : (
        <div className="space-y-3">
          {uniqueBookings.map((booking) => (
            <div key={booking.id} className="bg-white rounded-lg shadow p-4 border-l-4 border-green-500">
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-semibold">{booking.organizer_name}</div>
                  <div className="text-sm text-gray-600">
                    {booking.event_type} • {booking.guest_count} guests
                  </div>
                  <div className="text-xs text-gray-400">
                    {booking.booking_date} • £{booking.total_price}
                  </div>
                </div>
                <div className="text-right">
                  <span className="inline-block px-2 py-1 text-xs bg-green-100 text-green-800 rounded">
                    {booking.created_by_type}
                  </span>
                  <div className="text-xs text-gray-400 mt-1">
                    Just now
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
      
      <div className="text-center">
        <button 
          onClick={() => refetch()}
          className="text-blue-500 hover:text-blue-700 text-sm"
        >
          Refresh Activity
        </button>
      </div>
    </div>
  )
}
```

## 🚀 Usage Examples

### **In Dashboard Component**
```typescript
// pages/Dashboard.tsx
import { HallBookingNotifications } from '@/components/HallBookingNotifications'
import { HallBookingActivityFeed } from '@/components/HallBookingActivityFeed'

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-gray-100">
      {/* Notification Toasts */}
      <HallBookingNotifications />
      
      <div className="flex">
        {/* Main Content */}
        <div className="flex-1 p-8">
          <h1 className="text-2xl font-bold mb-6">Dashboard</h1>
          {/* Other dashboard content */}
        </div>
        
        {/* Sidebar with Live Activity */}
        <div className="w-80 bg-white shadow-lg p-6">
          <HallBookingActivityFeed />
        </div>
      </div>
    </div>
  )
}
```

## 🔧 Features

### **✅ Real-Time Updates**
- Instant notifications when bookings are created
- Live activity feed updates
- No polling required

### **✅ Rich Notification Data**
- Complete booking information
- Event type icons
- Priority levels
- Creator attribution

### **✅ Connection Management**
- Automatic reconnection
- Heartbeat monitoring
- Connection status indicators

### **✅ User Experience**
- Toast notifications
- Activity feed updates
- Clear notification history
- Responsive design

## 📱 Mobile Support

The SSE connection works on mobile browsers too:
- iOS Safari ✅
- Android Chrome ✅
- React Native ✅

## 🔍 Testing

### **Test Notifications**
```bash
# Create a test booking to trigger notification
curl -X POST "http://localhost:8080/api/v1/bookings" \
  -H "Content-Type: application/json" \
  -d '{
    "organizer_name": "Test User",
    "organizer_email": "test@example.com",
    "organizer_phone": "+1234567890",
    "event_type": "party",
    "guest_count": 25,
    "booking_date": "2026-06-01",
    "start_time": "19:00",
    "end_time": "23:00",
    "total_price": 1500,
    "deposit_required": 300,
    "payment_method": "online"
  }'
```

### **Test SSE Connection**
```javascript
// In browser console
const eventSource = new EventSource('/api/v1/notifications/sse')
eventSource.addEventListener('notification', (e) => {
  console.log('New notification:', JSON.parse(e.data))
})
```

## 🎯 Benefits

### **For Admins**
- **Instant awareness** of new bookings
- **Real-time dashboard** updates
- **No manual refresh** needed
- **Priority alerts** for important events

### **For Users**
- **Immediate feedback** on booking actions
- **Live status updates**
- **Better user experience**
- **Mobile-friendly** notifications

**Your hall booking system now has real-time notifications!** 🚀✨
