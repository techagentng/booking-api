package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"time"
)

// CalendarService defines calendar business logic operations
type CalendarService interface {
	// Public endpoints (no auth required)
	GetCalendarAvailability(year int, month time.Month) ([]models.CalendarAvailability, error)
	GetDailyAvailability(date time.Time) (*models.CalendarAvailability, error)
	GetTimeSlots(date time.Time) ([]models.TimeSlot, error)
	CheckSlotAvailability(date time.Time, startTime, endTime string) (bool, error)

	// Admin endpoints (auth required)
	UpdateDailyAvailability(date time.Time, status string, notes string) error
	UpdateTimeSlotStatus(date time.Time, startTime, endTime, status string) error
	GetMonthlyStats(year int, month time.Month) (*models.CalendarStatsResponse, error)
	GenerateCalendar(year int, month time.Month) error
	BulkUpdateAvailability(updates []models.CalendarUpdateRequest) error
}

// calendarService implements CalendarService
type calendarService struct {
	calendarRepo db.CalendarRepository
	bookingRepo  db.HallBookingRepository
}

// NewCalendarService creates a new calendar service
func NewCalendarService(calendarRepo db.CalendarRepository, bookingRepo db.HallBookingRepository) CalendarService {
	return &calendarService{
		calendarRepo: calendarRepo,
		bookingRepo:  bookingRepo,
	}
}

// GetCalendarAvailability returns availability for a month
func (s *calendarService) GetCalendarAvailability(year int, month time.Month) ([]models.CalendarAvailability, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	availabilities, err := s.calendarRepo.GetAvailabilityByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Update availability based on current bookings for each date
	for i := range availabilities {
		err := s.calendarRepo.UpdateAvailabilityFromBookings(availabilities[i].Date)
		if err != nil {
			// Log error but don't fail the entire operation
			continue
		}

		// Refresh the availability data
		updated, err := s.calendarRepo.GetAvailabilityByDate(availabilities[i].Date)
		if err == nil {
			availabilities[i] = *updated
		}
	}

	return availabilities, nil
}

// GetDailyAvailability returns availability for a specific date
func (s *calendarService) GetDailyAvailability(date time.Time) (*models.CalendarAvailability, error) {
	// Update availability based on current bookings
	err := s.calendarRepo.UpdateAvailabilityFromBookings(date)
	if err != nil {
		return nil, err
	}

	return s.calendarRepo.GetAvailabilityByDate(date)
}

// GetTimeSlots returns available time slots for a date
func (s *calendarService) GetTimeSlots(date time.Time) ([]models.TimeSlot, error) {
	// First update availability from bookings
	err := s.calendarRepo.UpdateAvailabilityFromBookings(date)
	if err != nil {
		return nil, err
	}

	slots, err := s.calendarRepo.GetTimeSlotsByDate(date)
	if err != nil {
		return nil, err
	}

	// Update slot status based on bookings
	bookings, err := s.bookingRepo.GetHallBookingsByDate(date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	// Update each slot's status based on bookings
	for i, slot := range slots {
		for _, booking := range bookings {
			if booking.Status != "cancelled" && s.timeSlotConflicts(slot.StartTime, slot.EndTime, booking.StartTime, booking.EndTime) {
				slots[i].Status = models.StatusBooked
				slots[i].CurrentUsage = 50 // Example usage calculation
				break
			}
		}
	}

	return slots, nil
}

// CheckSlotAvailability checks if a time slot is available
func (s *calendarService) CheckSlotAvailability(date time.Time, startTime, endTime string) (bool, error) {
	// Check if date is available
	dailyAvailability, err := s.calendarRepo.GetAvailabilityByDate(date)
	if err != nil {
		return false, err
	}

	if dailyAvailability.Status != models.StatusAvailable {
		return false, errors.New("date not available")
	}

	// Check time slots
	slots, err := s.calendarRepo.GetTimeSlotsByDate(date)
	if err != nil {
		return false, err
	}

	// Check if all requested slots are available
	for _, slot := range slots {
		if slot.StartTime >= startTime && slot.StartTime < endTime {
			if slot.Status != models.StatusAvailable {
				return false, nil
			}
		}
	}

	return true, nil
}

// UpdateDailyAvailability updates daily availability (admin only)
func (s *calendarService) UpdateDailyAvailability(date time.Time, status string, notes string) error {
	availability, err := s.calendarRepo.GetAvailabilityByDate(date)
	if err != nil {
		return err
	}

	availability.Status = status
	availability.Notes = notes

	// Update available slots based on status
	switch status {
	case models.StatusAvailable:
		availability.AvailableSlots = availability.TotalSlots
	case models.StatusBooked, models.StatusClosed, models.StatusMaintenance:
		availability.AvailableSlots = 0
	case models.StatusPending:
		availability.AvailableSlots = availability.TotalSlots - availability.BookedSlots
	}

	return s.calendarRepo.CreateOrUpdateAvailability(availability)
}

// UpdateTimeSlotStatus updates a specific time slot status (admin only)
func (s *calendarService) UpdateTimeSlotStatus(date time.Time, startTime, endTime, status string) error {
	slots, err := s.calendarRepo.GetTimeSlotsByDate(date)
	if err != nil {
		return err
	}

	// Find and update the specific slot
	for _, slot := range slots {
		if slot.StartTime == startTime && slot.EndTime == endTime {
			slot.Status = status
			return s.calendarRepo.UpdateTimeSlot(&slot)
		}
	}

	return errors.New("time slot not found")
}

// GetMonthlyStats returns statistics for a month
func (s *calendarService) GetMonthlyStats(year int, month time.Month) (*models.CalendarStatsResponse, error) {
	statsMap, err := s.calendarRepo.GetMonthlyStats(year, month)
	if err != nil {
		return nil, err
	}

	// Calculate occupancy rate
	totalSlots := statsMap["total_slots"].(int)
	bookedSlots := statsMap["booked_slots"].(int)
	var occupancyRate float64
	if totalSlots > 0 {
		occupancyRate = float64(bookedSlots) / float64(totalSlots) * 100
	}

	return &models.CalendarStatsResponse{
		TotalDays:       statsMap["total_days"].(int),
		AvailableDays:   statsMap["available_days"].(int),
		BookedDays:      statsMap["booked_days"].(int),
		PendingDays:     statsMap["pending_days"].(int),
		ClosedDays:      statsMap["closed_days"].(int),
		MaintenanceDays: statsMap["maintenance_days"].(int),
		TotalBookings:   statsMap["total_bookings"].(int),
		Revenue:         statsMap["revenue"].(float64),
		OccupancyRate:   occupancyRate,
	}, nil
}

// GenerateCalendar generates calendar for a month (admin only)
func (s *calendarService) GenerateCalendar(year int, month time.Month) error {
	return s.calendarRepo.GenerateMonthlyCalendar(year, month)
}

// BulkUpdateAvailability updates multiple calendar entries (admin only)
func (s *calendarService) BulkUpdateAvailability(updates []models.CalendarUpdateRequest) error {
	for _, update := range updates {
		err := s.UpdateDailyAvailability(update.Date, update.Status, update.Notes)
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper methods

// timeSlotConflicts checks if two time slots conflict
func (s *calendarService) timeSlotConflicts(start1, end1, start2, end2 string) bool {
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
