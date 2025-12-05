package db

import (
	"hotel/models"
	"time"

	"gorm.io/gorm"
)

// StaffRepository defines staff database operations
type StaffRepository interface {
	CreateStaff(staff *models.Staff) (*models.Staff, error)
	GetStaffByID(id uint) (*models.Staff, error)
	GetAllStaff(page, pageSize int, department, status, shift string) ([]models.Staff, int64, error)
	UpdateStaff(id uint, updates map[string]interface{}) error
	DeleteStaff(id uint) error
	GetStaffByEmployeeID(employeeID string) (*models.Staff, error)
	GetStaffByDepartment(department string) ([]models.Staff, error)
	GetStaffStats() ([]models.StaffStats, error)

	// Availability management
	ClockIn(id uint) error
	ClockOut(id uint) error
	SetAvailable(id uint, available bool) error
	GetOnDutyStaff(department string) ([]models.Staff, error)
	GetAvailableStaff(department string) ([]models.Staff, error)
	AssignTask(staffID uint, taskID uint) error
	CompleteTask(staffID uint) error
	ResetDailyStats() error
}

// staffRepository implements StaffRepository
type staffRepository struct {
	db *gorm.DB
}

// NewStaffRepository creates a new staff repository
func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

// CreateStaff creates a new staff member
func (r *staffRepository) CreateStaff(staff *models.Staff) (*models.Staff, error) {
	if err := r.db.Create(staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// GetStaffByID retrieves a staff member by ID
func (r *staffRepository) GetStaffByID(id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// GetAllStaff retrieves all staff with pagination and filters
func (r *staffRepository) GetAllStaff(page, pageSize int, department, status, shift string) ([]models.Staff, int64, error) {
	var staff []models.Staff
	var total int64

	query := r.db.Model(&models.Staff{})

	// Apply filters
	if department != "" && department != "all" {
		query = query.Where("department = ?", department)
	}
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if shift != "" && shift != "all" {
		query = query.Where("shift = ?", shift)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.
		Order("department, last_name, first_name").
		Offset(offset).
		Limit(pageSize).
		Find(&staff).Error; err != nil {
		return nil, 0, err
	}

	return staff, total, nil
}

// UpdateStaff updates a staff member
func (r *staffRepository) UpdateStaff(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Staff{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteStaff deletes a staff member
func (r *staffRepository) DeleteStaff(id uint) error {
	return r.db.Delete(&models.Staff{}, id).Error
}

// GetStaffByEmployeeID retrieves a staff member by employee ID
func (r *staffRepository) GetStaffByEmployeeID(employeeID string) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.Where("employee_id = ?", employeeID).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// GetStaffByDepartment retrieves all staff in a department
func (r *staffRepository) GetStaffByDepartment(department string) ([]models.Staff, error) {
	var staff []models.Staff
	if err := r.db.Where("department = ? AND status = ?", department, "active").
		Order("last_name, first_name").
		Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// GetStaffStats retrieves staff statistics by department
func (r *staffRepository) GetStaffStats() ([]models.StaffStats, error) {
	var stats []models.StaffStats

	rows, err := r.db.Model(&models.Staff{}).
		Select(`
			department,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status = 'on_leave' THEN 1 ELSE 0 END) as on_leave,
			SUM(CASE WHEN shift = 'morning' THEN 1 ELSE 0 END) as morning_shift,
			SUM(CASE WHEN shift = 'afternoon' THEN 1 ELSE 0 END) as afternoon_shift,
			SUM(CASE WHEN shift = 'night' THEN 1 ELSE 0 END) as night_shift
		`).
		Group("department").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat models.StaffStats
		if err := rows.Scan(
			&stat.Department,
			&stat.Total,
			&stat.Active,
			&stat.OnLeave,
			&stat.MorningShift,
			&stat.AfternoonShift,
			&stat.NightShift,
		); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// ClockIn marks a staff member as on duty
func (r *staffRepository) ClockIn(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Staff{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_on_duty":     true,
		"is_available":   true,
		"clock_in_time":  now,
		"clock_out_time": nil,
	}).Error
}

// ClockOut marks a staff member as off duty
func (r *staffRepository) ClockOut(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Staff{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_on_duty":      false,
		"is_available":    false,
		"clock_out_time":  now,
		"current_task_id": nil,
	}).Error
}

// SetAvailable updates staff availability
func (r *staffRepository) SetAvailable(id uint, available bool) error {
	return r.db.Model(&models.Staff{}).Where("id = ?", id).Update("is_available", available).Error
}

// GetOnDutyStaff retrieves all on-duty staff, optionally filtered by department
func (r *staffRepository) GetOnDutyStaff(department string) ([]models.Staff, error) {
	var staff []models.Staff
	query := r.db.Where("is_on_duty = ? AND status = ?", true, "active")
	if department != "" {
		query = query.Where("department = ?", department)
	}
	if err := query.Order("last_name, first_name").Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// GetAvailableStaff retrieves all available staff for assignment
func (r *staffRepository) GetAvailableStaff(department string) ([]models.Staff, error) {
	var staff []models.Staff
	query := r.db.Where("is_on_duty = ? AND is_available = ? AND status = ?", true, true, "active")
	if department != "" {
		query = query.Where("department = ?", department)
	}
	// Order by tasks today (ascending) then by last assigned (oldest first for round-robin)
	if err := query.Order("tasks_today ASC, last_assigned_at ASC NULLS FIRST").Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// AssignTask assigns a task to a staff member
func (r *staffRepository) AssignTask(staffID uint, taskID uint) error {
	now := time.Now()
	return r.db.Model(&models.Staff{}).Where("id = ?", staffID).Updates(map[string]interface{}{
		"is_available":     false,
		"current_task_id":  taskID,
		"last_assigned_at": now,
		"tasks_today":      gorm.Expr("tasks_today + 1"),
	}).Error
}

// CompleteTask marks a staff member's current task as complete
func (r *staffRepository) CompleteTask(staffID uint) error {
	return r.db.Model(&models.Staff{}).Where("id = ?", staffID).Updates(map[string]interface{}{
		"is_available":    true,
		"current_task_id": nil,
		"tasks_completed": gorm.Expr("tasks_completed + 1"),
	}).Error
}

// ResetDailyStats resets daily task counters (should be called at midnight)
func (r *staffRepository) ResetDailyStats() error {
	return r.db.Model(&models.Staff{}).Updates(map[string]interface{}{
		"tasks_today":     0,
		"tasks_completed": 0,
	}).Error
}
