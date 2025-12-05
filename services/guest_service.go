package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"net/mail"
)

// GuestService defines guest business logic
type GuestService interface {
	CreateGuest(req *models.CreateGuestRequest) (*models.Guest, error)
	GetGuestByID(id uint) (*models.Guest, error)
	GetGuestDetailsByID(id uint) (*models.GuestDetailsResponse, error)
	GetAllGuests(page, pageSize int, search string) (*models.GuestListResponse, error)
	UpdateGuest(id uint, req *models.UpdateGuestRequest) (*models.Guest, error)
	DeleteGuest(id uint) error
	GetGuestHistory(guestID uint) (*models.GuestHistory, error)
	GetGuestPreferences(guestID uint) (*models.GuestPreferences, error)
	GetGuestAIInsights(guestID uint) (*models.GuestAIInsights, error)
}

// guestService implements GuestService
type guestService struct {
	guestRepo db.GuestRepository
}

// NewGuestService creates a new guest service
func NewGuestService(guestRepo db.GuestRepository) GuestService {
	return &guestService{
		guestRepo: guestRepo,
	}
}

// CreateGuest creates a new guest with validation
func (s *guestService) CreateGuest(req *models.CreateGuestRequest) (*models.Guest, error) {
	// Validate input
	if err := s.validateCreateGuestRequest(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	existingGuest, err := s.guestRepo.GetGuestByEmail(req.Email)
	if err == nil && existingGuest != nil {
		return nil, errors.New("email already exists")
	}

	// Create guest
	guest := &models.Guest{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Nationality: req.Nationality,
		IDType:      req.IDType,
		IDNumber:    req.IDNumber,
	}

	return s.guestRepo.CreateGuest(guest)
}

// GetGuestByID retrieves a guest by ID
func (s *guestService) GetGuestByID(id uint) (*models.Guest, error) {
	if id == 0 {
		return nil, errors.New("invalid guest ID")
	}
	return s.guestRepo.GetGuestByID(id)
}

// GetAllGuests retrieves all guests with pagination and search
func (s *guestService) GetAllGuests(page, pageSize int, search string) (*models.GuestListResponse, error) {
	guests, total, err := s.guestRepo.GetAllGuests(page, pageSize, search)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.GuestListResponse{
		Guests: guests,
		Meta: models.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetGuestDetailsByID retrieves a guest with all details including statistics and service usage
func (s *guestService) GetGuestDetailsByID(id uint) (*models.GuestDetailsResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid guest ID")
	}
	return s.guestRepo.GetGuestDetailsByID(id)
}

// UpdateGuest updates a guest
func (s *guestService) UpdateGuest(id uint, req *models.UpdateGuestRequest) (*models.Guest, error) {
	if id == 0 {
		return nil, errors.New("invalid guest ID")
	}

	// Get existing guest
	existingGuest, err := s.guestRepo.GetGuestByID(id)
	if err != nil {
		return nil, err
	}

	// Validate email uniqueness if updating
	if req.Email != nil && *req.Email != existingGuest.Email {
		checkGuest, err := s.guestRepo.GetGuestByEmail(*req.Email)
		if err == nil && checkGuest != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Update fields
	updateGuest := &models.Guest{}
	if req.Name != nil {
		updateGuest.Name = *req.Name
	}
	if req.Email != nil {
		updateGuest.Email = *req.Email
	}
	if req.Phone != nil {
		updateGuest.Phone = *req.Phone
	}
	if req.Nationality != nil {
		updateGuest.Nationality = *req.Nationality
	}
	if req.IDType != nil {
		updateGuest.IDType = *req.IDType
	}
	if req.IDNumber != nil {
		updateGuest.IDNumber = *req.IDNumber
	}

	return s.guestRepo.UpdateGuest(id, updateGuest)
}

// DeleteGuest deletes a guest
func (s *guestService) DeleteGuest(id uint) error {
	if id == 0 {
		return errors.New("invalid guest ID")
	}
	return s.guestRepo.DeleteGuest(id)
}

// GetGuestHistory retrieves guest history
func (s *guestService) GetGuestHistory(guestID uint) (*models.GuestHistory, error) {
	if guestID == 0 {
		return nil, errors.New("invalid guest ID")
	}

	// Verify guest exists
	_, err := s.guestRepo.GetGuestByID(guestID)
	if err != nil {
		return nil, err
	}

	return s.guestRepo.GetGuestHistory(guestID)
}

// GetGuestPreferences retrieves guest preferences
func (s *guestService) GetGuestPreferences(guestID uint) (*models.GuestPreferences, error) {
	if guestID == 0 {
		return nil, errors.New("invalid guest ID")
	}

	// Verify guest exists
	_, err := s.guestRepo.GetGuestByID(guestID)
	if err != nil {
		return nil, err
	}

	return s.guestRepo.GetGuestPreferences(guestID)
}

// GetGuestAIInsights retrieves guest AI insights
func (s *guestService) GetGuestAIInsights(guestID uint) (*models.GuestAIInsights, error) {
	if guestID == 0 {
		return nil, errors.New("invalid guest ID")
	}

	// Verify guest exists
	_, err := s.guestRepo.GetGuestByID(guestID)
	if err != nil {
		return nil, err
	}

	return s.guestRepo.GetGuestAIInsights(guestID)
}

// validateCreateGuestRequest validates guest creation request
func (s *guestService) validateCreateGuestRequest(req *models.CreateGuestRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}

	if req.Phone == "" {
		return errors.New("phone is required")
	}

	if req.Nationality == "" {
		return errors.New("nationality is required")
	}

	if req.IDType == "" {
		return errors.New("id_type is required")
	}

	if req.IDNumber == "" {
		return errors.New("id_number is required")
	}

	return nil
}
