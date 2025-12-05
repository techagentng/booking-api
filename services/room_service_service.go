package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"time"

	"github.com/lib/pq"
)

// RoomServiceService defines room service business logic
type RoomServiceService interface {
	// Menu operations
	CreateMenuItem(req *models.CreateMenuItemRequest) (*models.MenuItem, error)
	GetMenuItemByID(id uint) (*models.MenuItem, error)
	GetAllMenuItems(category string, availableOnly bool) ([]models.MenuItem, error)
	UpdateMenuItem(id uint, req *models.UpdateMenuItemRequest) (*models.MenuItem, error)
	DeleteMenuItem(id uint) error

	// Order operations
	CreateOrder(req *models.CreateRoomServiceOrderRequest) (*models.RoomServiceOrder, error)
	GetOrderByID(id uint) (*models.RoomServiceOrderDetails, error)
	GetAllOrders(params db.RoomServiceQueryParams) (*models.RoomServiceOrderListResponse, error)
	GetActiveOrders() ([]models.ActiveOrderItem, error)
	UpdateOrderStatus(id uint, req *models.UpdateOrderStatusRequest) error
	DeliverOrder(id uint, req *models.DeliverOrderRequest) error
	CancelOrder(id uint, req *models.CancelOrderRequest) error
	GetStats(date string, period string) (*models.RoomServiceStats, error)
}

// roomServiceService implements RoomServiceService
type roomServiceService struct {
	roomServiceRepo db.RoomServiceRepository
	roomRepo        db.RoomRepository
	guestRepo       db.GuestRepository
}

// NewRoomServiceService creates a new room service service
func NewRoomServiceService(roomServiceRepo db.RoomServiceRepository, roomRepo db.RoomRepository, guestRepo db.GuestRepository) RoomServiceService {
	return &roomServiceService{
		roomServiceRepo: roomServiceRepo,
		roomRepo:        roomRepo,
		guestRepo:       guestRepo,
	}
}

// Tax and service charge rates
const (
	TaxRate           = 0.10 // 10%
	ServiceChargeRate = 0.05 // 5%
)

// ==================== Menu Operations ====================

// CreateMenuItem creates a new menu item
func (s *roomServiceService) CreateMenuItem(req *models.CreateMenuItemRequest) (*models.MenuItem, error) {
	if req.Name == "" {
		return nil, errors.New("menu item name is required")
	}
	if req.Category == "" {
		return nil, errors.New("menu item category is required")
	}
	if req.Price <= 0 {
		return nil, errors.New("menu item price must be greater than 0")
	}

	item := &models.MenuItem{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Price:           req.Price,
		PreparationTime: req.PreparationTime,
		Available:       req.Available,
		ImageURL:        req.ImageURL,
		Allergens:       pq.StringArray(req.Allergens),
		DietaryInfo:     pq.StringArray(req.DietaryInfo),
	}

	return s.roomServiceRepo.CreateMenuItem(item)
}

// GetMenuItemByID retrieves a menu item by ID
func (s *roomServiceService) GetMenuItemByID(id uint) (*models.MenuItem, error) {
	if id == 0 {
		return nil, errors.New("invalid menu item ID")
	}
	return s.roomServiceRepo.GetMenuItemByID(id)
}

// GetAllMenuItems retrieves all menu items
func (s *roomServiceService) GetAllMenuItems(category string, availableOnly bool) ([]models.MenuItem, error) {
	return s.roomServiceRepo.GetAllMenuItems(category, availableOnly)
}

// UpdateMenuItem updates a menu item
func (s *roomServiceService) UpdateMenuItem(id uint, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	if id == 0 {
		return nil, errors.New("invalid menu item ID")
	}

	// Get existing item
	existing, err := s.roomServiceRepo.GetMenuItemByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}
	if req.PreparationTime > 0 {
		existing.PreparationTime = req.PreparationTime
	}
	if req.Available != nil {
		existing.Available = *req.Available
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	if req.Allergens != nil {
		existing.Allergens = pq.StringArray(req.Allergens)
	}
	if req.DietaryInfo != nil {
		existing.DietaryInfo = pq.StringArray(req.DietaryInfo)
	}

	return s.roomServiceRepo.UpdateMenuItem(id, existing)
}

// DeleteMenuItem deletes a menu item
func (s *roomServiceService) DeleteMenuItem(id uint) error {
	if id == 0 {
		return errors.New("invalid menu item ID")
	}
	return s.roomServiceRepo.DeleteMenuItem(id)
}

// ==================== Order Operations ====================

// CreateOrder creates a new room service order
func (s *roomServiceService) CreateOrder(req *models.CreateRoomServiceOrderRequest) (*models.RoomServiceOrder, error) {
	// Validate room
	room, err := s.roomRepo.GetRoomByID(req.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if room.Status != "occupied" {
		return nil, errors.New("room is not currently occupied")
	}

	// Validate guest
	_, err = s.guestRepo.GetGuestByID(req.GuestID)
	if err != nil {
		return nil, errors.New("guest not found")
	}

	// Validate items
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Generate order number
	orderNumber, err := s.roomServiceRepo.GenerateOrderNumber()
	if err != nil {
		return nil, err
	}

	// Build order items and calculate totals
	var orderItems []models.OrderItem
	var subtotal float64

	for _, itemReq := range req.Items {
		menuItem, err := s.roomServiceRepo.GetMenuItemByID(itemReq.MenuItemID)
		if err != nil {
			return nil, errors.New("menu item not found: " + err.Error())
		}
		if !menuItem.Available {
			return nil, errors.New("menu item '" + menuItem.Name + "' is currently unavailable")
		}
		if itemReq.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}

		itemSubtotal := menuItem.Price * float64(itemReq.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			MenuItemID:          menuItem.ID,
			Name:                menuItem.Name,
			Category:            menuItem.Category,
			Quantity:            itemReq.Quantity,
			UnitPrice:           menuItem.Price,
			SpecialInstructions: itemReq.SpecialInstructions,
			Subtotal:            itemSubtotal,
		})
		subtotal += itemSubtotal
	}

	// Calculate tax and service charge
	tax := subtotal * TaxRate
	serviceCharge := subtotal * ServiceChargeRate
	totalAmount := subtotal + tax + serviceCharge

	// Parse estimated delivery time
	var estimatedDelivery *time.Time
	if req.EstimatedDeliveryTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.EstimatedDeliveryTime)
		if err == nil {
			estimatedDelivery = &parsed
		}
	} else {
		// Default: 30 minutes from now
		defaultTime := time.Now().Add(30 * time.Minute)
		estimatedDelivery = &defaultTime
	}

	// Set priority
	priority := "normal"
	if req.Priority != "" {
		if req.Priority == "normal" || req.Priority == "high" || req.Priority == "urgent" {
			priority = req.Priority
		}
	}

	order := &models.RoomServiceOrder{
		OrderNumber:           orderNumber,
		RoomID:                req.RoomID,
		GuestID:               req.GuestID,
		Subtotal:              subtotal,
		Tax:                   tax,
		ServiceCharge:         serviceCharge,
		TotalAmount:           totalAmount,
		Status:                "pending",
		Priority:              priority,
		SpecialRequests:       req.SpecialRequests,
		EstimatedDeliveryTime: estimatedDelivery,
		Items:                 orderItems,
	}

	return s.roomServiceRepo.CreateOrder(order)
}

// GetOrderByID retrieves order details
func (s *roomServiceService) GetOrderByID(id uint) (*models.RoomServiceOrderDetails, error) {
	if id == 0 {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.roomServiceRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	details := &models.RoomServiceOrderDetails{
		ID:                    order.ID,
		OrderNumber:           order.OrderNumber,
		Items:                 order.Items,
		Subtotal:              order.Subtotal,
		Tax:                   order.Tax,
		ServiceCharge:         order.ServiceCharge,
		TotalAmount:           order.TotalAmount,
		Status:                order.Status,
		Priority:              order.Priority,
		SpecialRequests:       order.SpecialRequests,
		OrderedAt:             order.OrderedAt,
		EstimatedDeliveryTime: order.EstimatedDeliveryTime,
		ActualDeliveryTime:    order.ActualDeliveryTime,
		PreparedBy:            order.PreparedBy,
		DeliveredBy:           order.DeliveredBy,
		PaymentStatus:         order.PaymentStatus,
		CreatedAt:             order.CreatedAt,
		UpdatedAt:             order.UpdatedAt,
	}

	if order.Room != nil {
		details.Room = models.RoomServiceRoomInfo{
			ID:         order.Room.ID,
			RoomNumber: order.Room.RoomNumber,
			RoomType:   order.Room.RoomType,
			Floor:      order.Room.Floor,
		}
	}

	if order.Guest != nil {
		details.Guest = models.RoomServiceGuestInfo{
			ID:    order.Guest.ID,
			Name:  order.Guest.Name,
			Phone: order.Guest.Phone,
			Email: order.Guest.Email,
		}
	}

	return details, nil
}

// GetAllOrders retrieves all orders with pagination
func (s *roomServiceService) GetAllOrders(params db.RoomServiceQueryParams) (*models.RoomServiceOrderListResponse, error) {
	orders, total, err := s.roomServiceRepo.GetAllOrders(params)
	if err != nil {
		return nil, err
	}

	var items []models.RoomServiceOrderListItem
	for _, order := range orders {
		item := s.toOrderListItem(&order)
		items = append(items, item)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.RoomServiceOrderListResponse{
		Data: items,
		Meta: models.PaginationMeta{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetActiveOrders retrieves all active orders
func (s *roomServiceService) GetActiveOrders() ([]models.ActiveOrderItem, error) {
	orders, err := s.roomServiceRepo.GetActiveOrders()
	if err != nil {
		return nil, err
	}

	var activeItems []models.ActiveOrderItem
	now := time.Now()

	for _, order := range orders {
		roomNumber := ""
		if order.Room != nil {
			roomNumber = order.Room.RoomNumber
		}

		timeElapsed := int(now.Sub(order.OrderedAt).Minutes())

		activeItems = append(activeItems, models.ActiveOrderItem{
			ID:                    order.ID,
			OrderNumber:           order.OrderNumber,
			RoomNumber:            roomNumber,
			Status:                order.Status,
			ItemsCount:            len(order.Items),
			TotalAmount:           order.TotalAmount,
			OrderedAt:             order.OrderedAt,
			EstimatedDeliveryTime: order.EstimatedDeliveryTime,
			TimeElapsed:           timeElapsed,
		})
	}

	return activeItems, nil
}

// UpdateOrderStatus updates order status
func (s *roomServiceService) UpdateOrderStatus(id uint, req *models.UpdateOrderStatusRequest) error {
	if id == 0 {
		return errors.New("invalid order ID")
	}

	// Get existing order
	order, err := s.roomServiceRepo.GetOrderByID(id)
	if err != nil {
		return err
	}

	// Validate status transition
	if err := s.validateStatusTransition(order.Status, req.Status); err != nil {
		return err
	}

	return s.roomServiceRepo.UpdateOrderStatus(id, req.Status, req.PreparedBy, req.Notes)
}

// DeliverOrder marks an order as delivered
func (s *roomServiceService) DeliverOrder(id uint, req *models.DeliverOrderRequest) error {
	if id == 0 {
		return errors.New("invalid order ID")
	}
	if req.DeliveredBy == "" {
		return errors.New("delivered by is required")
	}

	// Get existing order
	order, err := s.roomServiceRepo.GetOrderByID(id)
	if err != nil {
		return err
	}

	// Must be in delivering status (or ready for quick delivery)
	if order.Status != "delivering" && order.Status != "ready" {
		return errors.New("order must be in 'delivering' or 'ready' status to mark as delivered")
	}

	// Parse delivery time
	deliveryTime := time.Now()
	if req.ActualDeliveryTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.ActualDeliveryTime)
		if err == nil {
			deliveryTime = parsed
		}
	}

	return s.roomServiceRepo.DeliverOrder(id, req.DeliveredBy, deliveryTime, req.Notes)
}

// CancelOrder cancels an order
func (s *roomServiceService) CancelOrder(id uint, req *models.CancelOrderRequest) error {
	if id == 0 {
		return errors.New("invalid order ID")
	}
	if req.CancellationReason == "" {
		return errors.New("cancellation reason is required")
	}
	if req.CancelledBy == "" {
		return errors.New("cancelled by is required")
	}

	// Get existing order
	order, err := s.roomServiceRepo.GetOrderByID(id)
	if err != nil {
		return err
	}

	// Cannot cancel delivered orders
	if order.Status == "delivered" {
		return errors.New("cannot cancel a delivered order")
	}
	if order.Status == "cancelled" {
		return errors.New("order is already cancelled")
	}

	return s.roomServiceRepo.CancelOrder(id, req.CancellationReason, req.CancelledBy)
}

// GetStats retrieves room service statistics
func (s *roomServiceService) GetStats(date string, period string) (*models.RoomServiceStats, error) {
	targetDate := time.Now()
	if date != "" {
		parsed, err := time.Parse("2006-01-02", date)
		if err == nil {
			targetDate = parsed
		}
	}

	if period == "" {
		period = "today"
	}

	return s.roomServiceRepo.GetOrderStats(targetDate, period)
}

// Helper functions

func (s *roomServiceService) toOrderListItem(order *models.RoomServiceOrder) models.RoomServiceOrderListItem {
	item := models.RoomServiceOrderListItem{
		ID:                    order.ID,
		OrderNumber:           order.OrderNumber,
		GuestID:               order.GuestID,
		Items:                 order.Items,
		Subtotal:              order.Subtotal,
		Tax:                   order.Tax,
		ServiceCharge:         order.ServiceCharge,
		TotalAmount:           order.TotalAmount,
		Status:                order.Status,
		Priority:              order.Priority,
		SpecialRequests:       order.SpecialRequests,
		OrderedAt:             order.OrderedAt,
		EstimatedDeliveryTime: order.EstimatedDeliveryTime,
		ActualDeliveryTime:    order.ActualDeliveryTime,
		PreparedBy:            order.PreparedBy,
		DeliveredBy:           order.DeliveredBy,
		CreatedAt:             order.CreatedAt,
		UpdatedAt:             order.UpdatedAt,
	}

	if order.Room != nil {
		item.RoomNumber = order.Room.RoomNumber
	}
	if order.Guest != nil {
		item.GuestName = order.Guest.Name
	}

	return item
}

func (s *roomServiceService) validateStatusTransition(currentStatus, newStatus string) error {
	// Define allowed transitions
	allowedTransitions := map[string][]string{
		"pending":    {"preparing", "cancelled"},
		"preparing":  {"ready", "cancelled"},
		"ready":      {"delivering", "delivered", "cancelled"},
		"delivering": {"delivered", "cancelled"},
		"delivered":  {}, // No transitions allowed
		"cancelled":  {}, // No transitions allowed
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

	return errors.New("invalid status transition from '" + currentStatus + "' to '" + newStatus + "'")
}
