-- Delete seeded / dummy transactional data (keeps metadata, users, roles, questions)
-- Safe wipe for assessments and sample health workers used in demo seeding.

BEGIN;

-- FP proficiency assessment transactional data
DELETE FROM thematic_area_scores;
DELETE FROM assessment_responses;
DELETE FROM assessments;

-- Health workers created for demo/seeding (optional: comment out if you want to keep them)
DELETE FROM health_workers;

-- RH SPARS transactional data (if tables exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rhspars_action_plan_items') THEN
        DELETE FROM rhspars_action_plan_items;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rhspars_domain_scores') THEN
        DELETE FROM rhspars_domain_scores;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rhspars_section_scores') THEN
        DELETE FROM rhspars_section_scores;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rhspars_responses') THEN
        DELETE FROM rhspars_responses;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rhspars_attendees') THEN
        DELETE FROM rhspars_attendees;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rhspars_assessments') THEN
        DELETE FROM rhspars_assessments;
    END IF;
END $$;

-- Optional: clear non-auth audit noise from seeding (keep if you need logs)
-- DELETE FROM events;

COMMIT;

-- Does NOT delete:
--   regions/districts/subcounties/facilities
--   assessment_types / thematic_areas / questions
--   rhspars_domains / rhspars_thematic_areas / rhspars_questions
--   users / roles / permissions / user_roles / role_permissions / user_admin_areas
