package db

import (
	"fmt"
	"time"

	"hotel/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}

// CreateServiceBooking creates a new service booking
func (r *BookingRepository) CreateServiceBooking(customerID, providerID, serviceID uuid.UUID, req *models.CreateServiceBookingRequest) (*models.ServiceBooking, error) {
	// Parse dates if provided
	var checkInDate, checkOutDate *time.Time
	if req.CheckInDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.CheckInDate)
		if err == nil {
			checkInDate = &parsed
		}
	}
	if req.CheckOutDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.CheckOutDate)
		if err == nil {
			checkOutDate = &parsed
		}
	}

	booking := models.ServiceBooking{
		ID:              uuid.New(),
		CustomerID:      customerID,
		ProviderID:      providerID,
		ServiceID:       serviceID,
		BookingDate:     time.Now(),
		CheckInDate:     checkInDate,
		CheckOutDate:    checkOutDate,
		GuestCount:      req.GuestCount,
		SpecialRequests: req.SpecialRequests,
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := r.DB.Create(&booking).Error; err != nil {
		return nil, fmt.Errorf("failed to create service booking: %w", err)
	}

	return &booking, nil
}

// GetCustomerBookings gets bookings for a customer
func (r *BookingRepository) GetCustomerBookings(customerID uuid.UUID, status string) ([]models.ServiceBooking, error) {
	var bookings []models.ServiceBooking
	query := r.DB.Order("created_at DESC")
	
	query = query.Where("customer_id = ?", customerID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer bookings: %w", err)
	}

	return bookings, nil
}

// GetProviderBookings gets bookings for a provider
func (r *BookingRepository) GetProviderBookings(providerID uuid.UUID, status string) ([]models.ServiceBooking, error) {
	var bookings []models.ServiceBooking
	query := r.DB.Order("created_at DESC")
	
	query = query.Where("provider_id = ?", providerID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get provider bookings: %w", err)
	}

	return bookings, nil
}

// GetBookingByID gets a booking by ID
func (r *BookingRepository) GetBookingByID(bookingID uuid.UUID) (*models.ServiceBooking, error) {
	var booking models.ServiceBooking
	if err := r.DB.Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}
	return &booking, nil
}

// UpdateBookingStatus updates booking status
func (r *BookingRepository) UpdateBookingStatus(bookingID uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if err := r.DB.Model(&models.ServiceBooking{}).Where("id = ?", bookingID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	return nil
}

// CancelBooking cancels a booking
func (r *BookingRepository) CancelBooking(bookingID uuid.UUID) error {
	return r.UpdateBookingStatus(bookingID, "cancelled")
}

// AcceptBooking accepts a booking (provider action)
func (r *BookingRepository) AcceptBooking(bookingID uuid.UUID) error {
	return r.UpdateBookingStatus(bookingID, "confirmed")
}

// RejectBooking rejects a booking (provider action)
func (r *BookingRepository) RejectBooking(bookingID uuid.UUID) error {
	return r.UpdateBookingStatus(bookingID, "rejected")
}

// CompleteBooking completes a booking (provider action)
func (r *BookingRepository) CompleteBooking(bookingID uuid.UUID) error {
	return r.UpdateBookingStatus(bookingID, "completed")
}
