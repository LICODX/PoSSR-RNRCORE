package types

// Transaction (500 Bytes avg)
// Menggunakan array byte tetap untuk menghindari GC Overhead berlebih
type Transaction struct {
	ID        [32]byte // SHA-256 Hash
	Type      int      // Transaction type (transfer, token, contract, etc)
	Sender    [32]byte
	Receiver  [32]byte
	Amount    uint64
	Fee       uint64 // Transaction fee (prevents spam, goes to miner)
	Gas       uint64 // Gas limit for contract execution
	Nonce     uint64
	Signature [64]byte
	Payload   []byte // Data for Smart Contracts (Optional)
}

// Block Header (Ringan, disimpan selamanya)
type BlockHeader struct {
	Version       uint32
	PrevBlockHash [32]byte
	MerkleRoot    [32]byte // Root dari gabungan 10 Shard Roots
	Timestamp     int64
	Height        uint64
	Nonce         uint64       // Mining counter
	Difficulty    uint64       // Target for the hash
	Hash          [32]byte     // Block Hash
	WinningNodes  [10][32]byte // PubKey 10 Pemenang
	ShardRoots    [10][32]byte // Merkle Roots of each Shard (New for Distributed Validation)
	VRFSeed       [32]byte     // Seed untuk blok berikutnya
}

// ShardData mewakili kontribusi 1 node
type ShardData struct {
	NodeID    [32]byte
	AlgoUsed  string // e.g., "QuickSort"
	TxData    []Transaction
	ShardRoot [32]byte
}

// Full Block (Berat 1 GB, akan di-prune)
type Block struct {
	Header BlockHeader
	Shards [10]ShardData // Data 10 x 100 MB
}
