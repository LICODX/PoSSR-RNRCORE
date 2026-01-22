package economics

import (
	"testing"
)

// TestShardAssignment_10Validators tests equal distribution (1 shard each)
func TestShardAssignment_10Validators(t *testing.T) {
	validators := make([][32]byte, 10)
	for i := 0; i < 10; i++ {
		validators[i] = [32]byte{byte(i)}
	}

	assignment := AssignShards(validators, 10)

	// Each validator should have exactly 1 shard
	for _, validator := range validators {
		count := assignment.GetShardCount(validator)
		if count != 1 {
			t.Errorf("Expected 1 shard per validator, got %d", count)
		}
	}

	// Validate assignment is complete
	if err := assignment.ValidateAssignment(); err != nil {
		t.Errorf("Invalid assignment: %v", err)
	}
}

// TestShardAssignment_4Validators tests unequal distribution
func TestShardAssignment_4Validators(t *testing.T) {
	validators := make([][32]byte, 4)
	for i := 0; i < 4; i++ {
		validators[i] = [32]byte{byte(i)}
	}

	assignment := AssignShards(validators, 10)

	// Expected distribution: 3, 3, 2, 2 shards
	expectedCounts := map[int]int{
		3: 2, // 2 validators with 3 shards
		2: 2, // 2 validators with 2 shards
	}

	actualCounts := make(map[int]int)
	for _, validator := range validators {
		count := assignment.GetShardCount(validator)
		actualCounts[count]++
	}

	if actualCounts[3] != expectedCounts[3] || actualCounts[2] != expectedCounts[2] {
		t.Errorf("Expected distribution %v, got %v", expectedCounts, actualCounts)
	}

	// Validate assignment
	if err := assignment.ValidateAssignment(); err != nil {
		t.Errorf("Invalid assignment: %v", err)
	}
}

// TestShardAssignment_2Validators tests 50-50 split
func TestShardAssignment_2Validators(t *testing.T) {
	validators := make([][32]byte, 2)
	for i := 0; i < 2; i++ {
		validators[i] = [32]byte{byte(i)}
	}

	assignment := AssignShards(validators, 10)

	// Each validator should have exactly 5 shards
	for _, validator := range validators {
		count := assignment.GetShardCount(validator)
		if count != 5 {
			t.Errorf("Expected 5 shards per validator, got %d", count)
		}
	}

	// Validate assignment
	if err := assignment.ValidateAssignment(); err != nil {
		t.Errorf("Invalid assignment: %v", err)
	}
}

// TestCalculateShardRewards tests reward calculation
func TestCalculateShardRewards(t *testing.T) {
	validators := make([][32]byte, 4)
	for i := 0; i < 4; i++ {
		validators[i] = [32]byte{byte(i)}
	}

	assignment := AssignShards(validators, 10)
	totalReward := uint64(100)

	rewards := CalculateShardRewards(totalReward, assignment)

	// Sum of all rewards should equal total reward
	sum := uint64(0)
	for _, reward := range rewards {
		sum += reward
	}

	if sum != totalReward {
		t.Errorf("Expected total reward %d, got %d", totalReward, sum)
	}

	// Each validator should get reward proportional to shard count
	rewardPerShard := totalReward / 10
	for validator, shards := range assignment.ValidatorShards {
		expectedReward := rewardPerShard * uint64(len(shards))
		actualReward := rewards[validator]

		if actualReward != expectedReward {
			t.Errorf("Validator %x: expected %d, got %d",
				validator[:2], expectedReward, actualReward)
		}
	}
}

// TestShardAssignment_SingleValidator tests all shards to one validator
func TestShardAssignment_SingleValidator(t *testing.T) {
	validators := [][32]byte{{0x01}}

	assignment := AssignShards(validators, 10)

	// Single validator should get all shards
	count := assignment.GetShardCount(validators[0])
	if count != 10 {
		t.Errorf("Expected 10 shards, got %d", count)
	}

	// Reward should be 100%
	rewards := CalculateShardRewards(100, assignment)
	if rewards[validators[0]] != 100 {
		t.Errorf("Expected 100 RNR, got %d", rewards[validators[0]])
	}
}
