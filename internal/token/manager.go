package token

import (
	"fmt"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Manager handles token operations
type Manager struct {
	registry   *Registry
	tokenState *state.TokenState
}

// NewManager creates a new token manager
func NewManager(registry *Registry, tokenState *state.TokenState) *Manager {
	return &Manager{
		registry:   registry,
		tokenState: tokenState,
	}
}

// CreateToken creates a new RNR-20 token
func (m *Manager) CreateToken(metadata types.TokenMetadata, creator [32]byte) (*types.Token, error) {
	// Validate metadata
	if metadata.Name == "" || metadata.Symbol == "" {
		return nil, fmt.Errorf("name and symbol required")
	}
	if len(metadata.Symbol) < 2 || len(metadata.Symbol) > 6 {
		return nil, fmt.Errorf("symbol must be 2-6 characters")
	}
	if metadata.Decimals > 18 {
		return nil, fmt.Errorf("decimals cannot exceed 18")
	}

	// Generate unique token address
	tokenAddress := types.GenerateTokenAddress(creator, metadata.Name, metadata.Symbol)

	// Create token
	token := &types.Token{
		Address:     tokenAddress,
		Name:        metadata.Name,
		Symbol:      metadata.Symbol,
		Decimals:    metadata.Decimals,
		TotalSupply: metadata.InitialSupply,
		Creator:     creator,
		CreatedAt:   time.Now().Unix(),
		IsMintable:  metadata.IsMintable,
		IsBurnable:  metadata.IsBurnable,
		IsPaused:    false,
	}

	// Register token
	if err := m.registry.Register(token); err != nil {
		return nil, err
	}

	// Set initial balance for creator
	if metadata.InitialSupply > 0 {
		m.tokenState.SetBalance(tokenAddress, creator, metadata.InitialSupply)
	}

	return token, nil
}

// Transfer transfers tokens from one account to another
func (m *Manager) Transfer(tokenAddr, from, to [32]byte, amount uint64) error {
	// Get token
	token, err := m.registry.Get(tokenAddr)
	if err != nil {
		return err
	}

	// Check if paused
	if token.IsPaused {
		return fmt.Errorf("token transfers are paused")
	}

	// Get balances
	fromBalance := m.tokenState.GetBalance(tokenAddr, from)
	if fromBalance < amount {
		return fmt.Errorf("insufficient balance: have %d, need %d", fromBalance, amount)
	}

	// Update balances
	m.tokenState.SetBalance(tokenAddr, from, fromBalance-amount)
	toBalance := m.tokenState.GetBalance(tokenAddr, to)
	m.tokenState.SetBalance(tokenAddr, to, toBalance+amount)

	return nil
}

// GetBalance returns token balance for account
func (m *Manager) GetBalance(tokenAddr, account [32]byte) uint64 {
	return m.tokenState.GetBalance(tokenAddr, account)
}

// Approve sets allowance for spender
func (m *Manager) Approve(tokenAddr, owner, spender [32]byte, amount uint64) error {
	// Verify token exists
	if !m.registry.Exists(tokenAddr) {
		return fmt.Errorf("token not found")
	}

	// Set allowance
	m.tokenState.SetAllowance(tokenAddr, owner, spender, amount)
	return nil
}

// TransferFrom transfers tokens using allowance
func (m *Manager) TransferFrom(tokenAddr, spender, from, to [32]byte, amount uint64) error {
	// Check allowance
	allowance := m.tokenState.GetAllowance(tokenAddr, from, spender)
	if allowance < amount {
		return fmt.Errorf("insufficient allowance: have %d, need %d", allowance, amount)
	}

	// Transfer
	if err := m.Transfer(tokenAddr, from, to, amount); err != nil {
		return err
	}

	// Decrease allowance
	m.tokenState.SetAllowance(tokenAddr, from, spender, allowance-amount)
	return nil
}

// Mint creates new tokens (only if mintable)
func (m *Manager) Mint(tokenAddr, to [32]byte, amount uint64, minter [32]byte) error {
	token, err := m.registry.Get(tokenAddr)
	if err != nil {
		return err
	}

	// Check if mintable
	if !token.IsMintable {
		return fmt.Errorf("token is not mintable")
	}

	// Check if minter is creator
	if minter != token.Creator {
		return fmt.Errorf("only creator can mint tokens")
	}

	// Mint tokens
	balance := m.tokenState.GetBalance(tokenAddr, to)
	m.tokenState.SetBalance(tokenAddr, to, balance+amount)
	token.TotalSupply += amount

	return nil
}

// Burn destroys tokens (only if burnable)
func (m *Manager) Burn(tokenAddr, from [32]byte, amount uint64) error {
	token, err := m.registry.Get(tokenAddr)
	if err != nil {
		return err
	}

	// Check if burnable
	if !token.IsBurnable {
		return fmt.Errorf("token is not burnable")
	}

	// Check balance
	balance := m.tokenState.GetBalance(tokenAddr, from)
	if balance < amount {
		return fmt.Errorf("insufficient balance to burn")
	}

	// Burn tokens
	m.tokenState.SetBalance(tokenAddr, from, balance-amount)
	token.TotalSupply -= amount

	return nil
}
