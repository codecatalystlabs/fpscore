-- Script to delete assessment data only (not metadata)
-- This script deletes:
--   - Assessment responses and scores
--   - Assessments
--   - Health workers (since they're linked to assessments)
--
-- This script does NOT delete:
--   - Geographic hierarchy (regions, districts, subcounties, facilities)
--   - Assessment types, thematic areas, questions
--   - Users, roles, permissions
--   - User admin areas

-- Delete in reverse order of dependencies to avoid foreign key constraint violations

-- Delete assessment-related data
DELETE FROM thematic_area_scores;
DELETE FROM assessment_responses;
DELETE FROM assessments;

-- Delete health workers (since they're only used for assessments)
DELETE FROM health_workers;

