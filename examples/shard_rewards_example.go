package main

import (
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/internal/economics"
)

// ExampleShardRewards demonstrates shard-based reward distribution
func ExampleShardRewards() {
	fmt.Println("\n🎓 Educational Example: Shard-Based Reward Distribution\n")

	// Scenario 1: 10 Validators (Equal Distribution)
	fmt.Println("📊 Scenario 1: 10 Validators, 10 Shards")
	validators10 := make([][32]byte, 10)
	for i := 0; i < 10; i++ {
		validators10[i] = [32]byte{byte(i)}
	}

	assignment10 := economics.AssignShards(validators10, 10)
	assignment10.PrintAssignment()

	rewards10 := economics.CalculateShardRewards(100, assignment10)
	fmt.Println("\n💰 Rewards:")
	for i, validator := range validators10 {
		fmt.Printf("  Validator %d: %d RNR (1 shard)\n", i, rewards10[validator])
	}

	// Scenario 2: 4 Validators (Unequal Distribution)
	fmt.Println("\n📊 Scenario 2: 4 Validators, 10 Shards")
	validators4 := make([][32]byte, 4)
	for i := 0; i < 4; i++ {
		validators4[i] = [32]byte{byte(i)}
	}

	assignment4 := economics.AssignShards(validators4, 10)
	assignment4.PrintAssignment()

	rewards4 := economics.CalculateShardRewards(100, assignment4)
	fmt.Println("\n💰 Rewards:")
	for i, validator := range validators4 {
		shards := assignment4.GetShards(validator)
		fmt.Printf("  Validator %d: %d RNR (%d shards)\n",
			i, rewards4[validator], len(shards))
	}

	// Scenario 3: 2 Validators (50-50 Split)
	fmt.Println("\n📊 Scenario 3: 2 Validators, 10 Shards")
	validators2 := make([][32]byte, 2)
	for i := 0; i < 2; i++ {
		validators2[i] = [32]byte{byte(i)}
	}

	assignment2 := economics.AssignShards(validators2, 10)
	assignment2.PrintAssignment()

	rewards2 := economics.CalculateShardRewards(100, assignment2)
	fmt.Println("\n💰 Rewards:")
	for i, validator := range validators2 {
		shards := assignment2.GetShards(validator)
		fmt.Printf("  Validator %d: %d RNR (%d shards)\n",
			i, rewards2[validator], len(shards))
	}

	fmt.Println("\n✅ Key Insight: Rewards are proportional to work done (shards processed)")
	fmt.Println("   - More shards = more reward")
	fmt.Println("   - Fair distribution regardless of validator count")
	fmt.Println("   - Total always sums to 100% of block reward\n")
}
