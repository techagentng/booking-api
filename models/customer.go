package models

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer profile
type Customer struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            uint      `gorm:"not null" json:"user_id"`
	User              User      `gorm:"foreignKey:UserID" json:"-"`
	FullName          string    `gorm:"not null" json:"full_name"`
	Email             string    `gorm:"not null;unique" json:"email"`
	Phone             string    `json:"phone"`
	AvatarURL         string    `gorm:"type:text" json:"avatar_url"`
	PreferredCity     string    `json:"preferred_city"`
	EmailVerified     bool      `gorm:"default:false" json:"email_verified"`
	PhoneVerified     bool      `gorm:"default:false" json:"phone_verified"`
	IdentityVerified  bool      `gorm:"default:false" json:"identity_verified"`
	Status            string    `gorm:"default:'active'" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CustomerSavedServices represents services saved by a customer
type CustomerSavedServices struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID  uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	Customer    Customer  `gorm:"foreignKey:CustomerID" json:"-"`
	ServiceID   uuid.UUID  `json:"service_id"`
	ProviderID  uuid.UUID  `json:"provider_id"`
	ServiceType string    `gorm:"not null" json:"service_type"` // hotel, restaurant, transport, experience
	ServiceName string    `gorm:"not null" json:"service_name"`
	ImageURL    string    `gorm:"type:text" json:"image_url"`
	Location    string    `json:"location"`
	Price       float64   `json:"price"`
	Metadata    string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
}

// CustomerPreferences represents customer preferences
type CustomerPreferences struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID            uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	Customer              Customer  `gorm:"foreignKey:CustomerID" json:"-"`
	PreferredCategories   string     `gorm:"type:text" json:"preferred_categories"` // JSON array as string
	PreferredLocations    string     `gorm:"type:text" json:"preferred_locations"`   // JSON array as string
	NotificationEmail     bool       `gorm:"default:true" json:"notification_email"`
	NotificationSMS       bool       `gorm:"default:true" json:"notification_sms"`
	NotificationPush      bool       `gorm:"default:true" json:"notification_push"`
	Metadata              string     `gorm:"type:jsonb" json:"metadata"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// CustomerProfileRequest represents customer profile update request
type CustomerProfileRequest struct {
	FullName      string `json:"full_name" binding:"required"`
	Phone         string `json:"phone"`
	AvatarURL     string `json:"avatar_url"`
	PreferredCity string `json:"preferred_city"`
}

// CustomerProfileResponse represents customer profile response
type CustomerProfileResponse struct {
	ID               uuid.UUID `json:"id"`
	UserID           uint      `json:"user_id"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	AvatarURL        string    `json:"avatar_url"`
	PreferredCity    string    `json:"preferred_city"`
	EmailVerified    bool      `json:"email_verified"`
	PhoneVerified    bool      `json:"phone_verified"`
	IdentityVerified bool      `json:"identity_verified"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// SaveServiceRequest represents a service save request
type SaveServiceRequest struct {
	ServiceID   uuid.UUID `json:"service_id" binding:"required"`
	ProviderID  uuid.UUID `json:"provider_id" binding:"required"`
	ServiceType string    `json:"service_type" binding:"required"`
	ServiceName string    `json:"service_name" binding:"required"`
	ImageURL    string    `json:"image_url"`
	Location    string    `json:"location"`
	Price       float64   `json:"price"`
}

// SaveServiceResponse represents saved service response
type SaveServiceResponse struct {
	ID          uuid.UUID `json:"id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	ServiceID   uuid.UUID `json:"service_id"`
	ServiceType string    `json:"service_type"`
	ServiceName string    `json:"service_name"`
	ImageURL    string    `json:"image_url"`
	Location    string    `json:"location"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// CustomerPreferencesRequest represents customer preferences update request
type CustomerPreferencesRequest struct {
	PreferredCategories []string `json:"preferred_categories"`
	PreferredLocations  []string `json:"preferred_locations"`
	NotificationEmail   bool     `json:"notification_email"`
	NotificationSMS     bool     `json:"notification_sms"`
	NotificationPush    bool     `json:"notification_push"`
}

// CustomerPreferencesResponse represents customer preferences response
type CustomerPreferencesResponse struct {
	ID                  uuid.UUID `json:"id"`
	CustomerID          uuid.UUID `json:"customer_id"`
	PreferredCategories []string `json:"preferred_categories"`
	PreferredLocations  []string `json:"preferred_locations"`
	NotificationEmail   bool     `json:"notification_email"`
	NotificationSMS     bool     `json:"notification_sms"`
	NotificationPush    bool     `json:"notification_push"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
