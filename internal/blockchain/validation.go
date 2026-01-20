package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/params"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// ValidateTransaction verifies transaction signature and basic validity
func ValidateTransaction(tx types.Transaction) error {
	// 1. Check signature (Skip for Coinbase/System TX)
	// Coinbase has empty sender [0...0]
	isCoinbase := true
	for _, b := range tx.Sender {
		if b != 0 {
			isCoinbase = false
			break
		}
	}

	if !isCoinbase {
		message := types.SerializeTransaction(tx)
		if !utils.Verify(tx.Sender[:], message, tx.Signature[:]) {
			return fmt.Errorf("invalid signature for tx %x", tx.ID)
		}
	}

	// 2. Basic sanity checks
	if tx.Amount == 0 {
		return fmt.Errorf("zero amount transaction")
	}

	// Allow Self-Transfer for Coinbase
	// Reuse 'isCoinbase' from above (already calculated)

	if !isCoinbase && tx.Sender == tx.Receiver {
		return fmt.Errorf("sender and receiver are the same")
	}

	return nil
}

// ValidateTransactionAgainstState verifies tx against current state (Nonce & Balance)
func ValidateTransactionAgainstState(tx types.Transaction, stateDir *state.Manager) error {
	// 1. Basic validation first
	if err := ValidateTransaction(tx); err != nil {
		return err
	}

	// 2. Get Account State
	acc, err := stateDir.GetAccount(tx.Sender)
	if err != nil {
		return fmt.Errorf("failed to get account state: %v", err)
	}

	// 3. Check Nonce (CRITICAL for Replay Protection)
	// Skip for Coinbase
	isCoinbase := true
	for _, b := range tx.Sender {
		if b != 0 {
			isCoinbase = false
			break
		}
	}

	if !isCoinbase {
		// Expect Nonce = CurrentNonce + 1
		if tx.Nonce != acc.Nonce+1 {
			return fmt.Errorf("invalid nonce: expected %d, got %d", acc.Nonce+1, tx.Nonce)
		}

		// 4. Check Balance (Prevent Mempool Spam)
		if acc.Balance < tx.Amount {
			return fmt.Errorf("insufficient balance: have %d, want %d", acc.Balance, tx.Amount)
		}
	}

	return nil
}

// ValidateBlock performs comprehensive block validation
func ValidateBlock(block types.Block, prevHeader types.BlockHeader) error {
	// 1. Validate timestamp (not too far in future)
	now := time.Now().Unix()
	if block.Header.Timestamp > now+600 {
		return fmt.Errorf("block timestamp too far in future")
	}

	// 2. Validate previous block hash
	if block.Header.Height > 0 {
		// Use PoW hash for comparison (excludes VRFSeed and MerkleRoot)
		expectedPrevHash := types.HashBlockHeaderForPoW(prevHeader)
		if block.Header.PrevBlockHash != expectedPrevHash {
			return fmt.Errorf("invalid previous block hash")
		}
	}

	// 3. Validate block size
	blockSize := calculateBlockSize(block)
	if blockSize > params.MaxBlockSize {
		return fmt.Errorf("block too large: %d bytes (max %d)", blockSize, params.MaxBlockSize)
	}

	// 4. Validate Merkle Root (2-Layer Calculation for Sharding)
	var shardRoots [][32]byte
	for _, shard := range block.Shards {
		// Calculate root for this shard
		var txHashes [][32]byte
		for _, tx := range shard.TxData {
			txHashes = append(txHashes, tx.ID)
		}
		shardRoot := utils.CalculateMerkleRoot(txHashes)
		shardRoots = append(shardRoots, shardRoot)
	}

	// Calculate Global Root from Shard Roots
	calculatedRoot := utils.CalculateMerkleRoot(shardRoots)

	if block.Header.MerkleRoot != calculatedRoot {
		return fmt.Errorf("merkle root mismatch: expected %x, got %x",
			calculatedRoot, block.Header.MerkleRoot)
	}

	// 5. Validate all transactions AND Sorting Order
	for shardID, shard := range block.Shards {
		// A. Validate Transactions
		for _, tx := range shard.TxData {
			if err := ValidateTransaction(tx); err != nil {
				return fmt.Errorf("invalid transaction in block: %v", err)
			}
		}

		// B. VERIFY SORTING ORDER (O(N) - Linear Scan)
		// We do NOT re-sort. We just check if Tx[i] <= Tx[i+1]

		// 1. Determine Algorithm & Seeds
		// VRF Seed was used to select Algo (already in Header)
		// We need to reconstruct the "Sorting Key" for each tx.
		// Note: We need access to consensus package for SelectAlgorithm/MixHash
		// For now, we assume implicit knowledge or move helper functions to shared package.
		// To avoid cycle, we perform a loose check or duplicate MixHash logic here.

		// Re-deriving shard seed:
		// shardSeed := sha256.Sum256(append(block.Header.VRFSeed[:], byte(shardID)))

		// Optimization: Check order linearly
		if len(shard.TxData) > 1 {
			shardSeed := sha256.Sum256(append(block.Header.VRFSeed[:], byte(shardID)))
			prevKey := utils.MixHash(shard.TxData[0].ID, shardSeed)
			for i := 1; i < len(shard.TxData); i++ {
				currKey := utils.MixHash(shard.TxData[i].ID, shardSeed)
				if currKey < prevKey {
					return fmt.Errorf("shard %d is NOT sorted! Cheating detected at index %d", shardID, i)
				}
				prevKey = currKey
			}
		}
	}

	return nil
}

// calculateBlockSize estimates block size in bytes
func calculateBlockSize(block types.Block) uint64 {
	// Rough estimate: header + shards
	size := uint64(1024) // Header overhead
	for _, shard := range block.Shards {
		size += uint64(len(shard.TxData)) * 500 // Avg tx size
	}
	return size
}
