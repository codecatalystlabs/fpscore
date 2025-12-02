@echo off
REM Script to delete assessment data and reseed everything (Windows)
REM This script deletes assessment data and health workers, then reseeds them

REM Database connection parameters (modify as needed)
set DB_USER=postgres
set DB_NAME=fpscore
set DB_HOST=localhost

echo ==========================================
echo Resetting and Reseeding Assessment Data
echo ==========================================
echo.

echo Step 1: Deleting assessment data and health workers...
psql -h %DB_HOST% -U %DB_USER% -d %DB_NAME% -f reset-and-reseed-assessments.sql
if errorlevel 1 (
    echo Error: Failed to delete assessment data
    exit /b 1
)

echo.
echo Step 2: Seeding health workers...
psql -h %DB_HOST% -U %DB_USER% -d %DB_NAME% -f seed-health-workers.sql
if errorlevel 1 (
    echo Error: Failed to seed health workers
    exit /b 1
)

echo.
echo Step 3: Seeding assessments...
psql -h %DB_HOST% -U %DB_USER% -d %DB_NAME% -f seed-assessments.sql
if errorlevel 1 (
    echo Error: Failed to seed assessments
    exit /b 1
)

echo.
echo ==========================================
echo Reset and reseed completed successfully!
echo ==========================================

