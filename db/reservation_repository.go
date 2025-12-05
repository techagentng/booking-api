package db

import (
	"errors"
	"fmt"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// ReservationRepository defines reservation database operations
type ReservationRepository interface {
	CreateReservation(reservation *models.Reservation) (*models.Reservation, error)
	GetReservationByID(id uint) (*models.Reservation, error)
	GetAllReservations(params ReservationQueryParams) ([]models.Reservation, int64, error)
	UpdateReservation(id uint, reservation *models.Reservation) (*models.Reservation, error)
	UpdateReservationStatus(id uint, status string) error
	CancelReservation(id uint, reason string, refundAmount float64) error
	CheckIn(id uint, req *models.CheckInRequest) error
	CheckOut(id uint, checkedOutAt time.Time, additionalCharges float64, notes string) error
	GetReservationStats(fromDate, toDate *time.Time) (*models.ReservationStats, error)
	GetCheckInStats(date time.Time) (*models.CheckInStats, error)
	GetReservationsByGuestID(guestID uint) ([]models.Reservation, error)
	GetReservationsByRoomID(roomID uint) ([]models.Reservation, error)
	IsRoomAvailable(roomID uint, checkIn, checkOut time.Time, excludeReservationID uint) (bool, error)
	GenerateConfirmationNumber() (string, error)
	GetDashboardStats() (*models.DashboardStats, error)
}

// ReservationQueryParams holds query parameters for listing reservations
type ReservationQueryParams struct {
	Page        int
	PageSize    int
	Status      string
	GuestID     uint
	RoomID      uint
	CheckInFrom *time.Time
	CheckInTo   *time.Time
	CheckInDate *time.Time // Exact check-in date filter
	Search      string
}

// reservationRepository implements ReservationRepository
type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new reservation repository
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// CreateReservation creates a new reservation
func (r *reservationRepository) CreateReservation(reservation *models.Reservation) (*models.Reservation, error) {
	if err := r.db.Create(reservation).Error; err != nil {
		return nil, err
	}
	// Reload with relations
	return r.GetReservationByID(reservation.ID)
}

// GetReservationByID retrieves a reservation by ID with relations
func (r *reservationRepository) GetReservationByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.Preload("Guest").Preload("Room").First(&reservation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}
	return &reservation, nil
}

// GetAllReservations retrieves all reservations with filters and pagination
func (r *reservationRepository) GetAllReservations(params ReservationQueryParams) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	query := r.db.Model(&models.Reservation{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.GuestID > 0 {
		query = query.Where("guest_id = ?", params.GuestID)
	}
	if params.RoomID > 0 {
		query = query.Where("room_id = ?", params.RoomID)
	}
	if params.CheckInFrom != nil {
		query = query.Where("check_in_date >= ?", *params.CheckInFrom)
	}
	if params.CheckInTo != nil {
		query = query.Where("check_in_date <= ?", *params.CheckInTo)
	}
	// Exact check-in date filter (for today's arrivals)
	if params.CheckInDate != nil {
		startOfDay := time.Date(params.CheckInDate.Year(), params.CheckInDate.Month(), params.CheckInDate.Day(), 0, 0, 0, 0, params.CheckInDate.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay)
	}
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Joins("LEFT JOIN guests ON guests.id = reservations.guest_id").
			Where("confirmation_number ILIKE ? OR guests.name ILIKE ? OR guests.email ILIKE ?",
				searchPattern, searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and relations
	if err := query.
		Preload("Guest").
		Preload("Room").
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

// UpdateReservation updates a reservation
func (r *reservationRepository) UpdateReservation(id uint, reservation *models.Reservation) (*models.Reservation, error) {
	if err := r.db.Model(&models.Reservation{}).Where("id = ?", id).Updates(reservation).Error; err != nil {
		return nil, err
	}
	return r.GetReservationByID(id)
}

// UpdateReservationStatus updates only the status
func (r *reservationRepository) UpdateReservationStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}

// CancelReservation cancels a reservation
func (r *reservationRepository) CancelReservation(id uint, reason string, refundAmount float64) error {
	now := time.Now()
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":              "cancelled",
		"cancellation_reason": reason,
		"refund_amount":       refundAmount,
		"cancelled_at":        now,
	}).Error
}

// CheckIn performs check-in for a reservation
func (r *reservationRepository) CheckIn(id uint, req *models.CheckInRequest) error {
	checkedInAt := time.Now()
	if req.ActualCheckInTime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ActualCheckInTime); err == nil {
			checkedInAt = parsed
		}
	}

	updates := map[string]interface{}{
		"status":        "checked-in",
		"checked_in_at": checkedInAt,
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.DepositCollected > 0 {
		updates["deposit_collected"] = req.DepositCollected
	}
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Updates(updates).Error
}

// CheckOut performs check-out for a reservation
func (r *reservationRepository) CheckOut(id uint, checkedOutAt time.Time, additionalCharges float64, notes string) error {
	// Get current reservation to calculate final amount
	var reservation models.Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"status":             "checked-out",
		"checked_out_at":     checkedOutAt,
		"additional_charges": additionalCharges,
		"total_amount":       reservation.TotalAmount + additionalCharges,
	}
	if notes != "" {
		updates["notes"] = notes
	}
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Updates(updates).Error
}

// GetReservationStats retrieves reservation statistics
func (r *reservationRepository) GetReservationStats(fromDate, toDate *time.Time) (*models.ReservationStats, error) {
	stats := &models.ReservationStats{}

	query := r.db.Model(&models.Reservation{})
	if fromDate != nil {
		query = query.Where("created_at >= ?", *fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", *toDate)
	}

	// Total reservations
	query.Count(&stats.TotalReservations)

	// Count by status
	r.db.Model(&models.Reservation{}).Where("status = ?", "confirmed").Count(&stats.Confirmed)
	r.db.Model(&models.Reservation{}).Where("status = ?", "pending").Count(&stats.Pending)
	r.db.Model(&models.Reservation{}).Where("status = ?", "checked-in").Count(&stats.CheckedIn)
	r.db.Model(&models.Reservation{}).Where("status = ?", "checked-out").Count(&stats.CheckedOut)
	r.db.Model(&models.Reservation{}).Where("status = ?", "cancelled").Count(&stats.Cancelled)

	// Total revenue (from completed/checked-out reservations)
	r.db.Model(&models.Reservation{}).
		Where("status IN ?", []string{"checked-out", "confirmed", "checked-in"}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	// Average booking value
	if stats.TotalReservations > 0 {
		var avgValue float64
		r.db.Model(&models.Reservation{}).
			Where("status NOT IN ?", []string{"cancelled"}).
			Select("COALESCE(AVG(total_amount), 0)").
			Scan(&avgValue)
		stats.AverageBookingValue = avgValue
	}

	// Occupancy rate (checked-in rooms / total rooms)
	var totalRooms int64
	var occupiedRooms int64
	r.db.Model(&models.Room{}).Count(&totalRooms)
	r.db.Model(&models.Room{}).Where("status = ?", "occupied").Count(&occupiedRooms)
	if totalRooms > 0 {
		stats.OccupancyRate = float64(occupiedRooms) / float64(totalRooms) * 100
	}

	return stats, nil
}

// GetReservationsByGuestID retrieves all reservations for a guest
func (r *reservationRepository) GetReservationsByGuestID(guestID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Preload("Room").Where("guest_id = ?", guestID).Order("created_at DESC").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByRoomID retrieves all reservations for a room
func (r *reservationRepository) GetReservationsByRoomID(roomID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Preload("Guest").Where("room_id = ?", roomID).Order("created_at DESC").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// IsRoomAvailable checks if a room is available for the given date range
func (r *reservationRepository) IsRoomAvailable(roomID uint, checkIn, checkOut time.Time, excludeReservationID uint) (bool, error) {
	var count int64

	query := r.db.Model(&models.Reservation{}).
		Where("room_id = ?", roomID).
		Where("status NOT IN ?", []string{"cancelled", "checked-out"}).
		Where("(check_in_date < ? AND check_out_date > ?)", checkOut, checkIn)

	if excludeReservationID > 0 {
		query = query.Where("id != ?", excludeReservationID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}

// GenerateConfirmationNumber generates a unique confirmation number
func (r *reservationRepository) GenerateConfirmationNumber() (string, error) {
	year := time.Now().Year()

	// Get the count of reservations this year
	var count int64
	r.db.Model(&models.Reservation{}).
		Where("created_at >= ?", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)).
		Count(&count)

	// Generate confirmation number: BK-YYYY-NNN
	confirmationNumber := fmt.Sprintf("BK-%d-%03d", year, count+1)

	// Check if it exists (edge case)
	var existing models.Reservation
	for {
		err := r.db.Where("confirmation_number = ?", confirmationNumber).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		count++
		confirmationNumber = fmt.Sprintf("BK-%d-%03d", year, count+1)
	}

	return confirmationNumber, nil
}

// GetCheckInStats retrieves check-in statistics for a specific date
func (r *reservationRepository) GetCheckInStats(date time.Time) (*models.CheckInStats, error) {
	stats := &models.CheckInStats{
		Date: date.Format("2006-01-02"),
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Total arrivals (reservations with check-in date = today)
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status NOT IN ?", []string{"cancelled"}).
		Count(&stats.TotalArrivals)

	// Checked in today
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status = ?", "checked-in").
		Count(&stats.CheckedIn)

	// Pending check-in (confirmed but not checked in)
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status = ?", "confirmed").
		Count(&stats.PendingCheckIn)

	// Early check-ins (checked in before 14:00)
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status = ?", "checked-in").
		Where("checked_in_at IS NOT NULL").
		Where("EXTRACT(HOUR FROM checked_in_at) < 14").
		Count(&stats.EarlyCheckIns)

	// Late check-ins (checked in after 18:00)
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status = ?", "checked-in").
		Where("checked_in_at IS NOT NULL").
		Where("EXTRACT(HOUR FROM checked_in_at) >= 18").
		Count(&stats.LateCheckIns)

	// No shows (confirmed reservations past check-in date that weren't checked in)
	// For simplicity, count confirmed reservations for today that haven't checked in and it's past 23:00
	now := time.Now()
	if now.After(endOfDay.Add(-1 * time.Hour)) { // After 23:00
		r.db.Model(&models.Reservation{}).
			Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
			Where("status = ?", "confirmed").
			Count(&stats.NoShows)
	}

	// Average check-in time
	var avgHour float64
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status = ?", "checked-in").
		Where("checked_in_at IS NOT NULL").
		Select("COALESCE(AVG(EXTRACT(HOUR FROM checked_in_at) + EXTRACT(MINUTE FROM checked_in_at)/60), 14)").
		Scan(&avgHour)

	hours := int(avgHour)
	minutes := int((avgHour - float64(hours)) * 60)
	stats.AverageCheckInTime = fmt.Sprintf("%02d:%02d", hours, minutes)

	return stats, nil
}

// GetDashboardStats retrieves dashboard statistics for today
func (r *reservationRepository) GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// New bookings today (reservations created today)
	r.db.Model(&models.Reservation{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&stats.NewBookingsToday)

	// Scheduled bookings (confirmed reservations with future check-in dates)
	r.db.Model(&models.Reservation{}).
		Where("status IN (?, ?)", "confirmed", "pending").
		Where("check_in_date >= ?", startOfDay).
		Count(&stats.ScheduledBookings)

	// Check-ins today (reservations actually checked in today based on checked_in_at)
	r.db.Model(&models.Reservation{}).
		Where("checked_in_at >= ? AND checked_in_at < ?", startOfDay, endOfDay).
		Where("status IN (?, ?)", "checked-in", "checked-out").
		Count(&stats.CheckInsToday)

	// Check-outs today (reservations actually checked out today based on checked_out_at)
	r.db.Model(&models.Reservation{}).
		Where("checked_out_at >= ? AND checked_out_at < ?", startOfDay, endOfDay).
		Where("status = ?", "checked-out").
		Count(&stats.CheckOutsToday)

	// Pending check-ins (reservations with check_in_date = today, status = confirmed)
	r.db.Model(&models.Reservation{}).
		Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).
		Where("status = ?", "confirmed").
		Count(&stats.PendingCheckIns)

	// Pending check-outs (reservations with check_out_date = today, status = checked-in)
	r.db.Model(&models.Reservation{}).
		Where("check_out_date >= ? AND check_out_date < ?", startOfDay, endOfDay).
		Where("status = ?", "checked-in").
		Count(&stats.PendingCheckOuts)

	// Total rooms
	r.db.Model(&models.Room{}).Count(&stats.TotalRooms)

	// Available rooms today (rooms with status = available)
	r.db.Model(&models.Room{}).
		Where("status = ?", "available").
		Count(&stats.AvailableRoomsToday)

	// Occupied rooms today (rooms with status = occupied)
	r.db.Model(&models.Room{}).
		Where("status = ?", "occupied").
		Count(&stats.OccupiedRoomsToday)

	// Sold out rooms today = rooms that have reservations for today
	var bookedRoomCount int64
	r.db.Model(&models.Reservation{}).
		Where("check_in_date <= ? AND check_out_date > ?", startOfDay, startOfDay).
		Where("status IN (?, ?, ?)", "confirmed", "checked-in", "pending").
		Distinct("room_id").
		Count(&bookedRoomCount)
	stats.SoldOutRoomsToday = bookedRoomCount

	// Total guests in the system
	r.db.Model(&models.Guest{}).Count(&stats.TotalGuests)

	return stats, nil
}
