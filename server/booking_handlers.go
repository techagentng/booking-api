package server

import (
	"net/http"

	"hotel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateServiceBooking handles creating a service booking
func (s *Server) CreateServiceBooking(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Get customer profile
	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer profile not found"})
		return
	}

	var req models.CreateServiceBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Get provider_id from service
	providerID, err := s.ProviderRepository.GetProviderIDByServiceID(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Service not found"})
		return
	}

	booking, err := s.BookingRepository.CreateServiceBooking(customer.ID, providerID, req.ServiceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    booking,
		"message": "Booking created successfully",
	})
}

// GetCustomerServiceBookings handles getting customer's service bookings
func (s *Server) GetCustomerServiceBookings(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer profile not found"})
		return
	}

	status := c.Query("status")
	bookings, err := s.BookingRepository.GetCustomerBookings(customer.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
		"message": "Customer bookings retrieved",
	})
}

// GetProviderServiceBookings handles getting provider's service bookings
func (s *Server) GetProviderServiceBookings(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider profile not found"})
		return
	}

	status := c.Query("status")
	bookings, err := s.BookingRepository.GetProviderBookings(provider.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
		"message": "Provider bookings retrieved",
	})
}

// CancelServiceBooking handles cancelling a booking
func (s *Server) CancelServiceBooking(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid booking ID"})
		return
	}

	// Verify the booking belongs to the customer
	booking, err := s.BookingRepository.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking not found"})
		return
	}

	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil || booking.CustomerID != customer.ID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to cancel this booking"})
		return
	}

	if err := s.BookingRepository.CancelBooking(bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking cancelled successfully",
	})
}

// AcceptServiceBooking handles accepting a booking (provider action)
func (s *Server) AcceptServiceBooking(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid booking ID"})
		return
	}

	// Verify the booking belongs to the provider
	booking, err := s.BookingRepository.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking not found"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil || booking.ProviderID != provider.ID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to accept this booking"})
		return
	}

	if err := s.BookingRepository.AcceptBooking(bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking accepted successfully",
	})
}

// RejectServiceBooking handles rejecting a booking (provider action)
func (s *Server) RejectServiceBooking(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid booking ID"})
		return
	}

	// Verify the booking belongs to the provider
	booking, err := s.BookingRepository.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking not found"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil || booking.ProviderID != provider.ID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to reject this booking"})
		return
	}

	if err := s.BookingRepository.RejectBooking(bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking rejected successfully",
	})
}

// CompleteServiceBooking handles completing a booking (provider action)
func (s *Server) CompleteServiceBooking(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid booking ID"})
		return
	}

	// Verify the booking belongs to the provider
	booking, err := s.BookingRepository.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking not found"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil || booking.ProviderID != provider.ID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to complete this booking"})
		return
	}

	if err := s.BookingRepository.CompleteBooking(bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking completed successfully",
	})
}
