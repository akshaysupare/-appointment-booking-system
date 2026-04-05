package models

import "time"

// Coach represents a coach/provider in the system
type Coach struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Timezone  string    `json:"timezone" binding:"required"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// CoachAvailability defines when a coach is available
type CoachAvailability struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CoachID   uint      `gorm:"index;not null" json:"coach_id" binding:"required"`
	Coach     Coach     `gorm:"foreignKey:CoachID" json:"-"`
	DayOfWeek string    `json:"day_of_week" binding:"required"` // Monday, Tuesday, etc.
	StartTime string    `json:"start_time" binding:"required"`  // HH:MM format
	EndTime   string    `json:"end_time" binding:"required"`    // HH:MM format
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// User represents a customer/user booking appointments
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Timezone  string    `json:"timezone" binding:"required"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Booking represents a confirmed or cancelled booking
type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id" binding:"required"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	CoachID   uint      `gorm:"index;not null" json:"coach_id" binding:"required"`
	Coach     Coach     `gorm:"foreignKey:CoachID" json:"-"`
	SlotTime  time.Time `gorm:"index;not null" json:"slot_time" binding:"required"` // Stored in UTC
	Status    string    `gorm:"index" json:"status" binding:"required"`             // confirmed, cancelled
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// SetAvailabilityRequest handles coach availability setting
type SetAvailabilityRequest struct {
	CoachID   uint   `json:"coach_id" binding:"required"`
	Day       string `json:"day" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// CreateBookingRequest handles booking creation
type CreateBookingRequest struct {
	UserID   uint      `json:"user_id" binding:"required"`
	CoachID  uint      `json:"coach_id" binding:"required"`
	DateTime time.Time `json:"datetime" binding:"required"`
}

// BookingResponse for GET bookings
type BookingResponse struct {
	BookingID uint      `json:"booking_id"`
	CoachID   uint      `json:"coach_id"`
	CoachName string    `json:"coach_name"`
	SlotTime  time.Time `json:"slot_time"`
	Status    string    `json:"status"`
}

// ErrorResponse for consistent error handling
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse for successful operations
type SuccessResponse struct {
	Message string `json:"message"`
}

// BookingIDResponse for booking creation
type BookingIDResponse struct {
	Message   string `json:"message"`
	BookingID uint   `json:"booking_id"`
}
