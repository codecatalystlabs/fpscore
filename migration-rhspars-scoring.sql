-- Align RH SPARS assessments with FP proficiency scoring metadata
ALTER TABLE rhspars_assessments
    ADD COLUMN IF NOT EXISTS total_possible_score INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS achieved_score INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS performance_level VARCHAR(50) NOT NULL DEFAULT 'Not Acceptable';

UPDATE rhspars_assessments a
SET
    total_possible_score = COALESCE((
        SELECT SUM(possible_score) FROM rhspars_section_scores ss WHERE ss.assessment_id = a.id
    ), 0),
    achieved_score = COALESCE((
        SELECT SUM(achieved_score) FROM rhspars_section_scores ss WHERE ss.assessment_id = a.id
    ), 0),
    performance_level = CASE
        WHEN grand_percentage > 90 THEN 'Proficient'
        WHEN grand_percentage >= 70 THEN 'Competent'
        ELSE 'Not Acceptable'
    END;
