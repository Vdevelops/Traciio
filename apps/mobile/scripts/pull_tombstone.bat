@echo off
REM Pull tombstone after native crash (run after app has crashed once)
REM Requires: adb in PATH, device/emulator connected

set OUTPUT=%~dp0..\tombstone_00
echo Checking tombstone files...
adb shell ls -l /data/tombstones/ 2>nul
echo.
echo Pulling latest tombstone to %OUTPUT%
adb pull /data/tombstones/tombstone_00 "%OUTPUT%" 2>nul
if %ERRORLEVEL% neq 0 (
  echo Pull failed. Trying logcat backup...
  adb logcat -d > "%~dp0..\logcat_crash.txt"
  echo Saved logcat to logcat_crash.txt
) else (
  echo Done. Open tombstone_00 and search for "com.example.mobile" and "signal"
)
REM Jalankan dengan argumen apa pun untuk skip pause (mis. pull_tombstone.bat no-pause)
if "%1"=="" pause
