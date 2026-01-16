package consensus

import (
	"fmt"
	"sync"
)

// Vote represents a node's vote on block validity
type Vote struct {
	NodeID    [32]byte
	BlockHash [32]byte
	Height    uint64
	Approve   bool // true = accept, false = reject
	Signature [64]byte
}

// VoteTracker tracks votes for blocks
type VoteTracker struct {
	votes map[uint64]map[[32]byte][]Vote // height -> blockHash -> votes
	mu    sync.RWMutex
}

// NewVoteTracker creates a new vote tracker
func NewVoteTracker() *VoteTracker {
	return &VoteTracker{
		votes: make(map[uint64]map[[32]byte][]Vote),
	}
}

// SubmitVote adds a vote to the tracker
func (vt *VoteTracker) SubmitVote(vote Vote) error {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	// Initialize maps if needed
	if vt.votes[vote.Height] == nil {
		vt.votes[vote.Height] = make(map[[32]byte][]Vote)
	}

	// Check for duplicate vote from same node
	existingVotes := vt.votes[vote.Height][vote.BlockHash]
	for _, v := range existingVotes {
		if v.NodeID == vote.NodeID {
			return fmt.Errorf("duplicate vote from node %x", vote.NodeID)
		}
	}

	// Add vote
	vt.votes[vote.Height][vote.BlockHash] = append(existingVotes, vote)
	return nil
}

// CountVotes returns vote counts for a block
func (vt *VoteTracker) CountVotes(height uint64, blockHash [32]byte) (approve int, reject int) {
	vt.mu.RLock()
	defer vt.mu.RUnlock()

	votes, ok := vt.votes[height][blockHash]
	if !ok {
		return 0, 0
	}

	for _, vote := range votes {
		if vote.Approve {
			approve++
		} else {
			reject++
		}
	}
	return
}

// HasMajority checks if block has 7/10 approval votes
func (vt *VoteTracker) HasMajority(height uint64, blockHash [32]byte) bool {
	approve, _ := vt.CountVotes(height, blockHash)
	return approve >= 7
}

// GetWinningBlock returns the block with majority votes at height
func (vt *VoteTracker) GetWinningBlock(height uint64) (*[32]byte, bool) {
	vt.mu.RLock()
	defer vt.mu.RUnlock()

	blocksAtHeight, ok := vt.votes[height]
	if !ok {
		return nil, false
	}

	for blockHash := range blocksAtHeight {
		approve, _ := vt.CountVotes(height, blockHash)
		if approve >= 7 {
			hash := blockHash
			return &hash, true
		}
	}
	return nil, false
}

// Cleanup removes old votes
func (vt *VoteTracker) Cleanup(currentHeight uint64) {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	for height := range vt.votes {
		if height < currentHeight-10 {
			delete(vt.votes, height)
		}
	}
}
