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
- **[PROJECT PRESENTATION](./PROJECT_PRESENTATION.md)**: Executive summary, architecture, and high-level overview.
- **[Public Testnet Manual](./docs/PUBLIC_TESTNET.md)**: 🌍 **Start Here for the 25-Node Simulation!**
- **[Real Network Setup](./docs/REAL_NETWORK_SETUP.md)**: 🌐 **Connect Computers via Internet (WAN)**.
- **[Whitepaper](./docs/whitepapers/PoSSR_Whitepaper.pdf)**: The theoretical foundation of PoSSR.
- **[Blueprints](./docs/whitepapers/PoSSR_Blueprint.pdf)**: Technical architecture diagrams.

### 🛠️ User Guides
- **[Installation & Mining](./docs/MINING.md)**: How to set up a node and start mining.
- **[Smart Contracts](./docs/SMART_CONTRACTS.md)**: Writing and deploying WASM contracts.
- **[Security Protections](./docs/SECURITY_PROTECTIONS.md)**: Deep dive into the VM security layer.
- **[Dashboard Manual](./docs/DASHBOARD_V2.2.md)**: Using the Material Design explorer.

### 📊 Reports & Specs
- **[Final Implementation Report](./docs/reports/FINAL_IMPLEMENTATION_REPORT.md)**: Validation of Whitepaper compliance.
- **[Hardware Test Report](./docs/whitepapers/Hardware_Test_Report.pdf)**: Performance metrics on different hardware.
- **[Sharding Specification](./docs/SHARDING.md)**
- **[API Reference](./docs/API.md)**

---

## ⚡ Quick Start (Public Testnet)

We are currently in the **Public Testnet Phase**, simulating a 25-node adversarial environment to prove BFT consensus.

### 1. Prerequisites
- **OS**: Windows (optimized for batch scripts), Linux, or macOS.
- **Go**: Version 1.20+

### 2. Run the 25-Node Test
This script spins up 18 honest nodes and 7 malicious nodes to test network resilience.

```bash
.\RUN_25_NODES.bat
```

> **Note**: This will open multiple terminal windows and a browser dashboard.
> See [docs/PUBLIC_TESTNET.md](./docs/PUBLIC_TESTNET.md) for full details.

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
