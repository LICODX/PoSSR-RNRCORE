package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/p2p"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

func main() {
	target := flag.String("target", "localhost:4001", "Target node multiaddr")
	duration := flag.Duration("duration", 60*time.Second, "Test duration")
	workers := flag.Int("workers", 10, "Number of concurrent workers")
	tps := flag.Int("tps", 1000, "Target transactions per second")
	flag.Parse()

	fmt.Printf("🔥 Load Testing Tool\n")
	fmt.Printf("Target: %s\n", *target)
	fmt.Printf("Duration: %v\n", *duration)
	fmt.Printf("Workers: %d\n", *workers)
	fmt.Printf("Target TPS: %d\n\n", *tps)

	ctx := context.Background()

	// Create GossipSub node for testing
	node, err := p2p.NewGossipSubNode(ctx, 5000)
	if err != nil {
		fmt.Printf("Failed to create node: %v\n", err)
		return
	}
	defer node.Close()

	// Connect to target
	fmt.Printf("Connecting to %s...\n", *target)
	if err := node.ConnectToPeer(*target); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}

	// Stats
	var totalTx uint64
	var mu sync.Mutex
	stop := make(chan struct{})

	// Start workers
	for i := 0; i < *workers; i++ {
		go func(id int) {
			ticker := time.NewTicker(time.Second / time.Duration(*tps / *workers))
			defer ticker.Stop()

			for {
				select {
				case <-stop:
					return
				case <-ticker.C:
					// Create mock transaction
					tx := types.Transaction{
						Amount: 100,
						Nonce:  uint64(time.Now().UnixNano()),
					}

					// Serialize and send
					data := []byte(fmt.Sprintf("TX-%d-%d", id, tx.Nonce))
					if err := node.PublishTransaction(data); err != nil {
						fmt.Printf("Worker %d error: %v\n", id, err)
						continue
					}

					mu.Lock()
					totalTx++
					mu.Unlock()
				}
			}
		}(i)
	}

	// Report stats
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			mu.Lock()
			count := totalTx
			mu.Unlock()

			elapsed := time.Since(time.Now().Add(-5 * time.Second))
			actualTPS := float64(count) / elapsed.Seconds()

			fmt.Printf("📊 Stats: %d tx sent (%.2f TPS)\n", count, actualTPS)
		}
	}()

	// Run for duration
	time.Sleep(*duration)
	close(stop)

	mu.Lock()
	final := totalTx
	mu.Unlock()

	fmt.Printf("\n✅ Load test complete!\n")
	fmt.Printf("Total transactions: %d\n", final)
	fmt.Printf("Average TPS: %.2f\n", float64(final)/duration.Seconds())
}
