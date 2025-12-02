-- Seed health workers data
-- This script creates sample health workers for testing

-- Insert sample health workers
-- Mulago National Referral Hospital
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Sarah Nakato', 'sarah.nakato@example.com', '+256700000001', 
     (SELECT id FROM facilities WHERE name = 'Mulago National Referral Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. James Ochieng', 'james.ochieng@example.com', '+256700000002',
     (SELECT id FROM facilities WHERE name = 'Mulago National Referral Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Mary Achieng', 'mary.achieng@example.com', '+256700000003',
     (SELECT id FROM facilities WHERE name = 'Mulago National Referral Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

-- Kampala City Clinic
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Peter Okello', 'peter.okello@example.com', '+256700000004',
     (SELECT id FROM facilities WHERE name = 'Kampala City Clinic' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Grace Atim', NULL, '+256700000005',
     (SELECT id FROM facilities WHERE name = 'Kampala City Clinic' LIMIT 1))
ON CONFLICT (phone_number) DO NOTHING;

-- Makindye Health Centre
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. John Mukasa', 'john.mukasa@example.com', '+256700000006',
     (SELECT id FROM facilities WHERE name = 'Makindye Health Centre' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Jane Nakato', NULL, '+256700000007',
     (SELECT id FROM facilities WHERE name = 'Makindye Health Centre' LIMIT 1))
ON CONFLICT (phone_number) DO NOTHING;

-- Entebbe Hospital
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Alice Namukasa', 'alice.namukasa@example.com', '+256700000008',
     (SELECT id FROM facilities WHERE name = 'Entebbe Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Rose Nalubega', 'rose.nalubega@example.com', NULL,
     (SELECT id FROM facilities WHERE name = 'Entebbe Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

-- Nansana Health Centre
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Robert Ssemwogerere', 'robert.ssemwogerere@example.com', '+256700000009',
     (SELECT id FROM facilities WHERE name = 'Nansana Health Centre' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

-- Jinja Regional Hospital
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Susan Nalubega', 'susan.nalubega@example.com', '+256700000010',
     (SELECT id FROM facilities WHERE name = 'Jinja Regional Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Betty Nakawuki', NULL, '+256700000011',
     (SELECT id FROM facilities WHERE name = 'Jinja Regional Hospital' LIMIT 1))
ON CONFLICT (phone_number) DO NOTHING;

-- Mbale Regional Hospital
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. David Waiswa', 'david.waiswa@example.com', '+256700000012',
     (SELECT id FROM facilities WHERE name = 'Mbale Regional Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Florence Nabukeera', 'florence.nabukeera@example.com', NULL,
     (SELECT id FROM facilities WHERE name = 'Mbale Regional Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

-- Gulu Regional Hospital
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Michael Ocen', 'michael.ocen@example.com', '+256700000013',
     (SELECT id FROM facilities WHERE name = 'Gulu Regional Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

-- Mbarara Regional Hospital
INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Dr. Elizabeth Tumwine', 'elizabeth.tumwine@example.com', '+256700000014',
     (SELECT id FROM facilities WHERE name = 'Mbarara Regional Hospital' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

INSERT INTO health_workers (full_name, email, phone_number, facility_id) VALUES
    ('Nurse Agnes Kyomugisha', NULL, '+256700000015',
     (SELECT id FROM facilities WHERE name = 'Mbarara Regional Hospital' LIMIT 1))
ON CONFLICT (phone_number) DO NOTHING;

-- Note: Some health workers have only email, some have only phone, some have both
-- This demonstrates the flexibility of the unique identifier system

