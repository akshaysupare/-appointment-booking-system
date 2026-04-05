package main

import (
	"appointment-booking/database"
	"appointment-booking/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Create Gin engine
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Add CORS middleware (optional but good for production)
	router.Use(corsMiddleware())

	log.Println("Server starting on :8080...")
	log.Println("Database: appointments.db")

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
