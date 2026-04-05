package handlers

import (
	"appointment-booking/database"
	"appointment-booking/models"
	"appointment-booking/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAvailableSlots returns all available 30-minute slots for a coach on a given date
// GET /users/slots?coach_id=1&date=2025-10-28
func GetAvailableSlots(c *gin.Context) {
	coachIDStr := c.Query("coach_id")
	date := c.Query("date")

	// Validate required query params
	if coachIDStr == "" || date == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing required query parameters: coach_id, date",
		})
		return
	}

	// Parse coach ID
	coachID, err := strconv.ParseUint(coachIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid coach_id format",
		})
		return
	}

	// Validate coach exists
	var coach models.Coach
	if err := database.DB.First(&coach, uint(coachID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Coach not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Get available slots
	slots, err := services.GetAvailableSlots(uint(coachID), date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if slots == nil {
		slots = []string{}
	}

	c.JSON(http.StatusOK, slots)
}

// CreateBooking books a 30-minute slot for a user
// POST /users/bookings
func CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	// Create booking (handles validation and concurrency)
	bookingID, err := services.CreateBooking(req.UserID, req.CoachID, req.DateTime)
	if err != nil {
		// Determine HTTP status code based on error message
		statusCode := http.StatusInternalServerError
		errMsg := err.Error()

		if errMsg == "user not found" || errMsg == "coach not found" {
			statusCode = http.StatusNotFound
		} else if errMsg == "slot already booked" {
			statusCode = http.StatusConflict
		} else if errMsg == "slot not within coach availability" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error: errMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, models.BookingIDResponse{
		Message:   "Booking confirmed",
		BookingID: bookingID,
	})
}

// GetUserBookings retrieves all confirmed bookings for a user
// GET /users/bookings?user_id=101
func GetUserBookings(c *gin.Context) {
	userIDStr := c.Query("user_id")

	// Validate required query param
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing required query parameter: user_id",
		})
		return
	}

	// Parse user ID
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user_id format",
		})
		return
	}

	// Get user bookings
	bookings, err := services.GetUserBookings(uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if bookings == nil {
		bookings = []models.BookingResponse{}
	}

	c.JSON(http.StatusOK, bookings)
}

// CancelBooking cancels a booking
// DELETE /users/bookings/:booking_id?user_id=101
func CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	userIDStr := c.Query("user_id")

	// Validate required params
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing required query parameter: user_id",
		})
		return
	}

	// Parse IDs
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid booking_id format",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user_id format",
		})
		return
	}

	// Cancel booking
	err = services.CancelBooking(uint(bookingID), uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errMsg := err.Error()

		if errMsg == "booking not found" {
			statusCode = http.StatusNotFound
		} else if errMsg == "unauthorized" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error: errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Booking cancelled successfully",
	})
}
