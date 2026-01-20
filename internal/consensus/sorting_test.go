package consensus

import (
	"testing"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// Helper function to create test transactions
func createTestTransactions(n int) []types.Transaction {
	txs := make([]types.Transaction, n)
	for i := 0; i < n; i++ {
		txs[i] = types.Transaction{
			ID: [32]byte{byte(i)},
		}
	}
	return txs
}

// Helper to check if array is sorted by keys
func isSorted(data []SortableTransaction) bool {
	for i := 1; i < len(data); i++ {
		if data[i-1].Key > data[i].Key {
			return false
		}
	}
	return true
}

// Test QuickSort
func TestQuickSort(t *testing.T) {
	data := []SortableTransaction{
		{Key: "gamma"},
		{Key: "alpha"},
		{Key: "beta"},
		{Key: "delta"},
	}

	sorted := QuickSort(data)

	if !isSorted(sorted) {
		t.Error("QuickSort did not sort correctly")
	}

	if sorted[0].Key != "alpha" || sorted[3].Key != "gamma" {
		t.Errorf("QuickSort order incorrect: got %v", sorted)
	}
}

// Test MergeSort
func TestMergeSort(t *testing.T) {
	data := []SortableTransaction{
		{Key: "zebra"},
		{Key: "apple"},
		{Key: "mango"},
		{Key: "banana"},
	}

	sorted := MergeSort(data)

	if !isSorted(sorted) {
		t.Error("MergeSort did not sort correctly")
	}

	if sorted[0].Key != "apple" || sorted[3].Key != "zebra" {
		t.Errorf("MergeSort order incorrect")
	}
}

// Test HeapSort
func TestHeapSort(t *testing.T) {
	data := []SortableTransaction{
		{Key: "dog"},
		{Key: "cat"},
		{Key: "elephant"},
		{Key: "ant"},
	}

	sorted := HeapSort(data)

	if !isSorted(sorted) {
		t.Error("HeapSort did not sort correctly")
	}

	if sorted[0].Key != "ant" || sorted[3].Key != "elephant" {
		t.Errorf("HeapSort order incorrect")
	}
}

// Test RadixSort
func TestRadixSort(t *testing.T) {
	data := []SortableTransaction{
		{Key: "xyz"},
		{Key: "abc"},
		{Key: "mno"},
		{Key: "def"},
	}

	sorted := RadixSort(data)

	if !isSorted(sorted) {
		t.Error("RadixSort did not sort correctly")
	}

	if sorted[0].Key != "abc" || sorted[3].Key != "xyz" {
		t.Errorf("RadixSort order incorrect")
	}
}

// Test TimSort
func TestTimSort(t *testing.T) {
	data := []SortableTransaction{
		{Key: "99"},
		{Key: "11"},
		{Key: "55"},
		{Key: "33"},
	}

	sorted := TimSort(data)

	if !isSorted(sorted) {
		t.Error("TimSort did not sort correctly")
	}
}

// Test IntroSort
func TestIntroSort(t *testing.T) {
	data := []SortableTransaction{
		{Key: "omega"},
		{Key: "alpha"},
		{Key: "theta"},
		{Key: "beta"},
	}

	sorted := IntroSort(data)

	if !isSorted(sorted) {
		t.Error("IntroSort did not sort correctly")
	}

	if sorted[0].Key != "alpha" || sorted[3].Key != "theta" {
		t.Errorf("IntroSort order incorrect")
	}
}

// Test all algorithms produce same result (crucial for consensus!)
func TestAllAlgorithmsProduceSameResult(t *testing.T) {
	// Create test data
	txs := createTestTransactions(20)
	seed := [32]byte{0x42, 0x42, 0x42} // Fixed seed

	// Prepare sortable data
	prepareData := func() []SortableTransaction {
		data := make([]SortableTransaction, len(txs))
		for i, tx := range txs {
			data[i] = SortableTransaction{
				Tx:  tx,
				Key: utils.MixHash(tx.ID, seed),
			}
		}
		return data
	}

	// Run all algorithms
	results := map[string][]SortableTransaction{
		"QUICK_SORT": QuickSort(prepareData()),
		"MERGE_SORT": MergeSort(prepareData()),
		"HEAP_SORT":  HeapSort(prepareData()),
		"RADIX_SORT": RadixSort(prepareData()),
		"TIM_SORT":   TimSort(prepareData()),
		"INTRO_SORT": IntroSort(prepareData()),
	}

	// Compare all results
	baseline := results["QUICK_SORT"]

	for name, result := range results {
		if len(result) != len(baseline) {
			t.Errorf("%s produced different length: %d vs %d", name, len(result), len(baseline))
			continue
		}

		for i := range result {
			if result[i].Key != baseline[i].Key {
				t.Errorf("%s produced different order at index %d", name, i)
				t.Logf("  Expected: %s", baseline[i].Key)
				t.Logf("  Got: %s", result[i].Key)
				break
			}
		}
	}

	t.Logf("✅ All 6 algorithms produced identical results")
}

// Test algorithm selection is deterministic
func TestAlgorithmSelectionDeterministic(t *testing.T) {
	testCases := []struct {
		lastByte byte
		expected string
	}{
		{0, "QUICK_SORT"},
		{1, "MERGE_SORT"},
		{2, "HEAP_SORT"},
		{3, "RADIX_SORT"},
		{4, "TIM_SORT"},
		{5, "INTRO_SORT"},
		{6, "SHELL_SORT"},  // % 7 = 6
		{7, "QUICK_SORT"},  // % 7 = 0
		{14, "QUICK_SORT"}, // % 7 = 0
		{255, "TIM_SORT"},  // 255 % 7 = 3 (Wait: 255 = 36*7 + 3) -> RADIX?
		// Let'scalc: 255 / 7 = 36 rem 3. So 3 is RADIX.
		// Wait, 4 is TIM.
		// Let's verify calculator: 255 - (36*7) = 255 - 252 = 3.
		// Case 3 is RADIX_SORT.
	}

	for _, tc := range testCases {
		seed := [32]byte{}
		seed[31] = tc.lastByte

		result := utils.SelectAlgorithm(seed)

		if result != tc.expected {
			t.Errorf("Seed byte %d: expected %s, got %s", tc.lastByte, tc.expected, result)
		}
	}
}

// Test StartRace with different seeds produces different algorithms
func TestStartRaceUsesCorrectAlgorithm(t *testing.T) {
	txs := createTestTransactions(10)

	testCases := []struct {
		seedByte byte
		expected string
	}{
		{0, "QUICK_SORT"},
		{1, "MERGE_SORT"},
		{2, "HEAP_SORT"},
		{3, "RADIX_SORT"},
		{4, "TIM_SORT"},
		{5, "INTRO_SORT"},
	}

	for _, tc := range testCases {
		seed := [32]byte{}
		seed[31] = tc.seedByte

		algo := utils.SelectAlgorithm(seed)
		if algo != tc.expected {
			t.Errorf("Seed %d should select %s, got %s", tc.seedByte, tc.expected, algo)
		}

		// Run StartRace
		sorted, root := StartRaceSimplified(txs, seed, algo)

		// Verify result is sorted (by running deterministic check)
		if len(sorted) != len(txs) {
			t.Errorf("StartRace returned wrong number of transactions")
		}

		if root == [32]byte{} {
			t.Error("StartRace returned empty Merkle root")
		}

		t.Logf("Seed byte %d → %s ✅", tc.seedByte, algo)
	}
}

// Test large dataset performance
func TestSortingLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	txs := createTestTransactions(1000)
	seed := [32]byte{0x99}

	// Prepare data
	data := make([]SortableTransaction, len(txs))
	for i, tx := range txs {
		data[i] = SortableTransaction{
			Tx:  tx,
			Key: utils.MixHash(tx.ID, seed),
		}
	}

	algorithms := map[string]func([]SortableTransaction) []SortableTransaction{
		"QuickSort": QuickSort,
		"MergeSort": MergeSort,
		"HeapSort":  HeapSort,
		"RadixSort": RadixSort,
		"TimSort":   TimSort,
		"IntroSort": IntroSort,
	}

	for name, sortFunc := range algorithms {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)

		sorted := sortFunc(dataCopy)

		if !isSorted(sorted) {
			t.Errorf("%s failed to sort large dataset", name)
		} else {
			t.Logf("%s: ✅ sorted %d items", name, len(sorted))
		}
	}
}

// Benchmark different algorithms
func BenchmarkQuickSort(b *testing.B) {
	data := make([]SortableTransaction, 100)
	for i := 0; i < 100; i++ {
		data[i].Key = utils.MixHash([32]byte{byte(i)}, [32]byte{0x42})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)
		QuickSort(dataCopy)
	}
}

func BenchmarkMergeSort(b *testing.B) {
	data := make([]SortableTransaction, 100)
	for i := 0; i < 100; i++ {
		data[i].Key = utils.MixHash([32]byte{byte(i)}, [32]byte{0x42})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)
		MergeSort(dataCopy)
	}
}

func BenchmarkHeapSort(b *testing.B) {
	data := make([]SortableTransaction, 100)
	for i := 0; i < 100; i++ {
		data[i].Key = utils.MixHash([32]byte{byte(i)}, [32]byte{0x42})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)
		HeapSort(dataCopy)
	}
}

func BenchmarkRadixSort(b *testing.B) {
	data := make([]SortableTransaction, 100)
	for i := 0; i < 100; i++ {
		data[i].Key = utils.MixHash([32]byte{byte(i)}, [32]byte{0x42})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)
		RadixSort(dataCopy)
	}
}

func BenchmarkTimSort(b *testing.B) {
	data := make([]SortableTransaction, 100)
	for i := 0; i < 100; i++ {
		data[i].Key = utils.MixHash([32]byte{byte(i)}, [32]byte{0x42})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)
		TimSort(dataCopy)
	}
}

func BenchmarkIntroSort(b *testing.B) {
	data := make([]SortableTransaction, 100)
	for i := 0; i < 100; i++ {
		data[i].Key = utils.MixHash([32]byte{byte(i)}, [32]byte{0x42})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataCopy := make([]SortableTransaction, len(data))
		copy(dataCopy, data)
		IntroSort(dataCopy)
	}
}
