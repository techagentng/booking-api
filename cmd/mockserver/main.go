package main

import (
	"log"

	"hotel/tripsbook/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Create mock handler
	mockHandler := handlers.NewMockAPIHandler()

	// TripsBook Public API Routes
	public := r.Group("/api/v1/public")
	{
		// Categories
		public.GET("/categories", mockHandler.GetCategories)

		// Services by category
		public.GET("/services/category/:category", mockHandler.GetServicesByCategory)

		// Explore endpoints
		public.GET("/explore/featured", mockHandler.GetExploreFeatured)
		public.GET("/explore/destinations", mockHandler.GetExploreDestinations)

		// Nearby services
		public.GET("/services/nearby", mockHandler.GetNearbyServices)
		public.GET("/services/nearby/filters", mockHandler.GetDistanceFilters)

		// Trending
		public.GET("/services/trending", mockHandler.GetTrendingServices)
		public.GET("/categories/trending", mockHandler.GetTrendingCategories)

		// Service Providers
		public.GET("/providers", mockHandler.GetAllProviders)
		public.GET("/providers/category/:category", mockHandler.GetProvidersByCategory)
		public.GET("/providers/category/:category/featured", mockHandler.GetFeaturedProviders)
		public.GET("/providers/:id", mockHandler.GetProviderByID)

		// Search
		public.GET("/search", mockHandler.SearchServices)
		public.GET("/search/suggestions", mockHandler.GetSearchSuggestions)

		// Location services
		public.GET("/location/current", mockHandler.GetCurrentLocation)
		public.POST("/location/update", mockHandler.UpdateLocation)
		public.GET("/locations/popular", mockHandler.GetPopularLocations)

		// Service details
		public.GET("/services/:id", mockHandler.GetServiceDetails)

		// Map view
		public.GET("/services/map", mockHandler.GetMapViewServices)

		// Analytics tracking
		public.POST("/analytics/track", mockHandler.TrackAnalytics)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "TripsBook Mock API Server is running",
		})
	})

	// Start server
	log.Println("🚀 TripsBook Mock API Server starting on http://localhost:8081")
	log.Println("📱 Available endpoints:")
	log.Println("   GET /api/v1/public/categories")
	log.Println("   GET /api/v1/public/services/trending")
	log.Println("   GET /api/v1/public/categories/trending")
	log.Println("   GET /api/v1/public/explore/featured")
	log.Println("   GET /api/v1/public/services/nearby")
	log.Println("   GET /api/v1/public/search")
	log.Println("   And more...")
	log.Println("🔗 Frontend can now connect to localhost:8081")

	if err := r.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
