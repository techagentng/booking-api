package services

import (
	"errors"
	"hotel/models"
	"regexp"
	"strings"
	"time"
)

// HallBookingService defines hall booking business logic operations
type HallBookingService interface {
	CreateHallBooking(req *models.CreateHallBookingRequest) (*models.HallBookingResponse, error)
	GetHallBookingByID(id uint) (*models.HallBookingResponse, error)
	GetHallBookingByBookingID(bookingID string) (*models.HallBookingResponse, error)
	GetAllHallBookings(page, pageSize int, search string) (*models.HallBookingListResponse, error)
	UpdateHallBooking(id uint, req *models.UpdateHallBookingRequest) (*models.HallBookingResponse, error)
	DeleteHallBooking(id uint) error
	CheckHallAvailability(date, startTime, endTime string) (bool, error)
	GetHallAvailability(date string) (*models.HallAvailability, error)
	GetHallBookingsByDate(date string) ([]models.HallBookingResponse, error)
	GetHallBookingsByDateRange(startDate, endDate string) ([]models.HallBookingResponse, error)
}

// hallBookingService implements HallBookingService
type hallBookingService struct {
	repo HallBookingRepository
}

// HallBookingRepository interface for dependency injection
type HallBookingRepository interface {
	CreateHallBooking(booking *models.HallBooking) (*models.HallBooking, error)
	GetHallBookingByID(id uint) (*models.HallBooking, error)
	GetHallBookingByBookingID(bookingID string) (*models.HallBooking, error)
	GetAllHallBookings(page, pageSize int, search string) ([]models.HallBooking, int64, error)
	UpdateHallBooking(id uint, booking *models.HallBooking) (*models.HallBooking, error)
	DeleteHallBooking(id uint) error
	CheckHallAvailability(date string, startTime, endTime string) (bool, error)
	GetHallAvailability(date string) (*models.HallAvailability, error)
	GetHallBookingsByDate(date string) ([]models.HallBooking, error)
	GetHallBookingsByDateRange(startDate, endDate string) ([]models.HallBooking, error)
}

// NewHallBookingService creates a new hall booking service
func NewHallBookingService(repo HallBookingRepository) HallBookingService {
	return &hallBookingService{repo: repo}
}

// CreateHallBooking creates a new hall booking with validation
func (s *hallBookingService) CreateHallBooking(req *models.CreateHallBookingRequest) (*models.HallBookingResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Parse booking date
	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return nil, errors.New("invalid booking date format")
	}

	// Create hall booking model
	booking := &models.HallBooking{
		OrganizerName:   strings.TrimSpace(req.OrganizerName),
		OrganizerEmail:  strings.ToLower(strings.TrimSpace(req.OrganizerEmail)),
		OrganizerPhone:  strings.TrimSpace(req.OrganizerPhone),
		EventType:       req.EventType,
		GuestCount:      req.GuestCount,
		SpecialRequests: strings.TrimSpace(req.SpecialRequests),
		BookingDate:     bookingDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		TotalPrice:      req.TotalPrice,
		DepositRequired: req.DepositRequired,
		PaymentMethod:   req.PaymentMethod,
		Status:          "pending",
	}

	// Create booking
	createdBooking, err := s.repo.CreateHallBooking(booking)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(createdBooking), nil
}

// GetHallBookingByID retrieves a hall booking by ID
func (s *hallBookingService) GetHallBookingByID(id uint) (*models.HallBookingResponse, error) {
	booking, err := s.repo.GetHallBookingByID(id)
	if err != nil {
		return nil, err
	}
	return s.convertToResponse(booking), nil
}

// GetHallBookingByBookingID retrieves a hall booking by booking ID
func (s *hallBookingService) GetHallBookingByBookingID(bookingID string) (*models.HallBookingResponse, error) {
	booking, err := s.repo.GetHallBookingByBookingID(bookingID)
	if err != nil {
		return nil, err
	}
	return s.convertToResponse(booking), nil
}

// GetAllHallBookings retrieves all hall bookings with pagination
func (s *hallBookingService) GetAllHallBookings(page, pageSize int, search string) (*models.HallBookingListResponse, error) {
	bookings, total, err := s.repo.GetAllHallBookings(page, pageSize, search)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var bookingResponses []models.HallBookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, *s.convertToResponse(&booking))
	}

	// Calculate pagination
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.HallBookingListResponse{
		Data: bookingResponses,
		Meta: models.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateHallBooking updates a hall booking
func (s *hallBookingService) UpdateHallBooking(id uint, req *models.UpdateHallBookingRequest) (*models.HallBookingResponse, error) {
	// Get existing booking
	existing, err := s.repo.GetHallBookingByID(id)
	if err != nil {
		return nil, err
	}

	// Validate update request
	if err := s.validateUpdateRequest(req, existing); err != nil {
		return nil, err
	}

	// Create update model
	update := &models.HallBooking{
		OrganizerName:   existing.OrganizerName,
		OrganizerEmail:  existing.OrganizerEmail,
		OrganizerPhone:  existing.OrganizerPhone,
		EventType:       existing.EventType,
		GuestCount:      existing.GuestCount,
		SpecialRequests: existing.SpecialRequests,
		BookingDate:     existing.BookingDate,
		StartTime:       existing.StartTime,
		EndTime:         existing.EndTime,
		TotalPrice:      existing.TotalPrice,
		DepositRequired: existing.DepositRequired,
		PaymentMethod:   existing.PaymentMethod,
		Status:          existing.Status,
	}

	// Apply updates
	if req.OrganizerName != nil {
		update.OrganizerName = strings.TrimSpace(*req.OrganizerName)
	}
	if req.OrganizerEmail != nil {
		update.OrganizerEmail = strings.ToLower(strings.TrimSpace(*req.OrganizerEmail))
	}
	if req.OrganizerPhone != nil {
		update.OrganizerPhone = strings.TrimSpace(*req.OrganizerPhone)
	}
	if req.EventType != nil {
		update.EventType = *req.EventType
	}
	if req.GuestCount != nil {
		update.GuestCount = *req.GuestCount
	}
	if req.SpecialRequests != nil {
		update.SpecialRequests = strings.TrimSpace(*req.SpecialRequests)
	}
	if req.BookingDate != nil {
		bookingDate, err := time.Parse("2006-01-02", *req.BookingDate)
		if err != nil {
			return nil, errors.New("invalid booking date format")
		}
		update.BookingDate = bookingDate
	}
	if req.StartTime != nil {
		update.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		update.EndTime = *req.EndTime
	}
	if req.TotalPrice != nil {
		update.TotalPrice = *req.TotalPrice
	}
	if req.DepositRequired != nil {
		update.DepositRequired = *req.DepositRequired
	}
	if req.PaymentMethod != nil {
		update.PaymentMethod = *req.PaymentMethod
	}
	if req.Status != nil {
		update.Status = *req.Status
	}

	// Update booking
	updatedBooking, err := s.repo.UpdateHallBooking(id, update)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(updatedBooking), nil
}

// DeleteHallBooking deletes a hall booking
func (s *hallBookingService) DeleteHallBooking(id uint) error {
	return s.repo.DeleteHallBooking(id)
}

// CheckHallAvailability checks if the hall is available for the given date and time
func (s *hallBookingService) CheckHallAvailability(date, startTime, endTime string) (bool, error) {
	// Validate inputs
	if err := s.validateDateTimeInputs(date, startTime, endTime); err != nil {
		return false, err
	}

	return s.repo.CheckHallAvailability(date, startTime, endTime)
}

// GetHallAvailability gets availability for a specific date
func (s *hallBookingService) GetHallAvailability(date string) (*models.HallAvailability, error) {
	// Validate date
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, errors.New("invalid date format")
	}

	return s.repo.GetHallAvailability(date)
}

// GetHallBookingsByDate retrieves all hall bookings for a specific date
func (s *hallBookingService) GetHallBookingsByDate(date string) ([]models.HallBookingResponse, error) {
	// Validate date
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, errors.New("invalid date format")
	}

	bookings, err := s.repo.GetHallBookingsByDate(date)
	if err != nil {
		return nil, err
	}

	var responses []models.HallBookingResponse
	for _, booking := range bookings {
		responses = append(responses, *s.convertToResponse(&booking))
	}

	return responses, nil
}

// GetHallBookingsByDateRange retrieves all hall bookings within a date range
func (s *hallBookingService) GetHallBookingsByDateRange(startDate, endDate string) ([]models.HallBookingResponse, error) {
	// Validate dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	if start.After(end) {
		return nil, errors.New("start date cannot be after end date")
	}

	bookings, err := s.repo.GetHallBookingsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var responses []models.HallBookingResponse
	for _, booking := range bookings {
		responses = append(responses, *s.convertToResponse(&booking))
	}

	return responses, nil
}

// Helper functions

// validateCreateRequest validates the create request
func (s *hallBookingService) validateCreateRequest(req *models.CreateHallBookingRequest) error {
	// Validate booking date is in the future
	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return errors.New("invalid booking date format")
	}

	if bookingDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("booking date must be in the future")
	}

	// Validate time logic
	if req.StartTime >= req.EndTime {
		return errors.New("end time must be after start time")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.OrganizerEmail) {
		return errors.New("invalid email format")
	}

	// Validate phone number (basic validation)
	if len(strings.ReplaceAll(req.OrganizerPhone, " ", "")) < 10 {
		return errors.New("phone number must be at least 10 digits")
	}

	// Validate pricing
	if req.DepositRequired > req.TotalPrice {
		return errors.New("deposit cannot be greater than total price")
	}

	if req.TotalPrice <= 0 {
		return errors.New("total price must be greater than 0")
	}

	return nil
}

// validateUpdateRequest validates the update request
func (s *hallBookingService) validateUpdateRequest(req *models.UpdateHallBookingRequest, existing *models.HallBooking) error {
	// Validate email if provided
	if req.OrganizerEmail != nil {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(*req.OrganizerEmail) {
			return errors.New("invalid email format")
		}
	}

	// Validate phone if provided
	if req.OrganizerPhone != nil {
		if len(strings.ReplaceAll(*req.OrganizerPhone, " ", "")) < 10 {
			return errors.New("phone number must be at least 10 digits")
		}
	}

	// Validate date and time if provided
	if req.BookingDate != nil || req.StartTime != nil || req.EndTime != nil {
		date := existing.BookingDate.Format("2006-01-02")
		startTime := existing.StartTime
		endTime := existing.EndTime

		if req.BookingDate != nil {
			date = *req.BookingDate
		}
		if req.StartTime != nil {
			startTime = *req.StartTime
		}
		if req.EndTime != nil {
			endTime = *req.EndTime
		}

		if err := s.validateDateTimeInputs(date, startTime, endTime); err != nil {
			return err
		}
	}

	// Validate pricing if provided
	if req.TotalPrice != nil && req.DepositRequired != nil {
		if *req.DepositRequired > *req.TotalPrice {
			return errors.New("deposit cannot be greater than total price")
		}
		if *req.TotalPrice <= 0 {
			return errors.New("total price must be greater than 0")
		}
	}

	return nil
}

// validateDateTimeInputs validates date and time inputs
func (s *hallBookingService) validateDateTimeInputs(date, startTime, endTime string) error {
	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return errors.New("invalid date format")
	}

	// Validate time format
	if _, err := time.Parse("15:04", startTime); err != nil {
		return errors.New("invalid start time format")
	}

	if _, err := time.Parse("15:04", endTime); err != nil {
		return errors.New("invalid end time format")
	}

	// Validate time logic
	if startTime >= endTime {
		return errors.New("end time must be after start time")
	}

	return nil
}

// convertToResponse converts a model to response format
func (s *hallBookingService) convertToResponse(booking *models.HallBooking) *models.HallBookingResponse {
	return &models.HallBookingResponse{
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
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
		Payments:        booking.Payments,
		Invoice:         booking.Invoice,
	}
}
