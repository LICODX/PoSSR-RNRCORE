package types

// Gas costs for various operations
// These are calibrated to prevent DoS while keeping costs low

const (
	// Basic operations
	GasBalance        uint64 = 50
	GasTransfer       uint64 = 500
	GasEmitEvent      uint64 = 200
	GasGetCaller      uint64 = 50
	GasGetBlockHeight uint64 = 50

	// Storage operations (expensive)
	GasStorageRead  uint64 = 100
	GasStorageWrite uint64 = 1000

	// Smart contract operations
	GasContractDeploy uint64 = 100000 // Deploy contract
	GasContractCall   uint64 = 10000  // Call contract method
)
