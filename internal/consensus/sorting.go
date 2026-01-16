package consensus

import (
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// SortableTransaction wraps a transaction with its sorting key
type SortableTransaction struct {
	Tx  types.Transaction
	Key string // MixHash(TxID, Seed)
}

// ByKey implements sort.Interface for []SortableTransaction based on Key field
type ByKey []SortableTransaction

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// =============================================================================
// 1. QUICK SORT - Divide and Conquer
// =============================================================================

// QuickSort implements the quicksort algorithm
// Average: O(n log n), Worst: O(n²), Space: O(log n)
func QuickSort(data []SortableTransaction) []SortableTransaction {
	if len(data) <= 1 {
		return data
	}

	result := make([]SortableTransaction, len(data))
	copy(result, data)
	quickSortRecursive(result, 0, len(result)-1)
	return result
}

func quickSortRecursive(arr []SortableTransaction, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSortRecursive(arr, low, pi-1)
		quickSortRecursive(arr, pi+1, high)
	}
}

func partition(arr []SortableTransaction, low, high int) int {
	pivot := arr[high].Key
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j].Key < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// =============================================================================
// 2. MERGE SORT - Stable Divide and Conquer
// =============================================================================

// MergeSort implements the merge sort algorithm
// Complexity: O(n log n) guaranteed, Space: O(n)
// Stable sorting algorithm
func MergeSort(data []SortableTransaction) []SortableTransaction {
	if len(data) <= 1 {
		return data
	}

	result := make([]SortableTransaction, len(data))
	copy(result, data)
	return mergeSortRecursive(result)
}

func mergeSortRecursive(arr []SortableTransaction) []SortableTransaction {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := mergeSortRecursive(arr[:mid])
	right := mergeSortRecursive(arr[mid:])

	return merge(left, right)
}

func merge(left, right []SortableTransaction) []SortableTransaction {
	result := make([]SortableTransaction, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].Key <= right[j].Key {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

// =============================================================================
// 3. HEAP SORT - In-Place Selection Sort Variant
// =============================================================================

// HeapSort implements the heap sort algorithm
// Complexity: O(n log n) guaranteed, Space: O(1)
// In-place sorting, not stable
func HeapSort(data []SortableTransaction) []SortableTransaction {
	result := make([]SortableTransaction, len(data))
	copy(result, data)

	n := len(result)

	// Build max heap
	for i := n/2 - 1; i >= 0; i-- {
		heapify(result, n, i)
	}

	// Extract elements from heap one by one
	for i := n - 1; i > 0; i-- {
		result[0], result[i] = result[i], result[0]
		heapify(result, i, 0)
	}

	return result
}

func heapify(arr []SortableTransaction, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[left].Key > arr[largest].Key {
		largest = left
	}

	if right < n && arr[right].Key > arr[largest].Key {
		largest = right
	}

	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

// =============================================================================
// 4. RADIX SORT - Non-Comparison Based (for strings)
// =============================================================================

// RadixSort implements radix sort for string keys
// Complexity: O(d*n) where d is max key length, Space: O(n+k)
// Stable, non-comparison based
func RadixSort(data []SortableTransaction) []SortableTransaction {
	if len(data) <= 1 {
		return data
	}

	result := make([]SortableTransaction, len(data))
	copy(result, data)

	// Find max length
	maxLen := 0
	for _, item := range result {
		if len(item.Key) > maxLen {
			maxLen = len(item.Key)
		}
	}

	// Sort from least significant to most significant character
	for pos := maxLen - 1; pos >= 0; pos-- {
		result = countingSortByChar(result, pos)
	}

	return result
}

func countingSortByChar(arr []SortableTransaction, pos int) []SortableTransaction {
	n := len(arr)
	output := make([]SortableTransaction, n)
	count := make([]int, 256) // ASCII characters

	// Count occurrences
	for _, item := range arr {
		char := getCharAt(item.Key, pos)
		count[char]++
	}

	// Calculate cumulative count
	for i := 1; i < 256; i++ {
		count[i] += count[i-1]
	}

	// Build output array (backwards for stability)
	for i := n - 1; i >= 0; i-- {
		char := getCharAt(arr[i].Key, pos)
		output[count[char]-1] = arr[i]
		count[char]--
	}

	return output
}

func getCharAt(s string, pos int) byte {
	if pos >= len(s) {
		return 0 // Null character for padding
	}
	return s[pos]
}

// =============================================================================
// 5. TIM SORT - Hybrid (Merge + Insertion)
// =============================================================================

const minMerge = 32

// TimSort implements a simplified version of Timsort
// Complexity: O(n log n), Space: O(n)
// Stable, optimized for real-world data (used in Python)
func TimSort(data []SortableTransaction) []SortableTransaction {
	result := make([]SortableTransaction, len(data))
	copy(result, data)

	n := len(result)

	// Sort individual subarrays of size minMerge using insertion sort
	for start := 0; start < n; start += minMerge {
		end := start + minMerge
		if end > n {
			end = n
		}
		insertionSort(result, start, end-1)
	}

	// Merge sorted runs
	size := minMerge
	for size < n {
		for start := 0; start < n; start += size * 2 {
			mid := start + size - 1
			end := start + size*2 - 1
			if end > n-1 {
				end = n - 1
			}
			if mid < end {
				mergeRuns(result, start, mid, end)
			}
		}
		size *= 2
	}

	return result
}

func insertionSort(arr []SortableTransaction, left, right int) {
	for i := left + 1; i <= right; i++ {
		key := arr[i]
		j := i - 1
		for j >= left && arr[j].Key > key.Key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func mergeRuns(arr []SortableTransaction, left, mid, right int) {
	len1 := mid - left + 1
	len2 := right - mid

	leftArr := make([]SortableTransaction, len1)
	rightArr := make([]SortableTransaction, len2)

	copy(leftArr, arr[left:left+len1])
	copy(rightArr, arr[mid+1:mid+1+len2])

	i, j, k := 0, 0, left

	for i < len1 && j < len2 {
		if leftArr[i].Key <= rightArr[j].Key {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}

	for i < len1 {
		arr[k] = leftArr[i]
		i++
		k++
	}

	for j < len2 {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

// =============================================================================
// 6. INTRO SORT - Hybrid (Quick + Heap)
// =============================================================================

// IntroSort implements introspective sort (used in C++ STL)
// Complexity: O(n log n) guaranteed, Space: O(log n)
// Starts with quicksort, switches to heapsort if depth limit exceeded
func IntroSort(data []SortableTransaction) []SortableTransaction {
	result := make([]SortableTransaction, len(data))
	copy(result, data)

	maxDepth := 2 * logBase2(len(result))
	introSortRecursive(result, 0, len(result)-1, maxDepth)

	return result
}

func introSortRecursive(arr []SortableTransaction, low, high, depthLimit int) {
	for high-low > 16 { // Use insertion sort for small arrays
		if depthLimit == 0 {
			// Switch to heapsort
			heapSortRange(arr, low, high)
			return
		}
		depthLimit--

		// Partition using quicksort
		p := partition(arr, low, high)

		// Recur on smaller partition, iterate on larger
		if p-low < high-p {
			introSortRecursive(arr, low, p-1, depthLimit)
			low = p + 1
		} else {
			introSortRecursive(arr, p+1, high, depthLimit)
			high = p - 1
		}
	}

	// Use insertion sort for small subarrays
	insertionSort(arr, low, high)
}

func heapSortRange(arr []SortableTransaction, low, high int) {
	n := high - low + 1

	// Build heap
	for i := n/2 - 1; i >= 0; i-- {
		heapifyRange(arr, low, n, i)
	}

	// Extract elements
	for i := n - 1; i > 0; i-- {
		arr[low], arr[low+i] = arr[low+i], arr[low]
		heapifyRange(arr, low, i, 0)
	}
}

func heapifyRange(arr []SortableTransaction, offset, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[offset+left].Key > arr[offset+largest].Key {
		largest = left
	}

	if right < n && arr[offset+right].Key > arr[offset+largest].Key {
		largest = right
	}

	if largest != i {
		arr[offset+i], arr[offset+largest] = arr[offset+largest], arr[offset+i]
		heapifyRange(arr, offset, n, largest)
	}
}

func logBase2(n int) int {
	log := 0
	for n > 1 {
		n >>= 1
		log++
	}
	return log
}
