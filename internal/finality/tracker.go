package finality

import (
	"fmt"
	"sync"
)

// FinalityTracker tracks finalized blocks (irreversible commits)
type FinalityTracker struct {
	mu sync.RWMutex

	// Finalized state
	FinalizedHeight uint64
	FinalizedHash   [32]byte

	// Checkpoints (every N blocks)
	Checkpoints        map[uint64][32]byte
	CheckpointInterval uint64
}

// NewFinalityTracker creates a new finality tracker
func NewFinalityTracker(checkpointInterval uint64) *FinalityTracker {
	return &FinalityTracker{
		FinalizedHeight:    0,
		CheckpointInterval: checkpointInterval,
		Checkpoints:        make(map[uint64][32]byte),
	}
}

// MarkFinalized marks a block as finalized (irreversible)
// This happens when a block receives 2/3+ precommit votes
func (ft *FinalityTracker) MarkFinalized(height uint64, hash [32]byte) error {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	// Can only advance finality forward
	if height <= ft.FinalizedHeight {
		return fmt.Errorf("cannot finalize height %d: already at %d", height, ft.FinalizedHeight)
	}

	ft.FinalizedHeight = height
	ft.FinalizedHash = hash

	// Add checkpoint if at interval
	if height%ft.CheckpointInterval == 0 {
		ft.Checkpoints[height] = hash
		fmt.Printf("[Finality] Checkpoint created at height %d\n", height)
	}

	fmt.Printf("[Finality] Block finalized: Height %d, Hash %x\n", height, hash[:8])

	return nil
}

// CanReorg checks if a height can be reorganized
// Finalized blocks cannot be reorganized
func (ft *FinalityTracker) CanReorg(height uint64) bool {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	return height > ft.FinalizedHeight
}

// GetFinalizedHeight returns the current finalized height
func (ft *FinalityTracker) GetFinalizedHeight() uint64 {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	return ft.FinalizedHeight
}

// GetFinalizedHash returns the hash of the finalized block
func (ft *FinalityTracker) GetFinalizedHash() [32]byte {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	return ft.FinalizedHash
}

// GetCheckpoint returns checkpoint hash at given height
func (ft *FinalityTracker) GetCheckpoint(height uint64) ([32]byte, bool) {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	hash, exists := ft.Checkpoints[height]
	return hash, exists
}

// PruneOldCheckpoints removes checkpoints older than keepDepth
func (ft *FinalityTracker) PruneOldCheckpoints(keepDepth uint64) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	minHeight := uint64(0)
	if ft.FinalizedHeight > keepDepth {
		minHeight = ft.FinalizedHeight - keepDepth
	}

	for height := range ft.Checkpoints {
		if height < minHeight {
			delete(ft.Checkpoints, height)
		}
	}
}

// IsFinalized checks if a specific block is finalized
func (ft *FinalityTracker) IsFinalized(height uint64, hash [32]byte) bool {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	if height > ft.FinalizedHeight {
		return false
	}

	if height == ft.FinalizedHeight {
		return ft.FinalizedHash == hash
	}

	// Check checkpoint
	if checkpointHash, exists := ft.Checkpoints[height]; exists {
		return checkpointHash == hash
	}

	// Older than finalized = assumed finalized
	return true
}
