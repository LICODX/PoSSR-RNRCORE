package state

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/vm"
	"github.com/syndtr/goleveldb/leveldb"
)

// Account represents an account's state
type Account struct {
	Balance uint64
	Nonce   uint64
}

// Manager manages account state, contracts, and tokens
type Manager struct {
	db    *leveldb.DB
	cache map[[32]byte]*Account
	mu    sync.RWMutex

	// Sub-managers
	contractState *ContractState
	tokenState    *TokenState
	executor      *vm.ContractExecutor
}

// NewManager creates a new state manager
func NewManager(db *leveldb.DB) *Manager {
	contractState := NewContractState(db)
	tokenState := NewTokenState(db)
	executor := vm.NewContractExecutor(contractState)

	return &Manager{
		db:            db,
		cache:         make(map[[32]byte]*Account),
		contractState: contractState,
		tokenState:    tokenState,
		executor:      executor,
	}
}

// GetAccount retrieves account state
func (m *Manager) GetAccount(pubkey [32]byte) (*Account, error) {
	m.mu.RLock()
	if acc, ok := m.cache[pubkey]; ok {
		m.mu.RUnlock()
		return acc, nil
	}
	m.mu.RUnlock()

	// Load from DB
	key := append([]byte("account-"), pubkey[:]...)
	data, err := m.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		// New account
		return &Account{Balance: 0, Nonce: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	var acc Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return nil, err
	}

	// Cache it
	m.mu.Lock()
	m.cache[pubkey] = &acc
	m.mu.Unlock()

	return &acc, nil
}

// UpdateAccount saves account state
func (m *Manager) UpdateAccount(pubkey [32]byte, acc *Account) error {
	m.mu.Lock()
	m.cache[pubkey] = acc
	m.mu.Unlock()

	// Persist to DB
	key := append([]byte("account-"), pubkey[:]...)
	data, _ := json.Marshal(acc)
	return m.db.Put(key, data, nil)
}

// ApplyTransaction validates and applies a transaction to state
func (m *Manager) ApplyTransaction(tx types.Transaction) error {
	// Get sender account
	sender, err := m.GetAccount(tx.Sender)
	if err != nil {
		return err
	}

	// Detect Coinbase (Sender = all zeros)
	isCoinbase := true
	for _, b := range tx.Sender {
		if b != 0 {
			isCoinbase = false
			break
		}
	}

	// Skip state checks for Coinbase (System TX)
	if !isCoinbase {
		// Check nonce (replay protection)
		if tx.Nonce != sender.Nonce+1 {
			return fmt.Errorf("invalid nonce: expected %d, got %d", sender.Nonce+1, tx.Nonce)
		}

		// Check balance
		if sender.Balance < tx.Amount {
			return fmt.Errorf("insufficient balance: has %d, needs %d", sender.Balance, tx.Amount)
		}
	}

	// Get receiver account
	receiver, err := m.GetAccount(tx.Receiver)
	if err != nil {
		return err
	}

	// Apply changes
	sender.Balance -= tx.Amount
	sender.Nonce++
	receiver.Balance += tx.Amount

	// Save both accounts
	if err := m.UpdateAccount(tx.Sender, sender); err != nil {
		return err
	}
	if err := m.UpdateAccount(tx.Receiver, receiver); err != nil {
		return err
	}

	return nil
}

// GetContractState returns the contract state manager
func (m *Manager) GetContractState() *ContractState {
	return m.contractState
}

// GetTokenState returns the token state manager
func (m *Manager) GetTokenState() *TokenState {
	return m.tokenState
}

// GetContractExecutor returns the contract executor
func (m *Manager) GetContractExecutor() *vm.ContractExecutor {
	return m.executor
}
