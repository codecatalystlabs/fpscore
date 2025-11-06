-- Sample data for testing the Family Planning Score Tool
-- Run this after creating the schema

-- Insert sample regions
INSERT INTO regions (name) VALUES
    ('Central Region'),
    ('Eastern Region'),
    ('Northern Region'),
    ('Western Region')
ON CONFLICT (name) DO NOTHING;

-- Insert sample districts
INSERT INTO districts (region_id, name) VALUES
    ((SELECT id FROM regions WHERE name = 'Central Region'), 'Kampala District'),
    ((SELECT id FROM regions WHERE name = 'Central Region'), 'Wakiso District'),
    ((SELECT id FROM regions WHERE name = 'Eastern Region'), 'Jinja District'),
    ((SELECT id FROM regions WHERE name = 'Eastern Region'), 'Mbale District'),
    ((SELECT id FROM regions WHERE name = 'Northern Region'), 'Gulu District'),
    ((SELECT id FROM regions WHERE name = 'Western Region'), 'Mbarara District')
ON CONFLICT (region_id, name) DO NOTHING;

-- Insert sample subcounties
INSERT INTO subcounties (district_id, name) VALUES
    ((SELECT id FROM districts WHERE name = 'Kampala District'), 'Kampala Central'),
    ((SELECT id FROM districts WHERE name = 'Kampala District'), 'Makindye'),
    ((SELECT id FROM districts WHERE name = 'Wakiso District'), 'Entebbe'),
    ((SELECT id FROM districts WHERE name = 'Wakiso District'), 'Nansana'),
    ((SELECT id FROM districts WHERE name = 'Jinja District'), 'Jinja Central'),
    ((SELECT id FROM districts WHERE name = 'Mbale District'), 'Mbale Central'),
    ((SELECT id FROM districts WHERE name = 'Gulu District'), 'Gulu Central'),
    ((SELECT id FROM districts WHERE name = 'Mbarara District'), 'Mbarara Central')
ON CONFLICT (district_id, name) DO NOTHING;

-- Insert sample facilities
INSERT INTO facilities (subcounty_id, name) VALUES
    ((SELECT id FROM subcounties WHERE name = 'Kampala Central'), 'Mulago National Referral Hospital'),
    ((SELECT id FROM subcounties WHERE name = 'Kampala Central'), 'Kampala City Clinic'),
    ((SELECT id FROM subcounties WHERE name = 'Makindye'), 'Makindye Health Centre'),
    ((SELECT id FROM subcounties WHERE name = 'Entebbe'), 'Entebbe Hospital'),
    ((SELECT id FROM subcounties WHERE name = 'Nansana'), 'Nansana Health Centre'),
    ((SELECT id FROM subcounties WHERE name = 'Jinja Central'), 'Jinja Regional Hospital'),
    ((SELECT id FROM subcounties WHERE name = 'Mbale Central'), 'Mbale Regional Hospital'),
    ((SELECT id FROM subcounties WHERE name = 'Gulu Central'), 'Gulu Regional Hospital'),
    ((SELECT id FROM subcounties WHERE name = 'Mbarara Central'), 'Mbarara Regional Hospital')
ON CONFLICT (subcounty_id, name) DO NOTHING;

-- Note: You'll need to manually insert thematic areas and questions based on your scoring key documents
-- Example structure for adding questions:

-- For Counselling Assessment Type (ID 1):
-- INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
--     (1, 'Maintains privacy and confidentiality', 1),
--     (1, 'Provides comprehensive and correct information', 2),
--     (1, 'Explains how the chosen service would be provided', 3),
--     (1, 'Provides information about other SRHR services', 4),
--     (1, 'Assesses the client''s medical eligibility', 5)
-- ON CONFLICT DO NOTHING;

-- INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order) VALUES
--     ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
--      'Greets and employs a client-centred style of communication when speaking to clients', 2, false, false, 1),
--     ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
--      'Uses language the client is comfortable with', 2, false, false, 2),
--     ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
--      'Follows a structured counselling approach like REDI (Rapport, Explore, Decide and Implement)', 5, false, true, 3)
-- ON CONFLICT DO NOTHING;

