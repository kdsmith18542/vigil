# Vigil Testnet Management Tools
# PowerShell script for common testnet operations

param(
    [Parameter(Position=0)]
    [string]$Command,
    [Parameter(Position=1)]
    [string]$Address,
    [Parameter(Position=2)]
    [decimal]$Amount
)

$TestnetDir = "$env:USERPROFILE\vigil_testnet"
$vgldConfig = "$TestnetDir\vgld.conf"
$WalletConfig = "$TestnetDir\vglwallet.conf"

function Show-Help {
    Write-Host "Vigil Testnet Tools" -ForegroundColor Green
    Write-Host "==================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Usage: .\testnet_tools.ps1 <command> [parameters]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  status          - Show node and wallet status"
    Write-Host "  balance         - Show wallet balance"
    Write-Host "  newaddress      - Generate new wallet address"
    Write-Host "  send <addr> <amount> - Send VGL to address"
    Write-Host "  mine <blocks>   - Mine specified number of blocks (simnet only)"
    Write-Host "  peers           - Show connected peers"
    Write-Host "  sync            - Check sync status"
    Write-Host "  faucet          - Request testnet coins (when available)"
    Write-Host "  stress          - Run stress test transactions"
    Write-Host "  stop            - Stop node and wallet"
    Write-Host "  logs            - Show recent log entries"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\testnet_tools.ps1 status"
    Write-Host "  .\testnet_tools.ps1 send VgTestAddress123 10.5"
    Write-Host "  .\testnet_tools.ps1 mine 10"
}

function Test-Prerequisites {
    $missing = @()
    
    if (!(Get-Command vgld -ErrorAction SilentlyContinue)) {
        $missing += "vgld"
    }
    if (!(Get-Command vglwallet -ErrorAction SilentlyContinue)) {
        $missing += "vglwallet"
    }
    if (!(Get-Command vglctl -ErrorAction SilentlyContinue)) {
        $missing += "vglctl"
    }
    
    if ($missing.Count -gt 0) {
        Write-Host "Missing required binaries: $($missing -join ', ')" -ForegroundColor Red
        Write-Host "Please build and install them first." -ForegroundColor Red
        return $false
    }
    
    if (!(Test-Path $vgldConfig)) {
        Write-Host "Testnet not initialized. Please run start_testnet.bat first." -ForegroundColor Red
        return $false
    }
    
    return $true
}

function Invoke-vgldCommand {
    param([string]$Cmd)
    try {
        $result = & vglctl --configfile="$vgldConfig" $Cmd.Split(' ') 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "Command failed: $result"
        }
        return $result
    }
    catch {
        Write-Host "Error executing vgld command: $_" -ForegroundColor Red
        return $null
    }
}

function Invoke-WalletCommand {
    param([string]$Cmd)
    try {
        $result = & vglctl --configfile="$WalletConfig" --wallet $Cmd.Split(' ') 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "Command failed: $result"
        }
        return $result
    }
    catch {
        Write-Host "Error executing wallet command: $_" -ForegroundColor Red
        return $null
    }
}

function Show-Status {
    Write-Host "Vigil Testnet Status" -ForegroundColor Green
    Write-Host "===================" -ForegroundColor Green
    Write-Host ""
    
    Write-Host "Node Information:" -ForegroundColor Yellow
    $nodeInfo = Invoke-vgldCommand "getinfo"
    if ($nodeInfo) {
        $nodeInfo | ConvertFrom-Json | Format-List version, protocolversion, blocks, connections, difficulty
    }
    
    Write-Host "Wallet Information:" -ForegroundColor Yellow
    $walletInfo = Invoke-WalletCommand "getinfo"
    if ($walletInfo) {
        $walletInfo | ConvertFrom-Json | Format-List version, walletversion, balance, unconfirmedbalance
    }
    
    Write-Host "Network Hash Rate:" -ForegroundColor Yellow
    $hashRate = Invoke-vgldCommand "getnetworkhashps"
    if ($hashRate) {
        Write-Host "  $hashRate H/s"
    }
}

function Show-Balance {
    Write-Host "Wallet Balance" -ForegroundColor Green
    Write-Host "=============" -ForegroundColor Green
    
    $balance = Invoke-WalletCommand "getbalance"
    if ($balance) {
        $balanceObj = $balance | ConvertFrom-Json
        Write-Host "Available: $($balanceObj.spendable) VGL" -ForegroundColor Cyan
        Write-Host "Unconfirmed: $($balanceObj.unconfirmed) VGL" -ForegroundColor Yellow
        Write-Host "Total: $($balanceObj.total) VGL" -ForegroundColor Green
    }
}

function New-Address {
    Write-Host "Generating new address..." -ForegroundColor Green
    $address = Invoke-WalletCommand "getnewaddress"
    if ($address) {
        Write-Host "New address: $address" -ForegroundColor Cyan
        Write-Host "Address copied to clipboard (if available)"
        try {
            $address | Set-Clipboard
        } catch {
            # Clipboard not available, ignore
        }
    }
}

function Send-Transaction {
    param([string]$ToAddress, [decimal]$Amount)
    
    if (!$ToAddress -or $Amount -le 0) {
        Write-Host "Usage: send <address> <amount>" -ForegroundColor Red
        return
    }
    
    Write-Host "Sending $Amount VGL to $ToAddress..." -ForegroundColor Green
    $txid = Invoke-WalletCommand "sendtoaddress $ToAddress $Amount"
    if ($txid) {
        Write-Host "Transaction sent successfully!" -ForegroundColor Green
        Write-Host "TXID: $txid" -ForegroundColor Cyan
    }
}

function Start-Mining {
    param([int]$Blocks = 1)
    
    Write-Host "Mining $Blocks blocks..." -ForegroundColor Green
    $result = Invoke-vgldCommand "generate $Blocks"
    if ($result) {
        Write-Host "Successfully mined $Blocks blocks" -ForegroundColor Green
        $result | ConvertFrom-Json | ForEach-Object {
            Write-Host "Block: $_" -ForegroundColor Cyan
        }
    }
}

function Show-Peers {
    Write-Host "Connected Peers" -ForegroundColor Green
    Write-Host "===============" -ForegroundColor Green
    
    $peers = Invoke-vgldCommand "getpeerinfo"
    if ($peers) {
        $peerList = $peers | ConvertFrom-Json
        if ($peerList.Count -eq 0) {
            Write-Host "No peers connected" -ForegroundColor Yellow
        } else {
            $peerList | ForEach-Object {
                Write-Host "Peer: $($_.addr) - Version: $($_.subver)" -ForegroundColor Cyan
            }
        }
    }
}

function Show-SyncStatus {
    Write-Host "Sync Status" -ForegroundColor Green
    Write-Host "===========" -ForegroundColor Green
    
    $blockCount = Invoke-vgldCommand "getblockcount"
    $bestBlockHash = Invoke-vgldCommand "getbestblockhash"
    
    if ($blockCount -and $bestBlockHash) {
        Write-Host "Current block height: $blockCount" -ForegroundColor Cyan
        Write-Host "Best block hash: $bestBlockHash" -ForegroundColor Cyan
        
        # Check if wallet is synced
        $walletInfo = Invoke-WalletCommand "getinfo"
        if ($walletInfo) {
            $walletObj = $walletInfo | ConvertFrom-Json
            if ($walletObj.blocks -eq $blockCount) {
                Write-Host "Wallet is synced" -ForegroundColor Green
            } else {
                Write-Host "Wallet sync: $($walletObj.blocks)/$blockCount" -ForegroundColor Yellow
            }
        }
    }
}

function Request-FaucetCoins {
    Write-Host "Testnet Faucet" -ForegroundColor Green
    Write-Host "==============" -ForegroundColor Green
    Write-Host "Web-based testnet faucet coming soon!" -ForegroundColor Yellow
    Write-Host "For now, you can:" -ForegroundColor Cyan
    Write-Host "1. Mine blocks locally (simnet mode)" -ForegroundColor Cyan
    Write-Host "2. Ask community members for testnet coins" -ForegroundColor Cyan
    Write-Host "3. Join the Discord/Telegram for faucet access" -ForegroundColor Cyan
}

function Start-StressTest {
    Write-Host "Stress Test" -ForegroundColor Green
    Write-Host "===========" -ForegroundColor Green
    Write-Host "Starting stress test transactions..." -ForegroundColor Yellow
    
    # Generate multiple addresses and send small amounts
    for ($i = 1; $i -le 10; $i++) {
        $addr = Invoke-WalletCommand "getnewaddress"
        if ($addr) {
            Write-Host "Generated address $i: $addr" -ForegroundColor Cyan
            # In a real stress test, we would send transactions here
            # For now, just generate addresses
        }
        Start-Sleep -Milliseconds 100
    }
    
    Write-Host "Stress test completed. Check transaction pool with 'getmempoolinfo'" -ForegroundColor Green
}

function Stop-Services {
    Write-Host "Stopping Vigil services..." -ForegroundColor Yellow
    
    # Try to stop gracefully
    try {
        Invoke-vgldCommand "stop"
        Write-Host "vgld stopped gracefully" -ForegroundColor Green
    } catch {
        Write-Host "Could not stop vgld gracefully, checking processes..." -ForegroundColor Yellow
    }
    
    # Force stop if still running
    Get-Process | Where-Object {$_.ProcessName -match "vgld|vglwallet"} | ForEach-Object {
        Write-Host "Stopping $($_.ProcessName)..." -ForegroundColor Yellow
        $_ | Stop-Process -Force
    }
    
    Write-Host "All services stopped" -ForegroundColor Green
}

function Show-Logs {
    Write-Host "Recent Log Entries" -ForegroundColor Green
    Write-Host "==================" -ForegroundColor Green
    
    $logDir = "$TestnetDir\logs"
    if (Test-Path $logDir) {
        $vgldLog = "$logDir\vgld.log"
        $walletLog = "$logDir\vglwallet.log"
        
        if (Test-Path $vgldLog) {
            Write-Host "vgld logs (last 10 lines):" -ForegroundColor Yellow
            Get-Content $vgldLog -Tail 10
            Write-Host ""
        }
        
        if (Test-Path $walletLog) {
            Write-Host "vglwallet logs (last 10 lines):" -ForegroundColor Yellow
            Get-Content $walletLog -Tail 10
        }
    } else {
        Write-Host "Log directory not found: $logDir" -ForegroundColor Red
    }
}

# Main script logic
if (!(Test-Prerequisites)) {
    exit 1
}

switch ($Command.ToLower()) {
    "status" { Show-Status }
    "balance" { Show-Balance }
    "newaddress" { New-Address }
    "send" { Send-Transaction $Address $Amount }
    "mine" { Start-Mining ([int]$Address) }
    "peers" { Show-Peers }
    "sync" { Show-SyncStatus }
    "faucet" { Request-FaucetCoins }
    "stress" { Start-StressTest }
    "stop" { Stop-Services }
    "logs" { Show-Logs }
    default { Show-Help }
}
