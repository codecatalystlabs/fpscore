-- Tool access permissions + starter roles
-- Run after permissions initialize (or rely on app InitializePermissions on boot).

BEGIN;

-- Ensure tool + RH SPARS permissions exist
INSERT INTO permissions (code, description)
SELECT v.code, v.description FROM (VALUES
  ('tools.proficiency.access', 'Access the FP Proficiency Score Tool'),
  ('tools.rh_spars.access', 'Access the RH SPARS (Integrated Reproductive Health Support Supervision) Tool'),
  ('rhspars.create', 'Create RH SPARS assessments'),
  ('rhspars.view', 'View RH SPARS assessments'),
  ('rhspars.edit', 'Edit RH SPARS assessments'),
  ('rhspars.delete', 'Delete RH SPARS assessments')
) AS v(code, description)
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.code = v.code);

-- Role: Proficiency Score Tool only
INSERT INTO roles (name, description)
SELECT 'Proficiency Score Assessor', 'Access FP proficiency assessments only'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'Proficiency Score Assessor');

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

-- Role: RH SPARS Tool only
INSERT INTO roles (name, description)
SELECT 'RH SPARS Assessor', 'Access RH SPARS supervision assessments only'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'RH SPARS Assessor');

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

-- Role: Both tools
INSERT INTO roles (name, description)
SELECT 'Integrated RH Tools Assessor', 'Access both Proficiency Score and RH SPARS tools'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'Integrated RH Tools Assessor');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Integrated RH Tools Assessor'
  AND p.code IN (
    'tools.proficiency.access', 'tools.rh_spars.access',
    'assessments.create', 'assessments.view',
    'rhspars.create', 'rhspars.view',
    'dashboard.view', 'reports.view', 'reports.export',
    'health_workers.view', 'health_workers.create', 'health_workers.edit',
    'admin.areas.view', 'facilities.view'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- Grant both tool access permissions to any existing admin-like roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE (LOWER(r.name) LIKE '%admin%' OR LOWER(r.name) LIKE '%administrator%')
  AND p.code IN (
    'tools.proficiency.access', 'tools.rh_spars.access',
    'rhspars.create', 'rhspars.view', 'rhspars.edit', 'rhspars.delete'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

COMMIT;
