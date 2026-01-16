package economics

import (
	"math"

	"github.com/LICODX/PoSSR-RNRCORE/internal/params"
)

// GetBlockReward menghitung reward total untuk 1 blok (untuk dibagi ke 10 node)
func GetBlockReward(height uint64) float64 {
	// Tentukan fase halving (integer division)
	phase := float64(height / params.HalvingInterval)

	// Rumus Decay: Initial * (0.93 ^ Phase)
	decayFactor := math.Pow(1.0-params.DecayRate, phase)

	totalReward := params.InitialReward * decayFactor

	return totalReward
}

// CalculateNodeShare menghitung jatah per node
func CalculateNodeShare(totalReward float64) float64 {
	return totalReward / float64(params.NumShards)
}
