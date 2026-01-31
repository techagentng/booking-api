package models

import (
	"time"
)

// StripePayment represents a Stripe payment transaction
type StripePayment struct {
	ID              string    `gorm:"primaryKey" json:"id"` // Stripe PaymentIntent ID
	BookingID       uint      `gorm:"not null" json:"booking_id"`
	Amount          int64     `gorm:"not null" json:"amount"` // Amount in pence
	Currency        string    `gorm:"default:'gbp'" json:"currency"`
	Status          string    `gorm:"not null" json:"status"` // succeeded, pending, failed, cancelled, refunded
	PaymentMethodID string    `json:"payment_method_id"`
	CustomerID      string    `json:"customer_id"`
	StripeSessionID string    `json:"stripe_session_id"`
	ReceiptURL      string    `json:"receipt_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	Booking HallBooking `gorm:"foreignKey:BookingID" json:"booking"`
}

// StripePaymentIntent represents a Stripe payment intent
type StripePaymentIntent struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	BookingID    uint      `gorm:"not null" json:"booking_id"`
	Amount       int64     `gorm:"not null" json:"amount"`
	Currency     string    `gorm:"default:'gbp'" json:"currency"`
	Status       string    `gorm:"not null" json:"status"`
	ClientSecret string    `json:"client_secret"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Booking HallBooking `gorm:"foreignKey:BookingID" json:"booking"`
}

// WebhookEvent represents a Stripe webhook event
type WebhookEvent struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	StripeEventID string    `gorm:"unique" json:"stripe_event_id"`
	EventType     string    `gorm:"not null" json:"event_type"`
	Processed     bool      `gorm:"default:false" json:"processed"`
	Payload       string    `gorm:"type:text" json:"payload"`
	ErrorMessage  string    `json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentRequest represents a payment creation request
type PaymentRequest struct {
	BookingID uint `json:"booking_id" binding:"required"`
}

// CheckoutSessionResponse represents checkout session response
type CheckoutSessionResponse struct {
	SessionID string `json:"id"`
	URL       string `json:"url"`
}

// PaymentIntentResponse represents payment intent response
type PaymentIntentResponse struct {
	ClientSecret string `json:"client_secret"`
	PaymentID    string `json:"payment_id"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
	Amount    int64  `json:"amount"` // Amount in pence, 0 for full refund
}

// PaymentStatus represents payment status constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
	PaymentStatusCancelled = "cancelled"
	PaymentStatusRefunded  = "refunded"
)

// PaymentEventTypes represents Stripe webhook event types
const (
	EventCheckoutSessionCompleted = "checkout.session.completed"
	EventPaymentIntentSucceeded   = "payment_intent.succeeded"
	EventPaymentIntentFailed      = "payment_intent.payment_failed"
	EventPaymentIntentCanceled    = "payment_intent.canceled"
	EventChargeSucceeded          = "charge.succeeded"
	EventChargeFailed             = "charge.failed"
)
