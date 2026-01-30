package models

import (
	"time"
)

// TimeSlot represents available time slots for a specific date
type TimeSlot struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Date         time.Time `gorm:"not null;index" json:"date"`
	StartTime    string    `gorm:"not null" json:"start_time"` // "09:00"
	EndTime      string    `gorm:"not null" json:"end_time"`   // "10:00"
	Status       string    `gorm:"not null" json:"status"`     // available, booked, pending
	MaxCapacity  int       `gorm:"default:100" json:"max_capacity"`
	CurrentUsage int       `gorm:"default:0" json:"current_usage"`
	Price        float64   `gorm:"default:0" json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	CalendarAvailabilityID uint                  `gorm:"index" json:"calendar_availability_id"`
	CalendarAvailability   *CalendarAvailability `gorm:"foreignKey:CalendarAvailabilityID" json:"calendar_availability,omitempty"`
}

// TableName specifies the table name
func (TimeSlot) TableName() string {
	return "time_slots"
}

// TimeSlotResponse is the response format for time slots
type TimeSlotResponse struct {
	ID           uint      `json:"id"`
	Date         string    `json:"date"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	Status       string    `json:"status"`
	MaxCapacity  int       `json:"max_capacity"`
	CurrentUsage int       `json:"current_usage"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TimeSlotUpdateRequest is the request payload for updating time slots
type TimeSlotUpdateRequest struct {
	Date        time.Time `json:"date" binding:"required"`
	StartTime   string    `json:"start_time" binding:"required"`
	EndTime     string    `json:"end_time" binding:"required"`
	Status      string    `json:"status" binding:"required,oneof=available booked pending"`
	MaxCapacity int       `json:"max_capacity" binding:"omitempty,min=1,max=500"`
	Price       float64   `json:"price" binding:"omitempty,min=0"`
}

// TimeSlotAvailabilityRequest is the request payload for checking slot availability
type TimeSlotAvailabilityRequest struct {
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// TimeSlotAvailabilityResponse is the response for slot availability check
type TimeSlotAvailabilityResponse struct {
	Available bool               `json:"available"`
	Date      string             `json:"date"`
	StartTime string             `json:"start_time"`
	EndTime   string             `json:"end_time"`
	Message   string             `json:"message,omitempty"`
	Slots     []TimeSlotResponse `json:"slots,omitempty"`
}
