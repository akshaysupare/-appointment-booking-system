# Postman Collection Import Guide

## Quick Import (30 seconds)

### Step 1: Download the Files
You have 2 files:
- `Appointment_Booking_API.postman_collection.json` - API endpoints
- `Appointment_Booking_API.postman_environment.json` - Environment variables

### Step 2: Import in Postman

**Import Collection:**
1. Open Postman
2. Click **Import** (top left)
3. Select `Appointment_Booking_API.postman_collection.json`
4. Click **Import**

**Import Environment:**
1. Click **Import** again
2. Select `Appointment_Booking_API.postman_environment.json`
3. Click **Import**

### Step 3: Select Environment
1. Top right, click environment dropdown
2. Select **Appointment Booking - Local**
3. Done! ✅

---

## Collection Overview

### 📁 Folders

1. **Base URL Setup**
   - Info on environment variables

2. **Coach Management**
   - `POST /coaches/availability` - Set coach's weekly schedule

3. **User Slots & Bookings**
   - `GET /users/slots` - View available 30-min slots
   - `POST /users/bookings` - Book a slot
   - `GET /users/bookings` - View my bookings
   - `DELETE /users/bookings/:id` - Cancel booking

4. **Test Scenarios**
   - Complete booking flow walkthrough
   - Concurrency testing guide
   - Error handling examples

---

## Environment Variables

Pre-configured with seed data defaults:

| Variable | Default | Purpose |
|----------|---------|---------|
| `base_url` | `http://localhost:8080` | API server address |
| `coach_id` | `1` | Coach A (from seed data) |
| `user_id` | `1` | Alice (from seed data) |
| `booking_id` | `5` | Example booking (update after booking) |
| `booking_date` | `2025-10-28` | Test date (must be valid for coach) |

**To modify:**
1. Top right environment dropdown
2. Click eye icon to see values
3. Click **Edit** to change
4. Save

---

## Quick Start: Complete Booking Flow

### 1️⃣ Set Coach Availability
```
POST /coaches/availability
```
Body:
```json
{
  "coach_id": 1,
  "day": "Monday",
  "start_time": "09:00",
  "end_time": "17:00"
}
```

### 2️⃣ View Available Slots
```
GET /users/slots?coach_id=1&date=2025-10-28
```

### 3️⃣ Book a Slot
```
POST /users/bookings
```
Body:
```json
{
  "user_id": 1,
  "coach_id": 1,
  "datetime": "2025-10-28T09:30:00Z"
}
```

### 4️⃣ View My Bookings
```
GET /users/bookings?user_id=1
```

### 5️⃣ Cancel Booking (Optional)
```
DELETE /users/bookings/5?user_id=1
```

---

## Seed Data (Pre-loaded)

**Coaches:**
- ID 1: Coach A (Asia/Kolkata)
  - Monday: 10:00-15:00
  - Wednesday: 09:00-12:00
- ID 2: Coach B (America/New_York)
  - Tuesday: 09:00-14:00
  - Friday: 13:00-17:00

**Users:**
- ID 1: Alice (America/New_York)
- ID 2: Bob (Asia/Kolkata)

---

## Testing Features

### ✅ Response Examples
Each endpoint has multiple response examples:
- 200/201 Success responses
- 400 Bad Request
- 403 Forbidden
- 404 Not Found
- 409 Conflict

Click on response names to see examples.

### ✅ Test Scenarios
Built-in documentation for:
1. **Complete Booking Flow** - End-to-end walkthrough
2. **Concurrency Testing** - Race condition prevention
3. **Error Handling** - Common error cases

### ✅ Dynamic Variables
Use `{{variable_name}}` in requests:
- `{{base_url}}` - Server address
- `{{coach_id}}` - Coach ID
- `{{user_id}}` - User ID
- `{{booking_id}}` - Booking ID
- `{{booking_date}}` - Test date

---

## Advanced Usage

### Run Collection with Newman (CLI)
```bash
# Install Newman (if not already)
npm install -g newman

# Run collection
newman run "Appointment_Booking_API.postman_collection.json" \
  -e "Appointment_Booking_API.postman_environment.json"
```

### Export as Tests
1. Click **...** on collection
2. **Export**
3. Choose format for your testing framework

### Team Collaboration
1. **Export collection** (right-click → Export)
2. **Share .json files** with team
3. **Team imports** both collection and environment
4. **Update** `base_url` if using different server

---

## Troubleshooting

**"Could not send request"**
- Ensure API server is running: `go run .`
- Check `base_url` points to correct address (default: http://localhost:8080)

**"404 Not Found"**
- Seed data might not be loaded
- Check database is initialized
- Verify coach_id/user_id exist (defaults: 1 and 1)

**"409 Conflict when booking"**
- Slot already booked
- Try a different datetime from the slots list
- Or cancel existing booking first

**Variables not working**
- Ensure environment is selected (top right dropdown)
- Check variable names are spelled correctly
- Use correct syntax: `{{variable_name}}`

---

## API Documentation

Full detailed documentation in [README.md](./README.md):
- Complete endpoint reference
- Request/response formats
- Concurrency approach explanation
- Error codes and meanings
- Setup instructions

---

## File Locations

In your project folder:
```
appointment-booking/
├── Appointment_Booking_API.postman_collection.json    ← Import this
├── Appointment_Booking_API.postman_environment.json    ← And this
├── README.md                                           ← Full docs
└── SETUP.md                                            ← Setup guide
```

---

## Support

- **API Issues**: Check README.md for full documentation
- **Database**: Ensure PostgreSQL is running with correct credentials
- **Server**: Verify `go run .` output shows "Listening on :8080"
- **Postman Help**: https://learning.postman.com/

Enjoy testing! 🚀
