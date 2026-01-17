package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus"
	"github.com/LICODX/PoSSR-RNRCORE/internal/dashboard"
	"github.com/LICODX/PoSSR-RNRCORE/internal/params"

	// "github.com/LICODX/PoSSR-RNRCORE/internal/dashboard"

	// 5. Start GUI Dashboard (Disabled for Headless Build)
	// dashboard.StartServer("8080", chain, nil)
	"github.com/LICODX/PoSSR-RNRCORE/internal/p2p"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

func main() {
	port := flag.Int("port", 3000, "P2P listening port")
	datadir := flag.String("datadir", "./data/chaindata", "Data directory for LevelDB")
	peers := flag.String("peers", "", "Comma-separated peer addresses")
	useGossipSub := flag.Bool("gossipsub", true, "Use GossipSub (recommended)")
	flag.Parse()

	fmt.Println("🚀 Starting rnr-core Mainnet Node...")
	fmt.Printf("Config: Port=%d | DataDir=%s\n", *port, *datadir)
	fmt.Println("Consensus: PoSSR | Block Size: 1 GB | Pruning: ON")

	// 1. Initialize Database
	db, err := storage.NewLevelDB(*datadir)
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		return
	}
	defer db.GetDB().Close()

	// 2. Initialize Blockchain State
	chain := blockchain.NewBlockchain(db)
	tip := chain.GetTip()
	fmt.Printf("⛓️  Current Tip: Block #%d\n", tip.Height)

	// 3. Setup context
	ctx := context.Background()

	// Create a channel to stop mining if needed (e.g. new block received)
	stopMining := make(chan struct{})

	// 4. Start P2P Network
	var node *p2p.GossipSubNode

	if *useGossipSub {
		var err error
		node, err = p2p.NewGossipSubNode(ctx, *port)
		if err != nil {
			fmt.Printf("Failed to start GossipSub: %v\n", err)
			return
		}
		defer node.Close()

		if *peers != "" {
			fmt.Printf("Connecting to peers: %s\n", *peers)
		}

		node.DiscoverPeers()

		node.ListenForBlocks(func(data []byte) {
			fmt.Println("📦 Received block from network")
			var block types.Block
			if err := json.Unmarshal(data, &block); err != nil {
				fmt.Printf("Failed to decode block: %v\n", err)
				return
			}
			if err := chain.AddBlock(block); err == nil {
				// Signal to restart mining on new head
				select {
				case stopMining <- struct{}{}:
				default:
				}
			}
		})

		node.ListenForTransactions(func(data []byte) {
			fmt.Println("💸 Received transaction from network")
			var tx types.Transaction
			if err := json.Unmarshal(data, &tx); err != nil {
				fmt.Printf("Failed to decode tx: %v\n", err)
				return
			}

			// SECURITY CHECK: Validate against State (Nonce & Balance)
			// This prevents Replay Attacks and insufficient balance spam
			if err := blockchain.ValidateTransactionAgainstState(tx, chain.GetStateManager()); err != nil {
				fmt.Printf("⚠️ Invalid transaction rejected: %v\n", err)
				return
			}

			node.AddToMempool(tx)
		})

		node.ListenForProofs(func(data []byte) {
			fmt.Println("✅ Received proof from network")
		})

	} else {
		fmt.Println("Legacy TCP not supported in this version.")
		return
	}

	// 5. Start GUI Dashboard
	// Pass 'node' which implements MempoolSource
	dashboard.StartServer("8080", chain, node)

	// 6. Mining Loop (Proof of Repeated Sorting)
	fmt.Println("🏁 Mining Loop Started. Searching for a valid block...")

	for {
		lastHeader := chain.GetTip()
		difficulty := uint64(1000)

		// Get transactions from P2P mempool
		txs := node.GetMempoolShard()

		if len(txs) == 0 {
			// For Demo: If empty, add a mock transaction so we can mine
			// In production, we would wait.
			// fmt.Println("Waiting for transactions...")
			time.Sleep(2 * time.Second)
			node.AddToMempool(types.Transaction{
				ID:     [32]byte{1, 2, 3},
				Amount: 100,
				Nonce:  uint64(time.Now().UnixNano()),
			})
			continue
		}

		fmt.Printf("🔨 Mining on top of Block #%d [Diff: %d] with %d txs\n", lastHeader.Height, difficulty, len(txs))

		newBlock, err := consensus.MineBlock(txs, lastHeader, difficulty, stopMining)

		if err != nil {
			if err.Error() == "mining interrupted" {
				fmt.Println("Mining interrupted! Restarting...")
				continue
			}
			fmt.Println("Mining error:", err)
			continue
		}

		fmt.Printf("💎 Block Found! Nonce: %d | Hash: %x\n", newBlock.Header.Nonce, newBlock.Header.Hash)

		// Add to local chain
		if err := chain.AddBlock(*newBlock); err != nil {
			fmt.Printf("Failed to add block: %v\n", err)
			continue
		}

		fmt.Printf("✅ Block Accepted! Height: %d\n", newBlock.Header.Height)

		// THROTTLE: Wait for BlockTime (6s) to ensure consistent heartbeat
		fmt.Printf("⏳ Waiting %d seconds for next round...\n", params.BlockTime)
		time.Sleep(time.Duration(params.BlockTime) * time.Second)

		// Broadcast Block
		if blockBytes, err := json.Marshal(*newBlock); err == nil {
			node.PublishBlock(blockBytes)
		}

		// Clear mempool (Simplified)
		node.ClearMempool()
	}
}
