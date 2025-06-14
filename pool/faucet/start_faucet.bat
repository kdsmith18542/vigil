@echo off
REM Vigil Testnet Faucet Server Startup Script

echo ========================================
echo Vigil Testnet Faucet Server
echo ========================================
echo.

REM Check if Python is installed
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Python not found in PATH
    echo Please install Python 3.8+ and add it to your PATH
    pause
    exit /b 1
)

echo Python found!
echo.

REM Check if pip is available
pip --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: pip not found
    echo Please ensure pip is installed with Python
    pause
    exit /b 1
)

echo pip found!
echo.

REM Check if virtual environment exists
if not exist "venv" (
    echo Creating Python virtual environment...
    python -m venv venv
    if %errorlevel% neq 0 (
        echo ERROR: Failed to create virtual environment
        pause
        exit /b 1
    )
    echo Virtual environment created!
echo.
)

REM Activate virtual environment
echo Activating virtual environment...
call venv\Scripts\activate.bat
if %errorlevel% neq 0 (
    echo ERROR: Failed to activate virtual environment
    pause
    exit /b 1
)

echo Virtual environment activated!
echo.

REM Install dependencies
echo Installing Python dependencies...
pip install -r requirements.txt
if %errorlevel% neq 0 (
    echo ERROR: Failed to install dependencies
    echo Please check your internet connection and try again
    pause
    exit /b 1
)

echo Dependencies installed!
echo.

REM Check if testnet is configured
set TESTNET_CONFIG=%USERPROFILE%\vigil_testnet\vglwallet.conf
if not exist "%TESTNET_CONFIG%" (
    echo WARNING: Testnet not configured!
    echo Please run start_testnet.bat first to initialize the testnet.
    echo.
    echo The faucet server will start but may not function properly
    echo until the testnet is running and wallet is created.
    echo.
    pause
)

REM Check if wallet is running
echo Checking if vglwallet is running...
tasklist /FI "IMAGENAME eq vglwallet.exe" 2>NUL | find /I /N "vglwallet.exe">NUL
if %errorlevel% neq 0 (
    echo WARNING: vglwallet is not running!
    echo The faucet requires a running vglwallet instance.
    echo Please start your testnet wallet before using the faucet.
    echo.
)

echo ========================================
echo Starting Faucet Server
echo ========================================
echo.
echo Faucet Configuration:
echo - Amount per request: 10 VGL
echo - Cooldown period: 24 hours
echo - Rate limit: 5 requests per IP per day
echo - Server URL: http://localhost:5000
echo.
echo Starting server...
echo.
echo Press Ctrl+C to stop the server
echo.

REM Start the faucet server
python faucet_server.py

echo.
echo Faucet server stopped.
pause
