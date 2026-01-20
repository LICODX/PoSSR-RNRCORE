# 25-Node Adversarial Network Test

## 🎯 Overview

This is a **Byzantine Fault Tolerance test** for the RNR blockchain with:
- **25 total nodes** (18 honest, 7 malicious)
- **28% Byzantine ratio** (below the 33% theoretical threshold ✅)
- **Automated transactions** simulating realistic network activity
- **Malicious behavior simulation** to test network resilience

## 🚀 Quick Start

### Prerequisites
- Built `rnr-node.exe` in the root directory
- At least 8GB RAM (25 nodes + transactions)
- Windows OS (PowerShell/CMD)

### Run Test

```bash
# Start all 25 nodes + automated transactions
RUN_25_NODES.bat
```

**What happens:**
1. Cleans previous test data
2. Starts 18 honest nodes (ports 8001-8018)
3. Starts 7 malicious nodes (ports 8019-8025)
4. Waits 30 seconds for network formation
5. Starts automated transaction system
6. Opens monitoring dashboard

### Monitor Test

**Dashboards:**
- Node 1 (Primary): http://localhost:9101
- Node 10 (Mid): http://localhost:9110
- Node 18 (Last Honest): http://localhost:9118

**Logs:**
- `node1.log` - `node25.log` in root directory
- Real-time in console windows

## 📊 Network Composition

### Honest Nodes (1-18)

Normal validators with proper behavior:
- ✅ Valid block proposals
- ✅ Correct transaction validation
- ✅ Honest consensus participation
- ✅ Proper block propagation

### Malicious Nodes (19-25)

Each node exhibits different Byzantine behavior:

| Node | Type | Attack Behavior |
|------|------|-----------------|
| **19** | Double Spend | Attempts to spend same coins twice |
| **20** | Invalid TX | Invalid signatures, wrong nonces |
| **21** | Block Spam | Excessive block proposals |
| **22** | TX Spam | Floods network with transactions |
| **23** | Selfish Mining | Withholds blocks strategically |
| **24** | Fork Creator | Attempts to fork the chain |
| **25** | Silent Attack | Accepts but doesn't propagate |

## 🔄 Automated Transactions

### Every 3 Blocks: Random Transfers

**Trigger**: Blocks 3, 6, 9, 12, 15, ...

**Action**:
- 5 random honest nodes (from 1-18)
- Each sends 2 RNR
- To random addresses discovered on blockchain

**Example Output:**
```
[Block 9] Triggering random transfers...
  ✅ Node 7 sent 2.0 RNR to a3f2...
  ✅ Node 12 sent 2.0 RNR to 5b89...
  ✅ Node 3 sent 2.0 RNR to c1d4...
  ✅ Node 15 sent 2.0 RNR to 9e27...
  ✅ Node 1 sent 2.0 RNR to 72f1...
```

### Every 25 Blocks: Token Creation

**Trigger**: Blocks 25, 50, 75, 100, ...

**Action**:
- 3 random honest nodes
- Create new RNR-20 tokens
- Various types: stablecoins, utility, governance

**Token Types**:
- USD Reward (USDR)
- EUR Reward (EURR)
- JPY Reward (JPYR)
- Game Token (GAME)
- Data Token (DATA)
- Point Token (POINT)
- DAO Tokens (DAO1, DAO2, DAO3)

**Example Output:**
```
[Block 25] Triggering token creation...
  ✅ Node 8 created token: Game Token-B25 (GAME1)
  ✅ Node 14 created token: DAO Token 1-B25 (DAO11)
  ✅ Node 5 created token: USD Reward-B25 (USDR1)
```

### Every 55 Blocks: Smart Contract Deployment

**Trigger**: Blocks 55, 110, 165

**Action**:
- 7 random honest nodes
- Deploy WASM smart contracts
- Contract types: counter (3), vesting (2), custom (2)

**Current Status**:
⚠️ Smart contract deployment not yet fully integrated into block processing. This feature shows placeholder output for now.

## 📈 Test Metrics

### Success Criteria

**Network Health**:
- ✅ Network remains operational throughout test
- ✅ Block production continues (target: 10s avg)
- ✅ No network partitions

**Byzantine Resistance**:
- ✅ >90% of honest blocks accepted
- ✅ <1% malicious blocks in main chain
- ✅ 100% double spend attempts rejected
- ✅ Invalid transactions rejected

**Performance**:
- ⚡ TPS > 50
- ⏱️ Block time ≈ 10s average
- 🔄 Chain reorgs < 3 blocks deep
- 📦 Mempool < 1,000 transactions

### What to Watch

**Good Signs** ✅:
- Steady block production
- Malicious node rejections logged
- Network peer count stable
- Dashboard shows healthy metrics

**Warning Signs** ⚠️:
- Block time > 20s consistently
- Chain reorganizations > 5 blocks
- Mempool > 2,000 transactions
- Network partitioning

**Critical Issues** 🚨:
- Block production stops
- Successful double spend
- Invalid blocks in main chain
- Network split

## 🔍 Monitoring

### During Test

**Console Windows**:
- Each node has its own window
- Watch for error messages
- Malicious nodes labeled in title

**Log Files**:
```bash
# View specific node log
type node1.log | more

# Search for errors
findstr /i "error" node*.log

# Search for malicious activity
findstr /i "attack\|malicious\|invalid" node*.log
```

### After Test

**Collect Results**:
```bash
# All logs in one file
copy /b node*.log all_nodes.log

# Count blocks produced
findstr /i "new block" all_nodes.log | find /c /v ""

# Count transactions
findstr /i "transaction" all_nodes.log | find /c /v ""

# Count rejections
findstr /i "rejected\|invalid" all_nodes.log | find /c /v ""
```

## 🎬 Test Scenarios

### Scenario 1: Baseline (Blocks 1-20)

**Objective**: Establish baseline metrics

**Actions**:
- Let network stabilize
- Only honest nodes producing blocks
- No attacks activated yet

**Expected**:
- Smooth block production
- Normal transaction processing
- Peer discovery complete

### Scenario 2: Byzantine Activation (Blocks 21-100)

**Objective**: Test resilience against attacks

**Actions**:
- Malicious nodes activate (simulated)
- Automated transfers every 3 blocks
- Token creation every 25 blocks

**Expected**:
- Network continues operating
- Malicious behavior detected
- Valid transactions still processed

### Scenario 3: High Load (Blocks 101-200)

**Objective**: Stress test under adversarial conditions

**Actions**:
- Continued attacks
- Contract deployments (every 55 blocks)
- High transaction volume

**Expected**:
- Network handles load
- Security protections active
- Circuit breaker works if needed

## 🛠️ Troubleshooting

### Nodes Won't Start

**Problem**: Ports already in use

**Solution**:
```bash
# Kill existing node processes
taskkill /F /IM rnr-node.exe

# Or restart computer
```

### Network Not Forming

**Problem**: Nodes not discovering peers

**Solution**:
- Wait longer (up to 60 seconds)
- Check if Node 1 (bootstrap) is running
- Verify firewall allows localhost connections

### High CPU/Memory Usage

**Problem**: 25 nodes too resource-intensive

**Solution**:
- Close other applications
- Reduce logging verbosity
- Run test on more powerful machine
- Or reduce to 15 nodes for testing

### Transaction Automation Not Working

**Problem**: Automated TX script errors

**Solution**:
```bash
# Check if nodes are accepting RPC
curl http://localhost:9001/api/blockchain

# Manually verify node is responsive
start http://localhost:9101
```

## 📊 Expected Results

### Network Stats (200 blocks @ 10s each)

- **Duration**: ~33 minutes
- **Total Blocks**: 200
- **Total Transfers**: ~330 (66 cycles × 5 nodes)
- **Total Tokens**: 24 (8 cycles × 3 nodes)
- **Total Contracts**: 21 (3 cycles × 7 nodes)
- **Total RNR Transferred**: ~660 RNR

### Byzantine Behavior

- **Double Spends Attempted**: ~20+
- **Invalid TXs Sent**: ~50+
- **Block Spam Attempts**: ~30+
- **TX Spam Count**: ~1000+
- **Selfish Mining Events**: ~10+
- **Fork Attempts**: ~5+

**All should be REJECTED** ✅

## 🎓 Learning Objectives

### Byzantine Fault Tolerance

**Question**: Can the network survive with 28% malicious nodes?

**Answer**: Yes! Byzantine threshold is 33%, so 28% should be tolerable.

**Observation**:
- How does network detect malicious behavior?
- How quickly are attacks mitigated?
- Does throughput degrade?

### Consensus Under Attack

**Question**: Does consensus still work with active attacks?

**Observation**:
- Do honest nodes agree on chain state?
- Are there chain reorganizations?
- How does PoSSR handle conflicts?

### Security Protections

**Question**: Do smart contract security protections work?

**Observation**:
- Are malicious contracts rejected?
- Do execution limits trigger?
- Does circuit breaker activate if needed?

## 📝 Test Report Template

After test completes, document:

### Network Health
- [ ] Block production rate:  blocks/min
- [ ] Average block time: __ seconds
- [ ] Chain reorganizations: __
- [ ] Network partitions: __

### Attack Resistance
- [ ] Double spends attempted: __
- [ ] Double spends successful: __ (should be 0!)
- [ ] Invalid blocks proposed: __
- [ ] Invalid blocks accepted: __ (should be 0!)

### Performance
- [ ] Transactions processed: __
- [ ] Average TPS: __
- [ ] Peak TPS: __
- [ ] Tokens created: __

### Anomalies
- [ ] Unexpected errors: __
- [ ] Node crashes: __
- [ ] Data corruption: __

## 🔗 Files

- `RUN_25_NODES.bat` - Main launcher
- `simulation/automated_tx/main.go` - Transaction automation
- `simulation/malicious_nodes/` - Attack simulations (TBD)
- `node1.log` - `node25.log` - Individual node logs

## ⚡ Quick Commands

```bash
# Start test
RUN_25_NODES.bat

# View primary dashboard
start http://localhost:9101

# Check all logs for errors
findstr /i "error\|fail\|panic" node*.log

# Count successful transactions
findstr /i "transaction confirmed" node*.log | find /c /v ""

# Stop all nodes
taskkill /F /IM rnr-node.exe
```

---

**Byzantine Test**: Proving RNR's resilience against adversarial conditions! 🛡️
