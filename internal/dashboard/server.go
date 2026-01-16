package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// MempoolSource defines interface for accessing mempool data
type MempoolSource interface {
	GetMempoolShard() []types.Transaction
}

type Server struct {
	Chain *blockchain.Blockchain
	P2P   MempoolSource
}

func StartServer(port string, chain *blockchain.Blockchain, p2p MempoolSource) {
	srv := &Server{Chain: chain, P2P: p2p}

	// Serve static files
	fs := http.FileServer(http.Dir("./internal/dashboard/static"))
	http.Handle("/", fs)

	// API Endpoints
	http.HandleFunc("/api/stats", srv.handleStats)

	fmt.Printf("📊 Dashboard available at http://localhost:%s\n", port)
	go http.ListenAndServe(":"+port, nil)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	tip := s.Chain.GetTip()
	mempool := s.P2P.GetMempoolShard()

	// Mock TPS calculation (just for demo)
	tps := len(mempool) * 5 // Random multiplier to look cool

	stats := map[string]interface{}{
		"height":      tip.Height,
		"mempoolSize": len(mempool),
		"tps":         tps,
		"lastHash":    fmt.Sprintf("%x", tip.MerkleRoot), // Using Merkle Root as visual proxy
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
