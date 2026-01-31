package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"

	"hotel/config"
	"hotel/db"
	"hotel/models"
)

// PaymentService defines payment business logic operations
type PaymentService interface {
	CreatePaymentIntent(booking *models.HallBooking) (*stripe.PaymentIntent, error)
	ProcessWebhookEvent(payload []byte, signature string) error
	GetStripePaymentDetails(paymentID string) (*models.StripePayment, error)
	RefundPayment(paymentID string, amount int64) (*stripe.Refund, error)
	GetStripePaymentsByBookingID(bookingID uint) ([]models.StripePayment, error)
}

type paymentService struct {
	paymentRepo db.PaymentRepository
	bookingRepo db.HallBookingRepository
	config      *config.Config
}

// NewPaymentService creates a new payment service
func NewPaymentService(paymentRepo db.PaymentRepository, bookingRepo db.HallBookingRepository, cfg *config.Config) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		bookingRepo: bookingRepo,
		config:      cfg,
	}
}

// InitStripe initializes Stripe with the secret key
func InitStripe(cfg *config.Config) {
	stripeKey := cfg.StripeSecretKey
	if stripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY environment variable is not set")
	}
	stripe.Key = stripeKey
}

// GetFrontendURL returns the frontend URL from environment
func GetFrontendURL(cfg *config.Config) string {
	frontendURL := cfg.FRONTEND_URL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	return frontendURL
}

// GetWebhookSecret returns the Stripe webhook secret from environment
func GetWebhookSecret(cfg *config.Config) string {
	return cfg.StripeWebhookSecret
}

// CreatePaymentIntent creates a Stripe payment intent for a booking
func (s *paymentService) CreatePaymentIntent(booking *models.HallBooking) (*stripe.PaymentIntent, error) {
	amount := int64(booking.TotalPrice * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("gbp"),
		Metadata: map[string]string{
			"booking_id":  strconv.FormatUint(uint64(booking.ID), 10),
			"booking_ref": booking.BookingID,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Store payment intent
	paymentIntent := &models.StripePaymentIntent{
		ID:           intent.ID,
		BookingID:    booking.ID,
		Amount:       amount,
		Currency:     "gbp",
		Status:       string(intent.Status),
		ClientSecret: intent.ClientSecret,
	}

	if err := s.paymentRepo.CreateStripePaymentIntent(paymentIntent); err != nil {
		log.Printf("Warning: Failed to store payment intent: %v", err)
	}

	return intent, nil
}

// ProcessWebhookEvent processes a Stripe webhook event
func (s *paymentService) ProcessWebhookEvent(payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, GetWebhookSecret(s.config))
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	// Store webhook event for processing
	webhookEvent := &models.WebhookEvent{
		ID:            generateID(),
		StripeEventID: string(event.ID),
		EventType:     string(event.Type),
		Processed:     false,
		Payload:       string(payload),
	}

	if err := s.paymentRepo.CreateWebhookEvent(webhookEvent); err != nil {
		return fmt.Errorf("failed to store webhook event: %w", err)
	}

	// Process the event
	return s.processStripeEvent(event)
}

// processStripeEvent processes a Stripe event
func (s *paymentService) processStripeEvent(event stripe.Event) error {
	switch event.Type {
	case models.EventPaymentIntentSucceeded:
		return s.handlePaymentIntentSucceeded(event)
	case models.EventPaymentIntentFailed:
		return s.handlePaymentIntentFailed(event)
	case models.EventPaymentIntentCanceled:
		return s.handlePaymentIntentCanceled(event)
	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
		return nil
	}
}

// handlePaymentIntentSucceeded handles payment intent success
func (s *paymentService) handlePaymentIntentSucceeded(event stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	// Update payment intent record
	if err := s.paymentRepo.UpdateStripePaymentIntent(string(intent.ID), string(intent.Status)); err != nil {
		log.Printf("Warning: Failed to update payment intent: %v", err)
	}

	// Create payment record
	bookingID, err := strconv.ParseUint(intent.Metadata["booking_id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid booking_id in metadata: %w", err)
	}

	payment := &models.StripePayment{
		ID:        string(intent.ID),
		BookingID: uint(bookingID),
		Amount:    intent.Amount,
		Currency:  string(intent.Currency),
		Status:    models.PaymentStatusSucceeded,
	}

	// Get customer info if available
	if intent.Customer != nil {
		payment.CustomerID = string(intent.Customer.ID)
	}

	if err := s.paymentRepo.CreateStripePayment(payment); err != nil {
		return fmt.Errorf("failed to create payment record: %w", err)
	}

	// Update booking status
	return s.updateBookingAfterPayment(uint(bookingID), "confirmed")
}

// handlePaymentIntentFailed handles payment intent failure
func (s *paymentService) handlePaymentIntentFailed(event stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	// Update payment intent record
	if err := s.paymentRepo.UpdateStripePaymentIntent(string(intent.ID), string(intent.Status)); err != nil {
		log.Printf("Warning: Failed to update payment intent: %v", err)
	}

	// Log payment failure
	if intent.LastPaymentError != nil {
		log.Printf("Payment failed for intent %s: %v", string(intent.ID), intent.LastPaymentError)
	}

	return nil
}

// handlePaymentIntentCanceled handles payment intent cancellation
func (s *paymentService) handlePaymentIntentCanceled(event stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	// Update payment intent record
	if err := s.paymentRepo.UpdateStripePaymentIntent(intent.ID, string(intent.Status)); err != nil {
		log.Printf("Warning: Failed to update payment intent: %v", err)
	}

	return nil
}

// updateBookingAfterPayment updates booking status after successful payment
func (s *paymentService) updateBookingAfterPayment(bookingID uint, status string) error {
	booking, err := s.bookingRepo.GetHallBookingByID(bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	booking.Status = status
	booking.UpdatedAt = time.Now()

	_, err = s.bookingRepo.UpdateHallBooking(bookingID, booking)
	return err
}

// GetStripePaymentDetails retrieves payment details by ID
func (s *paymentService) GetStripePaymentDetails(paymentID string) (*models.StripePayment, error) {
	return s.paymentRepo.GetStripePaymentByID(paymentID)
}

// RefundPayment processes a refund for a payment
func (s *paymentService) RefundPayment(paymentID string, amount int64) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentID),
	}

	if amount > 0 {
		params.Amount = stripe.Int64(amount)
	}

	refund, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// Update payment status if full refund
	if amount <= 0 || amount >= refund.Amount {
		if err := s.paymentRepo.UpdateStripePaymentStatus(paymentID, models.PaymentStatusRefunded); err != nil {
			log.Printf("Warning: Failed to update payment status after refund: %v", err)
		}
	}

	return refund, nil
}

// GetStripePaymentsByBookingID retrieves all payments for a booking
func (s *paymentService) GetStripePaymentsByBookingID(bookingID uint) ([]models.StripePayment, error) {
	return s.paymentRepo.GetStripePaymentsByBookingID(bookingID)
}

// generateID generates a unique ID for webhook events
func generateID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
