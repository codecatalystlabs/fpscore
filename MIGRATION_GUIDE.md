# Migration Guide: Adding Health Workers Support

This guide explains how to migrate an existing database to support health workers.

## Step 1: Run the Migration Script

Run the migration script to add health workers support to your existing database:

```bash
psql -U your_username -d your_database -f migration-add-health-workers.sql
```

This script will:
1. Create the `health_workers` table
2. Add `health_worker_id` column to `assessments` table
3. Create default health workers for each facility
4. Assign existing assessments to default health workers
5. Make `health_worker_id` NOT NULL
6. Add necessary indexes
7. Add health workers permissions

## Step 2: Initialize Permissions

The migration script automatically adds the health workers permissions. However, if you need to initialize them manually:

1. Via API (if you have admin access):
   - POST to `/api/permissions/initialize`

2. Or via SQL:
   ```sql
   INSERT INTO permissions (code, description) VALUES
       ('health_workers.view', 'View health workers'),
       ('health_workers.create', 'Create health workers'),
       ('health_workers.edit', 'Edit health workers'),
       ('health_workers.delete', 'Delete health workers')
   ON CONFLICT (code) DO NOTHING;
   ```

## Step 3: Assign Permissions to Roles

After permissions are created, assign them to appropriate roles:

```sql
-- Example: Give Admin role all health worker permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Admin'
AND p.code IN ('health_workers.view', 'health_workers.create', 'health_workers.edit', 'health_workers.delete')
ON CONFLICT DO NOTHING;
```

## Step 4: Replace Default Health Workers

After migration, you should replace the default health workers with real ones:

1. Use the Health Workers management page (`health-workers.html`)
2. Or use the seed script: `seed-health-workers.sql`

## Step 5: Restart the Application

Restart your Go application to load the new permissions and routes.

## Notes

- The migration creates default health workers with email format: `default_<facility_id>@example.com`
- All existing assessments are assigned to these default health workers
- You can delete default health workers after creating real ones and reassigning assessments
- The `delete-data.sql` script now only deletes assessment data, preserving metadata (regions, districts, facilities, etc.)

