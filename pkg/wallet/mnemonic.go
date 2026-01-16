package wallet

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

// MnemonicStrength defines the entropy strength for mnemonic generation
const (
	// MnemonicStrength12Words represents 128 bits of entropy (12 words)
	MnemonicStrength12Words = 128
	// MnemonicStrength24Words represents 256 bits of entropy (24 words)
	MnemonicStrength24Words = 256
)

// GenerateMnemonic creates a new BIP39 mnemonic phrase
// Default strength is 128 bits (12 words)
func GenerateMnemonic() (string, error) {
	return GenerateMnemonicWithStrength(MnemonicStrength12Words)
}

// GenerateMnemonicWithStrength creates a BIP39 mnemonic with custom strength
func GenerateMnemonicWithStrength(strength int) (string, error) {
	// Generate entropy
	entropy, err := generateEntropy(strength)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Convert entropy to mnemonic
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to create mnemonic: %w", err)
	}

	return mnemonic, nil
}

// ValidateMnemonic checks if a mnemonic phrase is valid according to BIP39
func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// MnemonicToSeed converts a mnemonic phrase to a seed
// Password is optional and can be empty string
func MnemonicToSeed(mnemonic, password string) ([]byte, error) {
	if !ValidateMnemonic(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	// BIP39 seed derivation with optional password
	seed := bip39.NewSeed(mnemonic, password)
	return seed, nil
}

// generateEntropy creates cryptographically secure random entropy
func generateEntropy(bitSize int) ([]byte, error) {
	if bitSize%32 != 0 || bitSize < 128 || bitSize > 256 {
		return nil, fmt.Errorf("invalid entropy bit size: must be 128, 160, 192, 224, or 256")
	}

	entropy := make([]byte, bitSize/8)
	_, err := rand.Read(entropy)
	if err != nil {
		return nil, err
	}

	return entropy, nil
}

// WordCount returns the number of words in a mnemonic
func WordCount(mnemonic string) int {
	words := strings.Fields(strings.TrimSpace(mnemonic))
	return len(words)
}
