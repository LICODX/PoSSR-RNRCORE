# 🌍 PoSSR Public Testnet Manual

> **Welcome to the RNR-CORE Public Testnet!**
> This manual guides you through running the 25-Node Adversarial Stability Test, our final verification step before Mainnet launch.

---

## 🧪 Network Scenario: "The Byzantine Arena"

We are simulating a **25-node network** on your local machine to prove the robustness of the PoSSR consensus.

### 🎭 Cast of Characters

| Role | Count | Description |
|------|-------|-------------|
| **Honest Miners** | 18 | Follow protocol, mine blocks, validate correctly. (Nodes 1-18) |
| **Malicious Actors** | 7 | Attempt to break the network (Double Spends, Spam, Forks). (Nodes 19-25) |

**Total Nodes:** 25
**Byzantine Ratio:** 28% (Safe within the 33% BFT threshold)

---

## 🚀 How to Run the Test

### Step 1: Start the Network
Double-click `RUN_25_NODES.bat` in the root directory.

OR run via terminal:
```powershell
.\RUN_25_NODES.bat
```

### Step 2: Watch the Chaos
The script will spawn 25 terminal windows (one for each node) and a transaction automation engine.

**What happens next?**
1. **Genesis Phase:** Node 1 initializes the blockchain.
2. **Discovery:** Nodes 2-25 discover Node 1 via P2P.
3. **Mining Race:** All nodes compete to sort data.
4. **Attacks Begin:** Nodes 19-25 will start attempting invalid actions (logging errors in their respective windows).
5. **Stability:** The Honest Majority (18 nodes) should reject the bad blocks and keep the chain moving.

---

## 📊 Monitoring

### Web Dashboard
Open your browser to view the status of honest nodes:
- **Node 1 (Genesis):** [http://localhost:9101](http://localhost:9101)
- **Node 10 (Miner):** [http://localhost:9110](http://localhost:9110)
- **Node 18 (Miner):** [http://localhost:9118](http://localhost:9118)

### Log Files
All output is redirected to `nodeX.log` files in the root directory.

**To check sync status (PowerShell):**
```powershell
# See if blocks are being accepted
Select-String "Block Accepted" node1.log -Tail 10
```

**To see malicious attempts:**
```powershell
# Check for rejected blocks
Select-String "Block Rejected" node*.log
```

---

## 🛑 Stopping the Test
Close the main terminal window running the batch script, or press `Ctrl+C`. You may need to manually close the spawned node windows if they persist.

To clean up data for a fresh run:
```powershell
.\CLEANUP.bat
```

---

## 📝 Success Criteria
The test is considered successful if:
1. **Chain Growth:** Block height increases consistently (1 block/minute).
2. **Consensus:** All honest nodes agree on the same Block Hash at the same Height.
3. **Resilience:** The network does *not* crash despite the 7 malicious nodes.
