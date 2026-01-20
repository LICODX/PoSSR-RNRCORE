package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Test configuration
const (
	BasePort = 9101 // Dashboard/API Port
	NumNodes = 25

	// Transaction schedules
	TransferInterval = 3  // Every 3 blocks
	TokenInterval    = 25 // Every 25 blocks
	ContractInterval = 55 // Every 55 blocks

	// Transaction amounts
	TransferAmount    = 2.0 // 2 RNR per transfer
	TransferSenders   = 5   // 5 random nodes send
	TokenCreators     = 3   // 3 nodes create tokens
	ContractDeployers = 7   // 7 nodes deploy contracts
)

// Global list of discovered address strings
var discoveredAddresses []string

func main() {
	fmt.Println("========================================")
	fmt.Println("RNR Automated Transaction System")
	fmt.Println("========================================")
	fmt.Println()

	rand.Seed(time.Now().UnixNano())

	// Wait for nodes to be ready with retries
	fmt.Println("Waiting for nodes to be ready...")
	if !waitForNode() {
		fmt.Println("❌ Failed to connect to nodes after retries. Exiting.")
		return
	}
	fmt.Println("✅ Nodes are ready!")

	// Discover REAL addresses
	discoverRealAddresses()
	if len(discoveredAddresses) == 0 {
		fmt.Println("⚠️  Warning: No addresses discovered. Transfers may fail.")
	} else {
		fmt.Printf("✅ Discovered %d real wallet addresses from running nodes\n", len(discoveredAddresses))
	}

	// Track block number
	var currentBlock int64 = 0
	var lastBlock int64 = -1

	fmt.Println()
	fmt.Println("Starting automated transactions...")
	fmt.Println("- Every 3 blocks: 5 nodes send 2 RNR to RANDOM REAL ADDRESSES")
	fmt.Println("- Every 25 blocks: 3 nodes create new tokens")
	fmt.Println("- Every 55 blocks: 7 nodes deploy smart contracts")
	fmt.Println()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	consecutiveErrors := 0

	for range ticker.C {
		// Get current block height
		block := getCurrentBlockHeight()
		if block == 0 {
			consecutiveErrors++
			if consecutiveErrors > 20 {
				fmt.Println("⚠️  Too many errors getting block height. Continuing to retry...")
				consecutiveErrors = 0
			}
			continue
		}

		consecutiveErrors = 0

		if block <= lastBlock {
			continue // No new block
		}

		currentBlock = block
		lastBlock = block

		fmt.Printf("\n[Block %d] New block detected\n", currentBlock)

		// Check for random transfers (every 3 blocks)
		if currentBlock%TransferInterval == 0 {
			fmt.Printf("[Block %d] Triggering random transfers...\n", currentBlock)
			sendRandomTransfers()
		}

		// Check for token creation (every 25 blocks)
		if currentBlock%TokenInterval == 0 {
			fmt.Printf("[Block %d] Triggering token creation...\n", currentBlock)
			createTokens(currentBlock)
		}

		// Check for contract deployment (every 55 blocks)
		if currentBlock%ContractInterval == 0 {
			fmt.Printf("[Block %d] Triggering contract deployments...\n", currentBlock)
			deployContracts(currentBlock)
		}
	}
}

// waitForNode waits for a node to be ready with retries
func waitForNode() bool {
	for i := 0; i < 30; i++ {
		// Try accessing /api/wallet as a health check
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/wallet", BasePort))
		if err == nil {
			resp.Body.Close()
			return true
		}
		fmt.Printf("  Retry %d/30: Waiting for node API (%v)...\n", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return false
}

// getCurrentBlockHeight gets the current block height
func getCurrentBlockHeight() int64 {
	client := &http.Client{Timeout: 3 * time.Second}
	// Use /api/blockchain endpoint
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/api/blockchain", BasePort))
	if err != nil {
		// Just fail silently for now unless debug needed
		return 0
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0
	}

	if height, ok := data["height"].(float64); ok {
		return int64(height)
	}
	return 0
}

// discoverRealAddresses queries all running nodes for their wallet addresses
func discoverRealAddresses() {
	discoveredAddresses = make([]string, 0)
	fmt.Println("🔍 Scanning for active node wallets...")

	for port := 9101; port <= 9125; port++ {
		client := &http.Client{Timeout: 3 * time.Second} // Increased timeout to 3s
		resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/api/wallet", port))
		if err != nil {
			if port == 9101 {
				fmt.Printf("DEBUG: Port 9101 connection failed: %v\n", err)
			}
			continue
		}

		var data struct {
			Address string `json:"address"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			if len(data.Address) > 5 {
				discoveredAddresses = append(discoveredAddresses, data.Address)
				fmt.Printf("  Found: %s (Port %d)\n", data.Address, port)
			} else if port == 9101 {
				fmt.Printf("DEBUG: Port 9101 JSON decoded but address empty. Data: %+v\n", data)
			}
		} else if port == 9101 {
			fmt.Printf("DEBUG: Port 9101 JSON decode failed: %v\n", err)
		}
		resp.Body.Close()
	}
}

// sendRandomTransfers sends 2 RNR from 5 random nodes to REAL addresses
func sendRandomTransfers() {
	if len(discoveredAddresses) == 0 {
		discoverRealAddresses()
		if len(discoveredAddresses) == 0 {
			fmt.Println("  ⚠️  No addresses available for transfers")
			return
		}
	}

	// Select 5 random sender nodes (from honest nodes 1-18)
	senderNodes := make([]int, 0, TransferSenders)
	used := make(map[int]bool)

	for len(senderNodes) < TransferSenders {
		node := rand.Intn(18) + 1 // Nodes 1-18
		if !used[node] {
			senderNodes = append(senderNodes, node)
			used[node] = true
		}
	}

	// Send from each selected node
	for _, nodeNum := range senderNodes {
		// Select random recipient from REAL addresses
		recipient := discoveredAddresses[rand.Intn(len(discoveredAddresses))]

		// Create transfer transaction for /api/wallet/send
		tx := map[string]interface{}{
			"to":     recipient,
			"amount": TransferAmount,
			"fee":    0.01,
		}

		// Send via node's RPC (Dashboard API)
		port := BasePort + (nodeNum - 1)
		if err := sendTransaction(port, tx); err != nil {
			fmt.Printf("  ⚠️  Node %d transfer failed: %v\n", nodeNum, err)
		} else {
			fmt.Printf("  ✅ Node %d sent %.1f RNR to %s...\n", nodeNum, TransferAmount, recipient)
		}
	}
}

// createTokens creates new tokens from 3 random nodes
func createTokens(blockNum int64) {
	tokenTemplates := []map[string]interface{}{
		{"name": "USD Reward", "symbol": "USDR", "decimals": 6, "supply": 1000000},
		{"name": "EUR Reward", "symbol": "EURR", "decimals": 6, "supply": 1000000},
		{"name": "JPY Reward", "symbol": "JPYR", "decimals": 2, "supply": 100000000},
		{"name": "Game Token", "symbol": "GAME", "decimals": 18, "supply": 5000000},
		{"name": "Data Token", "symbol": "DATA", "decimals": 18, "supply": 10000000},
		{"name": "Point Token", "symbol": "POINT", "decimals": 0, "supply": 1000000},
		{"name": "DAO Token 1", "symbol": "DAO1", "decimals": 18, "supply": 1000000},
		{"name": "DAO Token 2", "symbol": "DAO2", "decimals": 18, "supply": 2000000},
		{"name": "DAO Token 3", "symbol": "DAO3", "decimals": 18, "supply": 3000000},
	}
	_ = tokenTemplates // Silence unused variable error

	creatorNodes := make([]int, 0, TokenCreators)
	used := make(map[int]bool)

	for len(creatorNodes) < TokenCreators {
		node := rand.Intn(18) + 1
		if !used[node] {
			creatorNodes = append(creatorNodes, node)
			used[node] = true
		}
	}

	for i, nodeNum := range creatorNodes {
		_ = i
		_ = nodeNum
		// Skip token creation for now
		port := BasePort + (nodeNum - 1)
		fmt.Printf("  ℹ️  Node %d skipping token creation (focusing on transfers)\n", port)
	}
}

// deployContracts deploys smart contracts from 7 random nodes
func deployContracts(blockNum int64) {
	fmt.Println("  ℹ️  Skipping contract deployment (focusing on transfers)")
}

// sendTransaction sends a transaction via /api/wallet/send
func sendTransaction(port int, tx map[string]interface{}) error {
	jsonData, _ := json.Marshal(tx)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(
		fmt.Sprintf("http://127.0.0.1:%d/api/wallet/send", port),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// createToken creates a new token (legacy)
func createToken(port int, token map[string]interface{}) error {
	return nil
}
