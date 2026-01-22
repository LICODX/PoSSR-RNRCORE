# PoSSR-RNRCORE: Proof of Repeated Sorting & Sharding (Alpha)

> **⚠️ WARNING: DEVELOPMENT PREVIEW (DEVNET)**
> This project is currently in **ALPHA STAGE**. It is NOT ready for production use with real assets.
> The "Mainnet" configuration currently refers to a **local development network** (Seed Nodes = Localhost).
> Use at your own risk. Code audits are pending.

**PoSSR (Proof of Shared Sorting Rarity)** is an experimental Layer-1 Blockchain Consensus mechanism that combines:
1.  **Proof-of-Work (PoW)** as a Sybil resistance mechanism (Spam Prevention).
2.  **Verifiable Random Function (VRF)** based on **Miner Signatures** for unpredictable entropy.
3.  **Proof-of-Sorting (PoS)** for deterministic winner selection using 7 sorting algorithms.
4.  **True Sharding** via P2P Topic Splitting for horizontal scalability.

## 🚀 Corrected Status (Jan 2026)

| Feature | Status | Notes |
|---------|--------|-------|
| **Consensus** | ✅ Active | Hybrid PoW + Signed VRF + Sorting |
| **Throughput** | ⚠️ Alpha | Target: 100k TPS (Simulated), Real: Untested on Public Net |
| **Sharding** | ✅ Functional | 10 Shards w/ Distributed P2P Topics |
| **Security** | ⚠️ Audit Pending | Internal "Self-Audit" only. **NOT AUDITED BY 3RD PARTY** |
| **Network** | 🚧 Devnet | Running on Local/Private IP. Public Testnet: TBA |
| **Block Size** | ⚠️ Experimental | 1GB Limit (Requires High Bandwidth) |

## 🛠️ Technical Highlights (Latest Updates)

-   **Signed VRF**: Randomness seeded by digital signatures (`Sign(PrivKey, PoWHash)`), verifiable and unpredictable.
-   **In-Place Sorting**: Zero-Copy memory optimization for sorting algorithms to reduce GC pressure.
-   **O(N) Validation**: Linear scan validation (no re-sorting) for maximum verification speed.
-   **Gas Metering**: Computational cost tracking for Smart Contract execution.

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
