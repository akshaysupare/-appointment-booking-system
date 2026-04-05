# Setup Guide - Appointment Booking System with PostgreSQL

## Quick Start (Choose One)

### **Option 1: PostgreSQL Local Installation (Windows)**

**Step 1: Install PostgreSQL**
- Download: https://www.postgresql.org/download/windows/
- Run installer, choose default settings
- Remember the password you set for `postgres` user

**Step 2: Create Database**

Open pgAdmin (comes with PostgreSQL) or use command line:

```bash
# Using psql command line (Windows PowerShell)
psql -U postgres

# Then run these commands:
CREATE DATABASE "appointment-booking";
\q
```

**Step 3: Update Connection (Optional)**

If you used a different password or host, set environment variables:

```powershell
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "your-password"
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_NAME = "appointment-booking"
```

---

### **Option 2: PostgreSQL via Docker (Easy, No Installation)**

**Step 1: Install Docker Desktop** (if not installed)
- Download: https://www.docker.com/products/docker-desktop
- Run installer and restart computer

**Step 2: Start PostgreSQL Container**

```powershell
docker run --name postgres-app ^
  -e POSTGRES_PASSWORD=root ^
  -e POSTGRES_DB=appointment-booking ^
  -p 5432:5432 ^
  -d postgres:latest
```

This starts PostgreSQL in the background with:
- Password: `root`
- Database: `appointment-booking`
- Accessible at: `localhost:5432`

---

## Step 3: Run the Server

After PostgreSQL is ready:

```powershell
cd c:\Users\supar\OneDrive\Desktop\MY\my_projects\NudgeBee\-appointment-booking-system
go mod tidy
go run .
```

---

## Expected Output

```
Database connected successfully!
Database migrations completed!
Seed data inserted successfully!
[GIN-debug] Listening and serving HTTP on :8080
```

---

## Test the API

```powershell
curl "http://localhost:8080/users/slots?coach_id=1&date=2025-10-28"
```

---

## Connection Configuration

Default settings (no changes needed):
```
Host:     localhost
Port:     5432
Username: postgres
Password: root
Database: appointment-booking
```

Override with environment variables:
```powershell
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "root"
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_NAME = "appointment-booking"
go run .
```

---

## Advantages Over SQLite

✅ **No CGO needed** - Pure Go drivers, works on any system  
✅ **Better for production** - Handles concurrent connections  
✅ **Scalable** - Supports large datasets  
✅ **Multi-user** - Multiple apps can connect simultaneously  
✅ **Easy backup** - Standard database backup tools  

---

## Troubleshooting

**"Connection refused"**
- PostgreSQL not running
- Wrong password (should be `root` by default)
- Wrong host/port

**"Database does not exist"**
- Create database: `CREATE DATABASE "appointment-booking";`

**Docker errors**
- Ensure Docker Desktop is running
- Check port 5432 not already in use: `netstat -ano | findstr :5432`
- Kill existing container: `docker ps` then `docker stop <container-id>`

**Port already in use**
- Change port in docker: `-p 5433:5432` (connects to 5433 on host)
- Then set: `$env:DB_PORT = "5433"`

