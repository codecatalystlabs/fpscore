-- Apply RH SPARS schema + questions + tool roles in one go
-- Usage: psql -d fpscore -f apply-rhspars.sql

\i schema-rhspars.sql
\i seed-rhspars-questions.sql
\i seed-tool-roles.sql
