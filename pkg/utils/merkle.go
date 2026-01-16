package utils

import (
	"crypto/sha256"
)

// CalculateMerkleRoot computes the Merkle Root of a list of 32-byte hashes
func CalculateMerkleRoot(hashes [][32]byte) [32]byte {
	if len(hashes) == 0 {
		return [32]byte{}
	}

	// Copy hashes to avoid modifying the original slice
	currentLevel := make([][32]byte, len(hashes))
	copy(currentLevel, hashes)

	for len(currentLevel) > 1 {
		// If odd number of hashes, duplicate the last one
		if len(currentLevel)%2 != 0 {
			currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
		}

		var nextLevel [][32]byte
		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			right := currentLevel[i+1]

			// Hash(Left + Right)
			combined := append(left[:], right[:]...)
			hash := sha256.Sum256(combined)
			nextLevel = append(nextLevel, hash)
		}
		currentLevel = nextLevel
	}

	return currentLevel[0]
}
