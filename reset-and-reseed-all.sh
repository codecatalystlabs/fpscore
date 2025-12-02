#!/bin/bash
# Script to delete assessment data and reseed everything
# This script deletes assessment data and health workers, then reseeds them

# Database connection parameters (modify as needed)
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-fpscore}"
DB_HOST="${DB_HOST:-localhost}"

echo "=========================================="
echo "Resetting and Reseeding Assessment Data"
echo "=========================================="
echo ""

echo "Step 1: Deleting assessment data and health workers..."
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f reset-and-reseed-assessments.sql

if [ $? -ne 0 ]; then
    echo "Error: Failed to delete assessment data"
    exit 1
fi

echo ""
echo "Step 2: Seeding health workers..."
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f seed-health-workers.sql

if [ $? -ne 0 ]; then
    echo "Error: Failed to seed health workers"
    exit 1
fi

echo ""
echo "Step 3: Seeding assessments..."
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f seed-assessments.sql

if [ $? -ne 0 ]; then
    echo "Error: Failed to seed assessments"
    exit 1
fi

echo ""
echo "=========================================="
echo "Reset and reseed completed successfully!"
echo "=========================================="

