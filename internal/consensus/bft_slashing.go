package consensus

import (
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus/bft"
	"github.com/LICODX/PoSSR-RNRCORE/internal/slashing"
)

// voteCache stores votes per validator to detect double-signing
type voteCache struct {
	prevotes   map[[32]byte]*bft.Vote // validator address -> vote
	precommits map[[32]byte]*bft.Vote // validator address -> vote
}

func newVoteCache() *voteCache {
	return &voteCache{
		prevotes:   make(map[[32]byte]*bft.Vote),
		precommits: make(map[[32]byte]*bft.Vote),
	}
}

// detectDoubleSign checks if a vote conflicts with a previous vote (double-signing)
func (be *BFTEngine) detectDoubleSign(vote *bft.Vote, cache *voteCache) bool {
	var existing *bft.Vote
	var ok bool

	// Get existing vote based on type
	switch vote.Type {
	case bft.VoteTypePrevote:
		existing, ok = cache.prevotes[vote.ValidatorAddress]
	case bft.VoteTypePrecommit:
		existing, ok = cache.precommits[vote.ValidatorAddress]
	default:
		return false
	}

	// If no previous vote, cache this one
	if !ok {
		switch vote.Type {
		case bft.VoteTypePrevote:
			cache.prevotes[vote.ValidatorAddress] = vote
		case bft.VoteTypePrecommit:
			cache.precommits[vote.ValidatorAddress] = vote
		}
		return false
	}

	// Check if same vote (idempotent)
	if existing.BlockHash == vote.BlockHash {
		return false // Not double-sign, just duplicate
	}

	// DOUBLE SIGN DETECTED!
	// Same height, same round, same type, DIFFERENT block hash
	if existing.Height == vote.Height && existing.Round == vote.Round {
		// Submit evidence
		evidence := slashing.Evidence{
			Type:             slashing.DoubleSign,
			ValidatorAddress: vote.ValidatorAddress,
			Height:           vote.Height,
			SubmittedAt:      time.Now().Unix(),
			// Proof would contain both votes serialized
		}

		be.SlashingTracker.SubmitEvidence(evidence)

		// Get validator from set to get stake amount
		val := be.Validators.GetByAddress(vote.ValidatorAddress)
		if val != nil {
			// Slash 100% of stake
			be.SlashingTracker.Slash(vote.ValidatorAddress, slashing.DoubleSign, evidence, val.VotingPower)

			// Remove validator from active set
			be.Validators.Remove(vote.ValidatorAddress)
		}

		return true
	}

	return false
}
