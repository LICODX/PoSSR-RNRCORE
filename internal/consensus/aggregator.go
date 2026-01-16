package consensus

import (
	"fmt"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// ShardProof represents a proof submission from a shard winner
type ShardProof struct {
	SlotID    uint8
	NodeID    [32]byte
	Proof     [32]byte // Merkle root of sorted data
	Timestamp int64
	Signature [64]byte
}

// Aggregator collects and validates shard proofs
type Aggregator struct {
	proofs map[uint64]map[uint8]*ShardProof // blockHeight -> slotID -> proof
	mu     sync.RWMutex
}

// NewAggregator creates a new consensus aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		proofs: make(map[uint64]map[uint8]*ShardProof),
	}
}

// SubmitProof adds a shard proof to the aggregator
func (a *Aggregator) SubmitProof(height uint64, proof *ShardProof) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Validate slot ID
	if proof.SlotID >= 10 {
		return fmt.Errorf("invalid slot ID: %d", proof.SlotID)
	}

	// Verify signature on proof
	message := append(proof.Proof[:], byte(proof.SlotID))
	if !utils.Verify(proof.NodeID[:], message, proof.Signature[:]) {
		return fmt.Errorf("invalid proof signature")
	}

	// Initialize map for this height if needed
	if a.proofs[height] == nil {
		a.proofs[height] = make(map[uint8]*ShardProof)
	}

	// Check if slot already filled
	if existing, ok := a.proofs[height][proof.SlotID]; ok {
		// Keep earliest submission
		if proof.Timestamp > existing.Timestamp {
			return fmt.Errorf("slot %d already filled", proof.SlotID)
		}
	}

	a.proofs[height][proof.SlotID] = proof
	return nil
}

// GetProofs returns all proofs for a given height
func (a *Aggregator) GetProofs(height uint64) map[uint8]*ShardProof {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if proofs, ok := a.proofs[height]; ok {
		// Return copy
		result := make(map[uint8]*ShardProof)
		for k, v := range proofs {
			result[k] = v
		}
		return result
	}
	return nil
}

// IsComplete checks if we have all 10 shard proofs
func (a *Aggregator) IsComplete(height uint64) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	proofs, ok := a.proofs[height]
	if !ok {
		return false
	}
	return len(proofs) == 10
}

// CreateAggregatedBlock combines 10 shard proofs into a block
func (a *Aggregator) CreateAggregatedBlock(height uint64, prevHeader types.BlockHeader) (*types.Block, error) {
	proofs := a.GetProofs(height)
	if len(proofs) < 10 {
		return nil, fmt.Errorf("incomplete proofs: %d/10", len(proofs))
	}

	// Build aggregated block
	block := &types.Block{
		Header: types.BlockHeader{
			Version:       1,
			PrevBlockHash: types.HashBlockHeader(prevHeader),
			Height:        height,
			Timestamp:     time.Now().Unix(),
			VRFSeed:       prevHeader.VRFSeed, // TODO: Generate new seed
		},
	}

	// Collect all shard roots
	var shardRoots [][32]byte
	for i := uint8(0); i < 10; i++ {
		if proof, ok := proofs[i]; ok {
			block.Header.WinningNodes[i] = proof.NodeID
			shardRoots = append(shardRoots, proof.Proof)

			// Add shard data (simplified - in reality we need actual tx data)
			block.Shards[i] = types.ShardData{
				NodeID:    proof.NodeID,
				ShardRoot: proof.Proof,
			}
		}
	}

	// Calculate aggregated Merkle root
	block.Header.MerkleRoot = utils.CalculateMerkleRoot(shardRoots)

	return block, nil
}

// Cleanup removes old proofs to prevent memory leak
func (a *Aggregator) Cleanup(currentHeight uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Keep only last 10 blocks worth of proofs
	for height := range a.proofs {
		if height < currentHeight-10 {
			delete(a.proofs, height)
		}
	}
}
