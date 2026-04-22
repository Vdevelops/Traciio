@echo off
REM Script untuk create database indexes setelah Docker containers running
REM Usage: run_indexes.bat

echo ========================================
echo Creating Database Indexes
echo ========================================
echo.

cd /d "%~dp0"

echo Waiting for API to be ready...
:WAIT_API
curl -s http://localhost:8080/health >nul 2>&1
if errorlevel 1 (
    echo   Waiting...
    timeout /t 2 /nobreak >nul
    goto WAIT_API
)
echo API is ready!
echo.

echo Creating indexes...
docker-compose exec -T postgres psql -U postgres -d crm_healthcare -f /docker-entrypoint-initdb.d/01_indexes.sql

echo.
echo ========================================
echo Done! Indexes created.
echo ========================================
pause
