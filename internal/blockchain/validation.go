package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
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
func ValidateBlock(block types.Block, prevHeader types.BlockHeader, shardCfg config.ShardConfig) error {
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

	// 2a. Validate PoW (Difficuly Target)
	// Crucial after criticism about "fake PoW"
	powHash := types.HashBlockHeaderForPoW(block.Header)
	hashInt := new(big.Int).SetBytes(powHash[:])
	maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	targetVal := new(big.Int).Div(maxVal, big.NewInt(int64(block.Header.Difficulty)))
	if hashInt.Cmp(targetVal) != -1 {
		return fmt.Errorf("block hash does not meet difficulty target")
	}

	// 2b. Validate VRF (Miner Signature & Seed Derivation)
	// Seed MUST be H(Signature(Miner, PoWHash))
	if !ed25519.Verify(ed25519.PublicKey(block.Header.MinerPubKey[:]), powHash[:], block.Header.MinerSignature[:]) {
		return fmt.Errorf("invalid miner signature (VRF proof failed)")
	}
	expectedSeed := sha256.Sum256(block.Header.MinerSignature[:])
	if expectedSeed != block.Header.VRFSeed {
		return fmt.Errorf("VRF seed mismatch: does not match signature entropy")
	}

	// 3. Validate block size
	blockSize := calculateBlockSize(block)
	if blockSize > params.MaxBlockSize {
		return fmt.Errorf("block too large: %d bytes (max %d)", blockSize, params.MaxBlockSize)
	}

	// 4. Validate Merkle Root (2-Layer Calculation for Sharding)
	// A. Verify that Header.ShardRoots form Header.MerkleRoot
	// Convert [10][32]byte to [][32]byte for util function
	var shardRootSlice [][32]byte
	for _, root := range block.Header.ShardRoots {
		shardRootSlice = append(shardRootSlice, root)
	}

	recalculatedGlobalRoot := utils.CalculateMerkleRoot(shardRootSlice)
	if recalculatedGlobalRoot != block.Header.MerkleRoot {
		return fmt.Errorf("global merkle root mismatch: expected %x, got %x",
			block.Header.MerkleRoot, recalculatedGlobalRoot)
	}

	// 5. Validate Shards (Partial Validation based on Config)
	// Identify shards we MUST validate
	shardsToValidate := make(map[int]bool)
	if shardCfg.Role == "FullNode" {
		for i := 0; i < 10; i++ {
			shardsToValidate[i] = true
		}
	} else {
		for _, id := range shardCfg.ShardIDs {
			shardsToValidate[id] = true
		}
	}

	for shardID, shard := range block.Shards {
		// Only validate if we are responsible for this shard OR if data is opportunistically present
		// STRICT MODE: If we are responsible, we MUST have data.
		isResponsible := shardsToValidate[shardID]

		if isResponsible {
			if len(shard.TxData) == 0 && shard.ShardRoot != ([32]byte{}) {
				// Warn: We expected data but got none?
				// Empty shard is valid if ShardRoot corresponds to Empty List.
				// But if Root != EmptyRoot, then we are missing data!
				emptyRoot := utils.CalculateMerkleRoot(nil)
				if block.Header.ShardRoots[shardID] != emptyRoot {
					return fmt.Errorf("missing data for assigned shard %d", shardID)
				}
			}

			// A. Validate Transactions & Recalculate Shard Root
			var txHashes [][32]byte
			for _, tx := range shard.TxData {
				if err := ValidateTransaction(tx); err != nil {
					return fmt.Errorf("invalid transaction in shard %d: %v", shardID, err)
				}
				txHashes = append(txHashes, tx.ID)
			}

			calculatedShardRoot := utils.CalculateMerkleRoot(txHashes)
			if calculatedShardRoot != block.Header.ShardRoots[shardID] {
				return fmt.Errorf("shard %d root mismatch: expected %x, got %x",
					shardID, block.Header.ShardRoots[shardID], calculatedShardRoot)
			}

			// B. VERIFY SORTING ORDER (O(N) - Linear Scan)
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
		} else {
			// We are NOT responsible. Trust the Header's ShardRoot.
			// (Verified by Committee/Others)
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
