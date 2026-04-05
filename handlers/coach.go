package handlers

import (
	"appointment-booking/database"
	"appointment-booking/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetCoachAvailability sets a coach's weekly availability
// POST /coaches/availability
func SetCoachAvailability(c *gin.Context) {
	var req models.SetAvailabilityRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate coach exists
	var coach models.Coach
	if err := database.DB.First(&coach, req.CoachID).Error; err != nil {
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

	// Validate day is a valid weekday
	validDays := map[string]bool{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	if !validDays[req.Day] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid day. Use valid weekday (Monday-Sunday)",
		})
		return
	}

	// Validate time format HH:MM
	if !isValidTimeFormat(req.StartTime) || !isValidTimeFormat(req.EndTime) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid time format. Use HH:MM (e.g., 09:00)",
		})
		return
	}

	// Validate start_time < end_time
	if !isTimeBefore(req.StartTime, req.EndTime) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "start_time must be before end_time",
		})
		return
	}

	// Create or update availability
	availability := models.CoachAvailability{
		CoachID:   req.CoachID,
		DayOfWeek: req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Check if availability already exists for this day
	var existing models.CoachAvailability
	result := database.DB.Where("coach_id = ? AND day_of_week = ?", req.CoachID, req.Day).
		First(&existing)

	if result.Error == nil {
		// Update existing
		if err := database.DB.Model(&existing).Updates(&availability).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to update availability",
			})
			return
		}
	} else if result.Error == gorm.ErrRecordNotFound {
		// Create new
		if err := database.DB.Create(&availability).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to create availability",
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Availability set successfully",
	})
}

// isValidTimeFormat validates HH:MM format
func isValidTimeFormat(timeStr string) bool {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}

	// Simple validation: both should be 2 digits
	if len(parts[0]) != 2 || len(parts[1]) != 2 {
		return false
	}

	// Check if both are numbers
	for _, part := range parts {
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}

	return true
}

// isTimeBefore checks if time1 is before time2 (both in HH:MM format)
func isTimeBefore(time1, time2 string) bool {
	parts1 := strings.Split(time1, ":")
	parts2 := strings.Split(time2, ":")

	h1, h2 := 0, 0
	m1, m2 := 0, 0

	// Parse time1
	for _, c := range parts1[0] {
		h1 = h1*10 + int(c-'0')
	}
	for _, c := range parts1[1] {
		m1 = m1*10 + int(c-'0')
	}

	// Parse time2
	for _, c := range parts2[0] {
		h2 = h2*10 + int(c-'0')
	}
	for _, c := range parts2[1] {
		m2 = m2*10 + int(c-'0')
	}

	if h1 < h2 {
		return true
	}
	if h1 == h2 && m1 < m2 {
		return true
	}

	return false
}
