package state

import (
	"encoding/json"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

// TokenState manages token balances and allowances
type TokenState struct {
	// In-memory cache
	balances   map[[32]byte]map[[32]byte]uint64              // token -> account -> balance
	allowances map[[32]byte]map[[32]byte]map[[32]byte]uint64 // token -> owner -> spender -> amount
	mu         sync.RWMutex

	// Persistent storage
	db *leveldb.DB
}

// NewTokenState creates a new token state manager
func NewTokenState(db *leveldb.DB) *TokenState {
	return &TokenState{
		balances:   make(map[[32]byte]map[[32]byte]uint64),
		allowances: make(map[[32]byte]map[[32]byte]map[[32]byte]uint64),
		db:         db,
	}
}

// GetBalance returns token balance for account
func (ts *TokenState) GetBalance(tokenAddr, account [32]byte) uint64 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Check cache first
	if tokenBalances, exists := ts.balances[tokenAddr]; exists {
		if balance, exists := tokenBalances[account]; exists {
			return balance
		}
	}

	// Load from DB
	key := append([]byte("token-balance-"), tokenAddr[:]...)
	key = append(key, account[:]...)

	data, err := ts.db.Get(key, nil)
	if err != nil {
		return 0
	}

	var balance uint64
	json.Unmarshal(data, &balance)
	return balance
}

// SetBalance sets token balance for account
func (ts *TokenState) SetBalance(tokenAddr, account [32]byte, balance uint64) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Update cache
	if _, exists := ts.balances[tokenAddr]; !exists {
		ts.balances[tokenAddr] = make(map[[32]byte]uint64)
	}
	ts.balances[tokenAddr][account] = balance

	// Persist to DB
	key := append([]byte("token-balance-"), tokenAddr[:]...)
	key = append(key, account[:]...)
	data, _ := json.Marshal(balance)
	ts.db.Put(key, data, nil)
}

// GetAllowance returns spending allowance
func (ts *TokenState) GetAllowance(tokenAddr, owner, spender [32]byte) uint64 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Check cache
	if tokenAllowances, exists := ts.allowances[tokenAddr]; exists {
		if ownerAllowances, exists := tokenAllowances[owner]; exists {
			if allowance, exists := ownerAllowances[spender]; exists {
				return allowance
			}
		}
	}

	// Load from DB
	key := append([]byte("token-allowance-"), tokenAddr[:]...)
	key = append(key, owner[:]...)
	key = append(key, spender[:]...)

	data, err := ts.db.Get(key, nil)
	if err != nil {
		return 0
	}

	var allowance uint64
	json.Unmarshal(data, &allowance)
	return allowance
}

// SetAllowance sets spending allowance
func (ts *TokenState) SetAllowance(tokenAddr, owner, spender [32]byte, amount uint64) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Update cache
	if _, exists := ts.allowances[tokenAddr]; !exists {
		ts.allowances[tokenAddr] = make(map[[32]byte]map[[32]byte]uint64)
	}
	if _, exists := ts.allowances[tokenAddr][owner]; !exists {
		ts.allowances[tokenAddr][owner] = make(map[[32]byte]uint64)
	}
	ts.allowances[tokenAddr][owner][spender] = amount

	// Persist to DB
	key := append([]byte("token-allowance-"), tokenAddr[:]...)
	key = append(key, owner[:]...)
	key = append(key, spender[:]...)
	data, _ := json.Marshal(amount)
	ts.db.Put(key, data, nil)
}

// GetAllBalances returns all token balances for an account
func (ts *TokenState) GetAllBalances(account [32]byte) map[[32]byte]uint64 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make(map[[32]byte]uint64)
	for tokenAddr, tokenBalances := range ts.balances {
		if balance, exists := tokenBalances[account]; exists && balance > 0 {
			result[tokenAddr] = balance
		}
	}

	return result
}
