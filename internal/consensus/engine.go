package consensus

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/mempool"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// MineBlock runs the Proof of Repeated Sorting (PoRS)
// SECURITY: Algorithm selection happens AFTER block hash is found to prevent prediction attacks
func MineBlock(txs []types.Transaction, prevBlock types.BlockHeader, difficulty uint64, stopChan chan struct{}) (*types.Block, error) {
	// 1. Prepare base data
	var nonce uint64 = 0
	// target := big.NewInt(0).SetUint64(difficulty)
	// Note: In real PoW, target is huge (2^256 / diff). For simulation, we use small numbers.
	// Let's assume higher difficulty = harder to find (target is smaller).
	// Current mock: Target is just a threshold check.
	// REAL POW: Hash < (MaxHash / Difficulty)

	// Create simplified target for PoSSR:
	// We want the Hash to start with N zeroes.

	fmt.Printf("[MINING] Started. Difficulty: %d\n", difficulty)

	for {
		// Check for interrupt
		select {
		case <-stopChan:
			return nil, fmt.Errorf("mining interrupted")
		default:
		}

		// 2. Create candidate header (without algorithm/Merkle yet)
		header := types.BlockHeader{
			Version:       1,
			PrevBlockHash: prevBlock.Hash,
			MerkleRoot:    [32]byte{}, // Will be filled after algorithm selection
			Timestamp:     time.Now().Unix(),
			Height:        prevBlock.Height + 1,
			Nonce:         nonce,
			Difficulty:    difficulty,
			VRFSeed:       [32]byte{}, // Will be filled after hash found
		}

		// 3. Calculate Block Hash (PoW Step)
		// Use PoW-specific hash that excludes VRFSeed and MerkleRoot
		blockHash := types.HashBlockHeaderForPoW(header)

		// 4. Check against Target (PoW Difficulty)
		hashInt := new(big.Int).SetBytes(blockHash[:])
		maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
		targetVal := new(big.Int).Div(maxVal, big.NewInt(int64(difficulty)))

		if hashInt.Cmp(targetVal) == -1 {
			// PoW SUCCESS! Block hash meets difficulty target
			header.Hash = blockHash // Store the PoW hash immediately

			// 5. SECURITY: Derive VRF seed from FOUND hash (unpredictable!)
			// seed = SHA256(BlockHash + Timestamp)
			// This prevents attackers from knowing algorithm before mining
			hasher := sha256.New()
			hasher.Write(blockHash[:])
			binary.Write(hasher, binary.BigEndian, header.Timestamp)
			vrfSeed := hasher.Sum(nil)
			var seed [32]byte
			copy(seed[:], vrfSeed)

			// 6. NOW select algorithm (post-mining, unpredictable)
			algo := utils.SelectAlgorithm(seed)
			fmt.Printf("  [VRF] Post-Mining Algorithm: %s (Seed: %x...)\n", algo, seed[:4])

			// 7. Shard the mempool
			shardingMgr := mempool.NewShardingManager()
			for _, tx := range txs {
				shardingMgr.AddTransaction(tx)
			}

			// 8. Run Sorting Race in PARALLEL (10 Shards)
			var wg sync.WaitGroup
			var shardResults [10]types.ShardData
			var shardRoots [10][32]byte

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(shardID int) {
					defer wg.Done()

					shardTxs := shardingMgr.GetSlot(uint8(shardID))

					// Each shard uses unique seed variation
					shardSeed := sha256.Sum256(append(seed[:], byte(shardID)))

					// Run sorting with DERIVED algorithm
					sorted, root := StartRaceSimplified(shardTxs, shardSeed, algo)

					shardResults[shardID] = types.ShardData{
						NodeID:    [32]byte{byte(shardID)},
						TxData:    sorted,
						ShardRoot: root,
					}
					shardRoots[shardID] = root
				}(i)
			}
			wg.Wait()

			// 9. Calculate Global Merkle Root from 10 Shard Roots
			var allRoots [][32]byte
			for _, r := range shardRoots {
				allRoots = append(allRoots, r)
			}
			globalMerkleRoot := utils.CalculateMerkleRoot(allRoots)

			// 10. Update header with VRF seed and Merkle root
			// NOTE: Do NOT recalculate hash! PoW hash (line 53) is the final hash
			header.VRFSeed = seed
			header.MerkleRoot = globalMerkleRoot
			// Hash was already set at line 53 during PoW

			// 11. Construct Full Block
			block := &types.Block{
				Header: header,
				Shards: shardResults,
			}
			return block, nil
		}

		// 7. Increment Nonce and try again
		nonce++
	}
}

// StartRaceSimplified runs the sorting logic for the mining loop
func StartRaceSimplified(mempool []types.Transaction, seed [32]byte, algo string) ([]types.Transaction, [32]byte) {
	// Algorithm already selected by caller (post-mining)

	sortableData := make([]SortableTransaction, len(mempool))
	for i, tx := range mempool {
		sortableData[i] = SortableTransaction{
			Tx:  tx,
			Key: utils.MixHash(tx.ID, seed),
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
	// CRITICAL: Must use the SAME serialization as validation.go
	return types.HashBlockHeader(h)
}

// Functions moved to pkg/utils/consensus_utils.go to avoid import cycle
// - SelectAlgorithm
// - MixHash
