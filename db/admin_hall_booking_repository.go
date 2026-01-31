package db

import (
	"errors"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// AdminHallBookingRepository interface for admin booking operations
type AdminHallBookingRepository interface {
	// Booking management
	GetAllBookings(page, pageSize int, status []string, dateFrom, dateTo, search, sortBy, sortOrder string) ([]models.HallBooking, int64, error)
	GetBookingByID(id uint) (*models.HallBooking, error)
	UpdateBookingStatus(id uint, status string, adminID uint, notes string) (*models.HallBooking, error)

	// Status history
	GetBookingStatusHistory(bookingID uint) ([]models.BookingStatusHistory, error)
	CreateStatusHistory(history *models.BookingStatusHistory) error

	// Statistics
	GetBookingStats(period, dateFrom, dateTo string) (map[string]interface{}, error)
}

type adminHallBookingRepository struct {
	db *gorm.DB
}

// NewAdminHallBookingRepository creates a new admin hall booking repository
func NewAdminHallBookingRepository(db *gorm.DB) AdminHallBookingRepository {
	return &adminHallBookingRepository{db: db}
}

// GetAllBookings retrieves all bookings with filtering and pagination
func (r *adminHallBookingRepository) GetAllBookings(page, pageSize int, status []string, dateFrom, dateTo, search, sortBy, sortOrder string) ([]models.HallBooking, int64, error) {
	var bookings []models.HallBooking
	var total int64

	query := r.db.Model(&models.HallBooking{})
	// Preload("ConfirmedByUser").      // Temporarily disabled
	// Preload("CancelledByUser").      // Temporarily disabled
	// Preload("UpdatedByUser").        // Temporarily disabled
	// Preload("StatusHistory.ChangedByUser") // Temporarily disabled

	// Apply filters
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	if dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
			query = query.Where("booking_date >= ?", startOfDay)
		}
	}

	if dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
			query = query.Where("booking_date <= ?", endOfDay)
		}
	}

	if search != "" {
		query = query.Where("organizer_name ILIKE ? OR organizer_email ILIKE ? OR booking_id ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		orderClause := sortBy
		if sortOrder == "desc" {
			orderClause += " DESC"
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

// GetBookingByID retrieves a booking by ID with all relations
func (r *adminHallBookingRepository) GetBookingByID(id uint) (*models.HallBooking, error) {
	var booking models.HallBooking
	if err := r.db.
		// Preload("ConfirmedByUser").      // Temporarily disabled
		// Preload("CancelledByUser").      // Temporarily disabled
		// Preload("UpdatedByUser").        // Temporarily disabled
		// Preload("StatusHistory.ChangedByUser"). // Temporarily disabled
		Preload("Payments").
		Preload("Invoice").
		First(&booking, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	return &booking, nil
}

// UpdateBookingStatus updates booking status with admin tracking
func (r *adminHallBookingRepository) UpdateBookingStatus(id uint, status string, adminID uint, notes string) (*models.HallBooking, error) {
	var booking models.HallBooking

	// Get current booking
	if err := r.db.First(&booking, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	// Store old status for history
	oldStatus := booking.Status

	// Update booking status and tracking fields
	now := time.Now()
	booking.Status = status
	booking.UpdatedBy = &adminID

	switch status {
	case "confirmed":
		booking.ConfirmedBy = &adminID
		booking.ConfirmedAt = &now
	case "cancelled":
		booking.CancelledBy = &adminID
		booking.CancelledAt = &now
	}

	// Start transaction
	tx := r.db.Begin()

	// Update booking
	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create status history
	history := models.BookingStatusHistory{
		BookingID: id,
		OldStatus: &oldStatus,
		NewStatus: status,
		ChangedBy: adminID,
		ChangedAt: now,
		Notes:     notes,
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Return updated booking with relations
	return r.GetBookingByID(id)
}

// GetBookingStatusHistory retrieves status history for a booking
func (r *adminHallBookingRepository) GetBookingStatusHistory(bookingID uint) ([]models.BookingStatusHistory, error) {
	var history []models.BookingStatusHistory
	if err := r.db.
		// Preload("ChangedByUser"). // Temporarily disabled
		Where("booking_id = ?", bookingID).
		Order("changed_at DESC").
		Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// CreateStatusHistory creates a new status history entry
func (r *adminHallBookingRepository) CreateStatusHistory(history *models.BookingStatusHistory) error {
	return r.db.Create(history).Error
}

// GetBookingStats retrieves booking statistics
func (r *adminHallBookingRepository) GetBookingStats(period, dateFrom, dateTo string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	baseQuery := r.db.Model(&models.HallBooking{})

	// Apply date filters based on period
	now := time.Now()
	var dateCondition string
	var dateValue interface{}

	switch period {
	case "today":
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		dateCondition = "booking_date BETWEEN ? AND ?"
		dateValue = []interface{}{startOfDay, endOfDay}
	case "week":
		weekStart := now.AddDate(0, 0, -int(now.Weekday()))
		startOfWeek := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())
		dateCondition = "booking_date >= ?"
		dateValue = startOfWeek
	case "month":
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		dateCondition = "booking_date >= ?"
		dateValue = startOfMonth
	case "year":
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		dateCondition = "booking_date >= ?"
		dateValue = startOfYear
	case "custom":
		if dateFrom != "" {
			if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
				startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
				dateCondition = "booking_date >= ?"
				dateValue = startOfDay
			}
		}
		if dateTo != "" {
			if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
				endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
				if dateCondition == "" {
					dateCondition = "booking_date <= ?"
					dateValue = endOfDay
				} else {
					dateCondition = "booking_date >= ? AND booking_date <= ?"
					dateValue = []interface{}{dateValue, endOfDay}
				}
			}
		}
	default:
		dateCondition = "1=1" // No date filter
		dateValue = []interface{}{}
	}

	// Apply date condition to base query
	if dateCondition != "1=1" {
		baseQuery = baseQuery.Where(dateCondition, dateValue)
	}

	// Total bookings
	var totalBookings int64
	baseQuery.Count(&totalBookings)
	stats["total_bookings"] = totalBookings

	// Status counts
	var pendingCount, confirmedCount, completedCount, cancelledCount int64
	baseQuery.Where("status = ?", "pending").Count(&pendingCount)
	baseQuery.Where("status = ?", "confirmed").Count(&confirmedCount)
	baseQuery.Where("status = ?", "completed").Count(&completedCount)
	baseQuery.Where("status = ?", "cancelled").Count(&cancelledCount)

	stats["pending_bookings"] = pendingCount
	stats["confirmed_bookings"] = confirmedCount
	stats["completed_bookings"] = completedCount
	stats["cancelled_bookings"] = cancelledCount

	// Revenue calculations
	var confirmedRevenue, completedRevenue float64
	revenueQuery := r.db.Model(&models.HallBooking{})
	if dateCondition != "1=1" {
		revenueQuery = revenueQuery.Where(dateCondition, dateValue)
	}

	revenueQuery.Where("status = ?", "confirmed").Select("COALESCE(SUM(total_price), 0)").Scan(&confirmedRevenue)
	revenueQuery.Where("status = ?", "completed").Select("COALESCE(SUM(total_price), 0)").Scan(&completedRevenue)

	stats["total_revenue"] = confirmedRevenue + completedRevenue
	stats["revenue_by_status"] = map[string]float64{
		"confirmed": confirmedRevenue,
		"completed": completedRevenue,
	}

	// Popular event types
	var eventTypes []struct {
		EventType string `json:"event_type"`
		Count     int64  `json:"count"`
	}

	eventQuery := r.db.Model(&models.HallBooking{})
	if dateCondition != "1=1" {
		eventQuery = eventQuery.Where(dateCondition, dateValue)
	}

	eventQuery.Select("event_type, COUNT(*) as count").
		Group("event_type").
		Order("count DESC").
		Limit(5).
		Scan(&eventTypes)

	stats["popular_event_types"] = eventTypes

	// NEW: Calculate average guests
	var totalGuests int64
	var bookingCount int64
	avgGuestsQuery := r.db.Model(&models.HallBooking{})
	if dateCondition != "1=1" {
		avgGuestsQuery = avgGuestsQuery.Where(dateCondition, dateValue)
	}
	avgGuestsQuery.Select("COALESCE(SUM(guest_count), 0)").Scan(&totalGuests)
	avgGuestsQuery.Count(&bookingCount)

	var averageGuests *float64
	if bookingCount > 0 {
		avg := float64(totalGuests) / float64(bookingCount)
		averageGuests = &avg
	}
	stats["average_guests"] = averageGuests

	// NEW: Calculate occupancy rate (assuming 100 guests max per hall booking)
	var totalBookedGuests int64
	occupancyQuery := r.db.Model(&models.HallBooking{})
	if dateCondition != "1=1" {
		occupancyQuery = occupancyQuery.Where(dateCondition, dateValue)
	}
	occupancyQuery.Select("COALESCE(SUM(guest_count), 0)").Scan(&totalBookedGuests)

	var occupancyRate *float64
	if bookingCount > 0 {
		totalCapacity := float64(bookingCount) * 100.0 // 100 guests max per booking
		if totalCapacity > 0 {
			occupancy := (float64(totalBookedGuests) / totalCapacity) * 100.0
			occupancyRate = &occupancy
		}
	}
	stats["occupancy_rate"] = occupancyRate

	// NEW: Calculate monthly revenue data
	var monthlyRevenue []struct {
		Month   string  `json:"month"`
		Revenue float64 `json:"revenue"`
	}

	revenueByMonthQuery := r.db.Model(&models.HallBooking{})
	if dateCondition != "1=1" {
		revenueByMonthQuery = revenueByMonthQuery.Where(dateCondition, dateValue)
	}

	revenueByMonthQuery.Select("DATE_TRUNC('month', booking_date)::text as month, COALESCE(SUM(total_price), 0) as revenue").
		Where("status IN ?", []string{"confirmed", "completed"}).
		Group("DATE_TRUNC('month', booking_date)").
		Order("month ASC").
		Scan(&monthlyRevenue)

	stats["monthly_revenue"] = monthlyRevenue

	// NEW: Calculate monthly booking volume data
	var monthlyBookings []struct {
		Month    string `json:"month"`
		Bookings int    `json:"bookings"`
	}

	bookingsByMonthQuery := r.db.Model(&models.HallBooking{})
	if dateCondition != "1=1" {
		bookingsByMonthQuery = bookingsByMonthQuery.Where(dateCondition, dateValue)
	}

	bookingsByMonthQuery.Select("DATE_TRUNC('month', booking_date)::text as month, COUNT(*) as bookings").
		Group("DATE_TRUNC('month', booking_date)").
		Order("month ASC").
		Scan(&monthlyBookings)

	stats["monthly_bookings"] = monthlyBookings

	return stats, nil
}
