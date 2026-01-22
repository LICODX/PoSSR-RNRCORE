package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/config"
	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus"
	"github.com/LICODX/PoSSR-RNRCORE/internal/dashboard"
	"github.com/LICODX/PoSSR-RNRCORE/internal/params"

	// "github.com/LICODX/PoSSR-RNRCORE/internal/dashboard"

	// 5. Start GUI Dashboard (Disabled for Headless Build)
	// dashboard.StartServer("8080", chain, nil)
	"github.com/LICODX/PoSSR-RNRCORE/internal/p2p"

	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/wallet"
)

func main() {
	port := flag.Int("port", 3000, "P2P listening port")
	rpcPort := flag.Int("rpc-port", 9001, "RPC API port")
	dashboardPort := flag.Int("dashboard-port", 9101, "Dashboard web UI port")
	datadir := flag.String("datadir", "./data/chaindata", "Data directory for LevelDB")
	datadirAlias := flag.String("data-dir", "", "Alias for -datadir")
	peers := flag.String("peers", "", "Comma-separated peer addresses")
	peerAlias := flag.String("peer", "", "Alias for -peers (single peer)")
	_ = flag.String("log-file", "", "Log file path (optional, default: stdout)") // TODO: implement log file redirection
	useGossipSub := flag.Bool("gossipsub", true, "Use GossipSub (recommended)")
	isGenesis := flag.Bool("genesis", false, "Start as Genesis Node (Authority)")
	walletPassword := flag.String("wallet-password", "", "Password for wallet encryption (optional)")
	flag.Parse()

	configPath := flag.String("config", "config/mainnet.yaml", "Path to configuration file")
	flag.Parse()

	// Handle aliases
	if *datadirAlias != "" {
		*datadir = *datadirAlias
	}
	if *peerAlias != "" {
		*peers = *peerAlias
	}

	// 0. Load Configuration (if exists)
	// We prioritize Flags > Config File > Defaults
	var cfg *config.Config
	if *configPath != "" {
		loadedCfg, err := config.LoadConfig(*configPath)
		if err == nil {
			fmt.Printf("📄 Loaded config from %s\n", *configPath)
			cfg = loadedCfg
		} else {
			// Only warn if user explicitly provided a non-default path
			if *configPath != "config/mainnet.yaml" {
				fmt.Printf("⚠️ Warning: Failed to load config: %v\n", err)
			}
		}
	}

	// Apply Config Defaults if Flags are unset/default
	if cfg != nil {
		if *port == 3000 && cfg.Network.ListenPort != 0 {
			*port = cfg.Network.ListenPort
		}
		// Assuming dashboard port might be in config too, but for now focus on Network
	}

	fmt.Println("🚀 Starting rnr-core Mainnet Node...")
	fmt.Printf("Config: Port=%d | RPC=%d | Dashboard=%d | DataDir=%s\n", *port, *rpcPort, *dashboardPort, *datadir)
	fmt.Println("Consensus: PoSSR | Block Size: 1 GB | Pruning: ON")

	// 1a. Load or Create Node Wallet
	walletPath := filepath.Join(*datadir, "node_wallet.json")
	var nodeWallet *wallet.Wallet

	if *isGenesis {
		fmt.Println("[GENESIS] Starting as Genesis Authority Node")

		// SECURITY: Get Genesis mnemonic from environment variable (not hardcoded!)
		genesisMnemonic := os.Getenv("GENESIS_MNEMONIC")
		if genesisMnemonic == "" {
			// Fallback: Check for genesis.secret file
			secretPath := filepath.Join(*datadir, "genesis.secret")
			if data, err := os.ReadFile(secretPath); err == nil {
				genesisMnemonic = strings.TrimSpace(string(data))
			}
		}

		if genesisMnemonic == "" {
			fmt.Println("[ERROR] Genesis mnemonic not found!")
			fmt.Println("  Option 1: Set GENESIS_MNEMONIC environment variable")
			fmt.Println("  Option 2: Create data/genesis.secret file with mnemonic")
			fmt.Println("  Generate new mnemonic: go run cmd/genesis-wallet/main.go")
			return
		}

		w, err := wallet.CreateWalletFromMnemonic(genesisMnemonic)
		if err != nil {
			fmt.Printf("Failed to restore Genesis wallet: %v\n", err)
			return
		}
		nodeWallet = w

		// Save Genesis wallet (encrypted if password provided)
		if *walletPassword != "" {
			ks := wallet.NewKeyStore(walletPath)
			if err := ks.Save(nodeWallet, *walletPassword); err != nil {
				fmt.Printf("Failed to save encrypted Genesis wallet: %v\n", err)
				return
			}
			fmt.Println("[GENESIS] Wallet saved (ENCRYPTED)")
		} else {
			data, _ := json.MarshalIndent(nodeWallet, "", "  ")
			os.WriteFile(walletPath, data, 0600)
			fmt.Println("[GENESIS] Wallet saved (CLEARTEXT - use -wallet-password for encryption)")
		}
		fmt.Printf("[GENESIS] Wallet Loaded: %s\n", nodeWallet.Address)
	} else {
		// Try to load existing wallet
		if data, err := os.ReadFile(walletPath); err == nil {
			if err := json.Unmarshal(data, &nodeWallet); err == nil {
				fmt.Printf("[WALLET] Loaded: %s\n", nodeWallet.Address)
			}
		} else {
			// No wallet - auto-create new one
			fmt.Println("  [WALLET] Generating new wallet...")
			newWallet, err := wallet.CreateWallet()
			if err != nil {
				fmt.Printf(" Failed to create wallet: %v\n", err)
				return
			}
			nodeWallet = newWallet
			data, _ := json.MarshalIndent(nodeWallet, "", "  ")
			os.WriteFile(walletPath, data, 0600)
			fmt.Printf(" [WALLET] Created: %s\n", nodeWallet.Address)
		}
	}

	// 1. Initialize Database
	db, err := storage.NewLevelDB(*datadir)
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		return
	}
	defer db.GetDB().Close()

	// 2. Initialize Blockchain State
	// Default to FullNode if no config loaded
	shardCfg := config.ShardConfig{Role: "FullNode", ShardIDs: []int{}}
	if cfg != nil {
		shardCfg = cfg.Sharding
	}
	chain := blockchain.NewBlockchain(db, shardCfg)
	tip := chain.GetTip()
	fmt.Printf("⛓️  Current Tip: Block #%d\n", tip.Height)

	// 3. Setup context
	ctx := context.Background()

	// Create a channel to stop mining if needed (e.g. new block received)
	stopMining := make(chan struct{})

	// 4. Start P2P Network
	var node *p2p.GossipSubNode

	if *useGossipSub {
		var err error
		// Default config if nil
		shardCfg := config.ShardConfig{Role: "FullNode", ShardIDs: []int{}}
		if cfg != nil {
			shardCfg = cfg.Sharding
		}

		node, err = p2p.NewGossipSubNode(ctx, *port, shardCfg)
		if err != nil {
			fmt.Printf("Failed to start GossipSub: %v\n", err)
			return
		}
		defer node.Close()

		if *peers != "" {
			fmt.Printf("Connecting to peers: %s\n", *peers)
			peerList := strings.Split(*peers, ",")
			for _, p := range peerList {
				node.ConnectToPeer(strings.TrimSpace(p))
			}
		} else if cfg != nil && len(cfg.Network.SeedNodes) > 0 {
			fmt.Println("🌱 Connecting to Seed Nodes from config...")
			for _, seed := range cfg.Network.SeedNodes {
				// Skip placeholders
				if strings.Contains(seed, "seed1.rnr.network") {
					continue
				}
				if err := node.ConnectToPeer(seed); err != nil {
					fmt.Printf("   ⚠️ Failed to connect to seed %s: %v\n", seed, err)
				}
			}
		}

		node.DiscoverPeers()

		node.DiscoverPeers()

		// 4a. Listen for HEADERS
		node.ListenForHeaders(func(data []byte) {
			var header types.BlockHeader
			if err := json.Unmarshal(data, &header); err != nil {
				return
			}
			fmt.Printf("📦 Header Received: #%d (Hash: %x)\n", header.Height, header.Hash)

			// TODO (Distributed): We need a 'BlockAssembler' to wait for shards.
			// For now, we just pass a skeleton block to chain?
			// No, chain.AddBlock will reject it.
			// We delay full integration until M3 (Validation Update).
		})

		// 4b. Listen for SHARDS (Configured)
		// We loop through our interested shards
		listeningShards := []int{}
		if cfg != nil && cfg.Sharding.Role == "ShardNode" {
			listeningShards = cfg.Sharding.ShardIDs
		} else {
			for i := 0; i < 10; i++ {
				listeningShards = append(listeningShards, i)
			}
		}

		for _, shardID := range listeningShards {
			// Capture loop variable
			sID := shardID
			node.ListenForShards(sID, func(data []byte) {
				fmt.Printf("🧩 Shard Data Received for Shard %d\n", sID)
				// TODO: Process Shard Data
			})
		}

		node.ListenForTransactions(func(data []byte) {
			fmt.Println("💸 Received transaction from network")
			var tx types.Transaction
			if err := json.Unmarshal(data, &tx); err != nil {
				fmt.Printf("Failed to decode tx: %v\n", err)
				return
			}

			// SECURITY CHECK: Validate against State (Nonce & Balance)
			// This prevents Replay Attacks and insufficient balance spam
			if err := blockchain.ValidateTransactionAgainstState(tx, chain.GetStateManager()); err != nil {
				fmt.Printf("⚠️ Invalid transaction rejected: %v\n", err)
				return
			}

			node.AddToMempool(tx)
		})

		node.ListenForProofs(func(data []byte) {
			fmt.Println("✅ Received proof from network")
		})

	} else {
		fmt.Println("Legacy TCP not supported in this version.")
		return
	}

	// 5a. GUEST MODE: Auto-Register with Genesis if no wallet
	if nodeWallet == nil {
		fmt.Println("[GUEST] Attempting to register with Bootnode...")
		go func() {
			// Retry loop
			for {
				time.Sleep(5 * time.Second)
				resp, err := http.Post("http://127.0.0.1:8080/api/register", "application/json", nil)
				if err != nil {
					fmt.Printf("[FAIL] Connection Failed: %v. Retrying...\n", err)
					continue
				}

				if resp.StatusCode != 200 {
					fmt.Printf("[FAIL] Registration Denied: Status %d. Retrying...\n", resp.StatusCode)
					resp.Body.Close()
					continue
				}

				var newWallet wallet.Wallet
				if err := json.NewDecoder(resp.Body).Decode(&newWallet); err != nil {
					fmt.Printf("[FAIL] Decode Failed: %v\n", err)
					resp.Body.Close()
					continue
				}
				resp.Body.Close()

				// Save
				data, _ := json.MarshalIndent(newWallet, "", "  ")
				os.WriteFile(walletPath, data, 0600)

				fmt.Printf("[SUCCESS] REGISTRATION COMPLETE! Wallet: %s\n", newWallet.Address)
				nodeWallet = &newWallet
				return
			}
		}()

		// Wait for registration before mining
		fmt.Println("[WAIT] Waiting for registration...")
		for nodeWallet == nil {
			time.Sleep(1 * time.Second)
		}
		fmt.Println("[SUCCESS] Identity Confirmed. Starting Mining...")
	}

	// 5. Start GUI Dashboard
	// Pass 'node' which implements MempoolSource
	dashboardPortStr := fmt.Sprintf("%d", *dashboardPort)
	go dashboard.StartServer(dashboardPortStr, chain, node, nodeWallet) // Pass wallet

	// 6. Mining Loop (Proof of Repeated Sorting)
	fmt.Println("🏁 Mining Loop Started. Searching for a valid block...")

	for {
		lastHeader := chain.GetTip()
		difficulty := uint64(1000)

		// Get transactions from P2P mempool
		txs := node.GetMempoolShard()

		//MAINNET MODE: Mine even if mempool is empty (empty blocks are valid)

		// Create Coinbase Transaction (Block Reward)
		var minerAddress [32]byte
		copy(minerAddress[:], nodeWallet.PublicKey)

		coinbaseTx := types.Transaction{
			ID:        [32]byte{1, 1, 1, 1, byte(lastHeader.Height)},
			Sender:    [32]byte{}, // System
			Receiver:  minerAddress,
			Amount:    uint64(params.InitialReward),
			Nonce:     0, // System TX
			Signature: [64]byte{},
		}

		// Prepend Coinbase to transactions
		minableTxs := append([]types.Transaction{coinbaseTx}, txs...)

		var minerPubKey [32]byte
		copy(minerPubKey[:], nodeWallet.PublicKey)

		newBlock, err := consensus.MineBlock(minableTxs, lastHeader, difficulty, stopMining, minerPubKey, nodeWallet.PrivateKey)

		if err != nil {
			if err.Error() == "mining interrupted" {
				fmt.Println("Mining interrupted! Restarting...")
				continue
			}
			fmt.Println("Mining error:", err)
			continue
		}

		fmt.Printf("[SUCCESS] Block Found! Nonce: %d | Hash: %x\n", newBlock.Header.Nonce, newBlock.Header.Hash)

		// Add to local chain
		if err := chain.AddBlock(*newBlock); err != nil {
			fmt.Printf("Failed to add block: %v\n", err)
			continue
		}

		fmt.Printf("[OK] Block Accepted! Height: %d\n", newBlock.Header.Height)

		// THROTTLE: Wait for BlockTime (6s) to ensure consistent heartbeat
		fmt.Printf("[WAIT] Waiting %d seconds for next round...\n", params.BlockTime)
		time.Sleep(time.Duration(params.BlockTime) * time.Second)

		// Broadcast Block (Split into Header + Shards)
		// PublishBlock now takes types.Block struct
		node.PublishBlock(*newBlock)

		// Clear mempool (Simplified)
		node.ClearMempool()
	}
}
