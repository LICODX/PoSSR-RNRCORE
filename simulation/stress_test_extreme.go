package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

const (
	GENESIS_NODES   = 1
	HONEST_NODES    = 20
	HACKER_NODES    = 50
	TARGET_BLOCKS   = 30
	SIMULATION_TIME = 60 * time.Second // Hard limit
)

var (
	// Stats
	totalBlocksMined int32 = 0
	bytesSpammed     int64 = 0
	txsRejected      int64 = 0
	txsAccepted      int64 = 0

	// Chain State
	chainMtx    sync.Mutex
	simpleChain []types.Block
)

func main() {
	fmt.Println("🚀 STARTING EXTREME STRESS TEST (100GB LOAD CHALLENGE)")
	fmt.Println("---------------------------------------------------")
	fmt.Printf("   - Honest: %d | Hacker: %d\n", GENESIS_NODES+HONEST_NODES, HACKER_NODES)
	fmt.Printf("   - Goal: Mine %d Blocks under heavy load\n", TARGET_BLOCKS)

	// Genesis
	simpleChain = append(simpleChain, types.Block{
		Header: types.BlockHeader{Height: 0, Hash: [32]byte{}},
	})

	var wg sync.WaitGroup
	wg.Add(GENESIS_NODES + HONEST_NODES + HACKER_NODES)

	// Honest (Genesis + Others)
	for i := 0; i < GENESIS_NODES+HONEST_NODES; i++ {
		go runMiner(i, false, &wg)
	}

	// Hackers
	for i := 0; i < HACKER_NODES; i++ {
		go runHacker(i, &wg)
	}

	// Stats Reporter
	go func() {
		for {
			time.Sleep(2 * time.Second)
			printRealTimeStats()
			if atomic.LoadInt32(&totalBlocksMined) >= int32(TARGET_BLOCKS) {
				fmt.Println("\n✅ TARGET REACHED! 30 BLOCKS MINED.")
				printRealTimeStats()
				os.Exit(0)
			}
		}
	}()

	// Timeout
	go func() {
		time.Sleep(SIMULATION_TIME)
		fmt.Println("\n⏰ TIMEOUT REACHED.")
		printRealTimeStats()
		os.Exit(0)
	}()

	wg.Wait()
}

func runMiner(id int, isHacker bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Verify Tip (Simulate Consensus Check)
		chainMtx.Lock()
		tip := simpleChain[len(simpleChain)-1].Header
		chainMtx.Unlock()

		// Mine (Simulate 100ms block time for fast test)
		time.Sleep(100 * time.Millisecond)

		// Create Block
		block := types.Block{
			Header: types.BlockHeader{
				Height:        tip.Height + 1,
				PrevBlockHash: tip.Hash,
				Timestamp:     time.Now().Unix(),
				Hash: func() [32]byte {
					arr := [32]byte{byte(id), byte(rand.Int())}
					return utils.Hash(arr[:])
				}(),
			},
		}

		// Try Append
		chainMtx.Lock()
		if simpleChain[len(simpleChain)-1].Header.Hash == block.Header.PrevBlockHash {
			simpleChain = append(simpleChain, block)
			atomic.AddInt32(&totalBlocksMined, 1)
			fmt.Printf("⛏️  Block #%d Mined by Honest Node %d\n", block.Header.Height, id)
		}
		chainMtx.Unlock()

		if atomic.LoadInt32(&totalBlocksMined) >= int32(TARGET_BLOCKS) {
			return
		}
	}
}

func runHacker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Pre-generate a 1MB payload to spam
	payload := make([]byte, 1024*1024)
	rand.Read(payload)

	for {
		if atomic.LoadInt32(&totalBlocksMined) >= int32(TARGET_BLOCKS) {
			return
		}

		// ATTACK: Spam fake transactions/data
		// In a real node, this hits the API/P2P layer.
		// We simulate the *cost* of rejecting it.

		// Simulate Bandwidth Usage
		atomic.AddInt64(&bytesSpammed, int64(len(payload)))

		// Simulate Validation Cost (CPU cycle to check and reject)
		// We deliberately burn some CPU to model the load
		_ = utils.Hash(payload) // Hash check

		atomic.AddInt64(&txsRejected, 1)

		// Small sleep to allow context switching (so we don't freeze the OS completely)
		time.Sleep(1 * time.Millisecond)
	}
}

func printRealTimeStats() {
	blocks := atomic.LoadInt32(&totalBlocksMined)
	bytes := atomic.LoadInt64(&bytesSpammed)
	rejected := atomic.LoadInt64(&txsRejected)

	gb := float64(bytes) / 1024 / 1024 / 1024

	report := fmt.Sprintf("\r📊 Progress: %d/%d Blocks | Spam Load: %.2f GB | Txs Rejected: %d",
		blocks, TARGET_BLOCKS, gb, rejected)
	fmt.Print(report)
	os.WriteFile("stress_report.txt", []byte(report), 0644)
}
