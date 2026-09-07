-- Remove duplicate proficiency assessments created by double-submit.
-- Keeps the earliest row (lowest id) in each duplicate group.
-- Child rows (assessment_responses, thematic_area_scores) cascade-delete.
--
-- Preview first (recommended):
--   See the SELECT below, then run the DELETE in a transaction.

BEGIN;

-- Preview which rows would be deleted (later copies of the same visit)
SELECT
    dup.id AS duplicate_id,
    keep.id AS keep_id,
    dup.health_worker_id,
    dup.facility_id,
    dup.assessment_type_id,
    dup.percentage_score,
    dup.created_at AS duplicate_created_at,
    keep.created_at AS keep_created_at
FROM assessments dup
JOIN assessments keep
  ON keep.id < dup.id
 AND keep.health_worker_id = dup.health_worker_id
 AND keep.facility_id = dup.facility_id
 AND keep.assessment_type_id = dup.assessment_type_id
 AND keep.achieved_score = dup.achieved_score
 AND keep.total_possible_score = dup.total_possible_score
 AND ROUND(keep.percentage_score::numeric, 2) = ROUND(dup.percentage_score::numeric, 2)
 AND ABS(EXTRACT(EPOCH FROM (dup.created_at - keep.created_at))) <= 120
ORDER BY dup.id;

-- Delete later duplicates (same worker/facility/type/score within 2 minutes)
DELETE FROM assessments a
USING assessments earlier
WHERE earlier.id < a.id
  AND earlier.health_worker_id = a.health_worker_id
  AND earlier.facility_id = a.facility_id
  AND earlier.assessment_type_id = a.assessment_type_id
  AND earlier.achieved_score = a.achieved_score
  AND earlier.total_possible_score = a.total_possible_score
  AND ROUND(earlier.percentage_score::numeric, 2) = ROUND(a.percentage_score::numeric, 2)
  AND ABS(EXTRACT(EPOCH FROM (a.created_at - earlier.created_at))) <= 120;

COMMIT;

-- If the preview looks wrong, run ROLLBACK; instead of COMMIT.
