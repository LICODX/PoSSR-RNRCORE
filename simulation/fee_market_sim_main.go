package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Config
const (
	RealBlockTime = 60 * time.Second
	MaxTxPerBlock = 1_000_000 // UNLEASHED LIMIT
	MempoolCap    = 2_000_000
)

// Priority Queue for Txs based on Fee
type TxHeap []types.Transaction

func (h TxHeap) Len() int           { return len(h) }
func (h TxHeap) Less(i, j int) bool { return h[i].Fee > h[j].Fee } // Max-Heap (Highest Fee First)
func (h TxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *TxHeap) Push(x interface{}) {
	*h = append(*h, x.(types.Transaction))
}
func (h *TxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	fmt.Println("############################################################")
	fmt.Println("#      RNR MAINNET: UNLEASHED CAPACITY TEST (1 MIN)        #")
	fmt.Println("############################################################")
	fmt.Println("Phase: Stress Testing Maximum Throughput...")
	fmt.Println("Block Time: 60 Seconds")
	fmt.Printf("Block Cap:  %d TXs\n", MaxTxPerBlock)
	fmt.Println("------------------------------------------------------------")

	// Init Mempool (Heap)
	mempool := &TxHeap{}
	heap.Init(mempool)

	// Channels - Buffer increased for high throughput
	txChan := make(chan types.Transaction, 50000)
	stopChan := make(chan bool)

	// 1. Traffic Generator (Scaled Up)
	// Spawning 10 Parallel Generators to flood the mempool
	for i := 0; i < 10; i++ {
		go func(id int) {
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

			for {
				select {
				case <-stopChan:
					return
				default:
					feeBase := rng.Float64()
					var fee uint64
					if feeBase > 0.99 {
						fee = uint64(1000 + rng.Intn(5000))
					} else if feeBase > 0.8 {
						fee = uint64(100 + rng.Intn(900))
					} else {
						fee = uint64(1 + rng.Intn(50))
					}

					tx := types.Transaction{
						Fee:    fee,
						Amount: uint64(rng.Intn(1000)),
						Nonce:// Random nonce
						uint64(rng.Int63()),
					}
					txChan <- tx
					// NO SLEEP. MAX SPEED.
				}
			}
		}(i)
	}

	// 2. Node Main Loop
	blockTimer := time.NewTicker(RealBlockTime)
	blockHeight := 1

	totalRevenue := uint64(0)

	fmt.Printf("[%s] ⏳ Mining Genesis Block (Waiting 60s)...\n", time.Now().Format("15:04:05"))

	for {
		select {
		case tx := <-txChan:
			// Ingest to Mempool
			if len(*mempool) < MempoolCap {
				heap.Push(mempool, tx)
			}
			// Drops if full (implicit)

		case <-blockTimer.C:
			// BLOCK TIME!
			fmt.Printf("\n[%s] 🔨 PROPOSING BLOCK #%d\n", time.Now().Format("15:04:05"), blockHeight)

			// Pop best transactions
			blockTxs := make([]types.Transaction, 0, MaxTxPerBlock)
			blockFees := uint64(0)
			minFee := uint64(9999999)
			maxFee := uint64(0)

			count := 0

			for mempool.Len() > 0 && count < MaxTxPerBlock {
				tx := heap.Pop(mempool).(types.Transaction)
				blockTxs = append(blockTxs, tx)
				blockFees += tx.Fee
				if tx.Fee < minFee {
					minFee = tx.Fee
				}
				if tx.Fee > maxFee {
					maxFee = tx.Fee
				}
				count++
			}

			// Statistics
			fmt.Println("   📦 Block Fullness: ", count, "/", MaxTxPerBlock)
			if count > 0 {
				fmt.Printf("   💰 Block Revenue:  %d RNR (Fees)\n", blockFees)
				fmt.Printf("   📈 Fee Range:      Min %d | Max %d\n", minFee, maxFee)
				fmt.Printf("   🌊 Mempool Backlog:%d Pending Txs\n", mempool.Len())

				avgFee := blockFees / uint64(count)
				fmt.Printf("   📊 Avg Fee:        %d\n", avgFee)
			}

			totalRevenue += blockFees

			// Stop after 1 block only (to prove point quickly)
			if blockHeight >= 1 {
				close(stopChan)
				fmt.Println("\n------------------------------------------------------------")
				fmt.Println("Simulation Complete.")
				fmt.Printf("Total TPS Reached: %.2f TX/s\n", float64(count)/60.0)
				fmt.Printf("Total In Block:    %d TXs\n", count)
				return
			}
			blockHeight++
		}
	}
}
