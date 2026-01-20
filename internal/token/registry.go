package token

import (
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Registry stores all created tokens
type Registry struct {
	tokens      map[[32]byte]*types.Token // tokenAddress -> Token
	symbolIndex map[string][32]byte       // symbol -> tokenAddress
	mu          sync.RWMutex
}

// NewRegistry creates a new token registry
func NewRegistry() *Registry {
	return &Registry{
		tokens:      make(map[[32]byte]*types.Token),
		symbolIndex: make(map[string][32]byte),
	}
}

// Register adds a new token to the registry
func (r *Registry) Register(token *types.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if already exists
	if _, exists := r.tokens[token.Address]; exists {
		return fmt.Errorf("token already registered: %x", token.Address)
	}

	// Check symbol uniqueness
	if _, exists := r.symbolIndex[token.Symbol]; exists {
		return fmt.Errorf("token symbol already taken: %s", token.Symbol)
	}

	// Register
	r.tokens[token.Address] = token
	r.symbolIndex[token.Symbol] = token.Address

	return nil
}

// Get retrieves a token by address
func (r *Registry) Get(address [32]byte) (*types.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, exists := r.tokens[address]
	if !exists {
		return nil, fmt.Errorf("token not found: %x", address)
	}

	return token, nil
}

// GetBySymbol retrieves a token by symbol
func (r *Registry) GetBySymbol(symbol string) (*types.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	address, exists := r.symbolIndex[symbol]
	if !exists {
		return nil, fmt.Errorf("token symbol not found: %s", symbol)
	}

	return r.tokens[address], nil
}

// List returns all tokens
func (r *Registry) List() []*types.Token {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tokens := make([]*types.Token, 0, len(r.tokens))
	for _, token := range r.tokens {
		tokens = append(tokens, token)
	}

	return tokens
}

// Count returns number of registered tokens
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tokens)
}

// Exists checks if token exists
func (r *Registry) Exists(address [32]byte) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.tokens[address]
	return exists
}
