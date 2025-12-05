package models

import (
	"time"

	"github.com/lib/pq"
)

// Guest represents a hotel guest
type Guest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name" gorm:"not null;index"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Phone       string    `json:"phone" gorm:"not null"`
	Nationality string    `json:"nationality"`
	IDType      string    `json:"id_type"`
	IDNumber    string    `json:"id_number" gorm:"uniqueIndex:,composite:guest_id"`
	JoinDate    time.Time `json:"join_date" gorm:"autoCreateTime"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`

	// Relations
	Reservations    []Reservation     `json:"reservations,omitempty" gorm:"foreignKey:GuestID"`
	ServiceRequests []ServiceRequest  `json:"service_requests,omitempty" gorm:"foreignKey:GuestID"`
	Preferences     *GuestPreferences `json:"preferences,omitempty" gorm:"foreignKey:GuestID"`
	AIInsights      *GuestAIInsights  `json:"ai_insights,omitempty" gorm:"foreignKey:GuestID"`
}

// GuestPreferences stores guest preferences
type GuestPreferences struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	GuestID         uint           `gorm:"not null;uniqueIndex" json:"guest_id"`
	RoomFloors      pq.StringArray `gorm:"type:text[]" json:"room_floors"`
	MealTypes       pq.StringArray `gorm:"type:text[]" json:"meal_types"`
	RoomTypes       pq.StringArray `gorm:"type:text[]" json:"room_types"`
	SpecialRequests pq.StringArray `gorm:"type:text[]" json:"special_requests"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       time.Time      `json:"deleted_at" gorm:"index"`

	// Relation
	Guest *Guest `json:"-" gorm:"foreignKey:GuestID"`
}

// GuestAIInsights stores AI-generated insights about guest
type GuestAIInsights struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	GuestID         uint           `gorm:"not null;uniqueIndex" json:"guest_id"`
	MealPreference  string         `json:"meal_preference"`
	RoomPreference  string         `json:"room_preference"`
	ServicePattern  string         `json:"service_pattern"`
	RiskScore       string         `json:"risk_score"` // low, medium, high
	Recommendations pq.StringArray `gorm:"type:text[]" json:"recommendations"`
	Complaints      pq.StringArray `gorm:"type:text[]" json:"complaints"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       time.Time      `json:"deleted_at" gorm:"index"`

	// Relation
	Guest *Guest `json:"-" gorm:"foreignKey:GuestID"`
}

// CreateGuestRequest is the request payload for creating a guest
type CreateGuestRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"required"`
	Nationality string `json:"nationality" binding:"required"`
	IDType      string `json:"id_type" binding:"required"`
	IDNumber    string `json:"id_number" binding:"required"`
}

// UpdateGuestRequest is the request payload for updating a guest
type UpdateGuestRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Email       *string `json:"email" binding:"omitempty,email"`
	Phone       *string `json:"phone" binding:"omitempty"`
	Nationality *string `json:"nationality"`
	IDType      *string `json:"id_type"`
	IDNumber    *string `json:"id_number"`
}

// GuestStatistics represents calculated guest statistics
type GuestStatistics struct {
	TotalStays     int       `json:"total_stays"`
	TotalSpent     float64   `json:"total_spent"`
	AverageSpend   float64   `json:"average_spend"`
	LastVisit      time.Time `json:"last_visit"`
	MostCommonRoom string    `json:"most_common_room"`
}

// ServiceUsageItem represents service usage statistics
type ServiceUsageItem struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// GuestHistory represents guest stay history
type GuestHistory struct {
	TotalStays     int           `json:"total_stays"`
	TotalSpent     float64       `json:"total_spent"`
	AverageSpend   float64       `json:"average_spend"`
	LastVisitDate  time.Time     `json:"last_visit_date"`
	MostCommonRoom string        `json:"most_common_room"`
	Reservations   []Reservation `json:"reservations"`
}

// GuestDetailsResponse is the enriched guest response with statistics
type GuestDetailsResponse struct {
	ID           uint               `json:"id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	Nationality  string             `json:"nationality"`
	IDType       string             `json:"id_type"`
	IDNumber     string             `json:"id_number"`
	JoinDate     time.Time          `json:"join_date"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Reservations []Reservation      `json:"reservations"`
	Preferences  *GuestPreferences  `json:"preferences"`
	AIInsights   *GuestAIInsights   `json:"ai_insights"`
	Statistics   GuestStatistics    `json:"statistics"`
	ServiceUsage []ServiceUsageItem `json:"service_usage"`
}

// GuestListResponse is the paginated response for guest list
type GuestListResponse struct {
	Guests []Guest        `json:"data"`
	Meta   PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
