#!/usr/bin/env bash
# Apply schema + migrations against the fpscore database (Ubuntu/Linux/macOS)
# Usage:
#   chmod +x migrate.sh
#   ./migrate.sh
#   make migrate
#
# Optional env overrides:
#   DB_HOST=localhost DB_USER=postgres DB_NAME=fpscore ./migrate.sh
#   PGPASSWORD=yourpassword ./migrate.sh

set -euo pipefail

cd "$(dirname "$0")"

DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-fpscore}"
DB_HOST="${DB_HOST:-localhost}"

run_sql() {
  local f="$1"
  if [[ -f "$f" ]]; then
    echo
    echo "--- Applying $f ---"
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -f "$f"
  else
    echo "Skipping missing file: $f"
  fi
}

echo "=========================================="
echo "Applying schema + migrations to ${DB_NAME} @ ${DB_HOST}"
echo "=========================================="

# 1) Base relations (must come first on a fresh database)
run_sql schema.sql

# 2) Core seed data (geography + assessment questions)
run_sql seed-data.sql
run_sql seed-questions.sql

# 3) Incremental migrations / tool updates (idempotent where possible)
FILES=(
  migration-add-health-workers.sql
  update-fp-tool-2026-04-06.sql
  role-cleanup-and-facility-hierarchy-role.sql
  schema-rhspars.sql
  seed-rhspars-questions.sql
  seed-tool-roles.sql
  migration-rhspars-scoring.sql
  seed-users.sql
)

for f in "${FILES[@]}"; do
  run_sql "$f"
done

echo
echo "All available schema/migrations applied successfully."
