package db

import (
	"errors"
	"fmt"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// GuestServiceRepository defines guest service database operations
type GuestServiceRepository interface {
	// Service request operations
	CreateServiceRequest(request *models.GuestServiceRequest) (*models.GuestServiceRequest, error)
	GetServiceRequestByID(id uint) (*models.GuestServiceRequest, error)
	GetServiceRequestsByGuestID(guestID uint, status string) ([]models.GuestServiceRequest, error)
	GetServiceRequestsByRoomID(roomID uint, status string) ([]models.GuestServiceRequest, error)
	GetAllServiceRequests(page, pageSize int, status, requestType string) ([]models.GuestServiceRequest, int64, error)
	GetAssignedServiceRequests() ([]models.GuestServiceRequest, error)
	UpdateServiceRequestStatus(id uint, status, assignedTo string) error
	AssignStaffToRequest(requestID uint, staffID uint, staffName string) error
	CompleteServiceRequest(id uint, completedBy string) error
	CancelServiceRequest(id uint, reason string) error
	GenerateRequestNumber(requestType string) (string, error)

	// Guest room info
	GetGuestByRoomNumber(roomNumber string) (*models.GuestRoomInfo, error)

	// Menu categories
	GetMenuCategories() ([]models.MenuCategory, error)

	// Guest active orders
	GetGuestActiveOrders(guestID uint) ([]models.RoomServiceOrder, error)
}

// guestServiceRepository implements GuestServiceRepository
type guestServiceRepository struct {
	db *gorm.DB
}

// NewGuestServiceRepository creates a new guest service repository
func NewGuestServiceRepository(db *gorm.DB) GuestServiceRepository {
	return &guestServiceRepository{db: db}
}

// CreateServiceRequest creates a new service request
func (r *guestServiceRepository) CreateServiceRequest(request *models.GuestServiceRequest) (*models.GuestServiceRequest, error) {
	if err := r.db.Create(request).Error; err != nil {
		return nil, err
	}
	return r.GetServiceRequestByID(request.ID)
}

// GetServiceRequestByID retrieves a service request by ID
func (r *guestServiceRepository) GetServiceRequestByID(id uint) (*models.GuestServiceRequest, error) {
	var request models.GuestServiceRequest
	if err := r.db.Preload("Room").Preload("Guest").First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// GetServiceRequestsByGuestID retrieves service requests for a guest
func (r *guestServiceRepository) GetServiceRequestsByGuestID(guestID uint, status string) ([]models.GuestServiceRequest, error) {
	var requests []models.GuestServiceRequest
	query := r.db.Where("guest_id = ?", guestID)

	if status == "active" {
		query = query.Where("status NOT IN ?", []string{"completed", "cancelled"})
	} else if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestsByRoomID retrieves service requests for a room
func (r *guestServiceRepository) GetServiceRequestsByRoomID(roomID uint, status string) ([]models.GuestServiceRequest, error) {
	var requests []models.GuestServiceRequest
	query := r.db.Where("room_id = ?", roomID)

	if status == "active" {
		query = query.Where("status NOT IN ?", []string{"completed", "cancelled"})
	} else if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetAllServiceRequests retrieves all service requests with pagination and filters
func (r *guestServiceRepository) GetAllServiceRequests(page, pageSize int, status, requestType string) ([]models.GuestServiceRequest, int64, error) {
	var requests []models.GuestServiceRequest
	var total int64

	query := r.db.Model(&models.GuestServiceRequest{})

	// Apply filters
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if requestType != "" && requestType != "all" {
		query = query.Where("type = ?", requestType)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.
		Preload("Room").
		Preload("Guest").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetAssignedServiceRequests returns all service requests that have been assigned to staff
func (r *guestServiceRepository) GetAssignedServiceRequests() ([]models.GuestServiceRequest, error) {
	var requests []models.GuestServiceRequest

	if err := r.db.
		Where("assigned_staff_id IS NOT NULL").
		Where("status IN (?, ?)", "in_progress", "pending").
		Preload("Room").
		Preload("Guest").
		Preload("AssignedStaff").
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

// UpdateServiceRequestStatus updates the status of a service request
func (r *guestServiceRepository) UpdateServiceRequestStatus(id uint, status, assignedTo string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if assignedTo != "" {
		updates["assigned_to"] = assignedTo
	}
	return r.db.Model(&models.GuestServiceRequest{}).Where("id = ?", id).Updates(updates).Error
}

// AssignStaffToRequest assigns a staff member to a service request
func (r *guestServiceRepository) AssignStaffToRequest(requestID uint, staffID uint, staffName string) error {
	now := time.Now()
	return r.db.Model(&models.GuestServiceRequest{}).Where("id = ?", requestID).Updates(map[string]interface{}{
		"status":            "in_progress",
		"assigned_staff_id": staffID,
		"assigned_to":       staffName,
		"assigned_at":       now,
	}).Error
}

// CompleteServiceRequest marks a service request as completed
func (r *guestServiceRepository) CompleteServiceRequest(id uint, completedBy string) error {
	now := time.Now()
	return r.db.Model(&models.GuestServiceRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": now,
		"completed_by": completedBy,
	}).Error
}

// CancelServiceRequest cancels a service request
func (r *guestServiceRepository) CancelServiceRequest(id uint, reason string) error {
	return r.db.Model(&models.GuestServiceRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":              "cancelled",
		"cancellation_reason": reason,
	}).Error
}

// GenerateRequestNumber generates a unique request number
func (r *guestServiceRepository) GenerateRequestNumber(requestType string) (string, error) {
	year := time.Now().Year()
	prefix := "SR"
	if requestType == "housekeeping" {
		prefix = "HK"
	} else if requestType == "maintenance" {
		prefix = "MN"
	}

	var count int64
	r.db.Model(&models.GuestServiceRequest{}).
		Where("type = ? AND created_at >= ?", requestType, time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)).
		Count(&count)

	requestNumber := fmt.Sprintf("%s-%d-%03d", prefix, year, count+1)

	// Check if it exists
	var existing models.GuestServiceRequest
	for {
		err := r.db.Where("request_number = ?", requestNumber).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		count++
		requestNumber = fmt.Sprintf("%s-%d-%03d", prefix, year, count+1)
	}

	return requestNumber, nil
}

// GetGuestByRoomNumber retrieves guest info for a room
func (r *guestServiceRepository) GetGuestByRoomNumber(roomNumber string) (*models.GuestRoomInfo, error) {
	// Find the room
	var room models.Room
	if err := r.db.Where("room_number = ?", roomNumber).First(&room).Error; err != nil {
		return nil, errors.New("room not found")
	}

	// Find active reservation for this room (checked-in or confirmed)
	var reservation models.Reservation
	if err := r.db.
		Where("room_id = ? AND status IN ?", room.ID, []string{"checked-in", "confirmed"}).
		Order("check_in_date DESC").
		First(&reservation).Error; err != nil {
		return nil, errors.New("no active reservation found for this room")
	}

	// Get guest
	var guest models.Guest
	if err := r.db.First(&guest, reservation.GuestID).Error; err != nil {
		return nil, errors.New("guest not found")
	}

	// Calculate number of nights
	checkIn, _ := time.Parse("2006-01-02", reservation.CheckInDate.Format("2006-01-02"))
	checkOut, _ := time.Parse("2006-01-02", reservation.CheckOutDate.Format("2006-01-02"))
	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	// Split guest name into first and last
	firstName := guest.Name
	lastName := ""
	// Simple split - could be improved
	for i, c := range guest.Name {
		if c == ' ' {
			firstName = guest.Name[:i]
			lastName = guest.Name[i+1:]
			break
		}
	}

	return &models.GuestRoomInfo{
		Guest: models.GuestRoomInfoGuest{
			ID:        guest.ID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     guest.Email,
			Phone:     guest.Phone,
		},
		Reservation: models.GuestRoomInfoReservation{
			ID:                 reservation.ID,
			ConfirmationNumber: reservation.ConfirmationNumber,
			CheckInDate:        reservation.CheckInDate.Format("2006-01-02"),
			CheckOutDate:       reservation.CheckOutDate.Format("2006-01-02"),
			NumberOfNights:     nights,
		},
		Room: models.GuestRoomInfoRoom{
			ID:         room.ID,
			RoomNumber: room.RoomNumber,
			RoomType:   room.RoomType,
			Floor:      room.Floor,
		},
		IsCheckedIn: reservation.Status == "checked-in",
	}, nil
}

// GetMenuCategories retrieves menu categories with item counts
func (r *guestServiceRepository) GetMenuCategories() ([]models.MenuCategory, error) {
	var results []struct {
		Category string
		Count    int64
	}

	if err := r.db.Model(&models.MenuItem{}).
		Select("category, COUNT(*) as count").
		Where("available = ?", true).
		Group("category").
		Order("category").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// Map categories to icons
	iconMap := map[string]string{
		"Breakfast":   "sunrise",
		"Main Course": "utensils",
		"Salads":      "leaf",
		"Desserts":    "cake",
		"Beverages":   "coffee",
		"Appetizers":  "soup",
		"Soups":       "soup",
		"Snacks":      "cookie",
	}

	var categories []models.MenuCategory
	for _, r := range results {
		icon := iconMap[r.Category]
		if icon == "" {
			icon = "utensils"
		}
		categories = append(categories, models.MenuCategory{
			Name:       r.Category,
			ItemsCount: r.Count,
			Icon:       icon,
		})
	}

	return categories, nil
}

// GetGuestActiveOrders retrieves active orders for a guest
func (r *guestServiceRepository) GetGuestActiveOrders(guestID uint) ([]models.RoomServiceOrder, error) {
	var orders []models.RoomServiceOrder
	if err := r.db.
		Preload("Items").
		Where("guest_id = ? AND status NOT IN ?", guestID, []string{"delivered", "cancelled"}).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
