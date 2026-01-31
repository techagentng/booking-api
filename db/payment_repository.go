package db

import (
	"fmt"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// PaymentRepository defines payment database operations
type PaymentRepository interface {
	CreateStripePayment(payment *models.StripePayment) error
	GetStripePaymentByID(id string) (*models.StripePayment, error)
	GetStripePaymentsByBookingID(bookingID uint) ([]models.StripePayment, error)
	UpdateStripePaymentStatus(id, status string) error
	CreateStripePaymentIntent(intent *models.StripePaymentIntent) error
	GetStripePaymentIntentByBookingID(bookingID uint) (*models.StripePaymentIntent, error)
	UpdateStripePaymentIntent(id, status string) error
	CreateWebhookEvent(event *models.WebhookEvent) error
	GetUnprocessedWebhookEvents() ([]models.WebhookEvent, error)
	MarkWebhookEventProcessed(id string) error
	GetStripePaymentIntentByID(id string) (*models.StripePaymentIntent, error)
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// CreateStripePayment creates a new Stripe payment record
func (r *paymentRepository) CreateStripePayment(payment *models.StripePayment) error {
	return r.db.Create(payment).Error
}

// GetStripePaymentByID retrieves a Stripe payment by ID
func (r *paymentRepository) GetStripePaymentByID(id string) (*models.StripePayment, error) {
	var payment models.StripePayment
	err := r.db.Preload("Booking").First(&payment, "id = ?", id).Error
	return &payment, err
}

// GetStripePaymentsByBookingID retrieves all Stripe payments for a booking
func (r *paymentRepository) GetStripePaymentsByBookingID(bookingID uint) ([]models.StripePayment, error) {
	var payments []models.StripePayment
	err := r.db.Where("booking_id = ?", bookingID).Find(&payments).Error
	return payments, err
}

// UpdateStripePaymentStatus updates the status of a Stripe payment
func (r *paymentRepository) UpdateStripePaymentStatus(id, status string) error {
	return r.db.Model(&models.StripePayment{}).Where("id = ?", id).Update("status", status).Error
}

// CreateStripePaymentIntent creates a new Stripe payment intent record
func (r *paymentRepository) CreateStripePaymentIntent(intent *models.StripePaymentIntent) error {
	return r.db.Create(intent).Error
}

// GetStripePaymentIntentByBookingID retrieves a Stripe payment intent by booking ID
func (r *paymentRepository) GetStripePaymentIntentByBookingID(bookingID uint) (*models.StripePaymentIntent, error) {
	var intent models.StripePaymentIntent
	err := r.db.Where("booking_id = ?", bookingID).First(&intent).Error
	return &intent, err
}

// GetStripePaymentIntentByID retrieves a Stripe payment intent by ID
func (r *paymentRepository) GetStripePaymentIntentByID(id string) (*models.StripePaymentIntent, error) {
	var intent models.StripePaymentIntent
	err := r.db.Preload("Booking").First(&intent, "id = ?", id).Error
	return &intent, err
}

// UpdateStripePaymentIntent updates the status of a Stripe payment intent
func (r *paymentRepository) UpdateStripePaymentIntent(id, status string) error {
	return r.db.Model(&models.StripePaymentIntent{}).Where("id = ?", id).Update("status", status).Error
}

// CreateWebhookEvent creates a new webhook event record
func (r *paymentRepository) CreateWebhookEvent(event *models.WebhookEvent) error {
	return r.db.Create(event).Error
}

// GetUnprocessedWebhookEvents retrieves all unprocessed webhook events
func (r *paymentRepository) GetUnprocessedWebhookEvents() ([]models.WebhookEvent, error) {
	var events []models.WebhookEvent
	err := r.db.Where("processed = ?", false).Find(&events).Error
	return events, err
}

// MarkWebhookEventProcessed marks a webhook event as processed
func (r *paymentRepository) MarkWebhookEventProcessed(id string) error {
	return r.db.Model(&models.WebhookEvent{}).Where("id = ?", id).Update("processed", true).Error
}

// generateID generates a unique ID for webhook events
func generateID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
