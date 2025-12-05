package db

import (
	"errors"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// GuestRepository defines guest database operations
type GuestRepository interface {
	CreateGuest(guest *models.Guest) (*models.Guest, error)
	GetGuestByID(id uint) (*models.Guest, error)
	GetGuestDetailsByID(id uint) (*models.GuestDetailsResponse, error)
	GetAllGuests(page, pageSize int, search string) ([]models.Guest, int64, error)
	UpdateGuest(id uint, guest *models.Guest) (*models.Guest, error)
	DeleteGuest(id uint) error
	GetGuestByEmail(email string) (*models.Guest, error)
	GetGuestHistory(guestID uint) (*models.GuestHistory, error)
	GetGuestPreferences(guestID uint) (*models.GuestPreferences, error)
	UpdateGuestPreferences(guestID uint, prefs *models.GuestPreferences) (*models.GuestPreferences, error)
	GetGuestAIInsights(guestID uint) (*models.GuestAIInsights, error)
	UpdateGuestAIInsights(guestID uint, insights *models.GuestAIInsights) (*models.GuestAIInsights, error)
	GetGuestStatistics(guestID uint) (*models.GuestStatistics, error)
	GetGuestServiceUsage(guestID uint) ([]models.ServiceUsageItem, error)
}

// guestRepository implements GuestRepository
type guestRepository struct {
	db *gorm.DB
}

// NewGuestRepository creates a new guest repository
func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepository{db: db}
}

// CreateGuest creates a new guest
func (r *guestRepository) CreateGuest(guest *models.Guest) (*models.Guest, error) {
	if err := r.db.Create(guest).Error; err != nil {
		return nil, err
	}

	// Initialize preferences and AI insights
	prefs := &models.GuestPreferences{GuestID: guest.ID}
	if err := r.db.Create(prefs).Error; err != nil {
		return nil, err
	}

	insights := &models.GuestAIInsights{
		GuestID:   guest.ID,
		RiskScore: "low",
	}
	if err := r.db.Create(insights).Error; err != nil {
		return nil, err
	}

	return r.GetGuestByID(guest.ID)
}

// GetGuestByID retrieves a guest by ID with all relations
func (r *guestRepository) GetGuestByID(id uint) (*models.Guest, error) {
	var guest models.Guest
	if err := r.db.
		Preload("Reservations").
		Preload("ServiceRequests").
		Preload("Preferences").
		Preload("AIInsights").
		First(&guest, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guest not found")
		}
		return nil, err
	}
	return &guest, nil
}

// GetAllGuests retrieves all guests with pagination and search
func (r *guestRepository) GetAllGuests(page, pageSize int, search string) ([]models.Guest, int64, error) {
	var guests []models.Guest
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Guest{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&guests).Error; err != nil {
		return nil, 0, err
	}

	return guests, total, nil
}

// UpdateGuest updates a guest
func (r *guestRepository) UpdateGuest(id uint, guest *models.Guest) (*models.Guest, error) {
	if err := r.db.Model(&models.Guest{}).Where("id = ?", id).Updates(guest).Error; err != nil {
		return nil, err
	}
	return r.GetGuestByID(id)
}

// DeleteGuest deletes a guest (soft delete)
func (r *guestRepository) DeleteGuest(id uint) error {
	return r.db.Delete(&models.Guest{}, id).Error
}

// GetGuestByEmail retrieves a guest by email
func (r *guestRepository) GetGuestByEmail(email string) (*models.Guest, error) {
	var guest models.Guest
	if err := r.db.Where("email = ?", email).First(&guest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guest not found")
		}
		return nil, err
	}
	return &guest, nil
}

// GetGuestHistory retrieves guest history with statistics
func (r *guestRepository) GetGuestHistory(guestID uint) (*models.GuestHistory, error) {
	var reservations []models.Reservation

	// Get all reservations for guest
	if err := r.db.Where("guest_id = ?", guestID).
		Order("check_in_date DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	// Calculate statistics using separate queries to handle NULLs
	var totalStays int64
	r.db.Model(&models.Reservation{}).
		Where("guest_id = ? AND status IN ?", guestID, []string{"completed", "checked-out"}).
		Count(&totalStays)

	var totalSpent float64
	var lastVisit time.Time
	averageSpend := 0.0

	if totalStays > 0 {
		r.db.Model(&models.Reservation{}).
			Where("guest_id = ? AND status IN ?", guestID, []string{"completed", "checked-out"}).
			Select("COALESCE(SUM(total_price), 0)").
			Scan(&totalSpent)

		var lastVisitPtr *time.Time
		r.db.Model(&models.Reservation{}).
			Where("guest_id = ? AND status IN ?", guestID, []string{"completed", "checked-out"}).
			Select("MAX(check_out_date)").
			Scan(&lastVisitPtr)
		if lastVisitPtr != nil {
			lastVisit = *lastVisitPtr
		}

		averageSpend = totalSpent / float64(totalStays)
	}

	history := &models.GuestHistory{
		TotalStays:    int(totalStays),
		TotalSpent:    totalSpent,
		AverageSpend:  averageSpend,
		LastVisitDate: lastVisit,
		Reservations:  reservations,
	}

	return history, nil
}

// GetGuestPreferences retrieves guest preferences
func (r *guestRepository) GetGuestPreferences(guestID uint) (*models.GuestPreferences, error) {
	var prefs models.GuestPreferences
	if err := r.db.Where("guest_id = ?", guestID).First(&prefs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("preferences not found")
		}
		return nil, err
	}
	return &prefs, nil
}

// UpdateGuestPreferences updates guest preferences
func (r *guestRepository) UpdateGuestPreferences(guestID uint, prefs *models.GuestPreferences) (*models.GuestPreferences, error) {
	if err := r.db.Model(&models.GuestPreferences{}).
		Where("guest_id = ?", guestID).
		Updates(prefs).Error; err != nil {
		return nil, err
	}
	return r.GetGuestPreferences(guestID)
}

// GetGuestAIInsights retrieves guest AI insights
func (r *guestRepository) GetGuestAIInsights(guestID uint) (*models.GuestAIInsights, error) {
	var insights models.GuestAIInsights
	if err := r.db.Where("guest_id = ?", guestID).First(&insights).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("AI insights not found")
		}
		return nil, err
	}
	return &insights, nil
}

// UpdateGuestAIInsights updates guest AI insights
func (r *guestRepository) UpdateGuestAIInsights(guestID uint, insights *models.GuestAIInsights) (*models.GuestAIInsights, error) {
	if err := r.db.Model(&models.GuestAIInsights{}).
		Where("guest_id = ?", guestID).
		Updates(insights).Error; err != nil {
		return nil, err
	}
	return r.GetGuestAIInsights(guestID)
}

// GetGuestDetailsByID retrieves a guest with all relations, statistics, and service usage
func (r *guestRepository) GetGuestDetailsByID(id uint) (*models.GuestDetailsResponse, error) {
	var guest models.Guest
	if err := r.db.
		Preload("Reservations").
		Preload("ServiceRequests").
		Preload("Preferences").
		Preload("AIInsights").
		First(&guest, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guest not found")
		}
		return nil, err
	}

	// Get statistics
	stats, _ := r.GetGuestStatistics(id)
	if stats == nil {
		stats = &models.GuestStatistics{}
	}

	// Get service usage
	serviceUsage, _ := r.GetGuestServiceUsage(id)
	if serviceUsage == nil {
		serviceUsage = []models.ServiceUsageItem{}
	}

	return &models.GuestDetailsResponse{
		ID:           guest.ID,
		Name:         guest.Name,
		Email:        guest.Email,
		Phone:        guest.Phone,
		Nationality:  guest.Nationality,
		IDType:       guest.IDType,
		IDNumber:     guest.IDNumber,
		JoinDate:     guest.JoinDate,
		CreatedAt:    guest.CreatedAt,
		UpdatedAt:    guest.UpdatedAt,
		Reservations: guest.Reservations,
		Preferences:  guest.Preferences,
		AIInsights:   guest.AIInsights,
		Statistics:   *stats,
		ServiceUsage: serviceUsage,
	}, nil
}

// GetGuestStatistics calculates statistics from guest's reservation history
func (r *guestRepository) GetGuestStatistics(guestID uint) (*models.GuestStatistics, error) {
	var stats models.GuestStatistics

	// Get completed reservations count and totals using separate queries to handle NULLs
	var totalStays int64
	r.db.Model(&models.Reservation{}).
		Where("guest_id = ? AND status IN ?", guestID, []string{"completed", "checked-out"}).
		Count(&totalStays)

	stats.TotalStays = int(totalStays)

	if totalStays > 0 {
		// Only query for sum and max if there are reservations
		var totalSpent float64
		r.db.Model(&models.Reservation{}).
			Where("guest_id = ? AND status IN ?", guestID, []string{"completed", "checked-out"}).
			Select("COALESCE(SUM(total_price), 0)").
			Scan(&totalSpent)
		stats.TotalSpent = totalSpent

		var lastVisit *time.Time
		r.db.Model(&models.Reservation{}).
			Where("guest_id = ? AND status IN ?", guestID, []string{"completed", "checked-out"}).
			Select("MAX(check_out_date)").
			Scan(&lastVisit)
		if lastVisit != nil {
			stats.LastVisit = *lastVisit
		}

		stats.AverageSpend = stats.TotalSpent / float64(stats.TotalStays)
	}

	// Find most common room type (placeholder - requires Room relation)
	stats.MostCommonRoom = ""

	return &stats, nil
}

// GetGuestServiceUsage aggregates service request types for a guest
func (r *guestRepository) GetGuestServiceUsage(guestID uint) ([]models.ServiceUsageItem, error) {
	var results []struct {
		ServiceType string
		Count       int64
	}

	err := r.db.Model(&models.ServiceRequest{}).
		Select("service_type, COUNT(*) as count").
		Where("guest_id = ?", guestID).
		Group("service_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Map to friendly labels
	serviceLabels := map[string]string{
		"room-service":       "Room Service",
		"housekeeping":       "Housekeeping",
		"maintenance":        "Maintenance",
		"special-requests":   "Special Requests",
		"transportation":     "Transportation",
		"general-assistance": "General Assistance",
	}

	var usage []models.ServiceUsageItem
	for _, result := range results {
		label := serviceLabels[result.ServiceType]
		if label == "" {
			label = result.ServiceType
		}

		usage = append(usage, models.ServiceUsageItem{
			Type:  result.ServiceType,
			Label: label,
			Count: int(result.Count),
		})
	}

	return usage, nil
}
