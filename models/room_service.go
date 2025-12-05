package models

import (
	"time"

	"github.com/lib/pq"
)

// MenuItem represents a room service menu item
type MenuItem struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`
	Description     string         `json:"description"`
	Category        string         `json:"category" gorm:"index"`
	Price           float64        `json:"price" gorm:"not null"`
	PreparationTime int            `json:"preparation_time"` // in minutes
	Available       bool           `json:"available" gorm:"default:true"`
	ImageURL        string         `json:"image_url"`
	Allergens       pq.StringArray `json:"allergens" gorm:"type:text[]"`
	DietaryInfo     pq.StringArray `json:"dietary_info" gorm:"type:text[]"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// RoomServiceOrder represents a room service order
type RoomServiceOrder struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	OrderNumber           string     `json:"order_number" gorm:"uniqueIndex;not null"`
	RoomID                uint       `json:"room_id" gorm:"index;not null"`
	GuestID               uint       `json:"guest_id" gorm:"index;not null"`
	Subtotal              float64    `json:"subtotal"`
	Tax                   float64    `json:"tax"`
	ServiceCharge         float64    `json:"service_charge"`
	TotalAmount           float64    `json:"total_amount"`
	Status                string     `json:"status" gorm:"default:pending;index"` // pending, preparing, ready, delivering, delivered, cancelled
	Priority              string     `json:"priority" gorm:"default:normal"`      // normal, high, urgent
	SpecialRequests       string     `json:"special_requests"`
	OrderedAt             time.Time  `json:"ordered_at" gorm:"autoCreateTime"`
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time"`
	ActualDeliveryTime    *time.Time `json:"actual_delivery_time"`
	PreparedBy            string     `json:"prepared_by"`
	DeliveredBy           string     `json:"delivered_by"`
	CancellationReason    string     `json:"cancellation_reason"`
	CancelledBy           string     `json:"cancelled_by"`
	PaymentStatus         string     `json:"payment_status" gorm:"default:pending"` // pending, paid, charged_to_room
	CreatedAt             time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Room  *Room       `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	Guest *Guest      `json:"guest,omitempty" gorm:"foreignKey:GuestID"`
	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderItem represents an item in a room service order
type OrderItem struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	OrderID             uint      `json:"order_id" gorm:"index;not null"`
	MenuItemID          uint      `json:"menu_item_id" gorm:"index;not null"`
	Name                string    `json:"name"`
	Category            string    `json:"category"`
	Quantity            int       `json:"quantity" gorm:"not null"`
	UnitPrice           float64   `json:"unit_price"`
	SpecialInstructions string    `json:"special_instructions"`
	Subtotal            float64   `json:"subtotal"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	MenuItem *MenuItem `json:"menu_item,omitempty" gorm:"foreignKey:MenuItemID"`
}

// CreateMenuItemRequest is the request for creating a menu item
type CreateMenuItemRequest struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	Category        string   `json:"category" binding:"required"`
	Price           float64  `json:"price" binding:"required,gt=0"`
	PreparationTime int      `json:"preparation_time"`
	Available       bool     `json:"available"`
	ImageURL        string   `json:"image_url"`
	Allergens       []string `json:"allergens"`
	DietaryInfo     []string `json:"dietary_info"`
}

// UpdateMenuItemRequest is the request for updating a menu item
type UpdateMenuItemRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Price           float64  `json:"price"`
	PreparationTime int      `json:"preparation_time"`
	Available       *bool    `json:"available"`
	ImageURL        string   `json:"image_url"`
	Allergens       []string `json:"allergens"`
	DietaryInfo     []string `json:"dietary_info"`
}

// CreateOrderItemRequest is the request for an order item
type CreateOrderItemRequest struct {
	MenuItemID          uint   `json:"menu_item_id" binding:"required"`
	Quantity            int    `json:"quantity" binding:"required,gt=0"`
	SpecialInstructions string `json:"special_instructions"`
}

// CreateRoomServiceOrderRequest is the request for creating an order
type CreateRoomServiceOrderRequest struct {
	RoomID                uint                     `json:"room_id" binding:"required"`
	GuestID               uint                     `json:"guest_id" binding:"required"`
	Items                 []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	SpecialRequests       string                   `json:"special_requests"`
	Priority              string                   `json:"priority"`
	EstimatedDeliveryTime string                   `json:"estimated_delivery_time"`
}

// UpdateOrderStatusRequest is the request for updating order status
type UpdateOrderStatusRequest struct {
	Status     string `json:"status" binding:"required"`
	PreparedBy string `json:"prepared_by"`
	Notes      string `json:"notes"`
}

// DeliverOrderRequest is the request for marking order as delivered
type DeliverOrderRequest struct {
	DeliveredBy        string `json:"delivered_by" binding:"required"`
	ActualDeliveryTime string `json:"actual_delivery_time"`
	Notes              string `json:"notes"`
	GuestSignature     string `json:"guest_signature"`
}

// CancelOrderRequest is the request for cancelling an order
type CancelOrderRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
	CancelledBy        string `json:"cancelled_by" binding:"required"`
}

// RoomServiceOrderListItem is a flattened view for list display
type RoomServiceOrderListItem struct {
	ID                    uint        `json:"id"`
	OrderNumber           string      `json:"order_number"`
	RoomNumber            string      `json:"room_number"`
	GuestName             string      `json:"guest_name"`
	GuestID               uint        `json:"guest_id"`
	Items                 []OrderItem `json:"items"`
	Subtotal              float64     `json:"subtotal"`
	Tax                   float64     `json:"tax"`
	ServiceCharge         float64     `json:"service_charge"`
	TotalAmount           float64     `json:"total_amount"`
	Status                string      `json:"status"`
	Priority              string      `json:"priority"`
	SpecialRequests       string      `json:"special_requests"`
	OrderedAt             time.Time   `json:"ordered_at"`
	EstimatedDeliveryTime *time.Time  `json:"estimated_delivery_time"`
	ActualDeliveryTime    *time.Time  `json:"actual_delivery_time"`
	PreparedBy            string      `json:"prepared_by"`
	DeliveredBy           string      `json:"delivered_by"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
}

// RoomServiceOrderDetails is the detailed view of an order
type RoomServiceOrderDetails struct {
	ID                    uint                 `json:"id"`
	OrderNumber           string               `json:"order_number"`
	Room                  RoomServiceRoomInfo  `json:"room"`
	Guest                 RoomServiceGuestInfo `json:"guest"`
	Items                 []OrderItem          `json:"items"`
	Subtotal              float64              `json:"subtotal"`
	Tax                   float64              `json:"tax"`
	ServiceCharge         float64              `json:"service_charge"`
	TotalAmount           float64              `json:"total_amount"`
	Status                string               `json:"status"`
	Priority              string               `json:"priority"`
	SpecialRequests       string               `json:"special_requests"`
	OrderedAt             time.Time            `json:"ordered_at"`
	EstimatedDeliveryTime *time.Time           `json:"estimated_delivery_time"`
	ActualDeliveryTime    *time.Time           `json:"actual_delivery_time"`
	PreparedBy            string               `json:"prepared_by"`
	DeliveredBy           string               `json:"delivered_by"`
	PaymentStatus         string               `json:"payment_status"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

// RoomServiceRoomInfo is room info for order details
type RoomServiceRoomInfo struct {
	ID         uint   `json:"id"`
	RoomNumber string `json:"room_number"`
	RoomType   string `json:"room_type"`
	Floor      int    `json:"floor"`
}

// RoomServiceGuestInfo is guest info for order details
type RoomServiceGuestInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// RoomServiceOrderListResponse is the paginated response for order list
type RoomServiceOrderListResponse struct {
	Data []RoomServiceOrderListItem `json:"data"`
	Meta PaginationMeta             `json:"meta"`
}

// RoomServiceStats holds room service statistics
type RoomServiceStats struct {
	Period              string        `json:"period"`
	Date                string        `json:"date"`
	TotalOrders         int64         `json:"total_orders"`
	PendingOrders       int64         `json:"pending_orders"`
	PreparingOrders     int64         `json:"preparing_orders"`
	ReadyOrders         int64         `json:"ready_orders"`
	DeliveringOrders    int64         `json:"delivering_orders"`
	DeliveredOrders     int64         `json:"delivered_orders"`
	CancelledOrders     int64         `json:"cancelled_orders"`
	TotalRevenue        float64       `json:"total_revenue"`
	AverageDeliveryTime int           `json:"average_delivery_time"` // in minutes
	PopularItems        []PopularItem `json:"popular_items"`
	PeakHours           []PeakHour    `json:"peak_hours"`
}

// PopularItem represents a popular menu item in stats
type PopularItem struct {
	MenuItemID  uint    `json:"menu_item_id"`
	Name        string  `json:"name"`
	OrdersCount int64   `json:"orders_count"`
	Revenue     float64 `json:"revenue"`
}

// PeakHour represents order count by hour
type PeakHour struct {
	Hour        int   `json:"hour"`
	OrdersCount int64 `json:"orders_count"`
}

// ActiveOrderItem is a simplified order for active orders list
type ActiveOrderItem struct {
	ID                    uint       `json:"id"`
	OrderNumber           string     `json:"order_number"`
	RoomNumber            string     `json:"room_number"`
	Status                string     `json:"status"`
	ItemsCount            int        `json:"items_count"`
	TotalAmount           float64    `json:"total_amount"`
	OrderedAt             time.Time  `json:"ordered_at"`
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time"`
	TimeElapsed           int        `json:"time_elapsed"` // minutes since order
}
