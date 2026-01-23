# Security Model - Educational L1 Blockchain

> **🔐 How Security Actually Works**
>
> This document explains the security mechanisms in technical detail, including attack vectors, mitigations, and limitations.

> **✅ UPDATE (Jan 23, 2026)**: All security features described below are now **FULLY INTEGRATED** into the runtime. Run with `--bft-mode` to enable BFT consensus, slashing, and instant finality.

---

## Table of Contents
1. [Security Overview](#security-overview)
2. [Byzantine Fault Tolerance](#byzantine-fault-tolerance)
3. [Economic Security (Slashing)](#economic-security-slashing)
4. [Finality Guarantees](#finality-guarantees)
5. [Attack Vectors & Mitigations](#attack-vectors--mitigations)
6. [Security Assumptions](#security-assumptions)
7. [Limitations](#limitations)

---

## Security Overview

### Three-Layer Security Model

```
Layer 1: Consensus Security (BFT)
├── Guarantees: Safety + Liveness (if >2/3 honest validators)
├── Mechanism: Tendermint-style voting (Prevote → Precommit)
└── Tolerance: Up to 1/3 Byzantine (malicious) validators

Layer 2: Economic Security (Slashing)
├── Guarantees: Financial disincentive for attacks
├── Mechanism: Stake confiscation for provable misbehavior
└── Deterrent: Attack cost > Expected gain

Layer 3: Finality (Irreversibility)
├── Guarantees: Committed blocks cannot be reversed
├── Mechanism: 2/3+ validator signatures required
└── Protection: No long-range reorg attacks
```

---

## Byzantine Fault Tolerance

### How BFT Works

**Problem**: In a distributed network, some nodes may be **Byzantine** (malicious/faulty).

**Solution**: Use voting with **2/3+ majority** requirement.

#### Mathematical Proof (Simplified)

```
Total Validators: N
Byzantine Validators: f
Honest Validators: N - f

Safety Condition:
- Need 2/3+ votes to finalize
- Byzantine can vote arbitrarily
- Maximum Byzantine: f < N/3

Proof:
If f < N/3, then:
- Honest validators: N - f > 2N/3
- Byzantine cannot create 2/3+ on conflicting blocks
- At most ONE block can get 2/3+ (safety)

If f >= N/3:
- Byzantine + dishonest subset could create 2/3+ on two different blocks
- Safety violated (fork)
```

### Implementation

**File**: `internal/consensus/bft/voting.go`

```go
// Step 1: Prevote Phase
// Validators vote on proposed block
func (voteSet *VoteSet) HasTwoThirdsMajority() (bool, [32]byte) {
    for blockHash, votingPower := range voteSet.votesByBlock {
        // Check if this block has >2/3 voting power
        if votingPower > (voteSet.totalVotingPower * 2 / 3) {
            return true, blockHash  // SAFE: only one block can reach 2/3+
        }
    }
    return false, [32]byte{}
}

// Step 2: Precommit Phase
// If 2/3+ prevoted, validators precommit
// If 2/3+ precommit → block is FINALIZED

// Security Property: Two conflicting blocks cannot both get 2/3+ votes
// Why? Because 2/3 + 2/3 = 4/3 > total (impossible)
```

### Attack Resistance

**Attack 1: Double Spend**
```
Attacker tries to create conflicting transactions:
├── TX1: Pay Merchant (in Block A)
└── TX2: Pay Self (in Block B)

Defense:
1. Both blocks submitted to network
2. Validators prevote (only ONE can get 2/3+)
3. Block A gets 2/3+ → finalized
4. Block B rejected (conflicting TX)
Result: Double spend PREVENTED ✅
```

**Attack 2: Censorship**
```
Attacker (proposer) excludes victim's transactions

Defense (Limited):
1. Current round: Attacker can censor (if they're proposer)
2. Next round: Different proposer (round-robin)
3. Victim's TX included in future block
Limitation: Temporary censorship possible ⚠️
```

---

## Economic Security (Slashing)

### How Slashing Works

**Principle**: Make misbehavior economically irrational.

#### Slashing Conditions

**File**: `internal/slashing/tracker.go`

**1. Double-Signing**
```go
// Offense: Signing two different blocks at same height
type DoubleSignEvidence struct {
    Vote1 Vote  // Signature on Block A
    Vote2 Vote  // Signature on Block B (at same height)
}

// Penalty: 100% stake slashed (tombstoned forever)
slashAmount := validator.Stake * 1.00
```

**Why This Matters**:
```
Scenario: Validator tries to fork chain

Step 1: Validator signs Block A (height 100)
Step 2: Validator signs Block B (height 100, different hash)
Step 3: Other validators detect double-sign
Step 4: Submit DoubleSignEvidence
Step 5: Validator loses entire stake + banned forever

Economic Analysis:
- Attack gain: Potential double-spend profit
- Attack cost: 100% of validator stake
- Rational behavior: Don't attack (cost > gain)
```

**2. Downtime**
```go
// Offense: Missing too many votes
type DowntimeEvidence struct {
    MissedHeights []uint64
    Threshold     uint64  // e.g., 100 consecutive misses
}

// Penalty: 1% stake slashed (warning)
slashAmount := validator.Stake * 0.01
```

**Why This Matters**:
```
Liveness Requirement: Network needs >2/3 validators online

If validator consistently offline:
- Reduces available voting power
- Risks network halt (if too many offline)
- Small slash incentivizes uptime
```

### Economic Game Theory

**Attack Scenario**: 34% Attack (Just Above Safety Threshold)

```
Setup:
- Total stake: 1000 RNR
- Attacker stake: 340 RNR (34%)
- Attack: Try to double-spend $1M

Analysis:
Success Probability:
- Need 34% to prevent finality (2/3 threshold)
- Can censor blocks, but cannot create fake 2/3+
- Result: Network halts, no finality

Attacker Cost:
- Slashed: 340 RNR (100% of stake)
- Value lost: $340,000 (assuming 1 RNR = $1000)

Attacker Gain:
- Double-spend: $1M (if successful)

Problem:
- Gain > Cost ($1M > $340k)
- Attack is RATIONAL! ⚠️

Mitigation Required:
- Increase min stake requirement
- Ensure total stake value > 3x max transaction value
```

---

## Finality Guarantees

### How Finality Works

**Finality** = Block is irreversible (cannot be reorganized)

**File**: `internal/finality/tracker.go`

```go
// Once block gets 2/3+ precommits, it is FINAL
func (ft *FinalityTracker) MarkFinalized(height uint64, hash [32]byte) {
    ft.FinalizedHeight = height
    ft.FinalizedHash = hash
    fmt.Printf("Block %d FINALIZED (irreversible)\n", height)
}

// Reorg protection
func (ft *FinalityTracker) CanReorg(height uint64) bool {
    return height > ft.FinalizedHeight  // Cannot reorg finalized blocks
}
```

### Why This Matters

**Bitcoin (No Finality)**:
```
Block Confirmations: 6 blocks (~60 min)
Finality: Probabilistic (never 100%)
Reorg Risk: Small but non-zero

If attacker has 40% hashpower:
- Can reorg last 6 blocks with ~1% probability
- Expensive but possible
```

**Our L1 (Instant Finality)**:
```
Block Confirmations: 1 block (~6 sec)
Finality: Absolute (100% after 2/3+ votes)
Reorg Risk: ZERO (mathematically impossible)

Why impossible:
- Need 2/3+ validators to sign conflicting block
- Would require 2/3+ to double-sign
- All would be slashed (economic suicide)
```

---

## Attack Vectors & Mitigations

### 1. 51% Attack (PoW Style)

**Attack**: Control >50% hashpower, rewrite history

**Mitigation**: 
- N/A (we don't use PoW for consensus)
- PoW only used for spam prevention (lightweight)

### 2. 34% Attack (BFT Liveness)

**Attack**: Control 34% stake, prevent finality

```
Scenario:
- Total validators: 100
- Attacker: 34 validators
- Required for finality: 67 votes

Attack:
1. Attacker refuses to vote
2. Only 66 honest votes (< 67 needed)
3. Network cannot finalize blocks

Mitigation:
- Downtime slashing (1% per infraction)
- After 100 infractions: attacker loses 34% → 0%
- Network recovers liveness
```

### 3. Long-Range Attack

**Attack**: Rewrite history from genesis (with old validator keys)

```
Scenario:
- Attacker unbonds stake (exits validator set)
- Keeps private keys
- Years later: creates alternate chain from genesis
- Uses old keys to sign fake history

Mitigation (Finality):
- Blocks finalized with 2/3+ signatures
- Checkpoints stored (every 10,000 blocks)
- New nodes sync from checkpoint (not genesis)
- Old keys from pre-checkpoint are INVALID
```

### 4. Sybil Attack

**Attack**: Create many fake identities to gain influence

```
Attack:
- Create 1000 validator identities
- Try to gain >2/3 voting power

Mitigation (Stake Requirement):
- Each validator needs minimum 1000 RNR stake
- 1000 validators = 1,000,000 RNR required
- Attacker must buy >2/3 of total supply
- Economic: cost > $10M (prohibitive)
```

### 5. DDoS Attack

**Attack**: Flood network with spam transactions

```
Attack:
- Send 1 million spam TXs/second
- Overwhelm mempool

Mitigation:
- Transaction fees (anti-spam)
- Mempool size limit (50 MB)
- Block size limit (10 MB)
- Oldest/lowest-fee TXs dropped first
```

---

## Security Assumptions

### Trust Model

**We Assume**:
1. ✅ **>2/3 validators are honest** (BFT assumption)
2. ✅ **Cryptography is secure** (Ed25519 unbreakable)
3. ✅ **Network is partially synchronous** (messages delivered eventually)
4. ✅ **Stake has economic value** (slashing creates cost)

**We Do NOT Assume**:
1. ❌ All validators are honest (tolerate up to 1/3 Byzantine)
2. ❌ Network is perfectly synchronous (tolerate delays)
3. ❌ No validators go offline (liveness degrades gracefully)

### Failure Modes

**If >1/3 validators are Byzantine**:
- ❌ Safety violated (conflicting blocks possible)
- ❌ Finality violated (double-finality possible)
- ⚠️ **Network forks** (requires manual intervention)

**If >1/3 validators are offline**:
- ✅ Safety preserved (no conflicting blocks)
- ❌ Liveness violated (no new blocks finalized)
- ⚠️ **Network halts** (waits for validators to return)

---

## Limitations

### Known Security Gaps

1. **No Social Consensus Layer**
   - If chain forks, no automatic resolution
   - Requires community/governance to choose canonical chain

2. **Stake Grinding**
   - Validators could try to manipulate proposer selection
   - Mitigation: Weighted round-robin (deterministic, not stake-weighted)

3. **Weak Subjectivity**
   - New nodes must trust checkpoint (cannot sync from genesis)
   - Mitigation: Rely on social consensus for checkpoint

4. **Centralization Risks**
   - High stake requirement → fewer validators
   - Trade-off: security vs decentralization

---

## Educational Takeaways

### What This Teaches

1. **BFT ≠ PoW**: Different security model, different attacks
2. **Economic Security**: Slashing creates real-world costs
3. **Finality**: Not all blockchains have instant finality
4. **Trade-offs**: Security vs Decentralization vs Scalability

### Comparison with Production Chains

| Feature | This L1 | Cosmos | Ethereum PoS |
|---------|---------|--------|--------------|
| BFT Consensus | ✅ Tendermint-style | ✅ Tendermint | ✅ Casper FFG |
| Instant Finality | ✅ Yes | ✅ Yes | ⚠️ ~15 min |
| Slashing | ✅ Double-sign | ✅ Double-sign | ✅ Double-sign + Inactivity |
| Max Validators | ~100 | ~175 | ~1M |
| Security | Educational | Production | Production |

---

## Conclusion

**Security is achieved through**:
1. ✅ **Mathematics**: BFT theory (2/3+ majority)
2. ✅ **Economics**: Slashing (financial deterrent)
3. ✅ **Cryptography**: Ed25519 signatures (unforgeable)

**Security is LIMITED by**:
1. ⚠️ **Trust assumption**: >2/3 validators honest
2. ⚠️ **Economic scale**: Stake value must exceed attack value
3. ⚠️ **Centralization**: High stake requirement reduces validator count

**For Educational Use**: This demonstrates real L1 security mechanisms.
**For Production Use**: Would need additional hardening, audits, and economic analysis.

---

**Learn More**:
- [BFT Consensus Theory](https://pmg.csail.mit.edu/papers/osdi99.pdf) - Original PBFT paper
- [Tendermint Consensus](https://arxiv.org/abs/1807.04938) - Tendermint specification
- [Ethereum Casper](https://arxiv.org/abs/1710.09437) - PoS finality gadget
