package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Constants for the Stress Test
const (
	NumNodes       = 3
	TargetTxPerMin = 10_000_000 // 10 Million
	DurationSec    = 10         // Run for 10 seconds to extrapolate
)

// Mock Node with a Mempool
type Node struct {
	ID           int
	Mempool      map[[32]byte]types.Transaction
	MempoolLimit int
	mu           sync.RWMutex
	TxCount      int64 // Ingested Count
}

func NewNode(id int) *Node {
	return &Node{
		ID:           id,
		Mempool:      make(map[[32]byte]types.Transaction),
		MempoolLimit: 500_000, // Cap mempool to simulate RAM limits
	}
}

func (n *Node) AddTx(tx types.Transaction) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if len(n.Mempool) >= n.MempoolLimit {
		return false // Mempool Full (dropped)
	}

	n.Mempool[tx.ID] = tx
	// atomic.AddInt64(&n.TxCount, 1) // Do inside lock or use atomic outside?
	// Basic map insert is the heavy part we want to measure
	return true
}

func main() {
	fmt.Println("############################################################")
	fmt.Println("#       RNR MEMPOOL STRESS TEST: 3 NODES @ 10M TX/MIN      #")
	fmt.Println("############################################################")

	nodes := make([]*Node, NumNodes)
	for i := 0; i < NumNodes; i++ {
		nodes[i] = NewNode(i)
	}

	// Calculate Target Rate
	targetTPS := TargetTxPerMin / 60
	fmt.Printf("🎯 Target Load: %d TXs / minute\n", TargetTxPerMin)
	fmt.Printf("⚡ Target TPS:  %d TXs / second\n", targetTPS)
	fmt.Printf("⏱️  Duration:    %d seconds\n\n", DurationSec)

	var totalSent int64
	var totalAccepted int64
	var dropped int64

	start := time.Now()

	// Generators: 20 concurrent routines generating traffic
	var wg sync.WaitGroup
	quit := make(chan bool)

	// Firehose
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

			for {
				select {
				case <-quit:
					return
				default:
					// Create dummy TX
					var txID [32]byte
					// Fast random ID generation
					rng.Read(txID[:])

					tx := types.Transaction{
						ID:     txID,
						Amount: uint64(rng.Intn(1000)),
						Nonce:  uint64(rng.Intn(100000)),
					}

					// Send to a random node (Simulate RPC Load Balancing)
					targetNode := nodes[rng.Intn(NumNodes)]

					accepted := targetNode.AddTx(tx)

					atomic.AddInt64(&totalSent, 1)
					if accepted {
						atomic.AddInt64(&totalAccepted, 1)
						atomic.AddInt64(&targetNode.TxCount, 1)
					} else {
						atomic.AddInt64(&dropped, 1)
					}

					// Micro-sleep to avoid complete CPU lockup if needed,
					// but for stress test we want max speed.
					// time.Sleep(1 * time.Microsecond)
				}
			}
		}(i)
	}

	// Monitor & Timer
	ticker := time.NewTicker(1 * time.Second)
	stopTimer := time.NewTimer(time.Duration(DurationSec) * time.Second)

	fmt.Println("🚀 Starting Firehose...")

loop:
	for {
		select {
		case <-stopTimer.C:
			close(quit)
			break loop
		case <-ticker.C:
			currSent := atomic.LoadInt64(&totalSent)
			currAcc := atomic.LoadInt64(&totalAccepted)
			fmt.Printf("   [Status] Sent: %d | Accepted: %d | Mempools: [%d, %d, %d]\n",
				currSent, currAcc,
				len(nodes[0].Mempool), len(nodes[1].Mempool), len(nodes[2].Mempool))
		}
	}

	wg.Wait()
	duration := time.Since(start)

	// Results
	avgTPS := float64(totalAccepted) / duration.Seconds()
	projectedMinute := int64(avgTPS * 60)

	fmt.Println("\n############################################################")
	fmt.Println("#                    TEST RESULTS                          #")
	fmt.Println("############################################################")
	fmt.Printf("✅ Total Sent:      %d\n", totalSent)
	fmt.Printf("✅ Total Ingested:  %d (%.2f%%)\n", totalAccepted, (float64(totalAccepted)/float64(totalSent))*100)
	fmt.Printf("❌ Dropped (Full):  %d\n", dropped)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("⚡ Sustained TPS:   %.0f TX/s\n", avgTPS)
	fmt.Printf("📅 Projected/Min:   %d TX/min\n", projectedMinute)
	fmt.Printf("🎯 Target Met?:     ")

	if projectedMinute >= TargetTxPerMin {
		fmt.Printf("YES! (%.1fx Target)\n", float64(projectedMinute)/float64(TargetTxPerMin))
	} else {
		fmt.Printf("NO. (Reached %.1f%% of target)\n", (float64(projectedMinute)/float64(TargetTxPerMin))*100)
		fmt.Println("   Bottleneck likely: RAM Map Insertion Speed or CPU Mutex Contention")
	}

	fmt.Println("############################################################")
}
