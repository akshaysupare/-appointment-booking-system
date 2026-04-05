package services

import (
	"appointment-booking/database"
	"appointment-booking/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// GetAvailableSlots returns all available 30-minute slots for a coach on a given date
func GetAvailableSlots(coachID uint, date string) ([]string, error) {
	// Parse the date
	slotDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	// Get the day of week
	dayOfWeek := slotDate.Weekday().String()

	// Find coach availability for that day
	var availability models.CoachAvailability
	if err := database.DB.Where("coach_id = ? AND day_of_week = ?", coachID, dayOfWeek).
		First(&availability).Error; err != nil {
		// No availability for this day
		return []string{}, nil
	}

	// Parse start and end times
	startTimeParts := strings.Split(availability.StartTime, ":")
	startHour, _ := strconv.Atoi(startTimeParts[0])
	startMin, _ := strconv.Atoi(startTimeParts[1])

	endTimeParts := strings.Split(availability.EndTime, ":")
	endHour, _ := strconv.Atoi(endTimeParts[0])
	endMin, _ := strconv.Atoi(endTimeParts[1])

	// Generate all 30-minute slots
	allSlots := []time.Time{}
	currentTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), startHour, startMin, 0, 0, time.UTC)
	endTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), endHour, endMin, 0, 0, time.UTC)

	for currentTime.Before(endTime) {
		allSlots = append(allSlots, currentTime)
		currentTime = currentTime.Add(30 * time.Minute)
	}

	// Get booked slots for this coach on this date
	bookedSlots := make(map[time.Time]bool)
	var bookings []models.Booking
	dayStart := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	if err := database.DB.Where("coach_id = ? AND slot_time >= ? AND slot_time < ? AND status = ?",
		coachID, dayStart, dayEnd, "confirmed").Find(&bookings).Error; err != nil {
		log.Printf("Error fetching bookings: %v", err)
	}

	for _, booking := range bookings {
		bookedSlots[booking.SlotTime] = true
	}

	// Filter available slots (30-minute precision)
	availableSlots := []string{}
	for _, slot := range allSlots {
		if !bookedSlots[slot] {
			availableSlots = append(availableSlots, slot.Format(time.RFC3339))
		}
	}

	return availableSlots, nil
}

// IsSlotWithinAvailability checks if a given slot is within coach's availability
func IsSlotWithinAvailability(coachID uint, slotTime time.Time) (bool, error) {
	dayOfWeek := slotTime.Weekday().String()

	// Find coach availability for that day
	var availability models.CoachAvailability
	if err := database.DB.Where("coach_id = ? AND day_of_week = ?", coachID, dayOfWeek).
		First(&availability).Error; err != nil {
		return false, fmt.Errorf("coach not available on %s", dayOfWeek)
	}

	// Parse times
	startTimeParts := strings.Split(availability.StartTime, ":")
	startHour, _ := strconv.Atoi(startTimeParts[0])
	startMin, _ := strconv.Atoi(startTimeParts[1])

	endTimeParts := strings.Split(availability.EndTime, ":")
	endHour, _ := strconv.Atoi(endTimeParts[0])
	endMin, _ := strconv.Atoi(endTimeParts[1])

	// Create times for the same date
	startTime := time.Date(slotTime.Year(), slotTime.Month(), slotTime.Day(),
		startHour, startMin, 0, 0, time.UTC)
	endTime := time.Date(slotTime.Year(), slotTime.Month(), slotTime.Day(),
		endHour, endMin, 0, 0, time.UTC)

	// Check if slot is within range
	if slotTime.Before(startTime) || slotTime.After(endTime) || slotTime.Equal(endTime) {
		return false, fmt.Errorf("slot not within coach availability")
	}

	return true, nil
}
