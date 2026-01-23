# PoSSR Economic Model & Incentive Structure

> **Abstract**: This document outlines the economic incentives, reward distribution mechanisms, and penalty (slashing) conditions for the RnR Core blockchain network. The system is designed to incentivize honest validation, ensure liveness, and fairly distribute rewards based on computational and verification work.

---

## 1. Hybrid Reward Mechanism

The RnR protocol utilizes a dual-layer reward system that combines Proof-of-Work (PoW) for spam prevention with Proof-of-Sequential-Sorting-Race (PoSSR) for leader election and BFT consensus for finality.

### 1.1 Block Reward Components

Changes to the block reward are governed by a dynamic decay function, but the distribution logic follows this breakdown:

| Component | Share | Purpose |
|-----------|-------|---------|
| **Proposer Reward** | 40% | Incentivizes fast sorting and block proposal (Leader Election) |
| **Validator Reward** | 60% | Distributed among committee members for BFT signatures (Shards) |

> **Note**: In the current Genesis Phase config (`--bft-mode`), the Proposer is simply the first valid block generator, and rewards are distributed proportionally to shard processing assignments.

---

## 2. Shard-Based Reward Distribution

Unlike traditional monolithic blockchains where the winner takes all, RnR utilizes a **Proportional Shard Reward** system.

### 2.1 The Logic
- The network state is divided into **10 Shards**.
- Validators are assigned shards via a deterministic **Round-Robin** algorithm.
- A validator only earns rewards for the shards they successfully validate and sign.

### 2.2 Mathematical Formula

$$ R_{validator} = \frac{R_{total} \times S_{processed}}{S_{total}} $$

Where:
- $R_{total}$: Total block reward available for validation layer.
- $S_{processed}$: Number of shards assigned to and processed by the validator.
- $S_{total}$: Total shards in the network (fixed at 10).

**Example Scenario**:
- Total Reward: 100 RnR
- Validator A processes 3 shards: Earns 30 RnR
- Validator B processes 7 shards: Earns 70 RnR

This ensures fair compensation for actual work performed and disincentivizes "lazy" validators.

---

## 3. Economic Security (Slashing)

To ensure Byzantine Fault Tolerance (BFT), the network implements strict economic penalties for malicious behavior.

### 3.1 Double-Signing (Equivocation)
**Definition**: A validator signs two different block headers at the same height and round.
**Penalty**: **100% Slashing** (Total Stake Burn).
**Consequence**: Immediate removal from the active set ("Tombstoning").
**Rationale**: Double-signing is an explicit attack on consensus integrity (attempting to fork the chain).

### 3.2 Downtime (Liveness Faults)
**Definition**: Failing to participate in consensus votes for a defined window of blocks (e.g., 100 blocks).
**Penalty**: **1% Slashing** (Warning Penalty).
**Consequence**: Temporary freezing (Jail) from the active set.
**Rationale**: Validators must maintain high uptime to ensure the network can reach the 2/3+ threshold for finality.

---

## 4. Genesis Parameters (Phase 0)

For the current `Pre-Alpha` phase, the economic parameters are configured for stability over aggression:

- **Initial Stake Requirement**: 32 RnR (Proposed)
- **Unbonding Period**: 21 Days (To prevent Long-Range attacks)
- **Max Validators**: 100 (Capped for initial BFT network optimization)

## 5. Future Roadmap: DAO & Governance

In Phase 2 (Scalability), reward parameters will be adjustable via on-chain governance proposals, allowing the network to adapt to changing hardware costs and bandwidth availability.
