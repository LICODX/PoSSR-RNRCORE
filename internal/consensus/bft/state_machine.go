package bft

import (
	"fmt"
	"sync"
	"time"
)

// RoundStep represents the current step in consensus
type RoundStep int

const (
	RoundStepNewHeight RoundStep = 0 // Wait for new height commitment
	RoundStepPropose   RoundStep = 1 // Proposer creates block proposal
	RoundStepPrevote   RoundStep = 2 // Validators broadcast prevote
	RoundStepPrecommit RoundStep = 3 // Validators broadcast precommit
	RoundStepCommit    RoundStep = 4 // Block committed (finalized)
)

func (rs RoundStep) String() string {
	switch rs {
	case RoundStepNewHeight:
		return "NewHeight"
	case RoundStepPropose:
		return "Propose"
	case RoundStepPrevote:
		return "Prevote"
	case RoundStepPrecommit:
		return "Precommit"
	case RoundStepCommit:
		return "Commit"
	default:
		return "Unknown"
	}
}

// Proposal represents a block proposal
type Proposal struct {
	Height    uint64
	Round     int32
	BlockHash [32]byte
	Timestamp int64
	Proposer  [32]byte
}

// ConsensusState manages the state machine for BFT consensus
type ConsensusState struct {
	mu sync.RWMutex

	// Current state
	Height uint64
	Round  int32
	Step   RoundStep

	// Validator set
	Validators *ValidatorSet

	// Current round data
	Proposal      *Proposal
	ProposalBlock interface{} // Actual block data (type.Block)
	LockedRound   int32       // Round where we locked on a block
	LockedBlock   [32]byte    // Block we're locked on (prevents flip-flopping)
	ValidRound    int32       // Latest round with valid proposal
	ValidBlock    [32]byte    // Latest valid block

	// Votes
	Prevotes   *VoteSet
	Precommits *VoteSet

	// Timeouts
	ProposeTimeout   time.Duration
	PrevoteTimeout   time.Duration
	PrecommitTimeout time.Duration
	CommitTimeout    time.Duration
}

// NewConsensusState creates a new consensus state machine
func NewConsensusState(height uint64, validators *ValidatorSet) *ConsensusState {
	cs := &ConsensusState{
		Height:      height,
		Round:       0,
		Step:        RoundStepNewHeight,
		Validators:  validators,
		LockedRound: -1,
		ValidRound:  -1,

		// Timeouts (Tendermint default values)
		ProposeTimeout:   3 * time.Second,
		PrevoteTimeout:   1 * time.Second,
		PrecommitTimeout: 1 * time.Second,
		CommitTimeout:    1 * time.Second,
	}

	cs.Prevotes = NewVoteSet(height, 0, VoteTypePrevote, validators)
	cs.Precommits = NewVoteSet(height, 0, VoteTypePrecommit, validators)

	return cs
}

// EnterPropose enters the propose step
func (cs *ConsensusState) EnterPropose(height uint64, round int32) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.Height != height || cs.Round > round || (cs.Round == round && cs.Step != RoundStepNewHeight) {
		return
	}

	cs.Round = round
	cs.Step = RoundStepPropose

	fmt.Printf("[Consensus] Height %d Round %d: Entered PROPOSE\n", height, round)
}

// EnterPrevote enters the prevote step
func (cs *ConsensusState) EnterPrevote(height uint64, round int32) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.Height != height || cs.Round > round || (cs.Round == round && cs.Step >= RoundStepPrevote) {
		return
	}

	cs.Round = round
	cs.Step = RoundStepPrevote

	fmt.Printf("[Consensus] Height %d Round %d: Entered PREVOTE\n", height, round)
}

// EnterPrecommit enters the precommit step
func (cs *ConsensusState) EnterPrecommit(height uint64, round int32) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.Height != height || cs.Round > round || (cs.Round == round && cs.Step >= RoundStepPrecommit) {
		return
	}

	cs.Round = round
	cs.Step = RoundStepPrecommit

	fmt.Printf("[Consensus] Height %d Round %d: Entered PRECOMMIT\n", height, round)
}

// EnterCommit enters the commit step (finalization)
func (cs *ConsensusState) EnterCommit(height uint64, commitRound int32) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.Height != height || cs.Step == RoundStepCommit {
		return
	}

	cs.Step = RoundStepCommit

	fmt.Printf("[Consensus] Height %d: Entered COMMIT (Round %d)\n", height, commitRound)
}

// FinalizeCommit finalizes the commit and moves to next height
func (cs *ConsensusState) FinalizeCommit(height uint64) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.Height != height || cs.Step != RoundStepCommit {
		return
	}

	// Move to next height
	cs.Height = height + 1
	cs.Round = 0
	cs.Step = RoundStepNewHeight
	cs.Proposal = nil
	cs.ProposalBlock = nil
	cs.LockedRound = -1
	cs.LockedBlock = [32]byte{}
	cs.ValidRound = -1
	cs.ValidBlock = [32]byte{}

	// Reset vote sets for new height
	cs.Prevotes = NewVoteSet(cs.Height, 0, VoteTypePrevote, cs.Validators)
	cs.Precommits = NewVoteSet(cs.Height, 0, VoteTypePrecommit, cs.Validators)

	fmt.Printf("[Consensus] Finalized commit. Moving to Height %d\n", cs.Height)
}

// SetProposal sets the current proposal
func (cs *ConsensusState) SetProposal(proposal *Proposal, block interface{}) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if proposal.Height != cs.Height || proposal.Round != cs.Round {
		return fmt.Errorf("proposal height/round mismatch")
	}

	// Check if proposer is valid
	if cs.Validators.GetProposer().Address != proposal.Proposer {
		return fmt.Errorf("invalid proposer")
	}

	cs.Proposal = proposal
	cs.ProposalBlock = block

	return nil
}

// AddVote adds a vote to the appropriate vote set
func (cs *ConsensusState) AddVote(vote *Vote) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var voteSet *VoteSet

	switch vote.Type {
	case VoteTypePrevote:
		if vote.Height == cs.Height && vote.Round == cs.Round {
			voteSet = cs.Prevotes
		}
	case VoteTypePrecommit:
		if vote.Height == cs.Height && vote.Round == cs.Round {
			voteSet = cs.Precommits
		}
	}

	if voteSet == nil {
		return fmt.Errorf("vote for different height/round")
	}

	added, err := voteSet.AddVote(vote)
	if err != nil {
		return err
	}

	if added {
		fmt.Printf("[Consensus] Added %s vote from validator %x for block %x\n",
			vote.Type, vote.ValidatorAddress[:4], vote.BlockHash[:4])
	}

	return nil
}

// HasTwoThirdsPrevotes checks if we have 2/3+ prevotes for any block
func (cs *ConsensusState) HasTwoThirdsPrevotes() (bool, [32]byte) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.Prevotes.HasTwoThirdsMajority()
}

// HasTwoThirdsPrecommits checks if we have 2/3+ precommits for any block
func (cs *ConsensusState) HasTwoThirdsPrecommits() (bool, [32]byte) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.Precommits.HasTwoThirdsMajority()
}

// GetState returns current height, round, and step (thread-safe)
func (cs *ConsensusState) GetState() (uint64, int32, RoundStep) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.Height, cs.Round, cs.Step
}
