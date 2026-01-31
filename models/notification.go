package models

import (
	"encoding/json"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeSignIn           NotificationType = "sign_in"
	NotificationTypeSignUp           NotificationType = "sign_up"
	NotificationTypeNewGuest         NotificationType = "new_guest"
	NotificationTypeNewBooking       NotificationType = "new_booking"
	NotificationTypeServiceRequest   NotificationType = "service_request"
	NotificationTypeServiceCompleted NotificationType = "service_completed"
	NotificationTypeRoomServiceOrder NotificationType = "room_service_order"
	NotificationTypeCheckIn          NotificationType = "check_in"
	NotificationTypeCheckOut         NotificationType = "check_out"
	// NEW: Hall booking notification types
	NotificationTypeNewHallBooking       NotificationType = "new_hall_booking"
	NotificationTypeHallBookingUpdated   NotificationType = "hall_booking_updated"
	NotificationTypeHallBookingCancelled NotificationType = "hall_booking_cancelled"
)

// NotificationPriority represents the priority level
type NotificationPriority string

const (
	PriorityLow    NotificationPriority = "low"
	PriorityNormal NotificationPriority = "normal"
	PriorityHigh   NotificationPriority = "high"
	PriorityUrgent NotificationPriority = "urgent"
)

// Notification represents a push notification event
type Notification struct {
	ID        string               `json:"id"`
	Type      NotificationType     `json:"type"`
	Title     string               `json:"title"`
	Message   string               `json:"message"`
	Priority  NotificationPriority `json:"priority"`
	Data      interface{}          `json:"data,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
}

// NewNotification creates a new notification with a unique ID
func NewNotification(notifType NotificationType, title, message string, priority NotificationPriority, data interface{}) *Notification {
	return &Notification{
		ID:        generateNotificationID(),
		Type:      notifType,
		Title:     title,
		Message:   message,
		Priority:  priority,
		Data:      data,
		CreatedAt: time.Now(),
	}
}

// ToJSON converts the notification to JSON bytes
func (n *Notification) ToJSON() ([]byte, error) {
	return json.Marshal(n)
}

// generateNotificationID generates a unique notification ID
func generateNotificationID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of given length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

// SignInData contains data for sign-in notifications
type SignInData struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// SignUpData contains data for sign-up notifications
type SignUpData struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// NewGuestData contains data for new guest notifications
type NewGuestData struct {
	GuestID   uint   `json:"guest_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// NewBookingData contains data for new booking notifications
type NewBookingData struct {
	ReservationID uint      `json:"reservation_id"`
	GuestName     string    `json:"guest_name"`
	RoomNumber    string    `json:"room_number"`
	RoomType      string    `json:"room_type"`
	CheckInDate   time.Time `json:"check_in_date"`
	CheckOutDate  time.Time `json:"check_out_date"`
	TotalAmount   float64   `json:"total_amount"`
}

// ServiceRequestData contains data for service request notifications
type ServiceRequestData struct {
	RequestID   uint   `json:"request_id"`
	ServiceType string `json:"service_type"`
	RoomNumber  string `json:"room_number"`
	GuestName   string `json:"guest_name"`
	Priority    string `json:"priority"`
	Description string `json:"description"`
}

// ServiceCompletedData contains data for service completed notifications
type ServiceCompletedData struct {
	RequestID   uint   `json:"request_id"`
	ServiceType string `json:"service_type"`
	RoomNumber  string `json:"room_number"`
	CompletedBy string `json:"completed_by"`
}

// RoomServiceOrderData contains data for room service order notifications
type RoomServiceOrderData struct {
	OrderID     uint    `json:"order_id"`
	OrderNumber string  `json:"order_number"`
	RoomNumber  string  `json:"room_number"`
	GuestName   string  `json:"guest_name"`
	TotalAmount float64 `json:"total_amount"`
	ItemCount   int     `json:"item_count"`
}

// CheckInData contains data for check-in notifications
type CheckInData struct {
	ReservationID uint   `json:"reservation_id"`
	GuestName     string `json:"guest_name"`
	RoomNumber    string `json:"room_number"`
}

// CheckOutData contains data for check-out notifications
type CheckOutData struct {
	ReservationID uint    `json:"reservation_id"`
	GuestName     string  `json:"guest_name"`
	RoomNumber    string  `json:"room_number"`
	TotalBill     float64 `json:"total_bill"`
}

// NEW: Hall booking notification data structures

// HallBookingData contains data for hall booking notifications
type HallBookingData struct {
	BookingID     uint    `json:"booking_id"`
	BookingIDStr  string  `json:"booking_id_str"`
	OrganizerName string  `json:"organizer_name"`
	EventType     string  `json:"event_type"`
	BookingDate   string  `json:"booking_date"`
	GuestCount    int     `json:"guest_count"`
	TotalPrice    float64 `json:"total_price"`
	Status        string  `json:"status,omitempty"`
	CreatedByType string  `json:"created_by_type,omitempty"`
}
