package wallet

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Wallet represents an HD wallet
type Wallet struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Address    string // Bech32 format: rnr1...
	Mnemonic   string // BIP39 mnemonic (12 words)
	Path       string // BIP32 derivation path
}

// CreateWallet generates a new wallet with BIP39 mnemonic
func CreateWallet() (*Wallet, error) {
	// Generate 12-word mnemonic
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Create HD wallet from mnemonic (no password)
	hdWallet, err := NewHDWalletFromMnemonic(mnemonic, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create HD wallet: %w", err)
	}

	// Derive default wallet (m/44'/999'/0'/0/0)
	wallet, err := hdWallet.GetDefaultWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to derive wallet: %w", err)
	}

	return wallet, nil
}

// CreateWalletFromMnemonic creates a wallet from an existing BIP39 mnemonic
func CreateWalletFromMnemonic(mnemonic string) (*Wallet, error) {
	return CreateWalletFromMnemonicWithPassword(mnemonic, "")
}

// CreateWalletFromMnemonicWithPassword creates a wallet from mnemonic with optional password
func CreateWalletFromMnemonicWithPassword(mnemonic, password string) (*Wallet, error) {
	// Create HD wallet
	hdWallet, err := NewHDWalletFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create HD wallet: %w", err)
	}

	// Derive default wallet
	wallet, err := hdWallet.GetDefaultWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to derive wallet: %w", err)
	}

	return wallet, nil
}

// ImportPrivateKey creates wallet from existing key (legacy support)
func ImportPrivateKey(hexKey string) (*Wallet, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	priv := ed25519.PrivateKey(keyBytes)
	pub := priv.Public().(ed25519.PublicKey)

	// Generate Bech32 address
	address, err := PubKeyToAddress(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Wallet{
		PrivateKey: priv,
		PublicKey:  pub,
		Address:    address,
		Mnemonic:   "", // Legacy import has no mnemonic
		Path:       "",
	}, nil
}

// SignTransaction signs a transaction
func (w *Wallet) SignTransaction(tx *types.Transaction) error {
	// Serialize transaction
	message := types.SerializeTransaction(*tx)

	// Sign
	sig := ed25519.Sign(w.PrivateKey, message)
	copy(tx.Signature[:], sig)

	return nil
}

// CreateTransaction creates and signs a new transaction
func (w *Wallet) CreateTransaction(to string, amount uint64, nonce uint64) (*types.Transaction, error) {
	// Decode receiver address
	toBytes, err := hex.DecodeString(to)
	if err != nil {
		return nil, err
	}

	var sender, receiver [32]byte
	copy(sender[:], w.PublicKey)
	copy(receiver[:], toBytes)

	tx := &types.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
		Nonce:    nonce,
	}

	// Calculate TX ID
	tx.ID = types.HashTransaction(*tx)

	// Sign
	if err := w.SignTransaction(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// ExportPrivateKey exports private key as hex
func (w *Wallet) ExportPrivateKey() string {
	return hex.EncodeToString(w.PrivateKey)
}

// GetBalance queries balance (requires state manager access)
func (w *Wallet) GetBalance() (uint64, error) {
	// This would query the state manager
	// For now, return mock
	return 1000, fmt.Errorf("balance query not implemented - requires RPC client")
}
