package blockchain_test

import (
	"testing"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
)

func TestGenesisBlock(t *testing.T) {
	// Test mainnet genesis is deterministic
	gen1 := blockchain.CreateMainnetGenesis()
	gen2 := blockchain.CreateMainnetGenesis()

	if gen1.Header.Timestamp != gen2.Header.Timestamp {
		t.Error("Mainnet genesis timestamp should be fixed")
	}

	if gen1.Header.VRFSeed != gen2.Header.VRFSeed {
		t.Error("Mainnet genesis VRF seed should be fixed")
	}

	if gen1.Header.Height != 0 {
		t.Error("Genesis block height should be 0")
	}
}

func TestBlockchainInitialization(t *testing.T) {
	// Create temp database
	db, err := storage.NewLevelDB("./test-data")
	if err != nil {
		t.Fatal(err)
	}
	defer db.GetDB().Close()

	// Initialize blockchain
	chain := blockchain.NewBlockchain(db, config.ShardConfig{Role: "FullNode", ShardIDs: []int{}})
	tip := chain.GetTip()

	if tip.Height != 0 {
		t.Errorf("Expected genesis height 0, got %d", tip.Height)
	}
}

func TestDoubleSpendPrevention(t *testing.T) {
	t.Skip("TODO: Implement transaction replay test")
	// TODO: Add test that tries to submit same transaction twice
	// Should fail with "invalid nonce" error
}

func TestForkResolution(t *testing.T) {
	t.Skip("TODO: Implement fork resolution test")
	// TODO: Create two competing chains, verify longest chain wins
}

func TestStateTransitions(t *testing.T) {
	t.Skip("TODO: Implement state transition tests")
	// TODO: Test balance updates, nonce increments
}
