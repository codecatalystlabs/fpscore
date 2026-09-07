-- RH SPARS (Integrated Reproductive Health Support Supervision Tool) seed
-- Generated from Final integrated tool.xlsx structure
-- Scoring: Yes=1, No=0, NA excluded from numerator and denominator
-- section% = sum/(n-NA)*100; domain% = avg(section%); grand% = avg(domain%)

BEGIN;

INSERT INTO rhspars_domains (code, name, display_order) SELECT 'maternity', 'Maternity services', 1 WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='maternity');

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Health Information System - Maternity Register', 1 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Health Information System - Maternity Register');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Integrated Maternity Register is available', 1, 1),
    ('Ten random entries in the Integrated Maternity Register are completely filled', 1, 2),
    ('Monthly summary totals are present in the Integrated Maternity Register', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Health Information System - Maternity Register'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Infrastructure - Maternity confidentiality and space', 2 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Infrastructure - Maternity confidentiality and space');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Maternity: space for this service exists', 1, 1),
    ('Maternity: space/room ensures privacy for every client', 1, 2),
    ('Maternity: safe place to keep case notes so non-health professionals cannot read them', 1, 3),
    ('Facility has separate arrangements or waiting areas for young people', 1, 4),
    ('Clean toilet facilities clearly labelled male or female', 1, 5),
    ('Facility has a functional placenta pit', 1, 6)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Infrastructure - Maternity confidentiality and space'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Job aids and protocols - Maternity unit', 3 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Job aids and protocols - Maternity unit');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Essential maternal and newborn care guidelines available at maternity workstation', 1, 1),
    ('PPH protocol/job aid available', 1, 2),
    ('Hypertension in pregnancy protocol/job aid available', 1, 3),
    ('PPFP compendium available', 1, 4),
    ('PAC FP compendium/protocols available', 1, 5),
    ('HIV EMTCT guidelines available', 1, 6),
    ('Referral forms available', 1, 7),
    ('Consent forms available', 1, 8),
    ('Partograph available', 1, 9)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Job aids and protocols - Maternity unit'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Infection prevention practices - Maternity unit', 4 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Infection prevention practices - Maternity unit');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Colour-coded bins all present: non-infectious (black), infectious (yellow), highly infectious (red), pharmaceutical (brown), domestic (green)', 1, 1),
    ('Sharps boxes available', 1, 2),
    ('Handwashing facilities available', 1, 3),
    ('Sterilizer/autoclave available and functional for reusable equipment', 1, 4)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Infection prevention practices - Maternity unit'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Essential equipment availability - Maternity', 5 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Essential equipment availability - Maternity');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Neonatal ambu bag and mask available and functional', 1, 1),
    ('Weighing scale for newborns available and functional', 1, 2),
    ('Penguin sucker available and functional', 1, 3),
    ('Neonatal resuscitation table available and functional', 1, 4),
    ('Sterilization equipment available and functional', 1, 5),
    ('Vacuum aspirator or D and C kit available and functional', 1, 6),
    ('Suction apparatus available and functional', 1, 7),
    ('Manual vacuum extractor available and functional', 1, 8),
    ('Delivery bed available and functional', 1, 9),
    ('Delivery bed for disabled available and functional', 1, 10),
    ('Delivery kit available and functional', 1, 11),
    ('Examination light available and functional', 1, 12),
    ('Emergency transport available and functional', 1, 13),
    ('PPH kit available and functional', 1, 14),
    ('PET kit available and functional', 1, 15),
    ('Fridge or cold box/vaccine carrier with cold ice packs for storage of oxytocin available and functional', 1, 16)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Essential equipment availability - Maternity'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Emergency medicines - Maternity', 6 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Emergency medicines - Maternity');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Magnesium sulphate inj. (ampoule) available', 1, 1),
    ('Calcium gluconate inj. (ampoule) available', 1, 2),
    ('Dextrose 50% (bottle) available', 1, 3),
    ('Hydrocortisone inj. (vial) available', 1, 4),
    ('Misoprostol (tablet) available', 1, 5),
    ('Vitamin K available', 1, 6),
    ('Oxytocin available', 1, 7),
    ('Tranexamic acid available', 1, 8),
    ('Heat-stable carbetocin available', 1, 9)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Emergency medicines - Maternity'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Human resources - Staff training (Maternity)', 7 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Human resources - Staff training (Maternity)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Any staff received training in maternity (on-job or workshop) within the past year (evidence: training register, CME book, and/or MoH iHRIS)', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Human resources - Staff training (Maternity)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Functionality of committees - Department meeting (Maternity)', 8 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Functionality of committees - Department meeting (Maternity)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Active department committee with regular meetings (last 2 meetings)', 1, 1),
    ('Data issues discussed in these meetings', 1, 2),
    ('RH supplies and commodities discussed in these meetings', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Functionality of committees - Department meeting (Maternity)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Support supervision - Maternity', 9 FROM rhspars_domains d WHERE d.code='maternity' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Support supervision - Maternity');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Maternity supervised by facility in-charge or senior management at least once during the previous quarter', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'maternity' AND ta.name = 'Support supervision - Maternity'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_domains (code, name, display_order) SELECT 'fp', 'Family planning services', 2 WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='fp');

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Health Information System - FP Register', 1 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Health Information System - FP Register');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Integrated FP Register is available', 1, 1),
    ('Ten random entries in the Integrated FP Register are completely filled', 1, 2),
    ('Monthly summary totals are present in the Integrated FP Register', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Health Information System - FP Register'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Infrastructure - Family planning confidentiality and space', 2 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Infrastructure - Family planning confidentiality and space');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Family planning: space for this service exists', 1, 1),
    ('Family planning: consultation room ensures privacy for every client', 1, 2),
    ('Family planning: safe place to keep case notes so non-health professionals cannot read them', 1, 3),
    ('Facility has separate arrangements or waiting areas for young people', 1, 4),
    ('Clean toilet facilities clearly labelled male or female', 1, 5),
    ('Facilities for waste disposal (dump pit, incinerator) or disposal plan available', 1, 6)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Infrastructure - Family planning confidentiality and space'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Job aids and protocols - Family planning space(s)', 3 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Job aids and protocols - Family planning space(s)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Method mix chart available at FP workstation', 1, 1),
    ('Myths and misconceptions pocket book available at FP workstation', 1, 2),
    ('FP counselling flip chart available at FP workstation', 1, 3),
    ('PPFP compendium available at FP workstation', 1, 4),
    ('PAC FP compendium/protocols available at FP workstation', 1, 5),
    ('MEC wheel available at FP workstation', 1, 6),
    ('Demonstration models available at FP workstation', 1, 7),
    ('Referral forms available at FP workstation', 1, 8),
    ('Consent forms available at FP workstation', 1, 9)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Job aids and protocols - Family planning space(s)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Infection prevention practices - Family planning space(s)', 4 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Infection prevention practices - Family planning space(s)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Colour-coded bins all present: non-infectious (black), infectious (yellow), highly infectious (red), pharmaceutical (brown), domestic (green)', 1, 1),
    ('Sharps boxes available', 1, 2),
    ('Handwashing facilities available', 1, 3),
    ('Sterilizer/autoclave available and functional for reusable equipment', 1, 4)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Infection prevention practices - Family planning space(s)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Essential equipment availability - Family planning', 5 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Essential equipment availability - Family planning');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Weighing scale available and functional', 1, 1),
    ('BP machine available and functional', 1, 2),
    ('Implant insertion/removal kit complete (mosquito artery forceps curved and straight; small dissecting forceps with teeth; scissors small sharp; scalpel handle with blade No. 11 or 15)', 1, 3),
    ('IUD insertion/removal kit complete (Cusco speculum; tenaculum; uterine sound; sponge-holding forceps; long curved or ring forceps; scissors)', 1, 4),
    ('BTL kit available and functional', 1, 5),
    ('Vasectomy kit available and functional', 1, 6),
    ('Angle light / lamp available and functional', 1, 7),
    ('Cervical cancer screening and preventive treatment available', 1, 8),
    ('Breast examination service available', 1, 9)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Essential equipment availability - Family planning'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Human resources - Staff training (Family planning)', 6 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Human resources - Staff training (Family planning)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Any staff received training in family planning within the past year (evidence: training register, CME book, and/or MoH iHRIS)', 1, 1),
    ('Staff undergone proficiency testing for family planning', 1, 2)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Human resources - Staff training (Family planning)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Functionality of committees - Department meeting (FP)', 7 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Functionality of committees - Department meeting (FP)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Active department committee with regular meetings (last 2 meetings)', 1, 1),
    ('Data issues discussed in these meetings', 1, 2),
    ('RH supplies and commodities discussed in these meetings', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Functionality of committees - Department meeting (FP)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Support supervision - Family planning', 8 FROM rhspars_domains d WHERE d.code='fp' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Support supervision - Family planning');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Family planning supervised by facility in-charge or senior management at least once during the previous quarter', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'fp' AND ta.name = 'Support supervision - Family planning'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_domains (code, name, display_order) SELECT 'anc', 'Antenatal care services', 3 WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='anc');

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Health Information System - ANC Register', 1 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Health Information System - ANC Register');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Integrated ANC Register is available', 1, 1),
    ('Ten random entries in the Integrated ANC Register are completely filled', 1, 2),
    ('Monthly summary totals are present in the Integrated ANC Register', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Health Information System - ANC Register'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Infrastructure - Antenatal clinic confidentiality and space', 2 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Infrastructure - Antenatal clinic confidentiality and space');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Antenatal clinic: space for this service exists', 1, 1),
    ('Antenatal clinic: consultation room ensures privacy for every client', 1, 2),
    ('Antenatal clinic: safe place to keep case notes so non-health professionals cannot read them', 1, 3),
    ('Facility has separate arrangements or waiting areas for young people', 1, 4),
    ('Clean toilet facilities clearly labelled male or female', 1, 5),
    ('Facilities for waste disposal (dump pit, incinerator) or disposal plan available', 1, 6)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Infrastructure - Antenatal clinic confidentiality and space'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Job aids and protocols - Antenatal unit', 3 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Job aids and protocols - Antenatal unit');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Essential Maternal and Newborn Care guidelines available at ANC workstation', 1, 1),
    ('STI screening protocols available at ANC workstation', 1, 2),
    ('HIV EMTCT guidelines available at ANC workstation', 1, 3),
    ('PPH client-facing posters/messages (high-risk conditions) available at ANC workstation', 1, 4),
    ('PPFP compendium available at ANC workstation', 1, 5),
    ('Referral forms available at ANC workstation', 1, 6),
    ('Consent forms available at ANC workstation', 1, 7)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Job aids and protocols - Antenatal unit'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Infection prevention practices - Antenatal clinic', 4 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Infection prevention practices - Antenatal clinic');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Colour-coded bins all present: non-infectious (black), infectious (yellow), highly infectious (red), pharmaceutical (brown), domestic (green)', 1, 1),
    ('Sharps boxes available', 1, 2),
    ('Handwashing facilities available', 1, 3),
    ('Sterilizer/autoclave available and functional for reusable equipment', 1, 4)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Infection prevention practices - Antenatal clinic'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Essential equipment availability - ANC', 5 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Essential equipment availability - ANC');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('BP machine available and in use', 1, 1),
    ('Fetoscope available and in use', 1, 2),
    ('Ultrasound available and in use', 1, 3),
    ('MUAC tape available and in use', 1, 4),
    ('Weighing scale available and in use', 1, 5),
    ('Height meter available and in use', 1, 6)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Essential equipment availability - ANC'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Human resources - Staff training (ANC)', 6 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Human resources - Staff training (ANC)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Any staff received training in ANC within the past year (evidence: training register, CME book, and/or MoH iHRIS)', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Human resources - Staff training (ANC)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Functionality of committees - Department meeting (ANC)', 7 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Functionality of committees - Department meeting (ANC)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Active department committee with regular meetings (last 2 meetings)', 1, 1),
    ('Data issues discussed in these meetings', 1, 2),
    ('RH supplies and commodities discussed in these meetings', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Functionality of committees - Department meeting (ANC)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Support supervision - ANC services', 8 FROM rhspars_domains d WHERE d.code='anc' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Support supervision - ANC services');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('ANC services supervised by facility in-charge or senior management at least once during the previous quarter', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'anc' AND ta.name = 'Support supervision - ANC services'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_domains (code, name, display_order) SELECT 'commodities', 'Drugs and supplies', 4 WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='commodities');

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Availability of items (B1)', 1 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Availability of items (B1)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: item available on day of survey', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: item available on day of survey', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): item available on day of survey', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): item available on day of survey', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): item available on day of survey', 1, 5),
    ('Levonorgestrel 0.15mg + Ethinyl Estradiol 0.03mg Tablets (28 Tablets): item available on day of survey', 1, 6),
    ('Levonorgestrel 0.75mg, 2 Tablets (ECP): item available on day of survey', 1, 7),
    ('Hormonal IUD: item available on day of survey', 1, 8),
    ('Cu-T Intrauterine Device: item available on day of survey', 1, 9),
    ('Male Condom: item available on day of survey', 1, 10),
    ('Sulfadoxine 500 + Pyrimethamine 25mg Tablets: item available on day of survey', 1, 11),
    ('Ferrous Sulphate/Fumarate + Folic Acid Tablets: item available on day of survey', 1, 12),
    ('Mebendazole 500mg Tablets: item available on day of survey', 1, 13),
    ('Oxytocin Injection 10 i.u./mL: item available on day of survey', 1, 14),
    ('Magnesium Sulphate Injection 50%: item available on day of survey', 1, 15),
    ('Female condoms: item available on day of survey', 1, 16),
    ('Moon beads: item available on day of survey', 1, 17)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Availability of items (B1)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - No expired stock (B)', 2 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - No expired stock (B)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: no expired quantity in stock (score 1 if NO expired stock present)', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: no expired quantity in stock (score 1 if NO expired stock present)', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): no expired quantity in stock (score 1 if NO expired stock present)', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): no expired quantity in stock (score 1 if NO expired stock present)', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): no expired quantity in stock (score 1 if NO expired stock present)', 1, 5),
    ('Levonorgestrel 0.15mg + Ethinyl Estradiol 0.03mg Tablets (28 Tablets): no expired quantity in stock (score 1 if NO expired stock present)', 1, 6),
    ('Levonorgestrel 0.75mg, 2 Tablets (ECP): no expired quantity in stock (score 1 if NO expired stock present)', 1, 7),
    ('Hormonal IUD: no expired quantity in stock (score 1 if NO expired stock present)', 1, 8),
    ('Cu-T Intrauterine Device: no expired quantity in stock (score 1 if NO expired stock present)', 1, 9),
    ('Male Condom: no expired quantity in stock (score 1 if NO expired stock present)', 1, 10),
    ('Sulfadoxine 500 + Pyrimethamine 25mg Tablets: no expired quantity in stock (score 1 if NO expired stock present)', 1, 11),
    ('Ferrous Sulphate/Fumarate + Folic Acid Tablets: no expired quantity in stock (score 1 if NO expired stock present)', 1, 12),
    ('Mebendazole 500mg Tablets: no expired quantity in stock (score 1 if NO expired stock present)', 1, 13),
    ('Oxytocin Injection 10 i.u./mL: no expired quantity in stock (score 1 if NO expired stock present)', 1, 14),
    ('Magnesium Sulphate Injection 50%: no expired quantity in stock (score 1 if NO expired stock present)', 1, 15),
    ('Female condoms: no expired quantity in stock (score 1 if NO expired stock present)', 1, 16),
    ('Moon beads: no expired quantity in stock (score 1 if NO expired stock present)', 1, 17)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - No expired stock (B)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Availability of stock cards (B2)', 3 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Availability of stock cards (B2)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: stock card available', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: stock card available', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): stock card available', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): stock card available', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): stock card available', 1, 5),
    ('Levonorgestrel 0.15mg + Ethinyl Estradiol 0.03mg Tablets (28 Tablets): stock card available', 1, 6),
    ('Levonorgestrel 0.75mg, 2 Tablets (ECP): stock card available', 1, 7),
    ('Hormonal IUD: stock card available', 1, 8),
    ('Cu-T Intrauterine Device: stock card available', 1, 9),
    ('Male Condom: stock card available', 1, 10),
    ('Sulfadoxine 500 + Pyrimethamine 25mg Tablets: stock card available', 1, 11),
    ('Ferrous Sulphate/Fumarate + Folic Acid Tablets: stock card available', 1, 12),
    ('Mebendazole 500mg Tablets: stock card available', 1, 13),
    ('Oxytocin Injection 10 i.u./mL: stock card available', 1, 14),
    ('Magnesium Sulphate Injection 50%: stock card available', 1, 15),
    ('Female condoms: stock card available', 1, 16),
    ('Moon beads: stock card available', 1, 17)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Availability of stock cards (B2)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Physical count last 3 months (B3)', 4 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Physical count last 3 months (B3)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: physical count done for last three months and PC marked in stock card', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: physical count done for last three months and PC marked in stock card', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): physical count done for last three months and PC marked in stock card', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): physical count done for last three months and PC marked in stock card', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): physical count done for last three months and PC marked in stock card', 1, 5)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Physical count last 3 months (B3)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Correct filling of stock cards (B4)', 5 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Correct filling of stock cards (B4)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: stock card filled correctly with name, strength, dosage form, AMC', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: stock card filled correctly with name, strength, dosage form, AMC', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): stock card filled correctly with name, strength, dosage form, AMC', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): stock card filled correctly with name, strength, dosage form, AMC', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): stock card filled correctly with name, strength, dosage form, AMC', 1, 5)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Correct filling of stock cards (B4)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Updating/accuracy of stock card balance (B5)', 6 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Updating/accuracy of stock card balance (B5)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: stock card balance and physical count agree 100%', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: stock card balance and physical count agree 100%', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): stock card balance and physical count agree 100%', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): stock card balance and physical count agree 100%', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): stock card balance and physical count agree 100%', 1, 5)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Updating/accuracy of stock card balance (B5)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Accuracy of AMC (B6)', 7 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Accuracy of AMC (B6)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: calculated AMC same as recorded AMC within +/- 10%', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: calculated AMC same as recorded AMC within +/- 10%', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): calculated AMC same as recorded AMC within +/- 10%', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): calculated AMC same as recorded AMC within +/- 10%', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): calculated AMC same as recorded AMC within +/- 10%', 1, 5)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Accuracy of AMC (B6)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Stock management - Stock status optimum (P)', 8 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Stock management - Stock status optimum (P)');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 1),
    ('Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 2),
    ('Etonogestrel 68mg Implant (Implanon): stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 3),
    ('Levonorgestrel 2x75mg Implant (Levoplant): stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 4),
    ('Levonorgestrel 2x75mg Implant (Jadelle): stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 5),
    ('Levonorgestrel 0.15mg + Ethinyl Estradiol 0.03mg Tablets (28 Tablets): stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 6),
    ('Levonorgestrel 0.75mg, 2 Tablets (ECP): stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 7),
    ('Hormonal IUD: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 8),
    ('Cu-T Intrauterine Device: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 9),
    ('Male Condom: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 10),
    ('Sulfadoxine 500 + Pyrimethamine 25mg Tablets: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 11),
    ('Ferrous Sulphate/Fumarate + Folic Acid Tablets: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 12),
    ('Mebendazole 500mg Tablets: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 13),
    ('Oxytocin Injection 10 i.u./mL: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 14),
    ('Magnesium Sulphate Injection 50%: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 15),
    ('Female condoms: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 16),
    ('Moon beads: stock status is Optimum for the next 3 months (Yes=Optimum; No=Overstock or Understock)', 1, 17)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Stock management - Stock status optimum (P)'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Procurement planning and ordering - Awareness', 9 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Procurement planning and ordering - Awareness');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Facility undertook joint planning (e.g., with IP) and staff are aware of planned quarterly/monthly Family Planning Outreaches', 1, 1),
    ('Facility staff aware of emergency ordering or redistribution procedures for RH commodities (notify DHT of shortages or excesses)', 1, 2),
    ('Facility is using the online stock status facility system', 1, 3)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Procurement planning and ordering - Awareness'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Procurement plan contents', 10 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Procurement plan contents');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('EMHS Procurement Plan for the current financial year available at the health facility', 1, 1),
    ('Procurement plan includes: Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection', 1, 2),
    ('Procurement plan includes: Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection', 1, 3),
    ('Procurement plan includes: Etonogestrel 68mg Implant (Implanon)', 1, 4),
    ('Procurement plan includes: Levonorgestrel 2x75mg Implant (Levoplant)', 1, 5),
    ('Procurement plan includes: Levonorgestrel 2x75mg Implant (Jadelle)', 1, 6),
    ('Procurement plan includes: Levonorgestrel 0.15mg + Ethinyl Estradiol 0.03mg Tablets (28 Tablets)', 1, 7),
    ('Procurement plan includes: Levonorgestrel 0.75mg, 2 Tablets (ECP)', 1, 8),
    ('Procurement plan includes: Hormonal IUD', 1, 9),
    ('Procurement plan includes: Cu-T Intrauterine Device', 1, 10),
    ('Procurement plan includes: Male Condom', 1, 11),
    ('Procurement plan includes: Sulfadoxine 500 + Pyrimethamine 25mg Tablets', 1, 12),
    ('Procurement plan includes: Ferrous Sulphate/Fumarate + Folic Acid Tablets', 1, 13),
    ('Procurement plan includes: Mebendazole 500mg Tablets', 1, 14),
    ('Procurement plan includes: Oxytocin Injection 10 i.u./mL', 1, 15),
    ('Procurement plan includes: Magnesium Sulphate Injection 50%', 1, 16),
    ('Procurement plan includes: Female condoms', 1, 17),
    ('Procurement plan includes: Moon beads', 1, 18)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Procurement plan contents'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Order integration', 11 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Order integration');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Facility submitted RH commodity orders alongside EMHS, Lab, ART, and TB orders', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Order integration'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Order timeliness', 12 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Order timeliness');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Most recent RH commodity order was submitted timely (before NMS/JMS order deadline)', 1, 1)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Order timeliness'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Order quality - last supplied order cycle', 13 FROM rhspars_domains d WHERE d.code='commodities' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Order quality - last supplied order cycle');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Order quality - Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: opening balances on both order forms agree', 1, 1),
    ('Order quality - Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: number of clients recorded on last RH order form', 1, 2),
    ('Order quality - Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection: client counts on RH order form vs FP register agree within +/- 10%', 1, 3),
    ('Order quality - Etonogestrel 68mg Implant (Implanon): opening balances on both order forms agree', 1, 4),
    ('Order quality - Etonogestrel 68mg Implant (Implanon): number of clients recorded on last RH order form', 1, 5),
    ('Order quality - Etonogestrel 68mg Implant (Implanon): client counts on RH order form vs FP register agree within +/- 10%', 1, 6),
    ('Order quality - Cu-T Intrauterine Device: opening balances on both order forms agree', 1, 7),
    ('Order quality - Cu-T Intrauterine Device: number of clients recorded on last RH order form', 1, 8),
    ('Order quality - Cu-T Intrauterine Device: client counts on RH order form vs FP register agree within +/- 10%', 1, 9),
    ('Order quality - Ferrous Sulphate/Fumarate + Folic Acid Tablets: opening balances on both order forms agree', 1, 10),
    ('Order quality - Ferrous Sulphate/Fumarate + Folic Acid Tablets: number of clients recorded on last RH order form', 1, 11),
    ('Order quality - Ferrous Sulphate/Fumarate + Folic Acid Tablets: client counts on RH order form vs FP register agree within +/- 10%', 1, 12),
    ('Order quality - Magnesium Sulphate Injection 50%: opening balances on both order forms agree', 1, 13),
    ('Order quality - Magnesium Sulphate Injection 50%: number of clients recorded on last RH order form', 1, 14),
    ('Order quality - Magnesium Sulphate Injection 50%: client counts on RH order form vs FP register agree within +/- 10%', 1, 15)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'commodities' AND ta.name = 'Order quality - last supplied order cycle'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_domains (code, name, display_order) SELECT 'health_information', 'Health Information', 5 WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='health_information');

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Availability and utilisation of HMIS tools', 1 FROM rhspars_domains d WHERE d.code='health_information' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Availability and utilisation of HMIS tools');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Family Planning Register: available', 1, 1),
    ('Family Planning Register: ten random entries completely filled (where applicable)', 1, 2),
    ('Family Planning Register: monthly summary totals present (where applicable)', 1, 3),
    ('Integrated ANC Register: available', 1, 4),
    ('Integrated ANC Register: ten random entries completely filled (where applicable)', 1, 5),
    ('Integrated ANC Register: monthly summary totals present (where applicable)', 1, 6),
    ('Integrated Maternity Register: available', 1, 7),
    ('Integrated Maternity Register: ten random entries completely filled (where applicable)', 1, 8),
    ('Integrated Maternity Register: monthly summary totals present (where applicable)', 1, 9),
    ('Postnatal Register: available', 1, 10),
    ('Postnatal Register: ten random entries completely filled (where applicable)', 1, 11),
    ('Postnatal Register: monthly summary totals present (where applicable)', 1, 12),
    ('Immunisation child register: available', 1, 13),
    ('Immunisation child register: ten random entries completely filled (where applicable)', 1, 14),
    ('Immunisation child register: monthly summary totals present (where applicable)', 1, 15),
    ('Monthly summaries report forms 105: available', 1, 16),
    ('Monthly summaries report forms 105: ten random entries completely filled (where applicable)', 1, 17),
    ('Monthly summaries report forms 105: monthly summary totals present (where applicable)', 1, 18),
    ('Community service registers and Summary Report 097b or eCHIS: available', 1, 19),
    ('Community service registers and Summary Report 097b or eCHIS: ten random entries completely filled (where applicable)', 1, 20),
    ('Community service registers and Summary Report 097b or eCHIS: monthly summary totals present (where applicable)', 1, 21),
    ('Dispensing log: available', 1, 22),
    ('Dispensing log: ten random entries completely filled (where applicable)', 1, 23),
    ('Dispensing log: monthly summary totals present (where applicable)', 1, 24),
    ('Electronic Medical Records EMR (e.g., eAFYA, Clinic Master, Ug-EMR): available', 1, 25),
    ('Electronic Medical Records EMR (e.g., eAFYA, Clinic Master, Ug-EMR): ten random entries completely filled (where applicable)', 1, 26),
    ('Electronic Medical Records EMR (e.g., eAFYA, Clinic Master, Ug-EMR): monthly summary totals present (where applicable)', 1, 27)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'health_information' AND ta.name = 'Availability and utilisation of HMIS tools'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Reporting - Completeness of HMIS 105', 2 FROM rhspars_domains d WHERE d.code='health_information' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Reporting - Completeness of HMIS 105');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Facility submitted reports for the two months prior to the supervision', 1, 1),
    ('Facility submitted information for ANC services in HMIS 105 Section 2.1', 1, 2),
    ('Facility submitted information for Maternity services in HMIS 105 Section 2.2', 1, 3),
    ('Facility submitted information for Family Planning services in HMIS 105 Section 2.4', 1, 4),
    ('Facility submitted information about dispensing of Family Planning methods in HMIS 105 Section 2.4.2', 1, 5),
    ('Facility submitted all data for all these commodities in HMIS 105 Section 6 (qty consumed, days out of stock, stock on hand, qty expired for: DMPA, SP tablets, Misoprostol 200mcg, Implanon, Oxytocin, Chlorhexidine Gel, Mama Kits)', 1, 6)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'health_information' AND ta.name = 'Reporting - Completeness of HMIS 105'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Accuracy of HMIS 105 Report', 3 FROM rhspars_domains d WHERE d.code='health_information' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Accuracy of HMIS 105 Report');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('HMIS accuracy - Number of Pregnant Women receiving at least 30 tablets of Folic Acid and Iron Sulphate at ANC 1st contact/visit (AN35): information available from the last report', 1, 1),
    ('HMIS accuracy - Number of Pregnant Women receiving at least 30 tablets of Folic Acid and Iron Sulphate at ANC 1st contact/visit (AN35): data agree or differ by no more than +/- 10%', 1, 2),
    ('HMIS accuracy - Number of pregnant women dewormed or receiving Mebendazole (AN34): information available from the last report', 1, 3),
    ('HMIS accuracy - Number of pregnant women dewormed or receiving Mebendazole (AN34): data agree or differ by no more than +/- 10%', 1, 4),
    ('HMIS accuracy - Total number of deliveries (MA04): information available from the last report', 1, 5),
    ('HMIS accuracy - Total number of deliveries (MA04): data agree or differ by no more than +/- 10%', 1, 6),
    ('HMIS accuracy - Total number (new + revisits) of Injectable DMPA (IM e.g., Depo) users (FP06): information available from the last report', 1, 7),
    ('HMIS accuracy - Total number (new + revisits) of Injectable DMPA (IM e.g., Depo) users (FP06): data agree or differ by no more than +/- 10%', 1, 8),
    ('HMIS accuracy - Total number (new + revisits) of users of 3 year implant e.g., Implanon NXT, Levoplant (FP10): information available from the last report', 1, 9),
    ('HMIS accuracy - Total number (new + revisits) of users of 3 year implant e.g., Implanon NXT, Levoplant (FP10): data agree or differ by no more than +/- 10%', 1, 10),
    ('HMIS accuracy - Total number (new + revisits) of users of IUD-Copper-T (FP13): information available from the last report', 1, 11),
    ('HMIS accuracy - Total number (new + revisits) of users of IUD-Copper-T (FP13): data agree or differ by no more than +/- 10%', 1, 12),
    ('HMIS accuracy - Quantity of DMPA issued in the store (SS02): information available from the last report', 1, 13),
    ('HMIS accuracy - Quantity of DMPA issued in the store (SS02): data agree or differ by no more than +/- 10%', 1, 14),
    ('HMIS accuracy - Quantity of Sulfadoxine/Pyrimethamine tablets issued in the store (SS04): information available from the last report', 1, 15),
    ('HMIS accuracy - Quantity of Sulfadoxine/Pyrimethamine tablets issued in the store (SS04): data agree or differ by no more than +/- 10%', 1, 16),
    ('HMIS accuracy - Quantity of Oxytocin issued in the store (SS31): information available from the last report', 1, 17),
    ('HMIS accuracy - Quantity of Oxytocin issued in the store (SS31): data agree or differ by no more than +/- 10%', 1, 18)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'health_information' AND ta.name = 'Accuracy of HMIS 105 Report'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_domains (code, name, display_order) SELECT 'general_facility', 'General facility', 6 WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='general_facility');

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Signage / poster', 1 FROM rhspars_domains d WHERE d.code='general_facility' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Signage / poster');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Signage - Family planning services: poster clearly shows availability of these services (English and local language)', 1, 1),
    ('Signage - Family planning services: clearly shows the day and time these services are available', 1, 2),
    ('Signage - Family planning services: services provided at least 5 days a week for at least 6 hours (maternity should be 24 hours)', 1, 3),
    ('Signage - Family planning services: clear within-facility directional signage to the clinic', 1, 4),
    ('Signage - Antenatal services: poster clearly shows availability of these services (English and local language)', 1, 5),
    ('Signage - Antenatal services: clearly shows the day and time these services are available', 1, 6),
    ('Signage - Antenatal services: services provided at least 5 days a week for at least 6 hours (maternity should be 24 hours)', 1, 7),
    ('Signage - Antenatal services: clear within-facility directional signage to the clinic', 1, 8),
    ('Signage - Maternity unit: poster clearly shows availability of these services (English and local language)', 1, 9),
    ('Signage - Maternity unit: clearly shows the day and time these services are available', 1, 10),
    ('Signage - Maternity unit: services provided at least 5 days a week for at least 6 hours (maternity should be 24 hours)', 1, 11),
    ('Signage - Maternity unit: clear within-facility directional signage to the clinic', 1, 12)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'general_facility' AND ta.name = 'Signage / poster'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Functionality of facility committees', 2 FROM rhspars_domains d WHERE d.code='general_facility' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Functionality of facility committees');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Health unit management committee / Hospital Management board: active committee with regular meetings (last 2 meetings)', 1, 1),
    ('Health unit management committee / Hospital Management board: MCH issues discussed in these meetings', 1, 2),
    ('Health unit management committee / Hospital Management board: RH supplies and commodities discussed', 1, 3),
    ('Senior management meeting: active committee with regular meetings (last 2 meetings)', 1, 4),
    ('Senior management meeting: MCH issues discussed in these meetings', 1, 5),
    ('Senior management meeting: RH supplies and commodities discussed', 1, 6),
    ('QI meeting: active committee with regular meetings (last 2 meetings)', 1, 7),
    ('QI meeting: MCH issues discussed in these meetings', 1, 8),
    ('QI meeting: RH supplies and commodities discussed', 1, 9)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'general_facility' AND ta.name = 'Functionality of facility committees'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Human resources - Staffing levels', 3 FROM rhspars_domains d WHERE d.code='general_facility' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Human resources - Staffing levels');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Obstetrician and gynaecologists: staffing level at or above 80% of recommended norm (per duty roster)', 1, 1),
    ('Medical officers: staffing level at or above 80% of recommended norm (per duty roster)', 1, 2),
    ('Midwives: staffing level at or above 80% of recommended norm (per duty roster)', 1, 3),
    ('Nurses: staffing level at or above 80% of recommended norm (per duty roster)', 1, 4),
    ('Anaesthetic officers: staffing level at or above 80% of recommended norm (per duty roster)', 1, 5),
    ('Assistant inventory management officer: staffing level at or above 80% of recommended norm (per duty roster)', 1, 6),
    ('Health information staff: staffing level at or above 80% of recommended norm (per duty roster)', 1, 7),
    ('Laboratory staff: staffing level at or above 80% of recommended norm (per duty roster)', 1, 8),
    ('Radiography / radiology staff: staffing level at or above 80% of recommended norm (per duty roster)', 1, 9)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'general_facility' AND ta.name = 'Human resources - Staffing levels'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) SELECT d.id, 'Staff performance management - Balanced score cards', 4 FROM rhspars_domains d WHERE d.code='general_facility' AND NOT EXISTS (SELECT 1 FROM rhspars_thematic_areas ta WHERE ta.domain_id=d.id AND ta.name='Staff performance management - Balanced score cards');

INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.display_order
FROM rhspars_thematic_areas ta
JOIN rhspars_domains d ON ta.domain_id = d.id,
(VALUES
    ('Obstetrician and gynaecologists: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)', 1, 1),
    ('Medical officers: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)', 1, 2),
    ('Midwives: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)', 1, 3),
    ('Nurses: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)', 1, 4),
    ('Anaesthetic officers: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)', 1, 5),
    ('Assistant inventory management officer: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)', 1, 6)
) AS q(question_text, score_weight, display_order)
WHERE d.code = 'general_facility' AND ta.name = 'Staff performance management - Balanced score cards'
AND NOT EXISTS (SELECT 1 FROM rhspars_questions x WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);

COMMIT;

-- Seeded 6 domains, 45 thematic areas, 333 questions