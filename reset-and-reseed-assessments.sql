-- Script to delete assessment data and reseed with new data
-- This script:
--   1. Deletes all assessment-related data
--   2. Deletes health workers
--   3. Reseeds health workers
--   4. Reseeds assessments
--
-- This script does NOT delete:
--   - Geographic hierarchy (regions, districts, subcounties, facilities)
--   - Assessment types, thematic areas, questions
--   - Users, roles, permissions
--   - User admin areas
--
-- Usage:
--   psql -U your_username -d your_database -f reset-and-reseed-assessments.sql

BEGIN;

-- Step 1: Delete assessment-related data in reverse order of dependencies
RAISE NOTICE 'Deleting assessment data...';
DELETE FROM thematic_area_scores;
DELETE FROM assessment_responses;
DELETE FROM assessments;

-- Step 2: Delete health workers
RAISE NOTICE 'Deleting health workers...';
DELETE FROM health_workers;

COMMIT;

-- Note: After running this script, you need to run:
--   1. seed-health-workers.sql
--   2. seed-assessments.sql
--
-- Example:
--   psql -U your_username -d your_database -f seed-health-workers.sql
--   psql -U your_username -d your_database -f seed-assessments.sql

