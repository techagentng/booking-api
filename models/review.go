package models

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a review for a service or provider
type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;not null;index" json:"customer_id"`
	ProviderID uuid.UUID `gorm:"type:uuid;not null;index" json:"provider_id"`
	ServiceID  *uuid.UUID `gorm:"type:uuid;index" json:"service_id,omitempty"` // Nullable for provider-only reviews
	Rating     int       `gorm:"not null;index" json:"rating"` // 1-5 stars
	Title      string    `json:"title"`
	Comment    string    `gorm:"type:text" json:"comment"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Provider *ServiceProvider `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	Service  *ProviderService `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

// CreateReviewRequest represents a request to create a review
type CreateReviewRequest struct {
	ServiceID *uuid.UUID `json:"service_id,omitempty"`
	Rating    int        `json:"rating" binding:"required,min=1,max=5"`
	Title     string     `json:"title"`
	Comment   string     `json:"comment"`
}

// UpdateReviewRequest represents a request to update a review
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Title   *string `json:"title"`
	Comment *string `json:"comment"`
}

// ReviewResponse represents a review response
type ReviewResponse struct {
	ID         uuid.UUID  `json:"id"`
	CustomerID uuid.UUID  `json:"customer_id"`
	ProviderID uuid.UUID  `json:"provider_id"`
	ServiceID  *uuid.UUID `json:"service_id,omitempty"`
	Rating     int        `json:"rating"`
	Title      string     `json:"title"`
	Comment    string     `json:"comment"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CustomerName string `json:"customer_name,omitempty"`
}

// ProviderRatingStats represents rating statistics for a provider
type ProviderRatingStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalReviews  int     `json:"total_reviews"`
	RatingDistribution map[string]int `json:"rating_distribution"` // {"5": 10, "4": 5, ...}
}

// ServiceRatingStats represents rating statistics for a service
type ServiceRatingStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalReviews  int     `json:"total_reviews"`
	RatingDistribution map[string]int `json:"rating_distribution"`
}
