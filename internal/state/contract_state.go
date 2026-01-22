package state

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/syndtr/goleveldb/leveldb"
)

// ContractState manages smart contract storage
type ContractState struct {
	// In-memory cache: contract -> key -> value
	storage map[[32]byte]map[string][]byte

	// Deployed contracts: address -> Contract
	contracts map[[32]byte]*types.Contract

	mu sync.RWMutex
	db *leveldb.DB
}

// GasMeter tracks computational expenditure during contract execution
type GasMeter struct {
	Limit uint64
	Used  uint64
}

func NewGasMeter(limit uint64) *GasMeter {
	return &GasMeter{Limit: limit, Used: 0}
}

func (gm *GasMeter) Consume(amount uint64) error {
	if gm.Used+amount > gm.Limit {
		return fmt.Errorf("out of gas: limit %d, used %d, requested %d", gm.Limit, gm.Used, amount)
	}
	gm.Used += amount
	return nil
}

// NewContractState creates contract state manager
func NewContractState(db *leveldb.DB) *ContractState {
	return &ContractState{
		storage:   make(map[[32]byte]map[string][]byte),
		contracts: make(map[[32]byte]*types.Contract),
		db:        db,
	}
}

// DeployContract registers a new contract
func (cs *ContractState) DeployContract(contract *types.Contract) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Store in memory
	cs.contracts[contract.Address] = contract
	cs.storage[contract.Address] = make(map[string][]byte)

	// Persist to DB
	key := append([]byte("contract-"), contract.Address[:]...)
	data, _ := json.Marshal(contract)
	return cs.db.Put(key, data, nil)
}

// GetContract retrieves a deployed contract
func (cs *ContractState) GetContract(address [32]byte) (*types.Contract, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	// Check cache
	if contract, exists := cs.contracts[address]; exists {
		return contract, nil
	}

	// Load from DB
	key := append([]byte("contract-"), address[:]...)
	data, err := cs.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var contract types.Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		return nil, err
	}

	// Cache it
	cs.contracts[address] = &contract
	return &contract, nil
}

// StorageRead reads from contract storage
func (cs *ContractState) StorageRead(contract [32]byte, key string) []byte {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	// Check memory cache
	if contractStorage, exists := cs.storage[contract]; exists {
		if value, exists := contractStorage[key]; exists {
			return value
		}
	}

	// Load from DB
	dbKey := append([]byte("storage-"), contract[:]...)
	dbKey = append(dbKey, []byte(key)...)

	data, err := cs.db.Get(dbKey, nil)
	if err != nil {
		return nil
	}

	return data
}

// StorageWrite writes to contract storage
func (cs *ContractState) StorageWrite(contract [32]byte, key string, value []byte) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Update memory cache
	if _, exists := cs.storage[contract]; !exists {
		cs.storage[contract] = make(map[string][]byte)
	}
	cs.storage[contract][key] = value

	// Persist to DB
	dbKey := append([]byte("storage-"), contract[:]...)
	dbKey = append(dbKey, []byte(key)...)

	return cs.db.Put(dbKey, value, nil)
}

// StorageDelete removes from contract storage
func (cs *ContractState) StorageDelete(contract [32]byte, key string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Remove from cache
	if contractStorage, exists := cs.storage[contract]; exists {
		delete(contractStorage, key)
	}

	// Delete from DB
	dbKey := append([]byte("storage-"), contract[:]...)
	dbKey = append(dbKey, []byte(key)...)

	return cs.db.Delete(dbKey, nil)
}

// UpdateContractBalance updates contract's RNR balance
func (cs *ContractState) UpdateContractBalance(address [32]byte, balance uint64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	contract, exists := cs.contracts[address]
	if !exists {
		// Load from DB
		var err error
		contract, err = cs.GetContract(address)
		if err != nil {
			return err
		}
	}

	contract.Balance = balance

	// Persist
	key := append([]byte("contract-"), address[:]...)
	data, _ := json.Marshal(contract)
	return cs.db.Put(key, data, nil)
}

// GetContractBalance returns contract's RNR balance
func (cs *ContractState) GetContractBalance(address [32]byte) uint64 {
	contract, err := cs.GetContract(address)
	if err != nil {
		return 0
	}
	return contract.Balance
}

// ListContracts returns all deployed contracts
func (cs *ContractState) ListContracts() []*types.Contract {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	contracts := make([]*types.Contract, 0, len(cs.contracts))
	for _, contract := range cs.contracts {
		contracts = append(contracts, contract)
	}
	return contracts
}

// ExecuteContract simulates contract execution with gas metering
func (cs *ContractState) ExecuteContract(address [32]byte, payload []byte, gasLimit uint64) (uint64, error) {
	meter := NewGasMeter(gasLimit)

	// Sub-PHASE 1: Load Contract
	_, err := cs.GetContract(address)
	if err != nil {
		return 0, err
	}
	meter.Consume(100) // Base execution cost

	// Sub-PHASE 2: Process Payload (Logic depends on Contract Type)
	// Example: Token Transfer within contract
	// For now, we simulate complexity based on payload length
	complexity := uint64(len(payload)) * 10
	if err := meter.Consume(complexity); err != nil {
		return meter.Used, err
	}

	// Persist changes if successful...
	return meter.Used, nil
}
