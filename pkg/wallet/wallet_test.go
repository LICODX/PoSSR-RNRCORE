package wallet

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	// Should have 12 words
	wordCount := WordCount(mnemonic)
	if wordCount != 12 {
		t.Errorf("Expected 12 words, got %d", wordCount)
	}

	// Should be valid
	if !ValidateMnemonic(mnemonic) {
		t.Error("Generated mnemonic is not valid")
	}

	t.Logf("Generated mnemonic: %s", mnemonic)
}

func TestValidateMnemonic(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		want     bool
	}{
		{
			name:     "valid 12-word mnemonic",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			want:     true,
		},
		{
			name:     "invalid mnemonic - wrong word",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon invalidword",
			want:     false,
		},
		{
			name:     "invalid mnemonic - too few words",
			mnemonic: "abandon abandon abandon",
			want:     false,
		},
		{
			name:     "empty mnemonic",
			mnemonic: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateMnemonic(tt.mnemonic)
			if got != tt.want {
				t.Errorf("ValidateMnemonic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMnemonicToSeed(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	// Test without password
	seed1, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to convert mnemonic to seed: %v", err)
	}

	// Seed should be 64 bytes (512 bits)
	if len(seed1) != 64 {
		t.Errorf("Expected seed length 64, got %d", len(seed1))
	}

	// Test with password
	seed2, err := MnemonicToSeed(mnemonic, "mypassword")
	if err != nil {
		t.Fatalf("Failed to convert mnemonic to seed with password: %v", err)
	}

	// Seeds should be different with different passwords
	if hex.EncodeToString(seed1) == hex.EncodeToString(seed2) {
		t.Error("Seeds should differ when using different passwords")
	}

	// Test invalid mnemonic
	_, err = MnemonicToSeed("invalid mnemonic phrase", "")
	if err == nil {
		t.Error("Expected error for invalid mnemonic")
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		want     int
	}{
		{
			name:     "12 words",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			want:     12,
		},
		{
			name:     "1 word",
			mnemonic: "word",
			want:     1,
		},
		{
			name:     "empty",
			mnemonic: "",
			want:     0,
		},
		{
			name:     "with extra spaces",
			mnemonic: "  word1   word2  word3  ",
			want:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordCount(tt.mnemonic)
			if got != tt.want {
				t.Errorf("WordCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateWallet(t *testing.T) {
	wallet, err := CreateWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Check mnemonic
	if wallet.Mnemonic == "" {
		t.Error("Wallet mnemonic is empty")
	}

	if WordCount(wallet.Mnemonic) != 12 {
		t.Errorf("Expected 12-word mnemonic, got %d words", WordCount(wallet.Mnemonic))
	}

	// Check address format
	if !strings.HasPrefix(wallet.Address, "rnr1") {
		t.Errorf("Address should start with 'rnr1', got: %s", wallet.Address)
	}

	// Check derivation path
	expectedPath := "m/44'/999'/0'/0/0"
	if wallet.Path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, wallet.Path)
	}

	t.Logf("Created wallet:")
	t.Logf("  Mnemonic: %s", wallet.Mnemonic)
	t.Logf("  Address: %s", wallet.Address)
	t.Logf("  Path: %s", wallet.Path)
}

func TestCreateWalletFromMnemonic(t *testing.T) {
	// Use known mnemonic
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	wallet1, err := CreateWalletFromMnemonic(mnemonic)
	if err != nil {
		t.Fatalf("Failed to create wallet from mnemonic: %v", err)
	}

	// Create again with same mnemonic - should produce same address
	wallet2, err := CreateWalletFromMnemonic(mnemonic)
	if err != nil {
		t.Fatalf("Failed to create wallet from mnemonic (2nd time): %v", err)
	}

	if wallet1.Address != wallet2.Address {
		t.Error("Same mnemonic should produce same address")
	}

	if hex.EncodeToString(wallet1.PrivateKey) != hex.EncodeToString(wallet2.PrivateKey) {
		t.Error("Same mnemonic should produce same private key")
	}

	t.Logf("Deterministic wallet:")
	t.Logf("  Mnemonic: %s", mnemonic)
	t.Logf("  Address: %s", wallet1.Address)
}
