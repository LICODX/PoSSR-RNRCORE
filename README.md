# PoSSR-RNRCORE (Proof of Sequential Sorting Race)

> **⚠️ EXPERIMENTAL RESEARCH PROJECT - NOT PRODUCTION READY**
> 
> This is an **experimental blockchain research project** exploring alternative consensus mechanisms.
> Currently in **ALPHA** stage with **LOCAL DEVELOPMENT ONLY**.
> 
> **DO NOT USE WITH REAL ASSETS OR IN PRODUCTION ENVIRONMENTS.**

## What This Project Actually Is

PoSSR-RNRCORE is a research implementation of a hybrid consensus mechanism that combines:

1. **Proof-of-Work (PoW)** - Basic hash-based mining for spam prevention
2. **VRF-like Randomness** - Using miner signatures on PoW hash for unpredictable entropy
3. **Sorting Verification** - 7 sorting algorithms (QuickSort, MergeSort, HeapSort, RadixSort, TimSort, IntroSort, ShellSort)

The goal is to explore whether sorting can be used as part of blockchain consensus, rather than being a production-ready blockchain.

## What Actually Works

✅ **Core Features (Implemented & Tested)**:
- Basic PoW mining with difficulty targets
- Ed25519 signature-based VRF seed generation
- Parallel shard processing (10 shards using goroutines)
- In-place sorting algorithms (memory optimized)
- O(N) linear validation (no re-sorting)
- LibP2P GossipSub networking with topic-based sharding
- BadgerDB for state storage
- Basic transaction signing and verification
- Simple mempool management

⚠️ **Limitations (NOT Implemented)**:
- **No Byzantine Fault Tolerance** - This is NOT a proven-secure consensus
- **No Public Network** - "Mainnet" config points to localhost only
- **No Smart Contract Runtime** - WASM claims are not implemented
- **No Cross-Shard Communication** - Sharding is message-level only, not state-level
- **No Economic Security** - No incentive mechanism or game theory
- **No External Audit** - Code has never been professionally audited
- **No Production Deployment** - Never run in adversarial environment

## Quick Start (Local Development Only)

### Prerequisites
- Go 1.21+
- ~4GB RAM
- Windows/Linux/macOS

### Build and Run

```bash
# Clone repository
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE

# Build
go build -o rnr-node ./cmd/rnr-node

# Run (creates local test network)
./rnr-node
```

**Note**: This will start a single node that mines blocks locally. It is NOT connecting to any public network.

### Run Simulations

```bash
# Test with multiple local nodes (stress test)
go run simulation/mainnet_stress_test_main.go

# Test distributed sharding logic
go run simulation/distributed_sharding_main.go
```

## Architecture (What's Actually Implemented)

### Consensus Flow
1. **Mining Phase**: Node runs PoW to find hash below difficulty target
2. **Signature Phase**: Miner signs PoW hash with private key (creates VRF seed)
3. **Sorting Phase**: Mempool is sharded into 10 parts, sorted in parallel using randomly selected algorithm
4. **Validation Phase**: Other nodes verify PoW, signature, and sorting order in O(N) time

### File Structure
```
PoSSR-RNRCORE/
├── cmd/rnr-node/           # Node binary entry point
├── internal/
│   ├── blockchain/         # Block validation, chain management
│   ├── consensus/          # PoW + Sorting algorithms  
│   ├── mempool/            # Transaction pool + sharding
│   ├── p2p/                # LibP2P networking
│   └── state/              # State management (accounts, tokens)
├── pkg/
│   ├── types/              # Core data structures
│   ├── utils/              # Crypto utilities
│   └── wallet/             # Ed25519 key management
└── simulation/             # Test scripts
```

## Performance Characteristics

| Metric | Value | Note |
|--------|-------|------|
| Block Time | 60 seconds | Target (varies with difficulty) |
| Block Size | 1GB max | Theoretical (NOT tested on network) |
| TPS | Unknown | No real-world testing |
| Network | Local Only | No public peers |
| Consensus Security | **UNPROVEN** | Research only |

## Known Issues & Technical Debt

1. **Security**: No Byzantine fault tolerance proof
2. **Consensus**: Sorting doesn't contribute to security, only used for data ordering
3. **Scalability**: 1GB blocks are impractical for network propagation
4. **State**: Simple key-value store, no Merkle Patricia Trie
5. **P2P**: Basic GossipSub, no sophisticated peer management
6. **Testing**: Simulations only, no real adversarial testing

## Why This Exists

This project is a **research experiment** to explore:
- Can sorting algorithms be used in consensus? (Answer: Partially, but doesn't solve Byzantine issues)
- How does parallel shard processing perform? (Answer: Works in simulation, untested in real network)
- Is signature-based VRF sufficient? (Answer: More testing needed)

## Contributing

This is a **research project**, not a product. Contributions welcome for:
- Academic analysis of the consensus mechanism
- Performance benchmarking
- Code quality improvements  
- Documentation of findings

**NOT looking for**:
- Marketing materials
- Production deployment help
- Investment/tokenomics

## License

MIT License - See [LICENSE](LICENSE) file.

## Disclaimer

**THIS SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND.**

This is experimental research code. It has NOT been audited, NOT been tested in adversarial conditions, and is NOT suitable for any production use. The consensus mechanism is NOT proven to be Byzantine fault tolerant.

Do NOT use this code to handle real value or in any security-critical context.

---

**Developer**: LICODX Team  
**Status**: Experimental Research (Alpha)  
**Last Updated**: January 2026
