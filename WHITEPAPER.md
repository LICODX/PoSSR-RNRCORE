# Whitepaper: PoSSR (Proof of Repeated Sorting)

**v2.0 - "Run and Repeat" Core**

## Abstract
PoSSR (Proof of Repeated Sorting) is a novel consensus algorithm designed to address the centralization risks of Proof of Work (PoW) and the plutocratic tendencies of Proof of Stake (PoS). By utilizing a Time-Memory Trade-off mechanism based on computationally hard sorting problems, PoSSR ensures that network control cannot be monopolized by specialized hardware (ASICs) or massive capital accumulation alone.

## 1. Introduction
The fundamental flaw in modern blockchain consensus is the "Economies of Scale" vulnerability.
*   **PoW:** Specialized hardware (ASICs) allows large farms to mine millions of times more efficiently than CPUs.
*   **PoRS Solution:** Sorting is a memory-bound operation. Latency (RAM speed) limits the maximum throughput, creating a natural "Hardware Ceiling" that keeps the playing field level for consumer hardware.

## 2. Core Mechanism: The Sorting Race

### 2.1 The Seed
Every block header contains a `PrevBlockHash`. This hash serves as the seed for a Verifiable Random Function (VRF).

### 2.2 Algorithm Selection
The VRF output deterministically selects one of 6 sorting algorithms for the current block round:
1.  **QuickSort** (Average Case: O(n log n))
2.  **MergeSort** (Stable: O(n log n))
3.  **HeapSort** (In-place: O(n log n))
4.  **RadixSort** (Non-comparative: O(nk))
5.  **TimSort** (Hybrid: Real-world data optimized)
6.  **IntroSort** (Hybrid: Quick/Heap mix)

This unpredictability prevents optimization for a single algorithm (ASIC Resistance).

### 2.3 The "Work"
Miners must:
1.  Generate a list of `N` pseudo-random integers using the Block Header + Nonce.
2.  Sort the list using the selected algorithm.
3.  Hash the sorted list to produce a `MixHash`.
4.  Compare `MixHash` against the network `Difficulty Target`.

If `MixHash < Target`, the block is valid.

## 3. Fairness Validation
In simulated environments (see `docs/audit/massive_simulation_report.md`):
*   **Scenario:** 80% Malicious Nodes vs 20% Honest Nodes.
*   **Advantage:** Malicious nodes given 10% speed boost.
*   **Result:** Honest nodes maintained ~5% block production.
*   **Conclusion:** The probabilistic nature of the search space ensures that "Winner Takes All" does not apply.

## 4. Network Architecture
*   **P2P Layer:** LibP2P / GossipSub v1.1
*   **Serialization:** JSON (for transparency/debugging) -> Protobuf (Roadmap v3)
*   **Block Time:** Target 10 seconds.
*   **Max Block Size:** Dynamic (capped by bandwidth, default 2MB).

## 5. Security Model
*   **Replay Protection:** Nonce-based account state validation.
*   **Sybil Resistance:** CPU/Memory cost of sorting prevents zero-cost identity generation.
*   **Long Range Attacks:** Checkpointing + Max Reorg limit (k=6 blocks).

## 6. Conclusion
PoSSR RNRCORE represents a return to the original vision of "One CPU, One Vote", updated for the modern era with memory-hard algorithms to resist centralization.
