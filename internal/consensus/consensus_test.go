package consensus_test

import (
	"testing"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

func TestConsensusAggregator(t *testing.T) {
	agg := consensus.NewAggregator()

	// TODO: Submit 10 proofs and verify aggregation
	if !agg.IsComplete(1) {
		// Expected to be incomplete initially
	}
}

func TestVoting(t *testing.T) {
	vt := consensus.NewVoteTracker()

	// TODO: Submit 7 votes and verify majority
	if vt.HasMajority(1, [32]byte{}) {
		t.Error("Should not have majority with 0 votes")
	}
}

func TestSelectAlgorithm(t *testing.T) {
	seed := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	algo := utils.SelectAlgorithm(seed)

	validAlgorithms := map[string]bool{
		"QUICK_SORT": true,
		"MERGE_SORT": true,
		"HEAP_SORT":  true,
		"RADIX_SORT": true,
		"TIM_SORT":   true,
		"INTRO_SORT": true,
	}

	if !validAlgorithms[algo] {
		t.Errorf("Invalid algorithm selected: %s", algo)
	}
}
