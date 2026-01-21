package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const (
	// Topics
	TopicHeader       = "rnr/header/1.0.0" // Base header (Small)
	TopicShardPrefix  = "rnr/shard/"       // + shardID (e.g. rnr/shard/0/1.0.0)
	TopicTransactions = "rnr/transactions/1.0.0"
	TopicProofs       = "rnr/proofs/1.0.0"
)

// GossipSubNode wraps LibP2P host with GossipSub
type GossipSubNode struct {
	host   host.Host
	pubsub *pubsub.PubSub
	ctx    context.Context

	blockTopic  *pubsub.Topic // Deprecated, replaced by header + shards
	headerTopic *pubsub.Topic
	shardTopics map[int]*pubsub.Topic
	txTopic     *pubsub.Topic
	proofTopic  *pubsub.Topic

	headerSub *pubsub.Subscription
	shardSubs map[int]*pubsub.Subscription
	txSub     *pubsub.Subscription
	proofSub  *pubsub.Subscription

	shardConfig config.ShardConfig

	// Local Mempool
	Mempool []types.Transaction
	mu      sync.Mutex
}

func (n *GossipSubNode) GetShardConfig() config.ShardConfig {
	return n.shardConfig
}

func (n *GossipSubNode) GetHost() host.Host {
	return n.host
}

// NewGossipSubNode creates a new LibP2P node with GossipSub
func NewGossipSubNode(ctx context.Context, port int, shardConfig config.ShardConfig) (*GossipSubNode, error) {
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
		host:        h,
		pubsub:      ps,
		ctx:         ctx,
		shardConfig: shardConfig,
		shardTopics: make(map[int]*pubsub.Topic),
		shardSubs:   make(map[int]*pubsub.Subscription),
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

	// 1. Join HEADER Topic (ALL Nodes)
	n.headerTopic, err = n.pubsub.Join(TopicHeader)
	if err != nil {
		return err
	}
	n.headerSub, err = n.headerTopic.Subscribe()
	if err != nil {
		return err
	}

	// 2. Join SHARD Topics (Selective)
	// If FullNode, join ALL. If ShardNode, join specific.
	shardsToJoin := []int{}
	if n.shardConfig.Role == "ShardNode" {
		shardsToJoin = n.shardConfig.ShardIDs
		fmt.Printf("🔍 ShardNode Mode: Subscribing to Shards %v\n", shardsToJoin)
	} else {
		// Full Node: Join 0-9
		for i := 0; i < 10; i++ {
			shardsToJoin = append(shardsToJoin, i)
		}
		fmt.Printf("💪 FullNode Mode: Subscribing to ALL Shards (0-9)\n")
	}

	for _, id := range shardsToJoin {
		topicName := fmt.Sprintf("%s%d/1.0.0", TopicShardPrefix, id)
		t, err := n.pubsub.Join(topicName)
		if err != nil {
			return err
		}
		sub, err := t.Subscribe()
		if err != nil {
			return err
		}
		n.shardTopics[id] = t
		n.shardSubs[id] = sub
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

// PublishBlock publishes a block to the network (SPLIT into Header + Shards)
func (n *GossipSubNode) PublishBlock(block types.Block) error {
	// 1. Publish Header
	headerData, _ := json.Marshal(block.Header)
	if err := n.headerTopic.Publish(n.ctx, headerData); err != nil {
		return err
	}

	// 2. Publish Shards
	// Note: We need to publish ALL shards if we produced the block.
	// But we only have topic handles for shards we are subscribed to.
	// Miner MUST be a FullNode (subscribe to all) OR we need to join temporarily.
	// Assumption: Miner is FullNode.

	for i, shard := range block.Shards {
		// Only publish if we have reference to the topic
		if topic, ok := n.shardTopics[i]; ok {
			shardData, _ := json.Marshal(shard)
			if err := topic.Publish(n.ctx, shardData); err != nil {
				fmt.Printf("Error publishing shard %d: %v\n", i, err)
			}
		}
	}
	return nil
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
// ListenForHeaders starts listening for block headers
func (n *GossipSubNode) ListenForHeaders(handler func([]byte)) {
	go func() {
		for {
			msg, err := n.headerSub.Next(n.ctx)
			if err != nil {
				continue
			}
			handler(msg.Data)
		}
	}()
}

// ListenForShards starts listening for specific shard data
func (n *GossipSubNode) ListenForShards(shardID int, handler func([]byte)) {
	sub, ok := n.shardSubs[shardID]
	if !ok {
		return
	}
	go func() {
		for {
			msg, err := sub.Next(n.ctx)
			if err != nil {
				continue
			}
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
			fmt.Printf("[P2P] Connected to %d peers\n", len(peers))
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
