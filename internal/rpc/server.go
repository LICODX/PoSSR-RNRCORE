package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
)

// Server provides JSON-RPC API
type Server struct {
	chain *blockchain.Blockchain
	state *state.Manager
	port  string
}

// RPCRequest represents JSON-RPC 2.0 request
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

// RPCResponse represents JSON-RPC 2.0 response
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewServer creates a new RPC server
func NewServer(chain *blockchain.Blockchain, state *state.Manager, port string) *Server {
	return &Server{
		chain: chain,
		state: state,
		port:  port,
	}
}

// Start starts the RPC server
func (s *Server) Start() {
	http.HandleFunc("/", s.handleRequest)
	fmt.Printf("🌐 RPC Server listening on http://localhost:%s\n", s.port)
	go http.ListenAndServe(":"+s.port, nil)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, -32700, "Parse error", nil)
		return
	}

	// Route to method handler
	var result interface{}
	var err error

	switch req.Method {
	case "eth_blockNumber":
		result, err = s.getBlockNumber()
	case "eth_getBalance":
		result, err = s.getBalance(req.Params)
	case "eth_sendRawTransaction":
		result, err = s.sendRawTransaction(req.Params)
	case "eth_getBlockByNumber":
		result, err = s.getBlockByNumber(req.Params)
	default:
		s.sendError(w, -32601, "Method not found", req.ID)
		return
	}

	if err != nil {
		s.sendError(w, -32000, err.Error(), req.ID)
		return
	}

	s.sendResult(w, result, req.ID)
}

func (s *Server) getBlockNumber() (uint64, error) {
	tip := s.chain.GetTip()
	return tip.Height, nil
}

func (s *Server) getBalance(params []interface{}) (uint64, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("missing address parameter")
	}

	address := params[0].(string)
	// TODO: Convert address string to [32]byte and query state
	_ = address

	return 1000, nil // Mock
}

func (s *Server) sendRawTransaction(params []interface{}) (string, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("missing transaction data")
	}

	// TODO: Parse transaction, validate, add to mempool
	return "0x" + "mock_tx_hash", nil
}

func (s *Server) getBlockByNumber(params []interface{}) (interface{}, error) {
	// TODO: Load block from storage
	return map[string]interface{}{
		"number": "0x1",
		"hash":   "0x...",
	}, nil
}

func (s *Server) sendResult(w http.ResponseWriter, result interface{}, id interface{}) {
	resp := RPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) sendError(w http.ResponseWriter, code int, message string, id interface{}) {
	resp := RPCResponse{
		JSONRPC: "2.0",
		Error:   &RPCError{Code: code, Message: message},
		ID:      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
