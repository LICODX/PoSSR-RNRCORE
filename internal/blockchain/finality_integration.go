package blockchain

import (
	"github.com/LICODX/PoSSR-RNRCORE/internal/finality"
)

// GetFinalityTracker returns the finality tracker for external access (e.g., BFT engine)
func (bc *Blockchain) GetFinalityTracker() *finality.FinalityTracker {
	return bc.finalityTracker
}

// MarkBlockFinalized marks a block as finalized (called by BFT consensus on 2/3+ precommits)
func (bc *Blockchain) MarkBlockFinalized(height uint64, hash [32]byte) error {
	return bc.finalityTracker.MarkFinalized(height, hash)
}
