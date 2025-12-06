package server

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}

	r := gin.New()

	// Logger middleware with custom format
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "X-Client-State"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Increase memory limit for multipart forms
	r.MaxMultipartMemory = 32 << 20
	s.defineRoutes(r)

	return r
}

func (s *Server) defineRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", s.handleHealthCheck())

	// Frontend routes (without /api/v1 prefix)
	frontend := router.Group("")
	{
		// Dashboard routes
		frontend.GET("/reservations/dashboard", s.Authorize(), s.handleGetDashboardStats())
		frontend.GET("/reservations/recent-activity", s.Authorize(), s.handleGetRecentActivity())

		// Reservations with query parameters
		frontend.GET("/reservations", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/reservations"
			router.HandleContext(c)
		})

		frontend.POST("/reservations", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/reservations"
			router.HandleContext(c)
		})

		// Public reservation (no auth required)
		frontend.POST("/reservations/public", func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/reservations/public"
			router.HandleContext(c)
		})

		// Guest routes
		frontend.GET("/guests", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/guests"
			router.HandleContext(c)
		})

		frontend.POST("/guests", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/guests"
			router.HandleContext(c)
		})

		frontend.GET("/guests/:id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/guests/" + c.Param("id")
			router.HandleContext(c)
		})

		frontend.GET("/guests/:id/preferences", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/guests/" + c.Param("id") + "/preferences"
			router.HandleContext(c)
		})

		frontend.GET("/guests/:id/ai-insights", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/guests/" + c.Param("id") + "/ai-insights"
			router.HandleContext(c)
		})

		frontend.GET("/guests/:id/history", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/guests/" + c.Param("id") + "/history"
			router.HandleContext(c)
		})

		// Staff routes
		frontend.GET("/staff", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff"
			router.HandleContext(c)
		})

		frontend.GET("/staff/available", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff/available"
			router.HandleContext(c)
		})

		frontend.GET("/staff/stats", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff/stats"
			router.HandleContext(c)
		})

		frontend.GET("/staff/on-duty", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff/on-duty"
			router.HandleContext(c)
		})

		frontend.POST("/staff", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff"
			router.HandleContext(c)
		})

		frontend.PUT("/staff/:id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff/" + c.Param("id")
			router.HandleContext(c)
		})

		frontend.DELETE("/staff/:id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/staff/" + c.Param("id")
			router.HandleContext(c)
		})

		// Service requests with query parameters
		frontend.GET("/service-requests", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/service-requests"
			router.HandleContext(c)
		})

		frontend.POST("/service-requests/:id/auto-assign", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/service-requests/" + c.Param("id") + "/auto-assign"
			router.HandleContext(c)
		})

		frontend.POST("/service-requests/:id/assign", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/service-requests/" + c.Param("id") + "/assign"
			router.HandleContext(c)
		})

		frontend.POST("/service-requests/:id/complete", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/service-requests/" + c.Param("id") + "/complete"
			router.HandleContext(c)
		})

		// Room routes
		frontend.GET("/rooms/summary", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms/summary"
			router.HandleContext(c)
		})

		frontend.GET("/rooms", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms"
			router.HandleContext(c)
		})

		frontend.POST("/rooms", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms"
			router.HandleContext(c)
		})

		frontend.GET("/rooms/guest/:room_number", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms/guest/" + c.Param("room_number")
			router.HandleContext(c)
		})

		frontend.GET("/rooms/:id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms/" + c.Param("id")
			router.HandleContext(c)
		})

		frontend.PUT("/rooms/:id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms/" + c.Param("id")
			router.HandleContext(c)
		})

		frontend.DELETE("/rooms/:id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/rooms/" + c.Param("id")
			router.HandleContext(c)
		})

		// Room service routes
		frontend.GET("/room-service/menu", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/room-service/menu"
			router.HandleContext(c)
		})

		frontend.GET("/room-service/menu/categories", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/room-service/menu/categories"
			router.HandleContext(c)
		})

		frontend.POST("/room-service/orders", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/room-service/orders"
			router.HandleContext(c)
		})

		frontend.GET("/room-service/orders/guest/:guest_id/active", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/room-service/orders/guest/" + c.Param("guest_id") + "/active"
			router.HandleContext(c)
		})

		frontend.GET("/services/guest/:guest_id", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/services/guest/" + c.Param("guest_id")
			router.HandleContext(c)
		})

		frontend.POST("/services/housekeeping", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/services/housekeeping"
			router.HandleContext(c)
		})

		frontend.POST("/services/maintenance", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/services/maintenance"
			router.HandleContext(c)
		})

		// Reservation status update
		frontend.PUT("/reservations/:id/status", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/reservations/" + c.Param("id") + "/status"
			router.HandleContext(c)
		})

		// Reservation check-in/check-out
		frontend.POST("/reservations/:id/checkin", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/reservations/" + c.Param("id") + "/checkin"
			router.HandleContext(c)
		})

		frontend.POST("/reservations/:id/checkout", s.Authorize(), func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/reservations/" + c.Param("id") + "/checkout"
			router.HandleContext(c)
		})

		// Service request status update (public - no auth, for staff tablets)
		frontend.GET("/service-requests/assigned", func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/service-requests/assigned"
			router.HandleContext(c)
		})

		frontend.PUT("/service-requests/assigned/:id/status", func(c *gin.Context) {
			c.Request.URL.Path = "/api/v1/service-requests/assigned/" + c.Param("id") + "/status"
			router.HandleContext(c)
		})

		// Notification routes
		frontend.GET("/notifications/stream", s.handleSSENotificationsPublic())
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		v1.POST("/auth/signup", s.handleSignup())
		v1.POST("/auth/login", s.handleLogin())
		v1.POST("/google/user/login", s.handleGoogleLogin())

		// Protected routes (authentication required)
		authorized := v1.Group("/auth")
		authorized.Use(s.Authorize())
		{
			authorized.POST("/logout", s.handleLogout())
		}

		// Notification routes (SSE)
		v1.GET("/notifications/stream", s.handleSSENotificationsPublic()) // Public for easy testing
		notifications := v1.Group("/notifications")
		notifications.Use(s.Authorize())
		{
			notifications.GET("/stats", s.handleGetNotificationStats())
			notifications.POST("/test", s.handleTestNotification())
		}

		// Guest routes
		guests := v1.Group("/guests")
		guests.Use(s.Authorize())
		{
			guests.GET("", s.handleGetGuests())
			guests.GET("/:id", s.handleGetGuestByID())
			guests.POST("", s.handleCreateGuest())
			guests.PUT("/:id", s.handleUpdateGuest())
			guests.DELETE("/:id", s.handleDeleteGuest())
			guests.GET("/:id/history", s.handleGetGuestHistory())
			guests.GET("/:id/preferences", s.handleGetGuestPreferences())
			guests.GET("/:id/ai-insights", s.handleGetGuestAIInsights())
			guests.POST("/:id/ai-insights/refresh", s.handleRefreshGuestAIInsights())
		}

		// Reservation routes
		reservations := v1.Group("/reservations")
		reservations.Use(s.Authorize())
		{
			reservations.GET("", s.handleGetReservations())
			reservations.GET("/stats", s.handleGetReservationStats())
			reservations.GET("/checkin/stats", s.handleGetCheckInStats())
			reservations.GET("/dashboard", s.handleGetDashboardStats())
			reservations.GET("/recent-activity", s.handleGetRecentActivity())
			reservations.GET("/guest/:guestID", s.handleGetReservationsByGuest())
			reservations.GET("/room/:roomID", s.handleGetReservationsByRoom())
			reservations.GET("/:id", s.handleGetReservationByID())
			reservations.GET("/:id/payment", s.handleGetPaymentDetails())
			reservations.POST("", s.handleCreateReservation())
			reservations.PUT("/:id", s.handleUpdateReservation())
			reservations.PUT("/:id/status", s.handleUpdateReservationStatus())
			reservations.PUT("/:id/cancel", s.handleCancelReservation())
			reservations.POST("/:id/checkin", s.handleCheckIn())
			reservations.POST("/:id/checkout", s.handleCheckOut())
			reservations.DELETE("/:id", s.handleDeleteReservation())
		}

		// Public reservation routes (no auth required)
		v1.POST("/reservations/public", s.handleCreateReservation())

		// Room routes
		rooms := v1.Group("/rooms")
		rooms.Use(s.Authorize())
		{
			rooms.GET("", s.handleGetRooms())
			rooms.GET("/available", s.handleGetAvailableRooms())
			rooms.GET("/stats", s.handleGetRoomTypeStats())
			rooms.GET("/summary", s.handleGetRoomStats())
			rooms.GET("/type/:type", s.handleGetRoomsByType())
			rooms.GET("/guest/:room_number", s.handleGetGuestByRoomNumber()) // Get guest info by room number
			rooms.GET("/:id", s.handleGetRoomByID())
			rooms.POST("", s.handleCreateRoom())
			rooms.PUT("/:id", s.handleUpdateRoom())
			rooms.DELETE("/:id", s.handleDeleteRoom())
			rooms.POST("/availability", s.handleCheckRoomAvailability())
			rooms.PUT("/:id/status", s.handleUpdateRoomStatus())
		}

		// Room Service routes
		roomService := v1.Group("/room-service")
		roomService.Use(s.Authorize())
		{
			// Menu routes
			roomService.GET("/menu", s.handleGetMenuItems())
			roomService.GET("/menu/categories", s.handleGetMenuCategories())
			roomService.GET("/menu/:id", s.handleGetMenuItemByID())
			roomService.POST("/menu", s.handleCreateMenuItem())
			roomService.PUT("/menu/:id", s.handleUpdateMenuItem())
			roomService.DELETE("/menu/:id", s.handleDeleteMenuItem())

			// Order routes
			roomService.GET("/orders", s.handleGetRoomServiceOrders())
			roomService.GET("/orders/active", s.handleGetActiveOrders())
			roomService.GET("/orders/guest/:guest_id/active", s.handleGetGuestActiveOrders())
			roomService.GET("/orders/:id", s.handleGetRoomServiceOrderByID())
			roomService.GET("/orders/:id/status", s.handleGetOrderStatus())
			roomService.POST("/orders", s.handleCreateRoomServiceOrder())
			roomService.PUT("/orders/:id/status", s.handleUpdateRoomServiceOrderStatus())
			roomService.POST("/orders/:id/deliver", s.handleDeliverOrder())
			roomService.PUT("/orders/:id/cancel", s.handleCancelRoomServiceOrder())

			// Stats
			roomService.GET("/stats", s.handleGetRoomServiceStats())
		}

		// Guest Services routes (for tablet/guest-facing)
		services := v1.Group("/services")
		services.Use(s.Authorize())
		{
			services.POST("/housekeeping", s.handleCreateHousekeepingRequest())
			services.POST("/maintenance", s.handleCreateMaintenanceRequest())
			services.GET("/guest/:guest_id", s.handleGetGuestServiceRequests())
		}

		// Service Requests (admin/staff)
		serviceRequests := v1.Group("/service-requests")
		serviceRequests.Use(s.Authorize())
		{
			serviceRequests.GET("", s.handleGetAllServiceRequests())
			serviceRequests.POST("/:id/auto-assign", s.handleAutoAssignRequest())
			serviceRequests.POST("/:id/assign", s.handleManualAssignRequest())
			serviceRequests.POST("/:id/complete", s.handleCompleteServiceRequest())
		}

		// Public assigned tasks endpoints (no auth)
		v1.GET("/service-requests/assigned", s.handleGetAssignedTasks())
		v1.PUT("/service-requests/assigned/:id/status", s.handleUpdateTaskStatus())

		// Staff routes
		staff := v1.Group("/staff")
		staff.Use(s.Authorize())
		{
			staff.GET("", s.handleGetStaff())
			staff.GET("/stats", s.handleGetStaffStats())
			staff.GET("/on-duty", s.handleGetOnDutyStaff())
			staff.GET("/available", s.handleGetAvailableStaff())
			staff.GET("/department/:department", s.handleGetStaffByDepartment())
			staff.GET("/:id", s.handleGetStaffByID())
			staff.POST("", s.handleCreateStaff())
			staff.PUT("/:id", s.handleUpdateStaff())
			staff.DELETE("/:id", s.handleDeleteStaff())
			staff.POST("/:id/clock-in", s.handleClockIn())
			staff.POST("/:id/clock-out", s.handleClockOut())
			staff.PUT("/:id/availability", s.handleSetAvailability())
		}

		// Hotel Info (for tablet)
		v1.GET("/hotel/info", s.handleGetHotelInfo())

		// Tablet/Guest-facing routes
		tablet := v1.Group("/tablet")
		tablet.Use(s.Authorize())
		{
			tablet.GET("/room/:room_number/guest", s.handleGetGuestByRoomNumber())
		}

		// Express checkout
		reservations.POST("/:id/checkout/express", s.handleExpressCheckout())

		// 	// Housekeeping Request routes
		// 	housekeepingRequests := v1.Group("/housekeeping-requests")
		// 	housekeepingRequests.Use(s.Authorize())
		// 	{
		// 		housekeepingRequests.GET("", s.handleGetHousekeepingRequests())
		// 		housekeepingRequests.GET("/:id", s.handleGetHousekeepingRequestByID())
		// 		housekeepingRequests.POST("", s.handleCreateHousekeepingRequest())
		// 		housekeepingRequests.PUT("/:id", s.handleUpdateHousekeepingRequest())
		// 		housekeepingRequests.DELETE("/:id", s.handleDeleteHousekeepingRequest())
		// 		housekeepingRequests.GET("/reservation/:reservationID", s.handleGetHousekeepingRequestsByReservation())
		// 		housekeepingRequests.GET("/status/:status", s.handleGetHousekeepingRequestsByStatus())
		// 	}

		// 	// Maintenance Issue routes
		// 	maintenanceIssues := v1.Group("/maintenance-issues")
		// 	maintenanceIssues.Use(s.Authorize())
		// 	{
		// 		maintenanceIssues.GET("", s.handleGetMaintenanceIssues())
		// 		maintenanceIssues.GET("/:id", s.handleGetMaintenanceIssueByID())
		// 		maintenanceIssues.POST("", s.handleCreateMaintenanceIssue())
		// 		maintenanceIssues.PUT("/:id", s.handleUpdateMaintenanceIssue())
		// 		maintenanceIssues.DELETE("/:id", s.handleDeleteMaintenanceIssue())
		// 		maintenanceIssues.GET("/reservation/:reservationID", s.handleGetMaintenanceIssuesByReservation())
		// 		maintenanceIssues.GET("/status/:status", s.handleGetMaintenanceIssuesByStatus())
		// 	}

		// 	// Menu Item routes
		// 	menuItems := v1.Group("/menu-items")
		// 	menuItems.Use(s.Authorize())
		// 	{
		// 		menuItems.GET("", s.handleGetMenuItems())
		// 		menuItems.GET("/:id", s.handleGetMenuItemByID())
		// 		menuItems.POST("", s.handleCreateMenuItem())
		// 		menuItems.PUT("/:id", s.handleUpdateMenuItem())
		// 		menuItems.DELETE("/:id", s.handleDeleteMenuItem())
		// 		menuItems.GET("/category/:category", s.handleGetMenuItemsByCategory())
		// 	}

		// 	// Check-In routes
		// 	checkIns := v1.Group("/check-ins")
		// 	checkIns.Use(s.Authorize())
		// 	{
		// 		checkIns.GET("", s.handleGetCheckIns())
		// 		checkIns.GET("/:id", s.handleGetCheckInByID())
		// 		checkIns.POST("", s.handleCreateCheckIn())
		// 		checkIns.PUT("/:id", s.handleUpdateCheckIn())
		// 		checkIns.GET("/reservation/:reservationID", s.handleGetCheckInByReservation())
		// 	}

		// 	// Check-Out routes
		// 	checkOuts := v1.Group("/check-outs")
		// 	checkOuts.Use(s.Authorize())
		// 	{
		// 		checkOuts.GET("", s.handleGetCheckOuts())
		// 		checkOuts.GET("/:id", s.handleGetCheckOutByID())
		// 		checkOuts.POST("", s.handleCreateCheckOut())
		// 		checkOuts.PUT("/:id", s.handleUpdateCheckOut())
		// 		checkOuts.GET("/reservation/:reservationID", s.handleGetCheckOutByReservation())
		// 	}

		// 	// Dashboard routes
		// 	dashboard := v1.Group("/dashboard")
		// 	dashboard.Use(s.Authorize())
		// 	{
		// 		dashboard.GET("/stats", s.handleGetDashboardStats())
		// 		dashboard.GET("/room-status", s.handleGetRoomStatusSummary())
		// 		dashboard.GET("/service-requests-summary", s.handleGetServiceRequestsSummary())
		// 		dashboard.GET("/revenue", s.handleGetRevenueStats())
		// 	}

		// 	// In-Room Tablet routes
		// 	inRoomTablet := v1.Group("/in-room-tablet")
		// 	inRoomTablet.Use(s.Authorize())
		// 	{
		// 		inRoomTablet.GET("/reservation/:reservationID", s.handleGetInRoomTabletReservation())
		// 		inRoomTablet.GET("/menu", s.handleGetInRoomTabletMenu())
		// 		inRoomTablet.POST("/room-service-order", s.handleCreateInRoomTabletRoomServiceOrder())
		// 		inRoomTablet.POST("/housekeeping-request", s.handleCreateInRoomTabletHousekeepingRequest())
		// 		inRoomTablet.POST("/maintenance-issue", s.handleCreateInRoomTabletMaintenanceIssue())
		// 	}

		// 	// Staff routes
		// 	staff := v1.Group("/staff")
		// 	staff.Use(s.Authorize())
		// 	{
		// 		staff.GET("", s.handleGetStaff())
		// 		staff.GET("/:id", s.handleGetStaffByID())
		// 		staff.POST("", s.handleCreateStaff())
		// 		staff.PUT("/:id", s.handleUpdateStaff())
		// 		staff.DELETE("/:id", s.handleDeleteStaff())
		// 		staff.GET("/role/:role", s.handleGetStaffByRole())
		// 	}
		// }
	}
}
