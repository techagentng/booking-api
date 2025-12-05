package models

import "time"

// GuestServiceRequest represents a unified service request (housekeeping or maintenance)
type GuestServiceRequest struct {
	ID                 uint       `json:"id" gorm:"primaryKey"`
	RequestNumber      string     `json:"request_number" gorm:"uniqueIndex;not null"`
	RoomID             uint       `json:"room_id" gorm:"index;not null"`
	GuestID            uint       `json:"guest_id" gorm:"index;not null"`
	Type               string     `json:"type" gorm:"index;not null"` // housekeeping, maintenance
	ServiceType        string     `json:"service_type"`               // cleaning, towels, amenities, bedding, turndown
	IssueType          string     `json:"issue_type"`                 // air_conditioning, plumbing, electrical, appliances, furniture, other
	Priority           string     `json:"priority" gorm:"default:normal"`
	Status             string     `json:"status" gorm:"default:pending;index"` // pending, in_progress, completed, cancelled
	Description        string     `json:"description"`
	Notes              string     `json:"notes"`
	PreferredTime      *time.Time `json:"preferred_time"`
	EstimatedTime      *time.Time `json:"estimated_time"`
	AssignedTo         string     `json:"assigned_to"`       // Legacy: staff name
	AssignedStaffID    *uint      `json:"assigned_staff_id"` // FK to Staff
	AssignedAt         *time.Time `json:"assigned_at"`       // When staff was assigned
	CompletedAt        *time.Time `json:"completed_at"`
	CompletedBy        string     `json:"completed_by"`
	CancellationReason string     `json:"cancellation_reason"`
	RequestedAt        time.Time  `json:"requested_at" gorm:"autoCreateTime"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Room          *Room  `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	Guest         *Guest `json:"guest,omitempty" gorm:"foreignKey:GuestID"`
	AssignedStaff *Staff `json:"assigned_staff,omitempty" gorm:"foreignKey:AssignedStaffID"`
}

// HousekeepingRequestInput is the request for housekeeping service
type HousekeepingRequestInput struct {
	RoomID        uint   `json:"room_id" binding:"required"`
	GuestID       uint   `json:"guest_id" binding:"required"`
	ServiceType   string `json:"service_type" binding:"required"` // cleaning, towels, amenities, bedding, turndown
	Priority      string `json:"priority"`
	Notes         string `json:"notes"`
	PreferredTime string `json:"preferred_time"`
}

// MaintenanceRequestInput is the request for maintenance service
type MaintenanceRequestInput struct {
	RoomID        uint   `json:"room_id" binding:"required"`
	GuestID       uint   `json:"guest_id" binding:"required"`
	IssueType     string `json:"issue_type" binding:"required"` // air_conditioning, plumbing, electrical, appliances, furniture, other
	Priority      string `json:"priority"`
	Description   string `json:"description" binding:"required"`
	PreferredTime string `json:"preferred_time"`
}

// GuestServiceRequestResponse is the response for a service request
type GuestServiceRequestResponse struct {
	ID                    uint   `json:"id"`
	RequestNumber         string `json:"request_number"`
	Type                  string `json:"type"`
	ServiceType           string `json:"service_type,omitempty"`
	IssueType             string `json:"issue_type,omitempty"`
	Status                string `json:"status"`
	Priority              string `json:"priority,omitempty"`
	RequestedAt           string `json:"requested_at"`
	EstimatedTime         string `json:"estimated_time,omitempty"`
	EstimatedResponseTime string `json:"estimated_response_time,omitempty"`
}

// GuestServiceListItem is a simplified view for list display
type GuestServiceListItem struct {
	ID            uint   `json:"id"`
	RequestNumber string `json:"request_number"`
	Type          string `json:"type"`
	ServiceType   string `json:"service_type,omitempty"`
	IssueType     string `json:"issue_type,omitempty"`
	Status        string `json:"status"`
	Priority      string `json:"priority,omitempty"`
	RequestedAt   string `json:"requested_at"`
	EstimatedTime string `json:"estimated_time,omitempty"`
}

// GuestRoomInfo is the response for guest room information
type GuestRoomInfo struct {
	Guest       GuestRoomInfoGuest       `json:"guest"`
	Reservation GuestRoomInfoReservation `json:"reservation"`
	Room        GuestRoomInfoRoom        `json:"room"`
	IsCheckedIn bool                     `json:"is_checked_in"`
}

// GuestRoomInfoGuest is guest info for tablet
type GuestRoomInfoGuest struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// GuestRoomInfoReservation is reservation info for tablet
type GuestRoomInfoReservation struct {
	ID                 uint   `json:"id"`
	ConfirmationNumber string `json:"confirmation_number"`
	CheckInDate        string `json:"check_in_date"`
	CheckOutDate       string `json:"check_out_date"`
	NumberOfNights     int    `json:"number_of_nights"`
}

// GuestRoomInfoRoom is room info for tablet
type GuestRoomInfoRoom struct {
	ID         uint   `json:"id"`
	RoomNumber string `json:"room_number"`
	RoomType   string `json:"room_type"`
	Floor      int    `json:"floor"`
}

// MenuCategory represents a menu category with item count
type MenuCategory struct {
	Name       string `json:"name"`
	ItemsCount int64  `json:"items_count"`
	Icon       string `json:"icon"`
}

// OrderStatusResponse is the response for order status tracking
type OrderStatusResponse struct {
	ID                    uint                 `json:"id"`
	OrderNumber           string               `json:"order_number"`
	Status                string               `json:"status"`
	StatusHistory         []OrderStatusHistory `json:"status_history"`
	EstimatedDeliveryTime string               `json:"estimated_delivery_time"`
	CurrentStep           int                  `json:"current_step"`
	TotalSteps            int                  `json:"total_steps"`
}

// OrderStatusHistory is a status change entry
type OrderStatusHistory struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	UpdatedBy string `json:"updated_by,omitempty"`
}

// GuestActiveOrder is a simplified order for guest view
type GuestActiveOrder struct {
	ID                    uint                   `json:"id"`
	OrderNumber           string                 `json:"order_number"`
	Status                string                 `json:"status"`
	Items                 []GuestActiveOrderItem `json:"items"`
	TotalAmount           float64                `json:"total_amount"`
	OrderedAt             string                 `json:"ordered_at"`
	EstimatedDeliveryTime string                 `json:"estimated_delivery_time"`
	TimeElapsed           int                    `json:"time_elapsed"`
}

// GuestActiveOrderItem is an item in guest active order
type GuestActiveOrderItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// HotelInfo represents hotel information for guests
type HotelInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Facilities  []HotelFacility `json:"facilities"`
	Dining      []DiningOption  `json:"dining"`
	Contact     HotelContact    `json:"contact"`
	WiFi        HotelWiFi       `json:"wifi"`
}

// HotelFacility represents a hotel facility
type HotelFacility struct {
	Name     string `json:"name"`
	Hours    string `json:"hours"`
	Location string `json:"location"`
	Icon     string `json:"icon"`
}

// DiningOption represents a dining option
type DiningOption struct {
	Name     string `json:"name"`
	Cuisine  string `json:"cuisine,omitempty"`
	Type     string `json:"type,omitempty"`
	Hours    string `json:"hours"`
	Location string `json:"location"`
}

// HotelContact represents hotel contact info
type HotelContact struct {
	Reception   string `json:"reception"`
	RoomService string `json:"room_service"`
	Concierge   string `json:"concierge"`
	Emergency   string `json:"emergency"`
}

// HotelWiFi represents WiFi credentials
type HotelWiFi struct {
	Network  string `json:"network"`
	Password string `json:"password"`
}

// ExpressCheckoutRequest is the request for express checkout
type ExpressCheckoutRequest struct {
	GuestID      uint                     `json:"guest_id" binding:"required"`
	EmailReceipt bool                     `json:"email_receipt"`
	Feedback     *ExpressCheckoutFeedback `json:"feedback"`
}

// ExpressCheckoutFeedback is optional feedback during checkout
type ExpressCheckoutFeedback struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// ExpressCheckoutResponse is the response for express checkout
type ExpressCheckoutResponse struct {
	ReservationID           uint    `json:"reservation_id"`
	CheckoutStatus          string  `json:"checkout_status"`
	TotalCharges            float64 `json:"total_charges"`
	AdditionalCharges       float64 `json:"additional_charges"`
	FinalTotal              float64 `json:"final_total"`
	ReceiptSent             bool    `json:"receipt_sent"`
	EstimatedProcessingTime string  `json:"estimated_processing_time"`
}
