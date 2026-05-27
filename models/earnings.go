package models

import (
	"time"

	"github.com/google/uuid"
)

// Earning represents a provider's earning record
type Earning struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderID  uuid.UUID `gorm:"type:uuid;not null;index" json:"provider_id"`
	BookingID   uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	ServiceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Commission  float64   `gorm:"not null" json:"commission"` // Platform commission
	NetAmount   float64   `gorm:"not null" json:"net_amount"` // Amount after commission
	Status      string    `gorm:"default:'pending';index" json:"status"` // pending, available, paid
	EarnedAt    time.Time `json:"earned_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PayoutRequest represents a payout request from a provider
type PayoutRequest struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderID  uuid.UUID `gorm:"type:uuid;not null;index" json:"provider_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	BankAccount string    `gorm:"not null" json:"bank_account"`
	Status      string    `gorm:"default:'pending';index" json:"status"` // pending, processing, completed, rejected
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePayoutRequest represents a request to create a payout
type CreatePayoutRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	BankAccount string  `json:"bank_account" binding:"required"`
}

// ProviderAnalytics represents analytics data for a provider
type ProviderAnalytics struct {
	TotalRevenue     float64 `json:"total_revenue"`
	TotalBookings    int     `json:"total_bookings"`
	TotalServices    int     `json:"total_services"`
	AverageRating    float64 `json:"average_rating"`
	ResponseTime     float64 `json:"response_time"` // in hours
	PendingRequests  int     `json:"pending_requests"`
	CompletedBookings int    `json:"completed_bookings"`
	CancelledBookings int    `json:"cancelled_bookings"`
}

// DashboardAnalytics represents dashboard-specific analytics
type DashboardAnalytics struct {
	TotalRevenue      float64 `json:"total_revenue"`
	RevenueChange     float64 `json:"revenue_change"` // percentage change
	TotalBookings     int     `json:"total_bookings"`
	BookingsChange    float64 `json:"bookings_change"` // percentage change
	ActiveServices    int     `json:"active_services"`
	AverageRating     float64 `json:"average_rating"`
	RatingChange      float64 `json:"rating_change"` // percentage change
	ResponseTime      float64 `json:"response_time"`
	PendingRequests   int     `json:"pending_requests"`
}
