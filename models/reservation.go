package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Reservation represents a hotel reservation/booking
type Reservation struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	ConfirmationNumber string     `json:"confirmation_number" gorm:"uniqueIndex;not null"`
	GuestID            uint       `json:"guest_id" gorm:"not null;index"`
	RoomID             uint       `json:"room_id" gorm:"not null;index"`
	CheckInDate        time.Time  `json:"check_in_date" gorm:"index"`
	CheckOutDate       time.Time  `json:"check_out_date" gorm:"index"`
	NumberOfGuests     int        `json:"number_of_guests" gorm:"default:1"`
	NumberOfNights     int        `json:"number_of_nights"`
	PricePerNight      float64    `json:"price_per_night"`
	TotalAmount        float64    `json:"total_amount"`
	Status             string     `json:"status" gorm:"default:pending;index"`   // pending, confirmed, checked-in, checked-out, cancelled
	PaymentStatus      string     `json:"payment_status" gorm:"default:pending"` // pending, paid, partially_paid, refunded
	PaymentMethod      string     `json:"payment_method"`
	SpecialRequests    string     `json:"special_requests"`
	CheckedInAt        *time.Time `json:"checked_in_at"`
	CheckedOutAt       *time.Time `json:"checked_out_at"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancellationReason string     `json:"cancellation_reason"`
	RefundAmount       float64    `json:"refund_amount"`
	AdditionalCharges  float64    `json:"additional_charges"`
	DepositCollected   float64    `json:"deposit_collected"`
	Notes              string     `json:"notes"`
	IDDocumentURL      string     `json:"id_document_url"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Guest *Guest `json:"guest,omitempty" gorm:"foreignKey:GuestID"`
	Room  *Room  `json:"room,omitempty" gorm:"foreignKey:RoomID"`
}

// ReservationListItem is a flattened view for list display
type ReservationListItem struct {
	ID                 uint      `json:"id"`
	ConfirmationNumber string    `json:"confirmation_number"`
	GuestID            uint      `json:"guest_id"`
	GuestName          string    `json:"guest_name"`
	GuestEmail         string    `json:"guest_email"`
	GuestPhone         string    `json:"guest_phone"`
	RoomID             uint      `json:"room_id"`
	RoomNumber         string    `json:"room_number"`
	RoomType           string    `json:"room_type"`
	CheckInDate        string    `json:"check_in_date"`
	CheckOutDate       string    `json:"check_out_date"`
	NumberOfGuests     int       `json:"number_of_guests"`
	NumberOfNights     int       `json:"number_of_nights"`
	PricePerNight      float64   `json:"price_per_night"`
	TotalAmount        float64   `json:"total_amount"`
	Status             string    `json:"status"`
	PaymentStatus      string    `json:"payment_status"`
	SpecialRequests    string    `json:"special_requests"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ReservationDetails is the detailed view with nested guest and room
type ReservationDetails struct {
	ID                 uint                 `json:"id"`
	ConfirmationNumber string               `json:"confirmation_number"`
	Guest              ReservationGuestInfo `json:"guest"`
	Room               ReservationRoomInfo  `json:"room"`
	CheckInDate        string               `json:"check_in_date"`
	CheckOutDate       string               `json:"check_out_date"`
	NumberOfGuests     int                  `json:"number_of_guests"`
	NumberOfNights     int                  `json:"number_of_nights"`
	PricePerNight      float64              `json:"price_per_night"`
	TotalAmount        float64              `json:"total_amount"`
	Status             string               `json:"status"`
	PaymentStatus      string               `json:"payment_status"`
	PaymentMethod      string               `json:"payment_method"`
	SpecialRequests    string               `json:"special_requests"`
	CheckedInAt        *time.Time           `json:"checked_in_at"`
	CheckedOutAt       *time.Time           `json:"checked_out_at"`
	CancelledAt        *time.Time           `json:"cancelled_at"`
	AdditionalCharges  float64              `json:"additional_charges"`
	Notes              string               `json:"notes"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

// ReservationGuestInfo is guest info embedded in reservation details
type ReservationGuestInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	IDType   string `json:"id_type"`
	IDNumber string `json:"id_number"`
}

// ReservationRoomInfo is room info embedded in reservation details
type ReservationRoomInfo struct {
	ID            uint    `json:"id"`
	RoomNumber    string  `json:"room_number"`
	RoomType      string  `json:"room_type"`
	Floor         int     `json:"floor"`
	BedType       string  `json:"bed_type"`
	PricePerNight float64 `json:"price_per_night"`
}

// FlexInt handles both string and int JSON values
type FlexInt int

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	// Try as int first
	var i int
	if err := json.Unmarshal(b, &i); err == nil {
		*fi = FlexInt(i)
		return nil
	}
	// Try as string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		parsed, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*fi = FlexInt(parsed)
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into FlexInt", string(b))
}

// FlexUint handles both string and uint JSON values
type FlexUint uint

func (fu *FlexUint) UnmarshalJSON(b []byte) error {
	// Try as uint first
	var u uint
	if err := json.Unmarshal(b, &u); err == nil {
		*fu = FlexUint(u)
		return nil
	}
	// Try as string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		parsed, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return err
		}
		*fu = FlexUint(parsed)
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into FlexUint", string(b))
}

// CreateReservationRequest is the request payload for creating a reservation
type CreateReservationRequest struct {
	// Either GuestID (existing guest) OR guest details (new guest)
	GuestID FlexUint `json:"guest_id"`
	// New guest fields (used when GuestID is 0)
	GuestName        string `json:"guest_name"`
	GuestEmail       string `json:"guest_email"`
	GuestPhone       string `json:"guest_phone"`
	GuestNationality string `json:"guest_nationality"`
	GuestIDType      string `json:"guest_id_type"`
	GuestIDNumber    string `json:"guest_id_number"`

	// Reservation fields
	RoomID          FlexUint `json:"room_id" binding:"required"`
	CheckInDate     string   `json:"check_in_date" binding:"required"`
	CheckOutDate    string   `json:"check_out_date" binding:"required"`
	NumberOfGuests  FlexInt  `json:"number_of_guests" binding:"required"`
	SpecialRequests string   `json:"special_requests"`
	PaymentMethod   string   `json:"payment_method"`
}

// GetGuestID returns GuestID as uint
func (r *CreateReservationRequest) GetGuestID() uint {
	return uint(r.GuestID)
}

// GetRoomID returns RoomID as uint
func (r *CreateReservationRequest) GetRoomID() uint {
	return uint(r.RoomID)
}

// GetNumberOfGuests returns NumberOfGuests as int
func (r *CreateReservationRequest) GetNumberOfGuests() int {
	return int(r.NumberOfGuests)
}

// IsNewGuest returns true if this request is for creating a new guest
func (r *CreateReservationRequest) IsNewGuest() bool {
	return r.GetGuestID() == 0 && r.GuestEmail != ""
}

// UpdateReservationRequest is the request payload for updating a reservation
type UpdateReservationRequest struct {
	CheckInDate     *string `json:"check_in_date"`
	CheckOutDate    *string `json:"check_out_date"`
	NumberOfGuests  *int    `json:"number_of_guests"`
	SpecialRequests *string `json:"special_requests"`
	PaymentMethod   *string `json:"payment_method"`
}

// UpdateReservationStatusRequest is the request for status updates
type UpdateReservationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// CancelReservationRequest is the request for cancellation
type CancelReservationRequest struct {
	CancellationReason string  `json:"cancellation_reason"`
	RefundAmount       float64 `json:"refund_amount"`
}

// CheckInRequest is the request for check-in
type CheckInRequest struct {
	ActualCheckInTime string  `json:"actual_check_in_time"`
	RoomKeyIssued     bool    `json:"room_key_issued"`
	PaymentVerified   bool    `json:"payment_verified"`
	IDVerified        bool    `json:"id_verified"`
	DepositCollected  float64 `json:"deposit_collected"`
	Notes             string  `json:"notes"`
}

// CheckInResponse is the response after check-in
type CheckInResponse struct {
	ID               uint    `json:"id"`
	Status           string  `json:"status"`
	CheckedInAt      string  `json:"checked_in_at"`
	RoomStatus       string  `json:"room_status"`
	RoomKeyNumber    string  `json:"room_key_number,omitempty"`
	DepositCollected float64 `json:"deposit_collected"`
}

// CheckOutRequest is the request for check-out
type CheckOutRequest struct {
	ActualCheckOutTime string  `json:"actual_check_out_time"`
	AdditionalCharges  float64 `json:"additional_charges"`
	DamageReported     bool    `json:"damage_reported"`
	Notes              string  `json:"notes"`
}

// ReservationListResponse is the paginated response for reservation list
type ReservationListResponse struct {
	Data []ReservationListItem `json:"data"`
	Meta PaginationMeta        `json:"meta"`
}

// ReservationStats holds reservation statistics
type ReservationStats struct {
	TotalReservations   int64   `json:"total_reservations"`
	Confirmed           int64   `json:"confirmed"`
	Pending             int64   `json:"pending"`
	CheckedIn           int64   `json:"checked_in"`
	CheckedOut          int64   `json:"checked_out"`
	Cancelled           int64   `json:"cancelled"`
	TotalRevenue        float64 `json:"total_revenue"`
	AverageBookingValue float64 `json:"average_booking_value"`
	OccupancyRate       float64 `json:"occupancy_rate"`
}

// CheckInStats holds check-in statistics for a specific date
type CheckInStats struct {
	Date               string `json:"date"`
	TotalArrivals      int64  `json:"total_arrivals"`
	CheckedIn          int64  `json:"checked_in"`
	PendingCheckIn     int64  `json:"pending_checkin"`
	EarlyCheckIns      int64  `json:"early_checkins"`
	LateCheckIns       int64  `json:"late_checkins"`
	NoShows            int64  `json:"no_shows"`
	AverageCheckInTime string `json:"average_checkin_time"`
}

// PaymentDetails holds payment information for a reservation
type PaymentDetails struct {
	ReservationID    uint    `json:"reservation_id"`
	TotalAmount      float64 `json:"total_amount"`
	PaidAmount       float64 `json:"paid_amount"`
	PendingAmount    float64 `json:"pending_amount"`
	PaymentStatus    string  `json:"payment_status"`
	PaymentMethod    string  `json:"payment_method,omitempty"`
	PaymentDate      string  `json:"payment_date,omitempty"`
	DepositRequired  float64 `json:"deposit_required"`
	DepositCollected float64 `json:"deposit_collected"`
}

// DashboardStats holds daily dashboard statistics
type DashboardStats struct {
	// Booking stats
	NewBookingsToday  int64 `json:"new_bookings_today"`
	ScheduledBookings int64 `json:"scheduled_bookings"`

	// Check-in/out stats
	CheckInsToday    int64 `json:"check_ins_today"`
	CheckOutsToday   int64 `json:"check_outs_today"`
	PendingCheckIns  int64 `json:"pending_check_ins"`
	PendingCheckOuts int64 `json:"pending_check_outs"`

	// Room stats
	TotalRooms          int64 `json:"total_rooms"`
	AvailableRoomsToday int64 `json:"available_rooms_today"`
	OccupiedRoomsToday  int64 `json:"occupied_rooms_today"`
	SoldOutRoomsToday   int64 `json:"sold_out_rooms_today"`

	// Guest stats
	TotalGuests int64 `json:"total_guests"`
}
