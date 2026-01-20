package vm

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// ContractExecutor handles WASM contract execution
// This executor is run by ALL validator nodes to ensure consensus
type ContractExecutor struct {
	vm            *WASMVM
	contractState *state.ContractState
	tokenState    *state.TokenState

	// Security protections
	securityConfig   *SecurityConfig
	executionLimiter *ExecutionLimiter
	circuitBreaker   *CircuitBreaker

	// Execution context
	caller       [32]byte
	contractAddr [32]byte
	value        uint64
	gasLimit     uint64
	gasUsed      uint64

	mu sync.RWMutex
}

// NewContractExecutor creates a new contract executor with default security config
func NewContractExecutor(vm *WASMVM, contractState *state.ContractState, tokenState *state.TokenState) *ContractExecutor {
	config := DefaultSecurityConfig()
	return &ContractExecutor{
		vm:             vm,
		contractState:  contractState,
		tokenState:     tokenState,
		securityConfig: config,
		circuitBreaker: NewCircuitBreaker(config),
	}
}

// NewContractExecutorWithConfig creates executor with custom security config
func NewContractExecutorWithConfig(
	vm *WASMVM,
	contractState *state.ContractState,
	tokenState *state.TokenState,
	config *SecurityConfig,
) *ContractExecutor {
	return &ContractExecutor{
		vm:             vm,
		contractState:  contractState,
		tokenState:     tokenState,
		securityConfig: config,
		circuitBreaker: NewCircuitBreaker(config),
	}
}

// DeployContract deploys a new WASM contract
// Called by the node that creates the block, then verified by all other nodes
func (e *ContractExecutor) DeployContract(
	bytecode []byte,
	creator [32]byte,
	initArgs []byte,
	gasLimit uint64,
) ([32]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Use the VM to deploy (VM handles WASM validation internally)
	contractAddr, err := e.vm.Deploy(bytecode, creator, initArgs, gasLimit)
	if err != nil {
		return [32]byte{}, err
	}

	return contractAddr, nil
}

// ExecuteContractWithSecurity executes a contract with full security protections
// This is the main entry point for contract calls during block processing
func (e *ContractExecutor) ExecuteContractWithSecurity(
	contractAddr [32]byte,
	caller [32]byte,
	method string,
	args []byte,
	value uint64,
	gasLimit uint64,
) (*types.ContractResult, error) {
	// Security: Check circuit breaker
	if err := e.circuitBreaker.Allow(); err != nil {
		GlobalSecurityMetrics.RecordViolation("CircuitBreakerOpen")
		return nil, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Set execution context
	e.caller = caller
	e.contractAddr = contractAddr
	e.value = value
	e.gasLimit = gasLimit
	e.gasUsed = 0

	// Execute via VM
	result, err := e.vm.Call(contractAddr, method, args, gasLimit)
	if err != nil {
		e.circuitBreaker.RecordFailure()
		return nil, err
	}

	e.circuitBreaker.RecordSuccess()
	return result, nil
}

// ConsumeGas deducts gas from the available gas limit
func (e *ContractExecutor) ConsumeGas(amount uint64) error {
	e.gasUsed += amount
	if e.gasUsed > e.gasLimit {
		return errors.New("out of gas")
	}
	return nil
}

// GetSecurityMetrics returns current security metrics
func (e *ContractExecutor) GetSecurityMetrics() map[string]uint64 {
	return GlobalSecurityMetrics.GetMetrics()
}

// EmergencyPause manually opens the circuit breaker
func (e *ContractExecutor) EmergencyPause() error {
	if !e.securityConfig.EmergencyPauseEnabled {
		return fmt.Errorf("emergency pause not enabled in config")
	}
	e.circuitBreaker.ManualOpen()
	fmt.Println("⚠️  EMERGENCY PAUSE ACTIVATED - All contract execution stopped")
	return nil
}

// EmergencyResume manually closes the circuit breaker
func (e *ContractExecutor) EmergencyResume() error {
	if !e.securityConfig.EmergencyPauseEnabled {
		return fmt.Errorf("emergency pause not enabled in config")
	}
	e.circuitBreaker.ManualClose()
	fmt.Println("✅ EMERGENCY PAUSE RELEASED - Contract execution resumed")
	return nil
}

// Helper function to convert address to hex string for logging
func addressToHex(addr [32]byte) string {
	return hex.EncodeToString(addr[:])
}

// Helper function to encode uint64 to bytes
func uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

// Helper function to decode bytes to uint64
func bytesToUint64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}
