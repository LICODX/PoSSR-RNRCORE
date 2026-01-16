package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tyler-smith/go-bip32"
)

// Derivation paths for RNR
const (
	// CoinType for RNR (999 = custom coin type)
	CoinTypeRNR = 999

	// DefaultDerivationPath is m/44'/999'/0'/0/0
	DefaultDerivationPath = "m/44'/999'/0'/0/0"
)

// HDWallet represents a Hierarchical Deterministic wallet
type HDWallet struct {
	masterKey *bip32.Key
	mnemonic  string
}

// NewHDWalletFromMnemonic creates an HD wallet from a BIP39 mnemonic
func NewHDWalletFromMnemonic(mnemonic, password string) (*HDWallet, error) {
	// Validate mnemonic
	if !ValidateMnemonic(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	// Convert mnemonic to seed
	seed, err := MnemonicToSeed(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate seed: %w", err)
	}

	// Generate master key from seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %w", err)
	}

	return &HDWallet{
		masterKey: masterKey,
		mnemonic:  mnemonic,
	}, nil
}

// DeriveKey derives a child key from the master key using BIP32 derivation path
// Path format: m/44'/999'/0'/0/0
func (hd *HDWallet) DeriveKey(path string) (*Wallet, error) {
	// Parse derivation path
	segments, err := parsePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid derivation path: %w", err)
	}

	// Start with master key
	key := hd.masterKey

	// Derive each level
	for i, segment := range segments {
		key, err = key.NewChildKey(segment)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key at index %d: %w", i, err)
		}
	}

	// Convert BIP32 key to Ed25519 keypair
	// We use the key bytes as seed for Ed25519
	privateKey := ed25519.NewKeyFromSeed(key.Key[:32])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	// Generate Bech32 address
	address, err := PubKeyToAddress(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		Mnemonic:   hd.mnemonic,
		Path:       path,
	}, nil
}

// GetDefaultWallet returns the wallet at the default derivation path
func (hd *HDWallet) GetDefaultWallet() (*Wallet, error) {
	return hd.DeriveKey(DefaultDerivationPath)
}

// parsePath parses a BIP32 derivation path and returns the segments
// Example: "m/44'/999'/0'/0/0" -> [0x8000002C, 0x800003E7, 0x80000000, 0, 0]
func parsePath(path string) ([]uint32, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty derivation path")
	}

	// Remove "m/" prefix if present
	if len(path) >= 2 && path[:2] == "m/" {
		path = path[2:]
	}

	// Split by '/'
	parts := splitPath(path)
	segments := make([]uint32, len(parts))

	for i, part := range parts {
		// Check for hardened key (apostrophe)
		hardened := false
		if len(part) > 0 && part[len(part)-1] == '\'' {
			hardened = true
			part = part[:len(part)-1]
		}

		// Parse number
		var num uint32
		_, err := fmt.Sscanf(part, "%d", &num)
		if err != nil {
			return nil, fmt.Errorf("invalid path segment: %s", part)
		}

		// Add hardened bit if needed
		if hardened {
			num += 0x80000000
		}

		segments[i] = num
	}

	return segments, nil
}

// splitPath splits a path by '/' separator
func splitPath(path string) []string {
	var parts []string
	current := ""

	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
