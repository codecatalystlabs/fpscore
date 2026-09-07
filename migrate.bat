@echo off
REM Apply SQL migrations in order against the fpscore database (Windows)
REM Usage: migrate.bat
REM Optional env overrides: DB_HOST, DB_USER, DB_NAME

if "%DB_USER%"=="" set DB_USER=postgres
if "%DB_NAME%"=="" set DB_NAME=fpscore
if "%DB_HOST%"=="" set DB_HOST=localhost

echo ==========================================
echo Applying migrations to %DB_NAME% @ %DB_HOST%
echo ==========================================

set FAILED=0

for %%F in (
    migration-add-health-workers.sql
    update-fp-tool-2026-04-06.sql
    role-cleanup-and-facility-hierarchy-role.sql
    schema-rhspars.sql
    seed-rhspars-questions.sql
    seed-tool-roles.sql
    migration-rhspars-scoring.sql
) do (
    if exist "%%F" (
        echo.
        echo --- Applying %%F ---
        psql -h %DB_HOST% -U %DB_USER% -d %DB_NAME% -v ON_ERROR_STOP=1 -f "%%F"
        if errorlevel 1 (
            echo ERROR: Failed applying %%F
            set FAILED=1
            goto :done
        )
    ) else (
        echo Skipping missing file: %%F
    )
)

:done
echo.
if "%FAILED%"=="1" (
    echo Migrations failed.
    exit /b 1
)
echo All available migrations applied successfully.
exit /b 0
