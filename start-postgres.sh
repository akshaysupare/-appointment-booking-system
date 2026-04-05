#!/bin/bash
# Quick Start PostgreSQL with Docker for Appointment Booking System
# This script starts a PostgreSQL container ready to use

echo ""
echo "============================================"
echo "Starting PostgreSQL with Docker..."
echo "============================================"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "ERROR: Docker is not running!"
    echo "Please start Docker and try again."
    exit 1
fi

# Stop and remove existing container if it exists
echo "Cleaning up old containers..."
docker stop postgres-app > /dev/null 2>&1
docker rm postgres-app > /dev/null 2>&1

# Start new PostgreSQL container
echo "Creating PostgreSQL container..."
docker run --name postgres-app \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=appointment-booking \
  -p 5432:5432 \
  -d postgres:latest

if [ $? -eq 0 ]; then
    echo ""
    echo "============================================"
    echo "PostgreSQL is running!"
    echo "============================================"
    echo ""
    echo "Database Details:"
    echo "  Host:     localhost"
    echo "  Port:     5432"
    echo "  Username: postgres"
    echo "  Password: root"
    echo "  Database: appointment-booking"
    echo ""
    echo "Waiting for PostgreSQL to be ready (10 seconds)..."
    sleep 10
    echo ""
    echo "You can now start the API server:"
    echo "  go run ."
    echo ""
else
    echo "ERROR: Failed to start PostgreSQL container"
    exit 1
fi
