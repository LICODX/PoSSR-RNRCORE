# RNR Mining Guide

## Overview
RNR uses Proof-of-SSR (Sorting-Shard-Race) consensus - a unique algorithm combining PoW with sorting race competition.

---

## How RNR Mining Works

### 1. Proof-of-Work (PoW)
- Find valid block hash meeting difficulty target
- Uses standard SHA-256 hashing
- Difficulty adjusts automatically

### 2. VRF Algorithm Selection (Post-Mining)
- After PoW success, VRF seed determines sorting algorithm
- Prevents pre-optimization attacks
- Ensures fairness

### 3. Sorting Race
- 10 fastest nodes compete
- Each sorts transactions using assigned algorithm
- Winners share block reward

---

## Mining Rewards

### Block Reward Structure
- **Initial Reward:** 100 RNR per block
- **Distribution:** Split among 10 winning nodes
- **Halving:** Every 3.5M blocks (~7% decay)
- **Total Supply:** 5 Billion RNR

### Reward Calculation
```
Coinbase Reward = BaseReward / 10 nodes = 10 RNR per winner
Transaction Fees = Sum of all TX fees in block
Total Reward = Coinbase + Fees
```

**Example:**
- Block #12345 with 100 transactions @ 1 RNR fee each
- Coinbase: 10 RNR
- Fees: 100 RNR
- **Total per winner: 110 RNR**

---

## Mining Requirements

### Hardware
**Minimum:**
- CPU: 4 cores
- RAM: 8 GB
- Storage: 100 GB SSD
- Network: 100 Mbps

**Recommended:**
- CPU: 8+ cores (higher = better)
- RAM: 16 GB
- Storage: 500 GB NVMe SSD
- Network: 1 Gbps

### Software
- RNR Node (latest version)
- Wallet with RNR for fees
- Stable internet connection

---

## Start Mining

### Step 1: Install Node
See [NODE_SETUP.md](NODE_SETUP.md) for installation instructions.

### Step 2: Create Wallet
```bash
# Generate wallet with password
./rnr-node -wallet-password "SecurePassword123"

# Note your address
# rnr1pq03gqs8zg0sgqg7zsw3u8...
```

### Step 3: Get Initial RNR
**Testnet:**
- Use faucet: https://faucet.rnr.network
- Request 100 RNR for gas fees

**Mainnet:**
- Buy from exchange
- Receive from another wallet
- Minimum: 10 RNR for fees

### Step 4: Start Mining
```bash
./rnr-node \
  -datadir ~/.rnr/data \
  -p2pport 9900 \
  -wallet-password "SecurePassword123" \
  -peers "seed1.rnr.network:9900,seed2.rnr.network:9900"
```

**Mining starts automatically!**

---

## Mining Strategies

### Solo Mining
**Pros:**
- Keep 100% of rewards
- Full control

**Cons:**
- Unpredictable income
- Need good hardware

**Best For:** Enthusiasts with strong hardware

### Pool Mining (Future)
**Pros:**
- Steady income
- Lower hardware requirements

**Cons:**
- Pool fees (typically 1-3%)
- Share rewards

**Status:** Pools planned for Q2 2026

---

## Optimizing Mining Performance

### 1. Hardware Optimization
```bash
# CPU affinity (Linux)
taskset -c 0-7 ./rnr-node ...

# Increase file descriptors
ulimit -n 65536

# Disable CPU throttling
sudo cpupower frequency-set -g performance
```

### 2. Network Optimization
```bash
# Increase network buffers
sudo sysctl -w net.core.rmem_max=134217728
sudo sysctl -w net.core.wmem_max=134217728

# Enable TCP BBR
sudo sysctl -w net.ipv4.tcp_congestion_control=bbr
```

### 3. Storage Optimization
- Use NVMe SSD (not HDD!)
- Enable TRIM for SSD longevity
- Monitor disk usage regularly

### 4. Multiple Nodes
Run multiple nodes to increase winning chances:

```bash
# Node 1
./rnr-node -p2pport 9900 -dashboard 8080 &

# Node 2 (different ports!)
./rnr-node -p2pport 9901 -dashboard 8081 &

# Node 3
./rnr-node -p2pport 9902 -dashboard 8082 &
```

---

## Monitoring Mining Activity

### Dashboard
Access: `http://localhost:8080`

**Key Metrics:**
- Blocks found: Your mining success
- Hash rate: Mining speed
- Peer count: Network connectivity
- Balance: Accumulated rewards

### Command Line
```bash
# Check balance
curl http://localhost:8080/api/wallet | jq '.balance'

# Check latest block
curl http://localhost:8080/api/stats | jq '.blockHeight'

# Check mining status
curl http://localhost:8080/api/mining
```

### Logs
```bash
# Successful block find
grep "Block #.*added to chain" ~/.rnr/logs/node.log

# Mining rewards
grep "Coinbase" ~/.rnr/logs/node.log
```

---

## Expected Mining Returns

### Hashrate vs Rewards
Assuming 1000 total network nodes:

| Your Nodes | Win Rate | Daily Blocks | Daily Rewards |
|------------|----------|--------------|---------------|
| 1 | 1% | ~8 | ~80 RNR |
| 5 | 5% | ~43 | ~430 RNR |
| 10 | 10% | ~86 | ~860 RNR |
| 50 | 50% | ~432 | ~4,320 RNR |

**Note:** Actual returns vary with:
- Network difficulty
- Your hardware performance
- Network latency
- Transaction fee volume

### Profitability Calculator
```
Daily Reward = (Your Hashrate / Network Hashrate) × Daily Blocks × Reward
Daily Cost = (Power Consumption kW × 24h × Electricity Rate) + Hardware Depreciation
Daily Profit = Daily Reward - Daily Cost
```

**Example:**
- Hardware: 8-core server @ 150W
- Electricity: $0.10/kWh
- Power cost: 150W × 24h × $0.10 = $0.36/day
- Rewards: 80 RNR/day @ $10/RNR = $800/day
- **Profit: $799.64/day**

---

## Troubleshooting Mining Issues

### Not Finding Blocks
**Symptoms:** Zero blocks found after 24 hours

**Solutions:**
1. Check peer count (need 5+)
2. Verify wallet has RNR for fees
3. Check system time (must be synced)
4. Increase CPU cores allocated

### Low Hash Rate
**Symptoms:** Dashboard shows low hash rate

**Solutions:**
1. Close other applications
2. Disable CPU power saving
3. Upgrade hardware
4. Check thermal throttling

### Orphaned Blocks
**Symptoms:** Found blocks not in main chain

**Solutions:**
1. Improve network connection
2. Reduce network latency
3. Connect to closer peers
4. Check for clock sync

---

## Mining ROI Analysis

### Equipment Costs
| Item | Cost | Lifespan |
|------|------|----------|
| 8-core Server | $1,500 | 3 years |
| 500GB NVMe | $100 | 3 years |
| Power (annual) | $130 | - |
| Internet (annual) | $600 | - |

**Total Year 1:** $2,330

### Revenue Projections (Conservative)
- Blocks/day: 8 (1% network share)
- Reward/block: 10 RNR
- Fees/block: 10 RNR
- **Daily: 160 RNR**
- **Monthly: 4,800 RNR**
- **Annual: 57,600 RNR**

**At $10/RNR:**
- Annual Revenue: $576,000
- Annual Cost: $2,330
- **Annual Profit: $573,670**
- **ROI: 24,500%**

**Note:** These are PROJECTIONS. Actual results vary significantly.

---

## Advanced Mining Topics

### Custom Mining Strategies
```go
// Implement custom block selection
func CustomBlockSelection(mempool []Transaction) []Transaction {
    // Prioritize high-fee transactions
    sort.SliceStable(mempool, func(i, j int) bool {
        return mempool[i].Fee > mempool[j].Fee
    })
    return mempool[:100] // Top 100
}
```

### Mining Analytics
Track your performance:
```sql
-- Database query for mining stats
SELECT 
    DATE(timestamp) as date,
    COUNT(*) as blocks_found,
    SUM(reward) as total_rewards
FROM mining_history
GROUP BY DATE(timestamp)
ORDER BY date DESC;
```

---

## Mining Best Practices

1. **Keep Node Updated**
   - Check for updates weekly
   - Join Discord for announcements
   - Enable auto-restart on crashes

2. **Secure Your Rewards**
   - Move rewards to cold wallet weekly
   - Never store large amounts in hot wallet
   - Use hardware wallet for long-term storage

3. **Monitor 24/7**
   - Set up alerts (email/SMS)
   - Use UPS for power outages
   - Monitor from phone app

4. **Join Community**
   - Discord: https://discord.gg/rnr
   - Reddit: r/RNRNetwork
   - Share strategies and tips

---

## Support

- **Mining Forum:** https://forum.rnr.network/mining
- **Discord #mining:** https://discord.gg/rnr
- **Email:** mining-support@rnr.network

---

**Happy Mining! 🚀⛏️**

**Last Updated:** 2026-01-18  
**Version:** 1.0.0
