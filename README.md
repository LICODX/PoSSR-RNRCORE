# RnR Core: Proof of Sorting Race (PoSSR) Protocol

> **Status: Pre-Alpha (Research & Development)**
>
> ⚠️ **Warning**: This protocol is currently in active development. Features described below are implemented in the `main` branch but require specific flags (e.g., `--bft-mode`) to enable.

| Metric | Genesis Target (Current) | Future Target (Mainnet) |
| :--- | :--- | :--- |
| **Consensus** | PoSSR Leader Election + Committee BFT | Same |
| **Block Size** | 10 MB | 1 GB (Sharded) |
| **Target Hardware** | Consumer PC (8GB RAM) | High-End Workstation |
| **Development** | Core Logic Validation | Scalability Optimization |
| **TPS** | ~6,000 (Tested) | ~35,000 (Target) |

---

## 📚 Technical Documentation

Core references for developers and researchers:

- **[Technical Whitepaper v2.0](docs/TECHNICAL_WHITEPAPER_v2.md)**: Architectural specification, reconciling the high-throughput vision with current network constraints.
- **[Economic & Incentive Model](docs/INCENTIVES.md)**: Details on the hybrid reward mechanism, proportional shard distribution, and slashing conditions.
- **[Security Architecture](docs/SECURITY.md)**: Comprehensive security model including BFT guarantees, attack vectors, and mitigations.
- **[Academic Analysis](docs/ACADEMIC_PREPRINT.md)**: Original research on the properties of sorting-based consensus.
- **[Implementation Status](docs/IMPLEMENTATION_STATUS.md)**: Audit of currently running features vs. roadmap.

---

## 🏗️ Architecture

RnR Core implements a novel **Hybrid Consensus** mechanism:

1.  **Proof-of-Sequential-Sorting-Race (PoSSR)**:
    - Replaces energy-intensive hashing with useful sorting work.
    - Determines **Leader Election** (Block Proposal).
    - Algorithms: QuickSort, MergeSort, HeapSort (Running in parallel).

2.  **Tendermint-Style BFT**:
    - Ensures **Instant Finality** (2/3+ Majority).
    - Prevents forks via **Slashing** (100% penalty for equivocation).
    - Provides safety even if leader election is gamed.

### Component Stack
```
✅ Fully Implemented:
├── PoW Module               - Spam prevention (Hashcash style)
├── VRF Module               - Ed25519 seed generation
├── Sorting Engine           - 7 Parallel Algorithms
├── P2P Network              - LibP2P GossipSub (Clustered Topics)
├── BFT Engine               - 4-Phase Voting (Propose->Commit)
└── Shard Rewards            - Proportional Distribution Logic

🚧 In Development:
├── WASM Runtime             - Smart Contract Layer
└── State Sharding           - Cross-shard Atomicity
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Git

### Installation

```bash
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE
go build -o rnr-node ./cmd/rnr-node
```

### Running the Node

**Mode A: Single Node (Dev/Test)**
Standard mode for testing sorting algorithms locally.
```bash
./rnr-node
```

**Mode B: BFT Consensus (Pre-Alpha Network)**
Enables full consensus engine, P2P voting, and slashing enforcement.
```bash
./rnr-node --bft-mode --port 3000
```

---

## 🤝 Contributing

We welcome contributions from the community. Please review our [Incentive Model](docs/INCENTIVES.md) to understand the protocol's goals effectively.

1.  Fork the repository
2.  Create your feature branch (`git checkout -b feature/amazing-feature`)
3.  Commit your changes (`git commit -m 'Add some amazing feature'`)
4.  Push to the branch (`git push origin feature/amazing-feature`)
5.  Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
