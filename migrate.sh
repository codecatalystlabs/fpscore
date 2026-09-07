#!/usr/bin/env bash
# Apply SQL migrations in order against the fpscore database (Ubuntu/Linux/macOS)
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

echo "=========================================="
echo "Applying migrations to ${DB_NAME} @ ${DB_HOST}"
echo "=========================================="

FILES=(
  migration-add-health-workers.sql
  update-fp-tool-2026-04-06.sql
  role-cleanup-and-facility-hierarchy-role.sql
  schema-rhspars.sql
  seed-rhspars-questions.sql
  seed-tool-roles.sql
  migration-rhspars-scoring.sql
)

for f in "${FILES[@]}"; do
  if [[ -f "$f" ]]; then
    echo
    echo "--- Applying $f ---"
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -f "$f"
  else
    echo "Skipping missing file: $f"
  fi
done

echo
echo "All available migrations applied successfully."
