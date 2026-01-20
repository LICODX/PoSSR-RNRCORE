package utils

import (
	"crypto/sha256"
)

// SelectAlgorithm uses VRF Seed to determine sorting algorithm
func SelectAlgorithm(seed [32]byte) string {
	selector := seed[31] % 7 // Increased to 7 algorithms
	switch selector {
	case 0:
		return "QUICK_SORT"
	case 1:
		return "MERGE_SORT"
	case 2:
		return "HEAP_SORT"
	case 3:
		return "RADIX_SORT"
	case 4:
		return "TIM_SORT"
	case 5:
		return "INTRO_SORT"
	case 6:
		return "SHELL_SORT"
	default:
		return "QUICK_SORT"
	}
}

// MixHash combines ID and seed to create a sorting key
func MixHash(id [32]byte, seed [32]byte) string {
	h := sha256.New()
	h.Write(id[:])
	h.Write(seed[:])
	return string(h.Sum(nil))
}
