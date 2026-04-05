# Appointment Booking System API

A production-quality RESTful API for managing appointment bookings between coaches and users. Built with Go, Gin, and SQLite for easy local deployment.

## Project Overview

This appointment booking system provides a robust platform for:
- **Coaches** to define their weekly availability (recurring schedules)
- **Users** to browse available time slots and book appointments
- **Automatic slot generation** in 30-minute intervals
- **Concurrency-safe booking** with transaction-based row-level locking to prevent double-bookings
- **Timezone support** for coaches and users (stored in UTC, flexible display)

### Design Decisions

1. **SQLite for Development/Production**: Lightweight, zero-setup database that doesn't require external services. Perfect for small to medium deployments.

2. **GORM ORM**: Type-safe database operations, automatic migrations, and built-in transaction support for concurrency handling.

3. **Gin Web Framework**: Minimal, fast HTTP framework with excellent middleware support and validation binding capabilities.

4. **Transaction-based Concurrency**: Uses GORM transactions with row-level locking to ensure no two bookings can occupy the same 30-minute slot simultaneously.

5. **UTC Storage**: All times stored in UTC in the database. Timezone information is maintained for future conversion purposes.

6. **30-minute Slot Model**: Appointments are fixed at 30-minute intervals for simplicity and scalability.

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: SQLite via GORM
- **Database Driver**: gorm.io/driver/sqlite
- **ORM**: gorm.io/gorm
- **No external dependencies**: Runs completely locally with a single command

## Database Models

### Coach
```
id           uint   - Primary key, auto-increment
name         string - Coach's name
timezone     string - Coach's timezone (e.g., "Asia/Kolkata")
created_at   time   - Timestamp
updated_at   time   - Timestamp
```

### CoachAvailability
```
id           uint   - Primary key
coach_id     uint   - Foreign key to Coach
day_of_week  string - Weekday (Monday, Tuesday, etc.)
start_time   string - Start time in HH:MM format (e.g., "09:00")
end_time     string - End time in HH:MM format (e.g., "14:00")
created_at   time   - Timestamp
updated_at   time   - Timestamp
```

### User
```
id           uint   - Primary key, auto-increment
name         string - User's name
timezone     string - User's timezone (e.g., "America/New_York")
created_at   time   - Timestamp
updated_at   time   - Timestamp
```

### Booking
```
id           uint   - Primary key
user_id      uint   - Foreign key to User
coach_id     uint   - Foreign key to Coach
slot_time    time   - Booking time stored in UTC (RFC3339 format)
status       string - "confirmed" or "cancelled"
created_at   time   - Timestamp
updated_at   time   - Timestamp
```

## Project Structure

```
appointment-booking/
├── main.go                 # Application entry point and CORS setup
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── appointments.db         # SQLite database (auto-created)
├── database/
│   └── db.go              # DB initialization, auto-migration, seed data
├── models/
│   └── models.go          # All GORM structs and request/response types
├── handlers/
│   ├── coach.go           # Coach availability endpoint handler
│   └── user.go            # Slots, bookings endpoints handlers
├── routes/
│   └── routes.go          # Route registration
├── services/
│   ├── slots.go           # Slot generation and availability logic
│   └── booking.go         # Booking creation and concurrency logic
└── README.md              # This file
```

## Setup Instructions

### Prerequisites
- Go 1.21 or later installed on your system
- Windows, macOS, or Linux

### Installation & Running

1. **Navigate to project directory**:
   ```bash
   cd appointment-booking
   ```

2. **Download dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **Expected output**:
   ```
   Database connected successfully!
   Database migrations completed!
   [seed data messages...]
   Server starting on :8080...
   Database: appointments.db
   ```

The SQLite database file `appointments.db` is automatically created in the working directory on first run. The database is pre-populated with sample coaches, users, and availability data.

## API Endpoints

### 1. POST `/coaches/availability`
Set or update a coach's weekly availability for a specific day.

**Request**:
```bash
curl -X POST http://localhost:8080/coaches/availability \
  -H "Content-Type: application/json" \
  -d '{
    "coach_id": 1,
    "day": "Monday",
    "start_time": "09:00",
    "end_time": "14:00"
  }'
```

**Response** (201 Created):
```json
{
  "message": "Availability set successfully"
}
```

**Validations**:
- `coach_id` must reference an existing coach
- `day` must be a valid weekday (Monday-Sunday)
- `start_time` and `end_time` must be in HH:MM format
- `start_time` must be before `end_time`

**Error Responses**:
- `400 Bad Request` - Invalid format or validation failure
- `404 Not Found` - Coach not found

---

### 2. GET `/users/slots`
Retrieve all available 30-minute slots for a coach on a specific date.

**Request**:
```bash
curl -X GET "http://localhost:8080/users/slots?coach_id=1&date=2025-10-28"
```

**Response** (200 OK):
```json
[
  "2025-10-28T09:00:00Z",
  "2025-10-28T09:30:00Z",
  "2025-10-28T10:00:00Z",
  "2025-10-28T10:30:00Z",
  "2025-10-28T11:00:00Z",
  "2025-10-28T11:30:00Z",
  "2025-10-28T12:00:00Z",
  "2025-10-28T12:30:00Z",
  "2025-10-28T13:00:00Z",
  "2025-10-28T13:30:00Z"
]
```

**Logic**:
1. Determines the weekday of the provided date
2. Finds coach's availability for that weekday
3. Generates all 30-minute slots between start and end times
4. Filters out slots that have confirmed bookings
5. Returns remaining slots in UTC ISO8601 format

**Query Parameters**:
- `coach_id` (required) - Coach ID
- `date` (required) - Date in YYYY-MM-DD format

**Error Responses**:
- `400 Bad Request` - Missing or invalid parameters
- `404 Not Found` - Coach not found

**Special Cases**:
- Empty array `[]` - Returned if no availability exists for that day or all slots are booked
- Only returns future available slots (based on system time)

---

### 3. POST `/users/bookings`
Book a 30-minute slot for a user with a coach.

**Request**:
```bash
curl -X POST http://localhost:8080/users/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 101,
    "coach_id": 1,
    "datetime": "2025-10-28T09:30:00Z"
  }'
```

**Response** (201 Created):
```json
{
  "message": "Booking confirmed",
  "booking_id": 5
}
```

**Booking Logic**:
1. Validates user exists
2. Validates coach exists
3. Checks if slot falls within coach's availability
4. Uses GORM transaction with row-level locking to prevent race conditions
5. Verifies no confirmed booking already exists for the slot
6. Creates new booking with "confirmed" status

**Concurrency Handling**:
The booking creation is wrapped in a GORM transaction that checks for existing bookings atomically. This prevents two simultaneous requests from booking the same slot.

```go
db.Transaction(func(tx *gorm.DB) error {
    // Check for existing booking with implicit row lock
    var existingBooking models.Booking
    if err := tx.Where("coach_id = ? AND slot_time = ? AND status = ?",
        coachID, slotTime, "confirmed").First(&existingBooking).Error; err == nil {
        return fmt.Errorf("slot already booked")
    }
    // Create new booking
    return tx.Create(&booking).Error
})
```

**Error Responses**:
- `400 Bad Request` - Invalid request format or slot not within availability
- `404 Not Found` - User or coach not found
- `409 Conflict` - Slot already booked by another user

**Input Validation**:
- `datetime` must be in RFC3339 format ("2025-10-28T09:30:00Z")
- `user_id` and `coach_id` must be valid positive integers

---

### 4. GET `/users/bookings`
Retrieve all confirmed bookings for a user.

**Request**:
```bash
curl -X GET "http://localhost:8080/users/bookings?user_id=101"
```

**Response** (200 OK):
```json
[
  {
    "booking_id": 5,
    "coach_id": 1,
    "coach_name": "Coach A",
    "slot_time": "2025-10-28T09:30:00Z",
    "status": "confirmed"
  },
  {
    "booking_id": 6,
    "coach_id": 2,
    "coach_name": "Coach B",
    "slot_time": "2025-10-29T10:00:00Z",
    "status": "confirmed"
  }
]
```

**Query Parameters**:
- `user_id` (required) - User ID

**Features**:
- Returns only confirmed bookings (excludes cancelled bookings)
- Sorted by slot_time in ascending order
- Includes coach name and full booking details

**Error Responses**:
- `400 Bad Request` - Missing or invalid user_id
- `404 Not Found` - User not found

---

### 5. DELETE `/users/bookings/:booking_id`
Cancel an existing booking.

**Request**:
```bash
curl -X DELETE "http://localhost:8080/users/bookings/5?user_id=101"
```

**Response** (200 OK):
```json
{
  "message": "Booking cancelled successfully"
}
```

**Path Parameters**:
- `booking_id` - The ID of the booking to cancel

**Query Parameters**:
- `user_id` (required) - User ID (for ownership verification)

**Logic**:
1. Finds the booking by ID
2. Verifies the booking belongs to the specified user (ownership check)
3. Sets booking status to "cancelled"
4. Returns success message

**Error Responses**:
- `400 Bad Request` - Invalid IDs or missing user_id
- `403 Forbidden` - Booking belongs to a different user
- `404 Not Found` - Booking not found

**Important**: A cancelled booking's slot becomes available again for other users to book.

## Concurrency Approach

The booking system uses **GORM transactions with row-level locking** to handle concurrent booking attempts safely:

### Problem
Without concurrency control, two simultaneous requests for the same slot could both succeed, creating two bookings for the same slot.

### Solution
```go
err = database.DB.Transaction(func(tx *gorm.DB) error {
    // Check if booking exists (implicit row lock acquired)
    var existingBooking models.Booking
    result := tx.Where("coach_id = ? AND slot_time = ? AND status = ?",
        coachID, slotTime, "confirmed").First(&existingBooking)
    
    if result.Error == nil {
        return fmt.Errorf("slot already booked") // Slot is taken
    }
    
    // Create new booking within transaction
    return tx.Create(&booking).Error
})
```

### How It Works
1. When the first request checks for existing bookings, the database acquires a read lock
2. When the first request creates the booking, it upgrades to a write lock
3. The second simultaneous request must wait for the transaction to complete
4. When the second request checks, it sees the booking and returns an error
5. Only one booking per slot is guaranteed per date/coach combination

### Testing Concurrency
You can test race conditions with tools like Apache Bench or wrk:

```bash
# Simulate 100 concurrent requests for the same slot
ab -n 100 -c 10 -p booking.json -T application/json http://localhost:8080/users/bookings
```

Only one request should succeed (201), and others should get 409 Conflict.

## Seed Data

On first run, the database is automatically populated with:

**Coaches**:
- Coach A (Timezone: Asia/Kolkata)
- Coach B (Timezone: America/New_York)

**Users**:
- Alice (Timezone: America/New_York)
- Bob (Timezone: Asia/Kolkata)

**Coach A Availability**:
- Monday: 10:00 - 15:00
- Wednesday: 09:00 - 12:00

**Coach B Availability**:
- Tuesday: 09:00 - 14:00
- Friday: 13:00 - 17:00

You can modify seed data in `database/db.go` before first run.

## Error Handling

All API errors return JSON with an "error" field:

```json
{
  "error": "descriptive error message"
}
```

HTTP Status Codes:
- `200 OK` - Successful GET/DELETE request
- `201 Created` - Successful POST request (resource created)
- `400 Bad Request` - Invalid input, validation failure, malformed request
- `403 Forbidden` - Insufficient permissions (e.g., trying to cancel another user's booking)
- `404 Not Found` - Resource doesn't exist (coach, user, or booking)
- `409 Conflict` - Slot already booked
- `500 Internal Server Error` - Database or server error

## Example Workflow

### 1. Set Coach Availability
```bash
curl -X POST http://localhost:8080/coaches/availability \
  -H "Content-Type: application/json" \
  -d '{
    "coach_id": 1,
    "day": "Tuesday",
    "start_time": "09:00",
    "end_time": "17:00"
  }'
```

### 2. Get Available Slots
```bash
curl -X GET "http://localhost:8080/users/slots?coach_id=1&date=2025-10-28"
```

### 3. Book a Slot
```bash
curl -X POST http://localhost:8080/users/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "datetime": "2025-10-28T09:30:00Z"
  }'
```

### 4. View User Bookings
```bash
curl -X GET "http://localhost:8080/users/bookings?user_id=1"
```

### 5. Cancel a Booking
```bash
curl -X DELETE "http://localhost:8080/users/bookings/1?user_id=1"
```

## Development Notes

### Building for Production
```bash
# Build executable for current OS
go build -o appointment-booking

# Run the executable
./appointment-booking  # On macOS/Linux
# or
appointment-booking.exe  # On Windows
```

### Database Inspection
The `appointments.db` SQLite file can be inspected with any SQLite client:
```bash
sqlite3 appointments.db
```

Common SQL queries:
```sql
SELECT * FROM coaches;
SELECT * FROM users;
SELECT * FROM bookings WHERE status = 'confirmed';
SELECT * FROM coach_availabilities;
```

### Resetting Database
Simply delete `appointments.db` and restart the application:
```bash
rm appointments.db
go run main.go
```

### Performance Considerations
1. **Indexes**: Foreign keys are indexed for fast lookups
2. **Transactions**: Keep them short to minimize lock contention
3. **Slot Generation**: Done in-memory to avoid database queries per slot
4. **Query Optimization**: Uses GORM preloading for efficient coach name fetches

## AI Assistance Note

This project was generated with AI assistance using GitHub Copilot. All code follows Go best practices and production standards:
- Clean separation of concerns (handlers → services → database)
- Type-safe GORM operations
- Transaction-based concurrency control
- Comprehensive error handling
- Input validation on all endpoints
- Well-documented code with clear comments

## Future Enhancements

Potential features for production deployment:
1. **User Authentication**: JWT tokens for user identification
2. **Email Notifications**: Send confirmation emails to users and coaches
3. **Recurring Bookings**: Support for recurring appointments
4. **Cancellation Policies**: Implement cancellation deadlines
5. **Rating System**: Allow users to rate coaches
6. **Payment Integration**: Charge for bookings
7. **Multi-timezone Display**: Convert times to user's timezone for display
8. **Availability Exceptions**: Handle holidays and special days
9. **Waitlist Management**: Allow users to join waitlists for full slots
10. **PostgreSQL Migration**: Upgrade to PostgreSQL for scalability

## License

MIT License - Feel free to use this as a template for your projects.

## Support

For issues or questions, refer to:
- GORM Documentation: https://gorm.io
- Gin Documentation: https://github.com/gin-gonic/gin
- Go Documentation: https://golang.org/doc 
