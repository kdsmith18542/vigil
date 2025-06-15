# Vigil Testnet Launch Kit

Welcome to the Vigil Testnet! This guide will help you get started with testing the Vigil blockchain network based on the specifications in the vglbp.txt file.

## Overview

Vigil is a Vigil fork implementing the KawPoW mining algorithm for GPU mining. The testnet launched on January 1, 2025, with the following key features:

- **Algorithm**: KawPoW (GPU-friendly, ASIC-resistant)
- **Block Reward**: 20 VGL per block (testnet)
- **Fair Launch**: No premine, community-driven development
- **Testnet Genesis**: "Vigil Testnet 2025 - Fair Launch"

## Quick Start

### Prerequisites

1. **Go 1.23 or 1.24** installed
2. **Git** for cloning repositories
3. **Windows 10+** (this guide is for Windows)

### Building the Software

```powershell
# Build vgld (node)
cd vgld
go install .

# Build vglwallet (wallet)
cd ..\vglwallet
go install .

# Build vglctl (command-line tool)
cd ..\vglctl
go install .
```

### Starting the Testnet

1. **Run the startup script**:
   ```cmd
   start_testnet.bat
   ```

2. **Create a testnet wallet**:
   ```cmd
   vglwallet --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --create
   ```

3. **Start the wallet**:
   ```cmd
   vglwallet --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf"
   ```

## Network Configuration

### Testnet Parameters

- **Network Magic**: `0x48e7a065` (testnet)
- **Default P2P Port**: 19108
- **Default RPC Port**: 19109
- **Address Prefixes**:
  - P2PKH: `Vg` (testnet)
  - P2SH: `Vg` (testnet)
  - Private Keys: `Pt` (testnet)

### Configuration Files

The startup script creates configuration files in `%USERPROFILE%\vigil_testnet\`:

- `vgld.conf` - Node configuration
- `vglwallet.conf` - Wallet configuration

## Common Commands

### Node Operations

```cmd
# Check node status
vglctl --configfile="%USERPROFILE%\vigil_testnet\vgld.conf" getinfo

# Get current block height
vglctl --configfile="%USERPROFILE%\vigil_testnet\vgld.conf" getblockcount

# Get network hash rate
vglctl --configfile="%USERPROFILE%\vigil_testnet\vgld.conf" getnetworkhashps
```

### Wallet Operations

```cmd
# Get wallet balance
vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet getbalance

# Get new address
vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet getnewaddress

# Send transaction
vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet sendtoaddress <address> <amount>
```

## GPU Mining with KawPoW

### Mining Setup

1. **Get a mining address**:
   ```cmd
   vglctl --configfile="%USERPROFILE%\vigil_testnet\vglwallet.conf" --wallet getnewaddress
   ```

2. **Configure mining** (when GPU miners are available):
   - **Pool URL**: `stratum+tcp://127.0.0.1:19108`
   - **Algorithm**: KawPoW
   - **Wallet Address**: Use address from step 1

### Mining Pool Development

As per vglbp.txt Phase 4 requirements:
- Fork Vigil's `VGLpool`
- Modify for KawPoW algorithm
- Implement GPU mining support
- Create web interface for miners

## Testnet Features & Testing

### Phase 3: Public Unveiling & Testnet

- ✅ **Testnet Launch**: January 1, 2025
- 🔄 **Stress Gauntlet Competition**: Community stress testing
- 🔄 **Web-based Testnet Faucet**: Get testnet VGL coins
- 🔄 **Easy-to-use Installers**: Simplified setup tools

### Testing Scenarios

1. **Basic Transactions**:
   - Send/receive VGL coins
   - Multi-signature transactions
   - Time-locked transactions

2. **Mining Tests**:
   - Solo mining
   - Pool mining (when available)
   - Difficulty adjustments

3. **Network Stress Tests**:
   - High transaction volume
   - Large block propagation
   - Network partitioning recovery

## Troubleshooting

### Common Issues

1. **"vgld not found in PATH"**:
   - Ensure Go is installed and `$GOPATH/bin` is in your PATH
   - Run `go install .` in the vgld directory

2. **Connection refused**:
   - Check if vgld is running
   - Verify configuration files
   - Check firewall settings

3. **Wallet sync issues**:
   - Ensure vgld is fully synced first
   - Check wallet configuration
   - Restart wallet if needed

### Log Files

Logs are stored in `%USERPROFILE%\vigil_testnet\logs\`:
- `vgld.log` - Node logs
- `vglwallet.log` - Wallet logs

## Development & Contribution

### Code Structure

- `vgld/vigil/` - Vigil-specific configurations
- `vgld/vigil/chaincfg/` - Network parameters
- `vgld/vigil/mining/` - KawPoW mining implementation

### Key Files Modified for Vigil

- `chaincfg/params.go` - Network parameters
- `chaincfg/genesis.go` - Genesis block configuration
- `mining/gpuminer/` - GPU mining implementation

## Community & Support

### Getting Help

- **Discord**: Join the Vigil community Discord
- **Telegram**: Vigil official Telegram group
- **GitHub**: Report issues and contribute code

### Testnet Faucet

Once available, the web-based testnet faucet will provide:
- Free testnet VGL coins
- Rate limiting to prevent abuse
- Simple web interface

## Roadmap

### Phase 4: Mainnet Launch & Growth

- **Pre-launch Readiness**: Final testing and audits
- **Official Mining Pool**: Production-ready pool software
- **Mainnet Launch**: Full network deployment
- **Post-launch Growth**: Community expansion and governance

## Technical Specifications

### Consensus Parameters

- **Block Time**: ~5 minutes (inherited from Vigil)
- **Difficulty Adjustment**: ASERT algorithm
- **Proof of Work**: KawPoW
- **Proof of Stake**: Vigil's hybrid system

### Address Formats

| Type | Mainnet Prefix | Testnet Prefix |
|------|----------------|----------------|
| P2PKH | Vg | Vg |
| P2SH | Vg | Vg |
| Private Key | Pm | Pt |

---

**Note**: This is testnet software. Do not use real funds. Testnet coins have no monetary value.

For the latest updates and documentation, check the project repository and community channels.
