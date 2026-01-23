package blockchain

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

type Blockchain struct {
	store             *storage.Store
	stateManager      *state.Manager
	contractProcessor *ContractProcessor // NEW: for smart contract transactions
	mu                sync.RWMutex
	tip               types.BlockHeader
	shardConfig       config.ShardConfig
}

// NewBlockchain creates a new Blockchain instance
func NewBlockchain(db *storage.Store, shardCfg config.ShardConfig) *Blockchain {
	bc := &Blockchain{
		store:        db,
		stateManager: state.NewManager(db.GetDB()),
		shardConfig:  shardCfg,
	}

	// Initialize Contract Processor
	// TEMP DISABLED: Circular import issue with vm package
	// bc.contractProcessor = NewContractProcessor(
	// 	bc.stateManager.GetContractExecutor(),
	// 	bc.stateManager.GetContractState(),
	// )

	//Try to load tip from DB
	tipData, err := db.GetTip()
	if err == nil {
		// Existing chain
		var header types.BlockHeader
		if err := json.Unmarshal(tipData, &header); err == nil {
			bc.tip = header
			return bc
		}
	}

	// Check if genesis exists
	if !db.HasBlock(0) {
		// Initialize with Mainnet genesis by default
		genesis := CreateGenesisBlock(true)
		if err := db.SaveBlock(genesis); err != nil {
			return nil
		}
		// CRITICAL: Set tip to genesis header after creation!
		bc.tip = genesis.Header
		tipData, _ := json.Marshal(bc.tip)
		db.SaveTip(tipData)
		fmt.Printf("🌍 Genesis Block Created. Hash: %x\n", bc.tip.Hash)
	}

	return bc
}

func (bc *Blockchain) GetTip() types.BlockHeader {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.tip
}

// GetBlockByHeight retrieves a block header by its height
func (bc *Blockchain) GetBlockByHeight(height uint64) *types.BlockHeader {
	header, err := bc.store.GetBlockHeaderByHeight(height)
	if err != nil {
		return nil
	}
	return header
}

// GetStateManager returns the state manager for external access
func (bc *Blockchain) GetStateManager() *state.Manager {
	return bc.stateManager
}

// AddBlock validates and saves a block
func (bc *Blockchain) AddBlock(block types.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 1. Validate Height (critical: check AFTER acquiring lock)
	if block.Header.Height != bc.tip.Height+1 && block.Header.Height != 0 {
		return fmt.Errorf("invalid block height: expected %d, got %d", bc.tip.Height+1, block.Header.Height)
	}

	// 2. Comprehensive validation
	if block.Header.Height > 0 {
		if err := ValidateBlock(block, bc.tip, bc.shardConfig); err != nil {
			return fmt.Errorf("block validation failed: %v", err)
		}
	}

	// 3. Apply all transactions to state
	for _, shard := range block.Shards {
		for _, tx := range shard.TxData {
			// Handle contract transactions
			if tx.Type == types.TxTypeContractDeploy || tx.Type == types.TxTypeContractCall {
				if err := bc.contractProcessor.ProcessContractTransaction(tx); err != nil {
					return fmt.Errorf("failed to process contract tx: %v", err)
				}
			}

			// Apply regular state changes
			if err := bc.stateManager.ApplyTransaction(tx); err != nil {
				return fmt.Errorf("failed to apply tx: %v", err)
			}
		}
	}

	// 4. Save to Disk
	if err := bc.store.SaveBlock(block); err != nil {
		return err
	}

	// 5. Update Tip
	bc.tip = block.Header

	// Save tip to DB
	tipData, _ := json.Marshal(bc.tip)
	bc.store.SaveTip(tipData)

	// 6. Prune Old Blocks (synchronously to avoid race)
	if block.Header.Height > 25 {
		bc.store.PruneOldBlocks(block.Header.Height)
	}

	fmt.Printf("⛓️  Block #%d added to chain.\n", block.Header.Height)
	return nil
}
