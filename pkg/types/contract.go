package types

// Smart Contract Types for RNR Blockchain

// Contract represents a deployed smart contract
type Contract struct {
	Address     [32]byte // Unique contract address
	Creator     [32]byte // Deployer address
	Bytecode    []byte   // WASM bytecode
	CodeHash    [32]byte // Hash of bytecode
	CreatedAt   int64    // Unix timestamp
	Balance     uint64   // Contract's RNR balance
	StorageRoot [32]byte // Merkle root of contract storage
}

// ContractStorage represents contract's persistent storage
type ContractStorage struct {
	Contract [32]byte          // Contract address
	Data     map[string][]byte // Key-value storage
}

// Gas represents gas for contract execution
type Gas struct {
	Limit  uint64 // Maximum gas allowed
	Used   uint64 // Gas consumed
	Price  uint64 // Gas price in RNR (per unit)
	Refund uint64 // Gas to refund
}

// ContractCall represents a call to contract method
type ContractCall struct {
	Contract [32]byte // Target contract
	Method   string   // Method name
	Args     []byte   // Encoded arguments
	Value    uint64   // RNR to send
	Gas      Gas      // Gas parameters
	Caller   [32]byte // Caller address
}

// ContractResult represents execution result
type ContractResult struct {
	Success    bool    // Execution success
	ReturnData []byte  // Return value
	GasUsed    uint64  // Total gas consumed
	Events     []Event // Emitted events
	Error      string  // Error message if failed
}

// Event represents a contract event log
type Event struct {
	Contract [32]byte // Contract that emitted
	Topic    string   // Event topic/name
	Data     []byte   // Event data
	TxHash   [32]byte // Transaction hash
	Index    uint32   // Event index in TX
}

// HostFunction represents functions provided by blockchain to contracts
type HostFunction int

const (
	HostFnTransferRNR HostFunction = iota
	HostFnTransferToken
	HostFnGetBlockNumber
	HostFnGetBlockTimestamp
	HostFnGetCaller
	HostFnGetBalance
	HostFnEmitEvent
	HostFnStorageRead
	HostFnStorageWrite
	HostFnCallContract
	HostFnSHA256
)
