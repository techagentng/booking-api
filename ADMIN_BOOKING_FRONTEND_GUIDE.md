# 🎯 Admin Booking Frontend Integration Guide

## 📋 Overview

Complete guide for implementing the Admin Booking Management frontend. All backend endpoints are ready and tested.

---

## 🚀 Quick Start

### 1. API Service Setup

```typescript
// services/adminBookingService.ts
const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export interface BookingListParams {
  page?: number
  limit?: number
  status?: string[]
  dateFrom?: string
  dateTo?: string
  search?: string
  sortBy?: string
  sortOrder?: string
}

export interface AdminBookingResponse {
  id: number
  booking_id: string
  organizer_name: string
  organizer_email: string
  organizer_phone: string
  event_type: string
  guest_count: number
  special_requests: string
  booking_date: string
  start_time: string
  end_time: string
  total_price: number
  deposit_required: number
  payment_method: string
  status: string
  confirmed_by?: User
  confirmed_at?: string
  cancelled_by?: User
  cancelled_at?: string
  updated_by?: User
  status_history: BookingStatusHistory[]
  created_at: string
  updated_at: string
}

export interface BookingStatusHistory {
  id: number
  booking_id: number
  old_status?: string
  new_status: string
  changed_by: number
  changed_at: string
  notes: string
  changed_by_user: User
}

export interface BookingStats {
  total_bookings: number
  pending_bookings: number
  confirmed_bookings: number
  completed_bookings: number
  cancelled_bookings: number
  total_revenue: number
  revenue_by_status: {
    confirmed: number
    completed: number
  }
  popular_event_types: Array<{
    event_type: string
    count: number
  }>
  // Enhanced statistics
  average_guests?: number  // Average number of guests per booking
  occupancy_rate?: number // Occupancy rate as percentage
  
  // NEW: Chart data for frontend visualizations
  monthly_revenue?: Array<{ month: string; revenue: number }>   // Revenue Overview chart data
  monthly_bookings?: Array<{ month: string; bookings: number }>  // Booking Volume chart data
}

class AdminBookingService {
  private getAuthHeaders() {
    const token = localStorage.getItem('auth_token')
    return {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  }

  // 📊 Get all bookings with filtering
  async getBookings(params: BookingListParams = {}): Promise<{
    data: AdminBookingResponse[]
    meta: PaginationMeta
  }> {
    const queryParams = new URLSearchParams()
    
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.limit) queryParams.append('limit', params.limit.toString())
    if (params.status) params.status.forEach(s => queryParams.append('status', s))
    if (params.dateFrom) queryParams.append('date_from', params.dateFrom)
    if (params.dateTo) queryParams.append('date_to', params.dateTo)
    if (params.search) queryParams.append('search', params.search)
    if (params.sortBy) queryParams.append('sort_by', params.sortBy)
    if (params.sortOrder) queryParams.append('sort_order', params.sortOrder)

    const response = await fetch(
      `${API_BASE}/admin/bookings?${queryParams}`,
      { headers: this.getAuthHeaders() }
    )
    
    if (!response.ok) throw new Error('Failed to fetch bookings')
    return response.json()
  }

  // 🔍 Get booking by ID
  async getBooking(id: number): Promise<AdminBookingResponse> {
    const response = await fetch(
      `${API_BASE}/admin/bookings/${id}`,
      { headers: this.getAuthHeaders() }
    )
    
    if (!response.ok) throw new Error('Booking not found')
    return response.json()
  }

  // ✏️ Update booking status
  async updateStatus(id: number, status: string, notes?: string): Promise<AdminBookingResponse> {
    const response = await fetch(
      `${API_BASE}/admin/bookings/${id}/status`,
      {
        method: 'PUT',
        headers: this.getAuthHeaders(),
        body: JSON.stringify({ status, notes })
      }
    )
    
    if (!response.ok) throw new Error('Failed to update status')
    return response.json()
  }

  // 📜 Get booking history
  async getHistory(id: number): Promise<BookingStatusHistory[]> {
    const response = await fetch(
      `${API_BASE}/admin/bookings/${id}/history`,
      { headers: this.getAuthHeaders() }
    )
    
    if (!response.ok) throw new Error('Failed to fetch history')
    const result = await response.json()
    return result.data
  }

  // 📈 Get statistics
  async getStats(params: {
    period?: string
    dateFrom?: string
    dateTo?: string
  } = {}): Promise<BookingStats> {
    const queryParams = new URLSearchParams()
    if (params.period) queryParams.append('period', params.period)
    if (params.dateFrom) queryParams.append('date_from', params.dateFrom)
    if (params.dateTo) queryParams.append('date_to', params.dateTo)

    const response = await fetch(
      `${API_BASE}/admin/bookings/stats?${queryParams}`,
      { headers: this.getAuthHeaders() }
    )
    
    if (!response.ok) throw new Error('Failed to fetch stats')
    const result = await response.json()
    return result.data
  }

  // 📅 Get admin calendar
  async getCalendar(params: {
    year: number
    month: number
    includeBookings?: boolean
  }): Promise<CalendarAvailability[]> {
    const queryParams = new URLSearchParams()
    queryParams.append('year', params.year.toString())
    queryParams.append('month', params.month.toString())
    if (params.includeBookings) queryParams.append('include_bookings', 'true')

    const response = await fetch(
      `${API_BASE}/admin/calendar/availability?${queryParams}`,
      { headers: this.getAuthHeaders() }
    )
    
    if (!response.ok) throw new Error('Failed to fetch calendar')
    const result = await response.json()
    return result.data
  }
}

export const adminBookingService = new AdminBookingService()
```

---

## 🎨 Component Implementation

### 1. Main Admin Booking Page

```typescript
// pages/admin/bookings.tsx
import React, { useState, useEffect } from 'react'
import { BookingListTable } from '@/components/admin/BookingListTable'
import { BookingFilters } from '@/components/admin/BookingFilters'
import { BookingStats } from '@/components/admin/BookingStats'
import { StatusUpdateModal } from '@/components/admin/StatusUpdateModal'
import { BookingDetailsModal } from '@/components/admin/BookingDetailsModal'
import { adminBookingService, AdminBookingResponse } from '@/services/adminBookingService'

export default function AdminBookings() {
  const [bookings, setBookings] = useState<AdminBookingResponse[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedBooking, setSelectedBooking] = useState<AdminBookingResponse | null>(null)
  const [showStatusModal, setShowStatusModal] = useState(false)
  const [showDetailsModal, setShowDetailsModal] = useState(false)
  const [stats, setStats] = useState(null)

  // Filter states
  const [filters, setFilters] = useState({
    page: 1,
    limit: 20,
    status: [],
    dateFrom: '',
    dateTo: '',
    search: '',
    sortBy: 'created_at',
    sortOrder: 'desc'
  })

  // Load bookings
  const loadBookings = async () => {
    try {
      setLoading(true)
      const result = await adminBookingService.getBookings(filters)
      setBookings(result.data)
    } catch (error) {
      console.error('Failed to load bookings:', error)
    } finally {
      setLoading(false)
    }
  }

  // Load stats
  const loadStats = async () => {
    try {
      const stats = await adminBookingService.getStats({ period: 'month' })
      setStats(stats)
    } catch (error) {
      console.error('Failed to load stats:', error)
    }
  }

  useEffect(() => {
    loadBookings()
    loadStats()
  }, [filters])

  // Handle status update
  const handleStatusUpdate = async (status: string, notes?: string) => {
    if (!selectedBooking) return

    try {
      await adminBookingService.updateStatus(selectedBooking.id, status, notes)
      await loadBookings()
      setShowStatusModal(false)
      setSelectedBooking(null)
    } catch (error) {
      console.error('Failed to update status:', error)
    }
  }

  return (
    <div className="admin-bookings-page">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Admin Bookings</h1>
        <p className="text-gray-600 mt-2">Manage all hall bookings in the system</p>
      </div>

      {/* Statistics Dashboard */}
      {stats && <BookingStats stats={stats} />}

      {/* Filters */}
      <BookingFilters 
        filters={filters} 
        onFiltersChange={setFilters}
      />

      {/* Bookings Table */}
      <BookingListTable
        bookings={bookings}
        loading={loading}
        onStatusClick={(booking) => {
          setSelectedBooking(booking)
          setShowStatusModal(true)
        }}
        onDetailsClick={(booking) => {
          setSelectedBooking(booking)
          setShowDetailsModal(true)
        }}
      />

      {/* Status Update Modal */}
      <StatusUpdateModal
        booking={selectedBooking}
        isOpen={showStatusModal}
        onClose={() => {
          setShowStatusModal(false)
          setSelectedBooking(null)
        }}
        onUpdate={handleStatusUpdate}
      />

      {/* Details Modal */}
      <BookingDetailsModal
        booking={selectedBooking}
        isOpen={showDetailsModal}
        onClose={() => {
          setShowDetailsModal(false)
          setSelectedBooking(null)
        }}
      />
    </div>
  )
}
```

### 2. Booking Filters Component

```typescript
// components/admin/BookingFilters.tsx
import React from 'react'

interface BookingFiltersProps {
  filters: any
  onFiltersChange: (filters: any) => void
}

export function BookingFilters({ filters, onFiltersChange }: BookingFiltersProps) {
  const statusOptions = [
    { value: 'pending', label: 'Pending' },
    { value: 'confirmed', label: 'Confirmed' },
    { value: 'completed', label: 'Completed' },
    { value: 'cancelled', label: 'Cancelled' }
  ]

  const sortOptions = [
    { value: 'created_at', label: 'Created Date' },
    { value: 'booking_date', label: 'Booking Date' },
    { value: 'status', label: 'Status' },
    { value: 'organizer_name', label: 'Organizer Name' },
    { value: 'total_price', label: 'Total Price' }
  ]

  return (
    <div className="bg-white p-4 rounded-lg shadow mb-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        
        {/* Search */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Search
          </label>
          <input
            type="text"
            placeholder="Name, email, booking ID..."
            value={filters.search}
            onChange={(e) => onFiltersChange({ ...filters, search: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Status Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Status
          </label>
          <select
            multiple
            value={filters.status}
            onChange={(e) => {
              const selected = Array.from(e.target.selectedOptions, option => option.value)
              onFiltersChange({ ...filters, status: selected })
            }}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {statusOptions.map(option => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>

        {/* Date Range */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            From Date
          </label>
          <input
            type="date"
            value={filters.dateFrom}
            onChange={(e) => onFiltersChange({ ...filters, dateFrom: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            To Date
          </label>
          <input
            type="date"
            value={filters.dateTo}
            onChange={(e) => onFiltersChange({ ...filters, dateTo: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Sort Options */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Sort By
          </label>
          <select
            value={filters.sortBy}
            onChange={(e) => onFiltersChange({ ...filters, sortBy: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {sortOptions.map(option => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Order
          </label>
          <select
            value={filters.sortOrder}
            onChange={(e) => onFiltersChange({ ...filters, sortOrder: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="desc">Newest First</option>
            <option value="asc">Oldest First</option>
          </select>
        </div>

        {/* Clear Filters */}
        <div className="flex items-end">
          <button
            onClick={() => onFiltersChange({
              page: 1,
              limit: 20,
              status: [],
              dateFrom: '',
              dateTo: '',
              search: '',
              sortBy: 'created_at',
              sortOrder: 'desc'
            })}
            className="w-full px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
          >
            Clear Filters
          </button>
        </div>
      </div>
    </div>
  )
}
```

### 3. Booking List Table

```typescript
// components/admin/BookingListTable.tsx
import React from 'react'
import { AdminBookingResponse } from '@/services/adminBookingService'

interface BookingListTableProps {
  bookings: AdminBookingResponse[]
  loading: boolean
  onStatusClick: (booking: AdminBookingResponse) => void
  onDetailsClick: (booking: AdminBookingResponse) => void
}

export function BookingListTable({ 
  bookings, 
  loading, 
  onStatusClick, 
  onDetailsClick 
}: BookingListTableProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'bg-yellow-100 text-yellow-800'
      case 'confirmed': return 'bg-green-100 text-green-800'
      case 'completed': return 'bg-blue-100 text-blue-800'
      case 'cancelled': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow p-8">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-full mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-full mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-full"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Booking ID
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Organizer
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Event
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Price
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {bookings.map((booking) => (
            <tr key={booking.id} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {booking.booking_id}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div>
                  <div className="text-sm font-medium text-gray-900">
                    {booking.organizer_name}
                  </div>
                  <div className="text-sm text-gray-500">
                    {booking.organizer_email}
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div>
                  <div className="text-sm text-gray-900 capitalize">
                    {booking.event_type}
                  </div>
                  <div className="text-sm text-gray-500">
                    {booking.guest_count} guests
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                <div>
                  <div>{booking.booking_date}</div>
                  <div className="text-gray-500">
                    {booking.start_time} - {booking.end_time}
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(booking.status)}`}>
                  {booking.status}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                ${booking.total_price.toFixed(2)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button
                  onClick={() => onDetailsClick(booking)}
                  className="text-indigo-600 hover:text-indigo-900 mr-3"
                >
                  Details
                </button>
                <button
                  onClick={() => onStatusClick(booking)}
                  className="text-blue-600 hover:text-blue-900"
                >
                  Update Status
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      
      {bookings.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          No bookings found matching your criteria
        </div>
      )}
    </div>
  )
}
```

### 4. Status Update Modal

```typescript
// components/admin/StatusUpdateModal.tsx
import React, { useState } from 'react'
import { AdminBookingResponse } from '@/services/adminBookingService'

interface StatusUpdateModalProps {
  booking: AdminBookingResponse | null
  isOpen: boolean
  onClose: () => void
  onUpdate: (status: string, notes?: string) => void
}

export function StatusUpdateModal({ 
  booking, 
  isOpen, 
  onClose, 
  onUpdate 
}: StatusUpdateModalProps) {
  const [status, setStatus] = useState('')
  const [notes, setNotes] = useState('')
  const [loading, setLoading] = useState(false)

  const statusOptions = [
    { value: 'pending', label: 'Pending', color: 'yellow' },
    { value: 'confirmed', label: 'Confirmed', color: 'green' },
    { value: 'completed', label: 'Completed', color: 'blue' },
    { value: 'cancelled', label: 'Cancelled', color: 'red' }
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!booking || !status) return

    setLoading(true)
    try {
      await onUpdate(status, notes)
    } finally {
      setLoading(false)
    }
  }

  if (!isOpen || !booking) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-bold mb-4">Update Booking Status</h2>
        
        <div className="mb-4">
          <h3 className="font-medium text-gray-900">Booking Details</h3>
          <p className="text-sm text-gray-600">
            {booking.booking_id} - {booking.organizer_name}
          </p>
          <p className="text-sm text-gray-600">
            {booking.booking_date} - {booking.event_type}
          </p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              New Status
            </label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="">Select status...</option>
              {statusOptions.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Notes (optional)
            </label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Add notes about this status change..."
            />
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading || !status}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? 'Updating...' : 'Update Status'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
```

### 5. Statistics Dashboard

```typescript
// components/admin/BookingStats.tsx
import React from 'react'
import { BookingStats } from '@/services/adminBookingService'

interface BookingStatsProps {
  stats: BookingStats
}

export function BookingStats({ stats }: BookingStatsProps) {
  const statCards = [
    {
      label: 'Total Bookings',
      value: stats.total_bookings,
      color: 'bg-blue-500'
    },
    {
      label: 'Pending',
      value: stats.pending_bookings,
      color: 'bg-yellow-500'
    },
    {
      label: 'Confirmed',
      value: stats.confirmed_bookings,
      color: 'bg-green-500'
    },
    {
      label: 'Total Revenue',
      value: `$${stats.total_revenue.toLocaleString()}`,
      color: 'bg-purple-500'
    },
    // NEW: Enhanced statistics cards
    {
      label: 'Avg Guests',
      value: stats.average_guests ? stats.average_guests.toFixed(1) : 'N/A',
      color: 'bg-indigo-500'
    },
    {
      label: 'Occupancy Rate',
      value: stats.occupancy_rate ? `${stats.occupancy_rate.toFixed(1)}%` : 'N/A',
      color: 'bg-orange-500'
    },
    {
      label: 'Completed',
      value: stats.completed_bookings,
      color: 'bg-teal-500'
    },
    {
      label: 'Cancelled',
      value: stats.cancelled_bookings,
      color: 'bg-red-500'
    }
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      {statCards.map((stat, index) => (
        <div key={index} className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className={`w-12 h-12 ${stat.color} rounded-lg flex items-center justify-center text-white`}>
              <div className="text-2xl font-bold">{stat.value}</div>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">{stat.label}</p>
              <p className="text-lg font-semibold text-gray-900">{stat.value}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
```

### 6. Chart Components

```typescript
// components/admin/RevenueOverviewChart.tsx
import React from 'react'
import { BookingStats } from '@/services/adminBookingService'

interface RevenueOverviewChartProps {
  stats: BookingStats
}

export function RevenueOverviewChart({ stats }: RevenueOverviewChartProps) {
  if (!stats.monthly_revenue || stats.monthly_revenue.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Revenue Overview</h3>
        <div className="flex items-center justify-center h-64 text-gray-500">
          No revenue data available
        </div>
      </div>
    )
  }

  const chartData = stats.monthly_revenue.map(item => ({
    month: new Date(item.month).toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
    revenue: item.revenue
  }))

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Revenue Overview</h3>
      <div className="h-64">
        {/* Your chart library implementation */}
        <div className="text-sm text-gray-600">
          Revenue data for {chartData.length} months
        </div>
      </div>
    </div>
  )
}

// components/admin/BookingVolumeChart.tsx
import React from 'react'
import { BookingStats } from '@/services/adminBookingService'

interface BookingVolumeChartProps {
  stats: BookingStats
}

export function BookingVolumeChart({ stats }: BookingVolumeChartProps) {
  if (!stats.monthly_bookings || stats.monthly_bookings.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Booking Volume</h3>
        <div className="flex items-center justify-center h-64 text-gray-500">
          No booking volume data available
        </div>
      </div>
    )
  }

  const chartData = stats.monthly_bookings.map(item => ({
    month: new Date(item.month).toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
    bookings: item.bookings
  }))

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Booking Volume</h3>
      <div className="h-64">
        {/* Your chart library implementation */}
        <div className="text-sm text-gray-600">
          Booking data for {chartData.length} months
        </div>
      </div>
    </div>
  )
}
```

---

## 🎯 Implementation Steps

### Phase 1: Core Setup
1. ✅ Copy the `adminBookingService.ts` file
2. ✅ Install required dependencies: `axios`, `react-hot-toast`
3. ✅ Set up authentication headers
4. ✅ Test API connectivity

### Phase 2: Basic Components
1. 📝 Create the main booking page
2. 📝 Implement the filters component
3. 📝 Build the booking table
4. 📝 Add status update modal

### Phase 3: Advanced Features
1. 📝 Add statistics dashboard
2. 📝 Implement booking details modal
3. 📝 Add pagination
4. 📝 Add export functionality

### Phase 4: Calendar Integration
1. 📝 Build calendar view component
2. 📝 Add drag-and-drop status updates
3. 📝 Integrate with booking management
4. 📝 Add bulk operations

---

## 🔧 Quick Testing

```bash
# Test the API endpoints
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/v1/admin/bookings?page=1&limit=10"

# Test status update
curl -X PUT \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed","notes":"Test update"}' \
  "http://localhost:8080/api/v1/admin/bookings/1/status"
```

---

## 🚀 Next Steps

1. **Start with the API service** - Copy and configure it
2. **Build the main page** - Use the provided template
3. **Add components one by one** - Follow the implementation order
4. **Test each endpoint** - Verify API connectivity
5. **Style with Tailwind** - All components use Tailwind classes

**Ready to start building? Which component would you like to implement first?** 🎯
