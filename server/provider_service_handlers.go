package server

import (
	"encoding/json"
	"net/http"

	"hotel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetProviderServices handles getting all services for a provider
func (s *Server) GetProviderServices(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider profile not found"})
		return
	}

	services, err := s.ProviderRepository.GetProviderServices(provider.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Parse JSON arrays for response
	var servicesResponse []map[string]interface{}
	for _, service := range services {
		var features []string
		var images []string
		json.Unmarshal([]byte(service.Features), &features)
		json.Unmarshal([]byte(service.Images), &images)

		servicesResponse = append(servicesResponse, map[string]interface{}{
			"id":           service.ID,
			"provider_id":  service.ProviderID,
			"title":        service.Title,
			"description":  service.Description,
			"category_id":  service.CategoryID,
			"price_type":   service.PriceType,
			"base_price":   service.BasePrice,
			"duration":     service.Duration,
			"features":     features,
			"images":       images,
			"is_active":    service.IsActive,
			"is_available": service.IsAvailable,
			"created_at":   service.CreatedAt,
			"updated_at":   service.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    servicesResponse,
		"message": "Provider services retrieved",
	})
}

// CreateProviderService handles creating a new service
func (s *Server) CreateProviderService(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	provider, err := s.ProviderRepository.GetProviderByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Provider profile not found"})
		return
	}

	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	service, err := s.ProviderRepository.CreateProviderService(provider.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Parse JSON arrays for response
	var features []string
	var images []string
	json.Unmarshal([]byte(service.Features), &features)
	json.Unmarshal([]byte(service.Images), &images)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"id":           service.ID,
			"provider_id":  service.ProviderID,
			"title":        service.Title,
			"description":  service.Description,
			"category_id":  service.CategoryID,
			"price_type":   service.PriceType,
			"base_price":   service.BasePrice,
			"duration":     service.Duration,
			"features":     features,
			"images":       images,
			"is_active":    service.IsActive,
			"is_available": service.IsAvailable,
			"created_at":   service.CreatedAt,
			"updated_at":   service.UpdatedAt,
		},
		"message": "Service created successfully",
	})
}

// UpdateProviderService handles updating a service
func (s *Server) UpdateProviderService(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	serviceIDStr := c.Param("id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid service ID"})
		return
	}

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.ProviderRepository.UpdateProviderService(serviceID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service updated successfully",
	})
}

// DeleteProviderService handles deleting a service
func (s *Server) DeleteProviderService(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	serviceIDStr := c.Param("id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid service ID"})
		return
	}

	if err := s.ProviderRepository.DeleteProviderService(serviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service deleted successfully",
	})
}

// ToggleServiceAvailability handles toggling service availability
func (s *Server) ToggleServiceAvailability(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	serviceIDStr := c.Param("id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid service ID"})
		return
	}

	var req struct {
		IsAvailable bool `json:"is_available" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := s.ProviderRepository.ToggleServiceAvailability(serviceID, req.IsAvailable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service availability updated",
	})
}
