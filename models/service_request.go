package models

import "time"

// ServiceRequest represents a guest service request
type ServiceRequest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	GuestID       uint      `json:"guest_id" gorm:"not null;index"`
	ReservationID uint      `json:"reservation_id" gorm:"not null;index"`
	ServiceType   string    `json:"service_type"`
	Priority      string    `json:"priority"` // low, medium, high
	Description   string    `json:"description"`
	Status        string    `json:"status"` // pending, in_progress, completed, cancelled
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"index"`

	// Relations
	Guest *Guest `json:"-" gorm:"foreignKey:GuestID"`
}
