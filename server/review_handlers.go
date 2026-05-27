package server

import (
	"net/http"

	"hotel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetServiceReviews handles getting reviews for a service
func (s *Server) GetServiceReviews(c *gin.Context) {
	serviceIDStr := c.Param("id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid service ID"})
		return
	}

	reviews, err := s.ReviewRepository.GetServiceReviews(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reviews,
		"message": "Service reviews retrieved",
	})
}

// GetProviderReviews handles getting reviews for a provider
func (s *Server) GetProviderReviews(c *gin.Context) {
	providerIDStr := c.Param("id")
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid provider ID"})
		return
	}

	reviews, err := s.ReviewRepository.GetProviderReviews(providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reviews,
		"message": "Provider reviews retrieved",
	})
}

// CreateServiceReview handles creating a review for a service
func (s *Server) CreateServiceReview(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer profile not found"})
		return
	}

	serviceIDStr := c.Param("id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid service ID"})
		return
	}

	// Get provider_id from service
	providerID, err := s.ProviderRepository.GetProviderIDByServiceID(serviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Service not found"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	req.ServiceID = &serviceID

	review, err := s.ReviewRepository.CreateReview(customer.ID, providerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    review,
		"message": "Review created successfully",
	})
}

// CreateProviderReview handles creating a review for a provider
func (s *Server) CreateProviderReview(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer profile not found"})
		return
	}

	providerIDStr := c.Param("id")
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid provider ID"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	review, err := s.ReviewRepository.CreateReview(customer.ID, providerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    review,
		"message": "Review created successfully",
	})
}

// UpdateReview handles updating a review
func (s *Server) UpdateReview(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	reviewIDStr := c.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid review ID"})
		return
	}

	// Verify the review belongs to the customer
	review, err := s.ReviewRepository.GetReviewByID(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Review not found"})
		return
	}

	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil || review.CustomerID != customer.ID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to update this review"})
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.ReviewRepository.UpdateReview(reviewID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review updated successfully",
	})
}

// DeleteReview handles deleting a review
func (s *Server) DeleteReview(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	reviewIDStr := c.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid review ID"})
		return
	}

	// Verify the review belongs to the customer
	review, err := s.ReviewRepository.GetReviewByID(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Review not found"})
		return
	}

	customer, err := s.CustomerRepository.GetCustomerByUserID(userID)
	if err != nil || review.CustomerID != customer.ID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to delete this review"})
		return
	}

	if err := s.ReviewRepository.DeleteReview(reviewID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review deleted successfully",
	})
}

// GetServiceRatingStats handles getting rating statistics for a service
func (s *Server) GetServiceRatingStats(c *gin.Context) {
	serviceIDStr := c.Param("id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid service ID"})
		return
	}

	stats, err := s.ReviewRepository.GetServiceRatingStats(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Service rating stats retrieved",
	})
}

// GetProviderRatingStats handles getting rating statistics for a provider
func (s *Server) GetProviderRatingStats(c *gin.Context) {
	providerIDStr := c.Param("id")
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid provider ID"})
		return
	}

	stats, err := s.ReviewRepository.GetProviderRatingStats(providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Provider rating stats retrieved",
	})
}
