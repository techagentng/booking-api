package server

import (
	"fmt"
	"net/http"

	"hotel/models"

	"github.com/gin-gonic/gin"
)

// RegisterProvider handles provider registration
func (s *Server) RegisterProvider(c *gin.Context) {
	var req models.ProviderRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	provider, user, err := s.ProviderRepository.RegisterProvider(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	response := models.ProviderRegistrationResponse{
		ProviderID:             provider.ID.String(),
		UserID:                 fmt.Sprintf("%d", user.ID),
		EmailVerificationToken: user.EmailVerificationToken,
		OnboardingStatus:       provider.OnboardingStatus,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Registration successful. Please verify your email.",
	})
}

// VerifyEmail handles email verification
func (s *Server) VerifyEmail(c *gin.Context) {
	var req models.EmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := s.ProviderRepository.VerifyEmail(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"verified": user.EmailVerified,
		},
		"message": "Email verified successfully",
	})
}

// UpdateBusinessInfo handles business information update
func (s *Server) UpdateBusinessInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider not found"})
		return
	}

	var req models.BusinessInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.ProviderRepository.UpdateBusinessInfo(provider.ID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"provider_id":        provider.ID.String(),
			"onboarding_status":  "business_info",
			"documents_uploaded": 0,
		},
		"message": "Business information saved successfully",
	})
}

// CreateProviderServices handles service creation during onboarding
func (s *Server) CreateProviderServices(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider not found"})
		return
	}

	var req models.ServiceCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.ProviderRepository.CreateProviderServices(provider.ID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"provider_id":       provider.ID.String(),
			"onboarding_status": "services",
			"services_created":  len(req.Services),
		},
		"message": "Services created successfully",
	})
}

// GetVerificationStatus handles verification status retrieval
func (s *Server) GetVerificationStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider not found"})
		return
	}

	verification, err := s.ProviderRepository.GetVerificationStatus(provider.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	steps := []models.VerificationStep{
		{
			ID:          "document_review",
			Name:        "Document Review",
			Status:      verification.DocumentReviewStatus,
			CompletedAt: nil,
		},
		{
			ID:          "business_verification",
			Name:        "Business Verification",
			Status:      verification.BusinessVerificationStatus,
			CompletedAt: nil,
		},
		{
			ID:          "identity_verification",
			Name:        "Identity Verification",
			Status:      verification.IdentityVerificationStatus,
			CompletedAt: nil,
		},
		{
			ID:          "manual_review",
			Name:        "Manual Review",
			Status:      verification.ManualReviewStatus,
			CompletedAt: verification.CompletedAt,
		},
	}

	response := models.VerificationStatusResponse{
		ProviderID:    provider.ID.String(),
		OverallStatus: verification.Status,
		Steps:         steps,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Verification status retrieved",
	})
}

// CompleteTraining handles training module completion
func (s *Server) CompleteTraining(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	moduleID := c.Param("module_id")
	if moduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Module ID required"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider not found"})
		return
	}

	if err := s.ProviderRepository.CompleteTraining(provider.ID, moduleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"module_id":    moduleID,
			"completed":    true,
			"completed_at": "now",
		},
		"message": "Training module completed",
	})
}

// ActivateProvider handles provider account activation
func (s *Server) ActivateProvider(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider not found"})
		return
	}

	if err := s.ProviderRepository.ActivateProvider(provider.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"provider_id":       provider.ID.String(),
			"onboarding_status": "active",
			"activated_at":      "now",
		},
		"message": "Account activated successfully",
	})
}

// GetOnboardingStatus handles onboarding status retrieval
func (s *Server) GetOnboardingStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"provider_id":         provider.ID.String(),
			"onboarding_status":   provider.OnboardingStatus,
			"verification_status": provider.VerificationStatus,
			"is_active":           provider.IsActive,
		},
		"message": "Onboarding status retrieved",
	})
}
