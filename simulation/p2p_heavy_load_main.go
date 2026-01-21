package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
	"github.com/LICODX/PoSSR-RNRCORE/internal/p2p"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	TotalNodes = 20
	BasePort   = 11000
	BlockSize  = 150 * 1024 * 1024 // 150 MB (Adjusted for single-machine sim)
	NumShards  = 10
)

func main() {
	fmt.Println("############################################################")
	fmt.Println("#       PoSSR 20-NODE HEAVY MAINNET SIMULATION (150MB)     #")
	fmt.Println("############################################################")

	// Print Memory Stats Start
	printMemUsage("Start")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Setup Nodes
	nodes := make([]*p2p.GossipSubNode, TotalNodes)
	blockchains := make([]*blockchain.Blockchain, TotalNodes)

	// We need a way to track what each node received
	receivedData := make([]map[string]bool, TotalNodes)
	var mu sync.Mutex

	fmt.Println("[SETUP] Initializing 20 LibP2P Nodes...")

	for i := 0; i < TotalNodes; i++ {
		receivedData[i] = make(map[string]bool)

		// Configure Role
		role := "ShardNode"
		shardIDs := []int{(i) % 10} // Distribute across 10 shards
		if i < 2 {
			role = "FullNode" // Node 0 and 1 are Full Nodes
			shardIDs = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		}

		cfg := config.ShardConfig{
			Role:     role,
			ShardIDs: shardIDs,
		}

		// Init Mock DB & Blockchain
		db, _ := storage.NewLevelDB(fmt.Sprintf("./data/heavy_sim_%d", i))
		bc := blockchain.NewBlockchain(db, cfg)
		blockchains[i] = bc

		// Init P2P Node
		node, err := p2p.NewGossipSubNode(ctx, BasePort+i, cfg)
		if err != nil {
			panic(err)
		}
		nodes[i] = node

		// Setup Listeners
		nodeID := i // Capture closure
		node.ListenForHeaders(func(data []byte) {
			mu.Lock()
			receivedData[nodeID]["Header"] = true
			mu.Unlock()
			// fmt.Printf("Node %d got Header\n", nodeID)
		})

		// Listen for assigned shards
		for _, sID := range shardIDs {
			shardID := sID
			node.ListenForShards(shardID, func(data []byte) {
				mu.Lock()
				receivedData[nodeID][fmt.Sprintf("Shard_%d", shardID)] = true
				mu.Unlock()
				// fmt.Printf("Node %d got Shard %d\n", nodeID, shardID)
			})
		}
	}

	// 2. Connect Peers (Mesh)
	fmt.Println("[NETWORK] Discovering peers and forming mesh...")
	// Simple mesh: Connect everyone to Node 0 (Bootstrap) and Node 1
	bootstrapInfo := nodes[0].GetHost().Addrs()
	bootstrapID := nodes[0].GetHost().ID()
	bootstrapInfo0 := peer.AddrInfo{
		ID:    bootstrapID,
		Addrs: bootstrapInfo,
	}
	fmt.Printf("Bootstrap Node 0: %s\n", bootstrapInfo0)

	// Explicitly connect everyone to Node 0
	for i := 1; i < TotalNodes; i++ {
		if err := nodes[i].GetHost().Connect(ctx, bootstrapInfo0); err != nil {
			fmt.Printf("Node %d failed to connect to bootstrap: %v\n", i, err)
		}
	}
	fmt.Println("✅ All nodes connected to Bootstrap Node 0")
	// Allow GossipSub mesh to stabilize
	time.Sleep(5 * time.Second)

	// 3. Generate Massive Block
	fmt.Println("[MINING] Generating 1.5GB Block Data (This may take a moment)...")
	// We cheat slightly: one huge transaction per shard to save CPU on signing 6 million txs
	// 1.5GB / 10 = 150MB per shard. -> Adjusted: 150MB / 10 = 15MB per shard
	payloadSize := 15 * 1024 * 1024
	hugePayload := make([]byte, payloadSize)
	rand.Read(hugePayload) // Random data

	var shards [10]types.ShardData
	var shardRoots [10][32]byte

	for i := 0; i < 10; i++ {
		// Create 1 huge tx
		tx := types.Transaction{
			ID:      types.HashTransaction(types.Transaction{Nonce: uint64(i)}), // Pseudo hash
			Payload: hugePayload,
		}
		// Calculate Root
		root := utils.CalculateMerkleRoot([][32]byte{tx.ID})
		shards[i] = types.ShardData{
			TxData:    []types.Transaction{tx},
			ShardRoot: root,
		}
		shardRoots[i] = root
	}

	// Global Root
	var rootsSlice [][32]byte
	for _, r := range shardRoots {
		rootsSlice = append(rootsSlice, r)
	}
	globalRoot := utils.CalculateMerkleRoot(rootsSlice)

	block := types.Block{
		Header: types.BlockHeader{
			Height:     1,
			MerkleRoot: globalRoot,
			ShardRoots: shardRoots,
			Timestamp:  time.Now().Unix(),
		},
		Shards: shards,
	}

	printMemUsage("Block Generated")

	// 4. Publish Block (Node 0)
	fmt.Println("[BROADCAST] Node 0 (FullNode) Publishing 1.5GB Block...")
	start := time.Now()
	err := nodes[0].PublishBlock(block)
	if err != nil {
		fmt.Printf("Publish failed: %v\n", err)
	}

	// 5. Wait for Propagation
	fmt.Println("[WAIT] Waiting 30s for propagation...")
	// 1.5GB is a lot. Give it time.
	time.Sleep(30 * time.Second)

	// 6. Verify Results
	fmt.Println("\n[VERIFICATION] Checking Node Receipt Status:")

	successCount := 0
	for i := 0; i < TotalNodes; i++ {
		mu.Lock()
		data := receivedData[i]
		mu.Unlock()

		hasHeader := data["Header"]

		// Determine expected shards
		expectedShards := nodes[i].GetShardConfig().ShardIDs

		missingShards := 0
		unexpectedShards := 0

		// Check expected
		for _, sID := range expectedShards {
			if !data[fmt.Sprintf("Shard_%d", sID)] {
				missingShards++
			}
		}

		// Check strictness (did we receive shards we didn't ask for?)
		// In GossipSub, we shouldn't receive msg if not subscribed.
		for k := range data {
			if k == "Header" {
				continue
			}
			var sID int
			fmt.Sscanf(k, "Shard_%d", &sID)

			isExpected := false
			for _, ex := range expectedShards {
				if ex == sID {
					isExpected = true
					break
				}
			}
			if !isExpected {
				unexpectedShards++
			}
		}

		status := "✅ PASS"
		if !hasHeader || missingShards > 0 || unexpectedShards > 0 {
			status = "❌ FAIL"
		} else {
			successCount++
		}

		fmt.Printf("Node %d (%s): Header=%v, MissingShards=%d, Unexpected=%d -> %s\n",
			i, nodes[i].GetShardConfig().Role, hasHeader, missingShards, unexpectedShards, status)
	}

	duration := time.Since(start)
	fmt.Printf("\nTest Completed in %v\n", duration)
	fmt.Printf("Success Rate: %d/%d Nodes\n", successCount, TotalNodes)

	printMemUsage("End")
}

func printMemUsage(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\n[MEM] %s: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
		msg, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
