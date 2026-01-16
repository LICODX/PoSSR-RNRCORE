package p2p

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PeerInfo represents a network peer
type PeerInfo struct {
	ID        string
	Address   string
	Height    uint64
	LastSeen  time.Time
	Connected bool
}

// Discovery handles peer discovery
type Discovery struct {
	peers     map[string]*PeerInfo
	seedNodes []string
	mu        sync.RWMutex
}

// NewDiscovery creates a new peer discovery service
func NewDiscovery(seedNodes []string) *Discovery {
	return &Discovery{
		peers:     make(map[string]*PeerInfo),
		seedNodes: seedNodes,
	}
}

// DiscoverPeers initiates peer discovery
func (d *Discovery) DiscoverPeers() {
	fmt.Println("🔍 Starting peer discovery...")

	// Connect to seed nodes
	for _, seed := range d.seedNodes {
		d.connectToPeer(seed)
	}

	// Start periodic peer discovery
	go d.discoveryLoop()
}

func (d *Discovery) connectToPeer(address string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	peerID := fmt.Sprintf("peer-%d", rand.Intn(10000))
	d.peers[peerID] = &PeerInfo{
		ID:        peerID,
		Address:   address,
		Height:    0,
		LastSeen:  time.Now(),
		Connected: true,
	}

	fmt.Printf("  ✓ Connected to peer: %s (%s)\n", peerID, address)
}

func (d *Discovery) discoveryLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		d.refreshPeers()
	}
}

func (d *Discovery) refreshPeers() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Remove stale peers
	for id, peer := range d.peers {
		if time.Since(peer.LastSeen) > 5*time.Minute {
			fmt.Printf("  ✗ Removing stale peer: %s\n", id)
			delete(d.peers, id)
		}
	}
}

// GetPeers returns list of active peers
func (d *Discovery) GetPeers() []*PeerInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var result []*PeerInfo
	for _, peer := range d.peers {
		result = append(result, peer)
	}
	return result
}

// GetPeerCount returns number of connected peers
func (d *Discovery) GetPeerCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.peers)
}
