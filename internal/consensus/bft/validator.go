package bft

import (
	"crypto/ed25519"
	"fmt"
	"sort"
)

// Validator represents a single validator in the network
type Validator struct {
	Address          [32]byte          // Validator's address (public key hash)
	PubKey           ed25519.PublicKey // Ed25519 public key for signature verification
	VotingPower      uint64            // Stake-based voting power
	ProposerPriority int64             // Used for round-robin proposer selection
}

// ValidatorSet represents the active set of validators
type ValidatorSet struct {
	Validators       []*Validator
	Proposer         *Validator
	totalVotingPower uint64
}

// NewValidatorSet creates a new validator set
func NewValidatorSet(validators []*Validator) *ValidatorSet {
	if len(validators) == 0 {
		return &ValidatorSet{
			Validators: []*Validator{},
		}
	}

	// Sort validators by address for deterministic ordering
	sort.Slice(validators, func(i, j int) bool {
		for k := 0; k < 32; k++ {
			if validators[i].Address[k] != validators[j].Address[k] {
				return validators[i].Address[k] < validators[j].Address[k]
			}
		}
		return false
	})

	vs := &ValidatorSet{
		Validators: validators,
	}
	vs.computeTotalVotingPower()
	vs.IncrementProposerPriority(1) // Set initial proposer
	return vs
}

// computeTotalVotingPower calculates total voting power
func (vs *ValidatorSet) computeTotalVotingPower() {
	total := uint64(0)
	for _, val := range vs.Validators {
		total += val.VotingPower
	}
	vs.totalVotingPower = total
}

// TotalVotingPower returns the total voting power
func (vs *ValidatorSet) TotalVotingPower() uint64 {
	return vs.totalVotingPower
}

// IncrementProposerPriority increments proposer priority and selects new proposer
// This implements Tendermint's weighted round-robin algorithm
func (vs *ValidatorSet) IncrementProposerPriority(times int) {
	for i := 0; i < times; i++ {
		vs.incrementProposerPriorityOnce()
	}
}

func (vs *ValidatorSet) incrementProposerPriorityOnce() {
	if len(vs.Validators) == 0 {
		return
	}

	// Increment each validator's priority by their voting power
	for _, val := range vs.Validators {
		val.ProposerPriority += int64(val.VotingPower)
	}

	// Select validator with highest priority as proposer
	maxPriority := vs.Validators[0]
	for _, val := range vs.Validators {
		if val.ProposerPriority > maxPriority.ProposerPriority {
			maxPriority = val
		}
	}

	// Set as proposer and decrease priority by total voting power
	vs.Proposer = maxPriority
	maxPriority.ProposerPriority -= int64(vs.totalVotingPower)
}

// GetProposer returns the current proposer
func (vs *ValidatorSet) GetProposer() *Validator {
	if vs.Proposer == nil && len(vs.Validators) > 0 {
		vs.IncrementProposerPriority(1)
	}
	return vs.Proposer
}

// GetByAddress returns validator by address
func (vs *ValidatorSet) GetByAddress(address [32]byte) *Validator {
	for _, val := range vs.Validators {
		if val.Address == address {
			return val
		}
	}
	return nil
}

// HasAddress checks if validator exists in set
func (vs *ValidatorSet) HasAddress(address [32]byte) bool {
	return vs.GetByAddress(address) != nil
}

// HasTwoThirdsMajority checks if given voting power exceeds 2/3 threshold
// This is the key BFT safety guarantee
func (vs *ValidatorSet) HasTwoThirdsMajority(votingPower uint64) bool {
	// Need strictly more than 2/3 (not equal)
	// 2/3 * total = threshold
	// votingPower > threshold
	threshold := (vs.totalVotingPower * 2) / 3
	return votingPower > threshold
}

// Size returns number of validators
func (vs *ValidatorSet) Size() int {
	return len(vs.Validators)
}

// Copy creates a deep copy of the validator set
func (vs *ValidatorSet) Copy() *ValidatorSet {
	validators := make([]*Validator, len(vs.Validators))
	for i, v := range vs.Validators {
		valCopy := *v
		validators[i] = &valCopy
	}
	return NewValidatorSet(validators)
}

// Add adds a new validator to the set
func (vs *ValidatorSet) Add(val *Validator) error {
	if vs.HasAddress(val.Address) {
		return fmt.Errorf("validator already exists: %x", val.Address[:4])
	}
	vs.Validators = append(vs.Validators, val)
	vs.computeTotalVotingPower()
	return nil
}

// Remove removes a validator from the set
func (vs *ValidatorSet) Remove(address [32]byte) error {
	idx := -1
	for i, val := range vs.Validators {
		if val.Address == address {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("validator not found: %x", address[:4])
	}

	// Remove validator
	vs.Validators = append(vs.Validators[:idx], vs.Validators[idx+1:]...)
	vs.computeTotalVotingPower()

	// Reset proposer if removed
	if vs.Proposer != nil && vs.Proposer.Address == address {
		vs.Proposer = nil
	}

	return nil
}

// UpdateVotingPower updates a validator's voting power
func (vs *ValidatorSet) UpdateVotingPower(address [32]byte, newPower uint64) error {
	val := vs.GetByAddress(address)
	if val == nil {
		return fmt.Errorf("validator not found: %x", address[:4])
	}
	val.VotingPower = newPower
	vs.computeTotalVotingPower()
	return nil
}
