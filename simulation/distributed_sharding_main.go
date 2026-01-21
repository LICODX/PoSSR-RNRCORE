package main

import (
	"fmt"
	"os"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// Simulation Parameters
const (
	TotalNodes = 20
)

type Node struct {
	ID         int
	Config     config.ShardConfig
	Blockchain *blockchain.Blockchain
	State      *state.Manager
}

func main() {
	fmt.Println("############################################################")
	fmt.Println("#    PoSSR DISTRIBUTED SHARDING & VALIDATION SIMULATION    #")
	fmt.Println("############################################################")

	// 1. Setup Network
	nodes := setupNetwork()
	defer cleanupNetwork(nodes)

	fmt.Println("\n[PHASE 1] GENERATING BLOCK DATA...")

	// Create Genesis Header
	genesisHeader := types.BlockHeader{
		Height: 0,
		Hash:   [32]byte{0x01}, // Mock
	}

	// Create a VALID Full Block with transactions in Shard 0 and Shard 1
	txSh0 := types.Transaction{ID: [32]byte{0xAA}, Sender: [32]byte{1}, Amount: 10, Nonce: 1} // 0xAA
	txSh1 := types.Transaction{ID: [32]byte{0xBB}, Sender: [32]byte{2}, Amount: 10, Nonce: 1} // 0xBB

	// Calculate Roots
	root0 := utils.CalculateMerkleRoot([][32]byte{txSh0.ID})
	root1 := utils.CalculateMerkleRoot([][32]byte{txSh1.ID})
	emptyRoot := utils.CalculateMerkleRoot(nil)

	// Shard Roots
	var shardRoots [10][32]byte
	shardRoots[0] = root0
	shardRoots[1] = root1
	for i := 2; i < 10; i++ {
		shardRoots[i] = emptyRoot
	}

	// Global Root
	var shardRootSlice [][32]byte
	for _, r := range shardRoots {
		shardRootSlice = append(shardRootSlice, r)
	}
	globalRoot := utils.CalculateMerkleRoot(shardRootSlice)

	fullBlock := types.Block{
		Header: types.BlockHeader{
			Height:        1,
			PrevBlockHash: genesisHeader.Hash,
			MerkleRoot:    globalRoot,
			ShardRoots:    shardRoots,
			Timestamp:     time.Now().Unix(),
		},
		Shards: [10]types.ShardData{
			{TxData: []types.Transaction{txSh0}}, // Shard 0
			{TxData: []types.Transaction{txSh1}}, // Shard 1
			// Others empty
		},
	}

	fmt.Printf("✅ Generated Full Block. Global Root: %x\n", globalRoot)

	// 2. Validate with Full Node (Node 0)
	fmt.Println("\n[TEST 1] VALIDATION BY FULL NODE (Node 0)")
	fmt.Printf("   Node Config: %s, Shards: %v\n", nodes[0].Config.Role, nodes[0].Config.ShardIDs)

	err := blockchain.ValidateBlock(fullBlock, genesisHeader, nodes[0].Config)
	if err != nil {
		fmt.Printf("   ❌ FAILED: Full Node rejected valid block: %v\n", err)
	} else {
		fmt.Printf("   ✅ SUCCESS: Full Node accepted valid block.\n")
	}

	// 3. Validate with Shard Node 0 (Node 2) - Partial Data
	fmt.Println("\n[TEST 2] VALIDATION BY SHARD 0 NODE (Node 2)")
	fmt.Printf("   Node Config: %s, Shards: %v\n", nodes[2].Config.Role, nodes[2].Config.ShardIDs)

	// Simulate "Partial Block": Node 2 receives Shard 0 but NOT Shard 1
	partialBlockSh0 := fullBlock
	partialBlockSh0.Shards[1].TxData = nil // Remove Shard 1 data
	// Note: ShardRoots in Header are STILL PRESENT (Header is always full)

	err = blockchain.ValidateBlock(partialBlockSh0, genesisHeader, nodes[2].Config)
	if err != nil {
		fmt.Printf("   ❌ FAILED: Shard 0 Node rejected valid partial block: %v\n", err)
	} else {
		fmt.Printf("   ✅ SUCCESS: Shard 0 Node accepted partial block (Verified Sh0, Trusted Sh1).\n")
	}

	// 4. Validate with Shard Node 1 (Node 3 or similar) with DATA MISMATCH
	// Scenario: Node receives Shard 1 data, but it doesn't match the Root in header
	fmt.Println("\n[TEST 3] FRAUD DETECTION (Tampered Data)")

	// Create Tampered Block
	tamperedBlock := fullBlock
	// Modify Shard 0 transaction ID to change hash
	tamperedBlock.Shards[0].TxData[0].Amount = 999999
	// The Header.ShardRoots[0] matches the ORIGINAL data.
	// But the TxData provided is MODIFIED.
	// Local calculation of root will mismatch Header.

	fmt.Println("   📝 Simulating Attack: Miner sends valid Header but Modified Data to Local Node.")

	err = blockchain.ValidateBlock(tamperedBlock, genesisHeader, nodes[0].Config) // Full Node checks Shard 0
	if err != nil {
		fmt.Printf("   ✅ SUCCESS: Full Node DETECTED mismatch: %v\n", err)
	} else {
		fmt.Printf("   ❌ FAIL: Full Node ACCEPTED tampered block!\n")
	}

	// 5. Sorting Check
	fmt.Println("\n[TEST 4] SORTING VERIFICATION (O(N) Check)")
	// Create unsorted shard
	txA := types.Transaction{ID: [32]byte{0xBB}} // Big
	txB := types.Transaction{ID: [32]byte{0xAA}} // Small
	// Unsorted: [BB, AA]

	// Recalculate root for unsorted data (Root is agnostic to order usually? No, Merkle Tree depends on order!)
	// Wait, if order changes, Root changes.
	// Miner provides Header with Root(Sorted).
	// Attacker provides Unsorted Data that leads to Wrong Root -> Caught by Root Check.
	// Attacker provides Header with Root(Unsorted).
	// Node receives Unsorted Data -> Root Matches Header.
	// BUT ValidateBlock does Sorting Check explicitly!

	unsortedHashes := [][32]byte{txA.ID, txB.ID}
	unsortedRoot := utils.CalculateMerkleRoot(unsortedHashes)

	// Malicious Header matches Malicious Data
	badHeader := genesisHeader
	badHeader.ShardRoots[0] = unsortedRoot
	globalParts := fillRoots(unsortedRoot)
	badHeader.MerkleRoot = utils.CalculateMerkleRoot(globalParts)

	badBlock := types.Block{
		Header: badHeader,
		Shards: [10]types.ShardData{
			{TxData: []types.Transaction{txA, txB}},
		},
	}

	err = blockchain.ValidateBlock(badBlock, genesisHeader, nodes[0].Config)
	if err != nil {
		fmt.Printf("   ✅ SUCCESS: Unsorted Block Detected: %v\n", err)
	} else {
		fmt.Printf("   ❌ FAIL: Unsorted Block Accepted!\n")
	}
}

func fillRoots(r0 [32]byte) [][32]byte {
	var res [][32]byte
	res = append(res, r0)
	empty := utils.CalculateMerkleRoot(nil)
	for i := 1; i < 10; i++ {
		res = append(res, empty)
	}
	return res
}

func setupNetwork() []*Node {
	fmt.Println("--- Initializing 20 Nodes ---")
	nodes := make([]*Node, TotalNodes)

	for i := 0; i < TotalNodes; i++ {
		// Mock DB
		dbPath := fmt.Sprintf("./data/sim_shard_%d", i)
		os.RemoveAll(dbPath)
		db, _ := storage.NewLevelDB(dbPath)

		// Config
		var cfg config.ShardConfig
		if i < 2 {
			cfg = config.ShardConfig{Role: "FullNode", ShardIDs: []int{}}
		} else {
			// Shard Nodes
			// Distribute 18 nodes across 10 shards
			// Shard ID = (i-2) % 10
			shardID := (i - 2) % 10
			cfg = config.ShardConfig{Role: "ShardNode", ShardIDs: []int{shardID}}
		}

		nodes[i] = &Node{
			ID:         i,
			Config:     cfg,
			State:      state.NewManager(db.GetDB()),
			Blockchain: blockchain.NewBlockchain(db, cfg),
		}

		if i < 5 {
			fmt.Printf("Node %d: Role=%s, Shards=%v\n", i, cfg.Role, cfg.ShardIDs)
		}
	}
	fmt.Println("...")
	return nodes
}

func cleanupNetwork(nodes []*Node) {
	for _, n := range nodes {
		os.RemoveAll(fmt.Sprintf("./data/sim_shard_%d", n.ID))
	}
}
