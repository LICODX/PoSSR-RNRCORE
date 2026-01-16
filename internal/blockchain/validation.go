package blockchain

import (
	"fmt"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/params"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// ValidateTransaction verifies transaction signature and basic validity
func ValidateTransaction(tx types.Transaction) error {
	// 1. Check signature
	message := types.SerializeTransaction(tx)
	if !utils.Verify(tx.Sender[:], message, tx.Signature[:]) {
		return fmt.Errorf("invalid signature for tx %x", tx.ID)
	}

	// 2. Basic sanity checks
	if tx.Amount == 0 {
		return fmt.Errorf("zero amount transaction")
	}

	if tx.Sender == tx.Receiver {
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
	// Expect Nonce = CurrentNonce + 1
	if tx.Nonce != acc.Nonce+1 {
		return fmt.Errorf("invalid nonce: expected %d, got %d", acc.Nonce+1, tx.Nonce)
	}

	// 4. Check Balance (Prevent Mempool Spam)
	if acc.Balance < tx.Amount {
		return fmt.Errorf("insufficient balance: have %d, want %d", acc.Balance, tx.Amount)
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
		expectedPrevHash := types.HashBlockHeader(prevHeader)
		if block.Header.PrevBlockHash != expectedPrevHash {
			return fmt.Errorf("invalid previous block hash")
		}
	}

	// 3. Validate block size
	blockSize := calculateBlockSize(block)
	if blockSize > params.MaxBlockSize {
		return fmt.Errorf("block too large: %d bytes (max %d)", blockSize, params.MaxBlockSize)
	}

	// 4. Validate Merkle Root
	var allTxHashes [][32]byte
	for _, shard := range block.Shards {
		for _, tx := range shard.TxData {
			allTxHashes = append(allTxHashes, tx.ID)
		}
	}

	calculatedRoot := utils.CalculateMerkleRoot(allTxHashes)
	if block.Header.MerkleRoot != calculatedRoot {
		return fmt.Errorf("merkle root mismatch: expected %x, got %x",
			calculatedRoot, block.Header.MerkleRoot)
	}

	// 5. Validate all transactions
	for _, shard := range block.Shards {
		for _, tx := range shard.TxData {
			if err := ValidateTransaction(tx); err != nil {
				return fmt.Errorf("invalid transaction in block: %v", err)
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
