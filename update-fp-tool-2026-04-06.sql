-- FP proficiency tool data patch (06.04.2026)
-- Applies missing content for Implant Insertion, Implant Removal, IUD Insertion, and IUD Removal.
-- Note: "Score key" descriptive text (purple content) is not stored in a dedicated DB table in this schema.

BEGIN;

-- -------------------------------------------------------------------
-- IMPLANT INSERTION (assessment_type_id = 4)
-- Add distinct thematic areas for one-rod and two-rod insertion.
-- -------------------------------------------------------------------
INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 4, 'Procedure: One rod implant insertion', 3
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 4 AND name = 'Procedure: One rod implant insertion'
);

INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 4, 'Procedure: Two rod implant insertion', 4
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 4 AND name = 'Procedure: Two rod implant insertion'
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 4 AND name = 'Procedure: One rod implant insertion'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Visually verifies the presence of the implant', 10, true, false, 1),
    ('Stretches skin around the insertion site with thumb and index finger', 2, false, false, 2),
    ('Inserts tip of the needle at a slight angle until the tip of the bevel just barely goes under the skin.', 2, false, false, 3),
    ('Releases the skin and lowers the applicator to a horizontal position.', 2, false, false, 4),
    ('Lifts the skin with the tip of the needle.', 2, false, false, 5),
    ('While tenting the skin, gently inserts the needle to its full length (keeping the cannula parallel to the surface of the skin).', 2, false, false, 6),
    ('Unlocks the slider by pushing it fully down.', 2, false, false, 7),
    ('Checks the needle for absence of the implant.', 2, false, false, 8)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 4 AND name = 'Procedure: Two rod implant insertion'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Introduces trocar subdermally at anesthetized site.', 2, false, false, 1),
    ('Inserts the trocar and plunger at a shallow angle with the beveled tip of the trocar facing up.', 2, false, false, 2),
    ('Moves the trocar forward, stopping as soon as the tip is completely beneath the dermis.', 2, false, false, 3),
    ('While tenting the skin, advances trocar and plunger to mark (1) nearest hub of trocar.', 2, false, false, 4),
    ('Removes plunger.', 2, false, false, 5),
    ('Inserts first rod into trocar sleeve.', 2, false, false, 6),
    ('Reinserts plunger and advances it until resistance against the rod is felt.', 2, false, false, 7),
    ('Holds plunger firmly in place with one hand and slides trocar out of incision until it reaches plunger handle.', 2, false, false, 8),
    ('Withdraws trocar and plunger together until mark (2) nearest trocar tip just clears incision (does not remove trocar from skin).', 2, false, false, 9),
    ('Moves tip of trocar away from end of first rod and holds rod out of the path of the trocar.', 2, false, false, 10),
    ('Redirects trocar about 15 degrees and advances trocar and plunger to mark (1).', 2, false, false, 11),
    ('Inserts remaining rod using same technique.', 2, false, false, 12),
    ('Carefully withdraws the trocar and presses down on the incision with a gauzed finger for a minute or so to stop any bleeding.', 2, false, false, 13)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

-- -------------------------------------------------------------------
-- IMPLANT REMOVAL (assessment_type_id = 5)
-- Add distinct thematic areas for forceps vs U technique.
-- -------------------------------------------------------------------
INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 5, 'Removal of implant rod(s): Forceps technique for 1-rod implant', 3
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 5 AND name = 'Removal of implant rod(s): Forceps technique for 1-rod implant'
);

INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 5, 'Removal of implant rod(s): "U" technique for 2-rod implant', 4
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 5 AND name = 'Removal of implant rod(s): "U" technique for 2-rod implant'
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 5 AND name = 'Removal of implant rod(s): Forceps technique for 1-rod implant'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Knows what to do for a deeper insertion or a non palpable implant.', 10, true, false, 1),
    ('Pushes down the proximal tip to fix the implant.', 2, false, false, 2),
    ('Makes a longitudinal incision of 2 mm from below the distal tip of the implant toward the distal tip of the implant.', 2, false, false, 3),
    ('Gently pushes the implant towards the incision with finger-tip until the tip of the implant is visible.', 2, false, false, 4),
    ('Grasps the implant with forceps and removes it by gently pulling it toward the incision.', 2, false, false, 5),
    ('Ensures hemostasis by gentle pressure.', 2, false, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 5 AND name = 'Removal of implant rod(s): "U" technique for 2-rod implant'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Knows what to do for a deeper insertion or a non palpable implant.', 2, false, false, 1),
    ('Chooses a point for incision between the rods, about 5 mm from the ends of the rods nearest the elbow.', 2, false, false, 2),
    ('Makes a small (4 mm) vertical incision to (and between) the long axis of the rods.', 2, false, false, 3),
    ('Gently inserts the holding forceps through the incision at a right angle to the long axis of the nearest rod.', 2, false, false, 4),
    ('Stabilizes the rod that is closest to the incision by placing the index finger parallel to the rod.', 2, false, false, 5),
    ('Advances the forceps until the tip touches the rod.', 2, false, false, 6),
    ('Then opens the forceps and grasps the rod at a right angle to its long axis about 5 mm above the distal end.', 2, false, false, 7),
    ('Cleans off and opens the fibrous tissue sheath surrounding the rod by gently rubbing with sterile gauze to expose the rod for easy removal.', 2, false, false, 8),
    ('Grasps the exposed part of the rod. Releases the holding forceps and slowly and gently removes the rod.', 2, false, false, 9),
    ('Removes the remaining rod using the same technique.', 2, false, false, 10)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

-- -------------------------------------------------------------------
-- IUD/IUS INSERTION (assessment_type_id = 6)
-- Distinguish IUD and IUS insertion tracks and add missing procedure rows.
-- -------------------------------------------------------------------
INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 6, 'Procedure: IUD/IUS Insertion - Setup and loading', 5
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 6 AND name = 'Procedure: IUD/IUS Insertion - Setup and loading'
);

INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 6, 'Procedure: IUD insertion', 6
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 6 AND name = 'Procedure: IUD insertion'
);

INSERT INTO thematic_areas (assessment_type_id, name, display_order)
SELECT 6, 'Procedure: IUS insertion', 7
WHERE NOT EXISTS (
    SELECT 1 FROM thematic_areas
    WHERE assessment_type_id = 6 AND name = 'Procedure: IUS insertion'
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 6 AND name = 'Procedure: IUD/IUS Insertion - Setup and loading'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Cleans cervical os with an antiseptic.', 10, true, false, 1),
    ('Gently grasps cervix with tenaculum or vulsellum forceps and applies gentle traction to straighten the cervical canal.', 2, false, false, 2),
    ('Accurately assesses depth and position of uterine cavity with sound or cannula using no-touch technique.', 2, false, false, 3),
    ('Loads IUD/IUS in sterile package using aseptic technique.', 2, false, false, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 6 AND name = 'Procedure: IUD insertion'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Sets gauge of loaded IUD inserter to the assessed uterine depth.', 2, false, false, 1),
    ('Gently inserts and releases loaded IUD in the uterus using withdrawal technique.', 10, true, false, 2),
    ('Discards plunger and pushes insertion tube gently upwards for fundal placement.', 2, false, false, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 6 AND name = 'Procedure: IUS insertion'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Makes sure the arms of the IUS are horizontal.', 10, true, false, 1),
    ('Releases the threads and pulls them to retract arms of the IUS into the insertion tube.', 2, false, false, 2),
    ('Secures threads and sets flange to the assessed uterine depth.', 2, false, false, 3),
    ('Gently inserts IUS using non-touch technique till 1-2 cm short of assessed uterine depth.', 10, true, false, 4),
    ('Retracts slider to the mark and waits 10 seconds for the IUS arms to fully open.', 2, false, false, 5),
    ('Advances flange all the way to the cervix.', 2, false, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

-- Fix incomplete IUD insertion question text
UPDATE questions
SET question_text = 'Pulls down the introitus with two fingers and visualizes the interior of the vagina.'
WHERE question_text ILIKE 'Pulls down the introitus with two fingers and visualize the interior of the%'
  AND thematic_area_id IN (
      SELECT id FROM thematic_areas WHERE assessment_type_id = 6
  );

-- -------------------------------------------------------------------
-- IUD/IUS REMOVAL (assessment_type_id = 7)
-- Add missing rows 7, 9, and 10 from tool.
-- -------------------------------------------------------------------
WITH ta AS (
    SELECT id FROM thematic_areas
    WHERE assessment_type_id = 7 AND name = 'Procedure: IUD/IUS Removal'
)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM ta,
(VALUES
    ('Performs a bimanual examination and inserts speculum gently to look at length and position of strings.', 10, true, false, 2),
    ('Grasps IUD/IUS strings close to the os with forceps and applies steady gentle traction to remove the IUD/IUS.', 2, false, false, 4),
    ('Knows what to do when strings are missing.', 2, false, false, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE NOT EXISTS (
    SELECT 1 FROM questions x
    WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text
);

COMMIT;
