package db

import (
	"errors"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// RoomRepository defines room database operations
type RoomRepository interface {
	CreateRoom(room *models.Room) (*models.Room, error)
	GetRoomByID(id uint) (*models.Room, error)
	GetAllRooms(page, pageSize int, search, status string) ([]models.Room, int64, error)
	UpdateRoom(id uint, room *models.Room) (*models.Room, error)
	DeleteRoom(id uint) error
	GetRoomByNumber(roomNumber string) (*models.Room, error)
	GetAvailableRooms(roomType string, checkIn, checkOut time.Time) ([]models.Room, error)
	GetRoomsByType(roomType string) ([]models.Room, error)
	GetRoomsByStatus(status string) ([]models.Room, error)
	UpdateRoomStatus(id uint, status string) error
	GetRoomTypeStats() ([]models.RoomTypeStats, error)
	GetRoomStats() (*models.RoomStats, error)
}

// roomRepository implements RoomRepository
type roomRepository struct {
	db *gorm.DB
}

// NewRoomRepository creates a new room repository
func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

// CreateRoom creates a new room
func (r *roomRepository) CreateRoom(room *models.Room) (*models.Room, error) {
	if err := r.db.Create(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

// GetRoomByID retrieves a room by ID
func (r *roomRepository) GetRoomByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := r.db.First(&room, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return &room, nil
}

// GetAllRooms retrieves all rooms with pagination, search, and status filter
func (r *roomRepository) GetAllRooms(page, pageSize int, search, status string) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Room{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("room_number ILIKE ? OR room_type ILIKE ?",
			searchPattern, searchPattern)
	}

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("room_number ASC").
		Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

// UpdateRoom updates a room
func (r *roomRepository) UpdateRoom(id uint, room *models.Room) (*models.Room, error) {
	if err := r.db.Model(&models.Room{}).Where("id = ?", id).Updates(room).Error; err != nil {
		return nil, err
	}
	return r.GetRoomByID(id)
}

// DeleteRoom deletes a room
func (r *roomRepository) DeleteRoom(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

// GetRoomByNumber retrieves a room by room number
func (r *roomRepository) GetRoomByNumber(roomNumber string) (*models.Room, error) {
	var room models.Room
	if err := r.db.Where("room_number = ?", roomNumber).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return &room, nil
}

// GetAvailableRooms retrieves available rooms for a given type and date range
func (r *roomRepository) GetAvailableRooms(roomType string, checkIn, checkOut time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// Find rooms that:
	// 1. Match the room type (case-insensitive)
	// 2. Have status 'available'
	// 3. Don't have overlapping reservations for the date range
	query := r.db.Model(&models.Room{}).
		Where("LOWER(room_type) = LOWER(?) AND status = ?", roomType, "available").
		Where("id NOT IN (?)",
			r.db.Model(&models.Reservation{}).
				Select("room_id").
				Where("status NOT IN ?", []string{"cancelled", "completed", "checked-out"}).
				Where("(check_in_date < ? AND check_out_date > ?)", checkOut, checkIn),
		)

	if err := query.Find(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}

// GetRoomsByType retrieves all rooms of a specific type
func (r *roomRepository) GetRoomsByType(roomType string) ([]models.Room, error) {
	var rooms []models.Room
	if err := r.db.Where("LOWER(room_type) = LOWER(?)", roomType).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetRoomsByStatus retrieves all rooms with a specific status
func (r *roomRepository) GetRoomsByStatus(status string) ([]models.Room, error) {
	var rooms []models.Room
	if err := r.db.Where("status = ?", status).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// UpdateRoomStatus updates a room's status
func (r *roomRepository) UpdateRoomStatus(id uint, status string) error {
	return r.db.Model(&models.Room{}).Where("id = ?", id).Update("status", status).Error
}

// GetRoomTypeStats retrieves statistics for each room type
func (r *roomRepository) GetRoomTypeStats() ([]models.RoomTypeStats, error) {
	var stats []models.RoomTypeStats

	err := r.db.Model(&models.Room{}).
		Select(`
			room_type,
			COUNT(*) as total_rooms,
			SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END) as available,
			SUM(CASE WHEN status = 'occupied' THEN 1 ELSE 0 END) as occupied,
			SUM(CASE WHEN status = 'maintenance' THEN 1 ELSE 0 END) as maintenance,
			AVG(price_per_night) as price_per_night
		`).
		Group("room_type").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetRoomStats retrieves overall room statistics
func (r *roomRepository) GetRoomStats() (*models.RoomStats, error) {
	var stats models.RoomStats

	err := r.db.Model(&models.Room{}).
		Select(`
			COUNT(*) as total_rooms,
			SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END) as available,
			SUM(CASE WHEN status = 'occupied' THEN 1 ELSE 0 END) as occupied,
			SUM(CASE WHEN status = 'maintenance' THEN 1 ELSE 0 END) as maintenance,
			SUM(CASE WHEN status = 'cleaning' THEN 1 ELSE 0 END) as cleaning
		`).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
