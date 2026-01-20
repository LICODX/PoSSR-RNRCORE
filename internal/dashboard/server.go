package dashboard

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/wallet"
)

// MempoolSource defines interface for accessing mempool data
type MempoolSource interface {
	GetMempoolShard() []types.Transaction
	AddToMempool(tx types.Transaction) // Needed for registration TX
}

type Server struct {
	bc                   *blockchain.Blockchain
	source               MempoolSource
	registeredNodesCount int            // Track number of registered guests
	Wallet               *wallet.Wallet // Genesis Wallet (if active)
	stateManager         *state.Manager // For nonce management
}

func StartServer(port string, chain *blockchain.Blockchain, p2p MempoolSource, w *wallet.Wallet) {
	srv := &Server{
		bc:           chain,
		source:       p2p,
		Wallet:       w,
		stateManager: chain.GetStateManager(),
	}

	// Serve static files
	fs := http.FileServer(http.Dir("./internal/dashboard/static"))
	http.Handle("/", fs)

	// API Endpoints
	http.HandleFunc("/api/stats", srv.handleStats)
	http.HandleFunc("/api/register", srv.handleRegister)
	http.HandleFunc("/api/wallet", srv.handleWallet)
	http.HandleFunc("/api/blockchain", srv.handleBlockchainInfo)
	http.HandleFunc("/api/network", srv.handleNetworkInfo)
	http.HandleFunc("/api/blocks/recent", srv.handleRecentBlocks)      // Wallet info
	http.HandleFunc("/api/wallet/send", srv.handleSendTx) // NEW: Send transaction
	http.HandleFunc("/api/mining", srv.handleMining)      // Mining status

	// RNRScan Explorer APIs
	http.HandleFunc("/api/blocks", srv.handleBlocks)        // List blocks
	http.HandleFunc("/api/block/", srv.handleBlockDetail)   // Block detail
	http.HandleFunc("/api/transactions", srv.handleTxList)  // List TXs
	http.HandleFunc("/api/tx/", srv.handleTxDetail)         // TX detail
	http.HandleFunc("/api/address/", srv.handleAddressInfo) // Address info
	http.HandleFunc("/api/search", srv.handleSearch)        // Universal search

	fmt.Printf("[DASH] Dashboard available at http://localhost:%s\n", port)
	go http.ListenAndServe(":"+port, nil)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	tip := s.bc.GetTip()
	mempool := s.source.GetMempoolShard()

	// Mock TPS calculation (just for demo)
	tps := len(mempool) * 5 // Random multiplier to look cool

	stats := map[string]interface{}{
		"height":      tip.Height,
		"mempoolSize": len(mempool),
		"tps":         tps,
		"lastHash":    fmt.Sprintf("%x", tip.MerkleRoot), // Using Merkle Root as visual proxy
		"isGenesis":   (s.Wallet != nil),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleRegister creates a new wallet and funds it (Only works on Genesis Node)
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.Wallet == nil {
		http.Error(w, "Registration only allowed on Genesis Node", http.StatusForbidden)
		return
	}

	// Limit: Maximum 10 guest nodes allowed
	if s.registeredNodesCount >= 10 {
		http.Error(w, "Registration limit reached: Maximum 10 guest nodes", http.StatusTooManyRequests)
		return
	}

	// 1. Generate New Wallet for the Guest
	newWallet, err := wallet.CreateWallet()
	if err != nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	// 2. Create Funding Transaction (1 RNR - minimal balance for transactions)
	// Genesis -> New Wallet's PublicKey (hex-encoded)
	receiverHex := hex.EncodeToString(newWallet.PublicKey)
	tx, err := s.Wallet.CreateTransaction(receiverHex, 1, uint64(time.Now().UnixNano()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create funding TX: %v", err), http.StatusInternalServerError)
		return
	}

	// 3. Add to Mempool
	s.source.AddToMempool(*tx)

	// 4. Increment counter
	s.registeredNodesCount++

	// 5. Return new wallet to guest
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newWallet)

	fmt.Printf("[REG] New Node Registered! Sent 1 RNR to %s (Total: %d/10)\n", newWallet.Address, s.registeredNodesCount)
}

// handleWallet returns wallet information (balance, address, nonce)
func (s *Server) handleWallet(w http.ResponseWriter, r *http.Request) {
	if s.Wallet == nil {
		http.Error(w, "No wallet available", http.StatusNotFound)
		return
	}

	// Get account from state manager
	stateManager := s.bc.GetStateManager()
	var pubkey [32]byte
	copy(pubkey[:], s.Wallet.PublicKey)

	account, err := stateManager.GetAccount(pubkey)
	balance := uint64(0)
	nonce := uint64(0)
	if err == nil && account != nil {
		balance = account.Balance
		nonce = account.Nonce
	}

	walletInfo := map[string]interface{}{
		"address":   s.Wallet.Address,
		"balance":   balance,
		"nonce":     nonce,
		"publicKey": hex.EncodeToString(s.Wallet.PublicKey),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(walletInfo)
}

// handleMining returns mining status and statistics
func (s *Server) handleMining(w http.ResponseWriter, r *http.Request) {
	tip := s.bc.GetTip()

	// Count blocks mined by this wallet (simplified - just show current height)
	miningInfo := map[string]interface{}{
		"status":        "active",
		"currentBlock":  tip.Height,
		"difficulty":    tip.Difficulty,
		"lastBlockTime": time.Unix(tip.Timestamp, 0).Format("15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miningInfo)
}

// handleSendTx creates and broadcasts a transaction
func (s *Server) handleSendTx(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
		Fee    float64 `json:"fee"` // Transaction fee
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate wallet exists
	if s.Wallet == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Wallet not initialized",
		})
		return
	}

	// Validate recipient address (basic check)
	if len(req.To) < 10 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid recipient address",
		})
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Amount must be greater than 0",
		})
		return
	}

	// Decode recipient address from Bech32
	recipientBytes, err := hex.DecodeString(req.To)
	if err != nil || len(recipientBytes) != 32 {
		// Try direct conversion if hex fails
		var receiver [32]byte
		copy(receiver[:], []byte(req.To))
		recipientBytes = receiver[:]
	}

	var receiver [32]byte
	copy(receiver[:], recipientBytes)

	// Get sender's public key
	var sender [32]byte
	copy(sender[:], s.Wallet.PublicKey)

	// Get current nonce from state (SECURE: Sequential, no collision)
	account, err := s.stateManager.GetAccount(sender)
	if err != nil {
		// Account doesn't exist yet, start at nonce 0
		account = &state.Account{Nonce: 0}
	}
	nonce := account.Nonce + 1 // Next sequential nonce

	// Create transaction
	tx := types.Transaction{
		ID:       [32]byte{}, // Will be computed
		Sender:   sender,
		Receiver: receiver,
		Amount:   uint64(req.Amount),
		Fee:      uint64(req.Fee), // Include fee
		Nonce:    nonce,
		Payload:  []byte{},
	}

	// Compute TX ID (hash of transaction data)
	txID := types.HashTransaction(tx)
	tx.ID = txID

	// Sign transaction using wallet's SignTransaction method
	if err := s.Wallet.SignTransaction(&tx); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to sign transaction",
		})
		return
	}

	// Submit to mempool
	s.source.AddToMempool(tx)

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"txHash":  fmt.Sprintf("%x", tx.ID),
		"message": "Transaction submitted to mempool",
	})
}

