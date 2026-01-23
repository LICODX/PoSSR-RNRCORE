package main

import (
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus/bft"
	"github.com/LICODX/PoSSR-RNRCORE/internal/economics"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// ValidatorRewardManager handles validator set management and proportional reward distribution
type ValidatorRewardManager struct {
	// Current epoch's shard assignment
	currentAssignment *economics.ShardAssignment

	// Total shards in the network
	totalShards int
}

// NewValidatorRewardManager creates a new manager
func NewValidatorRewardManager(totalShards int) *ValidatorRewardManager {
	return &ValidatorRewardManager{
		totalShards: totalShards,
	}
}

// UpdateShardAssignment reassigns shards based on current validator set
// Called at each epoch or when validator set changes
func (vrm *ValidatorRewardManager) UpdateShardAssignment(validators *bft.ValidatorSet) {
	// Extract validator addresses
	validatorAddrs := make([][32]byte, 0, len(validators.Validators))
	for _, val := range validators.Validators {
		validatorAddrs = append(validatorAddrs, val.Address)
	}

	// Assign shards using round-robin
	vrm.currentAssignment = economics.AssignShards(validatorAddrs, vrm.totalShards)

	// Print assignment for debugging
	vrm.currentAssignment.PrintAssignment()
}

// DistributeRewards calculates and distributes block rewards proportionally
// Returns map of validator -> reward amount
func (vrm *ValidatorRewardManager) DistributeRewards(
	baseReward uint64,
	validators *bft.ValidatorSet,
) map[[32]byte]uint64 {

	// If no assignment yet, create one
	if vrm.currentAssignment == nil {
		vrm.UpdateShardAssignment(validators)
	}

	// Calculate proportional rewards based on shard processing
	rewards := economics.CalculateShardRewards(baseReward, vrm.currentAssignment)

	// Log reward distribution
	fmt.Println("\n💰 Block Reward Distribution:")
	for validator, amount := range rewards {
		shardCount := vrm.currentAssignment.GetShardCount(validator)
		percentage := float64(shardCount) / float64(vrm.totalShards) * 100
		fmt.Printf("  Validator %x: %d RNR (%.1f%% - %d shards)\n",
			validator[:4], amount, percentage, shardCount)
	}

	return rewards
}

// CreateCoinbaseTransactions creates multiple coinbase transactions for proportional rewards
// One transaction per validator based on shard processing
func (vrm *ValidatorRewardManager) CreateCoinbaseTransactions(
	height uint64,
	baseReward uint64,
	validators *bft.ValidatorSet,
) []types.Transaction {

	rewards := vrm.DistributeRewards(baseReward, validators)

	txs := make([]types.Transaction, 0, len(rewards))

	for validator, amount := range rewards {
		if amount == 0 {
			continue
		}

		coinbaseTx := types.Transaction{
			ID:        [32]byte{1, 1, 1, byte(height), byte(len(txs))}, // Unique ID per coinbase
			Sender:    [32]byte{},                                      // System
			Receiver:  validator,
			Amount:    amount,
			Nonce:     0,
			Signature: [64]byte{},
		}

		txs = append(txs, coinbaseTx)
	}

	return txs
}

// GetValidatorShards returns shard IDs assigned to a specific validator
func (vrm *ValidatorRewardManager) GetValidatorShards(validator [32]byte) []int {
	if vrm.currentAssignment == nil {
		return nil
	}
	return vrm.currentAssignment.GetShards(validator)
}
