package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"time"
)

// AdminHallBookingService interface for admin booking operations
type AdminHallBookingService interface {
	// Booking management
	GetAllBookings(page, pageSize int, status []string, dateFrom, dateTo, search, sortBy, sortOrder string) ([]models.HallBooking, int64, error)
	GetBookingByID(id uint) (*models.HallBooking, error)
	UpdateBookingStatus(id uint, req *models.AdminBookingRequest, adminID uint) (*models.HallBooking, error)

	// Status history
	GetBookingStatusHistory(bookingID uint) ([]models.BookingStatusHistory, error)

	// Statistics
	GetBookingStats(period, dateFrom, dateTo string) (map[string]interface{}, error)

	// Response conversion
	ConvertToAdminResponse(booking *models.HallBooking) *models.AdminBookingResponse
}

type adminHallBookingService struct {
	repo db.AdminHallBookingRepository
}

// NewAdminHallBookingService creates a new admin hall booking service
func NewAdminHallBookingService(repo db.AdminHallBookingRepository) AdminHallBookingService {
	return &adminHallBookingService{repo: repo}
}

// GetAllBookings retrieves all bookings with filtering and pagination
func (s *adminHallBookingService) GetAllBookings(page, pageSize int, status []string, dateFrom, dateTo, search, sortBy, sortOrder string) ([]models.HallBooking, int64, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Validate sort field
	validSortFields := map[string]bool{
		"date":           true,
		"created_at":     true,
		"status":         true,
		"organizer_name": true,
		"total_price":    true,
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	return s.repo.GetAllBookings(page, pageSize, status, dateFrom, dateTo, search, sortBy, sortOrder)
}

// GetBookingByID retrieves a booking by ID
func (s *adminHallBookingService) GetBookingByID(id uint) (*models.HallBooking, error) {
	if id == 0 {
		return nil, errors.New("invalid booking ID")
	}

	return s.repo.GetBookingByID(id)
}

// UpdateBookingStatus updates booking status with admin tracking
func (s *adminHallBookingService) UpdateBookingStatus(id uint, req *models.AdminBookingRequest, adminID uint) (*models.HallBooking, error) {
	if id == 0 {
		return nil, errors.New("invalid booking ID")
	}

	if adminID == 0 {
		return nil, errors.New("invalid admin user ID")
	}

	// Validate status transition
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return nil, err
	}

	// Check if status transition is valid
	if !s.isValidStatusTransition(booking.Status, req.Status) {
		return nil, errors.New("invalid status transition from " + booking.Status + " to " + req.Status)
	}

	// Update status
	return s.repo.UpdateBookingStatus(id, req.Status, adminID, req.Notes)
}

// GetBookingStatusHistory retrieves status history for a booking
func (s *adminHallBookingService) GetBookingStatusHistory(bookingID uint) ([]models.BookingStatusHistory, error) {
	if bookingID == 0 {
		return nil, errors.New("invalid booking ID")
	}

	return s.repo.GetBookingStatusHistory(bookingID)
}

// GetBookingStats retrieves booking statistics
func (s *adminHallBookingService) GetBookingStats(period, dateFrom, dateTo string) (map[string]interface{}, error) {
	// Validate period
	validPeriods := map[string]bool{
		"today":  true,
		"week":   true,
		"month":  true,
		"year":   true,
		"custom": true,
	}

	if !validPeriods[period] {
		period = "month" // default
	}

	// Validate custom date range
	if period == "custom" {
		if dateFrom == "" || dateTo == "" {
			return nil, errors.New("both date_from and date_to are required for custom period")
		}

		// Validate date format
		fromDate, err := parseDate(dateFrom)
		if err != nil {
			return nil, errors.New("invalid date_from format, use YYYY-MM-DD")
		}

		toDate, err := parseDate(dateTo)
		if err != nil {
			return nil, errors.New("invalid date_to format, use YYYY-MM-DD")
		}

		// Validate date range
		if fromDate.After(toDate) {
			return nil, errors.New("date_from must be before or equal to date_to")
		}
	}

	return s.repo.GetBookingStats(period, dateFrom, dateTo)
}

// ConvertToAdminResponse converts booking to admin response format
func (s *adminHallBookingService) ConvertToAdminResponse(booking *models.HallBooking) *models.AdminBookingResponse {
	return &models.AdminBookingResponse{
		ID:              booking.ID,
		BookingID:       booking.BookingID,
		OrganizerName:   booking.OrganizerName,
		OrganizerEmail:  booking.OrganizerEmail,
		OrganizerPhone:  booking.OrganizerPhone,
		EventType:       booking.EventType,
		GuestCount:      booking.GuestCount,
		SpecialRequests: booking.SpecialRequests,
		BookingDate:     booking.BookingDate.Format("2006-01-02"),
		StartTime:       booking.StartTime,
		EndTime:         booking.EndTime,
		TotalPrice:      booking.TotalPrice,
		DepositRequired: booking.DepositRequired,
		PaymentMethod:   booking.PaymentMethod,
		Status:          booking.Status,
		// ConfirmedBy:     booking.ConfirmedByUser,  // Temporarily disabled
		ConfirmedAt: booking.ConfirmedAt,
		// CancelledBy:     booking.CancelledByUser,  // Temporarily disabled
		CancelledAt: booking.CancelledAt,
		// UpdatedBy:       booking.UpdatedByUser,    // Temporarily disabled
		// StatusHistory:   booking.StatusHistory,    // Temporarily disabled
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}
}

// isValidStatusTransition checks if a status transition is valid
func (s *adminHallBookingService) isValidStatusTransition(oldStatus, newStatus string) bool {
	// Allow same status (idempotent)
	if oldStatus == newStatus {
		return true
	}

	// Define valid transitions
	validTransitions := map[string][]string{
		"pending":   {"confirmed", "cancelled"},
		"confirmed": {"completed", "cancelled"},
		"completed": {},          // No transitions from completed
		"cancelled": {"pending"}, // Can reactivate cancelled bookings
	}

	allowedStatuses, exists := validTransitions[oldStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
