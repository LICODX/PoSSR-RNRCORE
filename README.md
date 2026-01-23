# Sorting-Based Consensus Exploration (SBCE)

> **🔬 ACADEMIC RESEARCH PROJECT**
> 
> This repository contains experimental code for academic research exploring whether sorting algorithms can contribute to distributed consensus mechanisms.
> 
> **NOT A PRODUCTION BLOCKCHAIN. EDUCATIONAL TESTBED (BFT IMPLEMENTED).**

---

## 📖 **Important: Vision vs Reality**

**This README describes Phase 0 (Educational L1) - our current implementation.**

For the full story of how this project evolved from the original whitepaper vision (1GB blocks) to the current implementation (10MB blocks), and the phased roadmap to eventually achieve the full vision, please read:

**👉 [VISION_VS_REALITY.md](VISION_VS_REALITY.md)** — Reconciles whitepaper claims, technical criticism, and realistic implementation.

---

## Research Question

**Primary Question**: Can computational sorting races provide measurable security or consensus properties in distributed systems?

**Hypothesis Being Tested**: Combining Proof-of-Work (spam prevention) + VRF (randomness) + Sorting Verification (ordering) might create a partial consensus mechanism.

**Expected Result**: Likely NO - sorting alone does not solve Byzantine Generals Problem, but measuring the overhead and limitations is valuable.

## What This Code Does

This is a **proof-of-concept implementation** to measure and analyze:

1. **Performance Overhead**: How much time does parallel sorting take for various data sizes?
2. **Verification Efficiency**: Can O(N) linear checks replace O(N log N) re-sorting?
3. **Network Behavior**: How does topic-based message sharding perform in LibP2P?
4. **Algorithm Variance**: Do different sorting algorithms (Quick, Merge, Heap, etc.) create measurable randomness?

## Architecture (Educational L1 Testbed)

```
✅ Fully Implemented Components:
├── PoW Module               - Basic hash-based mining (spam prevention)
├── VRF Module               - Ed25519 signature-based seed generation
├── Sorting Engine           - 7 algorithms (parallel execution)
├── Linear Validator         - O(N) order verification
├── P2P Network              - LibP2P with GossipSub (10 shard topics + BFT topics)
├── State Store              - LevelDB key-value store
├── BFT Consensus Engine     - Tendermint-style consensus (Propose→Prevote→Precommit→Commit)
├── Finality Tracker         - Instant finality via 2/3+ votes
├── Slashing Enforcement     - Double-sign detection & automatic penalties
├── Validator Management     - Multi-validator support with shard assignment
└── Proportional Rewards     - Shard-based reward distribution

⚠️ Partially Implemented:
├── Smart Contract Runtime   - WASM VM present but disabled (circular import issue)
├── Cross-Shard Atomicity    - Message-level sharding only
└── Network Discovery        - mDNS for local, manual peering for WAN

❌ Not Yet Implemented:
├── Full State Sharding      - Only transaction/mempool sharding
├── Light Clients            - No SPV/merkle proofs
└── Advanced Cryptoeconomics - Basic slashing only, no complex incentives
```

## Running the Node

### Prerequisites
- Go 1.21+
- ~4GB RAM
- Port 3000 (P2P) and 9001 (RPC) available

### PoW Mode (Original - Single Node)
```bash
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE
go build -o rnr-node ./cmd/rnr-node
./rnr-node
```

### BFT Mode (NEW - Multi-Validator Consensus)
```bash
# Node with BFT consensus enabled
./rnr-node --bft-mode

# Custom ports
./rnr-node --bft-mode --port 3001 --rpc-port 9002
```

**Multi-Node Setup** (requires manual validator configuration - see [bft_integration_walkthrough.md](docs/bft_integration_walkthrough.md)):
```bash
# Node 1
./rnr-node --bft-mode --port 3000

# Node 2 (connect to Node 1)
./rnr-node --bft-mode --port 3001 --peer /ip4/<NODE1_IP>/tcp/3000/p2p/<PEER_ID>
```

## BFT Consensus Features (NEW)

The node now supports **Byzantine Fault Tolerant consensus** alongside the original sorting-based mechanism:

### Consensus Modes

| Mode | Command | Use Case |
|------|---------|----------|
| **PoW** | `./rnr-node` | Single-node testing, sorting research |
| **BFT** | `./rnr-node --bft-mode` | Multi-validator network with fault tolerance |

### BFT Features Implemented

1. **Tendermint-Style Consensus**
   - Four-phase voting: Propose → Prevote → Precommit → Commit
   - 2/3+ majority required at each phase
   - Automatic timeout and round progression

2. **Instant Finality**
   - Blocks finalized immediately upon 2/3+ precommits
   - No probabilistic finality (unlike Bitcoin/Ethereum PoW)
   - Checkpoint system every 100 blocks

3. **Economic Security (Slashing)**
   - Automatic detection of double-signing
   - 100% stake burned for Byzantine behavior
   - Malicious validators tombstoned (permanently removed)

4. **Fair Reward Distribution**
   - Validators rewarded based on shard processing
   - Round-robin shard assignment
   - Multiple coinbase transactions per block

### Limitation: Multi-Validator Setup

Currently, multi-validator networks require manual configuration. See the [BFT Integration Walkthrough](docs/bft_integration_walkthrough.md) for detailed setup instructions.

---

### Simulation Experiments
```bash
# Experiment 1: Sorting Performance with Large Datasets
go run simulation/mainnet_stress_test_main.go

# Experiment 2: Distributed Shard Communication Overhead
go run simulation/distributed_sharding_main.go

# Experiment 3: P2P Message Propagation
go run simulation/p2p_heavy_load_main.go
```

## Research Findings (Preliminary)

### Performance Metrics (20 Node Simulation, 1.5GB Dataset)

| Metric | Observed Value | Notes |
|--------|---------------|-------|
| Sorting Time (10 Shards) | 5-10 seconds | Parallel execution on 8-core CPU |
| Validation Time | ~1 second | O(N) linear scan |
| Memory Usage | 4.2 GB | 20 in-process nodes + mempool |
| Block Propagation | NOT TESTED | 1GB blocks are impractical for real network |

### Key Limitations Discovered

1. **Sorting Doesn't Prevent Byzantine Behavior**
   - A malicious node can easily provide correctly-sorted data but exclude/reorder critical transactions
   - Sorting only enforces ordering, not validity or safety

2. **1GB Blocks Are Impractical**
   - Network propagation would take minutes to hours
   - Forces extreme centralization (only data centers can participate)
   - Contradicts decentralization goals

3. **Topic-Based Sharding ≠ State Sharding**
   - Current implementation only routes messages by topic
   - Does NOT partition state or enable parallel execution
   - Scalability gain is minimal

4. **No Economic Security**
   - Without game theory (slashing, rewards, etc.), nodes have no incentive to behave honestly
   - Research prototype only

## Academic Context

### Related Work
- **Bitcoin PoW**: Uses hash preimage resistance, proven secure under assumptions
- **Algorand VRF**: Uses cryptographic VRF for leader selection
- **Practical BFT**: Solves consensus with 3f+1 nodes, proven Byzantine fault tolerant

### How This Differs
- **Not BFT**: Does not solve Byzantine agreement
- **Not Proven Secure**: No security proof, no game-theoretic analysis
- **Exploratory Only**: Measuring performance, not claiming security

### Potential Publication Venue
- **Workshop on Consensus Mechanisms** (if findings show interesting performance trade-offs)
- **Systems Performance Conferences** (if overhead analysis is novel)
- **Negative Results Track** (documenting why sorting-based consensus doesn't work)

## Citation

If you use this code in academic work, please cite:

```bibtex
@misc{sbce2026,
  title={Sorting-Based Consensus Exploration: Measuring Performance Overhead of Computational Races in Distributed Systems},
  author={LICODX Research Team},
  year={2026},
  note={Experimental research code, not peer-reviewed},
  url={https://github.com/LICODX/PoSSR-RNRCORE}
}
```

## Contributing

This is a **research experiment**. Contributions welcome for:
- ✅ Performance benchmarking
- ✅ Formal analysis of security properties (or lack thereof)
- ✅ Alternative experimental designs
- ✅ Documentation of failure modes

**NOT looking for**:
- ❌ Production features
- ❌ Marketing/tokenomics
- ❌ Deployment infrastructure

## Research Team

- **LICODX Team** - Experimental implementation
- **Status**: Independent research, no institutional affiliation
- **Funding**: None (unfunded research)

## License

MIT License - Code provided as-is for academic/research purposes only.

## Disclaimer

**THIS IS EXPERIMENTAL RESEARCH CODE.**

- ❌ NOT Byzantine fault tolerant
- ❌ NOT suitable for production
- ❌ NOT audited for security
- ❌ NO economic incentives
- ❌ NO warranty

This code is for **academic exploration only**. Do not use with real value or in any security-critical context.

---

**Repository**: [github.com/LICODX/PoSSR-RNRCORE](https://github.com/LICODX/PoSSR-RNRCORE)  
**Research Status**: Active Exploration  
**Last Updated**: January 2026
