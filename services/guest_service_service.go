package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"time"
)

// GuestServiceService defines guest service business logic
type GuestServiceService interface {
	// Service requests
	CreateHousekeepingRequest(req *models.HousekeepingRequestInput) (*models.GuestServiceRequestResponse, error)
	CreateMaintenanceRequest(req *models.MaintenanceRequestInput) (*models.GuestServiceRequestResponse, error)
	GetServiceRequestsByGuestID(guestID uint, status string) ([]models.GuestServiceListItem, error)
	GetServiceRequestByID(id uint) (*models.GuestServiceRequest, error)
	UpdateServiceRequestStatus(id uint, status, assignedTo string) error
	CompleteServiceRequest(id uint, completedBy string) error
	CancelServiceRequest(id uint, reason string) error

	// Guest room info
	GetGuestByRoomNumber(roomNumber string) (*models.GuestRoomInfo, error)

	// Menu
	GetMenuCategories() ([]models.MenuCategory, error)

	// Guest orders
	GetGuestActiveOrders(guestID uint) ([]models.GuestActiveOrder, error)
	GetOrderStatus(orderID uint) (*models.OrderStatusResponse, error)

	// Hotel info
	GetHotelInfo() *models.HotelInfo

	// Express checkout
	ExpressCheckout(reservationID uint, req *models.ExpressCheckoutRequest) (*models.ExpressCheckoutResponse, error)
}

// guestServiceService implements GuestServiceService
type guestServiceService struct {
	guestServiceRepo db.GuestServiceRepository
	roomServiceRepo  db.RoomServiceRepository
	reservationRepo  db.ReservationRepository
	roomRepo         db.RoomRepository
	guestRepo        db.GuestRepository
}

// NewGuestServiceService creates a new guest service service
func NewGuestServiceService(
	guestServiceRepo db.GuestServiceRepository,
	roomServiceRepo db.RoomServiceRepository,
	reservationRepo db.ReservationRepository,
	roomRepo db.RoomRepository,
	guestRepo db.GuestRepository,
) GuestServiceService {
	return &guestServiceService{
		guestServiceRepo: guestServiceRepo,
		roomServiceRepo:  roomServiceRepo,
		reservationRepo:  reservationRepo,
		roomRepo:         roomRepo,
		guestRepo:        guestRepo,
	}
}

// Valid service types
var validHousekeepingTypes = map[string]bool{
	"cleaning":  true,
	"towels":    true,
	"amenities": true,
	"bedding":   true,
	"turndown":  true,
}

var validMaintenanceTypes = map[string]bool{
	"air_conditioning": true,
	"plumbing":         true,
	"electrical":       true,
	"appliances":       true,
	"furniture":        true,
	"other":            true,
}

// CreateHousekeepingRequest creates a housekeeping service request
func (s *guestServiceService) CreateHousekeepingRequest(req *models.HousekeepingRequestInput) (*models.GuestServiceRequestResponse, error) {
	// Validate room
	_, err := s.roomRepo.GetRoomByID(req.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// Validate guest
	_, err = s.guestRepo.GetGuestByID(req.GuestID)
	if err != nil {
		return nil, errors.New("guest not found")
	}

	// Validate service type
	if !validHousekeepingTypes[req.ServiceType] {
		return nil, errors.New("invalid service type")
	}

	// Generate request number
	requestNumber, err := s.guestServiceRepo.GenerateRequestNumber("housekeeping")
	if err != nil {
		return nil, err
	}

	// Parse preferred time
	var preferredTime *time.Time
	if req.PreferredTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.PreferredTime)
		if err == nil {
			preferredTime = &parsed
		}
	}

	// Set estimated time (default: 30 minutes from now or preferred time)
	estimatedTime := time.Now().Add(30 * time.Minute)
	if preferredTime != nil && preferredTime.After(time.Now()) {
		estimatedTime = *preferredTime
	}

	// Set priority
	priority := "normal"
	if req.Priority == "high" {
		priority = "high"
	}

	request := &models.GuestServiceRequest{
		RequestNumber: requestNumber,
		RoomID:        req.RoomID,
		GuestID:       req.GuestID,
		Type:          "housekeeping",
		ServiceType:   req.ServiceType,
		Priority:      priority,
		Status:        "pending",
		Notes:         req.Notes,
		PreferredTime: preferredTime,
		EstimatedTime: &estimatedTime,
	}

	created, err := s.guestServiceRepo.CreateServiceRequest(request)
	if err != nil {
		return nil, err
	}

	return &models.GuestServiceRequestResponse{
		ID:            created.ID,
		RequestNumber: created.RequestNumber,
		Type:          created.Type,
		ServiceType:   created.ServiceType,
		Status:        created.Status,
		EstimatedTime: estimatedTime.Format(time.RFC3339),
	}, nil
}

// CreateMaintenanceRequest creates a maintenance service request
func (s *guestServiceService) CreateMaintenanceRequest(req *models.MaintenanceRequestInput) (*models.GuestServiceRequestResponse, error) {
	// Validate room
	_, err := s.roomRepo.GetRoomByID(req.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// Validate guest
	_, err = s.guestRepo.GetGuestByID(req.GuestID)
	if err != nil {
		return nil, errors.New("guest not found")
	}

	// Validate issue type
	if !validMaintenanceTypes[req.IssueType] {
		return nil, errors.New("invalid issue type")
	}

	if req.Description == "" {
		return nil, errors.New("description is required")
	}

	// Generate request number
	requestNumber, err := s.guestServiceRepo.GenerateRequestNumber("maintenance")
	if err != nil {
		return nil, err
	}

	// Parse preferred time
	var preferredTime *time.Time
	if req.PreferredTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.PreferredTime)
		if err == nil {
			preferredTime = &parsed
		}
	}

	// Set estimated response time based on priority
	estimatedTime := time.Now().Add(60 * time.Minute) // Default: 1 hour
	priority := "normal"
	if req.Priority == "high" {
		priority = "high"
		estimatedTime = time.Now().Add(30 * time.Minute)
	} else if req.Priority == "urgent" {
		priority = "urgent"
		estimatedTime = time.Now().Add(15 * time.Minute)
	}

	if preferredTime != nil && preferredTime.After(estimatedTime) {
		estimatedTime = *preferredTime
	}

	request := &models.GuestServiceRequest{
		RequestNumber: requestNumber,
		RoomID:        req.RoomID,
		GuestID:       req.GuestID,
		Type:          "maintenance",
		IssueType:     req.IssueType,
		Priority:      priority,
		Status:        "pending",
		Description:   req.Description,
		PreferredTime: preferredTime,
		EstimatedTime: &estimatedTime,
	}

	created, err := s.guestServiceRepo.CreateServiceRequest(request)
	if err != nil {
		return nil, err
	}

	return &models.GuestServiceRequestResponse{
		ID:                    created.ID,
		RequestNumber:         created.RequestNumber,
		Type:                  created.Type,
		IssueType:             created.IssueType,
		Status:                created.Status,
		Priority:              created.Priority,
		EstimatedResponseTime: estimatedTime.Format(time.RFC3339),
	}, nil
}

// GetServiceRequestsByGuestID retrieves service requests for a guest
func (s *guestServiceService) GetServiceRequestsByGuestID(guestID uint, status string) ([]models.GuestServiceListItem, error) {
	requests, err := s.guestServiceRepo.GetServiceRequestsByGuestID(guestID, status)
	if err != nil {
		return nil, err
	}

	var items []models.GuestServiceListItem
	for _, req := range requests {
		item := models.GuestServiceListItem{
			ID:            req.ID,
			RequestNumber: req.RequestNumber,
			Type:          req.Type,
			ServiceType:   req.ServiceType,
			IssueType:     req.IssueType,
			Status:        req.Status,
			Priority:      req.Priority,
			RequestedAt:   req.RequestedAt.Format(time.RFC3339),
		}
		if req.EstimatedTime != nil {
			item.EstimatedTime = req.EstimatedTime.Format(time.RFC3339)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetServiceRequestByID retrieves a service request by ID
func (s *guestServiceService) GetServiceRequestByID(id uint) (*models.GuestServiceRequest, error) {
	return s.guestServiceRepo.GetServiceRequestByID(id)
}

// UpdateServiceRequestStatus updates the status of a service request
func (s *guestServiceService) UpdateServiceRequestStatus(id uint, status, assignedTo string) error {
	return s.guestServiceRepo.UpdateServiceRequestStatus(id, status, assignedTo)
}

// CompleteServiceRequest marks a service request as completed
func (s *guestServiceService) CompleteServiceRequest(id uint, completedBy string) error {
	return s.guestServiceRepo.CompleteServiceRequest(id, completedBy)
}

// CancelServiceRequest cancels a service request
func (s *guestServiceService) CancelServiceRequest(id uint, reason string) error {
	return s.guestServiceRepo.CancelServiceRequest(id, reason)
}

// GetGuestByRoomNumber retrieves guest info for a room
func (s *guestServiceService) GetGuestByRoomNumber(roomNumber string) (*models.GuestRoomInfo, error) {
	return s.guestServiceRepo.GetGuestByRoomNumber(roomNumber)
}

// GetMenuCategories retrieves menu categories
func (s *guestServiceService) GetMenuCategories() ([]models.MenuCategory, error) {
	return s.guestServiceRepo.GetMenuCategories()
}

// GetGuestActiveOrders retrieves active orders for a guest
func (s *guestServiceService) GetGuestActiveOrders(guestID uint) ([]models.GuestActiveOrder, error) {
	orders, err := s.guestServiceRepo.GetGuestActiveOrders(guestID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var activeOrders []models.GuestActiveOrder

	for _, order := range orders {
		var items []models.GuestActiveOrderItem
		for _, item := range order.Items {
			items = append(items, models.GuestActiveOrderItem{
				Name:     item.Name,
				Quantity: item.Quantity,
			})
		}

		timeElapsed := int(now.Sub(order.OrderedAt).Minutes())
		estimatedDelivery := ""
		if order.EstimatedDeliveryTime != nil {
			estimatedDelivery = order.EstimatedDeliveryTime.Format(time.RFC3339)
		}

		activeOrders = append(activeOrders, models.GuestActiveOrder{
			ID:                    order.ID,
			OrderNumber:           order.OrderNumber,
			Status:                order.Status,
			Items:                 items,
			TotalAmount:           order.TotalAmount,
			OrderedAt:             order.OrderedAt.Format(time.RFC3339),
			EstimatedDeliveryTime: estimatedDelivery,
			TimeElapsed:           timeElapsed,
		})
	}

	return activeOrders, nil
}

// GetOrderStatus retrieves order status for tracking
func (s *guestServiceService) GetOrderStatus(orderID uint) (*models.OrderStatusResponse, error) {
	order, err := s.roomServiceRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	// Build status history (simplified - in production, you'd have a separate status_history table)
	statusHistory := []models.OrderStatusHistory{
		{
			Status:    "pending",
			Timestamp: order.CreatedAt.Format(time.RFC3339),
		},
	}

	// Determine current step based on status
	statusSteps := map[string]int{
		"pending":    1,
		"preparing":  2,
		"ready":      3,
		"delivering": 4,
		"delivered":  5,
	}

	currentStep := statusSteps[order.Status]
	if currentStep == 0 {
		currentStep = 1
	}

	// Add current status to history if not pending
	if order.Status != "pending" {
		statusHistory = append(statusHistory, models.OrderStatusHistory{
			Status:    order.Status,
			Timestamp: order.UpdatedAt.Format(time.RFC3339),
			UpdatedBy: order.PreparedBy,
		})
	}

	estimatedDelivery := ""
	if order.EstimatedDeliveryTime != nil {
		estimatedDelivery = order.EstimatedDeliveryTime.Format(time.RFC3339)
	}

	return &models.OrderStatusResponse{
		ID:                    order.ID,
		OrderNumber:           order.OrderNumber,
		Status:                order.Status,
		StatusHistory:         statusHistory,
		EstimatedDeliveryTime: estimatedDelivery,
		CurrentStep:           currentStep,
		TotalSteps:            5,
	}, nil
}

// GetHotelInfo returns static hotel information
func (s *guestServiceService) GetHotelInfo() *models.HotelInfo {
	return &models.HotelInfo{
		Name:        "Grand Hotel",
		Description: "Luxury accommodation in the heart of the city",
		Facilities: []models.HotelFacility{
			{
				Name:     "Swimming Pool",
				Hours:    "6:00 AM - 10:00 PM",
				Location: "Rooftop",
				Icon:     "pool",
			},
			{
				Name:     "Fitness Center",
				Hours:    "24/7",
				Location: "Ground Floor",
				Icon:     "dumbbell",
			},
			{
				Name:     "Spa",
				Hours:    "9:00 AM - 9:00 PM",
				Location: "2nd Floor",
				Icon:     "spa",
			},
			{
				Name:     "Business Center",
				Hours:    "24/7",
				Location: "Lobby Level",
				Icon:     "briefcase",
			},
		},
		Dining: []models.DiningOption{
			{
				Name:     "Main Restaurant",
				Cuisine:  "International",
				Hours:    "7:00 AM - 11:00 PM",
				Location: "Ground Floor",
			},
			{
				Name:     "Rooftop Bar",
				Type:     "Bar & Lounge",
				Hours:    "5:00 PM - 1:00 AM",
				Location: "Rooftop",
			},
			{
				Name:     "Pool Bar",
				Type:     "Casual Dining",
				Hours:    "10:00 AM - 8:00 PM",
				Location: "Rooftop Pool Area",
			},
		},
		Contact: models.HotelContact{
			Reception:   "ext. 0",
			RoomService: "ext. 1",
			Concierge:   "ext. 2",
			Emergency:   "ext. 911",
		},
		WiFi: models.HotelWiFi{
			Network:  "GrandHotel-Guest",
			Password: "welcome2024",
		},
	}
}

// ExpressCheckout initiates express checkout
func (s *guestServiceService) ExpressCheckout(reservationID uint, req *models.ExpressCheckoutRequest) (*models.ExpressCheckoutResponse, error) {
	// Get reservation
	reservation, err := s.reservationRepo.GetReservationByID(reservationID)
	if err != nil {
		return nil, errors.New("reservation not found")
	}

	// Verify guest
	if reservation.GuestID != req.GuestID {
		return nil, errors.New("guest ID does not match reservation")
	}

	// Must be checked-in
	if reservation.Status != "checked-in" {
		return nil, errors.New("reservation must be checked-in for express checkout")
	}

	// Calculate charges (simplified)
	totalCharges := reservation.TotalAmount
	additionalCharges := reservation.AdditionalCharges
	finalTotal := totalCharges + additionalCharges

	// In production, you would:
	// 1. Sum up all room service orders
	// 2. Add any minibar charges
	// 3. Process payment
	// 4. Send receipt email
	// 5. Update reservation status

	return &models.ExpressCheckoutResponse{
		ReservationID:           reservationID,
		CheckoutStatus:          "pending_review",
		TotalCharges:            totalCharges,
		AdditionalCharges:       additionalCharges,
		FinalTotal:              finalTotal,
		ReceiptSent:             req.EmailReceipt,
		EstimatedProcessingTime: "15 minutes",
	}, nil
}
