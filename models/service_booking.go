package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceBooking represents a service booking (non-hotel services)
type ServiceBooking struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"customer_id"`
	ProviderID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"provider_id"`
	ServiceID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"service_id"`
	BookingDate     time.Time  `json:"booking_date"`
	CheckInDate     *time.Time `json:"check_in_date"`
	CheckOutDate    *time.Time `json:"check_out_date"`
	GuestCount      int        `json:"guest_count"`
	SpecialRequests string     `gorm:"type:text" json:"special_requests"`
	Status          string     `gorm:"default:'pending';index" json:"status"` // pending, confirmed, cancelled, completed, rejected
	Price           float64    `json:"price"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateServiceBookingRequest represents a request to create a service booking
type CreateServiceBookingRequest struct {
	ServiceID       uuid.UUID `json:"service_id" binding:"required"`
	CheckInDate     string    `json:"check_in_date"`
	CheckOutDate    string    `json:"check_out_date"`
	GuestCount      int       `json:"guest_count"`
	SpecialRequests string    `json:"special_requests"`
}

// ServiceBookingResponse represents a service booking response
type ServiceBookingResponse struct {
	ID              uuid.UUID  `json:"id"`
	CustomerID      uuid.UUID  `json:"customer_id"`
	ProviderID      uuid.UUID  `json:"provider_id"`
	ServiceID       uuid.UUID  `json:"service_id"`
	BookingDate     time.Time  `json:"booking_date"`
	CheckInDate     *time.Time `json:"check_in_date"`
	CheckOutDate    *time.Time `json:"check_out_date"`
	GuestCount      int        `json:"guest_count"`
	SpecialRequests string     `json:"special_requests"`
	Status          string     `json:"status"`
	Price           float64    `json:"price"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
