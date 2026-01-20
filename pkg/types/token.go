package types

import (
	"crypto/sha256"
	"time"
)

// Token represents an RNR-20 token (similar to ERC-20)
type Token struct {
	// Unique identifier
	Address [32]byte `json:"address"`

	// Metadata
	Name     string `json:"name"`     // e.g., "MyToken"
	Symbol   string `json:"symbol"`   // e.g., "MTK"
	Decimals uint8  `json:"decimals"` // Typically 6 or 18

	// Supply
	TotalSupply uint64 `json:"totalSupply"` // Total tokens minted

	// Creator
	Creator   [32]byte `json:"creator"`   // Token creator address
	CreatedAt int64    `json:"createdAt"` // Unix timestamp

	// Capabilities
	IsMintable bool `json:"isMintable"` // Can mint more tokens
	IsBurnable bool `json:"isBurnable"` // Can burn tokens
	IsPaused   bool `json:"isPaused"`   // Pause/unpause transfers
}

// TokenBalance represents a token balance for an account
type TokenBalance struct {
	TokenAddress [32]byte `json:"tokenAddress"`
	Account      [32]byte `json:"account"`
	Balance      uint64   `json:"balance"`
}

// TokenAllowance represents spending allowance
type TokenAllowance struct {
	TokenAddress [32]byte `json:"tokenAddress"`
	Owner        [32]byte `json:"owner"`
	Spender      [32]byte `json:"spender"`
	Amount       uint64   `json:"amount"`
}

// TokenMetadata for token creation
type TokenMetadata struct {
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	Decimals      uint8  `json:"decimals"`
	InitialSupply uint64 `json:"initialSupply"`
	IsMintable    bool   `json:"isMintable"`
	IsBurnable    bool   `json:"isBurnable"`
}

// GenerateTokenAddress creates unique token address
func GenerateTokenAddress(creator [32]byte, name, symbol string) [32]byte {
	// Hash creator + name + symbol + timestamp for uniqueness
	data := append(creator[:], []byte(name+symbol)...)
	data = append(data, []byte(time.Now().String())...)
	return sha256.Sum256(data)
}
