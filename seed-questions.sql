-- Comprehensive seeding script for all assessment questions
-- Based on the Family Planning Score Tool scoring keys

-- ============================================
-- COUNSELLING (Assessment Type ID: 1)
-- ============================================

-- Thematic Area 1: Maintains privacy and confidentiality
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Maintains privacy and confidentiality', 1)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Greets and employs a client-centred style of communication when speaking to clients', 2, false, false, 1),
    ('Uses language the client is comfortable with', 2, false, false, 2),
    ('Follows a structured counselling approach like REDI (Rapport, Explore, Decide and Implement)', 5, false, true, 3),
    ('Asks client about the service(s) they are seeking and if they have something specific in mind', 2, false, false, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Maintains privacy and confidentiality'
ON CONFLICT DO NOTHING;

-- Thematic Area 2: Provides comprehensive and correct information
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Provides comprehensive and correct information on service options that fits with client''s reproductive health needs and lifestyle preferences using integrated flipchart or other job aids appropriately', 2)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Explains condom use for dual protection to prevent STIs/HIV', 2, false, false, 1),
    ('Uses integrated flipchart or other job aids appropriately to explain service options', 5, false, true, 2),
    ('Supports clients to make own decisions after weighing up all information including advantages, disadvantages and consequences of each option without coercion', 5, false, true, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.display_order = 2
ON CONFLICT DO NOTHING;

-- Thematic Area 3: Explains how the chosen service would be provided
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Explains how the chosen service would be provided; its side effects and what to do in case of side effects which may require referral to higher level facility in rare cases', 3)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Checks that client fully understands the chosen option by asking them to repeat key points', 2, false, false, 1),
    ('Takes informed consent appropriately with consideration for vulnerable groups such as young people; illiterate; clients with history of sexual abuse or violence, long term physical, mental, intellectual or sensory impairments and mental health illness', 5, false, true, 2),
    ('Indicates how and where to access the chosen service', 2, false, false, 3),
    ('If appropriate, discusses scheduling; and interim or partner contraception', 2, false, false, 4),
    ('Informs the client that they could switch to another method if they wanted to or needed to', 2, false, false, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.display_order = 3
ON CONFLICT DO NOTHING;

-- Thematic Area 4: Provides information about other SRHR services
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Provides information about other SRHR services including (maternity care, skilled attendance at birth, postpartum FP, post abortion care, prevention of GBV, etc)?', 4)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Asks the client about her current or prior experience with family planning', 5, false, false, 1),
    ('Provider sensitively asks client whether she has experienced intimate partner violence or reproductive coercion (only asked in the absence of their partner)', 5, false, true, 2),
    ('Screens Client for NCDs (DM, Hypertension, sickle cell, cervical and breast cancer, prostate cancer, Tobacco use and exposure, alcohol and substance use, health diet and physical activity)', 5, false, true, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.display_order = 4
ON CONFLICT DO NOTHING;

-- Thematic Area 5: Assesses the client''s medical eligibility
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Assesses the client''s medical eligibility for the chosen method using the appropriate checklist', 5)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client history and documents all relevant conditions using the medical eligibility checklist', 5, false, true, 1),
    ('Rules out pregnancy using the WHO/MOH pregnancy checklist before proceeding', 5, false, true, 2),
    ('Screens for contraindications (hypertension, postpartum status, thrombosis, etc.) and documents findings', 5, false, true, 3),
    ('Records the required observations (e.g., BP, temperature, weight) needed to confirm eligibility', 2, false, false, 4),
    ('Documents the eligibility decision, counselling provided, and any referral or follow-up plan', 2, false, false, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 
  AND ta.name = 'Assesses the client''s medical eligibility for the chosen method using the appropriate checklist'
ON CONFLICT DO NOTHING;

-- Client Scenarios
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Return visit of client satisfied with FP method', 6),
(1, 'Post-abortion client', 7),
(1, 'Returning client who is not satisfied with her FP method', 8),
(1, 'Post-partum client', 9)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Assesses client''s satisfaction with chosen method and asks about her experience with common side effects and takes appropriate action', 5, false, true, 1),
    ('Probes and assesses for any new health condition that might contradict the continued use of the current method (Reference the MEC wheel)', 5, false, true, 2),
    ('Informs client about other methods', 2, false, false, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Return visit of client satisfied with FP method'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Assesses date of abortion and ensures appropriate Post abortion care', 2, false, false, 1),
    ('Rules out new pregnancy', 2, false, false, 2),
    ('Advises client to wait three menstrual cycles before trying to become pregnant', 2, false, false, 3),
    ('Informs the client that contraception can be started on the same day of abortion/ PAC treatment (Reference the MEC wheel)', 5, false, true, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Post-abortion client'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Takes history on FP method use and current concerns including GBV; Notes medications the client may be on (and potential drug interactions)', 5, false, true, 1),
    ('Provides information and explanation regarding experienced problems/side effects or concerns and manages appropriately', 5, false, true, 2),
    ('If client chooses to keep method, advises her when to return for a follow up appointment to ensure method satisfaction and continuation', 2, false, false, 3),
    ('Discusses method switch and gives appropriate options', 2, false, false, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Returning client who is not satisfied with her FP method'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Assesses date of child birth and rules out pregnancy', 2, false, false, 1),
    ('Explains that fertility can return as early as four weeks after child birth even if she is breastfeeding or has not seen her periods', 5, false, true, 2),
    ('Advises client to wait for atleast 2 years after childbirth before the next pregnancy and emphasizes the benefits', 2, false, false, 3),
    ('Informs client that she can take an FP method on the same day of child birth and takes them through the FP options (Refer to the FP compedium)', 5, false, true, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Post-partum client'
ON CONFLICT DO NOTHING;

-- Thematic Area 10: Self-injecting client
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Self-injecting client', 10)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Actively listens and uses open-ended questions to uncover client barriers to self-injection (Observe and see if the provider uses active listening techniques like paraphrasing, affirming, clarifying and empathy; non-verbal communication and open-ended questions to uncover client barriers to self inject e.g fear of pain, fear of the needle, lack of confidence to correctly inject, product safety concerns, covert use due to partner opposition )', 10, true, false, 1),
    ('Uses targeted messaging to address the identified barrier to self-injection expressed by the client (Observe and see if the provider uses empathetic messages to avert client''s fear to self injection)', 5, false, true, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Self-injecting client'
ON CONFLICT DO NOTHING;

-- Thematic Area 11: Storage and Disposal at Home
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(1, 'Storage and Disposal at Home: Assesses whether health workers are giving correct guidance about storage and disposal', 11)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Discusses with client home storage information of the units dispensed (Observe and listen: 1.Store at room temperature (do not refrigerate); 2.Store out of direct sunlight and heat; 3.Store out of reach of children and animals)', 5, false, true, 1),
    ('Discusses with client how to manage product waste at home (Look out for: 1. Do not touch the needle; 2. Do not recap the needle; 3. Dispose used needle in puncture-proof container at home e.g plastic bottle; 4. Return the container with used needles to the health facility when coming for refill)', 5, false, true, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 1 AND ta.name = 'Storage and Disposal at Home: Assesses whether health workers are giving correct guidance about storage and disposal'
ON CONFLICT DO NOTHING;

-- ============================================
-- OCPs (Assessment Type ID: 2)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(2, 'Initial Assessment Steps', 1),
(2, 'Client Eligibility for OCP (General)', 2),
(2, 'Additional Eligibility Requirements for COC Clients', 3),
(2, 'Instructional and Documentation Steps', 4)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client information and confirms client has been counselled', 2, false, false, 1),
    ('Completes a relevant history and clinical examination', 5, false, true, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 2 AND ta.name = 'Initial Assessment Steps'
ON CONFLICT DO NOTHING;

-- Eligibility questions (grouped - 10 points total)
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Is not currently pregnant or at risk of being pregnant using the Pregnancy Checklist', 10, true, false, 1),
    ('Does not have history of any vaginal bleeding that is unusual for them', 10, true, false, 2),
    ('Does not have liver disease', 10, true, false, 3),
    ('Does not have/has not had breast cancer', 10, true, false, 4),
    ('Does not have/has not had a blood clot in the legs or lung', 10, true, false, 5),
    ('Enquires about use of other medicines and correctly identifies those reducing effectiveness of the OCP', 10, true, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 2 AND ta.name = 'Client Eligibility for OCP (General)'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Has not just given birth or is not breastfeeding a baby less than 6 months', 10, true, false, 1),
    ('Does not get bad migraine headaches affecting vision or hearing', 10, true, false, 2),
    ('Does not have or is not being treated for high blood pressure and/or understands risk of taking COC without knowing blood pressure', 10, true, false, 3),
    ('Does not have complicated or long standing diabetes', 10, true, false, 4),
    ('Does not have heart disease or history of stroke', 10, true, false, 5),
    ('Is not a smoker aged over 35 years', 10, true, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 2 AND ta.name = 'Additional Eligibility Requirements for COC Clients'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Correctly advises start time of OCP depending on client''s state (last menstrual period, post abortion, post-partum with or without breastfeeding)', 5, false, true, 1),
    ('Checks the tablet pack for expiry date and ensures that its not damaged', 5, false, true, 2),
    ('Explains when client needs to return and also reviews warning signs', 5, false, true, 3),
    ('Confirms client understands what to expect, the need for daily use preferably at the same time for the pills to be effective, what to do if pills are missed and the importance of telling another prescriber about OCP use', 2, false, false, 4),
    ('Completes documentation of client records', 2, false, false, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 2 AND ta.name = 'Instructional and Documentation Steps'
ON CONFLICT DO NOTHING;

-- ============================================
-- INJECTABLE PROGESTERONE ONLY (Assessment Type ID: 3)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(3, 'Client Information and Eligibility (Absolute Contraindications)', 1),
(3, 'Administering Injection (Preparation Steps)', 2),
(3, 'Identifies Correct Site for the Injection Type', 3),
(3, 'Site Preparation', 4),
(3, 'Injection Procedure', 5),
(3, 'Post-Injection Care and Documentation', 6)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client information and confirms client has been counselled', 2, false, false, 1),
    ('Takes relevant detailed history and confirms that client is eligible for POI by checking that the client (Absolute contraindications): Is not currently pregnant or at risk of being pregnant using the Pregnancy Checklist', 10, true, false, 2),
    ('Does not have unexplained vaginal bleeding', 10, true, false, 3),
    ('Is not breast feeding a baby that is less than 6 weeks old', 10, true, false, 4),
    ('Does not currently have a deep vein thrombosis (blood clot in their leg) or pulmonary embolism (blood clot in their lung)', 10, true, false, 5),
    ('Does not have high blood pressure (systolic ≥160 or diastolic ≥100 mm Hg)', 10, true, false, 6),
    ('Does not have multiple risk factors for cardiovascular disease (such as smoking, diabetes, high blood pressure, obesity or high cholesterol)', 10, true, false, 7),
    ('Does not have history of ischaemic heart disease or a cerebrovascular accident (stroke)', 10, true, false, 8),
    ('Does not have systemic lupus erythematosus (SLE) and positive (or unknown) antiphospholipid antibodies or severe thrombocytopenia', 10, true, false, 9),
    ('Does not have diabetes of > 20 years'' duration or diabetes with nephropathy/retinopathy/neuropathy', 10, true, false, 10),
    ('Does not have/has not had severe decompensated cirrhosis or liver cancer (hepatocellular adenoma or malignant hepatoma)', 10, true, false, 11),
    ('Does not have/has not had breast cancer', 10, true, false, 12)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Client Information and Eligibility (Absolute Contraindications)'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Checks expiry date and integrity of injection', 2, false, false, 1),
    ('If DMPA IM: Mixes solution by rolling it between palms or shakes gently until content is frothy', 5, false, true, 2),
    ('IF DMPA IM: Removes cap and exposes rubber cover; empties content of bottle fully into syringe and ensures that all air is expelled from syringe', 2, false, false, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Administering Injection (Preparation Steps)'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Identifies correct site for Intramuscular: upper outer quadrant of buttocks/three fingers down from the top of the shoulder', 10, true, false, 1),
    ('Identifies correct site for Subcutaneous: thigh/abdomen/upper arm', 10, true, false, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Identifies Correct Site for the Injection Type'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('IM: Cleans area with antiseptic', 5, false, true, 1),
    ('SC: cleans the site if visibly dirty', 5, false, true, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Site Preparation'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('IM: Gently shakes the vial before loading syringe, aspirates and injects the medicine slowly with needle at 90°to the skin taking 5-7 seconds before withdrawing', 10, true, false, 1),
    ('SC: Holds the port of injector and shakes the injector for 30 seconds, activates the injector correctly by pushing the needle shield and port together, closing the gap, gently grasps and squeezes a large area of skin, pushes the needle straight into the skin with the needle pointing down and empties the reservoir by squeezing firmly but slowly for 5 - 7 seconds', 10, true, false, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Injection Procedure'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Presses injection site gently with clean cotton wool after removal of needle (no massaging of the injection site)', 2, false, false, 1),
    ('Follows Universal infection prevention principles at all times', 10, true, false, 2),
    ('Completes documentation of client records', 5, false, true, 3),
    ('Provides post injection counselling, confirms client understands duration of use, what to expect, need for regular injections and date of next injection', 2, false, false, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Post-Injection Care and Documentation'
ON CONFLICT DO NOTHING;

-- Thematic Area 7: Self-injection (SI) training and follow-up
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(3, 'Self-injection (SI) training and follow-up', 7)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('The provider walks the client through all the 4 critical steps (Mix, Activate, Pinch, slow pressing to release the drug) for SI provision (Observe the provider takes the client through all the 4 steps)', 10, true, false, 1),
    ('Provider uses models and instruction sheets to demonstrate SI (Check if the provider uses salt filled condom models and the self inject instructional sheet to demonstrate)', 2, false, false, 2),
    ('*Provider trains client to use the calendar to calculate reinjection dates (Observe)', 5, false, true, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Self-injection (SI) training and follow-up'
ON CONFLICT DO NOTHING;

-- Thematic Area 8: Handling clients not ready for independent SI
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(3, 'Health provider appropriately handles client who is not ready for independent self-injection', 8)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Health worker gives the DMPA-SC injection, then asks the client to return for SI retraining at their next injection.', 2, false, true, 1),
    ('Health worker does not give out units until the client has demonstrated readiness.', 2, false, true, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Health provider appropriately handles client who is not ready for independent self-injection'
ON CONFLICT DO NOTHING;

-- Thematic Area 9: Returning self-injection clients
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(3, 'Information and discussion areas for returning self-injection clients', 9)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Evaluates whether the client is experiencing any problems, including side effects.', 2, false, false, 1),
    ('Screens for eligibility and willingness to continue with self-injection.', 2, false, false, 2),
    ('Reviews the 4 critical injection steps and any questions about self-injection.', 2, false, false, 3),
    ('Provides additional training/guidance on injection or reinjection timing as needed.', 2, false, false, 4),
    ('Provides recommended number of DMPA-SC units if adequate stock available and any other applicable supplies.', 2, false, false, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 3 AND ta.name = 'Information and discussion areas for returning self-injection clients'
ON CONFLICT DO NOTHING;

-- ============================================
-- IMPLANT INSERTION (Assessment Type ID: 4)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(4, 'Pre-Procedure', 1),
(4, 'Procedure: Implant Insertion', 2),
(4, 'Post-Procedure', 3)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client information; and confirms client has been counselled and informed consent', 2, false, false, 1),
    ('Takes relevant detailed history and confirms that client is eligible for implant by checking that the client: Is not currently pregnant or at risk of being pregnant using the Pregnancy Checklist', 10, true, false, 2),
    ('Does not have any history of vaginal bleeding that is unusual for them', 10, true, false, 3),
    ('Does not have a serious liver disease (yellow eyes and skin) needing treatment', 10, true, false, 4),
    ('Does not have/has not had breast cancer', 10, true, false, 5),
    ('Does not have problems with a blood clot in the legs or lung', 10, true, false, 6),
    ('Does not have a rheumatic disease such as lupus', 10, true, false, 7),
    ('Ensures necessary equipment, processed instruments and supplies are ready', 2, false, false, 8),
    ('Ensures implant packs are undamaged and not beyond expiry date', 5, false, true, 9)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 4 AND ta.name = 'Pre-Procedure'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Identifies and Marks insertion points appropriately (The non-dominant arm, 8-10 cm from the medial epicondyle and 3-5 cm below the sulcus)', 10, true, false, 1),
    ('Cleans area with antiseptic and waits for it to dry', 5, false, true, 2),
    ('Covers area with sterile drape', 2, false, false, 3),
    ('Injects local anesthetic (1% without epinephrine) just under skin; raises a small wheal. Advances needle about 4 cm and injects 1mls of local anesthetic for one rod implant and 2 mls for two rods subdermal tracks', 10, true, false, 4),
    ('Palpates ends of rods to be sure the rods are placed correctly ("V") & palpates incision to check that ends of rods are about 5mm away from incision', 10, true, false, 5),
    ('Verifies presence of implant by palpation', 10, true, false, 6),
    ('Brings edges of incision together and closes it with band aids or plaster with sterile cotton', 5, false, true, 7),
    ('Lets the client verify that she can feel the Implant herself', 10, true, false, 8)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 4 AND ta.name = 'Procedure: Implant Insertion'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Assesses the need and provides pain management accordingly', 5, false, true, 1),
    ('Follows universal Infection prevention principles at all times', 10, true, false, 2),
    ('Completes documentation of client records', 5, false, true, 3),
    ('Ensures completed Implant card is given to client', 2, false, false, 4),
    ('Post insertion counselling; Confirms client understands post procedure instructions, especially what to expect (including common side-effects and what do), warning signs and duration of use; and ensures client has contact details for any emergency related to the service', 5, false, true, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 4 AND ta.name = 'Post-Procedure'
ON CONFLICT DO NOTHING;

-- ============================================
-- IMPLANT REMOVAL (Assessment Type ID: 5)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(5, 'Pre-removal counseling', 1),
(5, 'Removal of implant rod(s)', 2),
(5, 'Post-removal', 3)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Asks client her reason for removal and answers any questions', 2, false, false, 1),
    ('Establishes reason for removal and where necessary explores options for switching', 2, false, false, 2),
    ('Describes the removal procedure and what to expect', 2, false, false, 3),
    ('Prepares removal equipment', 2, false, false, 4),
    ('Confirms the position of each rod by making a mark at both ends of the rod(s)', 10, true, false, 5),
    ('Puts sterile or high level disinfected gloves on both hands', 5, false, true, 6),
    ('Prepares removal site with antiseptic solution', 2, false, false, 7),
    ('Places sterile drape over arm', 2, false, false, 8),
    ('Injects 1mls for 1 rod, 2mls for 2 rods of local anesthetic (1% without epinephrine) at the incision site and under the end of the capsule(s) and waits 2 mins. Checks for anesthetic effect before making skin incision', 5, false, true, 9)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 5 AND ta.name = 'Pre-removal counseling'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Removes the remaining rod using the same technique', 10, true, false, 1),
    ('Shows client the implant(s) and disposes as hazardous waste', 2, false, false, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 5 AND ta.name = 'Removal of implant rod(s)'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Processes instruments and other consumables appropriately', 2, false, false, 1),
    ('Assesses and manages pain appropriately', 2, false, false, 2),
    ('Instructs client regarding wound care and makes return visit appointment, if necessary', 2, false, false, 3),
    ('Offers another contraceptive method to the client if necessary', 2, false, false, 4),
    ('Completes documentation of client records', 5, false, true, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 5 AND ta.name = 'Post-removal'
ON CONFLICT DO NOTHING;

-- ============================================
-- IUD/IUS INSERTION (Assessment Type ID: 6)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(6, 'Pre-Procedure: IUD/IUS Insertion', 1),
(6, 'For post partum IUD', 2),
(6, 'Pre insertion', 3),
(6, 'Procedure: IUD/IUS Insertion', 4),
(6, 'Post-Procedure', 5)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client information; and confirms client has been counselled appropriately', 2, false, false, 1),
    ('Is not currently pregnant or at risk of being pregnant using the Pregnancy Checklist', 10, true, false, 2),
    ('Does not have any history of vaginal bleeding that is unusual for them', 10, true, false, 3),
    ('Has not had a baby over 48 hours ago and in the last 4 weeks', 10, true, false, 4),
    ('Does not have or is unlikely to have a genital or pelvic infection', 10, true, false, 5),
    ('Does not have/has not had genital cancer', 10, true, false, 6),
    ('Does not have/has not had breast cancer (IUS only)', 10, true, false, 7),
    ('Does not have/has not had blood clots (IUS only)', 10, true, false, 8),
    ('Does not have/has not had liver disease (IUS only)', 10, true, false, 9)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 6 AND ta.name = 'Pre-Procedure: IUD/IUS Insertion'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews the course of her labor/delivery & confirms there was no: prolonged Rupture of membranes (>24 hours), prolonged labor (>24 hours), fever (>38 c/100.4f), Intrapartum hemorrhage, extensive genital trauma and chorioamnionitis', 10, true, false, 1),
    ('Ensures necessary equipment, processed instruments and supplies are ready, IUD/IUS packs are undamaged and not beyond expiry date', 5, false, true, 2)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 6 AND ta.name = 'For post partum IUD'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Makes sure that the client has emptied her bladder', 2, false, false, 1),
    ('Performs abdominal exam - palpates abdomen and checks for suprapubic or pelvic tenderness', 2, false, false, 2),
    ('Performs vaginal exam - checking for ulcers, lesions, sores, or discharge vaginally or at the vulva', 2, false, false, 3),
    ('Performs bimanual exam checking for cervical, adnexal, or uterine abnormalities that would preclude insertion', 2, false, false, 4),
    ('Removes and disposes of gloves correctly, then puts new sterile exam gloves on both hands', 2, false, false, 5),
    ('Inserts speculum and performs exam, locates cervix checking for any signs of cervical or vaginal problems that might preclude insertion at this time including cervical cancer screening', 2, false, false, 6),
    ('Makes appropriate decision on proceeding with insertion and communicates with client', 2, false, false, 7)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 6 AND ta.name = 'Pre insertion'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Gently inserts and releases loaded IUD in the uterus using withdrawal technique', 10, true, false, 1),
    ('Gently inserts IUS using non-touch technique till 1-2 cm short of assessed uterine depth', 10, true, false, 2),
    ('Grasps the anterior lip of the cervix with a ring forceps. (Does not use a toothed tenaculum, as it may tear the cervix.) Closes one notch only', 10, true, false, 3),
    ('Pulls down the introitus with two fingers and visualize the interior of the', 10, true, false, 4),
    ('Cuts strings with scissors 3-4 cm from cervical os', 2, false, false, 5),
    ('Examines cervix for bleeding before removing speculum', 2, false, false, 6),
    ('Explains and demonstrates self-checking of IUD/IUS threads', 2, false, false, 7)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 6 AND ta.name = 'Procedure: IUD/IUS Insertion'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Assesses the need and provides appropriate pain management', 2, false, false, 1),
    ('Follows universal infection prevention principles at all times', 10, true, false, 2),
    ('Completes documentation of client records', 5, false, true, 3),
    ('Confirms client understands post procedure instructions, especially what to expect (including common side effects and what to do), warning signs and duration of use and ensures client has contact details for any emergency related to the service', 5, false, true, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 6 AND ta.name = 'Post-Procedure'
ON CONFLICT DO NOTHING;

-- ============================================
-- IUD/IUS REMOVAL (Assessment Type ID: 7)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(7, 'Pre-Procedure: IUD/IUS Removal', 1),
(7, 'Procedure: IUD/IUS Removal', 2),
(7, 'Post-Procedure', 3)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client information and asks client their reason for removal; confirms they have been counselled', 2, false, false, 1),
    ('Completes a relevant history', 2, false, false, 2),
    ('Ensures necessary equipment, processed instruments and supplies are ready', 5, false, true, 3)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 7 AND ta.name = 'Pre-Procedure: IUD/IUS Removal'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Confirms client has emptied bladder recently', 2, false, false, 1),
    ('Use Stopes forceps to stabilize cervix if required after cleaning cervix with an antiseptic', 10, true, false, 2),
    ('Shows client the removed IUD/IUS', 2, false, false, 3),
    ('Offers another contraceptive method to the client if necessary', 2, false, false, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 7 AND ta.name = 'Procedure: IUD/IUS Removal'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Assesses the need and provides appropriate pain management', 2, false, false, 1),
    ('Follows universal infection prevention principles at all times', 10, true, false, 2),
    ('Completes documentation of client records: Procedure notes including date of removal and details of adverse events if any', 5, false, true, 3),
    ('Confirms client understands post procedure instructions', 2, false, false, 4)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 7 AND ta.name = 'Post-Procedure'
ON CONFLICT DO NOTHING;

-- ============================================
-- MINI-LAP (Assessment Type ID: 8)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(8, 'Pre-Procedure', 1),
(8, 'Procedure: Mini-Laparotomy Tubal Ligation', 2),
(8, 'Post-Procedure', 3)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Reviews client information; and confirms client has been counseled and informed consent has been documented', 2, false, false, 1),
    ('Confirms to the client the permanent nature of the operation', 5, false, true, 2),
    ('Takes/reviews relevant detailed history', 2, false, false, 3),
    ('Confirms that client is eligible for tubal ligation', 10, true, false, 4),
    ('Ensures necessary equipment, processed instruments and supplies are ready', 10, true, false, 5),
    ('Washes hands thoroughly and dries them', 2, false, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 8 AND ta.name = 'Pre-Procedure'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Gently retracts Stopes forceps/Millin Vulsellum forceps and inserts uterine elevator', 10, true, false, 1),
    ('Aspirates and infiltrates skin, subcutaneous tissue, fascia, and peritoneum with local anesthetic appropriately', 10, true, false, 2),
    ('Grasps nicked fascia with Allis forceps and extends opening on both sides', 10, true, false, 3),
    ('Retrieving Fallopian Tubes', 10, true, false, 4),
    ('Occluding Tubes', 10, true, false, 5),
    ('Closing Abdomen', 10, true, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 8 AND ta.name = 'Procedure: Mini-Laparotomy Tubal Ligation'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Ensures vital signs and client condition are monitored and recorded post-procedure', 5, false, true, 1),
    ('Provides appropriate pain management', 5, false, true, 2),
    ('Follows universal infection prevention principles at all times', 10, true, false, 3),
    ('Completes documentation of client records', 5, false, true, 4),
    ('Assesses client if ready for discharge', 2, false, false, 5)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 8 AND ta.name = 'Post-Procedure'
ON CONFLICT DO NOTHING;

-- ============================================
-- VASECTOMY (Assessment Type ID: 9)
-- ============================================

INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
(9, 'Getting Ready', 1),
(9, 'Pre-operative Tasks', 2),
(9, 'local anaesthesia', 3),
(9, 'Procedure: No-Scalpel Vasectomy', 4),
(9, 'Occlusion by ligation with excision and fascial interposition', 5),
(9, 'Occlusion by Cautery', 6),
(9, 'Post-operative Tasks', 7)
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Ensures that client has been appropriately counseled about the procedure', 2, false, false, 1),
    ('Verifies client''s identity and checks that informed consent was obtained', 5, false, true, 2),
    ('Explains to the client the effectiveness and permanence of the procedure', 5, false, true, 3),
    ('Ensures that client has emptied bladder', 2, false, false, 4),
    ('Takes medical history and performs heart, lung & abdominal examination', 2, false, false, 5),
    ('Asks the client if he is allergic to any antiseptic solution or local anesthetics', 2, false, false, 6),
    ('Checks that client has thoroughly washed and trimmed the front portion of the scrotum', 5, false, true, 7)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'Getting Ready'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Ensures warm temperature in the room to relax the scrotum', 2, false, false, 1),
    ('Determines that sterile instruments: (ringed clamp, dissecting forceps and straight scissors) and supplies are available', 10, true, false, 2),
    ('Performs genital examination: e.g. palpating the scrotum to assess the thickness of the scrotal skin and if spermatic cords are mobile', 5, false, true, 3),
    ('Gently washes the scrotum with a warm antiseptic solution', 5, false, true, 4),
    ('Puts the penis in a 12 o''clock position on the man''s abdomen, so that the median raphe is clearly visible', 2, false, false, 5),
    ('Ensures sterile drapes are properly placed', 2, false, false, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'Pre-operative Tasks'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Administers local anesthesia correctly using the 3-finger technique, creating skin wheal, and injecting anesthetic into spermatic fascia', 10, true, false, 1)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'local anaesthesia'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Performs no-scalpel vasectomy procedure correctly using three-finger method, ringed clamp, and proper technique', 10, true, false, 1)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'Procedure: No-Scalpel Vasectomy'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Performs occlusion by ligation with excision and fascial interposition correctly', 10, true, false, 1)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'Occlusion by ligation with excision and fascial interposition'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Performs occlusion by cautery correctly', 10, true, false, 1)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'Occlusion by Cautery'
ON CONFLICT DO NOTHING;

INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
SELECT ta.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
FROM thematic_areas ta,
(VALUES
    ('Ensures the disposal of waste materials and sharps is in accordance with infection prevention guidelines', 2, false, false, 1),
    ('Provides client with written post-operative instructions and information when and where to return for follow-up', 2, false, false, 2),
    ('Reviews instructions verbally and asks if client has any questions', 2, false, false, 3),
    ('Reviews the need for back-up contraception for at least 3 months and provides client with condoms if needed', 2, false, false, 4),
    ('Advises client to return for semen analysis after 3 months', 2, false, false, 5),
    ('Completes registers and other record keeping tools/documentation', 5, false, true, 6)
) AS q(question_text, score_weight, is_critical, is_important, display_order)
WHERE ta.assessment_type_id = 9 AND ta.name = 'Post-operative Tasks'
ON CONFLICT DO NOTHING;

