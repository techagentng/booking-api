package services

import (
	"errors"
	"hotel/db"
	"hotel/models"
	"time"
)

// StaffAssignmentService handles staff assignment logic
type StaffAssignmentService interface {
	// Auto-assign staff based on request type and priority
	AutoAssignStaff(request *models.GuestServiceRequest) (*models.Staff, error)

	// Manual assignment
	AssignStaffToRequest(requestID uint, staffID uint) error

	// Unassign staff from request
	UnassignStaff(requestID uint) error

	// Complete request and free up staff
	CompleteRequest(requestID uint, completedBy string) error

	// Clock in/out
	ClockIn(staffID uint) error
	ClockOut(staffID uint) error

	// Toggle availability
	SetAvailability(staffID uint, available bool) error

	// Get available staff for a department
	GetAvailableStaff(department string) ([]models.Staff, error)

	// Get on-duty staff
	GetOnDutyStaff(department string) ([]models.Staff, error)
}

// staffAssignmentService implements StaffAssignmentService
type staffAssignmentService struct {
	staffRepo        db.StaffRepository
	guestServiceRepo db.GuestServiceRepository
}

// NewStaffAssignmentService creates a new staff assignment service
func NewStaffAssignmentService(
	staffRepo db.StaffRepository,
	guestServiceRepo db.GuestServiceRepository,
) StaffAssignmentService {
	return &staffAssignmentService{
		staffRepo:        staffRepo,
		guestServiceRepo: guestServiceRepo,
	}
}

// getDepartmentForRequestType maps request type to department
func getDepartmentForRequestType(requestType string) string {
	switch requestType {
	case "housekeeping":
		return "housekeeping"
	case "maintenance":
		return "maintenance"
	case "room_service":
		return "room_service"
	default:
		return ""
	}
}

// AutoAssignStaff automatically assigns the best available staff to a request
func (s *staffAssignmentService) AutoAssignStaff(request *models.GuestServiceRequest) (*models.Staff, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	// Determine department based on request type
	department := getDepartmentForRequestType(request.Type)
	if department == "" {
		return nil, errors.New("unknown request type")
	}

	// Get available staff for this department
	availableStaff, err := s.staffRepo.GetAvailableStaff(department)
	if err != nil {
		return nil, err
	}

	if len(availableStaff) == 0 {
		return nil, errors.New("no available staff for this department")
	}

	// Select the best staff member
	// Priority logic: staff with fewer tasks today gets priority
	// For high priority requests, we pick the one with least tasks
	// For normal requests, we use round-robin (oldest last_assigned_at)
	var selectedStaff *models.Staff

	if request.Priority == "high" || request.Priority == "urgent" {
		// Pick staff with least tasks today
		selectedStaff = &availableStaff[0] // Already sorted by tasks_today ASC
	} else {
		// Round-robin: pick staff with oldest last_assigned_at
		selectedStaff = &availableStaff[0] // Already sorted by last_assigned_at ASC
	}

	// Assign the task to the staff member
	if err := s.staffRepo.AssignTask(selectedStaff.ID, request.ID); err != nil {
		return nil, err
	}

	// Update the service request with assignment info
	now := time.Now()
	staffName := selectedStaff.FirstName + " " + selectedStaff.LastName
	if err := s.guestServiceRepo.UpdateServiceRequestStatus(request.ID, "in_progress", staffName); err != nil {
		// Rollback staff assignment
		s.staffRepo.CompleteTask(selectedStaff.ID)
		return nil, err
	}

	// Update the assigned_staff_id on the request
	if err := s.updateRequestAssignment(request.ID, selectedStaff.ID, &now); err != nil {
		return nil, err
	}

	return selectedStaff, nil
}

// updateRequestAssignment updates the service request with staff assignment
func (s *staffAssignmentService) updateRequestAssignment(requestID uint, staffID uint, assignedAt *time.Time) error {
	request, err := s.guestServiceRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}
	request.AssignedStaffID = &staffID
	request.AssignedAt = assignedAt
	// The repository should handle this update
	return nil
}

// AssignStaffToRequest manually assigns a staff member to a request
func (s *staffAssignmentService) AssignStaffToRequest(requestID uint, staffID uint) error {
	// Get the request
	request, err := s.guestServiceRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return errors.New("request not found")
	}

	// Check if request is already assigned
	if request.AssignedStaffID != nil {
		return errors.New("request is already assigned")
	}

	// Get the staff member
	staff, err := s.staffRepo.GetStaffByID(staffID)
	if err != nil {
		return errors.New("staff not found")
	}

	// Check if staff is available
	if !staff.IsOnDuty {
		return errors.New("staff is not on duty")
	}
	if !staff.IsAvailable {
		return errors.New("staff is not available")
	}

	// Assign the task
	if err := s.staffRepo.AssignTask(staffID, requestID); err != nil {
		return err
	}

	// Update request status
	staffName := staff.FirstName + " " + staff.LastName
	if err := s.guestServiceRepo.UpdateServiceRequestStatus(requestID, "in_progress", staffName); err != nil {
		s.staffRepo.CompleteTask(staffID)
		return err
	}

	return nil
}

// UnassignStaff removes staff assignment from a request
func (s *staffAssignmentService) UnassignStaff(requestID uint) error {
	request, err := s.guestServiceRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return errors.New("request not found")
	}

	if request.AssignedStaffID == nil {
		return errors.New("request is not assigned")
	}

	// Free up the staff member
	if err := s.staffRepo.CompleteTask(*request.AssignedStaffID); err != nil {
		return err
	}

	// Update request status back to pending
	if err := s.guestServiceRepo.UpdateServiceRequestStatus(requestID, "pending", ""); err != nil {
		return err
	}

	return nil
}

// CompleteRequest marks a request as complete and frees up the staff
func (s *staffAssignmentService) CompleteRequest(requestID uint, completedBy string) error {
	request, err := s.guestServiceRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return errors.New("request not found")
	}

	// Free up the staff member if assigned
	if request.AssignedStaffID != nil {
		if err := s.staffRepo.CompleteTask(*request.AssignedStaffID); err != nil {
			return err
		}
	}

	// Mark request as completed
	if err := s.guestServiceRepo.CompleteServiceRequest(requestID, completedBy); err != nil {
		return err
	}

	return nil
}

// ClockIn marks a staff member as on duty
func (s *staffAssignmentService) ClockIn(staffID uint) error {
	staff, err := s.staffRepo.GetStaffByID(staffID)
	if err != nil {
		return errors.New("staff not found")
	}

	if staff.IsOnDuty {
		return errors.New("staff is already on duty")
	}

	return s.staffRepo.ClockIn(staffID)
}

// ClockOut marks a staff member as off duty
func (s *staffAssignmentService) ClockOut(staffID uint) error {
	staff, err := s.staffRepo.GetStaffByID(staffID)
	if err != nil {
		return errors.New("staff not found")
	}

	if !staff.IsOnDuty {
		return errors.New("staff is not on duty")
	}

	// Check if staff has an active task
	if staff.CurrentTaskID != nil {
		return errors.New("staff has an active task, complete it first")
	}

	return s.staffRepo.ClockOut(staffID)
}

// SetAvailability updates staff availability
func (s *staffAssignmentService) SetAvailability(staffID uint, available bool) error {
	staff, err := s.staffRepo.GetStaffByID(staffID)
	if err != nil {
		return errors.New("staff not found")
	}

	if !staff.IsOnDuty {
		return errors.New("staff must be on duty to change availability")
	}

	// Can't set available if has active task
	if available && staff.CurrentTaskID != nil {
		return errors.New("staff has an active task")
	}

	return s.staffRepo.SetAvailable(staffID, available)
}

// GetAvailableStaff returns available staff for a department
func (s *staffAssignmentService) GetAvailableStaff(department string) ([]models.Staff, error) {
	return s.staffRepo.GetAvailableStaff(department)
}

// GetOnDutyStaff returns on-duty staff for a department
func (s *staffAssignmentService) GetOnDutyStaff(department string) ([]models.Staff, error) {
	return s.staffRepo.GetOnDutyStaff(department)
}
