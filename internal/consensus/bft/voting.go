package bft

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// VoteType represents the type of vote
type VoteType byte

const (
	VoteTypePrevote   VoteType = 0x01 // First voting phase
	VoteTypePrecommit VoteType = 0x02 // Second voting phase (commitment)
)

// Vote represents a validator's vote in consensus
type Vote struct {
	Type      VoteType // Prevote or Precommit
	Height    uint64   // Block height being voted on
	Round     int32    // Consensus round number
	BlockHash [32]byte // Hash of the block being voted for (nil = vote for nil)
	Timestamp int64    // Unix timestamp of vote

	ValidatorAddress [32]byte // Address of validator
	ValidatorIndex   int32    // Index in validator set (for fast lookup)
	Signature        [64]byte // Ed25519 signature
}

// SignBytes returns the canonical byte representation for signing
func (vote *Vote) SignBytes() []byte {
	var buf bytes.Buffer

	buf.WriteByte(byte(vote.Type))
	binary.Write(&buf, binary.BigEndian, vote.Height)
	binary.Write(&buf, binary.BigEndian, vote.Round)
	buf.Write(vote.BlockHash[:])
	binary.Write(&buf, binary.BigEndian, vote.Timestamp)

	return buf.Bytes()
}

// Sign signs the vote with a private key
func (vote *Vote) Sign(privKey ed25519.PrivateKey) {
	signBytes := vote.SignBytes()
	sig := ed25519.Sign(privKey, signBytes)
	copy(vote.Signature[:], sig)
}

// Verify verifies the vote signature
func (vote *Vote) Verify(pubKey ed25519.PublicKey) bool {
	signBytes := vote.SignBytes()
	return ed25519.Verify(pubKey, signBytes, vote.Signature[:])
}

// IsNil returns true if this is a vote for nil
func (vote *Vote) IsNil() bool {
	return vote.BlockHash == [32]byte{}
}

// VoteSet manages a set of votes for a specific height/round
type VoteSet struct {
	Height   uint64
	Round    int32
	VoteType VoteType
	valSet   *ValidatorSet

	votes            map[int32]*Vote     // validatorIndex -> vote
	votesByBlock     map[[32]byte]uint64 // blockHash -> total voting power
	totalVotingPower uint64
}

// NewVoteSet creates a new vote set
func NewVoteSet(height uint64, round int32, voteType VoteType, valSet *ValidatorSet) *VoteSet {
	return &VoteSet{
		Height:           height,
		Round:            round,
		VoteType:         voteType,
		valSet:           valSet,
		votes:            make(map[int32]*Vote),
		votesByBlock:     make(map[[32]byte]uint64),
		totalVotingPower: valSet.TotalVotingPower(),
	}
}

// AddVote adds a vote to the set
func (voteSet *VoteSet) AddVote(vote *Vote) (bool, error) {
	// Validate vote
	if vote.Height != voteSet.Height {
		return false, fmt.Errorf("vote height mismatch: %d != %d", vote.Height, voteSet.Height)
	}
	if vote.Round != voteSet.Round {
		return false, fmt.Errorf("vote round mismatch: %d != %d", vote.Round, voteSet.Round)
	}
	if vote.Type != voteSet.VoteType {
		return false, fmt.Errorf("vote type mismatch: %d != %d", vote.Type, voteSet.VoteType)
	}

	// Get validator
	val := voteSet.valSet.GetByAddress(vote.ValidatorAddress)
	if val == nil {
		return false, fmt.Errorf("validator not in set: %x", vote.ValidatorAddress[:4])
	}

	// Verify signature
	if !vote.Verify(val.PubKey) {
		return false, fmt.Errorf("invalid vote signature")
	}

	// Check if already voted
	if existing, ok := voteSet.votes[vote.ValidatorIndex]; ok {
		// If same vote, ignore (idempotent)
		if existing.BlockHash == vote.BlockHash {
			return false, nil
		}
		// Double voting detected! This is slashable
		return false, fmt.Errorf("double vote detected from validator %x", vote.ValidatorAddress[:4])
	}

	// Add vote
	voteSet.votes[vote.ValidatorIndex] = vote
	voteSet.votesByBlock[vote.BlockHash] += val.VotingPower

	return true, nil
}

// HasTwoThirdsMajority checks if any block has 2/3+ votes
func (voteSet *VoteSet) HasTwoThirdsMajority() (bool, [32]byte) {
	for blockHash, votingPower := range voteSet.votesByBlock {
		if voteSet.valSet.HasTwoThirdsMajority(votingPower) {
			return true, blockHash
		}
	}
	return false, [32]byte{}
}

// HasTwoThirdsAny checks if 2/3+ validators have voted (regardless of block)
func (voteSet *VoteSet) HasTwoThirdsAny() bool {
	totalVoted := uint64(0)
	for _, vote := range voteSet.votes {
		val := voteSet.valSet.GetByAddress(vote.ValidatorAddress)
		if val != nil {
			totalVoted += val.VotingPower
		}
	}
	return voteSet.valSet.HasTwoThirdsMajority(totalVoted)
}

// GetVotingPowerFor returns voting power for a specific block
func (voteSet *VoteSet) GetVotingPowerFor(blockHash [32]byte) uint64 {
	return voteSet.votesByBlock[blockHash]
}

// Size returns number of votes
func (voteSet *VoteSet) Size() int {
	return len(voteSet.votes)
}

// HeightRoundKey generates a unique key for height/round combination
func HeightRoundKey(height uint64, round int32) [32]byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, height)
	binary.Write(&buf, binary.BigEndian, round)
	return sha256.Sum256(buf.Bytes())
}
