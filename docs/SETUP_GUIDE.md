# RnR Core: Setup & Installation Guide

Welcome to the **RnR Core** node setup guide. This document will walk you through the process of building the node from source, configuring it for local development, and joining the Pre-Alpha network.

---

## 📋 Prerequisites

Before you begin, ensure your environment meets the following requirements:

### Hardware (Minimum)
- **CPU**: 2+ Cores
- **RAM**: 4GB
- **Storage**: 10GB free space
- **Network**: Stable internet connection

### Software
- **Operating System**: Linux (Ubuntu 20.04+), macOS, or Windows
- **Go**: Version 1.21 or higher ([Download Go](https://go.dev/dl/))
- **Git**: Latest version

---

## 🛠️ Step 1: detailed Installation

### 1. Clone the Repository
```bash
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE
```

### 2. Build the Binary
We use standard Go toolchain for building.
```bash
# Download dependencies
go mod download

# Build the binary
go build -o ./bin/rnr-node ./cmd/rnr-node

# Verify installation
./bin/rnr-node --help
```

---

## 🚀 Step 2: Running a Local Node

### Mode A: Single Node (Dev/Test)
This is the default mode for testing sorting algorithms and mining locally.
```bash
./bin/rnr-node --port 3000 --datadir ./data/node1
```
*Output should show "PoW Mining Mode" active.*

### Mode B: BFT Consensus Mode (Advanced)
To test the new BFT consensus engine:
```bash
./bin/rnr-node --bft-mode --port 3000 --datadir ./data/validator1
```
*Output should show "BFT Consensus Mode Enabled".*

---

## 🌐 Step 3: Multi-Node Local Network

To simulate a network locally, you can run multiple instances on different ports.

### Node 1 (Genesis Validator)
```bash
./bin/rnr-node --bft-mode --port 3000 --rpc-port 9001 --datadir ./data/node1 --genesis
```

### Node 2 (Peer)
```bash
./bin/rnr-node --bft-mode --port 3001 --rpc-port 9002 --datadir ./data/node2 --peers /ip4/127.0.0.1/tcp/3000/p2p/<NODE1_PEER_ID>
```

---

## ⚙️ Configuration Flags

| Flag | Description | Default |
| :--- | :--- | :--- |
| `--port` | P2P listening port | 3000 |
| `--rpc-port` | JSON-RPC API port | 9001 |
| `--datadir` | Directory for chain data | `./data/chaindata` |
| `--bft-mode` | Enable BFT Consensus Engine | `false` |
| `--peers` | Comma-separated list of bootstrap peers | `""` |
| `--genesis` | Start as Genesis Authority | `false` |

---

## ❓ Troubleshooting

**Q: "Command not found: go"**
A: Ensure Go is in your system PATH. Try `export PATH=$PATH:/usr/local/go/bin`.

**Q: "Failed to open database: resource temporarily unavailable"**
A: You are trying to run two nodes pointing to the same `--datadir`. Use different directories for each node.

**Q: "Consensus stuck at height X"**
A: In BFT mode, consensus waits for 2/3 voting power. If you are running a single node, ensure it has 100% of the voting power (Genesis mode).
