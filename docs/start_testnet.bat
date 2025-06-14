@echo off
REM Vigil Testnet Startup Script
REM Based on vglbp.txt requirements for Phase 3: Public Unveiling & Testnet

echo ========================================
echo Vigil Testnet Launch Kit
echo ========================================
echo.

REM Check if required binaries exist
echo Checking for required binaries...
where vgld >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: vgld not found in PATH
    echo Please build and install vgld first:
    echo   cd vgld
    echo   go install .
    pause
    exit /b 1
)

where vglwallet >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: vglwallet not found in PATH
    echo Please build and install vglwallet first:
    echo   cd vglwallet
    echo   go install .
    pause
    exit /b 1
)

where vglctl >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: vglctl not found in PATH
    echo Please build and install vglctl first:
    echo   cd vglctl
    echo   go install .
    pause
    exit /b 1
)

echo All required binaries found!
echo.

REM Create testnet data directory
set TESTNET_DIR=%USERPROFILE%\vigil_testnet
if not exist "%TESTNET_DIR%" (
    echo Creating testnet directory: %TESTNET_DIR%
    mkdir "%TESTNET_DIR%"
)

REM Create vgld config for testnet
echo Creating vgld testnet configuration...
echo # Vigil Testnet Configuration > "%TESTNET_DIR%\vgld.conf"
echo testnet=1 >> "%TESTNET_DIR%\vgld.conf"
echo rpcuser=vigiluser >> "%TESTNET_DIR%\vgld.conf"
echo rpcpass=vigilpass123 >> "%TESTNET_DIR%\vgld.conf"
echo rpclisten=127.0.0.1:19109 >> "%TESTNET_DIR%\vgld.conf"
echo listen=0.0.0.0:19108 >> "%TESTNET_DIR%\vgld.conf"
echo datadir=%TESTNET_DIR%\vgld_data >> "%TESTNET_DIR%\vgld.conf"
echo logdir=%TESTNET_DIR%\logs >> "%TESTNET_DIR%\vgld.conf"
echo debuglevel=info >> "%TESTNET_DIR%\vgld.conf"
echo # Enable mining for testnet >> "%TESTNET_DIR%\vgld.conf"
echo generate=false >> "%TESTNET_DIR%\vgld.conf"
echo # KawPoW mining will be enabled when miners connect >> "%TESTNET_DIR%\vgld.conf"

REM Create vglwallet config for testnet
echo Creating vglwallet testnet configuration...
echo # Vigil Testnet Wallet Configuration > "%TESTNET_DIR%\vglwallet.conf"
echo testnet=1 >> "%TESTNET_DIR%\vglwallet.conf"
echo username=vigiluser >> "%TESTNET_DIR%\vglwallet.conf"
echo password=vigilpass123 >> "%TESTNET_DIR%\vglwallet.conf"
echo rpcconnect=127.0.0.1:19109 >> "%TESTNET_DIR%\vglwallet.conf"
echo rpclisten=127.0.0.1:19110 >> "%TESTNET_DIR%\vglwallet.conf"
echo appdata=%TESTNET_DIR%\wallet_data >> "%TESTNET_DIR%\vglwallet.conf"
echo logdir=%TESTNET_DIR%\logs >> "%TESTNET_DIR%\vglwallet.conf"
echo debuglevel=info >> "%TESTNET_DIR%\vglwallet.conf"

REM Create directories
if not exist "%TESTNET_DIR%\vgld_data" mkdir "%TESTNET_DIR%\vgld_data"
if not exist "%TESTNET_DIR%\wallet_data" mkdir "%TESTNET_DIR%\wallet_data"
if not exist "%TESTNET_DIR%\logs" mkdir "%TESTNET_DIR%\logs"

echo.
echo ========================================
echo Starting Vigil Testnet Node
echo ========================================
echo.
echo Testnet Configuration:
echo - Network: Vigil Testnet (January 1, 2025)
echo - Algorithm: KawPoW (GPU Mining)
echo - Initial Reward: 20 VGL per block
echo - RPC Port: 19109
echo - P2P Port: 19108
echo - Data Directory: %TESTNET_DIR%
echo.
echo Starting vgld in testnet mode...
echo.

REM Start vgld with testnet configuration
start "Vigil Testnet Node" vgld --configfile="%TESTNET_DIR%\vgld.conf"

REM Wait a moment for vgld to start
timeout /t 5 /nobreak >nul

echo.
echo ========================================
echo Testnet Node Started Successfully!
echo ========================================
echo.
echo Next Steps:
echo.
echo 1. Create a testnet wallet:
echo    vglwallet --configfile="%TESTNET_DIR%\vglwallet.conf" --create
echo.
echo 2. Start the wallet:
echo    vglwallet --configfile="%TESTNET_DIR%\vglwallet.conf"
echo.
echo 3. Check node status:
echo    vglctl --configfile="%TESTNET_DIR%\vgld.conf" getinfo
echo.
echo 4. Get a mining address:
echo    vglctl --configfile="%TESTNET_DIR%\vglwallet.conf" --wallet getnewaddress
echo.
echo 5. Start GPU mining (when available):
echo    Connect your KawPoW miner to: stratum+tcp://127.0.0.1:19108
echo.
echo Testnet Features:
echo - Fair Launch: No premine, community-driven
echo - GPU Mining: KawPoW algorithm for decentralization
echo - Stress Testing: Participate in the "Stress Gauntlet" competition
echo - Web Faucet: Get testnet VGL coins for testing
echo.
echo For support and community: Join the Vigil Discord/Telegram
echo.
echo Press any key to continue...
pause >nul
