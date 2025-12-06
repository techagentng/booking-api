package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hotel/db"
	"hotel/models"
	"hotel/server/response"
	"hotel/services"

	"github.com/gin-gonic/gin"
)

// handleGetReservations returns all reservations with pagination and filters
func (s *Server) handleGetReservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := db.ReservationQueryParams{
			Page:     1,
			PageSize: 10,
		}

		// Parse query parameters
		if p := c.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				params.Page = parsed
			}
		}

		if ps := c.Query("page_size"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
				params.PageSize = parsed
			}
		}

		params.Status = c.Query("status")
		params.Search = c.Query("search")

		if guestID := c.Query("guest_id"); guestID != "" {
			if parsed, err := strconv.ParseUint(guestID, 10, 32); err == nil {
				params.GuestID = uint(parsed)
			}
		}

		if roomID := c.Query("room_id"); roomID != "" {
			if parsed, err := strconv.ParseUint(roomID, 10, 32); err == nil {
				params.RoomID = uint(parsed)
			}
		}

		if checkInFrom := c.Query("check_in_from"); checkInFrom != "" {
			if parsed, err := time.Parse("2006-01-02", checkInFrom); err == nil {
				params.CheckInFrom = &parsed
			}
		}

		if checkInTo := c.Query("check_in_to"); checkInTo != "" {
			if parsed, err := time.Parse("2006-01-02", checkInTo); err == nil {
				params.CheckInTo = &parsed
			}
		}

		// Exact check-in date filter (for today's arrivals)
		if checkInDate := c.Query("check_in_date"); checkInDate != "" {
			if parsed, err := time.Parse("2006-01-02", checkInDate); err == nil {
				params.CheckInDate = &parsed
			}
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		result, err := reservationService.GetAllReservations(params)
		if err != nil {
			log.Printf("handleGetReservations: error fetching reservations: %v", err)
			response.JSON(c, "Failed to fetch reservations", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Reservations retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetReservationByID returns a reservation by ID
func (s *Server) handleGetReservationByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservation, err := reservationService.GetReservationByID(uint(id))
		if err != nil {
			log.Printf("handleGetReservationByID: error fetching reservation: %v", err)
			response.JSON(c, "Reservation not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Reservation retrieved successfully", http.StatusOK, reservation, nil)
	}
}

// handleCreateReservation creates a new reservation (supports JSON and multipart/form-data)
func (s *Server) handleCreateReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateReservationRequest
		var idDocumentURL string

		contentType := c.GetHeader("Content-Type")
		log.Printf("handleCreateReservation: Content-Type = %s", contentType)

		// Handle multipart/form-data (with file upload)
		isMultipart := strings.Contains(contentType, "multipart/form-data")
		log.Printf("handleCreateReservation: isMultipart = %v", isMultipart)

		if isMultipart {
			// Parse multipart form (max 10MB)
			if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
				log.Printf("handleCreateReservation: ParseMultipartForm error: %v", err)
				response.JSON(c, "Failed to parse form data", http.StatusBadRequest, nil, err)
				return
			}

			// Parse form data
			guestID, _ := strconv.ParseUint(c.PostForm("guest_id"), 10, 32)
			roomID, _ := strconv.ParseUint(c.PostForm("room_id"), 10, 32)
			numberOfGuests, _ := strconv.Atoi(c.PostForm("number_of_guests"))

			log.Printf("handleCreateReservation: Form data - guest_id=%d, room_id=%d, check_in=%s, check_out=%s, guests=%d",
				guestID, roomID, c.PostForm("check_in_date"), c.PostForm("check_out_date"), numberOfGuests)

			req = models.CreateReservationRequest{
				// Existing guest OR new guest fields
				GuestID:          models.FlexUint(guestID),
				GuestName:        c.PostForm("guest_name"),
				GuestEmail:       c.PostForm("guest_email"),
				GuestPhone:       c.PostForm("guest_phone"),
				GuestNationality: c.PostForm("guest_nationality"),
				GuestIDType:      c.PostForm("guest_id_type"),
				GuestIDNumber:    c.PostForm("guest_id_number"),
				// Reservation fields
				RoomID:          models.FlexUint(roomID),
				CheckInDate:     c.PostForm("check_in_date"),
				CheckOutDate:    c.PostForm("check_out_date"),
				NumberOfGuests:  models.FlexInt(numberOfGuests),
				SpecialRequests: c.PostForm("special_requests"),
				PaymentMethod:   c.PostForm("payment_method"),
			}

			// Handle file upload
			file, header, err := c.Request.FormFile("id_document")
			if err == nil {
				defer file.Close()

				// Validate file
				allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/webp", "application/pdf"}
				if err := services.ValidateFile(header, 5, allowedTypes); err != nil {
					response.JSON(c, "Invalid file", http.StatusBadRequest, nil, err)
					return
				}

				// Upload to S3
				s3Service, err := services.NewS3Service()
				if err != nil {
					log.Printf("handleCreateReservation: S3 service error: %v", err)
					// Continue without file upload if S3 is not configured
				} else {
					url, err := s3Service.UploadFile(file, header, "id-documents")
					if err != nil {
						log.Printf("handleCreateReservation: S3 upload error: %v", err)
						response.JSON(c, "Failed to upload ID document", http.StatusInternalServerError, nil, err)
						return
					}
					idDocumentURL = url
				}
			}
		} else {
			// Handle JSON request (backward compatible)
			if err := c.ShouldBindJSON(&req); err != nil {
				response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
				return
			}
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservation, err := reservationService.CreateReservation(&req)
		if err != nil {
			log.Printf("handleCreateReservation: error creating reservation: %v", err)
			statusCode := http.StatusBadRequest
			if err.Error() == "room is not available for selected dates" {
				statusCode = http.StatusConflict
			}
			response.JSON(c, "Reservation creation failed", statusCode, nil, err)
			return
		}

		// Update reservation with ID document URL if uploaded
		if idDocumentURL != "" {
			reservation.IDDocumentURL = idDocumentURL
			s.DB.Save(reservation)
		}

		// Send new booking notification
		guestName := ""
		roomNumber := ""
		roomType := ""
		if reservation.Guest != nil {
			guestName = reservation.Guest.Name
		}
		if reservation.Room != nil {
			roomNumber = reservation.Room.RoomNumber
			roomType = reservation.Room.RoomType
		}
		s.NotificationHub.NotifyNewBooking(
			reservation.ID,
			guestName,
			roomNumber,
			roomType,
			reservation.CheckInDate.Format("2006-01-02"),
			reservation.CheckOutDate.Format("2006-01-02"),
			reservation.TotalAmount,
		)

		response.JSON(c, "Reservation created successfully", http.StatusCreated, reservation, nil)
	}
}

// handleUpdateReservation updates a reservation
func (s *Server) handleUpdateReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateReservationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservation, err := reservationService.UpdateReservation(uint(id), &req)
		if err != nil {
			log.Printf("handleUpdateReservation: error updating reservation: %v", err)
			statusCode := http.StatusBadRequest
			if err.Error() == "room is not available for selected dates" {
				statusCode = http.StatusConflict
			}
			response.JSON(c, "Reservation update failed", statusCode, nil, err)
			return
		}

		response.JSON(c, "Reservation updated successfully", http.StatusOK, reservation, nil)
	}
}

// handleUpdateReservationStatus updates a reservation's status
func (s *Server) handleUpdateReservationStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateReservationStatusRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		err = reservationService.UpdateReservationStatus(uint(id), req.Status)
		if err != nil {
			log.Printf("handleUpdateReservationStatus: error updating status: %v", err)
			response.JSON(c, "Failed to update reservation status", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Reservation status updated successfully", http.StatusOK, nil, nil)
	}
}

// handleCancelReservation cancels a reservation
func (s *Server) handleCancelReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.CancelReservationRequest
		// Request body is optional
		c.ShouldBindJSON(&req)

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		err = reservationService.CancelReservation(uint(id), &req)
		if err != nil {
			log.Printf("handleCancelReservation: error cancelling reservation: %v", err)
			response.JSON(c, "Failed to cancel reservation", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Reservation cancelled successfully", http.StatusOK, nil, nil)
	}
}

// handleCheckIn performs check-in for a reservation
func (s *Server) handleCheckIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.CheckInRequest
		// Request body is optional
		c.ShouldBindJSON(&req)

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservation, err := reservationService.CheckIn(uint(id), &req)
		if err != nil {
			log.Printf("handleCheckIn: error checking in: %v", err)
			response.JSON(c, "Check-in failed", http.StatusBadRequest, nil, err)
			return
		}

		// Send check-in notification
		guestName := ""
		roomNumber := ""
		if reservation.Guest != nil {
			guestName = reservation.Guest.Name
		}
		if reservation.Room != nil {
			roomNumber = reservation.Room.RoomNumber
		}
		s.NotificationHub.NotifyCheckIn(reservation.ID, guestName, roomNumber)

		response.JSON(c, "Guest checked in successfully", http.StatusOK, gin.H{
			"id":            reservation.ID,
			"status":        reservation.Status,
			"checked_in_at": reservation.CheckedInAt,
			"room_status":   "occupied",
		}, nil)
	}
}

// handleCheckOut performs check-out for a reservation
func (s *Server) handleCheckOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.CheckOutRequest
		// Request body is optional
		c.ShouldBindJSON(&req)

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservation, err := reservationService.CheckOut(uint(id), &req)
		if err != nil {
			log.Printf("handleCheckOut: error checking out: %v", err)
			response.JSON(c, "Check-out failed", http.StatusBadRequest, nil, err)
			return
		}

		// Trigger guest insights generation asynchronously
		go func(guestID uint) {
			insightsService := services.NewGuestInsightsService(s.DB)
			if _, err := insightsService.GenerateInsights(guestID); err != nil {
				log.Printf("handleCheckOut: error generating insights for guest %d: %v", guestID, err)
			}
		}(reservation.GuestID)

		// Send check-out notification
		guestName := ""
		roomNumber := ""
		if reservation.Guest != nil {
			guestName = reservation.Guest.Name
		}
		if reservation.Room != nil {
			roomNumber = reservation.Room.RoomNumber
		}
		s.NotificationHub.NotifyCheckOut(reservation.ID, guestName, roomNumber, reservation.TotalAmount)

		response.JSON(c, "Guest checked out successfully", http.StatusOK, gin.H{
			"id":             reservation.ID,
			"status":         reservation.Status,
			"checked_out_at": reservation.CheckedOutAt,
			"final_amount":   reservation.TotalAmount,
			"room_status":    "cleaning",
		}, nil)
	}
}

// handleGetReservationStats returns reservation statistics
func (s *Server) handleGetReservationStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		fromDate := c.Query("from_date")
		toDate := c.Query("to_date")

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		stats, err := reservationService.GetReservationStats(fromDate, toDate)
		if err != nil {
			log.Printf("handleGetReservationStats: error fetching stats: %v", err)
			response.JSON(c, "Failed to fetch reservation statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Reservation statistics retrieved successfully", http.StatusOK, stats, nil)
	}
}

// handleGetReservationsByGuest returns all reservations for a guest
func (s *Server) handleGetReservationsByGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestID, err := strconv.ParseUint(c.Param("guestID"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid guest ID", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservations, err := reservationService.GetReservationsByGuestID(uint(guestID))
		if err != nil {
			log.Printf("handleGetReservationsByGuest: error fetching reservations: %v", err)
			response.JSON(c, "Failed to fetch reservations", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Reservations retrieved successfully", http.StatusOK, reservations, nil)
	}
}

// handleGetReservationsByRoom returns all reservations for a room
func (s *Server) handleGetReservationsByRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := strconv.ParseUint(c.Param("roomID"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid room ID", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		reservations, err := reservationService.GetReservationsByRoomID(uint(roomID))
		if err != nil {
			log.Printf("handleGetReservationsByRoom: error fetching reservations: %v", err)
			response.JSON(c, "Failed to fetch reservations", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Reservations retrieved successfully", http.StatusOK, reservations, nil)
	}
}

// handleDeleteReservation deletes (cancels) a reservation
func (s *Server) handleDeleteReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		err = reservationService.CancelReservation(uint(id), nil)
		if err != nil {
			log.Printf("handleDeleteReservation: error deleting reservation: %v", err)
			response.JSON(c, "Failed to delete reservation", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Reservation deleted successfully", http.StatusOK, nil, nil)
	}
}

// handleGetCheckInStats returns check-in statistics for a specific date
func (s *Server) handleGetCheckInStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Query("date")

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		stats, err := reservationService.GetCheckInStats(date)
		if err != nil {
			log.Printf("handleGetCheckInStats: error fetching stats: %v", err)
			response.JSON(c, "Failed to fetch check-in statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Check-in statistics retrieved", http.StatusOK, stats, nil)
	}
}

// handleGetPaymentDetails returns payment details for a reservation
func (s *Server) handleGetPaymentDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid reservation ID", http.StatusBadRequest, nil, err)
			return
		}

		reservationService := services.NewReservationService(s.ReservationRepository, s.RoomRepository, s.GuestRepository)
		payment, err := reservationService.GetPaymentDetails(uint(id))
		if err != nil {
			log.Printf("handleGetPaymentDetails: error fetching payment: %v", err)
			response.JSON(c, "Failed to fetch payment details", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Payment details retrieved", http.StatusOK, payment, nil)
	}
}

// handleGetDashboardStats returns dashboard statistics for today
func (s *Server) handleGetDashboardStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := s.ReservationRepository.GetDashboardStats()
		if err != nil {
			log.Printf("handleGetDashboardStats: error fetching stats: %v", err)
			response.JSON(c, "Failed to fetch dashboard statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Dashboard statistics retrieved successfully", http.StatusOK, stats, nil)
	}
}

// handleGetRecentActivity returns the 3 newest bookings and 3 guest preferences
func (s *Server) handleGetRecentActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get 3 newest reservations with guest and room details
		reservations, err := s.ReservationRepository.GetRecentReservations(3)
		if err != nil {
			log.Printf("handleGetRecentActivity: error fetching reservations: %v", err)
			response.JSON(c, "Failed to fetch recent activity", http.StatusInternalServerError, nil, err)
			return
		}

		// Format bookings
		type RecentBooking struct {
			ID          uint    `json:"id"`
			Date        string  `json:"date"`
			GuestName   string  `json:"guest_name"`
			RoomNumber  string  `json:"room_number"`
			RoomType    string  `json:"room_type"`
			CheckIn     string  `json:"check_in"`
			CheckOut    string  `json:"check_out"`
			Status      string  `json:"status"`
			TotalAmount float64 `json:"total_amount"`
		}

		recentBookings := make([]RecentBooking, 0, len(reservations))
		for _, r := range reservations {
			booking := RecentBooking{
				ID:          r.ID,
				Date:        r.CreatedAt.Format("2006-01-02"),
				Status:      string(r.Status),
				CheckIn:     r.CheckInDate.Format("2006-01-02"),
				CheckOut:    r.CheckOutDate.Format("2006-01-02"),
				TotalAmount: r.TotalAmount,
			}
			if r.Guest != nil {
				booking.GuestName = r.Guest.Name
			}
			if r.Room != nil {
				booking.RoomNumber = r.Room.RoomNumber
				booking.RoomType = r.Room.RoomType
			}
			recentBookings = append(recentBookings, booking)
		}

		// Get 3 recent guest preferences
		type GuestPreference struct {
			GuestID         uint     `json:"guest_id"`
			GuestName       string   `json:"guest_name"`
			RoomFloors      []string `json:"room_floors"`
			MealTypes       []string `json:"meal_types"`
			RoomTypes       []string `json:"room_types"`
			SpecialRequests []string `json:"special_requests"`
		}

		preferences, err := s.GuestRepository.GetRecentGuestPreferences(3)
		if err != nil {
			log.Printf("handleGetRecentActivity: error fetching preferences: %v", err)
			// Continue with empty preferences instead of failing
			preferences = nil
		}

		recentPreferences := make([]GuestPreference, 0)
		if preferences != nil {
			for _, p := range preferences {
				pref := GuestPreference{
					GuestID:         p.GuestID,
					RoomFloors:      []string(p.RoomFloors),
					MealTypes:       []string(p.MealTypes),
					RoomTypes:       []string(p.RoomTypes),
					SpecialRequests: []string(p.SpecialRequests),
				}
				if p.Guest != nil {
					pref.GuestName = p.Guest.Name
				}
				recentPreferences = append(recentPreferences, pref)
			}
		}

		response.JSON(c, "Recent activity retrieved successfully", http.StatusOK, gin.H{
			"recent_bookings":   recentBookings,
			"guest_preferences": recentPreferences,
		}, nil)
	}
}
