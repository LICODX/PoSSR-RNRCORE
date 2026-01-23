package consensus

import (
	"crypto/ed25519"
	"fmt"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus/bft"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// BFTEngine orchestrates the BFT consensus protocol
type BFTEngine struct {
	mu sync.RWMutex

	// Core components
	State      *bft.ConsensusState
	Validators *bft.ValidatorSet

	// Node identity
	ValidatorAddress [32]byte
	ValidatorPrivKey ed25519.PrivateKey
	ValidatorIndex   int32

	// Communication channels
	VoteChan     chan *bft.Vote
	ProposalChan chan *bft.Proposal
	BlockChan    chan *types.Block

	// Callbacks for P2P broadcasting
	BroadcastVote     func(*bft.Vote) error
	BroadcastProposal func(*bft.Proposal) error
	MarkFinalized     func(uint64, [32]byte) error // Called when block reaches 2/3+ precommits
}

// NewBFTEngine creates a new BFT consensus engine
func NewBFTEngine(
	height uint64,
	validators *bft.ValidatorSet,
	validatorAddr [32]byte,
	privKey ed25519.PrivateKey,
) *BFTEngine {
	state := bft.NewConsensusState(height, validators)

	// Find our validator index
	valIndex := int32(-1)
	for i, val := range validators.Validators {
		if val.Address == validatorAddr {
			valIndex = int32(i)
			break
		}
	}

	engine := &BFTEngine{
		State:            state,
		Validators:       validators,
		ValidatorAddress: validatorAddr,
		ValidatorPrivKey: privKey,
		ValidatorIndex:   valIndex,
		VoteChan:         make(chan *bft.Vote, 100),
		ProposalChan:     make(chan *bft.Proposal, 10),
		BlockChan:        make(chan *types.Block, 10),
	}

	return engine
}

// RunConsensusRound executes one round of BFT consensus
// Returns the committed block or error
func (be *BFTEngine) RunConsensusRound(
	height uint64,
	mempool []types.Transaction,
) (*types.Block, error) {
	fmt.Printf("\n[BFT] Starting consensus round for height %d\n", height)

	// Phase 1: NewHeight → Propose
	be.State.EnterPropose(height, 0)

	var proposedBlock *types.Block
	var proposal *bft.Proposal

	// Check if we are the proposer
	proposer := be.Validators.GetProposer()
	weAreProposer := proposer.Address == be.ValidatorAddress

	if weAreProposer {
		// We propose the block
		fmt.Printf("[BFT] We are the proposer (validator %x)\n", be.ValidatorAddress[:4])

		// Create block using existing mining logic (PoW + Sorting)
		// Get previous block header from blockchain
		// For now, create a simplified block
		block, err := be.createBlock(height, mempool)
		if err != nil {
			return nil, fmt.Errorf("failed to create proposal block: %w", err)
		}

		proposedBlock = block

		// Create proposal
		proposal = &bft.Proposal{
			Height:    height,
			Round:     0,
			BlockHash: block.Header.Hash,
			Timestamp: time.Now().Unix(),
			Proposer:  be.ValidatorAddress,
		}

		// Set proposal in consensus state
		if err := be.State.SetProposal(proposal, block); err != nil {
			return nil, fmt.Errorf("failed to set proposal: %w", err)
		}

		// Broadcast proposal to network
		if be.BroadcastProposal != nil {
			if err := be.BroadcastProposal(proposal); err != nil {
				fmt.Printf("[BFT] Warning: Failed to broadcast proposal: %v\n", err)
			}
		}

		// Broadcast block data
		if be.BlockChan != nil {
			be.BlockChan <- block
		}
	} else {
		// Wait for proposal from network
		fmt.Printf("[BFT] Waiting for proposal from proposer %x\n", proposer.Address[:4])

		select {
		case proposal = <-be.ProposalChan:
			if proposal.Proposer != proposer.Address {
				return nil, fmt.Errorf("proposal from wrong proposer")
			}

			// Wait for block data
			select {
			case proposedBlock = <-be.BlockChan:
				if err := be.State.SetProposal(proposal, proposedBlock); err != nil {
					return nil, fmt.Errorf("failed to set proposal: %w", err)
				}
			case <-time.After(be.State.ProposeTimeout):
				return nil, fmt.Errorf("timeout waiting for block data")
			}

		case <-time.After(be.State.ProposeTimeout):
			return nil, fmt.Errorf("timeout waiting for proposal")
		}
	}

	// Phase 2: Prevote
	be.State.EnterPrevote(height, 0)

	// Create and broadcast our prevote
	prevote := &bft.Vote{
		Type:             bft.VoteTypePrevote,
		Height:           height,
		Round:            0,
		BlockHash:        proposedBlock.Header.Hash,
		Timestamp:        time.Now().Unix(),
		ValidatorAddress: be.ValidatorAddress,
		ValidatorIndex:   be.ValidatorIndex,
	}
	prevote.Sign(be.ValidatorPrivKey)

	// Add our own vote
	if err := be.State.AddVote(prevote); err != nil {
		fmt.Printf("[BFT] Warning: Failed to add own prevote: %v\n", err)
	}

	// Broadcast prevote
	if be.BroadcastVote != nil {
		if err := be.BroadcastVote(prevote); err != nil {
			fmt.Printf("[BFT] Warning: Failed to broadcast prevote: %v\n", err)
		}
	}

	// Collect prevotes from network with timeout
	prevoteDeadline := time.Now().Add(be.State.PrevoteTimeout)
	prevoteReached := false

	for time.Now().Before(prevoteDeadline) {
		select {
		case vote := <-be.VoteChan:
			if vote.Type == bft.VoteTypePrevote && vote.Height == height {
				if err := be.State.AddVote(vote); err != nil {
					fmt.Printf("[BFT] Warning: Invalid prevote: %v\n", err)
					continue
				}

				// Check if we have 2/3+ prevotes
				if has2_3, blockHash := be.State.HasTwoThirdsPrevotes(); has2_3 {
					if blockHash == proposedBlock.Header.Hash {
						fmt.Printf("[BFT] Reached 2/3+ prevotes for block %x\n", blockHash[:4])
						prevoteReached = true
						break
					}
				}
			}

		case <-time.After(100 * time.Millisecond):
			// Check periodically
			if has2_3, blockHash := be.State.HasTwoThirdsPrevotes(); has2_3 {
				if blockHash == proposedBlock.Header.Hash {
					prevoteReached = true
					break
				}
			}
		}
	}

	if !prevoteReached {
		return nil, fmt.Errorf("failed to reach 2/3+ prevotes")
	}

	// Phase 3: Precommit
	be.State.EnterPrecommit(height, 0)

	// Create and broadcast our precommit
	precommit := &bft.Vote{
		Type:             bft.VoteTypePrecommit,
		Height:           height,
		Round:            0,
		BlockHash:        proposedBlock.Header.Hash,
		Timestamp:        time.Now().Unix(),
		ValidatorAddress: be.ValidatorAddress,
		ValidatorIndex:   be.ValidatorIndex,
	}
	precommit.Sign(be.ValidatorPrivKey)

	// Add our own vote
	if err := be.State.AddVote(precommit); err != nil {
		fmt.Printf("[BFT] Warning: Failed to add own precommit: %v\n", err)
	}

	// Broadcast precommit
	if be.BroadcastVote != nil {
		if err := be.BroadcastVote(precommit); err != nil {
			fmt.Printf("[BFT] Warning: Failed to broadcast precommit: %v\n", err)
		}
	}

	// Collect precommits from network
	precommitDeadline := time.Now().Add(be.State.PrecommitTimeout)
	commitReached := false

	for time.Now().Before(precommitDeadline) {
		select {
		case vote := <-be.VoteChan:
			if vote.Type == bft.VoteTypePrecommit && vote.Height == height {
				if err := be.State.AddVote(vote); err != nil {
					fmt.Printf("[BFT] Warning: Invalid precommit: %v\n", err)
					continue
				}

				// Check if we have 2/3+ precommits
				if has2_3, blockHash := be.State.HasTwoThirdsPrecommits(); has2_3 {
					if blockHash == proposedBlock.Header.Hash {
						fmt.Printf("[BFT] Reached 2/3+ precommits for block %x\n", blockHash[:4])
						commitReached = true
						break
					}
				}
			}

		case <-time.After(100 * time.Millisecond):
			// Check periodically
			if has2_3, blockHash := be.State.HasTwoThirdsPrecommits(); has2_3 {
				if blockHash == proposedBlock.Header.Hash {
					commitReached = true
					break
				}
			}
		}
	}

	if !commitReached {
		return nil, fmt.Errorf("failed to reach 2/3+ precommits")
	}

	// Phase 4: Commit (Block Finalized!)
	be.State.EnterCommit(height, 0)
	fmt.Printf("[BFT] ✅ Block %d COMMITTED (finalized with 2/3+ votes)\n", height)

	// Mark block as finalized (irreversible)
	if be.MarkFinalized != nil {
		if err := be.MarkFinalized(height, proposedBlock.Header.Hash); err != nil {
			fmt.Printf("[BFT] Warning: Failed to mark block as finalized: %v\n", err)
		}
	}

	// Finalize commit and prepare for next height
	be.State.FinalizeCommit(height)

	// Move proposer to next validator (round-robin)
	be.Validators.IncrementProposerPriority(1)

	return proposedBlock, nil
}

// createBlock creates a block proposal (using existing PoW + Sorting logic)
func (be *BFTEngine) createBlock(height uint64, txs []types.Transaction) (*types.Block, error) {
	// This is a simplified version - in full implementation, this would call
	// the existing MineBlock function or a modified version without PoW difficulty

	// For now, use the existing sorting-based block creation
	// Get previous block header from blockchain (TODO: pass as parameter)
	prevBlock := types.BlockHeader{
		Height: height - 1,
		Hash:   [32]byte{}, // TODO: Get from blockchain
	}

	difficulty := uint64(1) // Minimal PoW for block creation

	// Create a dummy channel for stop signal
	stopChan := make(chan struct{})

	// Mine block using existing logic
	var minerPubKey [32]byte
	copy(minerPubKey[:], be.ValidatorPrivKey.Public().(ed25519.PublicKey))

	block, err := MineBlock(txs, prevBlock, difficulty, stopChan, minerPubKey, be.ValidatorPrivKey)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// ProcessIncomingVote handles votes received from the P2P network
func (be *BFTEngine) ProcessIncomingVote(vote *bft.Vote) {
	// Send to vote channel for consensus processing
	select {
	case be.VoteChan <- vote:
	default:
		fmt.Printf("[BFT] Warning: Vote channel full, dropping vote\n")
	}
}

// ProcessIncomingProposal handles proposals received from the P2P network
func (be *BFTEngine) ProcessIncomingProposal(proposal *bft.Proposal) {
	// Send to proposal channel for consensus processing
	select {
	case be.ProposalChan <- proposal:
	default:
		fmt.Printf("[BFT] Warning: Proposal channel full, dropping proposal\n")
	}
}

// ProcessIncomingBlock handles block data received from the P2P network
func (be *BFTEngine) ProcessIncomingBlock(block *types.Block) {
	// Send to block channel for consensus processing
	select {
	case be.BlockChan <- block:
	default:
		fmt.Printf("[BFT] Warning: Block channel full, dropping block\n")
	}
}
