package server

import (
	"fmt"
	"hotel/models"
	"hotel/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// handleSSENotifications handles the SSE endpoint for real-time notifications
func (s *Server) handleSSENotifications() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set headers for SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Accel-Buffering", "no")

		// Get user info from context (set by Authorize middleware)
		userID := uint(0)
		userRole := "guest"
		if user, exists := c.Get("user"); exists {
			if u, ok := user.(map[string]interface{}); ok {
				if id, ok := u["id"].(float64); ok {
					userID = uint(id)
				}
				if role, ok := u["role"].(string); ok {
					userRole = role
				}
			}
		}

		// Create a new client
		client := &services.Client{
			ID:       uuid.New().String(),
			Channel:  make(chan []byte, 10),
			UserID:   userID,
			UserRole: userRole,
		}

		// Register the client
		s.NotificationHub.Register(client)

		// Send initial connection message
		c.SSEvent("connected", map[string]interface{}{
			"client_id": client.ID,
			"message":   "Connected to notification stream",
			"timestamp": time.Now(),
		})
		c.Writer.Flush()

		// Create a channel to detect client disconnect
		clientGone := c.Request.Context().Done()

		// Keep the connection alive
		for {
			select {
			case <-clientGone:
				// Client disconnected
				s.NotificationHub.Unregister(client)
				return

			case data := <-client.Channel:
				// Send notification to client
				c.SSEvent("notification", string(data))
				c.Writer.Flush()

			case <-time.After(30 * time.Second):
				// Send heartbeat to keep connection alive
				c.SSEvent("heartbeat", map[string]interface{}{
					"timestamp": time.Now(),
				})
				c.Writer.Flush()
			}
		}
	}
}

// handleSSENotificationsPublic handles SSE without auth (for testing)
func (s *Server) handleSSENotificationsPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set headers for SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Accel-Buffering", "no")

		// Create a new client
		client := &services.Client{
			ID:       uuid.New().String(),
			Channel:  make(chan []byte, 10),
			UserID:   0,
			UserRole: "public",
		}

		// Register the client
		s.NotificationHub.Register(client)

		// Send initial connection message
		c.SSEvent("connected", map[string]interface{}{
			"client_id": client.ID,
			"message":   "Connected to notification stream",
			"timestamp": time.Now(),
		})
		c.Writer.Flush()

		// Create a channel to detect client disconnect
		clientGone := c.Request.Context().Done()

		// Keep the connection alive
		for {
			select {
			case <-clientGone:
				// Client disconnected
				s.NotificationHub.Unregister(client)
				return

			case data := <-client.Channel:
				// Send notification to client
				c.SSEvent("notification", string(data))
				c.Writer.Flush()

			case <-time.After(30 * time.Second):
				// Send heartbeat to keep connection alive
				c.SSEvent("heartbeat", map[string]interface{}{
					"timestamp": time.Now(),
				})
				c.Writer.Flush()
			}
		}
	}
}

// handleGetNotificationStats returns stats about connected clients
func (s *Server) handleGetNotificationStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"connected_clients": s.NotificationHub.GetConnectedClients(),
			"status":            "active",
		})
	}
}

// handleTestNotification sends a test notification (for development)
func (s *Server) handleTestNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Type    string `json:"type"`
			Title   string `json:"title"`
			Message string `json:"message"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Type == "" {
			req.Type = "test"
		}
		if req.Title == "" {
			req.Title = "Test Notification"
		}
		if req.Message == "" {
			req.Message = fmt.Sprintf("Test notification sent at %s", time.Now().Format(time.RFC3339))
		}

		// Create and broadcast test notification using models
		notification := models.NewNotification(
			models.NotificationType(req.Type),
			req.Title,
			req.Message,
			models.PriorityNormal,
			map[string]string{"source": "test_endpoint"},
		)

		s.NotificationHub.Broadcast(notification)

		c.JSON(http.StatusOK, gin.H{
			"message":      "Test notification sent",
			"notification": notification,
			"clients":      s.NotificationHub.GetConnectedClients(),
		})
	}
}
