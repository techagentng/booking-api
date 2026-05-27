package models

import (
	"time"

	"github.com/google/uuid"
)

// ProviderService represents a service offered by a provider
type ProviderService struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderID  uuid.UUID       `gorm:"type:uuid;not null" json:"provider_id"`
	Provider    ServiceProvider `gorm:"foreignKey:ProviderID" json:"-"`
	Title       string          `gorm:"not null" json:"title"`
	Description string          `gorm:"type:text" json:"description"`
	CategoryID  string          `gorm:"not null" json:"category_id"`
	PriceType   string          `gorm:"not null" json:"price_type"` // fixed, hourly, per_item, custom
	BasePrice   float64         `gorm:"not null" json:"base_price"`
	Duration    int             `gorm:"not null" json:"duration"`
	Features    string          `gorm:"type:text" json:"features"` // JSON array as string
	Images      string          `gorm:"type:text" json:"images"`   // JSON array as string
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	IsAvailable bool            `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// CreateServiceRequest represents a request to create a service
type CreateServiceRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	PriceType   string   `json:"price_type" binding:"required"`
	BasePrice   float64  `json:"base_price" binding:"required"`
	Duration    int      `json:"duration"`
	Features    []string `json:"features"`
	Images      []string `json:"images"`
}

// UpdateServiceRequest represents a request to update a service
type UpdateServiceRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	PriceType   string   `json:"price_type"`
	BasePrice   float64  `json:"base_price"`
	Duration    int      `json:"duration"`
	Features    []string `json:"features"`
	Images      []string `json:"images"`
	IsActive    *bool    `json:"is_active"`
	IsAvailable *bool    `json:"is_available"`
}

// ServiceProvider represents a service provider profile
type ServiceProvider struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"-"`
	BusinessName    string    `gorm:"not null" json:"business_name"`
	BusinessType    string    `gorm:"not null" json:"business_type"` // individual, company, franchise
	PrimaryCategory string    `json:"primary_category"`
	CategoryID      string    `json:"category_id"`
	Description     string    `gorm:"type:text" json:"description"`
	Address         string    `gorm:"type:text" json:"address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	Country         string    `gorm:"default:'Nigeria'" json:"country"`
	PostalCode      string    `json:"postal_code"`
	BusinessPhone   string    `json:"business_phone"`
	BusinessEmail   string    `json:"business_email"`
	Website         string    `json:"website"`

	// Onboarding Status
	OnboardingStatus   string `gorm:"default:'registered'" json:"onboarding_status"` // registered, business_info, services, verification, training, active
	VerificationStatus string `gorm:"default:'pending'" json:"verification_status"`  // pending, in_progress, completed, failed
	IsActive           bool   `gorm:"default:false" json:"is_active"`

	// Admin Positioning (from existing system)
	AdminPosition    int    `gorm:"default:0" json:"admin_position"`
	PositionCategory string `json:"position_category"`
	IsFeatured       bool   `gorm:"default:false" json:"is_featured"`

	// Ratings (from existing system)
	AverageRating float64 `gorm:"default:0" json:"average_rating"`
	TotalReviews  int     `gorm:"default:0" json:"total_reviews"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProviderVerification tracks verification status for providers
type ProviderVerification struct {
	ID                         uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderID                 uuid.UUID       `gorm:"type:uuid;not null" json:"provider_id"`
	Provider                   ServiceProvider `gorm:"foreignKey:ProviderID" json:"-"`
	Status                     string          `gorm:"default:'pending'" json:"status"` // pending, in_progress, completed, failed
	DocumentReviewStatus       string          `gorm:"default:'pending'" json:"document_review_status"`
	BusinessVerificationStatus string          `gorm:"default:'pending'" json:"business_verification_status"`
	IdentityVerificationStatus string          `gorm:"default:'pending'" json:"identity_verification_status"`
	ManualReviewStatus         string          `gorm:"default:'pending'" json:"manual_review_status"`
	CompletedAt                *time.Time      `json:"completed_at"`
	CreatedAt                  time.Time       `json:"created_at"`
	UpdatedAt                  time.Time       `json:"updated_at"`
}

// ProviderDocument stores uploaded documents
type ProviderDocument struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderID        uuid.UUID       `gorm:"type:uuid;not null" json:"provider_id"`
	Provider          ServiceProvider `gorm:"foreignKey:ProviderID" json:"-"`
	DocumentType      string          `gorm:"not null" json:"document_type"` // business_registration, tax_id, proof_of_address, identity
	FileName          string          `gorm:"not null" json:"file_name"`
	FilePath          string          `gorm:"not null" json:"file_path"`
	FileSize          int             `json:"file_size"`
	MimeType          string          `json:"mime_type"`
	UploadStatus      string          `gorm:"default:'pending'" json:"upload_status"` // pending, uploaded, verified, rejected
	VerificationNotes string          `gorm:"type:text" json:"verification_notes"`
	UploadedAt        time.Time       `json:"uploaded_at"`
	VerifiedAt        *time.Time      `json:"verified_at"`
}

// ProviderTrainingProgress tracks training completion
type ProviderTrainingProgress struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderID  uuid.UUID       `gorm:"type:uuid;not null" json:"provider_id"`
	Provider    ServiceProvider `gorm:"foreignKey:ProviderID" json:"-"`
	ModuleID    string          `gorm:"not null" json:"module_id"` // dashboard, bookings, services, communication
	Completed   bool            `gorm:"default:false" json:"completed"`
	CompletedAt *time.Time      `json:"completed_at"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ProviderRegistrationRequest represents provider registration data
type ProviderRegistrationRequest struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	BusinessName string `json:"business_name" binding:"required"`
	BusinessType string `json:"business_type" binding:"required"` // individual, company, franchise
	CategoryID   string `json:"category_id" binding:"required"`
}

// ProviderRegistrationResponse represents registration response
type ProviderRegistrationResponse struct {
	ProviderID             string `json:"provider_id"`
	UserID                 string `json:"user_id"`
	EmailVerificationToken string `json:"email_verification_token"`
	OnboardingStatus       string `json:"onboarding_status"`
}

// EmailVerificationRequest represents email verification request
type EmailVerificationRequest struct {
	Token string `json:"token" binding:"required"`
}

// BusinessInfoRequest represents business information update
type BusinessInfoRequest struct {
	Address       string `json:"address" binding:"required"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state" binding:"required"`
	Country       string `json:"country"`
	PostalCode    string `json:"postal_code"`
	BusinessPhone string `json:"business_phone" binding:"required"`
	BusinessEmail string `json:"business_email" binding:"required,email"`
	Website       string `json:"website"`
	Description   string `json:"description" binding:"required"`
}

// ServiceCreationRequest represents service creation during onboarding
type ServiceCreationRequest struct {
	Categories []string              `json:"categories" binding:"required"`
	Services   []ProviderServiceData `json:"services" binding:"required"`
}

// ProviderServiceData represents service data for creation
type ProviderServiceData struct {
	Title        string         `json:"title" binding:"required"`
	CategoryID   string         `json:"category_id" binding:"required"`
	Description  string         `json:"description" binding:"required"`
	PriceType    string         `json:"price_type" binding:"required"` // fixed, hourly, per_item, custom
	BasePrice    float64        `json:"base_price" binding:"required"`
	Duration     int            `json:"duration" binding:"required"`
	Features     []string       `json:"features"`
	Availability []TimeSlotData `json:"availability"`
}

// TimeSlotData represents time slot data
type TimeSlotData struct {
	DayOfWeek string `json:"day_of_week" binding:"required"` // monday, tuesday, etc.
	OpenTime  string `json:"open_time" binding:"required"`   // "09:00"
	CloseTime string `json:"close_time" binding:"required"`  // "17:00"
}

// TrainingCompletionRequest represents training module completion
type TrainingCompletionRequest struct {
	Completed bool `json:"completed" binding:"required"`
}

// VerificationStep represents a verification step status
type VerificationStep struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"` // pending, in_progress, completed
	CompletedAt *time.Time `json:"completed_at"`
}

// VerificationStatusResponse represents verification status response
type VerificationStatusResponse struct {
	ProviderID    string             `json:"provider_id"`
	OverallStatus string             `json:"overall_status"`
	Steps         []VerificationStep `json:"steps"`
}
