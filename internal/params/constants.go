package params

const (
	// Dimensi Blok & Shard
	BlockTime    = 60                 // 60 Detik - Mainnet Production
	MaxBlockSize = 1024 * 1024 * 1024 // 1 GB Total (10 Shards x 100 MB)
	ShardSize    = 100 * 1024 * 1024  // 100 MB per Node
	NumShards    = 10                 // 10 Pemenang per Blok

	// Tokenomics (5 Billion Supply, 7% Decay / 3.5M Blocks)
	TotalSupply     = 5000000000
	InitialReward   = 100.0   // 10 koin x 10 node
	HalvingInterval = 3500000 // Blok
	DecayRate       = 0.07    // 7%

	// Storage
	PruningWindow = 25 // Keep 25 blocks (~25 minutes of data) - Whitepaper Spec

	// Transaction Fees (Anti-Spam)
	MinTxFee = 1 // Minimum 1 unit (0.000001 RNR) per transaction

	// Network
	// Network
	BootnodeIP   = "0.0.0.0" // Listen on ALL interfaces
	BootnodePort = "9900"

	// Genesis Config
	GenesisAddress = "rnr1pq03gqs8zg0sgqg7zsw3u8sgqqdp7rsrzuy3wxg7pyyqxrcspsr3cqq7qvqs78c2zyrpqzqdqvfq7xs8pcgq2m9d04"
	GenesisBalance = 0 // Coins from MINING
)
