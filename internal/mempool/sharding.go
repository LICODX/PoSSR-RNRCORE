package mempool

import (
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// ShardingManager implements hash-based slot assignment per PoSSR spec
type ShardingManager struct {
	slots [10][]types.Transaction
}

// NewShardingManager creates a new sharding manager
func NewShardingManager() *ShardingManager {
	return &ShardingManager{}
}

// AssignToSlot assigns transaction to appropriate slot based on hash prefix
// Slot 0: 0x0... to 0x1...
// Slot 1: 0x2... to 0x3...
// ...
// Slot 9: 0xE... to 0xF...
func (sm *ShardingManager) AssignToSlot(tx types.Transaction) uint8 {
	// Get first nibble (4 bits) of transaction hash
	firstNibble := tx.ID[0] >> 4

	// Map to slot (0-9)
	// 0x0-0x1 -> Slot 0
	// 0x2-0x3 -> Slot 1
	// ...
	// 0xE-0xF -> Slot 9
	slot := firstNibble / 2
	if slot > 9 {
		slot = 9
	}

	return uint8(slot)
}

// AddTransaction adds transaction to appropriate slot
func (sm *ShardingManager) AddTransaction(tx types.Transaction) {
	slot := sm.AssignToSlot(tx)
	sm.slots[slot] = append(sm.slots[slot], tx)
}

// GetSlot returns all transactions in a specific slot
func (sm *ShardingManager) GetSlot(slotID uint8) []types.Transaction {
	if slotID >= 10 {
		return nil
	}
	return sm.slots[slotID]
}

// GetSlotSize returns number of transactions in slot
func (sm *ShardingManager) GetSlotSize(slotID uint8) int {
	if slotID >= 10 {
		return 0
	}
	return len(sm.slots[slotID])
}

// ClearSlot clears all transactions from a slot
func (sm *ShardingManager) ClearSlot(slotID uint8) {
	if slotID < 10 {
		sm.slots[slotID] = nil
	}
}

// PrintDistribution prints transaction distribution across slots
func (sm *ShardingManager) PrintDistribution() {
	fmt.Println("📊 Mempool Shard Distribution:")
	for i := 0; i < 10; i++ {
		fmt.Printf("  Slot %d: %d transactions\n", i, len(sm.slots[i]))
	}
}
