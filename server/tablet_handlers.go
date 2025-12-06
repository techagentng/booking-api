package server

import (
	"hotel/models"
	"hotel/server/response"
	"hotel/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== Service Requests (Admin) ====================

// handleGetAllServiceRequests returns all service requests with pagination
func (s *Server) handleGetAllServiceRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		status := c.DefaultQuery("status", "")
		requestType := c.DefaultQuery("type", "")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}

		requests, total, err := s.GuestServiceRepository.GetAllServiceRequests(page, pageSize, status, requestType)
		if err != nil {
			log.Printf("handleGetAllServiceRequests: error: %v", err)
			response.JSON(c, "Failed to fetch service requests", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Service requests retrieved", http.StatusOK, gin.H{
			"data": requests,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		}, nil)
	}
}

// ==================== Guest Room Info ====================

// handleGetGuestByRoomNumber returns guest info for a room
func (s *Server) handleGetGuestByRoomNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomNumber := c.Param("room_number")
		if roomNumber == "" {
			response.JSON(c, "Room number is required", http.StatusBadRequest, nil, nil)
			return
		}

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		info, err := guestService.GetGuestByRoomNumber(roomNumber)
		if err != nil {
			log.Printf("handleGetGuestByRoomNumber: error: %v", err)
			response.JSON(c, "Guest information not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Guest information retrieved", http.StatusOK, info, nil)
	}
}

// ==================== Menu Categories ====================

// handleGetMenuCategories returns menu categories with item counts
func (s *Server) handleGetMenuCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		categories, err := guestService.GetMenuCategories()
		if err != nil {
			log.Printf("handleGetMenuCategories: error: %v", err)
			response.JSON(c, "Failed to fetch categories", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Categories retrieved", http.StatusOK, categories, nil)
	}
}

// ==================== Guest Active Orders ====================

// handleGetGuestActiveOrders returns active orders for a guest
func (s *Server) handleGetGuestActiveOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestID, err := strconv.ParseUint(c.Param("guest_id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		orders, err := guestService.GetGuestActiveOrders(uint(guestID))
		if err != nil {
			log.Printf("handleGetGuestActiveOrders: error: %v", err)
			response.JSON(c, "Failed to fetch active orders", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Active orders retrieved", http.StatusOK, orders, nil)
	}
}

// ==================== Order Status ====================

// handleGetOrderStatus returns order status for tracking
func (s *Server) handleGetOrderStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		status, err := guestService.GetOrderStatus(uint(orderID))
		if err != nil {
			log.Printf("handleGetOrderStatus: error: %v", err)
			response.JSON(c, "Order not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Order status retrieved", http.StatusOK, status, nil)
	}
}

// ==================== Housekeeping Request ====================

// handleCreateHousekeepingRequest creates a housekeeping service request
func (s *Server) handleCreateHousekeepingRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.HousekeepingRequestInput
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		result, err := guestService.CreateHousekeepingRequest(&req)
		if err != nil {
			log.Printf("handleCreateHousekeepingRequest: error: %v", err)
			response.JSON(c, "Failed to create request", http.StatusBadRequest, nil, err)
			return
		}

		// Send service request notification
		priority := result.Priority
		if priority == "" {
			priority = "normal"
		}
		s.NotificationHub.NotifyServiceRequest(
			result.ID,
			"Housekeeping",
			"", // Room number not in response
			"", // Guest name not in response
			priority,
			result.ServiceType,
		)

		response.JSON(c, "Service request created", http.StatusCreated, result, nil)
	}
}

// ==================== Maintenance Request ====================

// handleCreateMaintenanceRequest creates a maintenance service request
func (s *Server) handleCreateMaintenanceRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.MaintenanceRequestInput
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		result, err := guestService.CreateMaintenanceRequest(&req)
		if err != nil {
			log.Printf("handleCreateMaintenanceRequest: error: %v", err)
			response.JSON(c, "Failed to create request", http.StatusBadRequest, nil, err)
			return
		}

		// Send service request notification
		priority := result.Priority
		if priority == "" {
			priority = "normal"
		}
		s.NotificationHub.NotifyServiceRequest(
			result.ID,
			"Maintenance",
			"", // Room number not in response
			"", // Guest name not in response
			priority,
			result.IssueType,
		)

		response.JSON(c, "Maintenance request created", http.StatusCreated, result, nil)
	}
}

// ==================== Guest Service Requests ====================

// handleGetGuestServiceRequests returns service requests for a guest
func (s *Server) handleGetGuestServiceRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestID, err := strconv.ParseUint(c.Param("guest_id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		status := c.DefaultQuery("status", "active")

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		requests, err := guestService.GetServiceRequestsByGuestID(uint(guestID), status)
		if err != nil {
			log.Printf("handleGetGuestServiceRequests: error: %v", err)
			response.JSON(c, "Failed to fetch service requests", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Service requests retrieved", http.StatusOK, requests, nil)
	}
}

// ==================== Hotel Info ====================

// handleGetHotelInfo returns hotel information
func (s *Server) handleGetHotelInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		info := guestService.GetHotelInfo()
		response.JSON(c, "Hotel information retrieved", http.StatusOK, info, nil)
	}
}

// ==================== Express Checkout ====================

// handleExpressCheckout initiates express checkout
func (s *Server) handleExpressCheckout() gin.HandlerFunc {
	return func(c *gin.Context) {
		reservationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.ExpressCheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		guestService := services.NewGuestServiceService(
			s.GuestServiceRepository,
			s.RoomServiceRepository,
			s.ReservationRepository,
			s.RoomRepository,
			s.GuestRepository,
		)

		result, err := guestService.ExpressCheckout(uint(reservationID), &req)
		if err != nil {
			log.Printf("handleExpressCheckout: error: %v", err)
			response.JSON(c, "Express checkout failed", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Express checkout initiated", http.StatusOK, result, nil)
	}
}
