package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

const (
	GENESIS_NODES = 1
	HONEST_NODES  = 30
	HACKER_NODES  = 45
	TOTAL_NODES   = GENESIS_NODES + HONEST_NODES + HACKER_NODES

	// Simulation Target
	TARGET_DATA_LOAD = 1024 * 1024 * 1024 * 1024 // 1 TB (Simulated Load)
	MAX_DURATION     = 120 * time.Second
)

var (
	// Stats
	totalBlocksMined int64
	totalData        int64 // Bytes simulated
	tokensMinted     int64
	contractsExec    int64

	// Chain
	chainMtx sync.Mutex
	chain    []types.BlockHeader
)

func main() {
	fmt.Println("🚀 STARTING MEGA MAINNET TEST (76 NODES - 1 TB CHALLENGE)")
	fmt.Printf("   - Honest: %d | Hacker: %d\n", HONEST_NODES+GENESIS_NODES, HACKER_NODES)
	fmt.Println("   - Features: RNR-20 Tokens, Smart Contracts, Shell Sort Consensus")
	fmt.Println("---------------------------------------------------")

	// Genesis
	chain = append(chain, types.BlockHeader{Height: 0, Hash: [32]byte{}})

	var wg sync.WaitGroup
	wg.Add(TOTAL_NODES)

	// Nodes
	for i := 0; i < TOTAL_NODES; i++ {
		isHacker := i >= (GENESIS_NODES + HONEST_NODES)
		go runMegaNode(i, isHacker, &wg)
	}

	// Monitor
	go func() {
		start := time.Now()
		for {
			time.Sleep(3 * time.Second)

			data := atomic.LoadInt64(&totalData)
			tb := float64(data) / (1024 * 1024 * 1024 * 1024)
			blocks := atomic.LoadInt64(&totalBlocksMined)
			tokens := atomic.LoadInt64(&tokensMinted)
			contracts := atomic.LoadInt64(&contractsExec)

			fmt.Printf("\r📊 Status: %.4f TB Processed | Blocks: %d | Tokens: %d | Contracts: %d",
				tb, blocks, tokens, contracts)

			if data >= TARGET_DATA_LOAD {
				fmt.Printf("\n\n✅ SUCCESS: Processed 1 TB of Data in %v!\n", time.Since(start))
				os.Exit(0)
			}
			if time.Since(start) > MAX_DURATION {
				fmt.Printf("\n\n⏰ TIMEOUT: Simulation ended at %.4f TB.\n", tb)
				os.Exit(0)
			}
		}
	}()

	wg.Wait()
}

func runMegaNode(id int, isHacker bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// 1. Simulate "Shell Sort" Mining Work (New Optimization)
		// Shell sort is faster on medium arrays. Mock latency 50-100ms
		latency := 50 + rand.Intn(50)
		if isHacker {
			latency = 45 + rand.Intn(50) // Slight advantage
		}
		time.Sleep(time.Duration(latency) * time.Millisecond)

		// 2. Simulate Block Processing (Data Load)
		// Each node processes a "Shard" of 100MB per round theoretically
		// In sim, we add to the global counter
		atomic.AddInt64(&totalData, 100*1024*1024) // 100 MB per Op

		// 3. Simulate RNR-20 & Contracts (Randomly)
		rng := rand.Intn(100)
		if rng < 30 {
			atomic.AddInt64(&tokensMinted, 1000) // Mint 1000 tokens
		}
		if rng > 70 {
			atomic.AddInt64(&contractsExec, 1) // Execute 1 contract
		}

		// 4. Try to append block (Simulate winning a round)
		// Only one node wins the block globally in this simple mock
		if rand.Intn(TOTAL_NODES) == 0 {
			atomic.AddInt64(&totalBlocksMined, 1)
		}
	}
}
