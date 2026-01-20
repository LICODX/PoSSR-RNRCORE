package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// StressTest simulates high transaction load on RNR network
type StressTest struct {
	targetURL  string
	tps        int
	duration   time.Duration
	concurrent int

	// Metrics
	totalSent    uint64
	totalSuccess uint64
	totalFailed  uint64
	latencies    []time.Duration
	mu           sync.Mutex
}

// Transaction represents a test transaction
type Transaction struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Fee    float64 `json:"fee"`
}

// NewStressTest creates a new stress tester
func NewStressTest(url string, tps int, duration time.Duration, concurrent int) *StressTest {
	return &StressTest{
		targetURL:  url,
		tps:        tps,
		duration:   duration,
		concurrent: concurrent,
		latencies:  make([]time.Duration, 0, 10000),
	}
}

// Run executes the stress test
func (st *StressTest) Run() {
	fmt.Printf("🔥 Starting Stress Test\n")
	fmt.Printf("   Target: %s\n", st.targetURL)
	fmt.Printf("   TPS: %d\n", st.tps)
	fmt.Printf("   Duration: %v\n", st.duration)
	fmt.Printf("   Concurrent Workers: %d\n\n", st.concurrent)

	// Calculate transactions per worker
	txPerWorker := st.tps / st.concurrent
	interval := time.Second / time.Duration(txPerWorker)

	// Start workers
	var wg sync.WaitGroup
	endTime := time.Now().Add(st.duration)

	for i := 0; i < st.concurrent; i++ {
		wg.Add(1)
		go st.worker(i, interval, endTime, &wg)
	}

	// Progress reporter
	go st.reportProgress(endTime)

	// Wait for completion
	wg.Wait()

	// Final report
	st.printReport()
}

// worker sends transactions at specified rate
func (st *StressTest) worker(id int, interval time.Duration, endTime time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Now().After(endTime) {
				return
			}
			st.sendTransaction()
		}
	}
}

// sendTransaction sends a single test transaction
func (st *StressTest) sendTransaction() {
	atomic.AddUint64(&st.totalSent, 1)

	// Generate random recipient
	recipient := st.generateAddress()

	tx := Transaction{
		To:     recipient,
		Amount: 1.0,
		Fee:    1.0,
	}

	data, _ := json.Marshal(tx)

	start := time.Now()
	resp, err := http.Post(
		st.targetURL+"/api/wallet/send",
		"application/json",
		nil,
	)
	latency := time.Since(start)

	st.mu.Lock()
	st.latencies = append(st.latencies, latency)
	st.mu.Unlock()

	if err != nil || (resp != nil && resp.StatusCode != 200) {
		atomic.AddUint64(&st.totalFailed, 1)
		if resp != nil {
			resp.Body.Close()
		}
		return
	}

	atomic.AddUint64(&st.totalSuccess, 1)
	if resp != nil {
		resp.Body.Close()
	}
}

// generateAddress creates a random test address
func (st *StressTest) generateAddress() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "rnr1" + hex.EncodeToString(b)[:60]
}

// reportProgress prints real-time progress
func (st *StressTest) reportProgress(endTime time.Time) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		if time.Now().After(endTime) {
			return
		}

		<-ticker.C

		sent := atomic.LoadUint64(&st.totalSent)
		success := atomic.LoadUint64(&st.totalSuccess)
		failed := atomic.LoadUint64(&st.totalFailed)

		successRate := float64(success) / float64(sent) * 100
		currentTPS := float64(sent) / time.Since(time.Now().Add(-st.duration)).Seconds()

		fmt.Printf("📊 Progress: Sent=%d Success=%d (%.1f%%) Failed=%d TPS=%.0f\n",
			sent, success, successRate, failed, currentTPS)
	}
}

// printReport shows final test results
func (st *StressTest) printReport() {
	fmt.Println("\n" + "="*60)
	fmt.Println("📈 STRESS TEST RESULTS")
	fmt.Println("=" * 60)

	sent := atomic.LoadUint64(&st.totalSent)
	success := atomic.LoadUint64(&st.totalSuccess)
	failed := atomic.LoadUint64(&st.totalFailed)

	fmt.Printf("\n📤 Transactions Sent:     %d\n", sent)
	fmt.Printf("✅ Successful:            %d (%.2f%%)\n", success, float64(success)/float64(sent)*100)
	fmt.Printf("❌ Failed:                %d (%.2f%%)\n", failed, float64(failed)/float64(sent)*100)

	// Calculate latency statistics
	st.mu.Lock()
	latencies := st.latencies
	st.mu.Unlock()

	if len(latencies) > 0 {
		var total, min, max time.Duration
		min = latencies[0]
		max = latencies[0]

		for _, lat := range latencies {
			total += lat
			if lat < min {
				min = lat
			}
			if lat > max {
				max = lat
			}
		}

		avg := total / time.Duration(len(latencies))

		fmt.Printf("\n⏱️  Latency:\n")
		fmt.Printf("   Average:              %v\n", avg)
		fmt.Printf("   Min:                  %v\n", min)
		fmt.Printf("   Max:                  %v\n", max)
	}

	// Calculate actual TPS
	actualTPS := float64(sent) / st.duration.Seconds()
	fmt.Printf("\n🚀 Performance:\n")
	fmt.Printf("   Target TPS:           %d\n", st.tps)
	fmt.Printf("   Actual TPS:           %.2f\n", actualTPS)
	fmt.Printf("   Efficiency:           %.1f%%\n", (actualTPS/float64(st.tps))*100)

	fmt.Println("\n" + "="*60)
}

func main() {
	// Command-line flags
	url := flag.String("url", "http://localhost:8080", "RNR node URL")
	tps := flag.Int("tps", 100, "Target transactions per second")
	duration := flag.Duration("duration", 60*time.Second, "Test duration")
	concurrent := flag.Int("concurrent", 10, "Number of concurrent workers")

	flag.Parse()

	// Create and run stress test
	test := NewStressTest(*url, *tps, *duration, *concurrent)
	test.Run()
}
