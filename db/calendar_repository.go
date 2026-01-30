package db

import (
	"fmt"
	"hotel/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// CalendarRepository defines calendar database operations
type CalendarRepository interface {
	// Daily availability
	GetAvailabilityByDate(date time.Time) (*models.CalendarAvailability, error)
	GetAvailabilityByDateRange(startDate, endDate time.Time) ([]models.CalendarAvailability, error)
	CreateOrUpdateAvailability(availability *models.CalendarAvailability) error

	// Time slot management
	GetTimeSlotsByDate(date time.Time) ([]models.TimeSlot, error)
	CreateTimeSlots(slots []models.TimeSlot) error
	UpdateTimeSlot(slot *models.TimeSlot) error

	// Calendar generation
	GenerateMonthlyCalendar(year int, month time.Month) error
	UpdateAvailabilityFromBookings(date time.Time) error

	// Statistics
	GetMonthlyStats(year int, month time.Month) (map[string]interface{}, error)
	GetAvailabilitySummary(startDate, endDate time.Time) (map[string]interface{}, error)
}

// calendarRepository implements CalendarRepository
type calendarRepository struct {
	db *gorm.DB
}

// NewCalendarRepository creates a new calendar repository
func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{db: db}
}

// GetAvailabilityByDate retrieves availability for a specific date
func (r *calendarRepository) GetAvailabilityByDate(date time.Time) (*models.CalendarAvailability, error) {
	var availability models.CalendarAvailability

	// Normalize date to start of day
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	err := r.db.Where("date = ?", normalizedDate).First(&availability).Error
	if err == gorm.ErrRecordNotFound {
		// Create default availability if not found
		availability = models.CalendarAvailability{
			Date:           normalizedDate,
			Status:         models.StatusAvailable,
			TotalSlots:     r.calculateTotalSlots(normalizedDate),
			AvailableSlots: r.calculateTotalSlots(normalizedDate),
		}
		err = r.db.Create(&availability).Error
	}

	return &availability, err
}

// GetAvailabilityByDateRange retrieves availability for a date range
func (r *calendarRepository) GetAvailabilityByDateRange(startDate, endDate time.Time) ([]models.CalendarAvailability, error) {
	var availabilities []models.CalendarAvailability

	// Normalize dates
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date ASC").
		Find(&availabilities).Error

	if err == nil && len(availabilities) == 0 {
		for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
			availability := models.CalendarAvailability{
				Date:           date,
				Status:         models.StatusAvailable,
				TotalSlots:     r.calculateTotalSlots(date),
				AvailableSlots: r.calculateTotalSlots(date),
			}

			if err := r.db.Create(&availability).Error; err != nil {
			}

			availabilities = append(availabilities, availability)

			r.generateTimeSlotsForDate(date)
		}
	}

	return availabilities, err
}

// CreateOrUpdateAvailability creates or updates availability
func (r *calendarRepository) CreateOrUpdateAvailability(availability *models.CalendarAvailability) error {
	// Normalize date
	availability.Date = time.Date(availability.Date.Year(), availability.Date.Month(), availability.Date.Day(), 0, 0, 0, 0, availability.Date.Location())

	var existing models.CalendarAvailability
	err := r.db.Where("date = ?", availability.Date).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(availability).Error
	}

	availability.ID = existing.ID
	return r.db.Save(availability).Error
}

// GetTimeSlotsByDate retrieves all time slots for a date
func (r *calendarRepository) GetTimeSlotsByDate(date time.Time) ([]models.TimeSlot, error) {
	var slots []models.TimeSlot

	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	err := r.db.Where("date = ?", normalizedDate).
		Order("start_time ASC").
		Find(&slots).Error

	return slots, err
}

// CreateTimeSlots creates multiple time slots
func (r *calendarRepository) CreateTimeSlots(slots []models.TimeSlot) error {
	for _, slot := range slots {
		// Normalize date
		slot.Date = time.Date(slot.Date.Year(), slot.Date.Month(), slot.Date.Day(), 0, 0, 0, 0, slot.Date.Location())

		var existing models.TimeSlot
		err := r.db.Where("date = ? AND start_time = ?", slot.Date, slot.StartTime).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := r.db.Create(&slot).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateTimeSlot updates a time slot
func (r *calendarRepository) UpdateTimeSlot(slot *models.TimeSlot) error {
	// Normalize date
	slot.Date = time.Date(slot.Date.Year(), slot.Date.Month(), slot.Date.Day(), 0, 0, 0, 0, slot.Date.Location())

	var existing models.TimeSlot
	err := r.db.Where("date = ? AND start_time = ?", slot.Date, slot.StartTime).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(slot).Error
	}

	slot.ID = existing.ID
	return r.db.Save(slot).Error
}

// GenerateMonthlyCalendar creates availability records for a month
func (r *calendarRepository) GenerateMonthlyCalendar(year int, month time.Month) error {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	for date := firstDay; date.After(lastDay) == false; date = date.AddDate(0, 0, 1) {
		availability := models.CalendarAvailability{
			Date:           date,
			Status:         models.StatusAvailable,
			TotalSlots:     r.calculateTotalSlots(date),
			AvailableSlots: r.calculateTotalSlots(date),
		}

		// Create or update
		r.CreateOrUpdateAvailability(&availability)

		// Generate time slots
		r.generateTimeSlotsForDate(date)

		// Update availability based on existing bookings
		r.UpdateAvailabilityFromBookings(date)
	}

	return nil
}

// UpdateAvailabilityFromBookings updates availability based on existing bookings
func (r *calendarRepository) UpdateAvailabilityFromBookings(date time.Time) error {
	var bookings []models.HallBooking

	// Normalize date to start of day for comparison
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	log.Printf("DEBUG: Looking for bookings between %s and %s", normalizedDate.Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))

	err := r.db.Where("booking_date BETWEEN ? AND ? AND status NOT IN ?",
		normalizedDate, endOfDay, []string{"cancelled"}).
		Find(&bookings).Error
	if err != nil {
		return err
	}

	log.Printf("DEBUG: Found %d bookings for date %s", len(bookings), normalizedDate.Format("2006-01-02"))

	availability, err := r.GetAvailabilityByDate(date)
	if err != nil {
		return err
	}

	// Generate all possible time slots for this date
	allSlots := r.generateDefaultTimeSlots(date)

	// Check which slots are booked
	bookedSlotsCount := 0
	for _, slot := range allSlots {
		// Check if this slot conflicts with any booking
		for _, booking := range bookings {
			if r.timeSlotConflicts(slot.StartTime, slot.EndTime, booking.StartTime, booking.EndTime) {
				bookedSlotsCount++
				break // Slot is booked, move to next slot
			}
		}
	}

	// Update availability based on actual slot conflicts
	availability.BookedSlots = bookedSlotsCount
	availability.AvailableSlots = availability.TotalSlots - bookedSlotsCount

	// Update status based on availability
	if availability.AvailableSlots <= 0 {
		availability.Status = models.StatusBooked
	} else if availability.BookedSlots > 0 {
		availability.Status = models.StatusPending
	} else {
		availability.Status = models.StatusAvailable
	}

	return r.CreateOrUpdateAvailability(availability)
}

// GetMonthlyStats returns statistics for a month
func (r *calendarRepository) GetMonthlyStats(year int, month time.Month) (map[string]interface{}, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	availabilities, err := r.GetAvailabilityByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_days":       len(availabilities),
		"available_days":   0,
		"booked_days":      0,
		"pending_days":     0,
		"closed_days":      0,
		"maintenance_days": 0,
		"total_bookings":   0,
		"revenue":          0.0,
	}

	for _, availability := range availabilities {
		switch availability.Status {
		case models.StatusAvailable:
			stats["available_days"] = stats["available_days"].(int) + 1
		case models.StatusBooked:
			stats["booked_days"] = stats["booked_days"].(int) + 1
		case models.StatusPending:
			stats["pending_days"] = stats["pending_days"].(int) + 1
		case models.StatusClosed:
			stats["closed_days"] = stats["closed_days"].(int) + 1
		case models.StatusMaintenance:
			stats["maintenance_days"] = stats["maintenance_days"].(int) + 1
		}

		stats["total_bookings"] = stats["total_bookings"].(int) + availability.BookedSlots
	}

	return stats, nil
}

// GetAvailabilitySummary returns availability summary for a date range
func (r *calendarRepository) GetAvailabilitySummary(startDate, endDate time.Time) (map[string]interface{}, error) {
	availabilities, err := r.GetAvailabilityByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_days":       len(availabilities),
		"available_days":   0,
		"booked_days":      0,
		"pending_days":     0,
		"closed_days":      0,
		"maintenance_days": 0,
		"total_slots":      0,
		"available_slots":  0,
		"booked_slots":     0,
	}

	for _, availability := range availabilities {
		switch availability.Status {
		case models.StatusAvailable:
			summary["available_days"] = summary["available_days"].(int) + 1
		case models.StatusBooked:
			summary["booked_days"] = summary["booked_days"].(int) + 1
		case models.StatusPending:
			summary["pending_days"] = summary["pending_days"].(int) + 1
		case models.StatusClosed:
			summary["closed_days"] = summary["closed_days"].(int) + 1
		case models.StatusMaintenance:
			summary["maintenance_days"] = summary["maintenance_days"].(int) + 1
		}

		summary["total_slots"] = summary["total_slots"].(int) + availability.TotalSlots
		summary["available_slots"] = summary["available_slots"].(int) + availability.AvailableSlots
		summary["booked_slots"] = summary["booked_slots"].(int) + availability.BookedSlots
	}

	return summary, nil
}

// Helper functions

// calculateTotalSlots calculates the total number of slots for a date
func (r *calendarRepository) calculateTotalSlots(date time.Time) int {
	dayOfWeek := date.Weekday()

	switch dayOfWeek {
	case time.Sunday:
		return 6 // 12:00 - 18:00 (6 slots)
	case time.Saturday:
		return 11 // 12:00 - 23:00 (11 slots)
	default: // Monday - Friday
		return 14 // 08:00 - 22:00 (14 slots)
	}
}

// generateTimeSlotsForDate generates default time slots for a date
func (r *calendarRepository) generateTimeSlotsForDate(date time.Time) error {
	slots := r.generateDefaultTimeSlots(date)
	return r.CreateTimeSlots(slots)
}

// generateDefaultTimeSlots generates default time slots for a date
func (r *calendarRepository) generateDefaultTimeSlots(date time.Time) []models.TimeSlot {
	var slots []models.TimeSlot
	dayOfWeek := date.Weekday()

	// Get the calendar availability for this date
	availability, err := r.GetAvailabilityByDate(date)
	if err != nil || availability == nil {
		log.Printf("ERROR: Could not get availability for date %s: %v", date.Format("2006-01-02"), err)
		return slots // Return empty if availability not found
	}

	var startHour, endHour int

	switch dayOfWeek {
	case time.Sunday:
		startHour, endHour = 12, 18
	case time.Saturday:
		startHour, endHour = 12, 23
	default:
		startHour, endHour = 8, 22
	}

	for hour := startHour; hour < endHour; hour++ {
		slots = append(slots, models.TimeSlot{
			Date:                   date,
			StartTime:              fmt.Sprintf("%02d:00", hour),
			EndTime:                fmt.Sprintf("%02d:00", hour+1),
			Status:                 models.StatusAvailable,
			MaxCapacity:            100,
			Price:                  r.calculateHourlyPrice(date, hour),
			CalendarAvailabilityID: availability.ID,
		})
	}

	return slots
}

// calculateHourlyPrice calculates the price for a time slot based on date and hour
func (r *calendarRepository) calculateHourlyPrice(date time.Time, hour int) float64 {
	dayOfWeek := date.Weekday()

	// Weekend pricing
	if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
		return 25.0
	}

	// Peak hours (evening)
	if hour >= 17 && hour <= 21 {
		return 20.0
	}

	// Standard hours
	return 15.0
}

// timeSlotConflicts checks if two time slots conflict
func (r *calendarRepository) timeSlotConflicts(start1, end1, start2, end2 string) bool {
	// Convert strings to time objects for comparison
	t1, _ := time.Parse("15:04", start1)
	t2, _ := time.Parse("15:04", end1)
	t3, _ := time.Parse("15:04", start2)
	t4, _ := time.Parse("15:04", end2)

	// Check if slots overlap
	return (t1.Before(t4) && t2.After(t3)) || // Overlaps
		(t1.Equal(t3) && t2.Equal(t4)) || // Same slot
		(t1.Before(t3) && t2.After(t4)) // Contains
}
