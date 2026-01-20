package types

// Transaction type constants
const (
	TxTypeRNRTransfer   = 0 // Normal RNR transfer
	TxTypeTokenCreate   = 1 // Create new token
	TxTypeTokenTransfer = 2 // Transfer tokens
	TxTypeTokenApprove  = 3 // Approve spending allowance
	TxTypeTokenMint     = 4 // Mint new tokens (if mintable)
	TxTypeTokenBurn     = 5 // Burn tokens (if burnable)
)

// Token operation payloads (JSON encoded in Transaction.Payload)

// TokenCreatePayload for creating new tokens
type TokenCreatePayload struct {
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	Decimals      uint8  `json:"decimals"`
	InitialSupply uint64 `json:"initialSupply"`
	IsMintable    bool   `json:"isMintable"`
	IsBurnable    bool   `json:"isBurnable"`
}

// TokenTransferPayload for token transfers
type TokenTransferPayload struct {
	TokenAddress [32]byte `json:"tokenAddress"`
	To           [32]byte `json:"to"`
	Amount       uint64   `json:"amount"`
}

// TokenApprovePayload for approving allowances
type TokenApprovePayload struct {
	TokenAddress [32]byte `json:"tokenAddress"`
	Spender      [32]byte `json:"spender"`
	Amount       uint64   `json:"amount"`
}

// TokenMintPayload for minting tokens
type TokenMintPayload struct {
	TokenAddress [32]byte `json:"tokenAddress"`
	To           [32]byte `json:"to"`
	Amount       uint64   `json:"amount"`
}

// TokenBurnPayload for burning tokens
type TokenBurnPayload struct {
	TokenAddress [32]byte `json:"tokenAddress"`
	Amount       uint64   `json:"amount"`
}
