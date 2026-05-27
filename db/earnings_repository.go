package db

import (
	"fmt"
	"time"

	"hotel/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EarningsRepository struct {
	DB *gorm.DB
}

func NewEarningsRepository(db *gorm.DB) *EarningsRepository {
	return &EarningsRepository{DB: db}
}

// CreateEarning creates a new earning record
func (r *EarningsRepository) CreateEarning(providerID, bookingID, serviceID uuid.UUID, amount, commission float64) (*models.Earning, error) {
	earning := models.Earning{
		ID:         uuid.New(),
		ProviderID: providerID,
		BookingID:  bookingID,
		ServiceID:  serviceID,
		Amount:     amount,
		Commission: commission,
		NetAmount:  amount - commission,
		Status:     "pending",
		EarnedAt:   time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := r.DB.Create(&earning).Error; err != nil {
		return nil, fmt.Errorf("failed to create earning: %w", err)
	}

	return &earning, nil
}

// GetProviderEarnings gets earnings for a provider with optional date range
func (r *EarningsRepository) GetProviderEarnings(providerID uuid.UUID, startDate, endDate string) ([]models.Earning, error) {
	var earnings []models.Earning
	query := r.DB.Where("provider_id = ?", providerID).Order("earned_at DESC")

	if startDate != "" {
		query = query.Where("earned_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("earned_at <= ?", endDate)
	}

	if err := query.Find(&earnings).Error; err != nil {
		return nil, fmt.Errorf("failed to get provider earnings: %w", err)
	}

	return earnings, nil
}

// GetEarningByID gets a specific earning by ID
func (r *EarningsRepository) GetEarningByID(earningID uuid.UUID) (*models.Earning, error) {
	var earning models.Earning
	if err := r.DB.Where("id = ?", earningID).First(&earning).Error; err != nil {
		return nil, fmt.Errorf("earning not found: %w", err)
	}
	return &earning, nil
}

// CreatePayoutRequest creates a new payout request
func (r *EarningsRepository) CreatePayoutRequest(providerID uuid.UUID, req *models.CreatePayoutRequest) (*models.PayoutRequest, error) {
	payout := models.PayoutRequest{
		ID:          uuid.New(),
		ProviderID:  providerID,
		Amount:      req.Amount,
		BankAccount: req.BankAccount,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.DB.Create(&payout).Error; err != nil {
		return nil, fmt.Errorf("failed to create payout request: %w", err)
	}

	return &payout, nil
}

// GetProviderAnalytics gets analytics data for a provider
func (r *EarningsRepository) GetProviderAnalytics(providerID uuid.UUID) (*models.ProviderAnalytics, error) {
	var analytics models.ProviderAnalytics

	// Get total revenue
	var totalRevenue float64
	r.DB.Model(&models.Earning{}).
		Where("provider_id = ? AND status = ?", providerID, "available").
		Select("COALESCE(SUM(net_amount), 0)").
		Scan(&totalRevenue)
	analytics.TotalRevenue = totalRevenue

	// Get total bookings
	var totalBookings int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ?", providerID).
		Count(&totalBookings)
	analytics.TotalBookings = int(totalBookings)

	// Get total services
	var totalServices int64
	r.DB.Model(&models.ProviderService{}).
		Where("provider_id = ? AND is_active = ?", providerID, true).
		Count(&totalServices)
	analytics.TotalServices = int(totalServices)

	// Get average rating from reviews
	var avgRating float64
	r.DB.Model(&models.Review{}).
		Where("provider_id = ?", providerID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating)
	analytics.AverageRating = avgRating

	// Get completed bookings
	var completedBookings int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ? AND status = ?", providerID, "completed").
		Count(&completedBookings)
	analytics.CompletedBookings = int(completedBookings)

	// Get cancelled bookings
	var cancelledBookings int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ? AND status = ?", providerID, "cancelled").
		Count(&cancelledBookings)
	analytics.CancelledBookings = int(cancelledBookings)

	// Get pending bookings
	var pendingBookings int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ? AND status = ?", providerID, "pending").
		Count(&pendingBookings)
	analytics.PendingRequests = int(pendingBookings)

	// Response time (placeholder - would need actual response time tracking)
	analytics.ResponseTime = 2.5 // 2.5 hours average

	return &analytics, nil
}

// GetDashboardAnalytics gets dashboard-specific analytics for a provider
func (r *EarningsRepository) GetDashboardAnalytics(providerID uuid.UUID) (*models.DashboardAnalytics, error) {
	var analytics models.DashboardAnalytics

	// Get current period revenue (last 30 days)
	var currentRevenue float64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	r.DB.Model(&models.Earning{}).
		Where("provider_id = ? AND status = ? AND earned_at >= ?", providerID, "available", thirtyDaysAgo).
		Select("COALESCE(SUM(net_amount), 0)").
		Scan(&currentRevenue)
	analytics.TotalRevenue = currentRevenue

	// Get previous period revenue (30-60 days ago)
	var previousRevenue float64
	sixtyDaysAgo := time.Now().AddDate(0, 0, -60)
	r.DB.Model(&models.Earning{}).
		Where("provider_id = ? AND status = ? AND earned_at >= ? AND earned_at < ?", providerID, "available", sixtyDaysAgo, thirtyDaysAgo).
		Select("COALESCE(SUM(net_amount), 0)").
		Scan(&previousRevenue)

	// Calculate revenue change
	if previousRevenue > 0 {
		analytics.RevenueChange = ((currentRevenue - previousRevenue) / previousRevenue) * 100
	}

	// Get current period bookings
	var currentBookings int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ? AND created_at >= ?", providerID, thirtyDaysAgo).
		Count(&currentBookings)
	analytics.TotalBookings = int(currentBookings)

	// Get previous period bookings
	var previousBookings int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ? AND created_at >= ? AND created_at < ?", providerID, sixtyDaysAgo, thirtyDaysAgo).
		Count(&previousBookings)

	// Calculate bookings change
	if previousBookings > 0 {
		analytics.BookingsChange = ((float64(currentBookings) - float64(previousBookings)) / float64(previousBookings)) * 100
	}

	// Get active services
	var activeServices int64
	r.DB.Model(&models.ProviderService{}).
		Where("provider_id = ? AND is_active = ? AND is_available = ?", providerID, true, true).
		Count(&activeServices)
	analytics.ActiveServices = int(activeServices)

	// Get average rating
	var avgRating float64
	r.DB.Model(&models.Review{}).
		Where("provider_id = ?", providerID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating)
	analytics.AverageRating = avgRating

	// Rating change (placeholder - would need historical rating data)
	analytics.RatingChange = 0.0

	// Response time
	analytics.ResponseTime = 2.5

	// Pending requests
	var pendingRequests int64
	r.DB.Model(&models.ServiceBooking{}).
		Where("provider_id = ? AND status = ?", providerID, "pending").
		Count(&pendingRequests)
	analytics.PendingRequests = int(pendingRequests)

	return &analytics, nil
}
