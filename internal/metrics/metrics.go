package metrics

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Metrics tracks system metrics
type Metrics struct {
	BlocksProduced   uint64
	TransactionsProc uint64
	PeerCount        uint64
	MempoolSize      uint64
	StartTime        time.Time
}

var global = &Metrics{
	StartTime: time.Now(),
}

// IncrementBlocks increments block counter
func IncrementBlocks() {
	atomic.AddUint64(&global.BlocksProduced, 1)
}

// IncrementTransactions increments transaction counter
func IncrementTransactions(count uint64) {
	atomic.AddUint64(&global.TransactionsProc, count)
}

// SetPeerCount sets current peer count
func SetPeerCount(count uint64) {
	atomic.StoreUint64(&global.PeerCount, count)
}

// SetMempoolSize sets mempool size
func SetMempoolSize(size uint64) {
	atomic.StoreUint64(&global.MempoolSize, size)
}

// GetStats returns current statistics
func GetStats() map[string]interface{} {
	uptime := time.Since(global.StartTime).Seconds()
	blocksPerSec := float64(atomic.LoadUint64(&global.BlocksProduced)) / uptime

	return map[string]interface{}{
		"blocks_produced":    atomic.LoadUint64(&global.BlocksProduced),
		"transactions_total": atomic.LoadUint64(&global.TransactionsProc),
		"peer_count":         atomic.LoadUint64(&global.PeerCount),
		"mempool_size":       atomic.LoadUint64(&global.MempoolSize),
		"uptime_seconds":     uptime,
		"blocks_per_second":  blocksPerSec,
	}
}

// PrintStats prints metrics to console
func PrintStats() {
	stats := GetStats()
	fmt.Println("\n📊 Metrics:")
	fmt.Printf("  Blocks: %d (%.2f/sec)\n", stats["blocks_produced"], stats["blocks_per_second"])
	fmt.Printf("  Transactions: %d\n", stats["transactions_total"])
	fmt.Printf("  Peers: %d\n", stats["peer_count"])
	fmt.Printf("  Mempool: %d\n", stats["mempool_size"])
	fmt.Printf("  Uptime: %.0fs\n", stats["uptime_seconds"])
}
