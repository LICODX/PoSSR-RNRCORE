package wallet

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/bech32"
)

// AddressPrefix is the human-readable part for RNR addresses
const AddressPrefix = "rnr"

// Bech32Encoding specifies the encoding version
// We use Bech32 (not Bech32m) for compatibility
const Bech32Encoding = bech32.Version0

// PubKeyToAddress converts an Ed25519 public key to a Bech32-encoded address
// Format: rnr1qp... (similar to Bitcoin's bc1q...)
func PubKeyToAddress(pubKey ed25519.PublicKey) (string, error) {
	if len(pubKey) != ed25519.PublicKeySize {
		return "", fmt.Errorf("invalid public key size: got %d, expected %d", len(pubKey), ed25519.PublicKeySize)
	}

	// Hash the public key with SHA256
	// This creates a 32-byte hash from the 32-byte public key
	hash := sha256.Sum256(pubKey)

	// Convert to 5-bit groups for Bech32 encoding
	conv, err := bech32.ConvertBits(hash[:], 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits: %w", err)
	}

	// Encode with Bech32
	address, err := bech32.EncodeFromBase256(AddressPrefix, conv)
	if err != nil {
		return "", fmt.Errorf("failed to encode address: %w", err)
	}

	return address, nil
}

// AddressToHash decodes a Bech32 address back to the public key hash
func AddressToHash(address string) ([]byte, error) {
	// Decode Bech32
	hrp, data, err := bech32.DecodeToBase256(address)
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}

	// Verify prefix
	if hrp != AddressPrefix {
		return nil, fmt.Errorf("invalid address prefix: got %s, expected %s", hrp, AddressPrefix)
	}

	// Verify hash size (should be 32 bytes)
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid address hash size: got %d, expected 32", len(data))
	}

	return data, nil
}

// ValidateAddress checks if an address is valid Bech32 with correct prefix
func ValidateAddress(address string) bool {
	_, err := AddressToHash(address)
	return err == nil
}

// IsValidAddressFormat performs a quick format check without full decoding
func IsValidAddressFormat(address string) bool {
	// Check minimum length (prefix + separator + some data)
	if len(address) < len(AddressPrefix)+2 {
		return false
	}

	// Check prefix
	if len(address) < len(AddressPrefix)+1 || address[:len(AddressPrefix)] != AddressPrefix {
		return false
	}

	// Check separator
	if address[len(AddressPrefix)] != '1' {
		return false
	}

	return true
}

// GetAddressPrefix returns the current address prefix
func GetAddressPrefix() string {
	return AddressPrefix
}
