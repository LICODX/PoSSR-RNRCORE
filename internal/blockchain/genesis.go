package blockchain

import (
	"crypto/rand"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// MAINNET GENESIS BLOCK - FIXED AND IMMUTABLE

const (
	// Mainnet Genesis Timestamp (FIXED)
	MainnetGenesisTimestamp = 1735689600 // 2025-01-01 00:00:00 UTC

	// Mainnet Genesis Wallet (MUST BE SET BEFORE LAUNCH)
	MainnetGenesisWallet = "8150a6af22851558e96cb9faad6b7e9cd5961179deb84c784fdf5bbb5d57b263"

	// Chain ID
	MainnetChainID = 1
	TestnetChainID = 1337
)

var (
	// MainnetGenesisBlock - FIXED for production
	MainnetGenesisBlock = types.Block{
		Header: types.BlockHeader{
			Version:       1,
			PrevBlockHash: [32]byte{}, // All zeros
			MerkleRoot:    [32]byte{}, // Empty
			Timestamp:     MainnetGenesisTimestamp,
			Height:        0,
			VRFSeed:       [32]byte{0xde, 0xad, 0xbe, 0xef}, // FIXED seed for mainnet
		},
		Shards: [10]types.ShardData{},
	}

	// TestnetGenesisBlock - Random for testnet
	TestnetGenesisBlock = types.Block{
		Header: types.BlockHeader{
			Version:       1,
			PrevBlockHash: [32]byte{},
			MerkleRoot:    [32]byte{},
			Timestamp:     0, // Will use current time
			Height:        0,
			VRFSeed:       [32]byte{}, // Will be random
		},
		Shards: [10]types.ShardData{},
	}
)

// CreateMainnetGenesis returns FIXED genesis block for mainnet
func CreateMainnetGenesis() types.Block {
	// Return exact copy - NO RANDOMNESS
	block := MainnetGenesisBlock

	// CRITICAL: Calculate Hash using PoW hash function (consistent with all blocks)
	block.Header.Hash = types.HashBlockHeaderForPoW(block.Header)

	return block
}

// CreateTestnetGenesis returns genesis block for testnet (can be random)
func CreateTestnetGenesis() types.Block {
	// Generate random VRF seed for testnet
	var seed [32]byte
	rand.Read(seed[:])

	block := TestnetGenesisBlock
	block.Header.Timestamp = time.Now().Unix()
	block.Header.VRFSeed = seed
	return block
}

// CreateGenesisBlock creates genesis based on network type
func CreateGenesisBlock(isMainnet bool) types.Block {
	if isMainnet {
		return CreateMainnetGenesis()
	}
	return CreateTestnetGenesis()
}
