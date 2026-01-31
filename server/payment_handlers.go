package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hotel/models"
	"hotel/services"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	paymentService services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePaymentIntent creates a new payment intent for a booking
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get booking details (we need to inject the booking service)
	// For now, we'll create a simple response
	c.JSON(http.StatusOK, gin.H{
		"message":    "Payment intent creation not yet implemented",
		"booking_id": req.BookingID,
	})
}

// GetPaymentDetails retrieves payment details by ID
func (h *PaymentHandler) GetPaymentDetails(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentService.GetStripePaymentDetails(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment details retrieved successfully",
		"data":    payment,
	})
}

// HandleWebhook processes Stripe webhook events
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body", "details": err.Error()})
		return
	}

	signature := c.GetHeader("Stripe-Signature")

	if err := h.paymentService.ProcessWebhookEvent(body, signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook processing failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// RefundPayment processes a refund for a payment
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	var req models.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	refund, err := h.paymentService.RefundPayment(req.PaymentID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refund", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Refund processed successfully",
		"data":    refund,
	})
}

// GetBookingPayments retrieves all payments for a booking
func (h *PaymentHandler) GetBookingPayments(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID", "details": err.Error()})
		return
	}

	payments, err := h.paymentService.GetStripePaymentsByBookingID(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payments retrieved successfully",
		"data":    payments,
	})
}
