package services

import (
	"fmt"
	"hotel/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GuestInsightsService handles guest insight generation
type GuestInsightsService struct {
	db *gorm.DB
}

// NewGuestInsightsService creates a new guest insights service
func NewGuestInsightsService(db *gorm.DB) *GuestInsightsService {
	return &GuestInsightsService{db: db}
}

// AnalyzeMealPreference generates meal preference insights based on room service orders
func (s *GuestInsightsService) AnalyzeMealPreference(guestID uint) string {
	var insights []string

	// Get total room service order count
	var orderCount int64
	s.db.Model(&models.GuestServiceRequest{}).
		Where("guest_id = ?", guestID).
		Where("type = ? OR service_type = ?", "room_service", "room-service").
		Where("status = ?", "completed").
		Count(&orderCount)

	if orderCount == 0 {
		return "No meal order history available"
	}

	// Get most ordered items (by description/service_type)
	type OrderItem struct {
		Description string
		Count       int64
	}
	var topItems []OrderItem
	s.db.Model(&models.GuestServiceRequest{}).
		Select("COALESCE(NULLIF(description, ''), service_type) as description, COUNT(*) as count").
		Where("guest_id = ?", guestID).
		Where("type = ? OR service_type = ?", "room_service", "room-service").
		Where("status = ?", "completed").
		Group("COALESCE(NULLIF(description, ''), service_type)").
		Order("count DESC").
		Limit(2).
		Scan(&topItems)

	if len(topItems) > 0 {
		items := []string{}
		for _, item := range topItems {
			if item.Description != "" {
				items = append(items, item.Description)
			}
		}
		if len(items) > 0 {
			insights = append(insights, fmt.Sprintf("Prefers %s", strings.Join(items, ", ")))
		}
	}

	// Get average order hour to determine meal time preference
	var avgHour float64
	s.db.Model(&models.GuestServiceRequest{}).
		Select("AVG(EXTRACT(HOUR FROM requested_at))").
		Where("guest_id = ?", guestID).
		Where("type = ? OR service_type = ?", "room_service", "room-service").
		Where("status = ?", "completed").
		Scan(&avgHour)

	if avgHour > 0 {
		if avgHour < 11 {
			insights = append(insights, "Breakfast person")
		} else if avgHour < 15 {
			insights = append(insights, "Lunch person")
		} else {
			insights = append(insights, "Dinner person")
		}
	}

	// Classify frequency
	if orderCount > 10 {
		insights = append(insights, "Frequent room service user")
	} else if orderCount > 5 {
		insights = append(insights, "Occasional room service user")
	} else {
		insights = append(insights, "Rare room service user")
	}

	if len(insights) == 0 {
		return "Standard meal preferences"
	}

	return strings.Join(insights, " • ")
}

// AnalyzeRoomPreference generates room preference insights based on booking history
func (s *GuestInsightsService) AnalyzeRoomPreference(guestID uint) string {
	var insights []string

	// Get total booking count
	var bookingCount int64
	s.db.Model(&models.Reservation{}).
		Where("guest_id = ?", guestID).
		Count(&bookingCount)

	if bookingCount == 0 {
		return "No booking history available"
	}

	// Get most booked room type
	type RoomTypeStat struct {
		RoomType string
		Count    int64
	}
	var roomTypeStat RoomTypeStat
	s.db.Model(&models.Reservation{}).
		Select("rooms.room_type, COUNT(*) as count").
		Joins("JOIN rooms ON reservations.room_id = rooms.id").
		Where("reservations.guest_id = ?", guestID).
		Group("rooms.room_type").
		Order("count DESC").
		Limit(1).
		Scan(&roomTypeStat)

	if roomTypeStat.RoomType != "" {
		if roomTypeStat.Count > 1 {
			insights = append(insights, fmt.Sprintf("Prefers %s rooms", roomTypeStat.RoomType))
		} else {
			insights = append(insights, fmt.Sprintf("Booked %s room", roomTypeStat.RoomType))
		}
	}

	// Get average floor preference
	var avgFloor float64
	s.db.Model(&models.Reservation{}).
		Select("AVG(rooms.floor)").
		Joins("JOIN rooms ON reservations.room_id = rooms.id").
		Where("reservations.guest_id = ?", guestID).
		Scan(&avgFloor)

	if avgFloor > 0 {
		if avgFloor > 5 {
			insights = append(insights, "Prefers higher floors")
		} else if avgFloor < 3 {
			insights = append(insights, "Prefers lower floors")
		} else {
			insights = append(insights, "Mid-level floors")
		}
	}

	// Get most common bed type
	type BedTypeStat struct {
		BedType string
		Count   int64
	}
	var bedTypeStat BedTypeStat
	s.db.Model(&models.Reservation{}).
		Select("rooms.bed_type, COUNT(*) as count").
		Joins("JOIN rooms ON reservations.room_id = rooms.id").
		Where("reservations.guest_id = ?", guestID).
		Where("rooms.bed_type != ''").
		Group("rooms.bed_type").
		Order("count DESC").
		Limit(1).
		Scan(&bedTypeStat)

	if bedTypeStat.BedType != "" && bedTypeStat.Count > 1 {
		insights = append(insights, fmt.Sprintf("%s bed preference", bedTypeStat.BedType))
	}

	// Loyalty indicator
	if bookingCount > 5 {
		insights = append(insights, fmt.Sprintf("Loyal guest (%d stays)", bookingCount))
	} else if bookingCount > 2 {
		insights = append(insights, fmt.Sprintf("Returning guest (%d stays)", bookingCount))
	} else {
		insights = append(insights, "New guest")
	}

	if len(insights) == 0 {
		return "Standard room preferences"
	}

	return strings.Join(insights, " • ")
}

// AnalyzeServicePattern generates service pattern insights
func (s *GuestInsightsService) AnalyzeServicePattern(guestID uint) string {
	var insights []string

	// Get total service request count
	var requestCount int64
	s.db.Model(&models.GuestServiceRequest{}).
		Where("guest_id = ?", guestID).
		Count(&requestCount)

	if requestCount == 0 {
		return "No service request history"
	}

	// Get most common service type
	type ServiceTypeStat struct {
		Type  string
		Count int64
	}
	var serviceTypeStat ServiceTypeStat
	s.db.Model(&models.GuestServiceRequest{}).
		Select("type, COUNT(*) as count").
		Where("guest_id = ?", guestID).
		Group("type").
		Order("count DESC").
		Limit(1).
		Scan(&serviceTypeStat)

	if serviceTypeStat.Type != "" {
		insights = append(insights, fmt.Sprintf("Most requested: %s", serviceTypeStat.Type))
	}

	// Classify usage frequency
	if requestCount > 20 {
		insights = append(insights, "High service usage")
	} else if requestCount > 10 {
		insights = append(insights, "Moderate service usage")
	} else if requestCount > 5 {
		insights = append(insights, "Low service usage")
	} else {
		insights = append(insights, "Minimal service usage")
	}

	// Check average response time expectation
	var avgResponseHours float64
	s.db.Model(&models.GuestServiceRequest{}).
		Select("AVG(EXTRACT(EPOCH FROM (completed_at - requested_at))/3600)").
		Where("guest_id = ?", guestID).
		Where("completed_at IS NOT NULL").
		Scan(&avgResponseHours)

	if avgResponseHours > 0 {
		if avgResponseHours < 1 {
			insights = append(insights, "Expects quick service")
		} else if avgResponseHours < 3 {
			insights = append(insights, "Standard service expectations")
		}
	}

	// Check for priority patterns
	var highPriorityCount int64
	s.db.Model(&models.GuestServiceRequest{}).
		Where("guest_id = ?", guestID).
		Where("priority = ?", "high").
		Count(&highPriorityCount)

	if highPriorityCount > 3 {
		insights = append(insights, "Often requests urgent service")
	}

	if len(insights) == 0 {
		return "Standard service pattern"
	}

	return strings.Join(insights, " • ")
}

// GenerateInsights creates or updates all insights for a guest
func (s *GuestInsightsService) GenerateInsights(guestID uint) (*models.GuestAIInsights, error) {
	mealPref := s.AnalyzeMealPreference(guestID)
	roomPref := s.AnalyzeRoomPreference(guestID)
	servicePat := s.AnalyzeServicePattern(guestID)

	// Calculate risk score based on complaints and service issues
	riskScore := s.calculateRiskScore(guestID)

	// Generate recommendations
	recommendations := s.generateRecommendations(guestID, mealPref, roomPref, servicePat)

	// Get complaints
	complaints := s.getComplaints(guestID)

	// Check if insights already exist
	var existingInsights models.GuestAIInsights
	err := s.db.Where("guest_id = ?", guestID).First(&existingInsights).Error

	if err == gorm.ErrRecordNotFound {
		// Create new insights
		newInsights := &models.GuestAIInsights{
			GuestID:         guestID,
			MealPreference:  mealPref,
			RoomPreference:  roomPref,
			ServicePattern:  servicePat,
			RiskScore:       riskScore,
			Recommendations: recommendations,
			Complaints:      complaints,
		}
		if err := s.db.Create(newInsights).Error; err != nil {
			return nil, err
		}
		return newInsights, nil
	} else if err != nil {
		return nil, err
	}

	// Update existing insights
	existingInsights.MealPreference = mealPref
	existingInsights.RoomPreference = roomPref
	existingInsights.ServicePattern = servicePat
	existingInsights.RiskScore = riskScore
	existingInsights.Recommendations = recommendations
	existingInsights.Complaints = complaints
	existingInsights.UpdatedAt = time.Now()

	if err := s.db.Save(&existingInsights).Error; err != nil {
		return nil, err
	}

	return &existingInsights, nil
}

// calculateRiskScore determines guest risk level
func (s *GuestInsightsService) calculateRiskScore(guestID uint) string {
	// Count cancelled reservations
	var cancelledCount int64
	s.db.Model(&models.Reservation{}).
		Where("guest_id = ?", guestID).
		Where("status = ?", "cancelled").
		Count(&cancelledCount)

	// Count total reservations
	var totalReservations int64
	s.db.Model(&models.Reservation{}).
		Where("guest_id = ?", guestID).
		Count(&totalReservations)

	// Count cancelled service requests
	var cancelledRequests int64
	s.db.Model(&models.GuestServiceRequest{}).
		Where("guest_id = ?", guestID).
		Where("status = ?", "cancelled").
		Count(&cancelledRequests)

	// Calculate risk
	if totalReservations > 0 {
		cancellationRate := float64(cancelledCount) / float64(totalReservations)
		if cancellationRate > 0.5 {
			return "high"
		} else if cancellationRate > 0.25 || cancelledRequests > 5 {
			return "medium"
		}
	}

	return "low"
}

// generateRecommendations creates personalized recommendations
func (s *GuestInsightsService) generateRecommendations(guestID uint, mealPref, roomPref, servicePat string) []string {
	var recommendations []string

	// Based on room preference
	if strings.Contains(roomPref, "Suite") {
		recommendations = append(recommendations, "Offer suite upgrade on next visit")
	}
	if strings.Contains(roomPref, "higher floors") {
		recommendations = append(recommendations, "Prioritize high floor rooms")
	}

	// Based on meal preference
	if strings.Contains(mealPref, "Frequent room service") {
		recommendations = append(recommendations, "Send room service menu upon check-in")
	}
	if strings.Contains(mealPref, "Dinner person") {
		recommendations = append(recommendations, "Recommend dinner specials")
	}

	// Based on service pattern
	if strings.Contains(servicePat, "quick service") {
		recommendations = append(recommendations, "Prioritize service requests")
	}

	// Based on loyalty
	if strings.Contains(roomPref, "Loyal guest") {
		recommendations = append(recommendations, "Consider loyalty rewards")
		recommendations = append(recommendations, "Offer complimentary amenities")
	}

	return recommendations
}

// getComplaints retrieves guest complaints from service requests
func (s *GuestInsightsService) getComplaints(guestID uint) []string {
	var complaints []string

	// Get cancelled service requests with reasons
	var cancelledRequests []models.GuestServiceRequest
	s.db.Where("guest_id = ?", guestID).
		Where("status = ?", "cancelled").
		Where("cancellation_reason != ''").
		Find(&cancelledRequests)

	for _, req := range cancelledRequests {
		if req.CancellationReason != "" {
			complaints = append(complaints, req.CancellationReason)
		}
	}

	// Limit to last 5 complaints
	if len(complaints) > 5 {
		complaints = complaints[len(complaints)-5:]
	}

	return complaints
}

// ShouldRefreshInsights checks if insights need to be refreshed (older than 7 days)
func (s *GuestInsightsService) ShouldRefreshInsights(guestID uint) bool {
	var insights models.GuestAIInsights
	err := s.db.Where("guest_id = ?", guestID).First(&insights).Error

	if err != nil {
		return true // No insights exist, should generate
	}

	// Check if older than 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	return insights.UpdatedAt.Before(sevenDaysAgo)
}
