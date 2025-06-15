# Vigil Testnet Complete Setup Guide

This guide provides step-by-step instructions to get the Vigil testnet running based on the requirements in `vglbp.txt` Phase 3: Public Unveiling & Testnet.

## 🚀 Quick Start (TL;DR)

1. **Build the software**: `cd vgld && go install . && cd ../vglwallet && go install . && cd ../vglctl && go install .`
2. **Start testnet**: `start_testnet.bat`
3. **Create wallet**: Follow the prompts
4. **Start faucet**: `cd faucet && start_faucet.bat`
5. **Access faucet**: Open http://localhost:5000

## 📋 Prerequisites

### Required Software
- **Go 1.23 or 1.24** - [Download from golang.org](https://golang.org/dl/)
- **Python 3.8+** - [Download from python.org](https://python.org/downloads/)
- **Git** - [Download from git-scm.com](https://git-scm.com/downloads)
- **Windows 10+** with PowerShell

### Hardware Requirements
- **CPU**: 2+ cores recommended
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 10GB free space for blockchain data
- **Network**: Stable internet connection

## 🔧 Installation Steps

### Step 1: Verify Prerequisites

```powershell
# Check Go installation
go version
# Should show: go version go1.23.x or go1.24.x

# Check Python installation
python --version
# Should show: Python 3.8.x or higher

# Check Git installation
git --version
# Should show: git version x.x.x
```

### Step 2: Build Vigil Software

```powershell
# Navigate to project directory
cd c:\Users\Keith\vgl

# Build vgld (blockchain node)
cd vgld
go install .
cd ..

# Build vglwallet (wallet software)
cd vglwallet
go install .
cd ..

# Build vglctl (command-line interface)
cd vglctl
go install .
cd ..
```

### Step 3: Initialize Testnet

```cmd
# Run the testnet startup script
start_testnet.bat
```

This script will:
- ✅ Check for required binaries
- ✅ Create testnet configuration files
- ✅ Set up directory structure
- ✅ Start the vgld node in testnet mode

### Step 4: Create Testnet Wallet

After the node starts, create a wallet:

```cmd
# Create a new testnet wallet
vglwallet --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --create
```

Follow the prompts to:
1. Set a wallet passphrase
2. Optionally set a public passphrase
3. Write down your seed phrase (IMPORTANT!)

### Step 5: Start Wallet

```cmd
# Start the wallet (in a new terminal)
vglwallet --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf"
```

### Step 6: Setup Testnet Faucet

```cmd
# Navigate to faucet directory
cd faucet

# Start the faucet server
start_faucet.bat
```

The faucet will be available at: http://localhost:5000

## 🛠️ Using the Testnet

### Basic Operations

Use the PowerShell tools for easy management:

```powershell
# Check status
.\testnet_tools.ps1 status

# Check balance
.\testnet_tools.ps1 balance

# Generate new address
.\testnet_tools.ps1 newaddress

# Send coins
.\testnet_tools.ps1 send VgTestAddress123 10.5

# View connected peers
.\testnet_tools.ps1 peers

# Check sync status
.\testnet_tools.ps1 sync
```

### Manual Commands

Alternatively, use `vglctl` directly:

```cmd
# Node commands
vglctl --configfile="%USERPROFILE%\vigil_testnet\vgld.conf" getinfo
vglctl --configfile="%USERPROFILE%\vigil_testnet\vgld.conf" getblockcount
vglctl --configfile="%USERPROFILE%\vigil_testnet\vgld.conf" getnetworkhashps

# Wallet commands
vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet getbalance
vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet getnewaddress
vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet listtransactions
```

## 🎯 Testnet Features

### Network Parameters
- **Algorithm**: KawPoW (GPU-friendly)
- **Block Reward**: 20 VGL per block
- **Block Time**: ~5 minutes
- **Launch Date**: January 1, 2025
- **Genesis Message**: "Vigil Testnet 2025 - Fair Launch"

### Faucet Features
- **Amount**: 10 VGL per request
- **Cooldown**: 24 hours between requests
- **Rate Limiting**: 5 requests per IP per day
- **Web Interface**: User-friendly HTML interface
- **API**: RESTful API for integration

### Address Formats
- **Testnet addresses**: Start with `Vg`
- **Private keys**: Start with `Pt`
- **Example**: `VgTestnetAddress123456789abcdef`

### Miner Configuration

To mine on the Vigil testnet, you will need a KawPoW-compatible GPU miner. While specific miner configurations can vary, here's a general outline:

1.  **Choose a KawPoW Miner**: Select a miner software that supports the KawPoW algorithm (e.g., `teamredminer`, `nbminer`, `gminer`).
2.  **Configure the Miner**: Point your miner to your local `vgld` testnet node's RPC interface or a public testnet mining pool (if available).

    Example (using a hypothetical miner and local `vgld` RPC):
    ```
    miner.exe -a kawpow -o stratum+tcp://127.0.0.1:XXXX -u <YOUR_TESTNET_WALLET_ADDRESS> -p x
    ```
    *Replace `XXXX` with the appropriate RPC port for mining (usually 9109 for `vgld`'s mining RPC).* 
    *Replace `<YOUR_TESTNET_WALLET_ADDRESS>` with an address generated from your Vigil testnet wallet.*

3.  **Start Mining**: Run the configured miner. It should connect to your `vgld` node and start submitting shares.

**Note**: Detailed miner-specific configurations are outside the scope of this general setup guide. Refer to your chosen miner's documentation for precise instructions.

## 🧪 Testing Scenarios

### Phase 3 Requirements (from vglbp.txt)

1. **✅ Testnet Launch**: January 1, 2025
2. **✅ Web-based Testnet Faucet**: Available at localhost:5000
3. **✅ Easy-to-use Installers**: Batch scripts provided
4. **🔄 Stress Gauntlet Competition**: Community stress testing
5. **🔄 Testnet Explorer**: A live Testnet Explorer URL (e.g., `http://testnet-explorer.vigil.network`).

**Note**: Setting up a full Testnet Explorer (like a forked `vgldata` instance) is beyond the scope of this guide. You would typically deploy this on a public server.

### Stress Gauntlet Competition

The "Stress Gauntlet" is an incentivized event designed to battle-test the network under load. While the full competition infrastructure is not part of this setup guide, here's how you can participate in stress testing:

-   **Highest Hashrate Challenge**: Run your KawPoW miners at maximum capacity to contribute hashrate to the testnet. Monitor your `vgld` node's logs for accepted blocks and shares.
-   **Bug Bounty Blitz**: Actively use the testnet, perform various transactions, and try to identify any unexpected behavior or bugs. Report any issues through the designated channels (e.g., GitHub issues, Discord).
-   **Governance Grand Prix**: If a testnet Vigiliteia instance is available, practice creating and voting on proposals to test the governance system.

**Note**: Details on specific rewards and official participation for the Stress Gauntlet Competition would be announced separately by the Vigil project team.

### Recommended Tests

1. **Basic Functionality**:
   ```powershell
   # Test wallet creation and address generation
   .\testnet_tools.ps1 newaddress
   
   # Test faucet request
   # Visit http://localhost:5000 and request coins
   
   # Test transaction sending
   .\testnet_tools.ps1 send <address> 1.0
   ```

2. **Network Stress Testing**:
   ```powershell
   # Run stress test
   .\testnet_tools.ps1 stress
   
   # Monitor network performance
   .\testnet_tools.ps1 status
   ```

3. **Mining Preparation**:
   ```cmd
   # Get mining address
   vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet getnewaddress
   
   # Note: GPU miners for KawPoW will connect to the node
   # Mining pool development is planned for Phase 4
   ```

## 🔍 Troubleshooting

### Common Issues

#### "vgld not found in PATH"
**Solution**: Ensure Go is properly installed and `$GOPATH/bin` is in your PATH
```powershell
# Add Go bin to PATH (PowerShell)
$env:PATH += ";$env:GOPATH\bin"

# Or add permanently via System Properties > Environment Variables
```

#### "Connection refused"
**Solution**: Check if vgld is running
```cmd
# Check if vgld is running
tasklist | findstr vgld

# If not running, restart with:
start_testnet.bat
```

#### "Wallet sync issues"
**Solution**: Ensure vgld is fully synced first
```powershell
# Check sync status
.\testnet_tools.ps1 sync

# Wait for full sync before starting wallet
```

#### "Faucet not working"
**Solution**: Check wallet and server status
```cmd
# Ensure wallet is running
tasklist | findstr vglwallet

# Check faucet logs in terminal
# Restart faucet if needed
```

### Log Files

Logs are stored in `%USERPROFILE%\vigil_testnet\logs\`:
- `vgld.log` - Node operation logs
- `vglwallet.log` - Wallet operation logs

```powershell
# View recent logs
.\testnet_tools.ps1 logs
```

## 🌐 Network Information

### Ports
- **P2P Network**: 19108
- **RPC Server**: 19109
- **Wallet RPC**: 19110
- **Faucet Web**: 5000

### Configuration Files
- **Node Config**: `%USERPROFILE%\vigil_testnet\vgld.conf`
- **Wallet Config**: `%USERPROFILE%\vigil_testnet\vglwallet.conf`
- **Data Directory**: `%USERPROFILE%\vigil_testnet\`

## 🚀 Phase 4 Preparation

### Mining Pool Development (Upcoming)
As per vglbp.txt Phase 4 requirements:
- Fork Vigil's `VGLpool`
- Modify for KawPoW algorithm
- Implement GPU mining support
- Create web interface for miners

### Mainnet Preparation
- Final security audits
- Performance optimization
- Community feedback integration
- Production deployment scripts

## 📞 Support & Community

### Getting Help
- **Documentation**: Check `TESTNET_README.md`
- **Discord**: Join the Vigil community Discord
- **Telegram**: Vigil official Telegram group
- **GitHub**: Report issues and contribute code

### Contributing
- Test the network and report bugs
- Participate in stress testing
- Provide feedback on user experience
- Help with documentation improvements

## ⚠️ Important Notes

1. **Testnet Only**: This is testnet software. Testnet coins have no monetary value.
2. **Security**: Never use real funds or mainnet private keys on testnet.
3. **Data Loss**: Testnet data may be reset during development.
4. **Performance**: Testnet may have different performance characteristics than mainnet.

## 📈 Success Metrics

For Phase 3 completion, we aim for:
- ✅ Stable testnet operation
- ✅ Functional web faucet
- 🔄 Community participation in stress testing
- 🔄 Successful "Stress Gauntlet" competition
- 🔄 Positive community feedback

---

**Ready to test Vigil?** Start with `start_testnet.bat` and join the community!

For the latest updates, check the project repository and community channels.
