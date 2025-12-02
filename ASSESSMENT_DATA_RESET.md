# Assessment Data Reset and Reseed Guide

## Overview

This guide explains how to delete all assessment data and health workers, then reseed with new data.

## Database Structure

All assessment-related tables properly reference `health_worker_id`:

1. **`assessments`** - Directly references `health_workers(id)` via `health_worker_id`
2. **`assessment_responses`** - References `assessments(id)`, which has `health_worker_id` (indirect reference)
3. **`thematic_area_scores`** - References `assessments(id)`, which has `health_worker_id` (indirect reference)

## Scripts

### 1. `reset-and-reseed-assessments.sql`
This SQL script deletes:
- All thematic area scores
- All assessment responses
- All assessments
- All health workers

**Note:** This script does NOT delete:
- Geographic hierarchy (regions, districts, subcounties, facilities)
- Assessment types, thematic areas, questions
- Users, roles, permissions
- User admin areas

### 2. `reset-and-reseed-all.sh` (Linux/Mac)
Bash script that:
1. Runs `reset-and-reseed-assessments.sql` to delete data
2. Runs `seed-health-workers.sql` to create new health workers
3. Runs `seed-assessments.sql` to create new assessments

### 3. `reset-and-reseed-all.bat` (Windows)
Windows batch script that does the same as the shell script.

## Usage

### Option 1: Run SQL script manually

```bash
# Step 1: Delete assessment data
psql -U your_username -d your_database -f reset-and-reseed-assessments.sql

# Step 2: Seed health workers
psql -U your_username -d your_database -f seed-health-workers.sql

# Step 3: Seed assessments
psql -U your_username -d your_database -f seed-assessments.sql
```

### Option 2: Use the automated script (Linux/Mac)

```bash
# Make script executable
chmod +x reset-and-reseed-all.sh

# Run the script (modify DB_USER, DB_NAME, DB_HOST as needed)
./reset-and-reseed-all.sh
```

Or set environment variables:
```bash
export DB_USER=postgres
export DB_NAME=fpscore
export DB_HOST=localhost
./reset-and-reseed-all.sh
```

### Option 3: Use the automated script (Windows)

```cmd
REM Modify DB_USER, DB_NAME, DB_HOST in the .bat file, then run:
reset-and-reseed-all.bat
```

## Important Notes

1. **Backup First**: Always backup your database before running these scripts, especially in production.

2. **Metadata Preserved**: The reset script only deletes assessment data and health workers. All other metadata (geographic hierarchy, assessment types, questions, users, etc.) is preserved.

3. **Dependencies**: Make sure you have:
   - Geographic hierarchy data (regions, districts, subcounties, facilities)
   - Assessment types, thematic areas, and questions
   - These should already exist from your initial setup

4. **Order Matters**: Health workers must be seeded before assessments, as assessments reference health workers.

## Verification

After running the scripts, verify the data:

```sql
-- Check health workers count
SELECT COUNT(*) FROM health_workers;

-- Check assessments count
SELECT COUNT(*) FROM assessments;

-- Check assessment responses count
SELECT COUNT(*) FROM assessment_responses;

-- Verify assessments are linked to health workers
SELECT COUNT(*) 
FROM assessments a
JOIN health_workers hw ON a.health_worker_id = hw.id;
```

