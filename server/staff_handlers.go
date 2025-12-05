package server

import (
	"hotel/models"
	"hotel/server/response"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// handleGetStaff returns all staff with pagination and filters
func (s *Server) handleGetStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		department := c.DefaultQuery("department", "")
		status := c.DefaultQuery("status", "")
		shift := c.DefaultQuery("shift", "")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}

		staff, total, err := s.StaffRepository.GetAllStaff(page, pageSize, department, status, shift)
		if err != nil {
			log.Printf("handleGetStaff: error: %v", err)
			response.JSON(c, "Failed to fetch staff", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff retrieved", http.StatusOK, gin.H{
			"staff": staff,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		}, nil)
	}
}

// handleGetStaffByID returns a staff member by ID
func (s *Server) handleGetStaffByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid staff ID", http.StatusBadRequest, nil, err)
			return
		}

		staff, err := s.StaffRepository.GetStaffByID(uint(id))
		if err != nil {
			log.Printf("handleGetStaffByID: error: %v", err)
			response.JSON(c, "Staff not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Staff retrieved", http.StatusOK, staff, nil)
	}
}

// handleCreateStaff creates a new staff member
func (s *Server) handleCreateStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateStaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		// Check if employee ID already exists
		existing, _ := s.StaffRepository.GetStaffByEmployeeID(req.EmployeeID)
		if existing != nil {
			response.JSON(c, "Employee ID already exists", http.StatusConflict, nil, nil)
			return
		}

		// Parse hire date
		hireDate := time.Now()
		if req.HireDate != "" {
			if parsed, err := time.Parse("2006-01-02", req.HireDate); err == nil {
				hireDate = parsed
			}
		}

		staff := &models.Staff{
			EmployeeID: req.EmployeeID,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Phone:      req.Phone,
			Department: req.Department,
			Position:   req.Position,
			Status:     "active",
			Shift:      req.Shift,
			HireDate:   hireDate,
		}

		created, err := s.StaffRepository.CreateStaff(staff)
		if err != nil {
			log.Printf("handleCreateStaff: error: %v", err)
			response.JSON(c, "Failed to create staff", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff created", http.StatusCreated, created, nil)
	}
}

// handleUpdateStaff updates a staff member
func (s *Server) handleUpdateStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid staff ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateStaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		updates := make(map[string]interface{})
		if req.FirstName != "" {
			updates["first_name"] = req.FirstName
		}
		if req.LastName != "" {
			updates["last_name"] = req.LastName
		}
		if req.Email != "" {
			updates["email"] = req.Email
		}
		if req.Phone != "" {
			updates["phone"] = req.Phone
		}
		if req.Department != "" {
			updates["department"] = req.Department
		}
		if req.Position != "" {
			updates["position"] = req.Position
		}
		if req.Status != "" {
			updates["status"] = req.Status
		}
		if req.Shift != "" {
			updates["shift"] = req.Shift
		}
		if req.ProfileImage != "" {
			updates["profile_image"] = req.ProfileImage
		}

		if err := s.StaffRepository.UpdateStaff(uint(id), updates); err != nil {
			log.Printf("handleUpdateStaff: error: %v", err)
			response.JSON(c, "Failed to update staff", http.StatusInternalServerError, nil, err)
			return
		}

		staff, _ := s.StaffRepository.GetStaffByID(uint(id))
		response.JSON(c, "Staff updated", http.StatusOK, staff, nil)
	}
}

// handleDeleteStaff deletes a staff member
func (s *Server) handleDeleteStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid staff ID", http.StatusBadRequest, nil, err)
			return
		}

		if err := s.StaffRepository.DeleteStaff(uint(id)); err != nil {
			log.Printf("handleDeleteStaff: error: %v", err)
			response.JSON(c, "Failed to delete staff", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff deleted", http.StatusOK, nil, nil)
	}
}

// handleGetStaffByDepartment returns staff by department
func (s *Server) handleGetStaffByDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		department := c.Param("department")

		staff, err := s.StaffRepository.GetStaffByDepartment(department)
		if err != nil {
			log.Printf("handleGetStaffByDepartment: error: %v", err)
			response.JSON(c, "Failed to fetch staff", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff retrieved", http.StatusOK, staff, nil)
	}
}

// handleGetStaffStats returns staff statistics
func (s *Server) handleGetStaffStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := s.StaffRepository.GetStaffStats()
		if err != nil {
			log.Printf("handleGetStaffStats: error: %v", err)
			response.JSON(c, "Failed to fetch stats", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff stats retrieved", http.StatusOK, stats, nil)
	}
}

// handleClockIn clocks in a staff member
func (s *Server) handleClockIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid staff ID", http.StatusBadRequest, nil, err)
			return
		}

		if err := s.StaffRepository.ClockIn(uint(id)); err != nil {
			log.Printf("handleClockIn: error: %v", err)
			response.JSON(c, "Failed to clock in", http.StatusInternalServerError, nil, err)
			return
		}

		staff, _ := s.StaffRepository.GetStaffByID(uint(id))
		response.JSON(c, "Staff clocked in", http.StatusOK, staff, nil)
	}
}

// handleClockOut clocks out a staff member
func (s *Server) handleClockOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid staff ID", http.StatusBadRequest, nil, err)
			return
		}

		// Check if staff has active task
		staff, err := s.StaffRepository.GetStaffByID(uint(id))
		if err != nil {
			response.JSON(c, "Staff not found", http.StatusNotFound, nil, err)
			return
		}

		if staff.CurrentTaskID != nil {
			response.JSON(c, "Staff has an active task, complete it first", http.StatusBadRequest, nil, nil)
			return
		}

		if err := s.StaffRepository.ClockOut(uint(id)); err != nil {
			log.Printf("handleClockOut: error: %v", err)
			response.JSON(c, "Failed to clock out", http.StatusInternalServerError, nil, err)
			return
		}

		staff, _ = s.StaffRepository.GetStaffByID(uint(id))
		response.JSON(c, "Staff clocked out", http.StatusOK, staff, nil)
	}
}

// handleSetAvailability updates staff availability
func (s *Server) handleSetAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid staff ID", http.StatusBadRequest, nil, err)
			return
		}

		var req struct {
			Available bool `json:"available"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		if err := s.StaffRepository.SetAvailable(uint(id), req.Available); err != nil {
			log.Printf("handleSetAvailability: error: %v", err)
			response.JSON(c, "Failed to update availability", http.StatusInternalServerError, nil, err)
			return
		}

		staff, _ := s.StaffRepository.GetStaffByID(uint(id))
		response.JSON(c, "Availability updated", http.StatusOK, staff, nil)
	}
}

// handleGetOnDutyStaff returns on-duty staff
func (s *Server) handleGetOnDutyStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		department := c.DefaultQuery("department", "")

		staff, err := s.StaffRepository.GetOnDutyStaff(department)
		if err != nil {
			log.Printf("handleGetOnDutyStaff: error: %v", err)
			response.JSON(c, "Failed to fetch staff", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "On-duty staff retrieved", http.StatusOK, staff, nil)
	}
}

// handleGetAvailableStaff returns available staff
func (s *Server) handleGetAvailableStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		department := c.DefaultQuery("department", "")

		staff, err := s.StaffRepository.GetAvailableStaff(department)
		if err != nil {
			log.Printf("handleGetAvailableStaff: error: %v", err)
			response.JSON(c, "Failed to fetch staff", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Available staff retrieved", http.StatusOK, staff, nil)
	}
}

// handleAutoAssignRequest auto-assigns staff to a service request
func (s *Server) handleAutoAssignRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid request ID", http.StatusBadRequest, nil, err)
			return
		}

		// Get the request
		request, err := s.GuestServiceRepository.GetServiceRequestByID(uint(requestID))
		if err != nil {
			response.JSON(c, "Request not found", http.StatusNotFound, nil, err)
			return
		}

		// Check if already assigned
		if request.AssignedStaffID != nil {
			response.JSON(c, "Request is already assigned", http.StatusBadRequest, nil, nil)
			return
		}

		// Get department based on request type
		department := request.Type // housekeeping or maintenance

		// Get available staff - first try department-specific on-duty staff
		availableStaff, err := s.StaffRepository.GetAvailableStaff(department)
		if err != nil || len(availableStaff) == 0 {
			// Fallback 1: try any on-duty available staff
			availableStaff, err = s.StaffRepository.GetAvailableStaff("")
		}
		if err != nil || len(availableStaff) == 0 {
			// Fallback 2: get any active staff from the department (ignore on-duty status)
			availableStaff, err = s.StaffRepository.GetStaffByDepartment(department)
		}
		if err != nil || len(availableStaff) == 0 {
			// Fallback 3: get any active staff at all
			availableStaff, _, err = s.StaffRepository.GetAllStaff(1, 10, "", "active", "")
		}
		if err != nil || len(availableStaff) == 0 {
			response.JSON(c, "No available staff", http.StatusNotFound, nil, err)
			return
		}

		// Select best staff (first one - already sorted by workload)
		selectedStaff := availableStaff[0]

		// Assign task to staff
		if err := s.StaffRepository.AssignTask(selectedStaff.ID, uint(requestID)); err != nil {
			log.Printf("handleAutoAssignRequest: error assigning task: %v", err)
			response.JSON(c, "Failed to assign staff", http.StatusInternalServerError, nil, err)
			return
		}

		// Update request with staff assignment (sets assigned_staff_id, assigned_to, assigned_at)
		staffName := selectedStaff.FirstName + " " + selectedStaff.LastName
		if err := s.GuestServiceRepository.AssignStaffToRequest(uint(requestID), selectedStaff.ID, staffName); err != nil {
			// Rollback
			s.StaffRepository.CompleteTask(selectedStaff.ID)
			response.JSON(c, "Failed to update request", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff assigned successfully", http.StatusOK, gin.H{
			"request_id": requestID,
			"staff":      selectedStaff,
		}, nil)
	}
}

// handleManualAssignRequest manually assigns staff to a service request
func (s *Server) handleManualAssignRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid request ID", http.StatusBadRequest, nil, err)
			return
		}

		var req struct {
			StaffID uint `json:"staff_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		// Get the request
		request, err := s.GuestServiceRepository.GetServiceRequestByID(uint(requestID))
		if err != nil {
			response.JSON(c, "Request not found", http.StatusNotFound, nil, err)
			return
		}

		// Check if already assigned
		if request.AssignedStaffID != nil {
			response.JSON(c, "Request is already assigned", http.StatusBadRequest, nil, nil)
			return
		}

		// Get the staff member
		staff, err := s.StaffRepository.GetStaffByID(req.StaffID)
		if err != nil {
			response.JSON(c, "Staff not found", http.StatusNotFound, nil, err)
			return
		}

		// Check availability
		if !staff.IsOnDuty {
			response.JSON(c, "Staff is not on duty", http.StatusBadRequest, nil, nil)
			return
		}
		if !staff.IsAvailable {
			response.JSON(c, "Staff is not available", http.StatusBadRequest, nil, nil)
			return
		}

		// Assign task
		if err := s.StaffRepository.AssignTask(req.StaffID, uint(requestID)); err != nil {
			response.JSON(c, "Failed to assign staff", http.StatusInternalServerError, nil, err)
			return
		}

		// Update request with staff assignment (sets assigned_staff_id, assigned_to, assigned_at)
		staffName := staff.FirstName + " " + staff.LastName
		if err := s.GuestServiceRepository.AssignStaffToRequest(uint(requestID), req.StaffID, staffName); err != nil {
			s.StaffRepository.CompleteTask(req.StaffID)
			response.JSON(c, "Failed to update request", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Staff assigned successfully", http.StatusOK, gin.H{
			"request_id": requestID,
			"staff":      staff,
		}, nil)
	}
}

// handleCompleteServiceRequest completes a service request and frees staff
func (s *Server) handleCompleteServiceRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid request ID", http.StatusBadRequest, nil, err)
			return
		}

		var req struct {
			CompletedBy string `json:"completed_by"`
		}
		c.ShouldBindJSON(&req)

		// Get the request
		request, err := s.GuestServiceRepository.GetServiceRequestByID(uint(requestID))
		if err != nil {
			response.JSON(c, "Request not found", http.StatusNotFound, nil, err)
			return
		}

		// Free up staff if assigned
		if request.AssignedStaffID != nil {
			if err := s.StaffRepository.CompleteTask(*request.AssignedStaffID); err != nil {
				log.Printf("handleCompleteServiceRequest: error freeing staff: %v", err)
			}
		}

		// Complete the request
		completedBy := req.CompletedBy
		if completedBy == "" && request.AssignedTo != "" {
			completedBy = request.AssignedTo
		}

		if err := s.GuestServiceRepository.CompleteServiceRequest(uint(requestID), completedBy); err != nil {
			response.JSON(c, "Failed to complete request", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Request completed", http.StatusOK, nil, nil)
	}
}

// handleGetAssignedTasks returns all service requests that have been assigned to staff
func (s *Server) handleGetAssignedTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		requests, err := s.GuestServiceRepository.GetAssignedServiceRequests()
		if err != nil {
			log.Printf("handleGetAssignedTasks: error fetching assigned tasks: %v", err)
			response.JSON(c, "Failed to fetch assigned tasks", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Assigned tasks retrieved", http.StatusOK, requests, nil)
	}
}

// handleUpdateTaskStatus updates the status of an assigned task
func (s *Server) handleUpdateTaskStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid request ID", http.StatusBadRequest, nil, err)
			return
		}

		var req struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		// Validate status
		validStatuses := map[string]bool{
			"assigned":    true,
			"in_progress": true,
			"completed":   true,
		}

		if !validStatuses[req.Status] {
			response.JSON(c, "Invalid status. Must be: assigned, in_progress, or completed", http.StatusBadRequest, nil, nil)
			return
		}

		// Get the request to check if it exists
		request, err := s.GuestServiceRepository.GetServiceRequestByID(uint(requestID))
		if err != nil {
			response.JSON(c, "Task not found", http.StatusNotFound, nil, err)
			return
		}

		// If completing, use the complete flow
		if req.Status == "completed" {
			completedBy := ""
			if request.AssignedTo != "" {
				completedBy = request.AssignedTo
			}

			// Free up staff if assigned
			if request.AssignedStaffID != nil {
				if err := s.StaffRepository.CompleteTask(*request.AssignedStaffID); err != nil {
					log.Printf("handleUpdateTaskStatus: error freeing staff: %v", err)
				}
			} else {
				// Fallback: find staff who has this task as current_task_id
				allStaff, _, _ := s.StaffRepository.GetAllStaff(1, 100, "", "", "")
				for _, staff := range allStaff {
					if staff.CurrentTaskID != nil && *staff.CurrentTaskID == uint(requestID) {
						if err := s.StaffRepository.CompleteTask(staff.ID); err != nil {
							log.Printf("handleUpdateTaskStatus: error freeing staff (fallback): %v", err)
						}
						break
					}
				}
			}

			if err := s.GuestServiceRepository.CompleteServiceRequest(uint(requestID), completedBy); err != nil {
				response.JSON(c, "Failed to complete task", http.StatusInternalServerError, nil, err)
				return
			}
		} else {
			// Update status for assigned or in_progress
			if err := s.GuestServiceRepository.UpdateServiceRequestStatus(uint(requestID), req.Status, ""); err != nil {
				response.JSON(c, "Failed to update task status", http.StatusInternalServerError, nil, err)
				return
			}
		}

		// Get updated request
		updatedRequest, _ := s.GuestServiceRepository.GetServiceRequestByID(uint(requestID))

		response.JSON(c, "Task status updated", http.StatusOK, updatedRequest, nil)
	}
}
