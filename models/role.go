package models

import (
	"github.com/google/uuid"
)

// Role constants
const (
	RoleAdmin     = "Admin"
	RoleProvider  = "Provider"
	RoleCustomer  = "Customer"
	RoleModerator = "Moderator"
	RoleUser      = "User" // Keep for backward compatibility
)

// Role represents a user role in the system
type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"unique;not null" json:"name"`
}
