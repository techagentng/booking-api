package models

import (
	"time"
)

// CalendarAvailability represents daily availability status
type CalendarAvailability struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Date           time.Time `gorm:"uniqueIndex;not null" json:"date"`
	Status         string    `gorm:"not null" json:"status"` // available, booked, pending, closed, maintenance
	TotalSlots     int       `gorm:"default:0" json:"total_slots"`
	BookedSlots    int       `gorm:"default:0" json:"booked_slots"`
	AvailableSlots int       `gorm:"default:0" json:"available_slots"`
	Notes          string    `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Status constants
const (
	StatusAvailable   = "available"
	StatusBooked      = "booked"
	StatusPending     = "pending"
	StatusClosed      = "closed"
	StatusMaintenance = "maintenance"
)

// TableName specifies the table name
func (CalendarAvailability) TableName() string {
	return "calendar_availability"
}

// CalendarAvailabilityResponse is the response format for calendar availability
type CalendarAvailabilityResponse struct {
	ID             uint      `json:"id"`
	Date           string    `json:"date"`
	Status         string    `json:"status"`
	TotalSlots     int       `json:"total_slots"`
	BookedSlots    int       `json:"booked_slots"`
	AvailableSlots int       `json:"available_slots"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CalendarListResponse is the response format for calendar list
type CalendarListResponse struct {
	Data []CalendarAvailabilityResponse `json:"data"`
	Meta PaginationMeta                 `json:"meta"`
}

// CalendarUpdateRequest is the request payload for updating calendar availability
type CalendarUpdateRequest struct {
	Date   time.Time `json:"date" binding:"required"`
	Status string    `json:"status" binding:"required,oneof=available booked pending closed maintenance"`
	Notes  string    `json:"notes" binding:"max=500"`
}

// CalendarStatsResponse represents calendar statistics
type CalendarStatsResponse struct {
	TotalDays       int     `json:"total_days"`
	AvailableDays   int     `json:"available_days"`
	BookedDays      int     `json:"booked_days"`
	PendingDays     int     `json:"pending_days"`
	ClosedDays      int     `json:"closed_days"`
	MaintenanceDays int     `json:"maintenance_days"`
	TotalBookings   int     `json:"total_bookings"`
	Revenue         float64 `json:"revenue"`
	OccupancyRate   float64 `json:"occupancy_rate"`
}

// CalendarGenerateRequest is the request payload for generating calendar
type CalendarGenerateRequest struct {
	Year  int `json:"year" binding:"required,min=2020,max=2030"`
	Month int `json:"month" binding:"required,min=1,max=12"`
}
