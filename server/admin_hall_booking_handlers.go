package server

import (
	"hotel/models"
	"hotel/server/response"
	"hotel/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// handleGetAllAdminBookings retrieves all bookings with filtering and pagination
func (s *Server) handleGetAllAdminBookings() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Parse status filter
		var status []string
		if statusParam := c.Query("status"); statusParam != "" {
			status = []string{statusParam}
		} else if statusParams := c.QueryArray("status"); len(statusParams) > 0 {
			status = statusParams
		}

		dateFrom := c.Query("date_from")
		dateTo := c.Query("date_to")
		search := c.Query("search")
		sortBy := c.DefaultQuery("sort_by", "created_at")
		sortOrder := c.DefaultQuery("sort_order", "desc")

		// Create service
		adminService := services.NewAdminHallBookingService(s.AdminHallBookingRepository)

		// Get bookings
		bookings, total, err := adminService.GetAllBookings(page, pageSize, status, dateFrom, dateTo, search, sortBy, sortOrder)
		if err != nil {
			log.Printf("handleGetAllAdminBookings: error getting bookings: %v", err)
			response.JSON(c, "Failed to retrieve bookings", http.StatusInternalServerError, nil, err)
			return
		}

		// Convert to admin responses
		var adminResponses []models.AdminBookingResponse
		for _, booking := range bookings {
			adminResponses = append(adminResponses, *adminService.ConvertToAdminResponse(&booking))
		}

		// Create paginated response
		meta := models.PaginationMeta{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.JSON(c, "Bookings retrieved successfully", http.StatusOK, map[string]interface{}{
			"data": adminResponses,
			"meta": meta,
		}, nil)
	}
}

// handleGetAdminBookingByID retrieves a booking by ID
func (s *Server) handleGetAdminBookingByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleGetAdminBookingByID: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid booking ID", http.StatusBadRequest, nil, err)
			return
		}

		adminService := services.NewAdminHallBookingService(s.AdminHallBookingRepository)
		booking, err := adminService.GetBookingByID(uint(id))
		if err != nil {
			log.Printf("handleGetAdminBookingByID: error getting booking: %v", err)
			response.JSON(c, "Booking not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Booking retrieved successfully", http.StatusOK, adminService.ConvertToAdminResponse(booking), nil)
	}
}

// handleUpdateAdminBookingStatus updates booking status with admin tracking
func (s *Server) handleUpdateAdminBookingStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleUpdateAdminBookingStatus: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid booking ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.AdminBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleUpdateAdminBookingStatus: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		// Get admin user ID from context (should be set by auth middleware)
		adminID := uint(1) // TODO: Get from auth context

		adminService := services.NewAdminHallBookingService(s.AdminHallBookingRepository)
		booking, err := adminService.UpdateBookingStatus(uint(id), &req, adminID)
		if err != nil {
			log.Printf("handleUpdateAdminBookingStatus: error updating status: %v", err)
			response.JSON(c, "Failed to update booking status", http.StatusBadRequest, nil, err)
			return
		}

		// Send real-time notification for status update
		if req.Status == "cancelled" {
			s.NotificationHub.NotifyHallBookingCancelled(
				booking.ID,
				booking.BookingID,
				booking.OrganizerName,
				booking.EventType,
				booking.BookingDate.Format("2006-01-02"),
			)
		} else {
			s.NotificationHub.NotifyHallBookingUpdated(
				booking.ID,
				booking.BookingID,
				booking.OrganizerName,
				booking.EventType,
				req.Status,
			)
		}

		// Send email notification
		if s.MailService != nil && booking.OrganizerEmail != "" {
			go func() {
				var err error
				if req.Status == "cancelled" {
					_, err = s.MailService.SendBookingCancellation(
						booking.OrganizerEmail,
						booking.OrganizerName,
						booking.BookingID,
						booking.EventType,
						booking.BookingDate.Format("2006-01-02"),
					)
				} else {
					_, err = s.MailService.SendBookingStatusUpdate(
						booking.OrganizerEmail,
						booking.OrganizerName,
						booking.BookingID,
						booking.EventType,
						req.Status,
					)
				}
				if err != nil {
					log.Printf("Failed to send booking status update email: %v", err)
				}
			}()
		}

		response.JSON(c, "Booking status updated successfully", http.StatusOK, adminService.ConvertToAdminResponse(booking), nil)
	}
}

// handleGetAdminBookingStatusHistory retrieves status history for a booking
func (s *Server) handleGetAdminBookingStatusHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Printf("handleGetAdminBookingStatusHistory: invalid ID parameter: %v", err)
			response.JSON(c, "Invalid booking ID", http.StatusBadRequest, nil, err)
			return
		}

		adminService := services.NewAdminHallBookingService(s.AdminHallBookingRepository)
		history, err := adminService.GetBookingStatusHistory(uint(id))
		if err != nil {
			log.Printf("handleGetAdminBookingStatusHistory: error getting history: %v", err)
			response.JSON(c, "Failed to retrieve booking history", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Booking history retrieved successfully", http.StatusOK, map[string]interface{}{
			"data": history,
		}, nil)
	}
}

// handleGetAdminBookingStats retrieves booking statistics
func (s *Server) handleGetAdminBookingStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		period := c.DefaultQuery("period", "month")
		dateFrom := c.Query("date_from")
		dateTo := c.Query("date_to")

		adminService := services.NewAdminHallBookingService(s.AdminHallBookingRepository)
		stats, err := adminService.GetBookingStats(period, dateFrom, dateTo)
		if err != nil {
			log.Printf("handleGetAdminBookingStats: error getting stats: %v", err)
			response.JSON(c, "Failed to retrieve booking statistics", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Booking statistics retrieved successfully", http.StatusOK, map[string]interface{}{
			"data": stats,
		}, nil)
	}
}

// handleGetAdminCalendarAvailability retrieves calendar availability with booking details
func (s *Server) handleGetAdminCalendarAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Query("year"))
		month, _ := strconv.Atoi(c.Query("month"))
		includeBookings := c.DefaultQuery("include_bookings", "false") == "true"

		if year == 0 || month == 0 {
			log.Printf("handleGetAdminCalendarAvailability: missing year or month parameters")
			response.JSON(c, "Year and month parameters are required", http.StatusBadRequest, nil, nil)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		availabilities, err := calendarService.GetCalendarAvailability(year, time.Month(month))
		if err != nil {
			log.Printf("handleGetAdminCalendarAvailability: error getting calendar: %v", err)
			response.JSON(c, "Failed to retrieve calendar availability", http.StatusInternalServerError, nil, err)
			return
		}

		// If include_bookings is true, enrich with booking details
		if includeBookings {
			adminService := services.NewAdminHallBookingService(s.AdminHallBookingRepository)

			// Get all bookings for the month
			startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endOfMonth := startOfMonth.AddDate(0, 1, -1)

			bookings, _, err := adminService.GetAllBookings(1, 1000, []string{},
				startOfMonth.Format("2006-01-02"),
				endOfMonth.Format("2006-01-02"),
				"", "date", "asc")
			if err == nil {
				// Group bookings by date
				bookingsByDate := make(map[string][]models.HallBooking)
				for _, booking := range bookings {
					date := booking.BookingDate.Format("2006-01-02")
					bookingsByDate[date] = append(bookingsByDate[date], booking)
				}

				// Enrich availabilities with booking details
				for i, availability := range availabilities {
					date := availability.Date.Format("2006-01-02")
					if dateBookings, exists := bookingsByDate[date]; exists {
						// Convert bookings to admin responses
						var bookingResponses []models.AdminBookingResponse
						for _, booking := range dateBookings {
							bookingResponses = append(bookingResponses, *adminService.ConvertToAdminResponse(&booking))
						}
						// Add bookings to availability (you might need to extend the model)
						// For now, we'll just log it
						log.Printf("Date %s has %d bookings", date, len(dateBookings))
					}
					availabilities[i] = availability
				}
			}
		}

		response.JSON(c, "Admin calendar availability retrieved successfully", http.StatusOK, map[string]interface{}{
			"data": availabilities,
		}, nil)
	}
}
