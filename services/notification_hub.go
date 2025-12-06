package services

import (
	"encoding/json"
	"hotel/models"
	"log"
	"sync"
)

// Client represents a connected SSE client
type Client struct {
	ID       string
	Channel  chan []byte
	UserID   uint
	UserRole string
}

// NotificationHub manages SSE connections and broadcasts
type NotificationHub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *models.Notification
	mu         sync.RWMutex
}

// NewNotificationHub creates a new notification hub
func NewNotificationHub() *NotificationHub {
	hub := &NotificationHub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *models.Notification, 100),
	}
	go hub.run()
	return hub
}

// run starts the hub's main loop
func (h *NotificationHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %s (User: %d, Role: %s)", client.ID, client.UserID, client.UserRole)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				close(client.Channel)
				delete(h.clients, client.ID)
				log.Printf("Client disconnected: %s", client.ID)
			}
			h.mu.Unlock()

		case notification := <-h.broadcast:
			h.mu.RLock()
			data, err := json.Marshal(notification)
			if err != nil {
				log.Printf("Error marshaling notification: %v", err)
				h.mu.RUnlock()
				continue
			}

			for _, client := range h.clients {
				select {
				case client.Channel <- data:
				default:
					// Client channel is full, skip
					log.Printf("Client %s channel full, skipping notification", client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a new client to the hub
func (h *NotificationHub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *NotificationHub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast sends a notification to all connected clients
func (h *NotificationHub) Broadcast(notification *models.Notification) {
	h.broadcast <- notification
}

// BroadcastToRole sends a notification to clients with a specific role
func (h *NotificationHub) BroadcastToRole(notification *models.Notification, role string) {
	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.UserRole == role {
			select {
			case client.Channel <- data:
			default:
				log.Printf("Client %s channel full, skipping notification", client.ID)
			}
		}
	}
}

// GetConnectedClients returns the number of connected clients
func (h *NotificationHub) GetConnectedClients() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// NotifySignIn sends a sign-in notification
func (h *NotificationHub) NotifySignIn(userID uint, email, firstName, lastName string) {
	notification := models.NewNotification(
		models.NotificationTypeSignIn,
		"User Sign In",
		firstName+" "+lastName+" has signed in",
		models.PriorityLow,
		models.SignInData{
			UserID:    userID,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	h.Broadcast(notification)
}

// NotifySignUp sends a sign-up notification
func (h *NotificationHub) NotifySignUp(userID uint, email, firstName, lastName string) {
	notification := models.NewNotification(
		models.NotificationTypeSignUp,
		"New User Registration",
		"New user registered: "+firstName+" "+lastName,
		models.PriorityNormal,
		models.SignUpData{
			UserID:    userID,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	h.Broadcast(notification)
}

// NotifyNewGuest sends a new guest notification
func (h *NotificationHub) NotifyNewGuest(guestID uint, firstName, lastName, email, phone string) {
	notification := models.NewNotification(
		models.NotificationTypeNewGuest,
		"New Guest Added",
		"New guest: "+firstName+" "+lastName,
		models.PriorityNormal,
		models.NewGuestData{
			GuestID:   guestID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Phone:     phone,
		},
	)
	h.Broadcast(notification)
}

// NotifyNewBooking sends a new booking notification
func (h *NotificationHub) NotifyNewBooking(reservationID uint, guestName, roomNumber, roomType string, checkIn, checkOut string, total float64) {
	notification := models.NewNotification(
		models.NotificationTypeNewBooking,
		"New Reservation",
		"New booking for "+guestName+" in Room "+roomNumber,
		models.PriorityHigh,
		map[string]interface{}{
			"reservation_id": reservationID,
			"guest_name":     guestName,
			"room_number":    roomNumber,
			"room_type":      roomType,
			"check_in_date":  checkIn,
			"check_out_date": checkOut,
			"total_amount":   total,
		},
	)
	h.Broadcast(notification)
}

// NotifyServiceRequest sends a service request notification
func (h *NotificationHub) NotifyServiceRequest(requestID uint, serviceType, roomNumber, guestName, priority, description string) {
	notification := models.NewNotification(
		models.NotificationTypeServiceRequest,
		"New Service Request",
		serviceType+" request from Room "+roomNumber,
		models.NotificationPriority(priority),
		models.ServiceRequestData{
			RequestID:   requestID,
			ServiceType: serviceType,
			RoomNumber:  roomNumber,
			GuestName:   guestName,
			Priority:    priority,
			Description: description,
		},
	)
	h.Broadcast(notification)
}

// NotifyServiceCompleted sends a service completed notification
func (h *NotificationHub) NotifyServiceCompleted(requestID uint, serviceType, roomNumber, completedBy string) {
	notification := models.NewNotification(
		models.NotificationTypeServiceCompleted,
		"Service Completed",
		serviceType+" completed for Room "+roomNumber,
		models.PriorityNormal,
		models.ServiceCompletedData{
			RequestID:   requestID,
			ServiceType: serviceType,
			RoomNumber:  roomNumber,
			CompletedBy: completedBy,
		},
	)
	h.Broadcast(notification)
}

// NotifyRoomServiceOrder sends a room service order notification
func (h *NotificationHub) NotifyRoomServiceOrder(orderID uint, orderNumber, roomNumber, guestName string, totalAmount float64, itemCount int) {
	notification := models.NewNotification(
		models.NotificationTypeRoomServiceOrder,
		"New Room Service Order",
		"Order "+orderNumber+" from Room "+roomNumber,
		models.PriorityHigh,
		models.RoomServiceOrderData{
			OrderID:     orderID,
			OrderNumber: orderNumber,
			RoomNumber:  roomNumber,
			GuestName:   guestName,
			TotalAmount: totalAmount,
			ItemCount:   itemCount,
		},
	)
	h.Broadcast(notification)
}

// NotifyCheckIn sends a check-in notification
func (h *NotificationHub) NotifyCheckIn(reservationID uint, guestName, roomNumber string) {
	notification := models.NewNotification(
		models.NotificationTypeCheckIn,
		"Guest Check-In",
		guestName+" checked in to Room "+roomNumber,
		models.PriorityNormal,
		models.CheckInData{
			ReservationID: reservationID,
			GuestName:     guestName,
			RoomNumber:    roomNumber,
		},
	)
	h.Broadcast(notification)
}

// NotifyCheckOut sends a check-out notification
func (h *NotificationHub) NotifyCheckOut(reservationID uint, guestName, roomNumber string, totalBill float64) {
	notification := models.NewNotification(
		models.NotificationTypeCheckOut,
		"Guest Check-Out",
		guestName+" checked out from Room "+roomNumber,
		models.PriorityNormal,
		models.CheckOutData{
			ReservationID: reservationID,
			GuestName:     guestName,
			RoomNumber:    roomNumber,
			TotalBill:     totalBill,
		},
	)
	h.Broadcast(notification)
}
