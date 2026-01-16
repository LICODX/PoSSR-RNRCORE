package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
)

// GenerateKeypair creates a new ED25519 keypair
func GenerateKeypair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(rand.Reader)
}

// Sign signs a message using the private key
func Sign(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}

// Verify checks if the signature is valid for the message and public key
func Verify(publicKey ed25519.PublicKey, message []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}

// Hash256 calculates the Double-SHA256 hash of data (Bitcoin style)
func Hash256(data []byte) [32]byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2
}

// Hash calculates the Single-SHA256 hash of data
func Hash(data []byte) [32]byte {
	return sha256.Sum256(data)
}
