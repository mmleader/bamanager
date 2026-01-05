@echo off
setlocal enabledelayedexpansion

echo 🚀 Starting Build for Windows...

:: 1. Setup Environment
echo 🔍 Detecting Go environment...
for /f "tokens=*" %%i in ('go env GOPATH') do set "GOPATH=%%i"
if "%GOPATH%"=="" (
    set "GOPATH=%USERPROFILE%\go"
)
set "PATH=%PATH%;%GOPATH%\bin;%USERPROFILE%\go\bin"

:: 2. Check for Wails
where wails >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Error: 'wails' command not found even after searching GOPATH\bin.
    echo Please ensure Wails is installed (go install github.com/wailsapp/wails/v2/cmd/wails@latest)
    pause
    exit /b 1
)

:: 1. Compile project
echo 📂 Compiling application...
wails build

if %errorlevel% neq 0 (
    echo ❌ Build failed.
    pause
    exit /b %errorlevel%
)

echo.
echo ✅ Build complete!
echo The executable can be found in: build\bin\bsmanager.exe
echo.
pause
