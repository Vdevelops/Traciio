@echo off
REM ============================================================
REM  K6 Load Testing Suite — CRM Healthcare
REM  Run all test scripts sequentially with cooldown periods
REM ============================================================

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0scripts"
set "REPORT_DIR=%~dp0reports"

REM Use BASE_URL from environment if set, otherwise fallback to default
if "%BASE_URL%"=="" (
    set "BASE_URL=http://localhost:8080"
)

set "TIMESTAMP=%DATE:~-4%%DATE:~3,2%%DATE:~0,2%_%TIME:~0,2%%TIME:~3,2%"
set "TIMESTAMP=%TIMESTAMP: =0%"

REM Create reports directory if not exists
if not exist "%REPORT_DIR%" mkdir "%REPORT_DIR%"

echo.
echo  ============================================================
echo   CRM Healthcare - K6 Load Testing Suite
echo   Target: %BASE_URL%
echo   Started: %DATE% %TIME%
echo  ============================================================
echo.

REM Check K6 is installed
where k6 >nul 2>nul
if errorlevel 1 (
    echo  [ERROR] K6 is not installed!
    echo  Install: winget install grafana.k6
    echo  Or: choco install k6
    echo  Or: https://grafana.com/docs/k6/latest/set-up/install-k6/
    pause
    exit /b 1
)
echo  [OK] K6 found: & k6 version
echo.

REM ============================================================
REM  Parse arguments: all, smoke, health, auth, load, stress, spike, soak
REM ============================================================

if "%~1"=="" goto :run_all
if "%~1"=="all" goto :run_all
if "%~1"=="smoke" goto :run_smoke
if "%~1"=="health" goto :run_health
if "%~1"=="auth" goto :run_auth
if "%~1"=="load" goto :run_load
if "%~1"=="stress" goto :run_stress
if "%~1"=="spike" goto :run_spike
if "%~1"=="soak" goto :run_soak
if "%~1"=="highload" goto :run_highload
if "%~1"=="monitoring" goto :run_monitoring
if "%~1"=="quick" goto :run_quick

echo  [ERROR] Unknown test: %~1
echo  Usage: run_all_tests.bat [all^|smoke^|health^|auth^|load^|stress^|spike^|soak^|highload^|monitoring^|quick]
echo  - all       : Run all tests sequentially (default)
echo  - quick     : Run smoke + health + auth only (~6 min)
echo  - highload  : Run high load test dengan 1000-5000 VUs (~35 min)
echo  - monitoring: Run monitoring endpoints test (~5 min)
exit /b 1

:run_all
echo  [1/7] Running Smoke Test...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\01_smoke_test.js"
if errorlevel 1 (
    echo  [WARN] Smoke test had failures - but continuing...
)
echo  Cooldown 15s...
timeout /t 15 /nobreak >nul

echo.
echo  [2/7] Running Health Check Test...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\02_health_check.js"
echo  Cooldown 15s...
timeout /t 15 /nobreak >nul

echo.
echo  [3/7] Running Auth Flow Test...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\03_auth_flow.js"
echo  Cooldown 15s...
timeout /t 15 /nobreak >nul

echo.
echo  [4/7] Running Load Test (~16 min)...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\04_load_test.js"
echo  Cooldown 30s...
timeout /t 30 /nobreak >nul

echo.
echo  [5/7] Running Stress Test (~20 min)...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\05_stress_test.js"
echo  Cooldown 30s...
timeout /t 30 /nobreak >nul

echo.
echo  [6/7] Running Spike Test (~5 min)...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\06_spike_test.js"
echo  Cooldown 30s...
timeout /t 30 /nobreak >nul

echo.
echo  [7/7] Running Soak Test (~34 min)...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\07_soak_test.js"
goto :done

:run_smoke
echo  Running Smoke Test...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\01_smoke_test.js"
goto :done

:run_health
echo  Running Health Check Test...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\02_health_check.js"
goto :done

:run_auth
echo  Running Auth Flow Test...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\03_auth_flow.js"
goto :done

:run_load
echo  Running Load Test (~16 min)...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\04_load_test.js"
goto :done

:run_stress
echo  Running Stress Test (~20 min)...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\05_stress_test.js"
goto :done

:run_spike
echo  Running Spike Test (~5 min)...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\06_spike_test.js"
goto :done

:run_soak
echo  Running Soak Test (~34 min)...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\07_soak_test.js"
goto :done

:run_highload
echo  Running High Load Test (~35 min)...
echo  WARNING: This test uses 1000-5000 VUs. Ensure server has sufficient resources.
echo  Set MAX_VUS environment variable untuk customize (default: 1000)
k6 run -e BASE_URL=%BASE_URL% -e MAX_VUS=%MAX_VUS% "%SCRIPT_DIR%\08_high_load_test.js"
goto :done

:run_monitoring
echo  Running Monitoring Test (~5 min)...
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\09_monitoring_test.js"
goto :done

:run_quick
echo  Running Quick Suite (smoke + health + auth, ~6 min)...
echo  -------------------------------------------------------
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\01_smoke_test.js"
timeout /t 10 /nobreak >nul
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\02_health_check.js"
timeout /t 10 /nobreak >nul
k6 run -e BASE_URL=%BASE_URL% "%SCRIPT_DIR%\03_auth_flow.js"
goto :done

:done
echo.
echo  ============================================================
echo   All tests complete!
echo   Reports saved to: %REPORT_DIR%\
echo   Finished: %DATE% %TIME%
echo  ============================================================
echo.
echo  JSON reports:
dir /b "%REPORT_DIR%\*.json" 2>nul || echo   (no reports found)
echo.
pause
