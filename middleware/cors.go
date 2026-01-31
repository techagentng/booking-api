package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware configuration
func CORSConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the origin from the request
		origin := c.Request.Header.Get("Origin")
		
		// List of allowed origins (add your deployed frontend URLs)
		allowedOrigins := []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"https://booking-app-08gs.onrender.com",
		}
		
		// Check if the origin is allowed
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}
		
		// Set CORS headers
		if isAllowed || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		// Allow specific methods
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		
		// Allow specific headers
		c.Header("Access-Control-Allow-Headers", 
			"Origin, Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token")
		
		// Allow credentials
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		
		c.Next()
	}
}
