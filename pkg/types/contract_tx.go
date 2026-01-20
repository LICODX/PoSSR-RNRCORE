package types

import "crypto/sha256"

// Smart Contract Transaction Types

const (
	// Smart Contract types
	TxTypeContractDeploy = 10 // Deploy new contract
	TxTypeContractCall   = 11 // Execute contract method
)

// ContractDeployPayload for deploying contracts
type ContractDeployPayload struct {
	Bytecode     []byte `json:"bytecode"`     // WASM bytecode
	Constructor  []byte `json:"constructor"`  // Constructor arguments
	InitialValue uint64 `json:"initialValue"` // RNR to send to contract
	GasLimit     uint64 `json:"gasLimit"`     // Maximum gas
}

// ContractCallPayload for calling contracts
type ContractCallPayload struct {
	ContractAddress [32]byte `json:"contractAddress"` // Target contract
	Method          string   `json:"method"`          // Method name
	Args            []byte   `json:"args"`            // Method arguments (encoded)
	Value           uint64   `json:"value"`           // RNR to send
	GasLimit        uint64   `json:"gasLimit"`        // Maximum gas
}

// GenerateContractAddress creates unique contract address
func GenerateContractAddress(deployer [32]byte, nonce uint64) [32]byte {
	// Contract address = SHA256(deployer + nonce)
	data := append(deployer[:], byte(nonce))
	return sha256.Sum256(data)
}
