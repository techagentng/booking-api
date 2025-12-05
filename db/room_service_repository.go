package db

import (
	"errors"
	"fmt"
	"hotel/models"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// RoomServiceRepository defines room service database operations
type RoomServiceRepository interface {
	// Menu Item operations
	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	GetMenuItemByID(id uint) (*models.MenuItem, error)
	GetAllMenuItems(category string, availableOnly bool) ([]models.MenuItem, error)
	UpdateMenuItem(id uint, item *models.MenuItem) (*models.MenuItem, error)
	DeleteMenuItem(id uint) error

	// Order operations
	CreateOrder(order *models.RoomServiceOrder) (*models.RoomServiceOrder, error)
	GetOrderByID(id uint) (*models.RoomServiceOrder, error)
	GetAllOrders(params RoomServiceQueryParams) ([]models.RoomServiceOrder, int64, error)
	GetActiveOrders() ([]models.RoomServiceOrder, error)
	UpdateOrderStatus(id uint, status, preparedBy, notes string) error
	DeliverOrder(id uint, deliveredBy string, deliveryTime time.Time, notes string) error
	CancelOrder(id uint, reason, cancelledBy string) error
	GetOrderStats(date time.Time, period string) (*models.RoomServiceStats, error)
	GenerateOrderNumber() (string, error)
}

// RoomServiceQueryParams holds query parameters for listing orders
type RoomServiceQueryParams struct {
	Page       int
	PageSize   int
	Status     string
	RoomNumber string
	Date       *time.Time
	Search     string
}

// roomServiceRepository implements RoomServiceRepository
type roomServiceRepository struct {
	db *gorm.DB
}

// NewRoomServiceRepository creates a new room service repository
func NewRoomServiceRepository(db *gorm.DB) RoomServiceRepository {
	return &roomServiceRepository{db: db}
}

// ==================== Menu Item Operations ====================

// CreateMenuItem creates a new menu item
func (r *roomServiceRepository) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// GetMenuItemByID retrieves a menu item by ID
func (r *roomServiceRepository) GetMenuItemByID(id uint) (*models.MenuItem, error) {
	var item models.MenuItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetAllMenuItems retrieves all menu items with optional filters
func (r *roomServiceRepository) GetAllMenuItems(category string, availableOnly bool) ([]models.MenuItem, error) {
	var items []models.MenuItem
	query := r.db.Model(&models.MenuItem{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if availableOnly {
		query = query.Where("available = ?", true)
	}

	if err := query.Order("category, name").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateMenuItem updates a menu item
func (r *roomServiceRepository) UpdateMenuItem(id uint, item *models.MenuItem) (*models.MenuItem, error) {
	if err := r.db.Model(&models.MenuItem{}).Where("id = ?", id).Updates(item).Error; err != nil {
		return nil, err
	}
	return r.GetMenuItemByID(id)
}

// DeleteMenuItem deletes a menu item
func (r *roomServiceRepository) DeleteMenuItem(id uint) error {
	return r.db.Delete(&models.MenuItem{}, id).Error
}

// ==================== Order Operations ====================

// CreateOrder creates a new room service order
func (r *roomServiceRepository) CreateOrder(order *models.RoomServiceOrder) (*models.RoomServiceOrder, error) {
	if err := r.db.Create(order).Error; err != nil {
		return nil, err
	}
	return r.GetOrderByID(order.ID)
}

// GetOrderByID retrieves an order by ID with all relations
func (r *roomServiceRepository) GetOrderByID(id uint) (*models.RoomServiceOrder, error) {
	var order models.RoomServiceOrder
	if err := r.db.
		Preload("Room").
		Preload("Guest").
		Preload("Items").
		First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAllOrders retrieves all orders with filters and pagination
func (r *roomServiceRepository) GetAllOrders(params RoomServiceQueryParams) ([]models.RoomServiceOrder, int64, error) {
	var orders []models.RoomServiceOrder
	var total int64

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	offset := (params.Page - 1) * params.PageSize

	query := r.db.Model(&models.RoomServiceOrder{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.RoomNumber != "" {
		query = query.Joins("JOIN rooms ON rooms.id = room_service_orders.room_id").
			Where("rooms.room_number = ?", params.RoomNumber)
	}
	if params.Date != nil {
		startOfDay := time.Date(params.Date.Year(), params.Date.Month(), params.Date.Day(), 0, 0, 0, 0, params.Date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("ordered_at >= ? AND ordered_at < ?", startOfDay, endOfDay)
	}
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("order_number ILIKE ?", searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and relations
	if err := query.
		Preload("Room").
		Preload("Guest").
		Preload("Items").
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetActiveOrders retrieves all non-delivered/non-cancelled orders
func (r *roomServiceRepository) GetActiveOrders() ([]models.RoomServiceOrder, error) {
	var orders []models.RoomServiceOrder
	if err := r.db.
		Preload("Room").
		Preload("Guest").
		Preload("Items").
		Where("status NOT IN ?", []string{"delivered", "cancelled"}).
		Order("CASE WHEN priority = 'urgent' THEN 1 WHEN priority = 'high' THEN 2 ELSE 3 END, created_at ASC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrderStatus updates the status of an order
func (r *roomServiceRepository) UpdateOrderStatus(id uint, status, preparedBy, notes string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if preparedBy != "" {
		updates["prepared_by"] = preparedBy
	}
	// Notes could be stored in a separate field or log
	return r.db.Model(&models.RoomServiceOrder{}).Where("id = ?", id).Updates(updates).Error
}

// DeliverOrder marks an order as delivered
func (r *roomServiceRepository) DeliverOrder(id uint, deliveredBy string, deliveryTime time.Time, notes string) error {
	return r.db.Model(&models.RoomServiceOrder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":               "delivered",
		"delivered_by":         deliveredBy,
		"actual_delivery_time": deliveryTime,
		"payment_status":       "charged_to_room",
	}).Error
}

// CancelOrder cancels an order
func (r *roomServiceRepository) CancelOrder(id uint, reason, cancelledBy string) error {
	return r.db.Model(&models.RoomServiceOrder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":              "cancelled",
		"cancellation_reason": reason,
		"cancelled_by":        cancelledBy,
	}).Error
}

// GetOrderStats retrieves room service statistics
func (r *roomServiceRepository) GetOrderStats(date time.Time, period string) (*models.RoomServiceStats, error) {
	stats := &models.RoomServiceStats{
		Period: period,
		Date:   date.Format("2006-01-02"),
	}

	var startDate, endDate time.Time
	startDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	switch period {
	case "week":
		// Go back to start of week (Sunday)
		weekday := int(date.Weekday())
		startDate = startDate.AddDate(0, 0, -weekday)
		endDate = startDate.AddDate(0, 0, 7)
	case "month":
		startDate = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		endDate = startDate.AddDate(0, 1, 0)
	default: // today
		endDate = startDate.Add(24 * time.Hour)
	}

	baseQuery := r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate)

	// Total orders
	baseQuery.Count(&stats.TotalOrders)

	// Count by status
	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "pending").Count(&stats.PendingOrders)

	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "preparing").Count(&stats.PreparingOrders)

	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "ready").Count(&stats.ReadyOrders)

	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "delivering").Count(&stats.DeliveringOrders)

	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "delivered").Count(&stats.DeliveredOrders)

	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "cancelled").Count(&stats.CancelledOrders)

	// Total revenue (from delivered orders)
	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "delivered").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	// Average delivery time (in minutes)
	var avgMinutes float64
	r.db.Model(&models.RoomServiceOrder{}).
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Where("status = ?", "delivered").
		Where("actual_delivery_time IS NOT NULL").
		Select("COALESCE(AVG(EXTRACT(EPOCH FROM (actual_delivery_time - ordered_at)) / 60), 0)").
		Scan(&avgMinutes)
	stats.AverageDeliveryTime = int(avgMinutes)

	// Popular items
	var popularItems []struct {
		MenuItemID  uint
		Name        string
		OrdersCount int64
		Revenue     float64
	}
	r.db.Table("order_items").
		Select("order_items.menu_item_id, order_items.name, COUNT(*) as orders_count, SUM(order_items.subtotal) as revenue").
		Joins("JOIN room_service_orders ON room_service_orders.id = order_items.order_id").
		Where("room_service_orders.ordered_at >= ? AND room_service_orders.ordered_at < ?", startDate, endDate).
		Where("room_service_orders.status != ?", "cancelled").
		Group("order_items.menu_item_id, order_items.name").
		Order("orders_count DESC").
		Limit(5).
		Scan(&popularItems)

	for _, item := range popularItems {
		stats.PopularItems = append(stats.PopularItems, models.PopularItem{
			MenuItemID:  item.MenuItemID,
			Name:        item.Name,
			OrdersCount: item.OrdersCount,
			Revenue:     item.Revenue,
		})
	}

	// Peak hours
	var peakHours []struct {
		Hour        int
		OrdersCount int64
	}
	r.db.Model(&models.RoomServiceOrder{}).
		Select("EXTRACT(HOUR FROM ordered_at) as hour, COUNT(*) as orders_count").
		Where("ordered_at >= ? AND ordered_at < ?", startDate, endDate).
		Group("hour").
		Order("orders_count DESC").
		Limit(5).
		Scan(&peakHours)

	for _, ph := range peakHours {
		stats.PeakHours = append(stats.PeakHours, models.PeakHour{
			Hour:        ph.Hour,
			OrdersCount: ph.OrdersCount,
		})
	}

	return stats, nil
}

// GenerateOrderNumber generates a unique order number
func (r *roomServiceRepository) GenerateOrderNumber() (string, error) {
	year := time.Now().Year()

	// Get the count of orders this year
	var count int64
	r.db.Model(&models.RoomServiceOrder{}).
		Where("created_at >= ?", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)).
		Count(&count)

	// Generate order number: RS-YYYY-NNN
	orderNumber := fmt.Sprintf("RS-%d-%03d", year, count+1)

	// Check if it exists (edge case)
	var existing models.RoomServiceOrder
	for {
		err := r.db.Where("order_number = ?", orderNumber).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		count++
		orderNumber = fmt.Sprintf("RS-%d-%03d", year, count+1)
	}

	return orderNumber, nil
}

// Ensure pq is imported for array types
var _ = pq.StringArray{}
