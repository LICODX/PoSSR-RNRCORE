# 🌌 PoSSR: Proof of Sequential Sorting Race - RNR Core

> **The First Deterministic Consensus Blockchain Based on Sorting Algorithms**

[![Go Report Card](https://goreportcard.com/badge/github.com/LICODX/PoSSR-RNRCORE)](https://goreportcard.com/report/github.com/LICODX/PoSSR-RNRCORE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-windows%20%7C%20linux-lightgrey)]()

---

## 🚀 Overview

**RNR Core** is a next-generation Layer-1 blockchain implementing the novel **Proof of Sequential Sorting Race (PoSSR)** consensus. Instead of energy-wasteful hashing (PoW) or stake-centralized validation (PoS), PoSSR uses cryptographic randomness (VRF) to generate strict sorting challenges. Nodes race to sort data using optimal algorithms (QuickSort, MergeSort, RadixSort, etc.), rewarding computational efficiency and algorithmic skill.

### ✨ Key Features

- **Consensus**: Proof of Sequential Sorting Race (PoSSR)
- **Engine**: 100% Go (Golang)
- **Smart Contracts**: WASM Runtime (Rust/C++) with comprehensive security
- **Scalability**: Parallel Sharding (256 Shards)
- **Database**: BadgerDB (High Performance KV Store)
- **Network**: LibP2P with GossipSub

---

## 📚 Documentation

Detailed documentation has been consolidated into the [`docs/`](./docs/) directory.

### 🌟 Start Here
- **[RNR Revolution (Whitepaper)](./docs/RNR_Revolution_Whitepaper.md)**: 📄 The complete explanation of the RNR revolution.
- **[Real Network Setup](./docs/REAL_NETWORK_SETUP.md)**: 🌐 Connect to the Mainnet Genesis Node.
- **[Adversarial Simulation](./simulation/adversarial_net_main.go)**: ⚔️ Code for 20-node attack simulation.

### 🛠️ Developer Guides
- **[Technical Analysis](./docs/Analisis_Teknis.md)**: Deep dive into current metrics.
- **[Installation & Mining](./docs/MINING.md)**: How to set up a node.
- **[Smart Contracts](./docs/SMART_CONTRACTS.md)**: Writing WASM contracts.
- **[Dashboard Manual](./docs/DASHBOARD_V2.2.md)**: Using the new Explorer & Wallet.

---

## ⚡ Quick Start: Join the Mainnet

### 1. Connect to Genesis Node
To join the live network and sync with the Genesis Node:

```bash
# 1. Build the Node
go build -o rnr-node.exe ./cmd/rnr-node

# 2. Run (Auto-connects to seed nodes in config/mainnet.yaml)
./rnr-node.exe
```

### 2. Run Simulations (Standalone)
You can run adversarial simulations without connecting to the network to verify security:

```bash
# Run 20-Node Adversarial Simulation (13 Malicious vs 7 Honest)
go run simulation/adversarial_net_main.go

# Run Internal Security Audit (Replay/DoS Tests)
go run simulation/internal_audit_main.go
```

### 3. Build from Source
```bash
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE
go build -o rnr-node.exe ./cmd/rnr-node
```

---

## 🛡️ Security & Performance

- **Block Time**: 60 Seconds (Mainnet)
- **Max Block Size**: 1GB (Theoretical Cap)
- **Protection**: Circuit Breakers, Execution Timeouts (5s), Memory Limits (64MB)
- **Audit**: [Self-Audit Report](./docs/security_audit_report.md)

---

## 🤝 Contribution

Contributions are welcome! Please check the `docs/` folder for architectural details before submitting PRs.

---

**Built with ❤️ by the LICODX Team**
