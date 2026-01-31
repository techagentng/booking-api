package models

import (
	"time"
)

// HallBooking represents a hall booking for events
type HallBooking struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	BookingID       string    `json:"booking_id" gorm:"uniqueIndex;not null;type:varchar(50)"`
	OrganizerName   string    `json:"organizer_name" gorm:"not null"`
	OrganizerEmail  string    `json:"organizer_email" gorm:"not null"`
	OrganizerPhone  string    `json:"organizer_phone" gorm:"not null"`
	EventType       string    `json:"event_type" gorm:"not null"`
	GuestCount      int       `json:"guest_count" gorm:"not null"`
	SpecialRequests string    `json:"special_requests"`
	BookingDate     time.Time `json:"booking_date" gorm:"not null"`
	StartTime       string    `json:"start_time" gorm:"not null"`
	EndTime         string    `json:"end_time" gorm:"not null"`
	TotalPrice      float64   `json:"total_price" gorm:"not null"`
	DepositRequired float64   `json:"deposit_required" gorm:"not null"`
	PaymentMethod   string    `json:"payment_method" gorm:"not null"`
	Status          string    `json:"status" gorm:"default:'pending'"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       time.Time `json:"deleted_at" gorm:"index"`

	// Admin tracking fields
	ConfirmedBy *uint      `json:"confirmed_by"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	CancelledBy *uint      `json:"cancelled_by"`
	CancelledAt *time.Time `json:"cancelled_at"`
	UpdatedBy   *uint      `json:"updated_by"`

	// Creator tracking fields
	CreatedBy     *uint  `json:"-"`                                       // User ID (if authenticated)
	CreatedByType string `json:"created_by_type" gorm:"default:'public'"` // "user", "admin", or "public"
	CreatorEmail  string `json:"creator_email"`                           // Email for public users
	CreatorName   string `json:"creator_name"`                            // Name for public users

	// Relations - temporarily disabled for migration
	Payments []Payment `json:"payments,omitempty"`
	Invoice  *Invoice  `json:"invoice,omitempty"`
	// ConfirmedByUser *User                  `json:"confirmed_by_user,omitempty"`
	// CancelledByUser *User                  `json:"cancelled_by_user,omitempty"`
	// UpdatedByUser   *User                  `json:"updated_by_user,omitempty"`
	// StatusHistory   []BookingStatusHistory `json:"status_history,omitempty"`
}

// Payment represents a payment for a hall booking
type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	HallBookingID uint      `json:"hall_booking_id" gorm:"not null"`
	PaymentType   string    `json:"payment_type" gorm:"not null"`   // deposit, balance
	PaymentMethod string    `json:"payment_method" gorm:"not null"` // cash, onsite, online
	Amount        float64   `json:"amount" gorm:"not null"`
	Status        string    `json:"status" gorm:"default:'pending'"` // pending, paid, overdue
	DueDate       time.Time `json:"due_date" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"index"`

	// Relation
	HallBooking *HallBooking `json:"-" gorm:"foreignKey:HallBookingID"`
}

// Invoice represents an invoice for a hall booking
type Invoice struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	HallBookingID uint      `json:"hall_booking_id" gorm:"uniqueIndex;not null"`
	InvoiceNumber string    `json:"invoice_number" gorm:"uniqueIndex;not null"`
	InvoiceDate   time.Time `json:"invoice_date" gorm:"not null"`
	DueDate       time.Time `json:"due_date" gorm:"not null"`
	TotalAmount   float64   `json:"total_amount" gorm:"not null"`
	Status        string    `json:"status" gorm:"default:'sent'"` // sent, paid, overdue, cancelled
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"index"`

	// Relation
	HallBooking *HallBooking `json:"-" gorm:"foreignKey:HallBookingID"`
}

// CreateHallBookingRequest is the request payload for creating a hall booking
type CreateHallBookingRequest struct {
	OrganizerName   string  `json:"organizer_name" binding:"required,min=2,max=100"`
	OrganizerEmail  string  `json:"organizer_email" binding:"required,email,max=100"`
	OrganizerPhone  string  `json:"organizer_phone" binding:"required,min=10,max=20"`
	EventType       string  `json:"event_type" binding:"required,oneof=party wedding meeting christening funeral corporate fete other"`
	GuestCount      int     `json:"guest_count" binding:"required,min=1,max=100"`
	SpecialRequests string  `json:"special_requests" binding:"max=500"`
	BookingDate     string  `json:"booking_date" binding:"required,datetime=2006-01-02"`
	StartTime       string  `json:"start_time" binding:"required,oneof=08:00 09:00 10:00 11:00 12:00 13:00 14:00 15:00 16:00 17:00 18:00 19:00 20:00 21:00 22:00 23:00"`
	EndTime         string  `json:"end_time" binding:"required,oneof=08:00 09:00 10:00 11:00 12:00 13:00 14:00 15:00 16:00 17:00 18:00 19:00 20:00 21:00 22:00 23:00"`
	TotalPrice      float64 `json:"total_price" binding:"required,min=0"`
	DepositRequired float64 `json:"deposit_required" binding:"required,min=0"`
	PaymentMethod   string  `json:"payment_method" binding:"required,oneof=cash onsite online"`

	// Creator tracking fields (set by backend based on authentication)
	CreatedBy     *uint  `json:"-"` // Set by backend
	CreatedByType string `json:"-"` // Set by backend
	CreatorEmail  string `json:"-"` // Set by backend
	CreatorName   string `json:"-"` // Set by backend
}

// UpdateHallBookingRequest is the request payload for updating a hall booking
type UpdateHallBookingRequest struct {
	OrganizerName   *string  `json:"organizer_name" binding:"omitempty,min=2,max=100"`
	OrganizerEmail  *string  `json:"organizer_email" binding:"omitempty,email,max=100"`
	OrganizerPhone  *string  `json:"organizer_phone" binding:"omitempty,min=10,max=20"`
	EventType       *string  `json:"event_type" binding:"omitempty,oneof=party wedding meeting christening funeral corporate fete other"`
	GuestCount      *int     `json:"guest_count" binding:"omitempty,min=1,max=100"`
	SpecialRequests *string  `json:"special_requests" binding:"omitempty,max=500"`
	BookingDate     *string  `json:"booking_date" binding:"omitempty,datetime=2006-01-02"`
	StartTime       *string  `json:"start_time" binding:"omitempty,oneof=08:00 09:00 10:00 11:00 12:00 13:00 14:00 15:00 16:00 17:00 18:00 19:00 20:00 21:00 22:00 23:00"`
	EndTime         *string  `json:"end_time" binding:"omitempty,oneof=08:00 09:00 10:00 11:00 12:00 13:00 14:00 15:00 16:00 17:00 18:00 19:00 20:00 21:00 22:00 23:00"`
	TotalPrice      *float64 `json:"total_price" binding:"omitempty,min=0"`
	DepositRequired *float64 `json:"deposit_required" binding:"omitempty,min=0"`
	PaymentMethod   *string  `json:"payment_method" binding:"omitempty,oneof=cash onsite online"`
	Status          *string  `json:"status" binding:"omitempty,oneof=pending confirmed cancelled completed"`
}

// HallBookingResponse is the response payload for hall booking
type HallBookingResponse struct {
	ID              uint      `json:"id"`
	BookingID       string    `json:"booking_id"`
	OrganizerName   string    `json:"organizer_name"`
	OrganizerEmail  string    `json:"organizer_email"`
	OrganizerPhone  string    `json:"organizer_phone"`
	EventType       string    `json:"event_type"`
	GuestCount      int       `json:"guest_count"`
	SpecialRequests string    `json:"special_requests"`
	BookingDate     string    `json:"booking_date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	TotalPrice      float64   `json:"total_price"`
	DepositRequired float64   `json:"deposit_required"`
	PaymentMethod   string    `json:"payment_method"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Payments        []Payment `json:"payments"`
	Invoice         *Invoice  `json:"invoice"`

	// Creator tracking fields
	CreatedByType string `json:"created_by_type"`
	CreatorEmail  string `json:"creator_email"`
	CreatorName   string `json:"creator_name"`
}

// HallBookingListResponse is the paginated response for hall booking list
type HallBookingListResponse struct {
	Data []HallBookingResponse `json:"data"`
	Meta PaginationMeta        `json:"meta"`
}

// HallBookingActivityResponse is the response for recent hall booking activity
type HallBookingActivityResponse struct {
	RecentBookings []HallBookingResponse `json:"recent_bookings"`
	TotalCount     int64                 `json:"total_count"`
}

// BookingStatusHistory tracks status changes for hall bookings
type BookingStatusHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BookingID uint      `gorm:"not null" json:"booking_id"`
	OldStatus *string   `json:"old_status"`
	NewStatus string    `gorm:"not null" json:"new_status"`
	ChangedBy uint      `gorm:"not null" json:"changed_by"`
	ChangedAt time.Time `gorm:"not null" json:"changed_at"`
	Notes     string    `json:"notes"`

	// Relations - temporarily disabled for migration
	// Booking       HallBooking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	// ChangedByUser User        `json:"changed_by_user,omitempty" gorm:"foreignKey:ChangedBy"`
}

// AdminBookingRequest represents a request to update booking status
type AdminBookingRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed completed cancelled"`
	Notes  string `json:"notes" binding:"max=500"`
}

// AdminBookingResponse represents admin booking response with user info
type AdminBookingResponse struct {
	ID              uint                   `json:"id"`
	BookingID       string                 `json:"booking_id"`
	OrganizerName   string                 `json:"organizer_name"`
	OrganizerEmail  string                 `json:"organizer_email"`
	OrganizerPhone  string                 `json:"organizer_phone"`
	EventType       string                 `json:"event_type"`
	GuestCount      int                    `json:"guest_count"`
	SpecialRequests string                 `json:"special_requests"`
	BookingDate     string                 `json:"booking_date"`
	StartTime       string                 `json:"start_time"`
	EndTime         string                 `json:"end_time"`
	TotalPrice      float64                `json:"total_price"`
	DepositRequired float64                `json:"deposit_required"`
	PaymentMethod   string                 `json:"payment_method"`
	Status          string                 `json:"status"`
	ConfirmedBy     *User                  `json:"confirmed_by,omitempty"`
	ConfirmedAt     *time.Time             `json:"confirmed_at,omitempty"`
	CancelledBy     *User                  `json:"cancelled_by,omitempty"`
	CancelledAt     *time.Time             `json:"cancelled_at,omitempty"`
	UpdatedBy       *User                  `json:"updated_by,omitempty"`
	StatusHistory   []BookingStatusHistory `json:"status_history,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// HallAvailability represents hall availability for a specific date
type HallAvailability struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}
