package database

import (
	"appointment-booking/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the PostgreSQL database connection and performs migrations
func InitDB() {
	var err error

	// PostgreSQL connection string
	// Format: postgres://username:password@host:port/dbname
	// Environment variables (optional):
	// DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME
	databaseUser := getEnvOrDefault("DB_USER", "postgres")
	databasePassword := getEnvOrDefault("DB_PASSWORD", "root")
	databaseHost := getEnvOrDefault("DB_HOST", "localhost")
	databasePort := getEnvOrDefault("DB_PORT", "5432")
	databaseName := getEnvOrDefault("DB_NAME", "appointment-booking")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		databaseHost, databaseUser, databasePassword, databaseName, databasePort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&models.Coach{},
		&models.CoachAvailability{},
		&models.User{},
		&models.Booking{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	log.Println("Database migrations completed!")

	// Seed data if tables are empty
	seedData()
}

// seedData inserts initial data into the database
func seedData() {
	// Check if coaches already exist
	var coachCount int64
	DB.Model(&models.Coach{}).Count(&coachCount)
	if coachCount > 0 {
		return // Already seeded
	}

	// Create sample coaches
	coaches := []models.Coach{
		{Name: "Coach A", Timezone: "Asia/Kolkata"},
		{Name: "Coach B", Timezone: "America/New_York"},
	}

	for _, coach := range coaches {
		if err := DB.Create(&coach).Error; err != nil {
			log.Printf("Error seeding coach: %v", err)
		}
	}

	// Create sample users
	users := []models.User{
		{Name: "Alice", Timezone: "America/New_York"},
		{Name: "Bob", Timezone: "Asia/Kolkata"},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Error seeding user: %v", err)
		}
	}

	// Create availability for Coach A: Monday 10:00-15:00, Wednesday 09:00-12:00
	availabilities := []models.CoachAvailability{
		{
			CoachID:   1,
			DayOfWeek: "Monday",
			StartTime: "10:00",
			EndTime:   "15:00",
		},
		{
			CoachID:   1,
			DayOfWeek: "Wednesday",
			StartTime: "09:00",
			EndTime:   "12:00",
		},
		// Coach B: Tuesday 09:00-14:00, Friday 13:00-17:00
		{
			CoachID:   2,
			DayOfWeek: "Tuesday",
			StartTime: "09:00",
			EndTime:   "14:00",
		},
		{
			CoachID:   2,
			DayOfWeek: "Friday",
			StartTime: "13:00",
			EndTime:   "17:00",
		},
	}

	for _, avail := range availabilities {
		if err := DB.Create(&avail).Error; err != nil {
			log.Printf("Error seeding availability: %v", err)
		}
	}

	log.Println("Seed data inserted successfully!")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// getEnvOrDefault retrieves environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
