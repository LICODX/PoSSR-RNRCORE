package params

const (
	// Realistic Block Parameters for Educational L1
	BlockTime      = 6                // 6 seconds (fast finality)
	MaxBlockSize   = 10 * 1024 * 1024 // 10 MB (Realistic constraint, addressing debat/9.txt)
	MaxMessageSize = 10 * 1024 * 1024 // 10 MB (LibP2P Limit Override)
	ShardSize      = 1 * 1024 * 1024  // 1 MB per shard
	NumShards      = 10               // 10 Shards

	// Tokenomics (5 Billion Supply, 7% Decay / 3.5M Blocks)
	TotalSupply     = 5000000000
	InitialReward   = 100.0   // 10 koin x 10 node
	HalvingInterval = 3500000 // Blok
	DecayRate       = 0.07    // 7%

	// Storage
	PruningWindow = 100 // Keep 100 blocks (Hardened from 25)

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
