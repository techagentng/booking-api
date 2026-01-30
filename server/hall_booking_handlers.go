package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"hotel/models"
	"hotel/server/response"
	"hotel/services"

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

		response.JSON(c, "Hall booking status updated successfully", http.StatusOK, result, nil)
	}
}
