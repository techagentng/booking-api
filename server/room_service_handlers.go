package server

import (
	"hotel/db"
	"hotel/models"
	"hotel/server/response"
	"hotel/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== Menu Item Handlers ====================

// handleGetMenuItems returns all menu items
func (s *Server) handleGetMenuItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Query("category")
		availableOnly := c.Query("available") == "true"

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		items, err := roomServiceService.GetAllMenuItems(category, availableOnly)
		if err != nil {
			log.Printf("handleGetMenuItems: error fetching menu items: %v", err)
			response.JSON(c, "Failed to fetch menu items", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Menu items retrieved successfully", http.StatusOK, items, nil)
	}
}

// handleGetMenuItemByID returns a menu item by ID
func (s *Server) handleGetMenuItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid menu item ID", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		item, err := roomServiceService.GetMenuItemByID(uint(id))
		if err != nil {
			log.Printf("handleGetMenuItemByID: error fetching menu item: %v", err)
			response.JSON(c, "Menu item not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Menu item retrieved successfully", http.StatusOK, item, nil)
	}
}

// handleCreateMenuItem creates a new menu item
func (s *Server) handleCreateMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateMenuItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		item, err := roomServiceService.CreateMenuItem(&req)
		if err != nil {
			log.Printf("handleCreateMenuItem: error creating menu item: %v", err)
			response.JSON(c, "Failed to create menu item", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Menu item created successfully", http.StatusCreated, item, nil)
	}
}

// handleUpdateMenuItem updates a menu item
func (s *Server) handleUpdateMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid menu item ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateMenuItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		item, err := roomServiceService.UpdateMenuItem(uint(id), &req)
		if err != nil {
			log.Printf("handleUpdateMenuItem: error updating menu item: %v", err)
			response.JSON(c, "Failed to update menu item", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Menu item updated successfully", http.StatusOK, item, nil)
	}
}

// handleDeleteMenuItem deletes a menu item
func (s *Server) handleDeleteMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid menu item ID", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		err = roomServiceService.DeleteMenuItem(uint(id))
		if err != nil {
			log.Printf("handleDeleteMenuItem: error deleting menu item: %v", err)
			response.JSON(c, "Failed to delete menu item", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Menu item deleted successfully", http.StatusOK, nil, nil)
	}
}

// ==================== Order Handlers ====================

// handleGetRoomServiceOrders returns all room service orders
func (s *Server) handleGetRoomServiceOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := db.RoomServiceQueryParams{
			Page:       1,
			PageSize:   20,
			Status:     c.Query("status"),
			RoomNumber: c.Query("room_number"),
			Search:     c.Query("search"),
		}

		if page := c.Query("page"); page != "" {
			if p, err := strconv.Atoi(page); err == nil {
				params.Page = p
			}
		}

		if pageSize := c.Query("page_size"); pageSize != "" {
			if ps, err := strconv.Atoi(pageSize); err == nil {
				params.PageSize = ps
			}
		}

		if date := c.Query("date"); date != "" {
			if parsed, err := time.Parse("2006-01-02", date); err == nil {
				params.Date = &parsed
			}
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		result, err := roomServiceService.GetAllOrders(params)
		if err != nil {
			log.Printf("handleGetRoomServiceOrders: error fetching orders: %v", err)
			response.JSON(c, "Failed to fetch orders", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Room service orders retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetRoomServiceOrderByID returns an order by ID
func (s *Server) handleGetRoomServiceOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		order, err := roomServiceService.GetOrderByID(uint(id))
		if err != nil {
			log.Printf("handleGetRoomServiceOrderByID: error fetching order: %v", err)
			response.JSON(c, "Order not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Order retrieved successfully", http.StatusOK, order, nil)
	}
}

// handleGetActiveOrders returns all active orders
func (s *Server) handleGetActiveOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		orders, err := roomServiceService.GetActiveOrders()
		if err != nil {
			log.Printf("handleGetActiveOrders: error fetching active orders: %v", err)
			response.JSON(c, "Failed to fetch active orders", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Active orders retrieved", http.StatusOK, orders, nil)
	}
}

// handleCreateRoomServiceOrder creates a new room service order
func (s *Server) handleCreateRoomServiceOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateRoomServiceOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		order, err := roomServiceService.CreateOrder(&req)
		if err != nil {
			log.Printf("handleCreateRoomServiceOrder: error creating order: %v", err)
			response.JSON(c, "Order creation failed", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Order created successfully", http.StatusCreated, gin.H{
			"id":                      order.ID,
			"order_number":            order.OrderNumber,
			"status":                  order.Status,
			"total_amount":            order.TotalAmount,
			"estimated_delivery_time": order.EstimatedDeliveryTime,
		}, nil)
	}
}

// handleUpdateRoomServiceOrderStatus updates order status
func (s *Server) handleUpdateRoomServiceOrderStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateOrderStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		err = roomServiceService.UpdateOrderStatus(uint(id), &req)
		if err != nil {
			log.Printf("handleUpdateRoomServiceOrderStatus: error updating status: %v", err)
			response.JSON(c, "Invalid status transition", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Order status updated successfully", http.StatusOK, gin.H{
			"id":         id,
			"status":     req.Status,
			"updated_at": time.Now(),
		}, nil)
	}
}

// handleDeliverOrder marks an order as delivered
func (s *Server) handleDeliverOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.DeliverOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		err = roomServiceService.DeliverOrder(uint(id), &req)
		if err != nil {
			log.Printf("handleDeliverOrder: error delivering order: %v", err)
			response.JSON(c, "Delivery failed", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Order marked as delivered", http.StatusOK, gin.H{
			"id":                   id,
			"status":               "delivered",
			"actual_delivery_time": time.Now(),
			"delivered_by":         req.DeliveredBy,
		}, nil)
	}
}

// handleCancelRoomServiceOrder cancels an order
func (s *Server) handleCancelRoomServiceOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid order ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.CancelOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		err = roomServiceService.CancelOrder(uint(id), &req)
		if err != nil {
			log.Printf("handleCancelRoomServiceOrder: error cancelling order: %v", err)
			response.JSON(c, "Cancellation failed", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Order cancelled successfully", http.StatusOK, gin.H{
			"id":                  id,
			"status":              "cancelled",
			"cancellation_reason": req.CancellationReason,
		}, nil)
	}
}

// handleGetRoomServiceStats returns room service statistics
func (s *Server) handleGetRoomServiceStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Query("date")
		period := c.DefaultQuery("period", "today")

		roomServiceService := services.NewRoomServiceService(s.RoomServiceRepository, s.RoomRepository, s.GuestRepository)
		stats, err := roomServiceService.GetStats(date, period)
		if err != nil {
			log.Printf("handleGetRoomServiceStats: error fetching stats: %v", err)
			response.JSON(c, "Failed to fetch statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Statistics retrieved successfully", http.StatusOK, stats, nil)
	}
}
