package db

import (
	"fmt"
	"time"

	"hotel/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

// CreateReview creates a new review
func (r *ReviewRepository) CreateReview(customerID, providerID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error) {
	review := models.Review{
		ID:         uuid.New(),
		CustomerID: customerID,
		ProviderID: providerID,
		ServiceID:  req.ServiceID,
		Rating:     req.Rating,
		Title:      req.Title,
		Comment:    req.Comment,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := r.DB.Create(&review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return &review, nil
}

// GetServiceReviews gets reviews for a specific service
func (r *ReviewRepository) GetServiceReviews(serviceID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.DB.Where("service_id = ?", serviceID).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("failed to get service reviews: %w", err)
	}
	return reviews, nil
}

// GetProviderReviews gets reviews for a provider
func (r *ReviewRepository) GetProviderReviews(providerID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.DB.Where("provider_id = ?", providerID).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("failed to get provider reviews: %w", err)
	}
	return reviews, nil
}

// GetCustomerReviews gets reviews by a customer
func (r *ReviewRepository) GetCustomerReviews(customerID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.DB.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer reviews: %w", err)
	}
	return reviews, nil
}

// GetReviewByID gets a review by ID
func (r *ReviewRepository) GetReviewByID(reviewID uuid.UUID) (*models.Review, error) {
	var review models.Review
	if err := r.DB.Where("id = ?", reviewID).First(&review).Error; err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}
	return &review, nil
}

// UpdateReview updates a review
func (r *ReviewRepository) UpdateReview(reviewID uuid.UUID, req *models.UpdateReviewRequest) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}

	if err := r.DB.Model(&models.Review{}).Where("id = ?", reviewID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

// DeleteReview deletes a review
func (r *ReviewRepository) DeleteReview(reviewID uuid.UUID) error {
	if err := r.DB.Delete(&models.Review{ID: reviewID}).Error; err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}

// GetProviderRatingStats gets rating statistics for a provider
func (r *ReviewRepository) GetProviderRatingStats(providerID uuid.UUID) (*models.ProviderRatingStats, error) {
	var stats models.ProviderRatingStats
	
	// Get total reviews and average rating
	var result struct {
		AverageRating float64
		TotalReviews  int64
	}
	
	if err := r.DB.Model(&models.Review{}).
		Select("AVG(rating) as average_rating, COUNT(*) as total_reviews").
		Where("provider_id = ?", providerID).
		Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating stats: %w", err)
	}
	
	stats.AverageRating = result.AverageRating
	stats.TotalReviews = int(result.TotalReviews)
	
	// Get rating distribution
	var distribution []struct {
		Rating int
		Count  int64
	}
	
	if err := r.DB.Model(&models.Review{}).
		Select("rating, COUNT(*) as count").
		Where("provider_id = ?", providerID).
		Group("rating").
		Order("rating DESC").
		Scan(&distribution).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}
	
	stats.RatingDistribution = make(map[string]int)
	for _, d := range distribution {
		stats.RatingDistribution[fmt.Sprintf("%d", d.Rating)] = int(d.Count)
	}
	
	return &stats, nil
}

// GetServiceRatingStats gets rating statistics for a service
func (r *ReviewRepository) GetServiceRatingStats(serviceID uuid.UUID) (*models.ServiceRatingStats, error) {
	var stats models.ServiceRatingStats
	
	var result struct {
		AverageRating float64
		TotalReviews  int64
	}
	
	if err := r.DB.Model(&models.Review{}).
		Select("AVG(rating) as average_rating, COUNT(*) as total_reviews").
		Where("service_id = ?", serviceID).
		Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating stats: %w", err)
	}
	
	stats.AverageRating = result.AverageRating
	stats.TotalReviews = int(result.TotalReviews)
	
	var distribution []struct {
		Rating int
		Count  int64
	}
	
	if err := r.DB.Model(&models.Review{}).
		Select("rating, COUNT(*) as count").
		Where("service_id = ?", serviceID).
		Group("rating").
		Order("rating DESC").
		Scan(&distribution).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}
	
	stats.RatingDistribution = make(map[string]int)
	for _, d := range distribution {
		stats.RatingDistribution[fmt.Sprintf("%d", d.Rating)] = int(d.Count)
	}
	
	return &stats, nil
}
