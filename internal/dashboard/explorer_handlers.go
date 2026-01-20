package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RNRScan Explorer API Handlers

// handleBlocks returns paginated list of recent blocks
func (s *Server) handleBlocks(w http.ResponseWriter, r *http.Request) {
	tip := s.bc.GetTip()

	// Get latest 20 blocks (simplified - should support pagination)
	var blocks []map[string]interface{}
	currentHeight := tip.Height

	for i := 0; i < 20 && currentHeight >= 0; i++ {
		// Fetch real block header from database
		blockHeader := s.bc.GetBlockByHeight(currentHeight)
		if blockHeader == nil {
			currentHeight--
			continue
		}

		blockInfo := map[string]interface{}{
			"height":     currentHeight,
			"hash":       fmt.Sprintf("%x", blockHeader.Hash[:16]),
			"timestamp":  blockHeader.Timestamp,
			"txCount":    1,
			"difficulty": blockHeader.Difficulty,
		}
		blocks = append(blocks, blockInfo)
		currentHeight--
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"blocks": blocks,
		"total":  tip.Height,
	})
}

// handleBlockDetail returns detailed info for a specific block
func (s *Server) handleBlockDetail(w http.ResponseWriter, r *http.Request) {
	// Extract block height from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid block path", http.StatusBadRequest)
		return
	}

	// Parse height
	height, err := strconv.ParseUint(pathParts[3], 10, 64)
	if err != nil {
		http.Error(w, "Invalid block height", http.StatusBadRequest)
		return
	}

	// Get block from blockchain
	blockHeader := s.bc.GetBlockByHeight(height)
	if blockHeader == nil {
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}

	blockDetail := map[string]interface{}{
		"height":     blockHeader.Height,
		"hash":       fmt.Sprintf("%x", blockHeader.Hash),
		"prevHash":   fmt.Sprintf("%x", blockHeader.PrevBlockHash),
		"merkleRoot": fmt.Sprintf("%x", blockHeader.MerkleRoot),
		"timestamp":  blockHeader.Timestamp,
		"difficulty": blockHeader.Difficulty,
		"nonce":      blockHeader.Nonce,
		"txCount":    1, // TODO: Get actual TX count
		"vrfSeed":    fmt.Sprintf("%x", blockHeader.VRFSeed[:8]),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockDetail)
}

// handleTxList returns paginated list of recent transactions
func (s *Server) handleTxList(w http.ResponseWriter, r *http.Request) {
	mempool := s.source.GetMempoolShard()

	var txList []map[string]interface{}
	for i, tx := range mempool {
		if i >= 20 { // Limit to 20
			break
		}
		txInfo := map[string]interface{}{
			"hash":   fmt.Sprintf("%x", tx.ID),
			"from":   fmt.Sprintf("%x", tx.Sender[:8]),
			"to":     fmt.Sprintf("%x", tx.Receiver[:8]),
			"amount": tx.Amount,
			"status": "pending",
		}
		txList = append(txList, txInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": txList,
		"total":        len(mempool),
	})
}

// handleTxDetail returns detailed info for a specific transaction
func (s *Server) handleTxDetail(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid tx path", http.StatusBadRequest)
		return
	}

	// Mock TX detail
	txDetail := map[string]interface{}{
		"hash":      pathParts[3],
		"status":    "confirmed",
		"block":     123,
		"from":      "rnr1...",
		"to":        "rnr1...",
		"amount":    100,
		"nonce":     5,
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txDetail)
}

// handleAddressInfo returns balance and TX history for an address
func (s *Server) handleAddressInfo(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid address path", http.StatusBadRequest)
		return
	}

	addressInfo := map[string]interface{}{
		"address":    pathParts[3],
		"balance":    150.00,
		"txCount":    25,
		"lastActive": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addressInfo)
}

// handleSearch performs universal search (blocks/tx/address)
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter required", http.StatusBadRequest)
		return
	}

	// Simple search logic (check if numeric = block, else = address/tx)
	var result map[string]interface{}

	if _, err := strconv.ParseUint(query, 10, 64); err == nil {
		// Numeric - likely block height
		result = map[string]interface{}{
			"type":   "block",
			"result": "/block/" + query,
		}
	} else if strings.HasPrefix(query, "rnr1") {
		// Address
		result = map[string]interface{}{
			"type":   "address",
			"result": "/address/" + query,
		}
	} else {
		// TX hash
		result = map[string]interface{}{
			"type":   "transaction",
			"result": "/tx/" + query,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
