package models

import "time"

// Room represents a hotel room
type Room struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RoomNumber    string    `json:"room_number" gorm:"uniqueIndex;not null"`
	RoomType      string    `json:"room_type" gorm:"not null;index"` // Standard, Deluxe, Suite
	Floor         int       `json:"floor" gorm:"index"`
	Capacity      int       `json:"capacity" gorm:"default:2"` // Maximum guests
	PricePerNight float64   `json:"price_per_night" gorm:"not null"`
	Status        string    `json:"status" gorm:"default:available;index"` // available, occupied, maintenance, cleaning
	Description   string    `json:"description"`
	Amenities     string    `json:"amenities"` // JSON string of amenities
	BedType       string    `json:"bed_type"`  // Single, Double, Queen, King, Twin
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Reservations []Reservation `json:"reservations,omitempty" gorm:"foreignKey:RoomID"`
}

// RoomAvailabilityRequest is the request for checking room availability (POST version)
type RoomAvailabilityRequest struct {
	RoomType     string `json:"room_type" binding:"required"`
	CheckInDate  string `json:"check_in_date" binding:"required"`
	CheckOutDate string `json:"check_out_date" binding:"required"`
}

// CreateRoomRequest is the request payload for creating a room
type CreateRoomRequest struct {
	RoomNumber    string  `json:"room_number" binding:"required"`
	RoomType      string  `json:"room_type" binding:"required"`
	Floor         int     `json:"floor" binding:"required,gt=0"`
	Capacity      int     `json:"capacity" binding:"required,gt=0"`
	PricePerNight float64 `json:"price_per_night" binding:"required,gt=0"`
	Description   string  `json:"description"`
	Amenities     string  `json:"amenities"`
	BedType       string  `json:"bed_type"`
}

// UpdateRoomRequest is the request payload for updating a room
type UpdateRoomRequest struct {
	RoomNumber    *string  `json:"room_number"`
	RoomType      *string  `json:"room_type"`
	Floor         *int     `json:"floor"`
	Capacity      *int     `json:"capacity"`
	PricePerNight *float64 `json:"price_per_night"`
	Status        *string  `json:"status"`
	Description   *string  `json:"description"`
	Amenities     *string  `json:"amenities"`
	BedType       *string  `json:"bed_type"`
}

// RoomListResponse is the paginated response for room list
type RoomListResponse struct {
	Rooms []Room         `json:"data"`
	Meta  PaginationMeta `json:"meta"`
}

// RoomTypeStats represents statistics for a room type
type RoomTypeStats struct {
	RoomType      string  `json:"room_type"`
	TotalRooms    int     `json:"total_rooms"`
	Available     int     `json:"available"`
	Occupied      int     `json:"occupied"`
	Maintenance   int     `json:"maintenance"`
	PricePerNight float64 `json:"price_per_night"`
}

// RoomStats represents overall room statistics
type RoomStats struct {
	TotalRooms  int `json:"total_rooms"`
	Available   int `json:"available"`
	Occupied    int `json:"occupied"`
	Maintenance int `json:"maintenance"`
	Cleaning    int `json:"cleaning"`
}
