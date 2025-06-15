# Vigil (VGL)

**Fair Launch, Secure and Decentralized**

Vigil is a groundbreaking cryptocurrency that merges the fair, grassroots distribution of GPU mining with the ironclad security and on-chain governance of a hybrid PoS network. It is built for and governed by its community.

[![Build Status](https://github.com/Vigil-Labs/vgl/workflows/Build%20Status/badge.svg)](https://github.com/Vigil-Labs/vgl/actions?query=workflow%3A%22Build+Status%22)

## 🚀 Project Overview

- **Ticker Symbol**: VGL
- **Base Technology**: Fork of Decred (DCR)
- **Proof-of-Work Algorithm**: KawPoW (ASIC-Resistant, GPU-focused)
- **Consensus**: Hybrid PoW/PoS
- **Block Time**: ~2.5 minutes
- **Total Supply**: ~21,000,000 VGL

## 🔧 Technical Specifications

### Tokenomics
- **Initial Block Reward**: 20 VGL
- **Emission Decay**: Smooth 1% reduction every 6,144 blocks (~10.7 days)
- **Block Reward Split (50/40/10 Model)**:
  - 50% (10 VGL) to PoW Miners (KawPoW)
  - 40% (8 VGL) to PoS Stakers (Ticket Holders)
  - 10% (2 VGL) to Vigil Treasury (Vigiliteia)

### Key Features
- **KawPoW Mining**: ASIC-resistant algorithm ensuring fair distribution
- **Hybrid Security**: Combined PoW mining and PoS staking for maximum security
- **On-Chain Governance**: Vigiliteia system for community-driven decisions
- **Self-Funding Treasury**: Sustainable development through treasury allocation

## 📁 Repository Structure

```
vigil/
├── node/           # Core blockchain node implementation
├── wallet/         # Wallet software
├── pool/           # Mining pool software
├── explorer/       # Block explorer
├── docs/           # Documentation
└── README.md       # This file
```

## 🛠️ Building and Running

### Prerequisites
- Go 1.19 or later
- Git

### Building the Node
```bash
cd node
go mod tidy
go build .
```

### Building the Wallet
```bash
cd wallet
go mod tidy
go build .
```

## 🏗️ Development Phases

### Phase 1: Foundation & Pre-Production (Months 0-2)
- [x] Identity & Brand Package
- [x] Technical Scaffolding & Feasibility
- [x] KawPoW Integration Proof-of-Concept
- [ ] Genesis Contributor Program

### Phase 2: Core Development & Alpha (Months 3-6)
- [x] Full KawPoW Integration
- [x] VGL Tokenomics Implementation
- [ ] Software Suite Rebranding
- [ ] Vigil Staking Dashboard

### Phase 3: Public Unveiling & Testnet (Months 7-9)
- [ ] Go-to-Market & Community Activation
- [ ] Public Testnet Launch
- [ ] Stress Gauntlet Competition

### Phase 4: Mainnet Launch & Growth (Month 10+)
- [ ] Security Audit
- [ ] Official Mining Pool Development
- [ ] Mainnet Launch Sequence
- [ ] Vigiliteia Incubator Initiative

## 🔐 Security

Vigil employs a hybrid consensus mechanism combining:
- **PoW Security**: KawPoW algorithm prevents ASIC centralization
- **PoS Security**: Ticket-based staking system adds additional security layers
- **Governance Security**: On-chain voting prevents hostile takeovers

## 🌐 Community

- **Website**: vigil.network (coming soon)
- **Discord**: [Join our community](https://discord.gg/vigil) (coming soon)
- **Twitter**: [@VigilCrypto](https://twitter.com/vigilcrypto) (coming soon)

## 📄 Documentation

For detailed technical documentation, please refer to:
- [Whitepaper](docs/whitepaper.pdf) (coming soon)
- [API Documentation](docs/api.md) (coming soon)
- [Mining Guide](docs/mining.md) (coming soon)
- [Staking Guide](docs/staking.md) (coming soon)

## 🤝 Contributing

We welcome contributions from the community! Please read our contributing guidelines and join our Genesis Contributor Program for early involvement opportunities.

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📜 License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This software is currently in development. Use at your own risk. Always verify transactions and keep your private keys secure.

---

**Vigil - Building the future of decentralized, community-governed cryptocurrency.**



