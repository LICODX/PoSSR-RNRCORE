# PoSSR RNRCORE (Layer 1 Blockchain)

[![Go Report Card](https://goreportcard.com/badge/github.com/LICODX/PoSSR-RNRCORE)](https://goreportcard.com/report/github.com/LICODX/PoSSR-RNRCORE)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Consensus: PoRS](https://img.shields.io/badge/Consensus-PoRS-green)](WHITEPAPER.md)

**PoSSR (Proof of Repeated Sorting)** is a revolutionary Layer 1 blockchain protocol designed to solve the "Winner Takes All" problem inherent in Proof of Work (PoW) and Proof of Stake (PoS). It utilizes a **Time-Memory Trade-off** algorithm where consensus is achieved through computationally intensive sorting tasks rather than random hashing, ensuring a fairer distribution of rewards.

## 🚀 Key Features

*   **Fair Consensus**: The PoRS algorithm prevents hardware monopoly. 10x hash power does not guarantee 10x rewards.
*   **High Resilience**: Audited against **Replay Attacks**, **Sybil Attacks**, and **DoS/Spam** (verified with 100GB load tests).
*   **P2P GossipSub**: Robust peer-to-peer networking layer for fast block propagation.
*   **Built-in Wallet & Explorer**: Comes with a GUI Wallet and Local Dashboard Explorer out of the box.

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE

# Build the Node
go build -o rnr-node.exe ./cmd/rnr-node

# Build the Genesis Wallet (Optional)
go build -o genesis-wallet.exe ./cmd/genesis-wallet
```

## 🛠️ Usage

### Running a Full Node
```bash
./rnr-node.exe
```
*   **P2P Port**: `9900`
*   **Dashboard**: `http://localhost:8080` (View blocks/stats)

### Running the Wallet
The node includes a built-in GUI wallet. Launch the node, and the wallet interface will initialize (if configured) or use the CLI tools.

## 🛡️ Security Audit & Stress Tests

This project has undergone rigorous Red Team auditing.
*   [Security Audit Report](docs/audit/security_audit_report.md)
*   [Massive 100-Node Simulation](docs/audit/massive_simulation_report.md)
*   [Extreme 100GB Stress Test](docs/audit/extreme_stress_report.md)

## 📄 Documentation

*   **[Whitepaper](WHITEPAPER.md)** - Technical Deep Dive into PoRS.
*   **[Testnet Guide](TESTNET_MANUAL.md)** - How to run local simulations.
*   **[API Reference](docs/api.md)** - JSON-RPC and P2P Specs.

## 🤝 Contributing

Contributions are welcome! Please read `CONTRIBUTING.md` before submitting a Pull Request.

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.

---
*Built with ❤️ by the RNRCORE Team.*
