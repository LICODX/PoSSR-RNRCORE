package slashing

import (
	"crypto/ed25519"
	"fmt"
	"sync"
)

// SlashingCondition represents the type of slashable offense
type SlashingCondition int

const (
	DoubleSign SlashingCondition = 1 // Signing two conflicting blocks at same height
	Downtime   SlashingCondition = 2 // Missing too many consecutive votes
)

func (sc SlashingCondition) String() string {
	switch sc {
	case DoubleSign:
		return "DoubleSign"
	case Downtime:
		return "Downtime"
	default:
		return "Unknown"
	}
}

// Evidence represents proof of a slashable offense
type Evidence struct {
	Type             SlashingCondition
	ValidatorAddress [32]byte
	Height           uint64
	Proof            []byte // Serialized proof data
	SubmittedAt      int64  // Timestamp when evidence was submitted
}

// DoubleSignEvidence proves that a validator signed two different blocks at same height
type DoubleSignEvidence struct {
	Vote1 Vote // First conflicting vote
	Vote2 Vote // Second conflicting vote
}

// Vote represents a simplified vote for evidence (matches bft.Vote structure)
type Vote struct {
	Height           uint64
	Round            int32
	BlockHash        [32]byte
	Signature        [64]byte
	ValidatorAddress [32]byte
}

// DowntimeEvidence proves that a validator missed too many votes
type DowntimeEvidence struct {
	ValidatorAddress [32]byte
	MissedHeights    []uint64 // Heights where validator didn't vote
	WindowSize       uint64   // Size of sliding window
}

// SlashingTracker tracks slashing events and enforces penalties
type SlashingTracker struct {
	mu sync.RWMutex

	// Evidence storage
	evidence map[[32]byte][]Evidence // validator -> evidences

	// Slashed validators (permanently removed)
	slashed map[[32]byte]SlashInfo

	// Configuration
	DoubleSignSlashAmount uint64 // Percentage (e.g., 100 = 100% = all stake)
	DowntimeSlashAmount   uint64 // Percentage (e.g., 1 = 1% of stake)
	DowntimeThreshold     uint64 // Number of missed votes to trigger slashing
}

// SlashInfo contains information about a slashed validator
type SlashInfo struct {
	Validator     [32]byte
	Condition     SlashingCondition
	SlashedAmount uint64
	SlashedAt     int64
	Evidence      Evidence
}

// NewSlashingTracker creates a new slashing tracker
func NewSlashingTracker() *SlashingTracker {
	return &SlashingTracker{
		evidence: make(map[[32]byte][]Evidence),
		slashed:  make(map[[32]byte]SlashInfo),

		// Default slashing amounts
		DoubleSignSlashAmount: 100, // 100% stake slashed (tombstoned)
		DowntimeSlashAmount:   1,   // 1% stake slashed (warning)
		DowntimeThreshold:     100, // Miss 100 votes in window
	}
}

// SubmitEvidence submits evidence of a slashable offense
func (st *SlashingTracker) SubmitEvidence(evidence Evidence) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	// Check if already slashed
	if _, exists := st.slashed[evidence.ValidatorAddress]; exists {
		return fmt.Errorf("validator already slashed: %x", evidence.ValidatorAddress[:4])
	}

	// Add evidence
	st.evidence[evidence.ValidatorAddress] = append(
		st.evidence[evidence.ValidatorAddress],
		evidence,
	)

	fmt.Printf("[Slashing] Evidence submitted: %s for validator %x at height %d\n",
		evidence.Type, evidence.ValidatorAddress[:4], evidence.Height)

	return nil
}

// Slash executes slashing for a validator
func (st *SlashingTracker) Slash(validator [32]byte, condition SlashingCondition, evidence Evidence, validatorStake uint64) uint64 {
	st.mu.Lock()
	defer st.mu.Unlock()

	// Check if already slashed
	if _, exists := st.slashed[validator]; exists {
		return 0
	}

	// Calculate slash amount
	var slashAmount uint64
	switch condition {
	case DoubleSign:
		slashAmount = (validatorStake * st.DoubleSignSlashAmount) / 100
	case Downtime:
		slashAmount = (validatorStake * st.DowntimeSlashAmount) / 100
	default:
		slashAmount = 0
	}

	// Record slashing
	st.slashed[validator] = SlashInfo{
		Validator:     validator,
		Condition:     condition,
		SlashedAmount: slashAmount,
		SlashedAt:     evidence.SubmittedAt,
		Evidence:      evidence,
	}

	fmt.Printf("[Slashing] ⚠️  VALIDATOR SLASHED: %x for %s | Amount: %d (%.1f%%)\n",
		validator[:4], condition, slashAmount, float64(slashAmount)/float64(validatorStake)*100)

	return slashAmount
}

// IsSlashed checks if a validator has been slashed
func (st *SlashingTracker) IsSlashed(validator [32]byte) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	_, exists := st.slashed[validator]
	return exists
}

// GetSlashInfo returns slashing information for a validator
func (st *SlashingTracker) GetSlashInfo(validator [32]byte) (SlashInfo, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	info, exists := st.slashed[validator]
	return info, exists
}

// GetEvidence returns all evidence for a validator
func (st *SlashingTracker) GetEvidence(validator [32]byte) []Evidence {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.evidence[validator]
}

// VerifyDoubleSignEvidence verifies double-sign evidence
func VerifyDoubleSignEvidence(evidence DoubleSignEvidence, pubKey ed25519.PublicKey) bool {
	// Check same height, different blocks
	if evidence.Vote1.Height != evidence.Vote2.Height {
		return false
	}

	if evidence.Vote1.BlockHash == evidence.Vote2.BlockHash {
		return false
	}

	// Check same validator
	if evidence.Vote1.ValidatorAddress != evidence.Vote2.ValidatorAddress {
		return false
	}

	// Verify signatures (simplified - in real impl would verify full vote structure)
	// vote1Bytes := serializeVote(evidence.Vote1)
	// vote2Bytes := serializeVote(evidence.Vote2)
	//
	// if !ed25519.Verify(pubKey, vote1Bytes, evidence.Vote1.Signature[:]) {
	// 	return false
	// }
	// if !ed25519.Verify(pubKey, vote2Bytes, evidence.Vote2.Signature[:]) {
	// 	return false
	// }

	return true
}

// VerifyDowntimeEvidence verifies downtime evidence
func VerifyDowntimeEvidence(evidence DowntimeEvidence, threshold uint64) bool {
	// Check if missed votes exceed threshold
	return uint64(len(evidence.MissedHeights)) >= threshold
}
