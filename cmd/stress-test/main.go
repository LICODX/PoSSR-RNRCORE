package main

import (
	"flag"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

func main() {
	target := flag.String("target", "localhost:3000", "Target Node P2P Address")
	duration := flag.Duration("duration", 10*time.Second, "Test Duration")
	workers := flag.Int("workers", 50, "Number of concurrent workers")
	flag.Parse()

	fmt.Printf("🔥 Starting Stress Test on %s for %v\n", *target, *duration)

	var ops uint64
	stop := make(chan struct{})

	// Timer to stop test
	go func() {
		time.Sleep(*duration)
		close(stop)
	}()

	// Workers
	for i := 0; i < *workers; i++ {
		go func() {
			conn, err := net.Dial("tcp", *target)
			if err != nil {
				// Retry or ignore
				return
			}
			defer conn.Close()

			for {
				select {
				case <-stop:
					return
				default:
					// Send Garbage Tx Data
					payload := []byte("{\"jsonrpc\":\"2.0\",\"method\":\"add_tx\",\"params\":[\"stress_test\"]}\n")
					_, err := conn.Write(payload)
					if err == nil {
						atomic.AddUint64(&ops, 1)
					} else {
						// Reconnect if broken
						conn.Close()
						conn, _ = net.Dial("tcp", *target)
					}
					// Small sleep to not kill OS socket limit
					time.Sleep(1 * time.Millisecond)
				}
			}
		}()
	}

	<-stop
	fmt.Println("\n🛑 Stress Test Finished.")
	fmt.Printf("Total Requests: %d\n", ops)
	fmt.Printf("Avg TPS: %.2f\n", float64(ops)/duration.Seconds())
}
