package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"time"
)

// RoomService defines room business logic
type RoomService interface {
	CreateRoom(req *models.CreateRoomRequest) (*models.Room, error)
	GetRoomByID(id uint) (*models.Room, error)
	GetAllRooms(page, pageSize int, search, status string) (*models.RoomListResponse, error)
	UpdateRoom(id uint, req *models.UpdateRoomRequest) (*models.Room, error)
	DeleteRoom(id uint) error
	CheckAvailability(req *models.RoomAvailabilityRequest) ([]models.Room, error)
	GetAvailableRoomsQuery(checkIn, checkOut, roomType string) ([]models.Room, error)
	GetRoomsByType(roomType string) ([]models.Room, error)
	UpdateRoomStatus(id uint, status string) error
	GetRoomTypeStats() ([]models.RoomTypeStats, error)
}

// roomService implements RoomService
type roomService struct {
	roomRepo db.RoomRepository
}

// NewRoomService creates a new room service
func NewRoomService(roomRepo db.RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

// CreateRoom creates a new room with validation
func (s *roomService) CreateRoom(req *models.CreateRoomRequest) (*models.Room, error) {
	// Validate input
	if err := s.validateCreateRoomRequest(req); err != nil {
		return nil, err
	}

	// Check if room number already exists
	existingRoom, err := s.roomRepo.GetRoomByNumber(req.RoomNumber)
	if err == nil && existingRoom != nil {
		return nil, errors.New("room number already exists")
	}

	// Create room with "available" status
	room := &models.Room{
		RoomNumber:    req.RoomNumber,
		RoomType:      req.RoomType,
		Floor:         req.Floor,
		Capacity:      req.Capacity,
		PricePerNight: req.PricePerNight,
		Status:        "available",
		Description:   req.Description,
		Amenities:     req.Amenities,
		BedType:       req.BedType,
	}

	return s.roomRepo.CreateRoom(room)
}

// GetRoomByID retrieves a room by ID
func (s *roomService) GetRoomByID(id uint) (*models.Room, error) {
	if id == 0 {
		return nil, errors.New("invalid room ID")
	}
	return s.roomRepo.GetRoomByID(id)
}

// GetAllRooms retrieves all rooms with pagination and optional status filter
func (s *roomService) GetAllRooms(page, pageSize int, search, status string) (*models.RoomListResponse, error) {
	rooms, total, err := s.roomRepo.GetAllRooms(page, pageSize, search, status)
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

	return &models.RoomListResponse{
		Rooms: rooms,
		Meta: models.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateRoom updates a room
func (s *roomService) UpdateRoom(id uint, req *models.UpdateRoomRequest) (*models.Room, error) {
	if id == 0 {
		return nil, errors.New("invalid room ID")
	}

	// Get existing room
	existingRoom, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return nil, err
	}

	// Validate room number uniqueness if updating
	if req.RoomNumber != nil && *req.RoomNumber != existingRoom.RoomNumber {
		checkRoom, err := s.roomRepo.GetRoomByNumber(*req.RoomNumber)
		if err == nil && checkRoom != nil {
			return nil, errors.New("room number already exists")
		}
	}

	// Validate status if provided
	if req.Status != nil {
		validStatuses := map[string]bool{
			"available":   true,
			"occupied":    true,
			"maintenance": true,
			"cleaning":    true,
		}
		if !validStatuses[*req.Status] {
			return nil, errors.New("invalid status")
		}
	}

	// Update fields
	updateRoom := &models.Room{}
	if req.RoomNumber != nil {
		updateRoom.RoomNumber = *req.RoomNumber
	}
	if req.RoomType != nil {
		updateRoom.RoomType = *req.RoomType
	}
	if req.Floor != nil {
		updateRoom.Floor = *req.Floor
	}
	if req.PricePerNight != nil {
		updateRoom.PricePerNight = *req.PricePerNight
	}
	if req.Status != nil {
		updateRoom.Status = *req.Status
	}
	if req.Description != nil {
		updateRoom.Description = *req.Description
	}
	if req.Amenities != nil {
		updateRoom.Amenities = *req.Amenities
	}
	if req.Capacity != nil {
		updateRoom.Capacity = *req.Capacity
	}
	if req.BedType != nil {
		updateRoom.BedType = *req.BedType
	}

	return s.roomRepo.UpdateRoom(id, updateRoom)
}

// DeleteRoom deletes a room
func (s *roomService) DeleteRoom(id uint) error {
	if id == 0 {
		return errors.New("invalid room ID")
	}
	return s.roomRepo.DeleteRoom(id)
}

// CheckAvailability checks room availability for a date range
func (s *roomService) CheckAvailability(req *models.RoomAvailabilityRequest) ([]models.Room, error) {
	// Parse dates
	checkIn, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return nil, errors.New("invalid check_in_date format, use YYYY-MM-DD")
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return nil, errors.New("invalid check_out_date format, use YYYY-MM-DD")
	}

	// Validate date range
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		return nil, errors.New("check_out_date must be after check_in_date")
	}

	// Check if dates are in the past
	today := time.Now().Truncate(24 * time.Hour)
	if checkIn.Before(today) {
		return nil, errors.New("check_in_date cannot be in the past")
	}

	return s.roomRepo.GetAvailableRooms(req.RoomType, checkIn, checkOut)
}

// GetRoomsByType retrieves all rooms of a specific type
func (s *roomService) GetRoomsByType(roomType string) ([]models.Room, error) {
	if roomType == "" {
		return nil, errors.New("room type is required")
	}
	return s.roomRepo.GetRoomsByType(roomType)
}

// UpdateRoomStatus updates a room's status
func (s *roomService) UpdateRoomStatus(id uint, status string) error {
	if id == 0 {
		return errors.New("invalid room ID")
	}

	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"maintenance": true,
		"cleaning":    true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	return s.roomRepo.UpdateRoomStatus(id, status)
}

// GetRoomTypeStats retrieves statistics for each room type
func (s *roomService) GetRoomTypeStats() ([]models.RoomTypeStats, error) {
	return s.roomRepo.GetRoomTypeStats()
}

// GetAvailableRoomsQuery checks room availability using query params (GET version)
func (s *roomService) GetAvailableRoomsQuery(checkIn, checkOut, roomType string) ([]models.Room, error) {
	// Parse dates
	checkInDate, err := time.Parse("2006-01-02", checkIn)
	if err != nil {
		return nil, errors.New("invalid check_in date format, use YYYY-MM-DD")
	}

	checkOutDate, err := time.Parse("2006-01-02", checkOut)
	if err != nil {
		return nil, errors.New("invalid check_out date format, use YYYY-MM-DD")
	}

	// Validate date range
	if checkOutDate.Before(checkInDate) || checkOutDate.Equal(checkInDate) {
		return nil, errors.New("check_out must be after check_in")
	}

	// Check if dates are in the past
	today := time.Now().Truncate(24 * time.Hour)
	if checkInDate.Before(today) {
		return nil, errors.New("check_in cannot be in the past")
	}

	return s.roomRepo.GetAvailableRooms(roomType, checkInDate, checkOutDate)
}

// validateCreateRoomRequest validates room creation request
func (s *roomService) validateCreateRoomRequest(req *models.CreateRoomRequest) error {
	if req.RoomNumber == "" {
		return errors.New("room_number is required")
	}

	if req.RoomType == "" {
		return errors.New("room_type is required")
	}

	validRoomTypes := map[string]bool{
		"Standard": true,
		"Deluxe":   true,
		"Suite":    true,
		"standard": true,
		"deluxe":   true,
		"suite":    true,
	}
	if !validRoomTypes[req.RoomType] {
		return errors.New("invalid room_type, must be Standard, Deluxe, or Suite")
	}

	if req.Floor <= 0 {
		return errors.New("floor must be a positive integer")
	}

	if req.Capacity <= 0 {
		return errors.New("capacity must be a positive integer")
	}

	if req.PricePerNight <= 0 {
		return errors.New("price_per_night must be greater than 0")
	}

	return nil
}
