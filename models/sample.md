package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Guest represents a hotel guest
type Guest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Email       string         `gorm:"uniqueIndex" json:"email"`
	Phone       string         `json:"phone"`
	Nationality string         `json:"nationality"`
	IDType      string         `json:"id_type"` // Passport, Driver's License, etc.
	IDNumber    string         `json:"id_number"`
	JoinDate    time.Time      `json:"join_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relations
	Reservations     []Reservation     `gorm:"foreignKey:GuestID" json:"reservations,omitempty"`
	ServiceRequests  []ServiceRequest  `gorm:"foreignKey:GuestID" json:"service_requests,omitempty"`
	Preferences      GuestPreferences  `gorm:"foreignKey:GuestID" json:"preferences,omitempty"`
	AIInsights       GuestAIInsights   `gorm:"foreignKey:GuestID" json:"ai_insights,omitempty"`
}

// GuestPreferences stores guest preferences
type GuestPreferences struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	GuestID          uint            `gorm:"uniqueIndex" json:"guest_id"`
	RoomFloors       datatypes.JSONSlice `gorm:"type:jsonb" json:"room_floors"`
	MealTypes        datatypes.JSONSlice `gorm:"type:jsonb" json:"meal_types"`
	RoomTypes        datatypes.JSONSlice `gorm:"type:jsonb" json:"room_types"`
	SpecialRequests  datatypes.JSONSlice `gorm:"type:jsonb" json:"special_requests"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// GuestAIInsights stores AI-generated insights about guests
type GuestAIInsights struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	GuestID             uint      `gorm:"uniqueIndex" json:"guest_id"`
	MealPreference      string    `json:"meal_preference"`
	RoomPreference      string    `json:"room_preference"`
	ServicePattern      string    `json:"service_pattern"`
	RiskScore           string    `json:"risk_score"` // low, medium, high
	Recommendations     datatypes.JSONSlice `gorm:"type:jsonb" json:"recommendations"`
	Complaints          datatypes.JSONSlice `gorm:"type:jsonb" json:"complaints"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Room represents a hotel room
type Room struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RoomNumber    string    `gorm:"uniqueIndex" json:"room_number"`
	RoomType      string    `json:"room_type"` // Standard, Deluxe, Suite
	Floor         int       `json:"floor"`
	Capacity      int       `json:"capacity"`
	PricePerNight float64   `json:"price_per_night"`
	Status        string    `json:"status"` // available, occupied, maintenance, cleaning
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservations []Reservation `gorm:"foreignKey:RoomID" json:"reservations,omitempty"`
}

// Reservation represents a booking
type Reservation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BookingID     string    `gorm:"uniqueIndex" json:"booking_id"`
	GuestID       uint      `json:"guest_id"`
	RoomID        uint      `json:"room_id"`
	CheckInDate   time.Time `json:"check_in_date"`
	CheckOutDate  time.Time `json:"check_out_date"`
	Nights        int       `json:"nights"`
	TotalPrice    float64   `json:"total_price"`
	PaidAmount    float64   `json:"paid_amount"`
	Status        string    `json:"status"` // pending, confirmed, checked-in, checked-out, cancelled
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Guest         Guest              `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
	Room          Room               `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	ServiceRequests []ServiceRequest `gorm:"foreignKey:ReservationID" json:"service_requests,omitempty"`
}

// ServiceRequest represents a guest service request
type ServiceRequest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReservationID uint      `json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	ServiceType   string    `json:"service_type"` // room-service, housekeeping, maintenance, special-requests, transportation, general-assistance
	Status        string    `json:"status"` // pending, in-progress, completed, cancelled
	Priority      string    `json:"priority"` // low, medium, high
	Description   string    `gorm:"type:text" json:"description"`
	Notes         string    `gorm:"type:text" json:"notes"`
	AssignedTo    string    `json:"assigned_to"`
	RequestedAt   time.Time `json:"requested_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Guest       Guest       `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
}

// RoomServiceOrder represents a room service order
type RoomServiceOrder struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       string    `gorm:"uniqueIndex" json:"order_id"`
	ReservationID uint      `json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	Items         datatypes.JSONSlice `gorm:"type:jsonb" json:"items"`
	Subtotal      float64   `json:"subtotal"`
	DeliveryFee   float64   `json:"delivery_fee"`
	Total         float64   `json:"total"`
	Status        string    `json:"status"` // pending, preparing, delivered, cancelled
	SpecialNotes  string    `gorm:"type:text" json:"special_notes"`
	OrderedAt     time.Time `json:"ordered_at"`
	DeliveredAt   *time.Time `json:"delivered_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Guest       Guest       `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
}

// HousekeepingRequest represents a housekeeping service request
type HousekeepingRequest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReservationID uint      `json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	RequestType   string    `json:"request_type"` // cleaning, linens, schedule
	Description   string    `gorm:"type:text" json:"description"`
	ScheduleTime  string    `json:"schedule_time"` // morning, afternoon, immediate
	Status        string    `json:"status"` // pending, in-progress, completed
	AssignedTo    string    `json:"assigned_to"`
	RequestedAt   time.Time `json:"requested_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Guest       Guest       `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
}

// MaintenanceIssue represents a maintenance issue report
type MaintenanceIssue struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReservationID uint      `json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	IssueType     string    `json:"issue_type"` // ac, tv, wifi, plumbing, lights, door-lock, noise, other
	Description   string    `gorm:"type:text" json:"description"`
	Status        string    `json:"status"` // reported, in-progress, resolved
	Priority      string    `json:"priority"` // low, medium, high
	AssignedTo    string    `json:"assigned_to"`
	ReportedAt    time.Time `json:"reported_at"`
	ResolvedAt    *time.Time `json:"resolved_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Guest       Guest       `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
}

// MenuItem represents a room service menu item
type MenuItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"` // Food, Drinks
	Available   bool      `json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Staff represents hotel staff members
type Staff struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"` // manager, housekeeping, maintenance, front-desk, room-service
	Status    string    `json:"status"` // active, inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CheckIn represents a guest check-in record
type CheckIn struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReservationID uint      `gorm:"uniqueIndex" json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	RoomID        uint      `json:"room_id"`
	CheckInTime   time.Time `json:"check_in_time"`
	IDVerified    bool      `json:"id_verified"`
	KeyIssued     bool      `json:"key_issued"`
	DocumentsSigned bool    `json:"documents_signed"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Guest       Guest       `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
	Room        Room        `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

// CheckOut represents a guest check-out record
type CheckOut struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReservationID uint      `gorm:"uniqueIndex" json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	RoomID        uint      `json:"room_id"`
	CheckOutTime  time.Time `json:"check_out_time"`
	RoomCondition string    `json:"room_condition"` // excellent, good, fair, poor
	Charges       float64   `json:"charges"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Guest       Guest       `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
	Room        Room        `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}



// responses
package responses

import "time"

// Standard API Response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Pagination Response
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// ===== GUEST RESPONSES =====

type GuestResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Nationality string    `json:"nationality"`
	IDType      string    `json:"id_type"`
	IDNumber    string    `json:"id_number"`
	JoinDate    time.Time `json:"join_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GuestDetailResponse struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	Email        string                   `json:"email"`
	Phone        string                   `json:"phone"`
	Nationality  string                   `json:"nationality"`
	IDType       string                   `json:"id_type"`
	IDNumber     string                   `json:"id_number"`
	JoinDate     time.Time                `json:"join_date"`
	Reservations []ReservationResponse    `json:"reservations"`
	Preferences  GuestPreferencesResponse `json:"preferences"`
	AIInsights   GuestAIInsightsResponse  `json:"ai_insights"`
	Statistics   GuestStatisticsResponse  `json:"statistics"`
	ServiceUsage []ServiceUsageResponse   `json:"service_usage"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

type GuestPreferencesResponse struct {
	ID              uint     `json:"id"`
	RoomFloors      []string `json:"room_floors"`
	MealTypes       []string `json:"meal_types"`
	RoomTypes       []string `json:"room_types"`
	SpecialRequests []string `json:"special_requests"`
}

type GuestAIInsightsResponse struct {
	ID              uint     `json:"id"`
	MealPreference  string   `json:"meal_preference"`
	RoomPreference  string   `json:"room_preference"`
	ServicePattern  string   `json:"service_pattern"`
	RiskScore       string   `json:"risk_score"`
	Recommendations []string `json:"recommendations"`
	Complaints      []string `json:"complaints"`
}

type GuestStatisticsResponse struct {
	TotalStays     int       `json:"total_stays"`
	TotalSpent     float64   `json:"total_spent"`
	AverageSpend   float64   `json:"average_spend"`
	LastVisit      time.Time `json:"last_visit"`
	MostCommonRoom string    `json:"most_common_room"`
}

type ServiceUsageResponse struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
	Label string `json:"label"`
}

// ===== RESERVATION RESPONSES =====

type ReservationResponse struct {
	ID           uint      `json:"id"`
	BookingID    string    `json:"booking_id"`
	GuestID      uint      `json:"guest_id"`
	RoomID       uint      `json:"room_id"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	Nights       int       `json:"nights"`
	TotalPrice   float64   `json:"total_price"`
	PaidAmount   float64   `json:"paid_amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReservationDetailResponse struct {
	ID              uint                     `json:"id"`
	BookingID       string                   `json:"booking_id"`
	Guest           GuestResponse            `json:"guest"`
	Room            RoomResponse             `json:"room"`
	CheckInDate     time.Time                `json:"check_in_date"`
	CheckOutDate    time.Time                `json:"check_out_date"`
	CheckInTime     string                   `json:"check_in_time"`
	CheckOutTime    string                   `json:"check_out_time"`
	Nights          int                      `json:"nights"`
	TotalPrice      float64                  `json:"total_price"`
	PaidAmount      float64                  `json:"paid_amount"`
	Status          string                   `json:"status"`
	Preferences     []string                 `json:"preferences"`
	ServiceRequests []ServiceRequestResponse `json:"service_requests"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// ===== ROOM RESPONSES =====

type RoomResponse struct {
	ID            uint      `json:"id"`
	RoomNumber    string    `json:"room_number"`
	RoomType      string    `json:"room_type"`
	Floor         int       `json:"floor"`
	Capacity      int       `json:"capacity"`
	PricePerNight float64   `json:"price_per_night"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ===== SERVICE REQUEST RESPONSES =====

type ServiceRequestResponse struct {
	ID            uint       `json:"id"`
	ReservationID uint       `json:"reservation_id"`
	GuestID       uint       `json:"guest_id"`
	ServiceType   string     `json:"service_type"`
	Status        string     `json:"status"`
	Priority      string     `json:"priority"`
	Description   string     `json:"description"`
	Notes         string     `json:"notes"`
	AssignedTo    string     `json:"assigned_to"`
	RequestedAt   time.Time  `json:"requested_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ServiceRequestDetailResponse struct {
	ID            uint          `json:"id"`
	ReservationID uint          `json:"reservation_id"`
	Guest         GuestResponse `json:"guest"`
	Room          RoomResponse  `json:"room"`
	ServiceType   string        `json:"service_type"`
	Status        string        `json:"status"`
	Priority      string        `json:"priority"`
	Description   string        `json:"description"`
	Notes         string        `json:"notes"`
	AssignedTo    string        `json:"assigned_to"`
	RequestedAt   time.Time     `json:"requested_at"`
	CompletedAt   *time.Time    `json:"completed_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// ===== ROOM SERVICE ORDER RESPONSES =====

type RoomServiceOrderItemResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type RoomServiceOrderResponse struct {
	ID            uint                           `json:"id"`
	OrderID       string                         `json:"order_id"`
	ReservationID uint                           `json:"reservation_id"`
	GuestID       uint                           `json:"guest_id"`
	Items         []RoomServiceOrderItemResponse `json:"items"`
	Subtotal      float64                        `json:"subtotal"`
	DeliveryFee   float64                        `json:"delivery_fee"`
	Total         float64                        `json:"total"`
	Status        string                         `json:"status"`
	SpecialNotes  string                         `json:"special_notes"`
	OrderedAt     time.Time                      `json:"ordered_at"`
	DeliveredAt   *time.Time                     `json:"delivered_at"`
	CreatedAt     time.Time                      `json:"created_at"`
	UpdatedAt     time.Time                      `json:"updated_at"`
}

// ===== HOUSEKEEPING RESPONSES =====

type HousekeepingRequestResponse struct {
	ID            uint       `json:"id"`
	ReservationID uint       `json:"reservation_id"`
	GuestID       uint       `json:"guest_id"`
	RequestType   string     `json:"request_type"`
	Description   string     `json:"description"`
	ScheduleTime  string     `json:"schedule_time"`
	Status        string     `json:"status"`
	AssignedTo    string     `json:"assigned_to"`
	RequestedAt   time.Time  `json:"requested_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ===== MAINTENANCE RESPONSES =====

type MaintenanceIssueResponse struct {
	ID            uint       `json:"id"`
	ReservationID uint       `json:"reservation_id"`
	GuestID       uint       `json:"guest_id"`
	IssueType     string     `json:"issue_type"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	Priority      string     `json:"priority"`
	AssignedTo    string     `json:"assigned_to"`
	ReportedAt    time.Time  `json:"reported_at"`
	ResolvedAt    *time.Time `json:"resolved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ===== MENU ITEM RESPONSES =====

type MenuItemResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Available   bool      `json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ===== CHECK-IN/CHECK-OUT RESPONSES =====

type CheckInResponse struct {
	ID              uint      `json:"id"`
	ReservationID   uint      `json:"reservation_id"`
	GuestID         uint      `json:"guest_id"`
	RoomID          uint      `json:"room_id"`
	CheckInTime     time.Time `json:"check_in_time"`
	IDVerified      bool      `json:"id_verified"`
	KeyIssued       bool      `json:"key_issued"`
	DocumentsSigned bool      `json:"documents_signed"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CheckOutResponse struct {
	ID            uint      `json:"id"`
	ReservationID uint      `json:"reservation_id"`
	GuestID       uint      `json:"guest_id"`
	RoomID        uint      `json:"room_id"`
	CheckOutTime  time.Time `json:"check_out_time"`
	RoomCondition string    `json:"room_condition"`
	Charges       float64   `json:"charges"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ===== IN-ROOM TABLET RESPONSES =====

type InRoomTabletReservationResponse struct {
	RoomNumber   string    `json:"room_number"`
	RoomType     string    `json:"room_type"`
	GuestName    string    `json:"guest_name"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	CheckInTime  string    `json:"check_in_time"`
	CheckOutTime string    `json:"check_out_time"`
	Nights       int       `json:"nights"`
	BookingID    string    `json:"booking_id"`
	TotalPrice   float64   `json:"total_price"`
	PaidAmount   float64   `json:"paid_amount"`
	Preferences  []string  `json:"preferences"`
}

type InRoomTabletMenuResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Available   bool    `json:"available"`
}

// ===== DASHBOARD STATISTICS RESPONSES =====

type DashboardStatsResponse struct {
	TotalGuests            int64   `json:"total_guests"`
	TotalRooms             int64   `json:"total_rooms"`
	OccupiedRooms          int64   `json:"occupied_rooms"`
	AvailableRooms         int64   `json:"available_rooms"`
	MaintenanceRooms       int64   `json:"maintenance_rooms"`
	PendingCheckIns        int64   `json:"pending_check_ins"`
	PendingCheckOuts       int64   `json:"pending_check_outs"`
	PendingServiceRequests int64   `json:"pending_service_requests"`
	TodayRevenue           float64 `json:"today_revenue"`
	MonthRevenue           float64 `json:"month_revenue"`
	OccupancyRate          float64 `json:"occupancy_rate"`
}

type RoomStatusSummaryResponse struct {
	Available   int64 `json:"available"`
	Occupied    int64 `json:"occupied"`
	Maintenance int64 `json:"maintenance"`
	Cleaning    int64 `json:"cleaning"`
}

type ServiceRequestSummaryResponse struct {
	Pending    int64 `json:"pending"`
	InProgress int64 `json:"in_progress"`
	Completed  int64 `json:"completed"`
	Cancelled  int64 `json:"cancelled"`
}

// ===== REQUEST BODIES =====

type CreateGuestRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"required"`
	Nationality string `json:"nationality" binding:"required"`
	IDType      string `json:"id_type" binding:"required"`
	IDNumber    string `json:"id_number" binding:"required"`
}

type CreateReservationRequest struct {
	GuestID      uint      `json:"guest_id" binding:"required"`
	RoomID       uint      `json:"room_id" binding:"required"`
	CheckInDate  time.Time `json:"check_in_date" binding:"required"`
	CheckOutDate time.Time `json:"check_out_date" binding:"required"`
	TotalPrice   float64   `json:"total_price" binding:"required"`
}

type CreateServiceRequestRequest struct {
	ReservationID uint   `json:"reservation_id" binding:"required"`
	GuestID       uint   `json:"guest_id" binding:"required"`
	ServiceType   string `json:"service_type" binding:"required"`
	Priority      string `json:"priority" binding:"required"`
	Description   string `json:"description" binding:"required"`
}

type CreateRoomServiceOrderRequest struct {
	ReservationID uint                                `json:"reservation_id" binding:"required"`
	GuestID       uint                                `json:"guest_id" binding:"required"`
	Items         []CreateRoomServiceOrderItemRequest `json:"items" binding:"required"`
	SpecialNotes  string                              `json:"special_notes"`
}

type CreateRoomServiceOrderItemRequest struct {
	MenuItemID uint `json:"menu_item_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,min=1"`
}

type CreateHousekeepingRequestRequest struct {
	ReservationID uint   `json:"reservation_id" binding:"required"`
	GuestID       uint   `json:"guest_id" binding:"required"`
	RequestType   string `json:"request_type" binding:"required"`
	Description   string `json:"description" binding:"required"`
	ScheduleTime  string `json:"schedule_time"`
}

type CreateMaintenanceIssueRequest struct {
	ReservationID uint   `json:"reservation_id" binding:"required"`
	GuestID       uint   `json:"guest_id" binding:"required"`
	IssueType     string `json:"issue_type" binding:"required"`
	Description   string `json:"description" binding:"required"`
}

type CreateCheckInRequest struct {
	ReservationID   uint   `json:"reservation_id" binding:"required"`
	GuestID         uint   `json:"guest_id" binding:"required"`
	RoomID          uint   `json:"room_id" binding:"required"`
	IDVerified      bool   `json:"id_verified"`
	KeyIssued       bool   `json:"key_issued"`
	DocumentsSigned bool   `json:"documents_signed"`
	Notes           string `json:"notes"`
}

type UpdateServiceRequestRequest struct {
	Status     string `json:"status"`
	AssignedTo string `json:"assigned_to"`
	Notes      string `json:"notes"`
}

type UpdateMaintenanceIssueRequest struct {
	Status     string `json:"status"`
	AssignedTo string `json:"assigned_to"`
	Notes      string `json:"notes"`
}
