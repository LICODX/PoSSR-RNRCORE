# RnR Core Protocol: Technical White Paper
**Version 1.0.0**

## 1. Abstract
RnR Core is a Layer-1 blockchain protocol compliant with the Proof-of-Sequential-Sorting-Race (PoSSR) consensus mechanism. It integrates a hybrid architecture utilizing Proof-of-Work (PoW) for identity generation and Sybil resistance, PoSSR for deterministic leader election, and Tendermint-style Byzantine Fault Tolerance (BFT) for block finality. This paper defines the protocol specification, network topology, and economic model.

## 2. Consensus Architecture

### 2.1 Hybrid Consensus Design
The protocol decouples **Leader Election** from **Validation** to optimize for both throughput and security.

1.  **Identity Layer (PoW)**: Nodes generate identities using a memory-hard PoW puzzle. This serves as a Sybil resistance mechanism, imposing a computational cost on identity creation ($Cost > 0$).
2.  **Election Layer (PoSSR)**: A deterministic algorithm selects a block proposer. Nodes compete to sort a pseudo-random dataset seeded by the previous block's Verifiable Random Function (VRF) output. The efficiency of sorting ($O(N \log N)$) vs verification ($O(N)$) provides the consensus asymmetry.
3.  **Finality Layer (BFT)**: A committee of weighted validators votes on the proposal using a 3-phase commit protocol (Propose $\to$ Prevote $\to$ Precommit).

### 2.2 Consenus State Machine
The BFT engine operates as a state machine with the following transitions:
- **NewHeight**: Enter new block height.
- **Propose**: Leader broadcasts `ProposalBlock`.
- **Prevote**: Validators broadcast `Prevote(Hash)` upon verification. Lock on $2/3+$ majority.
- **Precommit**: Validators broadcast `Precommit(Hash)` upon seeing Polka ($2/3+$ Prevotes).
- **Commit**: Block is finalized upon seeing $2/3+$ Precommits.

## 3. Network Topology

### 3.1 P2P Layer
The network utilizes `libp2p` with GossipSub routing. Nodes subscribe to:
- **Global Topics**: `/rnr/consensus/vote`, `/rnr/consensus/proposal` (High priority).
- **Shard Topics**: `/rnr/shard/{id}` (Data availability).

### 3.2 Sharding Model
The state is partitioned into 10 fixed shards ($S_0 \dots S_9$).
- **Assignment**: Validators are assigned to subsets of shards via a deterministic Round-Robin schedule: $ShardID = ValidatorID \pmod{10}$.
- **Propagation**: Transactions are gossiped only within their respective shard topics to reduce global bandwidth consumption.

## 4. Economic Model

### 4.1 Incentive Structure
Block rewards ($R_{total}$) are distributed based on shard participation.
$$ R_{validator} = R_{base} \times \frac{S_{processed}}{S_{total}} $$
Where $S_{processed}$ is the count of valid shard proofs submitted by the validator.

### 4.2 Slashing Conditions
The protocol enforces economic security via automated slashing.

| Offense | Condition | Penalty | Action |
| :--- | :--- | :--- | :--- |
| **Equivocation** | Two conflicting signatures at same ($H, R$) | 100% Stake | Tombstone (Permaban) |
| **Availability** | Missed $N$ consecutive votes | 1% Stake | Jail (24h Suspension) |

## 5. Security Specifications

### 5.1 Threat Model
The system assumes a partially synchronous network with $N$ validators where $f < N/3$ are Byzantine.
- **Safety**: Guaranteed for $f < N/3$. Two conflicting blocks cannot reach $2/3+$ quorum.
- **Liveness**: Guaranteed for $N - f > 2N/3$. The network advances if supermajority is online and honest.

### 5.2 Algorithm Agility
To mitigate ASIC optimization for sorting, the protocol rotates sorting algorithms per block based on entropy:
$$ Algo = VRF(H_{prev}) \pmod 7 $$
Supported algorithms: QuickSort, MergeSort, HeapSort, RadixSort, TimSort, IntroSort, ShellSort.

## 6. References
1.  Castro, M., & Liskov, B. (1999). "Practical Byzantine Fault Tolerance".
2.  Kwon, J. (2014). "Tendermint: Consensus without Mining".
