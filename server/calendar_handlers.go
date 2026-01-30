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

// handleGetCalendarAvailability returns calendar availability for a month
func (s *Server) handleGetCalendarAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, err := strconv.Atoi(c.Query("year"))
		if err != nil {
			log.Printf("handleGetCalendarAvailability: invalid year parameter: %v", err)
			response.JSON(c, "Invalid year parameter", http.StatusBadRequest, nil, err)
			return
		}

		month, err := strconv.Atoi(c.Query("month"))
		if err != nil || month < 1 || month > 12 {
			log.Printf("handleGetCalendarAvailability: invalid month parameter: %v", err)
			response.JSON(c, "Invalid month parameter", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		availabilities, err := calendarService.GetCalendarAvailability(year, time.Month(month))
		if err != nil {
			log.Printf("handleGetCalendarAvailability: error fetching calendar availability: %v", err)
			response.JSON(c, "Failed to fetch calendar availability", http.StatusInternalServerError, nil, err)
			return
		}

		// Convert to response format
		var responses []models.CalendarAvailabilityResponse
		for _, availability := range availabilities {
			responses = append(responses, models.CalendarAvailabilityResponse{
				ID:             availability.ID,
				Date:           availability.Date.Format("2006-01-02"),
				Status:         availability.Status,
				TotalSlots:     availability.TotalSlots,
				BookedSlots:    availability.BookedSlots,
				AvailableSlots: availability.AvailableSlots,
				Notes:          availability.Notes,
				CreatedAt:      availability.CreatedAt,
				UpdatedAt:      availability.UpdatedAt,
			})
		}

		response.JSON(c, "Calendar availability retrieved successfully", http.StatusOK, responses, nil)
	}
}

// handleGetDailyAvailability returns availability for a specific date
func (s *Server) handleGetDailyAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		log.Printf("DEBUG: handleGetDailyAvailability called for date %s", dateStr)

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("handleGetDailyAvailability: invalid date format: %v", err)
			response.JSON(c, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		log.Printf("DEBUG: About to call GetDailyAvailability")

		availability, err := calendarService.GetDailyAvailability(date)
		if err != nil {
			log.Printf("handleGetDailyAvailability: error fetching daily availability: %v", err)
			response.JSON(c, "Failed to fetch daily availability", http.StatusInternalServerError, nil, err)
			return
		}

		availabilityResponse := models.CalendarAvailabilityResponse{
			ID:             availability.ID,
			Date:           availability.Date.Format("2006-01-02"),
			Status:         availability.Status,
			TotalSlots:     availability.TotalSlots,
			BookedSlots:    availability.BookedSlots,
			AvailableSlots: availability.AvailableSlots,
			Notes:          availability.Notes,
			CreatedAt:      availability.CreatedAt,
			UpdatedAt:      availability.UpdatedAt,
		}

		response.JSON(c, "Daily availability retrieved successfully", http.StatusOK, availabilityResponse, nil)
	}
}

// handleGetTimeSlots returns available time slots for a date
func (s *Server) handleGetTimeSlots() gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("handleGetTimeSlots: invalid date format: %v", err)
			response.JSON(c, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		slots, err := calendarService.GetTimeSlots(date)
		if err != nil {
			log.Printf("handleGetTimeSlots: error fetching time slots: %v", err)
			response.JSON(c, "Failed to fetch time slots", http.StatusInternalServerError, nil, err)
			return
		}

		// Convert to response format
		var responses []models.TimeSlotResponse
		for _, slot := range slots {
			responses = append(responses, models.TimeSlotResponse{
				ID:           slot.ID,
				Date:         slot.Date.Format("2006-01-02"),
				StartTime:    slot.StartTime,
				EndTime:      slot.EndTime,
				Status:       slot.Status,
				MaxCapacity:  slot.MaxCapacity,
				CurrentUsage: slot.CurrentUsage,
				Price:        slot.Price,
				CreatedAt:    slot.CreatedAt,
				UpdatedAt:    slot.UpdatedAt,
			})
		}

		response.JSON(c, "Time slots retrieved successfully", http.StatusOK, responses, nil)
	}
}

// handleCheckSlotAvailability checks if a time slot is available
func (s *Server) handleCheckSlotAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Query("date")
		startTime := c.Query("start_time")
		endTime := c.Query("end_time")

		if dateStr == "" || startTime == "" || endTime == "" {
			log.Printf("handleCheckSlotAvailability: missing required parameters")
			response.JSON(c, "date, start_time, and end_time are required", http.StatusBadRequest, nil, nil)
			return
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("handleCheckSlotAvailability: invalid date format: %v", err)
			response.JSON(c, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		available, err := calendarService.CheckSlotAvailability(date, startTime, endTime)
		if err != nil {
			log.Printf("handleCheckSlotAvailability: error checking availability: %v", err)
			response.JSON(c, "Failed to check slot availability", http.StatusBadRequest, nil, err)
			return
		}

		result := models.TimeSlotAvailabilityResponse{
			Available: available,
			Date:      dateStr,
			StartTime: startTime,
			EndTime:   endTime,
		}

		if !available {
			result.Message = "Requested time slot is not available"
		} else {
			result.Message = "Requested time slot is available"
		}

		response.JSON(c, "Slot availability checked successfully", http.StatusOK, result, nil)
	}
}

// Admin endpoints (require authentication)

// handleUpdateDailyAvailability updates daily availability
func (s *Server) handleUpdateDailyAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("handleUpdateDailyAvailability: invalid date format: %v", err)
			response.JSON(c, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest, nil, err)
			return
		}

		var req models.CalendarUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleUpdateDailyAvailability: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		err = calendarService.UpdateDailyAvailability(date, req.Status, req.Notes)
		if err != nil {
			log.Printf("handleUpdateDailyAvailability: error updating availability: %v", err)
			response.JSON(c, "Failed to update daily availability", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Daily availability updated successfully", http.StatusOK, nil, nil)
	}
}

// handleGetMonthlyStats returns monthly statistics
func (s *Server) handleGetMonthlyStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, err := strconv.Atoi(c.Query("year"))
		if err != nil {
			log.Printf("handleGetMonthlyStats: invalid year parameter: %v", err)
			response.JSON(c, "Invalid year parameter", http.StatusBadRequest, nil, err)
			return
		}

		month, err := strconv.Atoi(c.Query("month"))
		if err != nil || month < 1 || month > 12 {
			log.Printf("handleGetMonthlyStats: invalid month parameter: %v", err)
			response.JSON(c, "Invalid month parameter", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		stats, err := calendarService.GetMonthlyStats(year, time.Month(month))
		if err != nil {
			log.Printf("handleGetMonthlyStats: error fetching monthly stats: %v", err)
			response.JSON(c, "Failed to fetch monthly statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Monthly statistics retrieved successfully", http.StatusOK, stats, nil)
	}
}

// handleGenerateCalendar generates calendar for a month
func (s *Server) handleGenerateCalendar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CalendarGenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleGenerateCalendar: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		err := calendarService.GenerateCalendar(req.Year, time.Month(req.Month))
		if err != nil {
			log.Printf("handleGenerateCalendar: error generating calendar: %v", err)
			response.JSON(c, "Failed to generate calendar", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Calendar generated successfully", http.StatusOK, nil, nil)
	}
}

// handleUpdateTimeSlotStatus updates a specific time slot status
func (s *Server) handleUpdateTimeSlotStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("handleUpdateTimeSlotStatus: invalid date format: %v", err)
			response.JSON(c, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest, nil, err)
			return
		}

		var req models.TimeSlotUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("handleUpdateTimeSlotStatus: invalid request body: %v", err)
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		calendarService := services.NewCalendarService(s.CalendarRepository, s.HallBookingRepository)
		err = calendarService.UpdateTimeSlotStatus(date, req.StartTime, req.EndTime, req.Status)
		if err != nil {
			log.Printf("handleUpdateTimeSlotStatus: error updating time slot: %v", err)
			response.JSON(c, "Failed to update time slot status", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Time slot status updated successfully", http.StatusOK, nil, nil)
	}
}
