-- Seed demo users and roles for local / Ubuntu setups
-- Prerequisites: schema.sql, seed-data.sql (regions), and preferably seed-tool-roles.sql
--
-- Accounts (password in plain text for operators — stored hashed below):
--   admin@fpscore.local          / Admin@123          → full Administrator (all permissions)
--   proficiency@fpscore.local    / Proficiency@123    → FP Proficiency Score Tool only
--   rhspars@fpscore.local        / Rhspars@123        → RH SPARS Tool only

BEGIN;

-- ---------------------------------------------------------------------------
-- Permissions (idempotent)
-- ---------------------------------------------------------------------------
INSERT INTO permissions (code, description)
SELECT v.code, v.description FROM (VALUES
  ('assessments.create', 'Create new assessments'),
  ('assessments.view', 'View assessments'),
  ('assessments.edit', 'Edit assessments'),
  ('assessments.delete', 'Delete assessments'),
  ('users.create', 'Create new users'),
  ('users.view', 'View users'),
  ('users.edit', 'Edit users'),
  ('users.delete', 'Delete users'),
  ('roles.create', 'Create new roles'),
  ('roles.view', 'View roles'),
  ('roles.edit', 'Edit roles'),
  ('roles.delete', 'Delete roles'),
  ('reports.view', 'View reports'),
  ('reports.export', 'Export reports'),
  ('dashboard.view', 'View dashboard'),
  ('facilities.view', 'View facilities'),
  ('facilities.manage', 'Manage facilities'),
  ('admin.areas.view', 'View administrative areas'),
  ('admin.areas.manage', 'Manage administrative areas'),
  ('hierarchy.facilities.manage', 'Manage facility hierarchy operations (move facilities between subcounties)'),
  ('health_workers.view', 'View health workers'),
  ('health_workers.create', 'Create health workers'),
  ('health_workers.edit', 'Edit health workers'),
  ('health_workers.delete', 'Delete health workers'),
  ('logs.view', 'View system logs'),
  ('tools.proficiency.access', 'Access the FP Proficiency Score Tool'),
  ('tools.rh_spars.access', 'Access the RH SPARS (Integrated Reproductive Health Support Supervision) Tool'),
  ('rhspars.create', 'Create RH SPARS assessments'),
  ('rhspars.view', 'View RH SPARS assessments'),
  ('rhspars.edit', 'Edit RH SPARS assessments'),
  ('rhspars.delete', 'Delete RH SPARS assessments')
) AS v(code, description)
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.code = v.code);

-- ---------------------------------------------------------------------------
-- Roles
-- ---------------------------------------------------------------------------
INSERT INTO roles (name, description)
SELECT 'Administrator', 'Full system access — all permissions'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'Administrator');

INSERT INTO roles (name, description)
SELECT 'Proficiency Score Assessor', 'Access FP proficiency assessments only'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'Proficiency Score Assessor');

INSERT INTO roles (name, description)
SELECT 'RH SPARS Assessor', 'Access RH SPARS supervision assessments only'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'RH SPARS Assessor');

-- Administrator: every permission currently in the table
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Administrator'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- Proficiency Score Assessor
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Proficiency Score Assessor'
  AND p.code IN (
    'tools.proficiency.access',
    'assessments.create', 'assessments.view',
    'dashboard.view', 'reports.view', 'reports.export',
    'health_workers.view', 'health_workers.create', 'health_workers.edit',
    'admin.areas.view', 'facilities.view'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- RH SPARS Assessor
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'RH SPARS Assessor'
  AND p.code IN (
    'tools.rh_spars.access',
    'rhspars.create', 'rhspars.view',
    'admin.areas.view', 'facilities.view'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- ---------------------------------------------------------------------------
-- Users (bcrypt hashes for the passwords listed in the header)
-- ---------------------------------------------------------------------------
INSERT INTO users (name, email, password_hash, is_active)
SELECT 'System Administrator', 'admin@fpscore.local',
       '$2a$10$vDUe/5rN/9FPE2SZPuHbyexWEdHAsrYTTGcLoRJ8UUq5ryI0Rx47W', true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@fpscore.local');

INSERT INTO users (name, email, password_hash, is_active)
SELECT 'Proficiency Assessor', 'proficiency@fpscore.local',
       '$2a$10$LMFHYbS0LirWDZOlbKdlJOFNxBMm2ooD1ux3gg2Pi2m8x7dG4VcrW', true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'proficiency@fpscore.local');

INSERT INTO users (name, email, password_hash, is_active)
SELECT 'RH SPARS Assessor', 'rhspars@fpscore.local',
       '$2a$10$8aqbLJurDs0tP4xpN5Ny.uKMohV33AziAVP8bkzx1mpWoXohd6R/a', true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'rhspars@fpscore.local');

-- Role assignments
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.name = 'Administrator'
WHERE u.email = 'admin@fpscore.local'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = r.id
  );

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.name = 'Proficiency Score Assessor'
WHERE u.email = 'proficiency@fpscore.local'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = r.id
  );

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.name = 'RH SPARS Assessor'
WHERE u.email = 'rhspars@fpscore.local'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = r.id
  );

-- ---------------------------------------------------------------------------
-- Admin areas for non-admin users (required or they see no scoped data)
-- Proficiency → Central Region; RH SPARS → Eastern Region
-- ---------------------------------------------------------------------------
INSERT INTO user_admin_areas (user_id, region_id)
SELECT u.id, r.id
FROM users u
JOIN regions r ON r.name = 'Central Region'
WHERE u.email = 'proficiency@fpscore.local'
  AND NOT EXISTS (
    SELECT 1 FROM user_admin_areas uaa
    WHERE uaa.user_id = u.id AND uaa.region_id = r.id
      AND uaa.district_id IS NULL AND uaa.subcounty_id IS NULL AND uaa.facility_id IS NULL
  );

INSERT INTO user_admin_areas (user_id, region_id)
SELECT u.id, r.id
FROM users u
JOIN regions r ON r.name = 'Eastern Region'
WHERE u.email = 'rhspars@fpscore.local'
  AND NOT EXISTS (
    SELECT 1 FROM user_admin_areas uaa
    WHERE uaa.user_id = u.id AND uaa.region_id = r.id
      AND uaa.district_id IS NULL AND uaa.subcounty_id IS NULL AND uaa.facility_id IS NULL
  );

COMMIT;
