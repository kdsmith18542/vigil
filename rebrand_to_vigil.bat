@echo off
REM Vigil Rebranding Script - Windows Batch Launcher
REM This batch file launches the PowerShell rebranding script

echo ========================================
echo Vigil Rebranding Script
echo ========================================
echo.

REM Check if PowerShell is available
powershell -Command "Write-Host 'PowerShell is available'" >nul 2>&1
if errorlevel 1 (
    echo ERROR: PowerShell is not available or not in PATH
    echo Please install PowerShell or run the .ps1 script directly
    pause
    exit /b 1
)

echo Choose an option:
echo 1. Run DRY RUN (preview changes without modifying files)
echo 2. Run FULL REBRANDING (modify files)
echo 3. Run with VERBOSE output
echo 4. Exit
echo.
set /p choice="Enter your choice (1-4): "

if "%choice%"=="1" (
    echo.
    echo Running DRY RUN - No files will be modified
    echo.
    powershell -ExecutionPolicy Bypass -File "%~dp0rebrand_to_vigil.ps1" -DryRun -Verbose
) else if "%choice%"=="2" (
    echo.
    echo WARNING: This will modify files in the repository!
    set /p confirm="Are you sure you want to continue? (y/N): "
    if /i "%confirm%"=="y" (
        echo.
        echo Running FULL REBRANDING...
        echo.
        powershell -ExecutionPolicy Bypass -File "%~dp0rebrand_to_vigil.ps1"
    ) else (
        echo Operation cancelled.
    )
) else if "%choice%"=="3" (
    echo.
    echo Choose verbose mode:
    echo 1. Dry run with verbose output
    echo 2. Full rebranding with verbose output
    set /p verboseChoice="Enter choice (1-2): "
    
    if "%verboseChoice%"=="1" (
        echo.
        echo Running DRY RUN with VERBOSE output
        echo.
        powershell -ExecutionPolicy Bypass -File "%~dp0rebrand_to_vigil.ps1" -DryRun -Verbose
    ) else if "%verboseChoice%"=="2" (
        echo.
        echo WARNING: This will modify files with verbose output!
        set /p confirm="Are you sure? (y/N): "
        if /i "%confirm%"=="y" (
            echo.
            echo Running FULL REBRANDING with VERBOSE output...
            echo.
            powershell -ExecutionPolicy Bypass -File "%~dp0rebrand_to_vigil.ps1" -Verbose
        ) else (
            echo Operation cancelled.
        )
    ) else (
        echo Invalid choice.
    )
) else if "%choice%"=="4" (
    echo Exiting...
    exit /b 0
) else (
    echo Invalid choice. Please run the script again.
)

echo.
echo Press any key to exit...
pause >nul