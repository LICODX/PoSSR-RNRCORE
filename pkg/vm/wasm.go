package vm

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// VirtualMachine interface for contract execution
type VirtualMachine interface {
	Deploy(bytecode []byte, creator [32]byte, initArgs []byte, gasLimit uint64) ([32]byte, error)
	Call(contractAddr [32]byte, method string, args []byte, gasLimit uint64) (*types.ContractResult, error)
	Query(contractAddr [32]byte, method string, args []byte) ([]byte, error)
}

// WASMVM implements WebAssembly virtual machine
type WASMVM struct {
	contractState *state.ContractState
}

// NewWASMVM creates a new WASM VM
func NewWASMVM(contractState *state.ContractState) *WASMVM {
	return &WASMVM{
		contractState: contractState,
	}
}

// Deploy deploys a new contract
func (vm *WASMVM) Deploy(bytecode []byte, creator [32]byte, initArgs []byte, gasLimit uint64) ([32]byte, error) {
	// Generate deterministic contract address
	hash := sha256.Sum256(append(creator[:], bytecode...))
	contractAddr := hash

	// Create contract in state
	contract := &types.Contract{
		Address:   contractAddr,
		Creator:   creator,
		Bytecode:  bytecode,
		Balance:   0,
		CreatedAt: time.Now().Unix(),
	}

	// Store contract
	if err := vm.contractState.DeployContract(contract); err != nil {
		return [32]byte{}, fmt.Errorf("failed to store contract: %w", err)
	}

	return contractAddr, nil
}

// Call executes a contract method
func (vm *WASMVM) Call(contractAddr [32]byte, method string, args []byte, gasLimit uint64) (*types.ContractResult, error) {
	// Get contract
	_, err := vm.contractState.GetContract(contractAddr)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	// For now, return success placeholder
	// Full WASM execution with wasmer when we fix executor integration
	return &types.ContractResult{
		Success:    true,
		ReturnData: []byte("contract_executed"),
		GasUsed:    10000,
		Events:     []types.Event{}, // Fixed: use types.Event
		Error:      "",
	}, nil
}

// Query queries contract state (read-only)
func (vm *WASMVM) Query(contractAddr [32]byte, method string, args []byte) ([]byte, error) {
	_, err := vm.contractState.GetContract(contractAddr)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	return []byte("query_result"), nil
}
