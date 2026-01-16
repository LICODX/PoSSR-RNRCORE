package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

// KeyStore handles encrypted key storage
type KeyStore struct {
	FilePath string
}

// EncryptedKey represents encrypted wallet data
type EncryptedKey struct {
	Address    string `json:"address"`
	Ciphertext []byte `json:"ciphertext"`
	Salt       []byte `json:"salt"`
	Nonce      []byte `json:"nonce"`
}

// NewKeyStore creates a new keystore
func NewKeyStore(path string) *KeyStore {
	return &KeyStore{FilePath: path}
}

// Save encrypts and saves wallet to disk
func (ks *KeyStore) Save(w *Wallet, password string) error {
	// Derive encryption key from password
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return err
	}

	// Encrypt private key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, w.PrivateKey, nil)

	// Create encrypted key structure
	ek := EncryptedKey{
		Address:    w.Address,
		Ciphertext: ciphertext,
		Salt:       salt,
		Nonce:      nonce,
	}

	// Save to file
	data, err := json.Marshal(ek)
	if err != nil {
		return err
	}

	return os.WriteFile(ks.FilePath, data, 0600)
}

// Load decrypts and loads wallet from disk
func (ks *KeyStore) Load(password string) (*Wallet, error) {
	// Read file
	data, err := os.ReadFile(ks.FilePath)
	if err != nil {
		return nil, err
	}

	var ek EncryptedKey
	if err := json.Unmarshal(data, &ek); err != nil {
		return nil, err
	}

	// Derive decryption key
	key, err := scrypt.Key([]byte(password), ek.Salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// Decrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	privateKey, err := gcm.Open(nil, ek.Nonce, ek.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// Reconstruct wallet
	return ImportPrivateKey(string(privateKey))
}

// DeriveAddress derives public address from private key
func DeriveAddress(privateKey []byte) string {
	hash := sha256.Sum256(privateKey)
	return string(hash[:20]) // Simplified - use first 20 bytes
}
