package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hotel/models"
	"hotel/server/response"
	"hotel/services"
	"hotel/services/jwt"

	"github.com/gin-gonic/gin"
)

// handleCreateHallBooking creates a new hall booking
func (s *Server) handleCreateHallBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateHallBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleCreateHallBooking: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.CreateHallBooking(&req)
		if err != nil {
			log.Printf("handleCreateHallBooking: error creating hall booking: %v", err)
			response.JSON(c, "Failed to create hall booking", http.StatusBadRequest, nil, err)
			return
		}

		// Update calendar availability after booking is created
		bookingDate, _ := time.Parse("2006-01-02", req.BookingDate)
		if err := s.CalendarRepository.UpdateAvailabilityFromBookings(bookingDate); err != nil {
			log.Printf("handleCreateHallBooking: error updating calendar: %v", err)
			// Don't fail the request, just log the error
		}

		response.JSON(c, "Hall booking created successfully", http.StatusCreated, result, nil)
	}
}

// handleGetHallBookings returns all hall bookings with pagination and search
func (s *Server) handleGetHallBookings() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		pageSize := 10
		search := c.Query("search")

		if p := c.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				page = parsed
			}
		}

		if ps := c.Query("page_size"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
				pageSize = parsed
			}
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.GetAllHallBookings(page, pageSize, search)
		if err != nil {
			log.Printf("handleGetHallBookings: error fetching hall bookings: %v", err)
			response.JSON(c, "Failed to fetch hall bookings", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Hall bookings retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetHallBookingByID returns a hall booking by ID
func (s *Server) handleGetHallBookingByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleGetHallBookingByID: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid hall booking ID", http.StatusBadRequest, nil, err)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.GetHallBookingByID(uint(id))
		if err != nil {
			log.Printf("handleGetHallBookingByID: error fetching hall booking: %v", err)
			response.JSON(c, "Failed to fetch hall booking", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Hall booking retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetHallBookingByBookingID returns a hall booking by booking ID
func (s *Server) handleGetHallBookingByBookingID() gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingID := c.Param("booking_id")
		if bookingID == "" {
			log.Printf("handleGetHallBookingByBookingID: missing booking ID parameter")
			response.JSON(c, "Missing booking ID", http.StatusBadRequest, nil, nil)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.GetHallBookingByBookingID(bookingID)
		if err != nil {
			log.Printf("handleGetHallBookingByBookingID: error fetching hall booking: %v", err)
			response.JSON(c, "Failed to fetch hall booking", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Hall booking retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleUpdateHallBooking updates a hall booking
func (s *Server) handleUpdateHallBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleUpdateHallBooking: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid hall booking ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateHallBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleUpdateHallBooking: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.UpdateHallBooking(uint(id), &req)
		if err != nil {
			log.Printf("handleUpdateHallBooking: error updating hall booking: %v", err)
			response.JSON(c, "Failed to update hall booking", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Hall booking updated successfully", http.StatusOK, result, nil)
	}
}

// handleDeleteHallBooking deletes a hall booking
func (s *Server) handleDeleteHallBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleDeleteHallBooking: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid hall booking ID", http.StatusBadRequest, nil, err)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		err = hallBookingService.DeleteHallBooking(uint(id))
		if err != nil {
			log.Printf("handleDeleteHallBooking: error deleting hall booking: %v", err)
			response.JSON(c, "Failed to delete hall booking", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Hall booking deleted successfully", http.StatusOK, nil, nil)
	}
}

// handleCheckHallAvailability checks hall availability for a specific date and time
func (s *Server) handleCheckHallAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Query("date")
		startTime := c.Query("start_time")
		endTime := c.Query("end_time")

		if date == "" || startTime == "" || endTime == "" {
			log.Printf("handleCheckHallAvailability: missing required parameters")
			response.JSON(c, "Missing required parameters: date, start_time, end_time", http.StatusBadRequest, nil, nil)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		available, err := hallBookingService.CheckHallAvailability(date, startTime, endTime)
		if err != nil {
			log.Printf("handleCheckHallAvailability: error checking availability: %v", err)
			response.JSON(c, "Failed to check hall availability", http.StatusBadRequest, nil, err)
			return
		}

		result := gin.H{
			"available":  available,
			"date":       date,
			"start_time": startTime,
			"end_time":   endTime,
		}

		response.JSON(c, "Hall availability checked successfully", http.StatusOK, result, nil)
	}
}

// handleGetHallAvailability gets hall availability for a specific date
func (s *Server) handleGetHallAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		if date == "" {
			log.Printf("handleGetHallAvailability: missing date parameter")
			response.JSON(c, "Missing date parameter", http.StatusBadRequest, nil, nil)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.GetHallAvailability(date)
		if err != nil {
			log.Printf("handleGetHallAvailability: error getting hall availability: %v", err)
			response.JSON(c, "Failed to get hall availability", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Hall availability retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetHallBookingsByDate returns hall bookings for a specific date
func (s *Server) handleGetHallBookingsByDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		if date == "" {
			log.Printf("handleGetHallBookingsByDate: missing date parameter")
			response.JSON(c, "Missing date parameter", http.StatusBadRequest, nil, nil)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.GetHallBookingsByDate(date)
		if err != nil {
			log.Printf("handleGetHallBookingsByDate: error fetching hall bookings: %v", err)
			response.JSON(c, "Failed to fetch hall bookings", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Hall bookings retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetHallBookingsByDateRange returns hall bookings within a date range
func (s *Server) handleGetHallBookingsByDateRange() gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			log.Printf("handleGetHallBookingsByDateRange: missing required parameters")
			response.JSON(c, "Missing required parameters: start_date, end_date", http.StatusBadRequest, nil, nil)
			return
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.GetHallBookingsByDateRange(startDate, endDate)
		if err != nil {
			log.Printf("handleGetHallBookingsByDateRange: error fetching hall bookings: %v", err)
			response.JSON(c, "Failed to fetch hall bookings", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Hall bookings retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleUpdateHallBookingStatus updates a hall booking status
func (s *Server) handleUpdateHallBookingStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleUpdateHallBookingStatus: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid hall booking ID", http.StatusBadRequest, nil, err)
			return
		}

		var req struct {
			Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleUpdateHallBookingStatus: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		updateReq := &models.UpdateHallBookingRequest{
			Status: &req.Status,
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.UpdateHallBooking(uint(id), updateReq)
		if err != nil {
			log.Printf("handleUpdateHallBookingStatus: error updating hall booking status: %v", err)
			response.JSON(c, "Failed to update hall booking status", http.StatusBadRequest, nil, err)
			return
		}

		// Send real-time notification for status update
		if req.Status == "cancelled" {
			s.NotificationHub.NotifyHallBookingCancelled(
				result.ID,
				result.BookingID,
				result.OrganizerName,
				result.EventType,
				result.BookingDate,
			)
		} else {
			s.NotificationHub.NotifyHallBookingUpdated(
				result.ID,
				result.BookingID,
				result.OrganizerName,
				result.EventType,
				req.Status,
			)
		}

		// Send email notification
		if s.MailService != nil && result.OrganizerEmail != "" {
			go func() {
				var err error
				if req.Status == "cancelled" {
					_, err = s.MailService.SendBookingCancellation(
						result.OrganizerEmail,
						result.OrganizerName,
						result.BookingID,
						result.EventType,
						result.BookingDate,
					)
				} else {
					_, err = s.MailService.SendBookingStatusUpdate(
						result.OrganizerEmail,
						result.OrganizerName,
						result.BookingID,
						result.EventType,
						req.Status,
					)
				}
				if err != nil {
					log.Printf("Failed to send booking status update email: %v", err)
				}
			}()
		}

		response.JSON(c, "Hall booking status updated successfully", http.StatusOK, result, nil)
	}
}

// authenticateUser determines the user type and ID from the request
func (s *Server) authenticateUser(c *gin.Context) (*uint, string) {
	// Check for JWT token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, "public"
	}

	// Validate token format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, "public"
	}

	// Extract and validate token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	accessClaims, err := jwt.ValidateAndGetClaims(token, s.Config.JWTSecret)
	if err != nil {
		return nil, "public"
	}

	// Check if user is admin or regular user
	isAdmin, ok := accessClaims["is_admin"].(bool)
	if !ok {
		return nil, "public"
	}

	userIDFloat, ok := accessClaims["id"].(float64)
	if !ok {
		return nil, "public"
	}
	userID := uint(userIDFloat)

	if isAdmin {
		return &userID, "admin"
	}

	return &userID, "user"
}

// handleCreateUnifiedBooking creates a new booking for both authenticated and public users
func (s *Server) handleCreateUnifiedBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to authenticate user
		userID, userType := s.authenticateUser(c)

		// Parse booking data (same for both user types)
		var req models.CreateHallBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleCreateUnifiedBooking: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		// Set creator information based on user type
		switch userType {
		case "admin":
			req.CreatedBy = userID
			req.CreatedByType = "admin"
			req.CreatorEmail = req.OrganizerEmail
			req.CreatorName = req.OrganizerName
		case "user":
			req.CreatedBy = userID
			req.CreatedByType = "user"
			req.CreatorEmail = req.OrganizerEmail
			req.CreatorName = req.OrganizerName
		default: // public/non-authenticated
			req.CreatedBy = nil
			req.CreatedByType = "public"
			req.CreatorEmail = req.OrganizerEmail
			req.CreatorName = req.OrganizerName
		}

		// Create booking (same logic for all)
		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)
		result, err := hallBookingService.CreateHallBooking(&req)
		if err != nil {
			log.Printf("handleCreateUnifiedBooking: error creating hall booking: %v", err)
			response.JSON(c, "Failed to create hall booking", http.StatusBadRequest, nil, err)
			return
		}

		// Send real-time notification for new hall booking
		s.NotificationHub.NotifyNewHallBooking(
			result.ID,
			result.BookingID,
			result.OrganizerName,
			result.EventType,
			result.BookingDate,
			result.GuestCount,
			result.TotalPrice,
			result.CreatedByType,
		)

		// Send email confirmation
		if s.MailService != nil && result.OrganizerEmail != "" {
			go func() {
				_, err := s.MailService.SendBookingConfirmation(
					result.OrganizerEmail,
					result.OrganizerName,
					result.BookingID,
					result.EventType,
					result.BookingDate,
					result.GuestCount,
					result.TotalPrice,
				)
				if err != nil {
					log.Printf("Failed to send booking confirmation email: %v", err)
				}
			}()
		}

		response.JSON(c, "Hall booking created successfully", http.StatusCreated, result, nil)
	}
}

// handleGetHallBookingRecentActivity retrieves recent hall booking activity for dashboard
func (s *Server) handleGetHallBookingRecentActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get limit from query param (default to 10)
		limitStr := c.DefaultQuery("limit", "10")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
		if limit > 50 { // Cap at 50 to prevent excessive data
			limit = 50
		}

		hallBookingService := services.NewHallBookingService(s.HallBookingRepository)

		// Get recent bookings
		recentBookings, err := hallBookingService.GetRecentBookings(limit)
		if err != nil {
			log.Printf("handleGetHallBookingRecentActivity: error getting recent bookings: %v", err)
			response.JSON(c, "Failed to fetch recent activity", http.StatusInternalServerError, nil, err)
			return
		}

		// Get total count
		totalCount, err := hallBookingService.GetTotalBookingsCount()
		if err != nil {
			log.Printf("handleGetHallBookingRecentActivity: error getting total count: %v", err)
			response.JSON(c, "Failed to get total count", http.StatusInternalServerError, nil, err)
			return
		}

		activityResponse := models.HallBookingActivityResponse{
			RecentBookings: recentBookings,
			TotalCount:     totalCount,
		}

		response.JSON(c, "Recent activity retrieved successfully", http.StatusOK, activityResponse, nil)
	}
}

// Payment handler methods

// handleCreatePaymentIntent creates a payment intent for a booking
func (s *Server) handleCreatePaymentIntent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
			return
		}

		// Get booking details
		booking, err := s.HallBookingRepository.GetHallBookingByID(req.BookingID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "details": err.Error()})
			return
		}

		// Create payment service
		paymentService := services.NewPaymentService(s.PaymentRepository, s.HallBookingRepository, s.Config)

		// Create payment intent
		intent, err := paymentService.CreatePaymentIntent(booking)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent", "details": err.Error()})
			return
		}

		response := models.PaymentIntentResponse{
			ClientSecret: intent.ClientSecret,
			PaymentID:    intent.ID,
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Payment intent created successfully",
			"data":    response,
		})
	}
}

// handleGetStripePaymentDetails retrieves Stripe payment details
func (s *Server) handleGetStripePaymentDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("id")

		paymentService := services.NewPaymentService(s.PaymentRepository, s.HallBookingRepository, s.Config)
		payment, err := paymentService.GetStripePaymentDetails(paymentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Payment details retrieved successfully",
			"data":    payment,
		})
	}
}

// handleStripeWebhook processes Stripe webhook events
func (s *Server) handleStripeWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body", "details": err.Error()})
			return
		}

		signature := c.GetHeader("Stripe-Signature")

		paymentService := services.NewPaymentService(s.PaymentRepository, s.HallBookingRepository, s.Config)
		if err := paymentService.ProcessWebhookEvent(body, signature); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook processing failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
}

// handleRefundPayment processes a refund
func (s *Server) handleRefundPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RefundRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
			return
		}

		paymentService := services.NewPaymentService(s.PaymentRepository, s.HallBookingRepository, s.Config)
		refund, err := paymentService.RefundPayment(req.PaymentID, req.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refund", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Refund processed successfully",
			"data":    refund,
		})
	}
}

// handleGetBookingPayments retrieves payments for a booking
func (s *Server) handleGetBookingPayments() gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingIDStr := c.Param("booking_id")
		bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID", "details": err.Error()})
			return
		}

		paymentService := services.NewPaymentService(s.PaymentRepository, s.HallBookingRepository, s.Config)
		payments, err := paymentService.GetStripePaymentsByBookingID(uint(bookingID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Payments retrieved successfully",
			"data":    payments,
		})
	}
}
