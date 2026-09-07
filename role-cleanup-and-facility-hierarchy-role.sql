-- Role cleanup + facility hierarchy role setup
-- Safe to run in Postgres.

BEGIN;

-- ------------------------------------------------------------------
-- 1) Ensure dedicated permission exists for facility hierarchy moves
-- ------------------------------------------------------------------
INSERT INTO permissions (code, description)
SELECT 'hierarchy.facilities.manage', 'Manage facility hierarchy operations (move facilities between subcounties)'
WHERE NOT EXISTS (
  SELECT 1 FROM permissions WHERE code = 'hierarchy.facilities.manage'
);

-- ------------------------------------------------------------------
-- 2) Create dedicated role for facility hierarchy operations
-- ------------------------------------------------------------------
INSERT INTO roles (name, description)
SELECT 'Facility Hierarchy Manager', 'Can move facilities across subcounties and view area/facility hierarchy.'
WHERE NOT EXISTS (
  SELECT 1 FROM roles WHERE lower(trim(name)) = lower(trim('Facility Hierarchy Manager'))
);

WITH role_row AS (
  SELECT id AS role_id FROM roles WHERE name = 'Facility Hierarchy Manager'
),
perm_rows AS (
  SELECT id, code FROM permissions
  WHERE code IN ('hierarchy.facilities.manage', 'admin.areas.view', 'facilities.view')
)
INSERT INTO role_permissions (role_id, permission_id)
SELECT role_row.role_id, perm_rows.id
FROM role_row
CROSS JOIN perm_rows
WHERE NOT EXISTS (
  SELECT 1 FROM role_permissions rp
  WHERE rp.role_id = role_row.role_id
    AND rp.permission_id = perm_rows.id
);

-- ------------------------------------------------------------------
-- 3) De-duplicate roles by normalized name (trim + lower)
--    Keeps the smallest id as canonical, re-links assignments, deletes duplicates.
-- ------------------------------------------------------------------
WITH normalized AS (
  SELECT
    id,
    name,
    lower(trim(name)) AS normalized_name
  FROM roles
),
dupe_groups AS (
  SELECT normalized_name, MIN(id) AS keep_id
  FROM normalized
  GROUP BY normalized_name
  HAVING COUNT(*) > 1
),
dupe_rows AS (
  SELECT n.id AS duplicate_id, g.keep_id, n.normalized_name
  FROM normalized n
  JOIN dupe_groups g ON n.normalized_name = g.normalized_name
  WHERE n.id <> g.keep_id
)
-- Move user-role assignments from duplicate roles to canonical role
INSERT INTO user_roles (user_id, role_id)
SELECT ur.user_id, d.keep_id
FROM user_roles ur
JOIN dupe_rows d ON ur.role_id = d.duplicate_id
WHERE NOT EXISTS (
  SELECT 1 FROM user_roles ur2
  WHERE ur2.user_id = ur.user_id
    AND ur2.role_id = d.keep_id
);

WITH normalized AS (
  SELECT
    id,
    name,
    lower(trim(name)) AS normalized_name
  FROM roles
),
dupe_groups AS (
  SELECT normalized_name, MIN(id) AS keep_id
  FROM normalized
  GROUP BY normalized_name
  HAVING COUNT(*) > 1
),
dupe_rows AS (
  SELECT n.id AS duplicate_id, g.keep_id, n.normalized_name
  FROM normalized n
  JOIN dupe_groups g ON n.normalized_name = g.normalized_name
  WHERE n.id <> g.keep_id
)
-- Move role-permission assignments from duplicate roles to canonical role
INSERT INTO role_permissions (role_id, permission_id)
SELECT d.keep_id, rp.permission_id
FROM role_permissions rp
JOIN dupe_rows d ON rp.role_id = d.duplicate_id
WHERE NOT EXISTS (
  SELECT 1 FROM role_permissions rp2
  WHERE rp2.role_id = d.keep_id
    AND rp2.permission_id = rp.permission_id
);

WITH normalized AS (
  SELECT
    id,
    name,
    lower(trim(name)) AS normalized_name
  FROM roles
),
dupe_groups AS (
  SELECT normalized_name, MIN(id) AS keep_id
  FROM normalized
  GROUP BY normalized_name
  HAVING COUNT(*) > 1
),
dupe_rows AS (
  SELECT n.id AS duplicate_id
  FROM normalized n
  JOIN dupe_groups g ON n.normalized_name = g.normalized_name
  WHERE n.id <> g.keep_id
)
DELETE FROM roles r
USING dupe_rows d
WHERE r.id = d.duplicate_id;

COMMIT;

-- Optional verification query:
-- SELECT lower(trim(name)) AS normalized_name, COUNT(*) FROM roles GROUP BY 1 HAVING COUNT(*) > 1;
