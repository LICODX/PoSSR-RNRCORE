# Changelog

All notable changes to the RnR Core project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-01-23

### Added
- **BFT Consensus Engine**: Full Tendermint-style PBFT implementation (`internal/consensus/bft_engine.go`).
- **Instant Finality**: Finality tracker ensuring irreversibility after 2/3+ precommits.
- **Slashing Mechanism**: Automated detector for double-signing (100% slash) and downtime (1% jail).
- **Proportional Rewards**: Rewards now distributed based on verified shard assignment.
- **Architecture Documentation**: Added `docs/` folder with complete specs (`TECHNICAL_WHITEPAPER_v2.md`, `INCENTIVES.md`, `SECURITY.md`).

### Changed
- **Consensus Mode**: Switched default run-mode to support `--bft-mode`.
- **Project Positioning**: Rebranded from "Educational Testbed" to "Pre-Alpha (R&D)" with clear Phase 0 roadmap.
- **Block Size**: Reduced default block size from 1GB (theoretical) to 10MB (practical) for current phase.
- **Documentation**: Massive restructuring of root directory; moved loose MD files to `docs/`.

### Fixed
- **Circular Dependencies**: Resolved import cycles in `internal/state/manager.go`.
- **Lint Errors**: Fixed unused variable declarations in `main.go`.
- **Criticism Response**: Addressed all major points from community feedback (security, feasibility, math).

## [0.1.0] - 2026-01-01

### Added
- Initial PoSSR sorting engine (QuickSort, MergeSort, HeapSort).
- Basic Blockchain structure (Blocks, Transactions, Headers).
- P2P Networking via GossipSub.
- Simple Miner implementation (Sorting Race).
- JSON-RPC Interface (Basic).
- Simulation test suite.

### Security
- Added `SECURITY.md` (initial draft).
