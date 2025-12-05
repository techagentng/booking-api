package models

import (
	"time"
)

// Staff represents a hotel staff member
type Staff struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeID   string    `gorm:"uniqueIndex;not null" json:"employee_id"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Phone        string    `json:"phone"`
	Department   string    `gorm:"not null" json:"department"` // housekeeping, maintenance, front_desk, room_service, management
	Position     string    `gorm:"not null" json:"position"`
	Status       string    `gorm:"default:active" json:"status"` // active, inactive, on_leave
	Shift        string    `json:"shift"`                        // morning, afternoon, night
	HireDate     time.Time `json:"hire_date"`
	ProfileImage string    `json:"profile_image,omitempty"`

	// Availability tracking
	IsOnDuty       bool       `gorm:"default:false" json:"is_on_duty"`
	IsAvailable    bool       `gorm:"default:true" json:"is_available"`
	CurrentTaskID  *uint      `json:"current_task_id,omitempty"`
	LastAssignedAt *time.Time `json:"last_assigned_at,omitempty"`
	ClockInTime    *time.Time `json:"clock_in_time,omitempty"`
	ClockOutTime   *time.Time `json:"clock_out_time,omitempty"`

	// Shift timing (24h format: "08:00")
	ShiftStartTime string `json:"shift_start_time"`
	ShiftEndTime   string `json:"shift_end_time"`

	// Daily workload stats
	TasksToday     int `gorm:"default:0" json:"tasks_today"`
	TasksCompleted int `gorm:"default:0" json:"tasks_completed"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// StaffListResponse is the response for listing staff
type StaffListResponse struct {
	Staff      []Staff        `json:"staff"`
	Pagination PaginationMeta `json:"pagination"`
}

// StaffStats represents staff statistics by department
type StaffStats struct {
	Department     string `json:"department"`
	Total          int64  `json:"total"`
	Active         int64  `json:"active"`
	OnLeave        int64  `json:"on_leave"`
	MorningShift   int64  `json:"morning_shift"`
	AfternoonShift int64  `json:"afternoon_shift"`
	NightShift     int64  `json:"night_shift"`
}

// CreateStaffRequest is the request for creating a staff member
type CreateStaffRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone"`
	Department string `json:"department" binding:"required"`
	Position   string `json:"position" binding:"required"`
	Shift      string `json:"shift"`
	HireDate   string `json:"hire_date"`
}

// UpdateStaffRequest is the request for updating a staff member
type UpdateStaffRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Department   string `json:"department"`
	Position     string `json:"position"`
	Status       string `json:"status"`
	Shift        string `json:"shift"`
	ProfileImage string `json:"profile_image"`
}
