package routes

import (
	"appointment-booking/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all API routes
func SetupRoutes(router *gin.Engine) {
	// Coach routes
	router.POST("/coaches/availability", handlers.SetCoachAvailability)

	// User routes
	router.GET("/users/slots", handlers.GetAvailableSlots)
	router.POST("/users/bookings", handlers.CreateBooking)
	router.GET("/users/bookings", handlers.GetUserBookings)
	router.DELETE("/users/bookings/:booking_id", handlers.CancelBooking)
}
