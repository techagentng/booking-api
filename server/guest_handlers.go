package server

import (
	"log"
	"net/http"
	"strconv"

	"hotel/models"
	"hotel/server/response"
	"hotel/services"

	"github.com/gin-gonic/gin"
)

// handleGetGuests returns all guests with pagination and search
func (s *Server) handleGetGuests() gin.HandlerFunc {
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

		guestService := services.NewGuestService(s.GuestRepository)
		result, err := guestService.GetAllGuests(page, pageSize, search)
		if err != nil {
			log.Printf("handleGetGuests: error fetching guests: %v", err)
			response.JSON(c, "Failed to fetch guests", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Guests retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetGuestByID returns a guest by ID with full details, statistics, and service usage
func (s *Server) handleGetGuestByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestService(s.GuestRepository)
		guestDetails, err := guestService.GetGuestDetailsByID(uint(id))
		if err != nil {
			log.Printf("handleGetGuestByID: error fetching guest: %v", err)
			response.JSON(c, "Guest not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Guest retrieved successfully", http.StatusOK, guestDetails, nil)
	}
}

// handleCreateGuest creates a new guest
func (s *Server) handleCreateGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateGuestRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestService(s.GuestRepository)
		guest, err := guestService.CreateGuest(&req)
		if err != nil {
			log.Printf("handleCreateGuest: error creating guest: %v", err)
			if err.Error() == "email already exists" {
				response.JSON(c, "Guest creation failed", http.StatusConflict, nil, err)
				return
			}
			response.JSON(c, "Guest creation failed", http.StatusBadRequest, nil, err)
			return
		}

		// Send new guest notification
		s.NotificationHub.NotifyNewGuest(guest.ID, guest.Name, "", guest.Email, guest.Phone)

		response.JSON(c, "Guest created successfully", http.StatusCreated, guest, nil)
	}
}

// handleUpdateGuest updates a guest
func (s *Server) handleUpdateGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateGuestRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestService(s.GuestRepository)
		guest, err := guestService.UpdateGuest(uint(id), &req)
		if err != nil {
			log.Printf("handleUpdateGuest: error updating guest: %v", err)
			if err.Error() == "email already exists" {
				response.JSON(c, "Guest update failed", http.StatusConflict, nil, err)
				return
			}
			response.JSON(c, "Guest update failed", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Guest updated successfully", http.StatusOK, guest, nil)
	}
}

// handleDeleteGuest deletes a guest
func (s *Server) handleDeleteGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestService(s.GuestRepository)
		err = guestService.DeleteGuest(uint(id))
		if err != nil {
			log.Printf("handleDeleteGuest: error deleting guest: %v", err)
			response.JSON(c, "Guest deletion failed", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Guest deleted successfully", http.StatusOK, nil, nil)
	}
}

// handleGetGuestHistory returns guest history
func (s *Server) handleGetGuestHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestService(s.GuestRepository)
		history, err := guestService.GetGuestHistory(uint(id))
		if err != nil {
			log.Printf("handleGetGuestHistory: error fetching history: %v", err)
			response.JSON(c, "Failed to fetch guest history", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Guest history retrieved successfully", http.StatusOK, history, nil)
	}
}

// handleGetGuestPreferences returns guest preferences
func (s *Server) handleGetGuestPreferences() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestService(s.GuestRepository)
		prefs, err := guestService.GetGuestPreferences(uint(id))
		if err != nil {
			log.Printf("handleGetGuestPreferences: error fetching preferences: %v", err)
			response.JSON(c, "Failed to fetch guest preferences", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Guest preferences retrieved successfully", http.StatusOK, prefs, nil)
	}
}

// handleGetGuestAIInsights returns guest AI insights (generates on-demand if stale)
func (s *Server) handleGetGuestAIInsights() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		// Verify guest exists
		guestService := services.NewGuestService(s.GuestRepository)
		_, err = guestService.GetGuestByID(uint(id))
		if err != nil {
			response.JSON(c, "Guest not found", http.StatusNotFound, nil, err)
			return
		}

		// Check if insights need to be refreshed
		insightsService := services.NewGuestInsightsService(s.DB)
		if insightsService.ShouldRefreshInsights(uint(id)) {
			// Generate fresh insights
			insights, err := insightsService.GenerateInsights(uint(id))
			if err != nil {
				log.Printf("handleGetGuestAIInsights: error generating insights: %v", err)
				// Fall back to existing insights if generation fails
				existingInsights, fetchErr := guestService.GetGuestAIInsights(uint(id))
				if fetchErr != nil {
					response.JSON(c, "Failed to fetch guest AI insights", http.StatusInternalServerError, nil, err)
					return
				}
				response.JSON(c, "Guest AI insights retrieved successfully", http.StatusOK, existingInsights, nil)
				return
			}
			response.JSON(c, "Guest AI insights generated successfully", http.StatusOK, insights, nil)
			return
		}

		// Return existing insights
		insights, err := guestService.GetGuestAIInsights(uint(id))
		if err != nil {
			log.Printf("handleGetGuestAIInsights: error fetching insights: %v", err)
			response.JSON(c, "Failed to fetch guest AI insights", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Guest AI insights retrieved successfully", http.StatusOK, insights, nil)
	}
}

// handleRefreshGuestAIInsights forces regeneration of guest AI insights
func (s *Server) handleRefreshGuestAIInsights() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		// Verify guest exists
		guestService := services.NewGuestService(s.GuestRepository)
		_, err = guestService.GetGuestByID(uint(id))
		if err != nil {
			response.JSON(c, "Guest not found", http.StatusNotFound, nil, err)
			return
		}

		// Generate fresh insights
		insightsService := services.NewGuestInsightsService(s.DB)
		insights, err := insightsService.GenerateInsights(uint(id))
		if err != nil {
			log.Printf("handleRefreshGuestAIInsights: error generating insights: %v", err)
			response.JSON(c, "Failed to generate guest AI insights", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Guest AI insights refreshed successfully", http.StatusOK, insights, nil)
	}
}

// handleHealthCheck returns health status
func (s *Server) handleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.JSON(c, "Service is healthy", http.StatusOK, gin.H{"status": "ok"}, nil)
	}
}
