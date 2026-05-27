package db

import (
	"encoding/json"
	"fmt"
	"time"

	"hotel/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	DB *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

// GetOrCreateCustomer gets or creates a customer profile for a user
func (r *CustomerRepository) GetOrCreateCustomer(userID uint, email, fullname string) (*models.Customer, error) {
	var customer models.Customer
	err := r.DB.Where("user_id = ?", userID).First(&customer).Error

	if err == gorm.ErrRecordNotFound {
		// Create new customer
		customer = models.Customer{
			ID:               uuid.New(),
			UserID:           userID,
			FullName:         fullname,
			Email:            email,
			EmailVerified:    false,
			PhoneVerified:    false,
			IdentityVerified: false,
			Status:           "active",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := r.DB.Create(&customer).Error; err != nil {
			return nil, fmt.Errorf("failed to create customer: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

// GetCustomerByUserID gets customer by user ID
func (r *CustomerRepository) GetCustomerByUserID(userID uint) (*models.Customer, error) {
	var customer models.Customer
	if err := r.DB.Where("user_id = ?", userID).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return &customer, nil
}

// GetCustomerByID gets customer by ID
func (r *CustomerRepository) GetCustomerByID(customerID uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.DB.Where("id = ?", customerID).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return &customer, nil
}

// UpdateCustomerProfile updates customer profile
func (r *CustomerRepository) UpdateCustomerProfile(customerID uuid.UUID, req *models.CustomerProfileRequest) error {
	updates := map[string]interface{}{
		"full_name":      req.FullName,
		"phone":          req.Phone,
		"avatar_url":     req.AvatarURL,
		"preferred_city": req.PreferredCity,
		"updated_at":     time.Now(),
	}

	if err := r.DB.Model(&models.Customer{ID: customerID}).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update customer profile: %w", err)
	}

	return nil
}

// SaveService saves a service for a customer
func (r *CustomerRepository) SaveService(customerID uuid.UUID, req *models.SaveServiceRequest) (*models.CustomerSavedServices, error) {
	// Check if already saved
	var existing models.CustomerSavedServices
	err := r.DB.Where("customer_id = ? AND service_id = ?", customerID, req.ServiceID).First(&existing).Error

	if err == nil {
		// Already saved, return existing
		return &existing, nil
	}

	// Create new saved service
	metadata := map[string]interface{}{}
	metadataJSON, _ := json.Marshal(metadata)

	savedService := models.CustomerSavedServices{
		ID:          uuid.New(),
		CustomerID:  customerID,
		ServiceID:   req.ServiceID,
		ProviderID:  req.ProviderID,
		ServiceType: req.ServiceType,
		ServiceName: req.ServiceName,
		ImageURL:    req.ImageURL,
		Location:    req.Location,
		Price:       req.Price,
		Metadata:    string(metadataJSON),
		CreatedAt:   time.Now(),
	}

	if err := r.DB.Create(&savedService).Error; err != nil {
		return nil, fmt.Errorf("failed to save service: %w", err)
	}

	return &savedService, nil
}

// GetSavedServices gets all saved services for a customer
func (r *CustomerRepository) GetSavedServices(customerID uuid.UUID) ([]models.CustomerSavedServices, error) {
	var savedServices []models.CustomerSavedServices
	if err := r.DB.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&savedServices).Error; err != nil {
		return nil, fmt.Errorf("failed to get saved services: %w", err)
	}
	return savedServices, nil
}

// RemoveSavedService removes a saved service
func (r *CustomerRepository) RemoveSavedService(savedID uuid.UUID) error {
	if err := r.DB.Delete(&models.CustomerSavedServices{ID: savedID}).Error; err != nil {
		return fmt.Errorf("failed to remove saved service: %w", err)
	}
	return nil
}

// RemoveSavedServiceByServiceID removes a saved service by service ID
func (r *CustomerRepository) RemoveSavedServiceByServiceID(customerID, serviceID uuid.UUID) error {
	if err := r.DB.Where("customer_id = ? AND service_id = ?", customerID, serviceID).Delete(&models.CustomerSavedServices{}).Error; err != nil {
		return fmt.Errorf("failed to remove saved service: %w", err)
	}
	return nil
}

// GetCustomerBookings gets bookings for a customer
func (r *CustomerRepository) GetCustomerBookings(customerID uuid.UUID, status string) ([]models.Reservation, error) {
	var bookings []models.Reservation
	query := r.DB.Order("created_at DESC")

	// Filter by customer_id
	query = query.Where("customer_id = ?", customerID.String())

	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer bookings: %w", err)
	}

	return bookings, nil
}

// GetCustomerPreferences gets customer preferences
func (r *CustomerRepository) GetCustomerPreferences(customerID uuid.UUID) (*models.CustomerPreferences, error) {
	var preferences models.CustomerPreferences
	err := r.DB.Where("customer_id = ?", customerID).First(&preferences).Error

	if err == gorm.ErrRecordNotFound {
		// Create default preferences
		preferences = models.CustomerPreferences{
			ID:                  uuid.New(),
			CustomerID:          customerID,
			PreferredCategories: "[]",
			PreferredLocations:  "[]",
			NotificationEmail:   true,
			NotificationSMS:     true,
			NotificationPush:    true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}
		if err := r.DB.Create(&preferences).Error; err != nil {
			return nil, fmt.Errorf("failed to create customer preferences: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get customer preferences: %w", err)
	}

	return &preferences, nil
}

// UpdateCustomerPreferences updates customer preferences
func (r *CustomerRepository) UpdateCustomerPreferences(customerID uuid.UUID, req *models.CustomerPreferencesRequest) error {
	categoriesJSON, _ := json.Marshal(req.PreferredCategories)
	locationsJSON, _ := json.Marshal(req.PreferredLocations)

	updates := map[string]interface{}{
		"preferred_categories": string(categoriesJSON),
		"preferred_locations":  string(locationsJSON),
		"notification_email":   req.NotificationEmail,
		"notification_sms":     req.NotificationSMS,
		"notification_push":    req.NotificationPush,
		"updated_at":           time.Now(),
	}

	if err := r.DB.Model(&models.CustomerPreferences{}).Where("customer_id = ?", customerID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update customer preferences: %w", err)
	}

	return nil
}
