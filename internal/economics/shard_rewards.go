package economics

import (
	"fmt"
)

// ShardAssignment represents which shards are assigned to which validator
type ShardAssignment struct {
	ValidatorShards map[[32]byte][]int // validator address -> assigned shard IDs
	TotalShards     int
}

// NewShardAssignment creates a new shard assignment
func NewShardAssignment(totalShards int) *ShardAssignment {
	return &ShardAssignment{
		ValidatorShards: make(map[[32]byte][]int),
		TotalShards:     totalShards,
	}
}

// AssignShards assigns shards to validators using round-robin load balancing
// This ensures fair distribution regardless of validator count
func AssignShards(validators [][32]byte, totalShards int) *ShardAssignment {
	assignment := NewShardAssignment(totalShards)

	if len(validators) == 0 {
		return assignment
	}

	// Round-robin assignment
	for shardID := 0; shardID < totalShards; shardID++ {
		validatorIndex := shardID % len(validators)
		validator := validators[validatorIndex]

		assignment.ValidatorShards[validator] = append(
			assignment.ValidatorShards[validator],
			shardID,
		)
	}

	return assignment
}

// GetShardCount returns number of shards assigned to a validator
func (sa *ShardAssignment) GetShardCount(validator [32]byte) int {
	return len(sa.ValidatorShards[validator])
}

// GetShards returns the shard IDs assigned to a validator
func (sa *ShardAssignment) GetShards(validator [32]byte) []int {
	return sa.ValidatorShards[validator]
}

// PrintAssignment prints shard assignment for debugging
func (sa *ShardAssignment) PrintAssignment() {
	fmt.Println("\n📦 Shard Assignment:")
	for validator, shards := range sa.ValidatorShards {
		fmt.Printf("  Validator %x: Shards %v (%d shards)\n",
			validator[:4], shards, len(shards))
	}
}

// CalculateShardRewards calculates rewards based on shard processing
// Each shard contributes equally to total reward
func CalculateShardRewards(
	totalReward uint64,
	assignment *ShardAssignment,
) map[[32]byte]uint64 {

	if assignment.TotalShards == 0 {
		return make(map[[32]byte]uint64)
	}

	// Reward per shard (fixed)
	rewardPerShard := totalReward / uint64(assignment.TotalShards)

	// Calculate reward for each validator based on shard count
	rewards := make(map[[32]byte]uint64)

	for validator, shards := range assignment.ValidatorShards {
		numProcessed := len(shards)
		rewards[validator] = rewardPerShard * uint64(numProcessed)
	}

	return rewards
}

// ValidateAssignment checks if shard assignment is complete and valid
func (sa *ShardAssignment) ValidateAssignment() error {
	// Check if all shards are assigned
	assignedShards := make(map[int]bool)

	for _, shards := range sa.ValidatorShards {
		for _, shardID := range shards {
			if shardID < 0 || shardID >= sa.TotalShards {
				return fmt.Errorf("invalid shard ID: %d", shardID)
			}
			if assignedShards[shardID] {
				return fmt.Errorf("shard %d assigned multiple times", shardID)
			}
			assignedShards[shardID] = true
		}
	}

	// Check if all shards are covered
	if len(assignedShards) != sa.TotalShards {
		return fmt.Errorf("incomplete assignment: %d/%d shards assigned",
			len(assignedShards), sa.TotalShards)
	}

	return nil
}

// Example scenarios for testing:
//
// Scenario 1: 10 Validators, 10 Shards
// Each validator: 1 shard
// Reward: 10 RNR per validator (100 RNR / 10)
//
// Scenario 2: 4 Validators, 10 Shards
// Validator 0: Shards [0, 4, 8]     → 3 shards → 30 RNR
// Validator 1: Shards [1, 5, 9]     → 3 shards → 30 RNR
// Validator 2: Shards [2, 6]        → 2 shards → 20 RNR
// Validator 3: Shards [3, 7]        → 2 shards → 20 RNR
//
// Scenario 3: 2 Validators, 10 Shards
// Validator 0: Shards [0, 2, 4, 6, 8] → 5 shards → 50 RNR
// Validator 1: Shards [1, 3, 5, 7, 9] → 5 shards → 50 RNR
