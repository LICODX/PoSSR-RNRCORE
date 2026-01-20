package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

func main() {
	fmt.Println("🔍 INTERNAL SECURITY AUDIT STARTED 🔍")

	// Setup Environment
	dbPath := "./data/internal-audit-test"
	os.RemoveAll(dbPath) // Clean start
	db, err := storage.NewLevelDB(dbPath)
	if err != nil {
		panic(err)
	}
	stateMgr := state.NewManager(db.GetDB())

	TestReplayAttack(stateMgr)
	TestFuzzing()

	// Cleanup
	db.GetDB().Close()
	os.RemoveAll(dbPath)
}

func TestReplayAttack(stateMgr *state.Manager) {
	fmt.Println("\n[TEST 1] Replay Attack (Mempool Flooding)")

	// 1. Setup Wallet
	pubKey, privKey, err := utils.GenerateKeypair()
	if err != nil {
		panic(err)
	}
	var pubKeyBytes [32]byte
	copy(pubKeyBytes[:], pubKey)

	// 2. Fund Account (Mock)
	acc := &state.Account{Balance: 1000, Nonce: 0}
	stateMgr.UpdateAccount(pubKeyBytes, acc)

	// 3. Create Valid Transaction (Nonce 1)
	tx := types.Transaction{
		Sender:   pubKeyBytes,
		Receiver: [32]byte{},
		Amount:   10,
		Nonce:    1,
		Payload:  []byte("Replay Me"),
	}
	// Sign
	msg := types.SerializeTransaction(tx)
	sig := utils.Sign(privKey, msg)
	copy(tx.Signature[:], sig)

	// 4. Simulate P2P Validation (Current Logic)
	// Current Main.go logic:
	// if err := blockchain.ValidateTransaction(tx); err != nil { reject }

	err = blockchain.ValidateTransaction(tx)
	if err != nil {
		fmt.Printf("❌ Unexpected error: Valid tx rejected: %v\n", err)
		return
	}
	fmt.Println("✅ Step 1: Valid Transaction accepted by P2P Layer.")

	// 5. Simulate Mining (Tx is included in a block)
	// This updates the State. Nonce increments to 1.
	err = stateMgr.ApplyTransaction(tx)
	if err != nil {
		fmt.Printf("❌ Failed to apply tx: %v\n", err)
		return
	}
	fmt.Println("✅ Step 2: Transaction mined. Account Nonce is now 1.")

	// 6. ATTACK: Re-broadcast the SAME transaction
	// The attacker sends the same mined tx again to the P2P network.
	// Since P2P *only* calls ValidateTransaction (which checks Signature),
	// and checks internal correctness, but DOES NOT check against Global State...

	fmt.Println("⚡ ATTACK: Re-broadcasting old transaction...")
	err = blockchain.ValidateTransactionAgainstState(tx, stateMgr)

	if err == nil {
		fmt.Println("😱 VULNERABILITY CONFIRMED: P2P Layer accepted a REPLAY transaction!")
		fmt.Println("   Reason: ValidateTransaction() checks signatures but does not check Account Nonce against State.")
		fmt.Println("   Impact: Attacker can flood Mempool with 1,000,000 copies of old mined txs.")
	} else {
		fmt.Println("🛡️ SECURE: Replay transaction rejected efficiently.")
	}
}

func TestFuzzing() {
	fmt.Println("\n[TEST 2] Packet Fuzzing (DoS Protection)")

	// Simulate garbage data sent to JSON unmarshal
	garbage := []byte(`{ "sender": "GARBAGE", "amount": "NAN" }`)

	var tx types.Transaction
	err := json.Unmarshal(garbage, &tx)

	if err != nil {
		fmt.Println("🛡️ SECURE: JSON parser handled garbage gracefully.")
	} else {
		fmt.Println("⚠️ WARNING: JSON parser accepted garbage?")
	}

	// Simulate zero bytes
	zeroes := make([]byte, 100)
	err = json.Unmarshal(zeroes, &tx)
	if err != nil {
		fmt.Println("🛡️ SECURE: JSON parser handled zeroes gracefully.")
	}
}
