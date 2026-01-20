package blockchain

import (
	"encoding/json"
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/vm"
)

// ContractProcessor handles smart contract transaction processing
// This ensures all nodes execute contracts for consensus
type ContractProcessor struct {
	executor      *vm.ContractExecutor
	contractState *state.ContractState
}

// NewContractProcessor creates a new contract processor
func NewContractProcessor(executor *vm.ContractExecutor, contractState *state.ContractState) *ContractProcessor {
	return &ContractProcessor{
		executor:      executor,
		contractState: contractState,
	}
}

// ProcessContractTransaction processes contract deploy and call transactions
// This function is called by ALL nodes when processing a block
func (cp *ContractProcessor) ProcessContractTransaction(tx types.Transaction) error {
	switch tx.Type {
	case types.TxTypeContractDeploy:
		return cp.processContractDeploy(tx)
	case types.TxTypeContractCall:
		return cp.processContractCall(tx)
	default:
		return fmt.Errorf("not a contract transaction: type %d", tx.Type)
	}
}

// processContractDeploy deploys a new contract
func (cp *ContractProcessor) processContractDeploy(tx types.Transaction) error {
	// Parse deploy payload
	var payload types.ContractDeployPayload
	if err := json.Unmarshal(tx.Payload, &payload); err != nil {
		return fmt.Errorf("invalid deploy payload: %w", err)
	}

	// Deploy contract (all nodes execute this)
	contractAddr, err := cp.executor.DeployContract(
		payload.Bytecode,
		tx.Sender,
		payload.Constructor, // Correct field name
		tx.Gas,
	)
	if err != nil {
		return fmt.Errorf("contract deployment failed: %w", err)
	}

	fmt.Printf("✅ Contract deployed at %x by %x\n", contractAddr[:4], tx.Sender[:4])

	return nil
}

// processContractCall executes a contract method
func (cp *ContractProcessor) processContractCall(tx types.Transaction) error {
	// Parse call payload
	var payload types.ContractCallPayload
	if err := json.Unmarshal(tx.Payload, &payload); err != nil {
		return fmt.Errorf("invalid call payload: %w", err)
	}

	// Execute contract (all nodes execute this)
	result, err := cp.executor.ExecuteContractWithSecurity(
		payload.ContractAddress, // Correct field name
		tx.Sender,
		payload.Method,
		payload.Args,
		tx.Amount, // value
		tx.Gas,
	)
	if err != nil {
		return fmt.Errorf("contract execution failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("contract execution error: %s", result.Error)
	}

	fmt.Printf("✅ Contract %x.%s executed: %d gas used\n",
		payload.ContractAddress[:4], payload.Method, result.GasUsed)

	return nil
}

// ValidateContractTransaction validates a contract transaction before execution
func (cp *ContractProcessor) ValidateContractTransaction(tx types.Transaction) error {
	switch tx.Type {
	case types.TxTypeContractDeploy:
		var payload types.ContractDeployPayload
		if err := json.Unmarshal(tx.Payload, &payload); err != nil {
			return fmt.Errorf("invalid deploy payload: %w", err)
		}

		// Validate bytecode is not empty
		if len(payload.Bytecode) == 0 {
			return fmt.Errorf("empty bytecode")
		}

		// Validate gas limit
		if tx.Gas < types.GasContractDeploy {
			return fmt.Errorf("insufficient gas for deploy")
		}

	case types.TxTypeContractCall:
		var payload types.ContractCallPayload
		if err := json.Unmarshal(tx.Payload, &payload); err != nil {
			return fmt.Errorf("invalid call payload: %w", err)
		}

		// Validate contract exists
		_, err := cp.contractState.GetContract(payload.ContractAddress)
		if err != nil {
			return fmt.Errorf("contract not found: %w", err)
		}

		// Validate gas limit
		if tx.Gas < types.GasContractCall {
			return fmt.Errorf("insufficient gas for call")
		}
	}

	return nil
}
