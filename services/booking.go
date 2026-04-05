package services

import (
	"appointment-booking/database"
	"appointment-booking/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateBooking creates a new booking with concurrency handling using transactions
func CreateBooking(userID, coachID uint, slotTime time.Time) (uint, error) {
	// Validate user exists
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}

	// Validate coach exists
	var coach models.Coach
	if err := database.DB.First(&coach, coachID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("coach not found")
		}
		return 0, err
	}

	// Validate slot is within coach availability
	isWithin, err := IsSlotWithinAvailability(coachID, slotTime)
	if err != nil || !isWithin {
		return 0, fmt.Errorf("slot not within coach availability")
	}

	// Use transaction with row-level locking to prevent race conditions
	var bookingID uint
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Check if slot is already booked (with row lock)
		var existingBooking models.Booking
		result := tx.Clauses().
			Where("coach_id = ? AND slot_time = ? AND status = ?",
				coachID, slotTime, "confirmed").
			First(&existingBooking)

		if result.Error == nil {
			// Booking exists for this slot
			return fmt.Errorf("slot already booked")
		} else if result.Error != gorm.ErrRecordNotFound {
			// Database error
			return result.Error
		}

		// Create new booking
		booking := models.Booking{
			UserID:   userID,
			CoachID:  coachID,
			SlotTime: slotTime,
			Status:   "confirmed",
		}

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		bookingID = booking.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return bookingID, nil
}

// GetUserBookings retrieves all confirmed bookings for a user
func GetUserBookings(userID uint) ([]models.BookingResponse, error) {
	// Validate user exists
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	var bookings []models.Booking
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "confirmed").
		Preload("Coach").
		Order("slot_time ASC").
		Find(&bookings).Error; err != nil {
		return nil, err
	}

	responses := []models.BookingResponse{}
	for _, booking := range bookings {
		responses = append(responses, models.BookingResponse{
			BookingID: booking.ID,
			CoachID:   booking.CoachID,
			CoachName: booking.Coach.Name,
			SlotTime:  booking.SlotTime,
			Status:    booking.Status,
		})
	}

	return responses, nil
}

// CancelBooking cancels a booking by setting status to cancelled
func CancelBooking(bookingID, userID uint) error {
	// Find booking and verify it belongs to the user
	var booking models.Booking
	if err := database.DB.First(&booking, bookingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("booking not found")
		}
		return err
	}

	if booking.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Update status to cancelled
	if err := database.DB.Model(&booking).Update("status", "cancelled").Error; err != nil {
		return err
	}

	return nil
}
