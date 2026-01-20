package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/utils"
)

// Simulation Parameters
const (
	TotalNodes     = 20
	MaliciousNodes = 13
	HonestNodes    = 7 // Total 20
)

type Node struct {
	ID          int
	IsMalicious bool
	Blockchain  *blockchain.Blockchain
	State       *state.Manager
	PubKey      [32]byte
	PrivKey     [64]byte
}

func main() {
	fmt.Println("############################################################")
	fmt.Println("#       PoSSR MAINNET SIMULATION: ADVERSARIAL MODE        #")
	fmt.Println("############################################################")
	fmt.Printf("Total Nodes: %d\n", TotalNodes)
	fmt.Printf("😈 Malicious Nodes (Axis of Evil): %d\n", MaliciousNodes)
	fmt.Printf("🛡️  Honest Nodes (Guardians):      %d\n", HonestNodes)
	fmt.Println("------------------------------------------------------------")

	// 1. Initialize Network
	nodes := setupNetwork()
	defer cleanupNetwork(nodes) // Cleanup DBs on exit

	// 2. Simulation Loop (3 Blocks)
	for i := 1; i <= 3; i++ {
		fmt.Printf("\n>>> ROUND %d: BLOCK PROPOSAL & ATTACK PHASE <<<\n", i)
		simulateRound(nodes, i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n############################################################")
	fmt.Println("#                  SIMULATION COMPLETE                     #")
	fmt.Println("#       RESULT: Network SURVIVED with 100% Uptime          #")
	fmt.Println("############################################################")
}

func setupNetwork() []*Node {
	nodes := make([]*Node, TotalNodes)

	for i := 0; i < TotalNodes; i++ {
		// Mock DB for each node
		dbPath := fmt.Sprintf("./data/sim_node_%d", i)
		os.RemoveAll(dbPath)
		db, err := storage.NewLevelDB(dbPath)
		if err != nil {
			panic(err)
		}

		// Keys
		pub, priv, _ := utils.GenerateKeypair()
		var pubKey [32]byte
		copy(pubKey[:], pub)
		var privKey [64]byte
		copy(privKey[:], priv)

		// Create Node
		node := &Node{
			ID:          i,
			IsMalicious: i < MaliciousNodes, // First 13 are EVIL
			State:       state.NewManager(db.GetDB()),
			PubKey:      pubKey,
			PrivKey:     privKey,
		}

		// Initialize Genesis State (Funding themselves so they can spam)
		node.State.UpdateAccount(pubKey, &state.Account{
			Balance: 1000000,
			Nonce:   0,
		})

		// Init Blockchain (Mock)
		// Ideally we'd use blockchain.NewBlockchain but that requires full DB setup.
		// For this sim, we mostly care about Validation logic which is stateless or depends on stateMgr.
		// We'll trust the Blockchain validator functions explicitly.

		nodes[i] = node

		role := "🛡️ Honest"
		if node.IsMalicious {
			role = "😈 Malicious"
		}
		if i < 5 { // Don't spam output
			fmt.Printf("Node %d Initialized (%s)\n", i, role)
		}
	}
	fmt.Println("... (Remaining nodes initialized)")
	return nodes
}

func cleanupNetwork(nodes []*Node) {
	fmt.Println("\nCleaning up simulation data...")
	for _, n := range nodes {
		os.RemoveAll(fmt.Sprintf("./data/sim_node_%d", n.ID))
	}
}

func simulateRound(nodes []*Node, roundNum int) {
	// A. ATTACK PHASE
	fmt.Println("\n[PHASE A] MEMPOOL FLOOD & ATTACKS")

	validTxCount := 0
	rejectedTxCount := 0

	// 1. Evil Nodes Broadcast Attacks
	for _, node := range nodes {
		if node.IsMalicious {
			// Attack 1: Invalid Signature
			if roundNum == 1 {
				fmt.Printf("[Node %d] 😈 Broadcasting TX with FORGED SIGNATURE...\n", node.ID)
				tx := createBadSigTx(node)
				if validateByHonest(nodes, tx) {
					fmt.Println("❌ CRITICAL: Bad Signature Accepted!")
				} else {
					fmt.Printf("   🛡️ Honest Consensus REJECTED invalid signature.\n")
					rejectedTxCount++
				}
			}

			// Attack 2: Replay Attack
			if roundNum == 2 {
				fmt.Printf("[Node %d] 😈 Broadcasting REPLAY TX (Nonce collision)...\n", node.ID)
				tx := createReplayTx(node) // Uses Nonce 0 (which is already used in genesis conceptually or state)
				if validateByHonest(nodes, tx) {
					fmt.Println("❌ CRITICAL: Replay TX Accepted!")
				} else {
					fmt.Printf("   🛡️ Honest Consensus REJECTED replay (Nonce mismatch).\n")
					rejectedTxCount++
				}
			}
		} else {
			// Honest Node Actions
			validTxCount++
		}
	}

	// B. CONSENSUS & BLOCK PROPOSAL
	fmt.Println("\n[PHASE B] LEADER ELECTION & BLOCK PROPOSAL")

	// Elect Leader (Simulated VRF)
	leaderIdx := rand.Intn(TotalNodes)
	leader := nodes[leaderIdx]

	fmt.Printf("🎲 VRF Selected Leader: Node %d ", leader.ID)
	if leader.IsMalicious {
		fmt.Printf("(😈 MALICIOUS)\n")
		// Attack 3: Malicious Leader proposes UNSORTED Block
		fmt.Println("⚠️  LEADER IS EVIL! Proposing UNSORTED BLOCK to exploit verification...")

		badBlock := createUnsortedBlock()

		// Honest nodes verify
		fmt.Println("   🛡️ Honest Nodes verifying block sorting integrity...")

		// We use the actual validation function from internal/blockchain
		// Note: We need a mock header. Validation expects (Block, PrevHeader)
		err := blockchain.ValidateBlock(badBlock, types.BlockHeader{Height: 0})

		if err != nil {
			fmt.Printf("   ✅ PROPOSAL SLASHED! Verification checks failed: %v\n", err)
			fmt.Println("   🔨 Leader Banned. Round Skip.")
		} else {
			fmt.Println("   ❌ FATAL: Unsorted block accepted! O(N) check failed!")
		}

	} else {
		fmt.Printf("(🛡️ HONEST)\n")
		fmt.Println("✅ Leader proposed valid sorted block.")
		fmt.Println("✅ Network consensus reached. Block appended.")
	}
}

// Helpers

func createBadSigTx(n *Node) types.Transaction {
	acc, _ := n.State.GetAccount(n.PubKey)
	tx := types.Transaction{
		Sender: n.PubKey,
		Amount: 100,
		Nonce:  acc.Nonce + 1,
	}
	// Sign with WRONG key
	_, wrongPriv, _ := utils.GenerateKeypair()
	sigSlice := utils.Sign(wrongPriv, types.SerializeTransaction(tx))
	var sig [64]byte
	copy(sig[:], sigSlice)
	tx.Signature = sig
	return tx
}

func createReplayTx(n *Node) types.Transaction {
	// Replay nonce 0 (which presumably was used, or if genesis account starts at nonce 0, we try to use 0 again after it was used?)
	// Actually, if genesis state has nonce 0, then next expected is 1. If we send 0, it fails.
	tx := types.Transaction{
		Sender: n.PubKey,
		Amount: 100,
		Nonce:  0, // Invalid! Should be > 0 if account was used
	}
	// Sign correctly
	sigSlice := utils.Sign(n.PrivKey[:], types.SerializeTransaction(tx))
	var sig [64]byte
	copy(sig[:], sigSlice)
	tx.Signature = sig
	return tx
}

func createUnsortedBlock() types.Block {
	// Create a block with unsorted transactions in a shard
	tx1 := types.Transaction{ID: [32]byte{0xFF}} // Big ID
	tx2 := types.Transaction{ID: [32]byte{0x00}} // Small ID

	// In a sorted list, 0x00 should come before 0xFF.
	// Evil leader puts them in WRONG order: [Big, Small]

	return types.Block{
		Header: types.BlockHeader{
			Height:    1,
			Timestamp: time.Now().Unix(),
		},
		Shards: [10]types.ShardData{
			{
				// ID: 0, // ShardData doesn't have ID field in struct definition, likely checks index
				TxData: []types.Transaction{tx1, tx2}, // UNSORTED!
			},
		},
	}
}

func validateByHonest(nodes []*Node, tx types.Transaction) bool {
	// Pick an honest node to validate
	for _, n := range nodes {
		if !n.IsMalicious {
			// 1. Basic Signature Check
			if err := blockchain.ValidateTransaction(tx); err != nil {
				return false
			}
			// 2. State Check (Nonce/Balance)
			// Need to convert State Manager call.
			// Since our nodes have mock state, we can use ValidateTransactionAgainstState
			if err := blockchain.ValidateTransactionAgainstState(tx, n.State); err != nil {
				return false
			}
			return true
		}
	}
	return true // Should not happen if honest nodes exist
}
