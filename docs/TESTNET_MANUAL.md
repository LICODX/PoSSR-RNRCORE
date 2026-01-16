# PoSSR RNRCORE Testnet Manual

Since PoSSR is a Layer 1 protocol, the "Testnet" can be run locally on your machine to verify consensus mechanics, mining, and security before deploying to public servers.

## 1. Prerequisites
*   OS: Windows, Linux, or macOS
*   Go: Version 1.20+
*   RAM: 4GB+ (for stress tests)

## 2. Running the Simulations (Local Testnet)

We provide three levels of network simulation:

### A. Fairness Test (Consensus Verification)
Simulates a hostile environment where 80% of nodes are malicious. Run this to prove that honest miners can still earn rewards.

```bash
go run simulation/massive_attack_main.go
# Expected Output: ~5% block wins for Honest Nodes.
```

### B. Security / Red Team Audit
Tests specific attack vectors like Replay Attacks and Packet Fuzzing.

```bash
go run simulation/red_team_main.go
# Expected Output: "🛡️ SECURE" for all vectors.
```

### C. Extreme Stress Test (100GB Load)
**WARNING: HIGH CPU USAGE.**
Simulates a massive data flood attack.

```bash
go run simulation/stress_test_extreme.go
# Expected Output: Node remains stable, rejects invalid txs.
```

## 3. Running a Private Local Node
To start a single full node that acts as its own testnet:

1.  **Start the Node:**
    ```bash
    go run cmd/rnr-node/main.go
    ```
2.  **Access Dashboard:**
    Open `http://localhost:8080/` in your browser.
3.  **Check API:**
    ```bash
    curl http://localhost:8080/api/stats
    ```

## 4. Wallet Usage
Use the built-in GUI wallet (part of the node) or the CLI tool `cmd/genesis-wallet` to generate keys for testing.

```bash
go run cmd/genesis-wallet/main.go
```
