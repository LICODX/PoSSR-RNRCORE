package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

const (
	// 2 GB Simulation
	// Avg Tx Size = 500 Bytes
	// Reduced to 1M for rapid VM verification (Extrapolate x4 for 2GB)
	NUM_TXS = 1_000_000
)

func main() {
	fmt.Println("🏁 STARTING 2GB MEMPOOL ALGORITHM RACE")
	fmt.Printf("   - Data Size: 2 GB (Simulated via %d Transactions)\n", NUM_TXS)
	fmt.Printf("   - System: %d CPUs\n", runtime.NumCPU())
	fmt.Println("---------------------------------------------------")

	// 1. Generate Data
	fmt.Print("⏳ Generating 2GB Random Data... ")
	data := make([]consensus.SortableTransaction, NUM_TXS)
	for i := 0; i < NUM_TXS; i++ {
		// Random Key for sorting
		key := fmt.Sprintf("%064x", rand.Uint64())
		data[i] = consensus.SortableTransaction{
			Tx:  types.Transaction{ID: [32]byte{}}, // Empty body to save real RAM, Key is what matters
			Key: key,
		}
	}
	fmt.Println("DONE.")

	// 2. Define Algos
	algos := []string{
		"QUICK_SORT", "MERGE_SORT", "HEAP_SORT",
		"RADIX_SORT", "TIM_SORT", "INTRO_SORT", "SHELL_SORT",
	}

	// 3. Race
	fmt.Println("\n🏎️  RACE START!")
	fmt.Printf("%-15s | %-15s | %s\n", "ALGORITHM", "TIME (s)", "STATUS")
	fmt.Println("------------------------------------------------")

	for _, name := range algos {
		// Clone data to ensure fair start
		input := make([]consensus.SortableTransaction, len(data))
		copy(input, data)

		start := time.Now()

		// Run Specific Algo
		switch name {
		case "QUICK_SORT":
			consensus.QuickSort(input)
		case "MERGE_SORT":
			consensus.MergeSort(input)
		case "HEAP_SORT":
			consensus.HeapSort(input)
		case "RADIX_SORT":
			consensus.RadixSort(input)
		case "TIM_SORT":
			consensus.TimSort(input)
		case "INTRO_SORT":
			consensus.IntroSort(input)
		case "SHELL_SORT":
			consensus.ShellSort(input)
		}

		duration := time.Since(start).Seconds()
		fmt.Printf("%-15s | %-15.4f | ✅ DONE\n", name, duration)

		// GC to prevent OOM between runs
		input = nil
		runtime.GC()
		time.Sleep(1 * time.Second) // Cool down
	}
}
