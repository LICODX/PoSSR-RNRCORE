package consensus

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// MineBlock runs the Proof of Repeated Sorting (PoRS)
// It continuously hashes and sorts until a valid block hash < target is found.
func MineBlock(mempool []types.Transaction, prevBlock types.BlockHeader, difficulty uint64, stopChan chan struct{}) (*types.Block, error) {
	// 1. Prepare base data
	var nonce uint64 = 0
	// target := big.NewInt(0).SetUint64(difficulty)
	// Note: In real PoW, target is huge (2^256 / diff). For simulation, we use small numbers.
	// Let's assume higher difficulty = harder to find (target is smaller).
	// Current mock: Target is just a threshold check.
	// REAL POW: Hash < (MaxHash / Difficulty)

	// Create simplified target for PoSSR:
	// We want the Hash to start with N zeroes.

	fmt.Printf("⛏️  Mining started. Difficulty: %d\n", difficulty)

	for {
		// Check for interrupt
		select {
		case <-stopChan:
			return nil, fmt.Errorf("mining interrupted")
		default:
		}

		// 2. Generate Seed for this Nonce
		// Seed = SHA256(PrevHash + Nonce)
		hasher := sha256.New()
		hasher.Write(prevBlock.Hash[:])
		binary.Write(hasher, binary.BigEndian, nonce)
		seedBytes := hasher.Sum(nil)
		var seed [32]byte
		copy(seed[:], seedBytes)

		// 3. PoSSR: Execute Sorting Race with this specific seed
		// This consumes CPU.
		sortedTxs, merkleRoot := StartRaceSimplified(mempool, seed)

		// 4. Construct Candidate Header
		header := types.BlockHeader{
			Version:       1,
			PrevBlockHash: prevBlock.Hash,
			MerkleRoot:    merkleRoot,
			Timestamp:     time.Now().Unix(),
			Height:        prevBlock.Height + 1,
			Nonce:         nonce,
			Difficulty:    difficulty,
			VRFSeed:       seed, // Seed used for sorting becomes the VRF seed
		}

		// 5. Calculate Block Hash
		blockHash := CalculateBlockHash(header)
		header.Hash = blockHash

		// 6. Check against Target (Proof of Work check)
		// Convert hash to big int
		hashInt := new(big.Int).SetBytes(blockHash[:])

		// Real Check: hashInt < (2^256 / difficulty)
		// For demo simplicity: We check if hash ends with '0' bytes based on difficulty
		// Actually, let's just use the standard check.
		maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
		targetVal := new(big.Int).Div(maxVal, big.NewInt(int64(difficulty)))

		if hashInt.Cmp(targetVal) == -1 {
			// SUCCESS!
			// fmt.Printf("💎 Block Found! Nonce: %d | Hash: %x\n", nonce, blockHash)

			// Construct Full Block
			block := &types.Block{
				Header: header,
			}
			// Fill shards (simplified)
			block.Shards[0] = types.ShardData{
				TxData: sortedTxs,
			}

			return block, nil
		}

		// 7. Increment Nonce and try again
		nonce++
	}
}

// StartRaceSimplified runs the sorting logic for the mining loop
func StartRaceSimplified(mempool []types.Transaction, seed [32]byte) ([]types.Transaction, [32]byte) {
	// Re-use existing logic, but made internal
	algo := SelectAlgorithm(seed)

	sortableData := make([]SortableTransaction, len(mempool))
	for i, tx := range mempool {
		sortableData[i] = SortableTransaction{
			Tx:  tx,
			Key: MixHash(tx.ID, seed),
		}
	}

	var sorted []SortableTransaction
	switch algo {
	case "QUICK_SORT":
		sorted = QuickSort(sortableData)
	case "MERGE_SORT":
		sorted = MergeSort(sortableData)
	case "HEAP_SORT":
		sorted = HeapSort(sortableData)
	case "RADIX_SORT":
		sorted = RadixSort(sortableData)
	case "TIM_SORT":
		sorted = TimSort(sortableData)
	case "INTRO_SORT":
		sorted = IntroSort(sortableData)
	case "SHELL_SORT":
		sorted = ShellSort(sortableData)
	default:
		sorted = QuickSort(sortableData)
	}

	result := make([]types.Transaction, len(sorted))
	for i, st := range sorted {
		result[i] = st.Tx
	}

	var txHashes [][32]byte
	for _, tx := range result {
		txHashes = append(txHashes, tx.ID)
	}
	root := utils.CalculateMerkleRoot(txHashes)

	return result, root
}

func CalculateBlockHash(h types.BlockHeader) [32]byte {
	// Simple serialization for hashing
	record := fmt.Sprintf("%d%x%x%d%d%d%x", h.Version, h.PrevBlockHash, h.MerkleRoot, h.Timestamp, h.Height, h.Nonce, h.VRFSeed)
	return sha256.Sum256([]byte(record))
}

// SelectAlgorithm uses VRF Seed to determine sorting algorithm
func SelectAlgorithm(seed [32]byte) string {
	selector := seed[31] % 7 // Increased to 7 algorithms
	switch selector {
	case 0:
		return "QUICK_SORT"
	case 1:
		return "MERGE_SORT"
	case 2:
		return "HEAP_SORT"
	case 3:
		return "RADIX_SORT"
	case 4:
		return "TIM_SORT"
	case 5:
		return "INTRO_SORT"
	case 6:
		return "SHELL_SORT"
	default:
		return "QUICK_SORT"
	}
}

// MixHash combines ID and seed to create a sorting key
func MixHash(id [32]byte, seed [32]byte) string {
	h := sha256.New()
	h.Write(id[:])
	h.Write(seed[:])
	return string(h.Sum(nil))
}
