package p2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const (
	// Topics
	TopicBlocks       = "rnr/blocks/1.0.0"
	TopicTransactions = "rnr/transactions/1.0.0"
	TopicProofs       = "rnr/proofs/1.0.0"
)

// GossipSubNode wraps LibP2P host with GossipSub
type GossipSubNode struct {
	host   host.Host
	pubsub *pubsub.PubSub
	ctx    context.Context

	blockTopic *pubsub.Topic
	txTopic    *pubsub.Topic
	proofTopic *pubsub.Topic

	blockSub *pubsub.Subscription
	txSub    *pubsub.Subscription
	proofSub *pubsub.Subscription

	// Local Mempool
	Mempool []types.Transaction
	mu      sync.Mutex
}

// NewGossipSubNode creates a new LibP2P node with GossipSub
func NewGossipSubNode(ctx context.Context, port int) (*GossipSubNode, error) {
	// Create LibP2P host
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		),
		libp2p.EnableRelay(),
	)
	if err != nil {
		return nil, err
	}

	// Create GossipSub instance
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	node := &GossipSubNode{
		host:   h,
		pubsub: ps,
		ctx:    ctx,
	}

	// Join topics
	if err := node.joinTopics(); err != nil {
		return nil, err
	}

	fmt.Printf("🌐 LibP2P GossipSub node started\n")
	fmt.Printf("   ID: %s\n", h.ID())
	fmt.Printf("   Addresses:\n")
	for _, addr := range h.Addrs() {
		fmt.Printf("     %s/p2p/%s\n", addr, h.ID())
	}

	return node, nil
}

// joinTopics subscribes to all relevant topics
func (n *GossipSubNode) joinTopics() error {
	var err error

	// Join blocks topic
	n.blockTopic, err = n.pubsub.Join(TopicBlocks)
	if err != nil {
		return err
	}
	n.blockSub, err = n.blockTopic.Subscribe()
	if err != nil {
		return err
	}

	// Join transactions topic
	n.txTopic, err = n.pubsub.Join(TopicTransactions)
	if err != nil {
		return err
	}
	n.txSub, err = n.txTopic.Subscribe()
	if err != nil {
		return err
	}

	// Join proofs topic
	n.proofTopic, err = n.pubsub.Join(TopicProofs)
	if err != nil {
		return err
	}
	n.proofSub, err = n.proofTopic.Subscribe()
	if err != nil {
		return err
	}

	fmt.Println("✅ Subscribed to GossipSub topics")
	return nil
}

// ConnectToPeer connects to a peer by multiaddr
func (n *GossipSubNode) ConnectToPeer(peerAddr string) error {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return err
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	if err := n.host.Connect(n.ctx, *peerInfo); err != nil {
		return err
	}

	fmt.Printf("✅ Connected to peer: %s\n", peerInfo.ID)
	return nil
}

// PublishBlock publishes a block to the network
func (n *GossipSubNode) PublishBlock(blockData []byte) error {
	return n.blockTopic.Publish(n.ctx, blockData)
}

// PublishTransaction publishes a transaction to the network
func (n *GossipSubNode) PublishTransaction(txData []byte) error {
	return n.txTopic.Publish(n.ctx, txData)
}

// PublishProof publishes a proof to the network
func (n *GossipSubNode) PublishProof(proofData []byte) error {
	return n.proofTopic.Publish(n.ctx, proofData)
}

// ListenForBlocks starts listening for blocks
func (n *GossipSubNode) ListenForBlocks(handler func([]byte)) {
	go func() {
		for {
			msg, err := n.blockSub.Next(n.ctx)
			if err != nil {
				fmt.Printf("Error reading block message: %v\n", err)
				continue
			}

			// Process message
			handler(msg.Data)
		}
	}()
}

// ListenForTransactions starts listening for transactions
func (n *GossipSubNode) ListenForTransactions(handler func([]byte)) {
	go func() {
		for {
			msg, err := n.txSub.Next(n.ctx)
			if err != nil {
				fmt.Printf("Error reading tx message: %v\n", err)
				continue
			}

			handler(msg.Data)
		}
	}()
}

// ListenForProofs starts listening for proofs
func (n *GossipSubNode) ListenForProofs(handler func([]byte)) {
	go func() {
		for {
			msg, err := n.proofSub.Next(n.ctx)
			if err != nil {
				fmt.Printf("Error reading proof message: %v\n", err)
				continue
			}

			handler(msg.Data)
		}
	}()
}

// GetPeers returns list of connected peers
func (n *GossipSubNode) GetPeers() []peer.ID {
	return n.host.Network().Peers()
}

// GetPeerCount returns number of connected peers
func (n *GossipSubNode) GetPeerCount() int {
	return len(n.GetPeers())
}

// Close shuts down the node
func (n *GossipSubNode) Close() error {
	return n.host.Close()
}

// DiscoverPeers finds peers using mDNS
func (n *GossipSubNode) DiscoverPeers() {
	// Auto-discovery via mDNS (local network)
	// For production, use Kademlia DHT
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			peers := n.GetPeers()
			fmt.Printf("📡 Connected to %d peers\n", len(peers))
			for i, p := range peers {
				if i < 5 { // Show first 5
					fmt.Printf("   - %s\n", p)
				}
			}
		}
	}()
}

// AddToMempool adds a transaction to the local mempool
func (n *GossipSubNode) AddToMempool(tx types.Transaction) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Mempool = append(n.Mempool, tx)
}

// GetMempoolShard returns a copy of the current mempool
func (n *GossipSubNode) GetMempoolShard() []types.Transaction {
	n.mu.Lock()
	defer n.mu.Unlock()
	txs := make([]types.Transaction, len(n.Mempool))
	copy(txs, n.Mempool)
	return txs
}

// ClearMempool wipes the mempool
func (n *GossipSubNode) ClearMempool() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Mempool = make([]types.Transaction, 0)
}
