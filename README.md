# PoSSR RNR Core

This repository contains the reference implementation of the **RnR Blockchain**, featuring a hybrid **Proof-of-Sequential-Sorting-Race (PoSSR)** and **BFT** consensus mechanism.

> **Status**: Experimental / Pre-Alpha. Do not use in production with real value.

## ✅ Implemented Features

The following features are implemented in Go and verifiable via code:

1.  **Consensus Engine** (`internal/consensus`)
    - **Algorithm**: Hybrid PoSSR (Sorting Race) + Tendermint-style BFT.
    - **Finality**: Instant finality once 2/3+ validators precommit.
    - **Logic**: Random sorting algorithm selection (Quick/Merge/Radix/Heap) seeded by VRF.

2.  **Networking** (`internal/p2p`)
    - **Protocol**: LibP2P with GossipSub.
    - **Hardening**: Explicit 10MB message size limit to override default 1MB cap.
    - **Topics**: Directed gossiping for Shards vs Headers.

3.  **Storage** (`internal/storage`)
    - **Engine**: LevelDB.
    - **Schema**: Batch-writes for blocks (Shards stored individually to prevent OOM).
    - **Pruning**: Rolling window (default 100 blocks) to maintain lightweight footprint.

4.  **Security** (`internal/slashing`)
    - **Slashing**: Automated detection for Double-Signing (100% slash) and Downtime (1% slash).
    - **Cryptography**: Ed25519 for validator signatures.

## 🛠️ Build & Run

### Prerequisites
- Go 1.21+
- GCC (for LevelDB/BadgerDB CGO)

### 1. Build
```bash
go build -o bin/rnr-node.exe ./cmd/rnr-node
```

### 2. Run (Single Node - PoW Mode)
Default educational mode (Mini-PoW for identity).
```bash
./bin/rnr-node.exe --port 3000 --rpc-port 9000
```

### 3. Run (BFT Authority Mode)
Runs the node as a validator with BFT finality enabled.
```bash
# Set Genesis Mnemonic (Required for Authority)
$env:GENESIS_MNEMONIC="your mnemonic here"

# Run with --bft-mode
./bin/rnr-node.exe --bft-mode --genesis --port 3000
```

### 4. Run Stress Test (Verification)
Verify the internal logic without running a full node.
```bash
go run ./cmd/stress-test/main.go
```

## 📂 Project Structure

- **`cmd/`**: Entry points.
  - `rnr-node`: Main node daemon.
  - `stress-test`: Hardening verification script.
- **`internal/`**: Core logic (Private).
  - `consensus`: BFT engine & Sorting algorithms.
  - `p2p`: LibP2P wrapper & hardening.
  - `storage`: LevelDB manager & batching logic.
  - `slashing`: Offense detection & evidence handling.
- **`pkg/`**: Shared libraries.
  - `types`: Core data structures (Block, Tx).
  - `wallet`: Key management.

## ⚠️ Known Limits (Hardcoded)
- **Max Block Size**: 10 MB (Hard limit to prevent network DoS).
- **Max Message Size**: 10 MB (LibP2P GossipSub constant).
- **Pruning**: Retains last 25-100 blocks. 
- **Validators**: Currently optimized for static validator sets (Genesis-based).
