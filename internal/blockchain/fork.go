package blockchain

import (
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// ForkDetector handles chain forks
type ForkDetector struct {
	chains map[[32]byte]*ChainState // blockHash -> chain state
}

// ChainState represents alternative chain state
type ChainState struct {
	Height uint64
	Hash   [32]byte
	Weight uint64 // Total difficulty/work
}

// DetectFork checks if incoming block creates a fork
func DetectFork(newBlock types.Block, currentTip types.BlockHeader) bool {
	newPrevHash := newBlock.Header.PrevBlockHash
	currentHash := types.HashBlockHeader(currentTip)

	// Fork if new block doesn't build on current tip
	return newPrevHash != currentHash && newBlock.Header.Height <= currentTip.Height
}

// ResolveFork implements longest chain rule
func (bc *Blockchain) ResolveFork(alternativeChain []types.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(alternativeChain) == 0 {
		return fmt.Errorf("empty alternative chain")
	}

	// Check if alternative chain is longer
	altTip := alternativeChain[len(alternativeChain)-1]
	if altTip.Header.Height <= bc.tip.Height {
		return fmt.Errorf("alternative chain not longer")
	}

	// Validate entire alternative chain
	for i, block := range alternativeChain {
		var prevHeader types.BlockHeader
		if i == 0 {
			// First block should connect to our chain
			// TODO: Load previous header from DB
			prevHeader = bc.tip
		} else {
			prevHeader = alternativeChain[i-1].Header
		}

		// Use node's shard config (assuming full node for fork resolution or bc.shardConfig)
		// Since we are reorganizing, we should use our own config.
		if err := ValidateBlock(block, prevHeader, bc.shardConfig); err != nil {
			return fmt.Errorf("invalid block in alternative chain: %v", err)
		}
	}

	// Reorganize to alternative chain
	fmt.Printf("⚠️  Chain reorganization: switching to longer chain (height %d)\n", altTip.Header.Height)

	// TODO: Rollback transactions from current chain
	// TODO: Apply transactions from alternative chain

	// Update tip
	bc.tip = altTip.Header

	// Save new tip
	// TODO: Persist reorganization to DB

	return nil
}

// GetCommonAncestor finds the block where chains diverged
func GetCommonAncestor(chain1, chain2 []types.BlockHeader) *types.BlockHeader {
	// Simple implementation: check hashes
	for i := len(chain1) - 1; i >= 0; i-- {
		hash1 := types.HashBlockHeader(chain1[i])
		for j := len(chain2) - 1; j >= 0; j-- {
			hash2 := types.HashBlockHeader(chain2[j])
			if hash1 == hash2 {
				return &chain1[i]
			}
		}
	}
	return nil
}
