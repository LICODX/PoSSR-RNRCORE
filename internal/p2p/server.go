package p2p

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

const (
	MaxMessageSize    = 4096  // 4KB max per message
	RateLimit         = 100   // 100 msgs/sec per peer
	MaxMempoolSize    = 10000 // 10K transactions max
	ConnectionTimeout = 30 * time.Second
)

type Server struct {
	ListenAddr string
	Peers      []string
	Mempool    []types.Transaction
	mu         sync.Mutex
	rateLimits map[string]*rateLimiter
}

type rateLimiter struct {
	count   int
	resetAt time.Time
	mu      sync.Mutex
}

func NewServer(addr string) *Server {
	return &Server{
		ListenAddr: addr,
		Mempool:    make([]types.Transaction, 0),
		rateLimits: make(map[string]*rateLimiter),
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		fmt.Printf("Error listening: %v\n", err)
		return
	}
	fmt.Printf("P2P Server listening on %s\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	// Get peer IP for rate limiting
	peerIP := conn.RemoteAddr().String()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(ConnectionTimeout))

	scanner := bufio.NewScanner(conn)
	// Set max message size
	buffer := make([]byte, 0, MaxMessageSize)
	scanner.Buffer(buffer, MaxMessageSize)

	for scanner.Scan() {
		// Check rate limit
		if !s.checkRateLimit(peerIP) {
			fmt.Printf("Rate limit exceeded for %s\n", peerIP)
			return
		}

		msg := scanner.Text()
		// Simple protocol: just print received messages for now
		fmt.Printf("Received from %s: %s\n", peerIP, msg)

		// Reset timeout
		conn.SetReadDeadline(time.Now().Add(ConnectionTimeout))
	}
}

func (s *Server) checkRateLimit(peerIP string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, ok := s.rateLimits[peerIP]
	if !ok {
		limiter = &rateLimiter{
			count:   0,
			resetAt: time.Now().Add(time.Second),
		}
		s.rateLimits[peerIP] = limiter
	}

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// Reset counter if time window passed
	if time.Now().After(limiter.resetAt) {
		limiter.count = 0
		limiter.resetAt = time.Now().Add(time.Second)
	}

	limiter.count++
	return limiter.count <= RateLimit
}

func (s *Server) GetMempoolShard() []types.Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Return a copy
	txs := make([]types.Transaction, len(s.Mempool))
	copy(txs, s.Mempool)
	return txs
}

func (s *Server) BroadcastProof(proof [32]byte) {
	msg := fmt.Sprintf("PROOF:%x", proof)
	s.Broadcast(msg)
}

func (s *Server) Broadcast(msg string) {
	// In a real implementation, iterate over peers and send
	fmt.Printf("[BROADCAST] %s\n", msg)
}

// AddMockTx adds a mock transaction for testing
func (s *Server) AddMockTx() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check mempool limit
	if len(s.Mempool) >= MaxMempoolSize {
		return fmt.Errorf("mempool full (max %d)", MaxMempoolSize)
	}

	tx := types.Transaction{
		Amount: 100,
		Nonce:  uint64(len(s.Mempool)) + 1, // Proper nonce
	}
	// Set a mock ID
	tx.ID = [32]byte{byte(len(s.Mempool))}

	s.Mempool = append(s.Mempool, tx)
	fmt.Println("New transaction added to mempool")
	return nil
}

// ClearMempool removes transactions (e.g., after mining)
func (s *Server) ClearMempool() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Mempool = make([]types.Transaction, 0)
}
