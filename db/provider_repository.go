package db

import (
	"encoding/json"
	"fmt"
	"time"

	"hotel/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProviderRepository struct {
	DB *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{DB: db}
}

// RegisterProvider creates a new user and provider profile
func (r *ProviderRepository) RegisterProvider(req *models.ProviderRegistrationRequest) (*models.ServiceProvider, *models.User, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get Provider role
	var providerRole models.Role
	if err := tx.Where("name = ?", models.RoleProvider).First(&providerRole).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("provider role not found: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate email verification token
	emailVerificationToken := uuid.New().String()

	// Create user
	user := models.User{
		Fullname:               fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		Username:               req.Email, // Use email as username for simplicity
		Telephone:              req.Phone,
		Email:                  req.Email,
		HashedPassword:         string(hashedPassword),
		RoleID:                 providerRole.ID,
		EmailVerified:          false,
		EmailVerificationToken: emailVerificationToken,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create provider profile
	provider := models.ServiceProvider{
		ID:                 uuid.New(),
		UserID:             user.ID,
		BusinessName:       req.BusinessName,
		BusinessType:       req.BusinessType,
		CategoryID:         req.CategoryID,
		OnboardingStatus:   "registered",
		VerificationStatus: "pending",
		IsActive:           false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := tx.Create(&provider).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create provider profile: %w", err)
	}

	// Create verification record
	verification := models.ProviderVerification{
		ID:                         uuid.New(),
		ProviderID:                 provider.ID,
		Status:                     "pending",
		DocumentReviewStatus:       "pending",
		BusinessVerificationStatus: "pending",
		IdentityVerificationStatus: "pending",
		ManualReviewStatus:         "pending",
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	if err := tx.Create(&verification).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create verification record: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("transaction failed: %w", err)
	}

	return &provider, &user, nil
}

// VerifyEmail verifies the email using the token
func (r *ProviderRepository) VerifyEmail(token string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email_verification_token = ?", token).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	// Update user as verified
	now := time.Now()
	user.EmailVerified = true
	user.EmailVerificationToken = ""
	user.EmailVerifiedAt = &now

	if err := r.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// UpdateBusinessInfo updates provider business information
func (r *ProviderRepository) UpdateBusinessInfo(providerID uuid.UUID, req *models.BusinessInfoRequest) error {
	provider := models.ServiceProvider{ID: providerID}
	updates := map[string]interface{}{
		"address":           req.Address,
		"city":              req.City,
		"state":             req.State,
		"country":           req.Country,
		"postal_code":       req.PostalCode,
		"business_phone":    req.BusinessPhone,
		"business_email":    req.BusinessEmail,
		"website":           req.Website,
		"description":       req.Description,
		"onboarding_status": "business_info",
		"updated_at":        time.Now(),
	}

	if err := r.DB.Model(&provider).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update business info: %w", err)
	}

	return nil
}

// SaveDocument saves a provider document
func (r *ProviderRepository) SaveDocument(providerID uuid.UUID, docType, fileName, filePath string, fileSize int, mimeType string) (*models.ProviderDocument, error) {
	document := models.ProviderDocument{
		ID:           uuid.New(),
		ProviderID:   providerID,
		DocumentType: docType,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     fileSize,
		MimeType:     mimeType,
		UploadStatus: "uploaded",
		UploadedAt:   time.Now(),
	}

	if err := r.DB.Create(&document).Error; err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	return &document, nil
}

// CreateProviderServices creates services for a provider
func (r *ProviderRepository) CreateProviderServices(providerID uuid.UUID, req *models.ServiceCreationRequest) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update provider onboarding status
	if err := tx.Model(&models.ServiceProvider{ID: providerID}).Update("onboarding_status", "services").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update onboarding status: %w", err)
	}

	// Create services
	for _, serviceData := range req.Services {
		service := models.ProviderService{
			ID:          uuid.New(),
			ProviderID:  providerID,
			Title:       serviceData.Title,
			Description: serviceData.Description,
			CategoryID:  serviceData.CategoryID,
			PriceType:   serviceData.PriceType,
			BasePrice:   serviceData.BasePrice,
			Duration:    serviceData.Duration,
			IsActive:    true,
			IsAvailable: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Create(&service).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create service: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

// GetVerificationStatus returns verification status for a provider
func (r *ProviderRepository) GetVerificationStatus(providerID uuid.UUID) (*models.ProviderVerification, error) {
	var verification models.ProviderVerification
	if err := r.DB.Where("provider_id = ?", providerID).First(&verification).Error; err != nil {
		return nil, fmt.Errorf("verification record not found: %w", err)
	}

	return &verification, nil
}

// CompleteTraining marks a training module as completed
func (r *ProviderRepository) CompleteTraining(providerID uuid.UUID, moduleID string) error {
	now := time.Now()

	// Check if training progress already exists
	var progress models.ProviderTrainingProgress
	err := r.DB.Where("provider_id = ? AND module_id = ?", providerID, moduleID).First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		// Create new progress record
		progress = models.ProviderTrainingProgress{
			ID:          uuid.New(),
			ProviderID:  providerID,
			ModuleID:    moduleID,
			Completed:   true,
			CompletedAt: &now,
			CreatedAt:   time.Now(),
		}
		if err := r.DB.Create(&progress).Error; err != nil {
			return fmt.Errorf("failed to create training progress: %w", err)
		}
	} else if err == nil {
		// Update existing record
		progress.Completed = true
		progress.CompletedAt = &now
		if err := r.DB.Save(&progress).Error; err != nil {
			return fmt.Errorf("failed to update training progress: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check training progress: %w", err)
	}

	return nil
}

// ActivateProvider activates a provider account
func (r *ProviderRepository) ActivateProvider(providerID uuid.UUID) error {
	now := time.Now()

	updates := map[string]interface{}{
		"is_active":           true,
		"onboarding_status":   "active",
		"verification_status": "completed",
		"updated_at":          now,
	}

	if err := r.DB.Model(&models.ServiceProvider{ID: providerID}).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to activate provider: %w", err)
	}

	// Update verification record
	verificationUpdates := map[string]interface{}{
		"status":       "completed",
		"completed_at": now,
		"updated_at":   now,
	}

	if err := r.DB.Model(&models.ProviderVerification{}).Where("provider_id = ?", providerID).Updates(verificationUpdates).Error; err != nil {
		return fmt.Errorf("failed to update verification status: %w", err)
	}

	return nil
}

// GetProviderByUserID gets provider profile by user ID
func (r *ProviderRepository) GetProviderByUserID(userID uint) (*models.ServiceProvider, error) {
	var provider models.ServiceProvider
	if err := r.DB.Where("user_id = ?", userID).First(&provider).Error; err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	return &provider, nil
}

// GetProviderByID gets provider profile by ID
func (r *ProviderRepository) GetProviderByID(providerID uuid.UUID) (*models.ServiceProvider, error) {
	var provider models.ServiceProvider
	if err := r.DB.Where("id = ?", providerID).First(&provider).Error; err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	return &provider, nil
}

// CreateProviderService creates a new service for a provider
func (r *ProviderRepository) CreateProviderService(providerID uuid.UUID, req *models.CreateServiceRequest) (*models.ProviderService, error) {
	featuresJSON, _ := json.Marshal(req.Features)
	imagesJSON, _ := json.Marshal(req.Images)

	service := models.ProviderService{
		ID:          uuid.New(),
		ProviderID:  providerID,
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		PriceType:   req.PriceType,
		BasePrice:   req.BasePrice,
		Duration:    req.Duration,
		Features:    string(featuresJSON),
		Images:      string(imagesJSON),
		IsActive:    true,
		IsAvailable: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.DB.Create(&service).Error; err != nil {
		return nil, fmt.Errorf("failed to create provider service: %w", err)
	}

	return &service, nil
}

// GetProviderServices gets all services for a provider
func (r *ProviderRepository) GetProviderServices(providerID uuid.UUID) ([]models.ProviderService, error) {
	var services []models.ProviderService
	if err := r.DB.Where("provider_id = ?", providerID).Order("created_at DESC").Find(&services).Error; err != nil {
		return nil, fmt.Errorf("failed to get provider services: %w", err)
	}
	return services, nil
}

// GetProviderServiceByID gets a specific service by ID
func (r *ProviderRepository) GetProviderServiceByID(serviceID uuid.UUID) (*models.ProviderService, error) {
	var service models.ProviderService
	if err := r.DB.Where("id = ?", serviceID).First(&service).Error; err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}
	return &service, nil
}

// GetProviderIDByServiceID gets the provider ID for a given service ID
func (r *ProviderRepository) GetProviderIDByServiceID(serviceID uuid.UUID) (uuid.UUID, error) {
	var service models.ProviderService
	if err := r.DB.Where("id = ?", serviceID).Select("provider_id").First(&service).Error; err != nil {
		return uuid.Nil, fmt.Errorf("service not found: %w", err)
	}
	return service.ProviderID, nil
}

// UpdateProviderService updates a provider service
func (r *ProviderRepository) UpdateProviderService(serviceID uuid.UUID, req *models.UpdateServiceRequest) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CategoryID != "" {
		updates["category_id"] = req.CategoryID
	}
	if req.PriceType != "" {
		updates["price_type"] = req.PriceType
	}
	if req.BasePrice > 0 {
		updates["base_price"] = req.BasePrice
	}
	if req.Duration > 0 {
		updates["duration"] = req.Duration
	}
	if req.Features != nil {
		featuresJSON, _ := json.Marshal(req.Features)
		updates["features"] = string(featuresJSON)
	}
	if req.Images != nil {
		imagesJSON, _ := json.Marshal(req.Images)
		updates["images"] = string(imagesJSON)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	if err := r.DB.Model(&models.ProviderService{}).Where("id = ?", serviceID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update provider service: %w", err)
	}

	return nil
}

// DeleteProviderService deletes a provider service
func (r *ProviderRepository) DeleteProviderService(serviceID uuid.UUID) error {
	if err := r.DB.Delete(&models.ProviderService{ID: serviceID}).Error; err != nil {
		return fmt.Errorf("failed to delete provider service: %w", err)
	}
	return nil
}

// ToggleServiceAvailability toggles service availability
func (r *ProviderRepository) ToggleServiceAvailability(serviceID uuid.UUID, isAvailable bool) error {
	if err := r.DB.Model(&models.ProviderService{}).Where("id = ?", serviceID).Update("is_available", isAvailable).Error; err != nil {
		return fmt.Errorf("failed to toggle service availability: %w", err)
	}
	return nil
}
