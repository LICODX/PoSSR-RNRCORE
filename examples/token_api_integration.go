// RNR-20 API Integration Example
// Add this to your dashboard server to enable token endpoints

package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LICODX/PoSSR-RNRCORE/internal/token"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Initialize token system in your main.go or server setup:
/*
func SetupTokenAPI(db *leveldb.DB) {
	// Create token components
	tokenRegistry := token.NewRegistry()
	tokenState := state.NewTokenState(db)
	tokenManager := token.NewManager(tokenRegistry, tokenState)

	// Register API routes
	http.HandleFunc("/api/token/create", handleTokenCreate(tokenManager))
	http.HandleFunc("/api/token/transfer", handleTokenTransfer(tokenManager))
	http.HandleFunc("/api/token/info/", handleTokenInfo(tokenRegistry))
	http.HandleFunc("/api/tokens", handleTokensList(tokenRegistry))
	http.HandleFunc("/api/token/balance/", handleTokenBalance(tokenManager))
}
*/

// Handler: Create Token
func handleTokenCreate(manager *token.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name          string `json:"name"`
			Symbol        string `json:"symbol"`
			Decimals      uint8  `json:"decimals"`
			InitialSupply uint64 `json:"initialSupply"`
			IsMintable    bool   `json:"isMintable"`
			IsBurnable    bool   `json:"isBurnable"`
			CreatorHex    string `json:"creator"` // Hex encoded [32]byte
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decode creator
		creatorBytes, _ := hex.DecodeString(req.CreatorHex)
		var creator [32]byte
		copy(creator[:], creatorBytes)

		// Create token
		metadata := types.TokenMetadata{
			Name:          req.Name,
			Symbol:        req.Symbol,
			Decimals:      req.Decimals,
			InitialSupply: req.InitialSupply,
			IsMintable:    req.IsMintable,
			IsBurnable:    req.IsBurnable,
		}

		token, err := manager.CreateToken(metadata, creator)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":      true,
			"tokenAddress": hex.EncodeToString(token.Address[:]),
			"name":         token.Name,
			"symbol":       token.Symbol,
		})
	}
}

// Handler: Transfer Token
func handleTokenTransfer(manager *token.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			TokenAddress string `json:"tokenAddress"`
			FromHex      string `json:"from"`
			ToHex        string `json:"to"`
			Amount       uint64 `json:"amount"`
		}

		json.NewDecoder(r.Body).Decode(&req)

		// Decode addresses
		tokenBytes, _ := hex.DecodeString(req.TokenAddress)
		fromBytes, _ := hex.DecodeString(req.FromHex)
		toBytes, _ := hex.DecodeString(req.ToHex)

		var tokenAddr, from, to [32]byte
		copy(tokenAddr[:], tokenBytes)
		copy(from[:], fromBytes)
		copy(to[:], toBytes)

		// Transfer
		err := manager.Transfer(tokenAddr, from, to, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Handler: Get Token Info
func handleTokenInfo(registry *token.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract address from URL: /api/token/info/{address}
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 {
			http.Error(w, "Address required", http.StatusBadRequest)
			return
		}

		addrBytes, _ := hex.DecodeString(parts[4])
		var addr [32]byte
		copy(addr[:], addrBytes)

		token, err := registry.Get(addr)
		if err != nil {
			http.Error(w, "Token not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(token)
	}
}

// Handler: List All Tokens
func handleTokensList(registry *token.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokens := registry.List()

		result := make([]map[string]interface{}, 0)
		for _, t := range tokens {
			result = append(result, map[string]interface{}{
				"address":     hex.EncodeToString(t.Address[:]),
				"name":        t.Name,
				"symbol":      t.Symbol,
				"totalSupply": t.TotalSupply,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"tokens": result,
			"total":  len(result),
		})
	}
}

// Handler: Get Balance
func handleTokenBalance(manager *token.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL: /api/token/balance/{tokenAddr}/{account}
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 6 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		tokenBytes, _ := hex.DecodeString(parts[4])
		accountBytes, _ := hex.DecodeString(parts[5])

		var tokenAddr, account [32]byte
		copy(tokenAddr[:], tokenBytes)
		copy(account[:], accountBytes)

		balance := manager.GetBalance(tokenAddr, account)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": balance,
		})
	}
}
