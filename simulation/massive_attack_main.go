package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

const (
	TOTAL_NODES     = 100
	HONEST_NODES    = 20
	HACKER_NODES    = 80
	SIMULATION_TIME = 30 * time.Second
)

var (
	// Stats
	mu           sync.Mutex
	honestBlocks = 0
	hackerBlocks = 0
	totalBlocks  = 0

	// Simplified Chain for Fairness Test
	chainMtx    sync.Mutex
	simpleChain []types.Block
)

func main() {
	fmt.Println("🚀 STARTING MASSIVE MAINNET SIMULATION (100 NODES)")
	fmt.Printf("   - Honest Nodes: %d (Normal Hashrate)\n", HONEST_NODES)
	fmt.Printf("   - Hacker Nodes: %d (Aggressive/Spamming)\n", HACKER_NODES)
	fmt.Println("---------------------------------------------------")

	// Genesis
	simpleChain = append(simpleChain, types.Block{
		Header: types.BlockHeader{
			Height: 0,
			Hash:   [32]byte{},
		},
	})

	var wg sync.WaitGroup
	wg.Add(TOTAL_NODES)

	// Start Honest Nodes
	for i := 0; i < HONEST_NODES; i++ {
		go runNode(i, false, &wg)
	}

	// Start Hacker Nodes
	for i := 0; i < HACKER_NODES; i++ {
		go runNode(HONEST_NODES+i, true, &wg)
	}

	// Simulation Timer
	go func() {
		time.Sleep(SIMULATION_TIME)
		fmt.Println("\n⏰ SIMULATION TIME OVER. STOPPING NODES...")
		printStats()
		os.Exit(0)
	}()

	// Wait (indefinitely until timeout)
	wg.Wait()
}

func runNode(id int, isHacker bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// Helper to get tip safely
	getTip := func() types.BlockHeader {
		chainMtx.Lock()
		defer chainMtx.Unlock()
		if len(simpleChain) == 0 {
			return types.BlockHeader{}
		}
		return simpleChain[len(simpleChain)-1].Header
	}

	currentTip := getTip()
	stopMine := make(chan struct{})

	// Mining Loop
	for {
		// Update tip
		newTip := getTip()
		if newTip.Height > currentTip.Height {
			currentTip = newTip
			// Stop current mining job if new block found
			select {
			case stopMine <- struct{}{}:
			default:
			}
		}

		// MOCK MINING (Latency based)
		// Honest: 100-200ms
		// Hacker: 90-180ms (10% faster)

		delay := 100 + rand.Intn(100)
		if isHacker {
			// Hacker Advantage
			delay = 90 + rand.Intn(90)
		}

		select {
		case <-stopMine:
			continue
		case <-time.After(time.Duration(delay) * time.Millisecond):
		}

		// Create Block
		block := &types.Block{
			Header: types.BlockHeader{
				Height:        currentTip.Height + 1,
				PrevBlockHash: currentTip.Hash,
				Timestamp:     time.Now().Unix(),
				Hash: func() [32]byte {
					arr := [32]byte{byte(id), byte(rand.Int())}
					return utils.Hash(arr[:])
				}(),
			},
		}

		// Try to append
		chainMtx.Lock()
		tip := simpleChain[len(simpleChain)-1].Header
		if block.Header.PrevBlockHash == tip.Hash {
			simpleChain = append(simpleChain, *block)
			totalBlocks++
			if isHacker {
				hackerBlocks++
				fmt.Printf("🏴‍☠️ [Node %d] HACKER Won Block #%d\n", id, block.Header.Height)
			} else {
				honestBlocks++
				fmt.Printf("🛡️ [Node %d] HONEST Won Block #%d\n", id, block.Header.Height)
			}
			chainMtx.Unlock()
			printStats() // Update report
		} else {
			chainMtx.Unlock()
		}
	}
}

func printStats() {
	mu.Lock()
	defer mu.Unlock()
	honestPerc := 0.0
	if totalBlocks > 0 {
		honestPerc = float64(honestBlocks) / float64(totalBlocks) * 100
	}
	report := fmt.Sprintf("📊 Stats: Total=%d | Honest=%d (%.1f%%) | Hacker=%d (%.1f%%)\n",
		totalBlocks, honestBlocks, honestPerc, hackerBlocks, 100-honestPerc)
	fmt.Print(report)

	// Write to file (best effort)
	_ = os.WriteFile("sim_report.txt", []byte(report), 0644)
}
