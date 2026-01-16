package metrics

import (
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Metrics
	blocksProduced = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rnr_blocks_produced_total",
		Help: "Total number of blocks produced",
	})

	transactionsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rnr_transactions_processed_total",
		Help: "Total number of transactions processed",
	})

	peerCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rnr_peer_count",
		Help: "Number of connected peers",
	})

	mempoolSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rnr_mempool_size",
		Help: "Current mempool size",
	})

	blockHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rnr_block_height",
		Help: "Current blockchain height",
	})

	consensusTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "rnr_consensus_duration_seconds",
		Help:    "Consensus duration in seconds",
		Buckets: prometheus.DefBuckets,
	})
)

func init() {
	// Register metrics
	prometheus.MustRegister(blocksProduced)
	prometheus.MustRegister(transactionsProcessed)
	prometheus.MustRegister(peerCount)
	prometheus.MustRegister(mempoolSize)
	prometheus.MustRegister(blockHeight)
	prometheus.MustRegister(consensusTime)
}

// StartPrometheusServer starts HTTP server for metrics
func StartPrometheusServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+port, nil)
}

// RecordBlock increments block counter
func RecordBlock(height uint64) {
	blocksProduced.Inc()
	blockHeight.Set(float64(height))
	atomic.AddUint64(&global.BlocksProduced, 1)
}

// RecordTransactions increments transaction counter
func RecordTransactions(count uint64) {
	transactionsProcessed.Add(float64(count))
	atomic.AddUint64(&global.TransactionsProc, count)
}

// UpdatePeerCount updates peer count gauge
func UpdatePeerCount(count int) {
	peerCount.Set(float64(count))
	atomic.StoreUint64(&global.PeerCount, uint64(count))
}

// UpdateMempoolSize updates mempool size gauge
func UpdateMempoolSize(size int) {
	mempoolSize.Set(float64(size))
	atomic.StoreUint64(&global.MempoolSize, uint64(size))
}

// RecordConsensusTime records consensus duration
func RecordConsensusTime(seconds float64) {
	consensusTime.Observe(seconds)
}
