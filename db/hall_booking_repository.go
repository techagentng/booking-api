package db

import (
	"errors"
	"fmt"
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// HallBookingRepository defines hall booking database operations
type HallBookingRepository interface {
	CreateHallBooking(booking *models.HallBooking) (*models.HallBooking, error)
	GetHallBookingByID(id uint) (*models.HallBooking, error)
	GetHallBookingByBookingID(bookingID string) (*models.HallBooking, error)
	GetAllHallBookings(page, pageSize int, search string) ([]models.HallBooking, int64, error)
	UpdateHallBooking(id uint, booking *models.HallBooking) (*models.HallBooking, error)
	DeleteHallBooking(id uint) error
	CheckHallAvailability(date string, startTime, endTime string) (bool, error)
	GetHallAvailability(date string) (*models.HallAvailability, error)
	GetHallBookingsByDate(date string) ([]models.HallBooking, error)
	GetHallBookingsByDateRange(startDate, endDate string) ([]models.HallBooking, error)
}

// hallBookingRepository implements HallBookingRepository
type hallBookingRepository struct {
	db *gorm.DB
}

// NewHallBookingRepository creates a new hall booking repository
func NewHallBookingRepository(db *gorm.DB) HallBookingRepository {
	return &hallBookingRepository{db: db}
}

// CreateHallBooking creates a new hall booking
func (r *hallBookingRepository) CreateHallBooking(booking *models.HallBooking) (*models.HallBooking, error) {
	// Check if hall is available for the requested time
	available, err := r.CheckHallAvailability(
		booking.BookingDate.Format("2006-01-02"),
		booking.StartTime,
		booking.EndTime,
	)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("hall is not available for the requested time slot")
	}

	// Generate unique booking ID
	booking.BookingID = r.generateBookingID(booking.BookingDate)

	// Create the booking
	if err := r.db.Create(booking).Error; err != nil {
		return nil, err
	}

	// Create payments
	if err := r.createPaymentSchedule(booking); err != nil {
		return nil, err
	}

	// Create invoice
	if err := r.createInvoice(booking); err != nil {
		return nil, err
	}

	return r.GetHallBookingByID(booking.ID)
}

// GetHallBookingByID retrieves a hall booking by ID with all relations
func (r *hallBookingRepository) GetHallBookingByID(id uint) (*models.HallBooking, error) {
	var booking models.HallBooking
	if err := r.db.
		Preload("Payments").
		Preload("Invoice").
		First(&booking, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hall booking not found")
		}
		return nil, err
	}
	return &booking, nil
}

// GetHallBookingByBookingID retrieves a hall booking by booking ID
func (r *hallBookingRepository) GetHallBookingByBookingID(bookingID string) (*models.HallBooking, error) {
	var booking models.HallBooking
	if err := r.db.
		Preload("Payments").
		Preload("Invoice").
		Where("booking_id = ?", bookingID).
		First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hall booking not found")
		}
		return nil, err
	}
	return &booking, nil
}

// GetAllHallBookings retrieves all hall bookings with pagination and search
func (r *hallBookingRepository) GetAllHallBookings(page, pageSize int, search string) ([]models.HallBooking, int64, error) {
	var bookings []models.HallBooking
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.HallBooking{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("organizer_name ILIKE ? OR organizer_email ILIKE ? OR booking_id ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Payments").
		Preload("Invoice").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

// UpdateHallBooking updates a hall booking
func (r *hallBookingRepository) UpdateHallBooking(id uint, booking *models.HallBooking) (*models.HallBooking, error) {
	// If date or time is being updated, check availability
	var existingBooking models.HallBooking
	if err := r.db.First(&existingBooking, id).Error; err != nil {
		return nil, errors.New("hall booking not found")
	}

	// Check if date or time is being changed
	dateChanged := !existingBooking.BookingDate.Equal(booking.BookingDate)
	timeChanged := existingBooking.StartTime != booking.StartTime || existingBooking.EndTime != booking.EndTime

	if dateChanged || timeChanged {
		available, err := r.CheckHallAvailability(
			booking.BookingDate.Format("2006-01-02"),
			booking.StartTime,
			booking.EndTime,
		)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, errors.New("hall is not available for the requested time slot")
		}
	}

	if err := r.db.Model(&models.HallBooking{}).Where("id = ?", id).Updates(booking).Error; err != nil {
		return nil, err
	}
	return r.GetHallBookingByID(id)
}

// DeleteHallBooking deletes a hall booking (soft delete)
func (r *hallBookingRepository) DeleteHallBooking(id uint) error {
	return r.db.Delete(&models.HallBooking{}, id).Error
}

// CheckHallAvailability checks if the hall is available for the given date and time
func (r *hallBookingRepository) CheckHallAvailability(date string, startTime, endTime string) (bool, error) {
	var count int64

	// Parse the date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, errors.New("invalid date format")
	}

	// Find bookings that overlap with the requested time slot
	err = r.db.Model(&models.HallBooking{}).
		Where("booking_date = ? AND status NOT IN ?",
			parsedDate,
			[]string{"cancelled"}).
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?)",
			endTime, startTime, // Overlaps from left
			startTime, endTime, // Overlaps from right
			startTime, endTime). // Completely contains
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetHallAvailability gets availability for a specific date
func (r *hallBookingRepository) GetHallAvailability(date string) (*models.HallAvailability, error) {
	// Validate the date format
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Get all bookings for the date
	bookings, err := r.GetHallBookingsByDate(date)
	if err != nil {
		return nil, err
	}

	// Define all possible time slots (hourly from 08:00 to 23:00)
	var slots []models.TimeSlot
	for hour := 8; hour <= 22; hour++ {
		startTime := time.Time{}.Add(time.Duration(hour) * time.Hour).Format("15:04")
		endTime := time.Time{}.Add(time.Duration(hour+1) * time.Hour).Format("15:04")

		slot := models.TimeSlot{
			StartTime: startTime,
			EndTime:   endTime,
			Available: true,
		}

		// Check if this slot conflicts with any booking
		for _, booking := range bookings {
			if booking.Status != "cancelled" && r.timeSlotConflicts(startTime, endTime, booking.StartTime, booking.EndTime) {
				slot.Available = false
				break
			}
		}

		slots = append(slots, slot)
	}

	return &models.HallAvailability{
		Date:  date,
		Slots: slots,
	}, nil
}

// GetHallBookingsByDate retrieves all hall bookings for a specific date
func (r *hallBookingRepository) GetHallBookingsByDate(date string) ([]models.HallBooking, error) {
	var bookings []models.HallBooking

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	err = r.db.Where("booking_date = ?", parsedDate).
		Order("start_time").
		Find(&bookings).Error

	return bookings, err
}

// GetHallBookingsByDateRange retrieves all hall bookings within a date range
func (r *hallBookingRepository) GetHallBookingsByDateRange(startDate, endDate string) ([]models.HallBooking, error) {
	var bookings []models.HallBooking

	parsedStart, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	parsedEnd, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	err = r.db.Where("booking_date BETWEEN ? AND ?", parsedStart, parsedEnd).
		Order("booking_date, start_time").
		Find(&bookings).Error

	return bookings, err
}

// Helper functions

// generateBookingID generates a unique booking ID
func (r *hallBookingRepository) generateBookingID(date time.Time) string {
	dateStr := date.Format("20060102")

	// Find the count of bookings for this date
	var count int64
	r.db.Model(&models.HallBooking{}).
		Where("booking_date = ?", date).
		Count(&count)

	// Generate booking ID with sequential number
	return fmt.Sprintf("HB-%s-%03d", dateStr, count+1)
}

// createPaymentSchedule creates payment schedule for a hall booking
func (r *hallBookingRepository) createPaymentSchedule(booking *models.HallBooking) error {
	// Deposit payment (due 7 days before event)
	depositDueDate := booking.BookingDate.AddDate(0, 0, -7)
	deposit := &models.Payment{
		HallBookingID: booking.ID,
		PaymentType:   "deposit",
		PaymentMethod: booking.PaymentMethod,
		Amount:        booking.DepositRequired,
		Status:        "pending",
		DueDate:       depositDueDate,
	}

	// Balance payment (due on event day)
	balanceDueDate := booking.BookingDate
	balance := &models.Payment{
		HallBookingID: booking.ID,
		PaymentType:   "balance",
		PaymentMethod: booking.PaymentMethod,
		Amount:        booking.TotalPrice - booking.DepositRequired,
		Status:        "pending",
		DueDate:       balanceDueDate,
	}

	payments := []models.Payment{*deposit, *balance}
	return r.db.Create(&payments).Error
}

// createInvoice creates an invoice for a hall booking
func (r *hallBookingRepository) createInvoice(booking *models.HallBooking) error {
	// Generate invoice number
	invoiceNumber := r.generateInvoiceNumber()

	invoice := &models.Invoice{
		HallBookingID: booking.ID,
		InvoiceNumber: invoiceNumber,
		InvoiceDate:   time.Now(),
		DueDate:       booking.BookingDate.AddDate(0, 0, -7), // Due 7 days before event
		TotalAmount:   booking.TotalPrice,
		Status:        "sent",
	}

	return r.db.Create(invoice).Error
}

// generateInvoiceNumber generates a unique invoice number
func (r *hallBookingRepository) generateInvoiceNumber() string {
	var count int64
	r.db.Model(&models.Invoice{}).Count(&count)
	return fmt.Sprintf("INV-%d", count+1)
}

// timeSlotConflicts checks if two time slots conflict
func (r *hallBookingRepository) timeSlotConflicts(start1, end1, start2, end2 string) bool {
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
