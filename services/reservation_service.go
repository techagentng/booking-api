package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"strings"
	"time"
)

// ReservationService defines reservation business logic
type ReservationService interface {
	CreateReservation(req *models.CreateReservationRequest) (*models.Reservation, error)
	GetReservationByID(id uint) (*models.ReservationDetails, error)
	GetAllReservations(params db.ReservationQueryParams) (*models.ReservationListResponse, error)
	UpdateReservation(id uint, req *models.UpdateReservationRequest) (*models.Reservation, error)
	UpdateReservationStatus(id uint, status string) error
	CancelReservation(id uint, req *models.CancelReservationRequest) error
	CheckIn(id uint, req *models.CheckInRequest) (*models.Reservation, error)
	CheckOut(id uint, req *models.CheckOutRequest) (*models.Reservation, error)
	GetReservationStats(fromDate, toDate string) (*models.ReservationStats, error)
	GetCheckInStats(date string) (*models.CheckInStats, error)
	GetPaymentDetails(id uint) (*models.PaymentDetails, error)
	GetReservationsByGuestID(guestID uint) ([]models.ReservationListItem, error)
	GetReservationsByRoomID(roomID uint) ([]models.ReservationListItem, error)
}

// reservationService implements ReservationService
type reservationService struct {
	reservationRepo db.ReservationRepository
	roomRepo        db.RoomRepository
	guestRepo       db.GuestRepository
}

// NewReservationService creates a new reservation service
func NewReservationService(reservationRepo db.ReservationRepository, roomRepo db.RoomRepository, guestRepo db.GuestRepository) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		roomRepo:        roomRepo,
		guestRepo:       guestRepo,
	}
}

// CreateReservation creates a new reservation with validation
func (s *reservationService) CreateReservation(req *models.CreateReservationRequest) (*models.Reservation, error) {
	// Parse dates
	checkIn, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return nil, errors.New("invalid check_in_date format, use YYYY-MM-DD")
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return nil, errors.New("invalid check_out_date format, use YYYY-MM-DD")
	}

	// Validate date range
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		return nil, errors.New("check_out_date must be after check_in_date")
	}

	// Check if check-in is in the past
	today := time.Now().Truncate(24 * time.Hour)
	if checkIn.Before(today) {
		return nil, errors.New("check_in_date cannot be in the past")
	}

	// Get typed values
	guestID := req.GetGuestID()
	roomID := req.GetRoomID()
	numberOfGuests := req.GetNumberOfGuests()

	var guest *models.Guest

	// Check if creating new guest or using existing
	if req.IsNewGuest() {
		// Validate required guest fields
		if req.GuestName == "" {
			return nil, errors.New("guest_name is required for new guest")
		}
		if req.GuestEmail == "" {
			return nil, errors.New("guest_email is required for new guest")
		}
		if req.GuestPhone == "" {
			return nil, errors.New("guest_phone is required for new guest")
		}

		// Check if guest with email already exists
		existingGuest, err := s.guestRepo.GetGuestByEmail(req.GuestEmail)
		if err == nil && existingGuest != nil {
			// Use existing guest (found by email)
			guest = existingGuest
			guestID = guest.ID
		} else {
			// Try to create new guest
			newGuest := &models.Guest{
				Name:        req.GuestName,
				Email:       req.GuestEmail,
				Phone:       req.GuestPhone,
				Nationality: req.GuestNationality,
				IDType:      req.GuestIDType,
				IDNumber:    req.GuestIDNumber,
			}
			createdGuest, err := s.guestRepo.CreateGuest(newGuest)
			if err != nil {
				// Check if it's a duplicate ID number error
				if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
					return nil, errors.New("a guest with this ID number or email already exists")
				}
				return nil, errors.New("failed to create guest: " + err.Error())
			}
			guest = createdGuest
			guestID = guest.ID
		}
	} else {
		// Validate existing guest
		if guestID == 0 {
			return nil, errors.New("guest_id is required or provide guest details")
		}
		var err error
		guest, err = s.guestRepo.GetGuestByID(guestID)
		if err != nil {
			return nil, errors.New("guest not found")
		}
	}

	// Validate room exists
	room, err := s.roomRepo.GetRoomByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// Check room capacity
	if numberOfGuests > room.Capacity {
		return nil, errors.New("number of guests exceeds room capacity")
	}

	// Check room availability
	available, err := s.reservationRepo.IsRoomAvailable(roomID, checkIn, checkOut, 0)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("room is not available for selected dates")
	}

	// Generate confirmation number
	confirmationNumber, err := s.reservationRepo.GenerateConfirmationNumber()
	if err != nil {
		return nil, err
	}

	// Calculate nights and total
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	totalAmount := float64(nights) * room.PricePerNight

	// Create reservation
	reservation := &models.Reservation{
		ConfirmationNumber: confirmationNumber,
		GuestID:            guestID,
		RoomID:             roomID,
		CheckInDate:        checkIn,
		CheckOutDate:       checkOut,
		NumberOfGuests:     numberOfGuests,
		NumberOfNights:     nights,
		PricePerNight:      room.PricePerNight,
		TotalAmount:        totalAmount,
		Status:             "pending",
		PaymentStatus:      "pending",
		PaymentMethod:      req.PaymentMethod,
		SpecialRequests:    req.SpecialRequests,
	}

	createdReservation, err := s.reservationRepo.CreateReservation(reservation)
	if err != nil {
		return nil, err
	}

	// Attach guest and room for response
	createdReservation.Guest = guest
	createdReservation.Room = room

	return createdReservation, nil
}

// GetReservationByID retrieves a reservation by ID with full details
func (s *reservationService) GetReservationByID(id uint) (*models.ReservationDetails, error) {
	if id == 0 {
		return nil, errors.New("invalid reservation ID")
	}

	reservation, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	// Build detailed response
	details := &models.ReservationDetails{
		ID:                 reservation.ID,
		ConfirmationNumber: reservation.ConfirmationNumber,
		CheckInDate:        reservation.CheckInDate.Format("2006-01-02"),
		CheckOutDate:       reservation.CheckOutDate.Format("2006-01-02"),
		NumberOfGuests:     reservation.NumberOfGuests,
		NumberOfNights:     reservation.NumberOfNights,
		PricePerNight:      reservation.PricePerNight,
		TotalAmount:        reservation.TotalAmount,
		Status:             reservation.Status,
		PaymentStatus:      reservation.PaymentStatus,
		PaymentMethod:      reservation.PaymentMethod,
		SpecialRequests:    reservation.SpecialRequests,
		CheckedInAt:        reservation.CheckedInAt,
		CheckedOutAt:       reservation.CheckedOutAt,
		CancelledAt:        reservation.CancelledAt,
		AdditionalCharges:  reservation.AdditionalCharges,
		Notes:              reservation.Notes,
		CreatedAt:          reservation.CreatedAt,
		UpdatedAt:          reservation.UpdatedAt,
	}

	// Add guest info
	if reservation.Guest != nil {
		details.Guest = models.ReservationGuestInfo{
			ID:       reservation.Guest.ID,
			Name:     reservation.Guest.Name,
			Email:    reservation.Guest.Email,
			Phone:    reservation.Guest.Phone,
			IDType:   reservation.Guest.IDType,
			IDNumber: reservation.Guest.IDNumber,
		}
	}

	// Add room info
	if reservation.Room != nil {
		details.Room = models.ReservationRoomInfo{
			ID:            reservation.Room.ID,
			RoomNumber:    reservation.Room.RoomNumber,
			RoomType:      reservation.Room.RoomType,
			Floor:         reservation.Room.Floor,
			BedType:       reservation.Room.BedType,
			PricePerNight: reservation.Room.PricePerNight,
		}
	}

	return details, nil
}

// GetAllReservations retrieves all reservations with pagination
func (s *reservationService) GetAllReservations(params db.ReservationQueryParams) (*models.ReservationListResponse, error) {
	reservations, total, err := s.reservationRepo.GetAllReservations(params)
	if err != nil {
		return nil, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	// Convert to list items
	items := make([]models.ReservationListItem, len(reservations))
	for i, r := range reservations {
		items[i] = s.toListItem(&r)
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.ReservationListResponse{
		Data: items,
		Meta: models.PaginationMeta{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateReservation updates a reservation
func (s *reservationService) UpdateReservation(id uint, req *models.UpdateReservationRequest) (*models.Reservation, error) {
	if id == 0 {
		return nil, errors.New("invalid reservation ID")
	}

	// Get existing reservation
	existing, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	// Cannot update cancelled or checked-out reservations
	if existing.Status == "cancelled" || existing.Status == "checked-out" {
		return nil, errors.New("cannot update cancelled or checked-out reservation")
	}

	updateData := &models.Reservation{}
	needsRecalculation := false

	// Parse and validate dates if provided
	checkIn := existing.CheckInDate
	checkOut := existing.CheckOutDate

	if req.CheckInDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.CheckInDate)
		if err != nil {
			return nil, errors.New("invalid check_in_date format")
		}
		checkIn = parsed
		updateData.CheckInDate = parsed
		needsRecalculation = true
	}

	if req.CheckOutDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.CheckOutDate)
		if err != nil {
			return nil, errors.New("invalid check_out_date format")
		}
		checkOut = parsed
		updateData.CheckOutDate = parsed
		needsRecalculation = true
	}

	// Validate date range
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		return nil, errors.New("check_out_date must be after check_in_date")
	}

	// Check room availability for new dates
	if needsRecalculation {
		available, err := s.reservationRepo.IsRoomAvailable(existing.RoomID, checkIn, checkOut, id)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, errors.New("room is not available for selected dates")
		}

		// Recalculate nights and total
		nights := int(checkOut.Sub(checkIn).Hours() / 24)
		updateData.NumberOfNights = nights
		updateData.TotalAmount = float64(nights) * existing.PricePerNight
	}

	if req.NumberOfGuests != nil {
		// Validate against room capacity
		room, err := s.roomRepo.GetRoomByID(existing.RoomID)
		if err != nil {
			return nil, err
		}
		if *req.NumberOfGuests > room.Capacity {
			return nil, errors.New("number of guests exceeds room capacity")
		}
		updateData.NumberOfGuests = *req.NumberOfGuests
	}

	if req.SpecialRequests != nil {
		updateData.SpecialRequests = *req.SpecialRequests
	}

	if req.PaymentMethod != nil {
		updateData.PaymentMethod = *req.PaymentMethod
	}

	return s.reservationRepo.UpdateReservation(id, updateData)
}

// UpdateReservationStatus updates the status of a reservation
func (s *reservationService) UpdateReservationStatus(id uint, status string) error {
	if id == 0 {
		return errors.New("invalid reservation ID")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":     true,
		"confirmed":   true,
		"checked-in":  true,
		"checked-out": true,
		"cancelled":   true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// Get existing reservation to validate transition
	existing, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return err
	}

	// Validate status transition
	if err := s.validateStatusTransition(existing.Status, status); err != nil {
		return err
	}

	return s.reservationRepo.UpdateReservationStatus(id, status)
}

// CancelReservation cancels a reservation
func (s *reservationService) CancelReservation(id uint, req *models.CancelReservationRequest) error {
	if id == 0 {
		return errors.New("invalid reservation ID")
	}

	// Get existing reservation
	existing, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return err
	}

	// Cannot cancel already cancelled or checked-out reservations
	if existing.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}
	if existing.Status == "checked-out" {
		return errors.New("cannot cancel checked-out reservation")
	}

	reason := ""
	refundAmount := 0.0
	if req != nil {
		reason = req.CancellationReason
		refundAmount = req.RefundAmount
	}

	return s.reservationRepo.CancelReservation(id, reason, refundAmount)
}

// CheckIn performs check-in for a reservation
func (s *reservationService) CheckIn(id uint, req *models.CheckInRequest) (*models.Reservation, error) {
	if id == 0 {
		return nil, errors.New("invalid reservation ID")
	}

	// Get existing reservation
	existing, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	// Validate status - must be confirmed
	if existing.Status != "confirmed" {
		return nil, errors.New("reservation must be confirmed before check-in")
	}

	// Perform check-in
	if err := s.reservationRepo.CheckIn(id, req); err != nil {
		return nil, err
	}

	// Update room status to occupied
	if err := s.roomRepo.UpdateRoomStatus(existing.RoomID, "occupied"); err != nil {
		return nil, err
	}

	return s.reservationRepo.GetReservationByID(id)
}

// CheckOut performs check-out for a reservation
func (s *reservationService) CheckOut(id uint, req *models.CheckOutRequest) (*models.Reservation, error) {
	if id == 0 {
		return nil, errors.New("invalid reservation ID")
	}

	// Get existing reservation
	existing, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	// Validate status - must be checked-in
	if existing.Status != "checked-in" {
		return nil, errors.New("guest must be checked-in before check-out")
	}

	// Parse check-out time
	checkedOutAt := time.Now()
	if req.ActualCheckOutTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.ActualCheckOutTime)
		if err != nil {
			checkedOutAt = time.Now()
		} else {
			checkedOutAt = parsed
		}
	}

	// Perform check-out
	if err := s.reservationRepo.CheckOut(id, checkedOutAt, req.AdditionalCharges, req.Notes); err != nil {
		return nil, err
	}

	// Update room status to cleaning
	if err := s.roomRepo.UpdateRoomStatus(existing.RoomID, "cleaning"); err != nil {
		return nil, err
	}

	return s.reservationRepo.GetReservationByID(id)
}

// GetReservationStats retrieves reservation statistics
func (s *reservationService) GetReservationStats(fromDate, toDate string) (*models.ReservationStats, error) {
	var from, to *time.Time

	if fromDate != "" {
		parsed, err := time.Parse("2006-01-02", fromDate)
		if err == nil {
			from = &parsed
		}
	}

	if toDate != "" {
		parsed, err := time.Parse("2006-01-02", toDate)
		if err == nil {
			to = &parsed
		}
	}

	return s.reservationRepo.GetReservationStats(from, to)
}

// GetReservationsByGuestID retrieves all reservations for a guest
func (s *reservationService) GetReservationsByGuestID(guestID uint) ([]models.ReservationListItem, error) {
	reservations, err := s.reservationRepo.GetReservationsByGuestID(guestID)
	if err != nil {
		return nil, err
	}

	items := make([]models.ReservationListItem, len(reservations))
	for i, r := range reservations {
		items[i] = s.toListItem(&r)
	}

	return items, nil
}

// GetReservationsByRoomID retrieves all reservations for a room
func (s *reservationService) GetReservationsByRoomID(roomID uint) ([]models.ReservationListItem, error) {
	reservations, err := s.reservationRepo.GetReservationsByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	items := make([]models.ReservationListItem, len(reservations))
	for i, r := range reservations {
		items[i] = s.toListItem(&r)
	}

	return items, nil
}

// toListItem converts a Reservation to ReservationListItem
func (s *reservationService) toListItem(r *models.Reservation) models.ReservationListItem {
	item := models.ReservationListItem{
		ID:                 r.ID,
		ConfirmationNumber: r.ConfirmationNumber,
		GuestID:            r.GuestID,
		RoomID:             r.RoomID,
		CheckInDate:        r.CheckInDate.Format("2006-01-02"),
		CheckOutDate:       r.CheckOutDate.Format("2006-01-02"),
		NumberOfGuests:     r.NumberOfGuests,
		NumberOfNights:     r.NumberOfNights,
		PricePerNight:      r.PricePerNight,
		TotalAmount:        r.TotalAmount,
		Status:             r.Status,
		PaymentStatus:      r.PaymentStatus,
		SpecialRequests:    r.SpecialRequests,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}

	if r.Guest != nil {
		item.GuestName = r.Guest.Name
		item.GuestEmail = r.Guest.Email
		item.GuestPhone = r.Guest.Phone
	}

	if r.Room != nil {
		item.RoomNumber = r.Room.RoomNumber
		item.RoomType = r.Room.RoomType
	}

	return item
}

// validateStatusTransition validates if a status transition is allowed
func (s *reservationService) validateStatusTransition(currentStatus, newStatus string) error {
	// Define allowed transitions
	allowedTransitions := map[string][]string{
		"pending":     {"confirmed", "cancelled"},
		"confirmed":   {"checked-in", "cancelled"},
		"checked-in":  {"checked-out"},
		"checked-out": {}, // No transitions allowed
		"cancelled":   {}, // No transitions allowed
	}

	allowed, exists := allowedTransitions[currentStatus]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, s := range allowed {
		if s == newStatus {
			return nil
		}
	}

	return errors.New("invalid status transition from " + currentStatus + " to " + newStatus)
}

// GetCheckInStats retrieves check-in statistics for a specific date
func (s *reservationService) GetCheckInStats(date string) (*models.CheckInStats, error) {
	targetDate := time.Now()
	if date != "" {
		parsed, err := time.Parse("2006-01-02", date)
		if err == nil {
			targetDate = parsed
		}
	}

	return s.reservationRepo.GetCheckInStats(targetDate)
}

// GetPaymentDetails retrieves payment details for a reservation
func (s *reservationService) GetPaymentDetails(id uint) (*models.PaymentDetails, error) {
	if id == 0 {
		return nil, errors.New("invalid reservation ID")
	}

	reservation, err := s.reservationRepo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	// Calculate payment details
	paidAmount := 0.0
	pendingAmount := reservation.TotalAmount

	if reservation.PaymentStatus == "paid" {
		paidAmount = reservation.TotalAmount
		pendingAmount = 0
	} else if reservation.PaymentStatus == "partially_paid" {
		// For partial payments, we'd need a separate payments table
		// For now, assume 50% paid
		paidAmount = reservation.TotalAmount * 0.5
		pendingAmount = reservation.TotalAmount * 0.5
	}

	// Default deposit required is 20% of total or minimum $50
	depositRequired := reservation.TotalAmount * 0.2
	if depositRequired < 50 {
		depositRequired = 50
	}

	paymentDate := ""
	if reservation.PaymentStatus == "paid" {
		paymentDate = reservation.CreatedAt.Format(time.RFC3339)
	}

	return &models.PaymentDetails{
		ReservationID:    reservation.ID,
		TotalAmount:      reservation.TotalAmount,
		PaidAmount:       paidAmount,
		PendingAmount:    pendingAmount,
		PaymentStatus:    reservation.PaymentStatus,
		PaymentMethod:    reservation.PaymentMethod,
		PaymentDate:      paymentDate,
		DepositRequired:  depositRequired,
		DepositCollected: reservation.DepositCollected,
	}, nil
}
