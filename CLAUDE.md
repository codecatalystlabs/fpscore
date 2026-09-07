# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Family Planning Score Tool — a digital assessment system for scoring family-planning clinical procedures (counselling, OCPs, implants, IUDs, vasectomy, etc.). Go backend using the **Fiber v2** web framework with a **PostgreSQL** database, serving a vanilla-JS/Bootstrap frontend from `public/`. There is no build step or framework for the frontend; HTML/JS files are served statically.

Note: the README describes an older, smaller version of the app (just assessments). The codebase has since grown to include JWT auth, full RBAC (roles/permissions), admin-area scoping, health-worker management, reports (PDF/XLS), and audit logging. Trust the code over the README/QUICKSTART when they disagree (e.g. the README says run on port 3000 and `go run main.go handlers.go`; the real entry point is `./cmd` and the default port is 5000).

## Commands

```bash
# Run (from project root — NOT from inside cmd/)
go run ./cmd          # or: make run

# Build
go build -o bin/fpscore.exe ./cmd    # or: make build

# Dependencies
go mod download

# Test — there are currently no _test.go files in the repo.
# When adding tests: go test ./...   (single test: go test ./handlers -run TestName)
```

Configuration is via environment variables (see `config/config.go`):
- `DATABASE_URL` — defaults to `postgres://postgres:pwaiswa@localhost/fpscore?sslmode=disable`
- `PORT` — defaults to `5000`

### Database setup / reset
```bash
psql -d fpscore -f schema.sql           # schema
psql -d fpscore -f seed-questions.sql   # assessment questions + thematic areas
psql -d fpscore -f seed-data.sql        # sample geographic hierarchy
./reset-and-reseed-all.sh               # or reset-and-reseed-all.bat — wipe + reseed assessment data
```
SQL migrations are loose `.sql` files in the repo root applied manually with `psql` (e.g. `migration-add-health-workers.sql`, `update-fp-tool-2026-04-06.sql`, `role-cleanup-and-facility-hierarchy-role.sql`). There is no migration framework — apply them in order by date/intent and update `schema.sql` to match.

## Architecture

The app is a single Fiber process. `cmd/main.go` is the composition root: it loads config, connects the DB, installs middleware, and **declares every route inline** — this file is the authoritative map of the API surface and the permission required for each endpoint. Read it first to find which handler serves a given path.

Packages (flat, no internal/ layering):
- `cmd/` — entry point + route table
- `config/` — env-var config loading
- `database/` — global `database.DB *sql.DB` singleton; all handlers use raw `database/sql` queries (no ORM)
- `handlers/` — all HTTP handlers, grouped by domain (`auth`, `users`, `roles`, `permissions`, `access_control`, `facilities`, `hierarchy`, `health_workers`, `reports`, `events`, `navigation`, and `handlers.go` for assessments). Auth & permission middleware also live here.
- `middleware/` — `api_logger.go` request/response audit middleware (separate package from handler-level middleware)
- `models/` — plain structs (`models.go`, `auth_models.go`)
- `helpers/` — `logger.go` business-action audit logging
- `public/` — static frontend (one HTML + matching JS file per page)

### Authentication & authorization (central to almost every change)

Three layered checks, applied in `cmd/main.go`:

1. **`handlers.AuthMiddleware`** (`handlers/auth.go`) — validates the JWT `Bearer` token and sets `c.Locals("userID")`, `"userEmail"`, `"userName"`. Applied to the whole `/api` group except the public routes (`/api/auth/login`, `/api/auth/bootstrap`, `/api/events/log`).
2. **`handlers.CheckPermission("code")` / `CheckAnyPermission(...)`** (`handlers/access_control.go`) — per-route RBAC. Resolves user → roles → permissions. **Admins bypass all permission checks** via `IsAdmin()`, which returns true if the user has any role matching `%admin%` OR holds *every* permission in the `permissions` table.
3. **Admin-area scoping** — beyond having a permission, non-admin users only see data within their assigned geographic areas. `GetUserAdminAreas(userID)` (`handlers/auth.go`) returns the region/district/subcounty/facility IDs from `user_admin_areas`. Handlers (e.g. assessment/facility/health-worker lists) must filter query results by these IDs; **a non-admin with no admin areas assigned sees nothing**. When adding any list/read endpoint, replicate this scoping pattern or you will leak cross-area data.

Permissions are not free-form: the canonical set lives in `PredefinedPermissions` in `handlers/permissions.go`. Adding a new permission means adding it there and wiring it into the relevant route(s) in `cmd/main.go` (and usually `handlers/navigation.go`). `POST /api/permissions/initialize` upserts this list into the DB.

`jwtSecret` is currently a hardcoded constant in `handlers/auth.go` (`"change-this-secret"`) — treat as a known issue, not a pattern to copy.

### Frontend permission model

`GET /api/navigation` (`handlers/navigation.go`) returns the menu items the current user may see, computed from their permissions — the frontend renders nav from this, it does not hardcode it. Keep `navigation.go` in sync with the routes/permissions in `cmd/main.go` when adding a new admin page.

### Audit logging (two independent systems, both writing to the `events` table)

- **`middleware.APILogger()`** — wraps all `/api/*` calls, recording method/url/status/duration/IP plus a generated human-readable description into `events.data` (JSONB). Masks sensitive request fields. Skips `/api/events*`.
- **`helpers.Log*` functions** (`LogCreate`, `LogUpdate`, `LogLogin`, `LogMove`, …) — explicit business-action logging called from inside handlers. Writes are **async (`go func()`)** and best-effort — failures are printed, not returned. Use these for meaningful domain events; don't rely on them for transactional guarantees.

`GET /api/events*` exposes logs (requires `logs.view`). The frontend `event-logger.js` can also POST client-side events to `/api/events/log` (intentionally unauthenticated, to capture pre-login events).

### Assessment scoring (the domain core)

Scoring lives in the assessment-creation handler in `handlers/handlers.go`. Each question has a `score_weight`: **10 = critical (bold), 5 = important (asterisk), 2 = normal**. On submit, per response:
- `"Yes"` → earns the weight; counts toward both achieved and total-possible.
- `"No"` → earns 0; counts toward total-possible only.
- `"NA"` → ignored entirely (not counted in possible or achieved).

Final = `achieved / totalPossible * 100`. Performance level: **Proficient (>90%), Competent (70–89%), Not Acceptable (<70%)**. Scores are computed inside a DB transaction and stored denormalized on the `assessments` row plus per-area in `thematic_area_scores`. An assessment is tied to a `health_worker` (and uses the worker's *current* facility, since workers can be moved between facilities).

### Data model hierarchy

Geographic: `regions → districts → subcounties → facilities`. Health workers belong to a facility. Assessments belong to a health worker + assessment type. Assessment content: `assessment_types → thematic_areas → questions`, with responses in `assessment_responses`. RBAC: `users ↔ user_roles ↔ roles ↔ role_permissions ↔ permissions`, and `user_admin_areas` for geographic scoping.
