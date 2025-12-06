package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"hotel/models"
	"hotel/server/response"
	"hotel/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	// Import the errors package for enhanced error handling
)

// handleGetRooms returns all rooms with pagination, search, and status filter
func (s *Server) handleGetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		pageSize := 10
		search := c.Query("search")
		status := c.Query("status")

		if p := c.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				page = parsed
			}
		}

		if ps := c.Query("page_size"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
				pageSize = parsed
			}
		}

		roomService := services.NewRoomService(s.RoomRepository)
		result, err := roomService.GetAllRooms(page, pageSize, search, status)
		if err != nil {
			log.Printf("handleGetRooms: error fetching rooms: %v", err)
			response.JSON(c, "Failed to fetch rooms", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Rooms retrieved successfully", http.StatusOK, result, nil)
	}
}

// handleGetRoomByID returns a room by ID
func (s *Server) handleGetRoomByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid room ID", http.StatusBadRequest, nil, err)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		room, err := roomService.GetRoomByID(uint(id))
		if err != nil {
			log.Printf("handleGetRoomByID: error fetching room: %v", err)
			response.JSON(c, "Room not found", http.StatusNotFound, nil, err)
			return
		}

		response.JSON(c, "Room retrieved successfully", http.StatusOK, room, nil)
	}
}

// handleCreateRoom creates a new room
func (s *Server) handleCreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateRoomRequest

		// Log raw request body for debugging
		body, _ := c.GetRawData()
		log.Printf("Create room request body: %s", string(body))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the request body

		if err := c.ShouldBindJSON(&req); err != nil {
			errMsg := "Invalid request: "
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				for _, fieldErr := range validationErrs {
					errMsg += fmt.Sprintf("Field '%s' failed on '%s' validation; ", fieldErr.Field(), fieldErr.Tag())
				}
			} else {
				errMsg += err.Error()
			}
			log.Printf("Validation error: %s", errMsg)
			response.JSON(c, errMsg, http.StatusBadRequest, nil, err)
			return
		}

		log.Printf("Creating room with data: %+v", req)

		roomService := services.NewRoomService(s.RoomRepository)
		room, err := roomService.CreateRoom(&req)
		if err != nil {
			errMsg := fmt.Sprintf("Room creation failed: %v", err)
			log.Println(errMsg)
			if err.Error() == "room number already exists" {
				response.JSON(c, errMsg, http.StatusConflict, nil, err)
				return
			}
			response.JSON(c, errMsg, http.StatusBadRequest, nil, err)
			return
		}

		log.Printf("Room created successfully: %+v", room)
		response.JSON(c, "Room created successfully", http.StatusCreated, room, nil)
	}
}

// handleUpdateRoom updates a room
func (s *Server) handleUpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid room ID", http.StatusBadRequest, nil, err)
			return
		}

		var req models.UpdateRoomRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		room, err := roomService.UpdateRoom(uint(id), &req)
		if err != nil {
			log.Printf("handleUpdateRoom: error updating room: %v", err)
			if err.Error() == "room number already exists" {
				response.JSON(c, "Room update failed", http.StatusConflict, nil, err)
				return
			}
			response.JSON(c, "Room update failed", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Room updated successfully", http.StatusOK, room, nil)
	}
}

// handleDeleteRoom deletes a room
func (s *Server) handleDeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid room ID", http.StatusBadRequest, nil, err)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		err = roomService.DeleteRoom(uint(id))
		if err != nil {
			log.Printf("handleDeleteRoom: error deleting room: %v", err)
			response.JSON(c, "Room deletion failed", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Room deleted successfully", http.StatusOK, nil, nil)
	}
}

// handleCheckRoomAvailability checks room availability for a date range (POST version)
func (s *Server) handleCheckRoomAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RoomAvailabilityRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		rooms, err := roomService.CheckAvailability(&req)
		if err != nil {
			log.Printf("handleCheckRoomAvailability: error checking availability: %v", err)
			response.JSON(c, "Failed to check availability", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Available rooms retrieved", http.StatusOK, rooms, nil)
	}
}

// handleGetAvailableRooms checks room availability using query params (GET version)
func (s *Server) handleGetAvailableRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkIn := c.Query("check_in")
		checkOut := c.Query("check_out")
		roomType := c.Query("room_type")

		if checkIn == "" || checkOut == "" {
			response.JSON(c, "check_in and check_out are required", http.StatusBadRequest, nil, nil)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		rooms, err := roomService.GetAvailableRoomsQuery(checkIn, checkOut, roomType)
		if err != nil {
			log.Printf("handleGetAvailableRooms: error checking availability: %v", err)
			response.JSON(c, "Failed to check availability", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Available rooms retrieved successfully", http.StatusOK, rooms, nil)
	}
}

// handleGetRoomsByType returns all rooms of a specific type
func (s *Server) handleGetRoomsByType() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomType := c.Param("type")
		if roomType == "" {
			response.JSON(c, "Room type is required", http.StatusBadRequest, nil, nil)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		rooms, err := roomService.GetRoomsByType(roomType)
		if err != nil {
			log.Printf("handleGetRoomsByType: error fetching rooms: %v", err)
			response.JSON(c, "Failed to fetch rooms", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Rooms retrieved successfully", http.StatusOK, rooms, nil)
	}
}

// handleUpdateRoomStatus updates a room's status
func (s *Server) handleUpdateRoomStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			response.JSON(c, "Invalid room ID", http.StatusBadRequest, nil, err)
			return
		}

		var req struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSON(c, "Invalid request body", http.StatusBadRequest, nil, err)
			return
		}

		roomService := services.NewRoomService(s.RoomRepository)
		err = roomService.UpdateRoomStatus(uint(id), req.Status)
		if err != nil {
			log.Printf("handleUpdateRoomStatus: error updating status: %v", err)
			response.JSON(c, "Failed to update room status", http.StatusBadRequest, nil, err)
			return
		}

		response.JSON(c, "Room status updated successfully", http.StatusOK, nil, nil)
	}
}

// handleGetRoomTypeStats returns statistics for each room type
func (s *Server) handleGetRoomTypeStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomService := services.NewRoomService(s.RoomRepository)
		stats, err := roomService.GetRoomTypeStats()
		if err != nil {
			log.Printf("handleGetRoomTypeStats: error fetching stats: %v", err)
			response.JSON(c, "Failed to fetch room statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Room statistics retrieved successfully", http.StatusOK, stats, nil)
	}
}

// handleGetRoomStats returns overall room statistics (total, available, occupied, maintenance)
func (s *Server) handleGetRoomStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := s.RoomRepository.GetRoomStats()
		if err != nil {
			log.Printf("handleGetRoomStats: error fetching stats: %v", err)
			response.JSON(c, "Failed to fetch room statistics", http.StatusInternalServerError, nil, err)
			return
		}

		response.JSON(c, "Room statistics retrieved successfully", http.StatusOK, stats, nil)
	}
}
