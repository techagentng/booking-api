package server

import (
	"encoding/json"
	"net/http"

	"hotel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCustomerProfile handles getting current customer profile
func (s *Server) GetCustomerProfile(c *gin.Context) {
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

	response := models.CustomerProfileResponse{
		ID:               customer.ID,
		UserID:           customer.UserID,
		FullName:         customer.FullName,
		Email:            customer.Email,
		Phone:            customer.Phone,
		AvatarURL:        customer.AvatarURL,
		PreferredCity:    customer.PreferredCity,
		EmailVerified:    customer.EmailVerified,
		PhoneVerified:    customer.PhoneVerified,
		IdentityVerified: customer.IdentityVerified,
		Status:           customer.Status,
		CreatedAt:        customer.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Customer profile retrieved",
	})
}

// UpdateCustomerProfile handles updating customer profile
func (s *Server) UpdateCustomerProfile(c *gin.Context) {
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

	var req models.CustomerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.CustomerRepository.UpdateCustomerProfile(customer.ID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer profile updated",
	})
}

// GetSavedServices handles getting customer's saved services
func (s *Server) GetSavedServices(c *gin.Context) {
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

	savedServices, err := s.CustomerRepository.GetSavedServices(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    savedServices,
		"message": "Saved services retrieved",
	})
}

// SaveService handles saving a service
func (s *Server) SaveService(c *gin.Context) {
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

	var req models.SaveServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	savedService, err := s.CustomerRepository.SaveService(customer.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	response := models.SaveServiceResponse{
		ID:          savedService.ID,
		CustomerID:  savedService.CustomerID,
		ServiceID:   savedService.ServiceID,
		ServiceType: savedService.ServiceType,
		ServiceName: savedService.ServiceName,
		ImageURL:    savedService.ImageURL,
		Location:    savedService.Location,
		Price:       savedService.Price,
		CreatedAt:   savedService.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Service saved successfully",
	})
}

// RemoveSavedService handles removing a saved service
func (s *Server) RemoveSavedService(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	savedIDStr := c.Param("id")
	if savedIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Saved service ID required"})
		return
	}

	savedID, err := uuid.Parse(savedIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid saved service ID"})
		return
	}

	if err := s.CustomerRepository.RemoveSavedService(savedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Saved service removed",
	})
}

// GetCustomerBookings handles getting customer bookings
func (s *Server) GetCustomerBookings(c *gin.Context) {
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

	status := c.Query("status")
	bookings, err := s.CustomerRepository.GetCustomerBookings(customer.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
		"message": "Customer bookings retrieved",
	})
}

// GetCustomerPreferences handles getting customer preferences
func (s *Server) GetCustomerPreferences(c *gin.Context) {
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

	preferences, err := s.CustomerRepository.GetCustomerPreferences(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Parse JSON arrays
	var categories []string
	var locations []string
	json.Unmarshal([]byte(preferences.PreferredCategories), &categories)
	json.Unmarshal([]byte(preferences.PreferredLocations), &locations)

	response := models.CustomerPreferencesResponse{
		ID:                  preferences.ID,
		CustomerID:          preferences.CustomerID,
		PreferredCategories: categories,
		PreferredLocations:  locations,
		NotificationEmail:   preferences.NotificationEmail,
		NotificationSMS:     preferences.NotificationSMS,
		NotificationPush:    preferences.NotificationPush,
		CreatedAt:           preferences.CreatedAt,
		UpdatedAt:           preferences.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Customer preferences retrieved",
	})
}

// UpdateCustomerPreferences handles updating customer preferences
func (s *Server) UpdateCustomerPreferences(c *gin.Context) {
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

	var req models.CustomerPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.CustomerRepository.UpdateCustomerPreferences(customer.ID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer preferences updated",
	})
}
