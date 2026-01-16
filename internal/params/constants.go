package params

const (
	// Dimensi Blok & Shard
	BlockTime    = 60                // 1 Menit (detik)
	MaxBlockSize = 100 * 1024 * 1024 // 100 MB per Block
	ShardSize    = 100 * 1024 * 1024 // 100 MB per Node
	NumShards    = 10                // 10 Pemenang per Blok

	// Tokenomics (5 Billion Supply, 7% Decay / 3.5M Blocks)
	TotalSupply     = 5000000000
	InitialReward   = 100.0   // 10 koin x 10 node
	HalvingInterval = 3500000 // Blok
	DecayRate       = 0.07    // 7%

	// Storage
	PruningWindow = 2880 // Keep 48 hours of blocks (was 25)
)
