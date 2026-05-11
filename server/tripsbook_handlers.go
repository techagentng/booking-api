package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hotel/tripsbook/handlers"
)

// TripsBookHandlers contains the TripsBook public API handlers
type TripsBookHandlers struct {
	publicHandler *handlers.PublicAPIHandler
}

// NewTripsBookHandlers creates new TripsBook handlers
func NewTripsBookHandlers(db *gorm.DB) *TripsBookHandlers {
	return &TripsBookHandlers{
		publicHandler: handlers.NewPublicAPIHandler(db),
	}
}

// Add TripsBook handlers to the server
func (s *Server) InitTripsBookHandlers() {
	// Initialize TripsBook handlers
	tripsbookHandlers := NewTripsBookHandlers(s.DB)

	// Store handlers for use in router
	s.tripsbookHandlers = tripsbookHandlers
}

// Add tripsbookHandlers field to Server struct
func (s *Server) handleGetCategories() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetCategories
}

func (s *Server) handleGetServicesByCategory() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetServicesByCategory
}

func (s *Server) handleGetExploreFeatured() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetExploreFeatured
}

func (s *Server) handleGetExploreDestinations() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetExploreDestinations
}

func (s *Server) handleGetNearbyServices() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetNearbyServices
}

func (s *Server) handleGetDistanceFilters() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetDistanceFilters
}

func (s *Server) handleGetTrendingServices() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetTrendingServices
}

func (s *Server) handleGetTrendingCategories() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetTrendingCategories
}

func (s *Server) handleSearchServices() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.SearchServices
}

func (s *Server) handleGetSearchSuggestions() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetSearchSuggestions
}

func (s *Server) handleGetCurrentLocation() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetCurrentLocation
}

func (s *Server) handleUpdateLocation() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.UpdateLocation
}

func (s *Server) handleGetPopularLocations() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetPopularLocations
}

func (s *Server) handleGetServiceDetails() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetServiceDetails
}

func (s *Server) handleGetMapViewServices() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.GetMapViewServices
}

func (s *Server) handleTrackAnalytics() gin.HandlerFunc {
	return s.tripsbookHandlers.publicHandler.TrackAnalytics
}
