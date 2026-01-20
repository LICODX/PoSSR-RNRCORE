package main

import (
	"fmt"
	"log"

	"github.com/LICODX/PoSSR-RNRCORE/internal/state"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
	"github.com/LICODX/PoSSR-RNRCORE/internal/token"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Example: Create custom tokens on RNR
func main() {
	// 1. Setup (connect to your blockchain DB)
	db, err := storage.NewLevelDB("./data/testnet")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Initialize token system
	registry := token.NewRegistry()
	tokenState := state.NewTokenState(db.DB)
	manager := token.NewManager(registry, tokenState)

	// 3. Create your wallet address (example)
	var myAddress [32]byte
	copy(myAddress[:], []byte("my-test-wallet-address-32bytes"))

	// 4. Create Token #1: Stablecoin
	fmt.Println("Creating USDR (Stablecoin)...")
	usdrMetadata := types.TokenMetadata{
		Name:          "USD on RNR",
		Symbol:        "USDR",
		Decimals:      6,             // 6 decimals like USDC
		InitialSupply: 1000000000000, // 1M USDR (with 6 decimals)
		IsMintable:    true,          // Can mint more
		IsBurnable:    true,          // Can burn
	}

	usdrToken, err := manager.CreateToken(usdrMetadata, myAddress)
	if err != nil {
		log.Fatal("Failed to create USDR:", err)
	}
	fmt.Printf("✅ USDR created! Address: %x\n\n", usdrToken.Address)

	// 5. Create Token #2: GameFi Token
	fmt.Println("Creating GAME token...")
	gameMetadata := types.TokenMetadata{
		Name:          "GameFi Rewards",
		Symbol:        "GAME",
		Decimals:      18,
		InitialSupply: 10000000000000000000000, // 10k GAME
		IsMintable:    true,                    // For player rewards
		IsBurnable:    false,                   // Cannot burn
	}

	gameToken, err := manager.CreateToken(gameMetadata, myAddress)
	if err != nil {
		log.Fatal("Failed to create GAME:", err)
	}
	fmt.Printf("✅ GAME created! Address: %x\n\n", gameToken.Address)

	// 6. Create Token #3: Governance Token
	fmt.Println("Creating DAO governance token...")
	daoMetadata := types.TokenMetadata{
		Name:          "RNR DAO Governance",
		Symbol:        "RNRDAO",
		Decimals:      18,
		InitialSupply: 100000000000000000000000000, // 100M tokens
		IsMintable:    false,                       // Fixed supply
		IsBurnable:    true,                        // Deflationary
	}

	daoToken, err := manager.CreateToken(daoMetadata, myAddress)
	if err != nil {
		log.Fatal("Failed to create RNRDAO:", err)
	}
	fmt.Printf("✅ RNRDAO created! Address: %x\n\n", daoToken.Address)

	// 7. Check balances
	fmt.Println("=== My Token Balances ===")
	usdrBalance := manager.GetBalance(usdrToken.Address, myAddress)
	gameBalance := manager.GetBalance(gameToken.Address, myAddress)
	daoBalance := manager.GetBalance(daoToken.Address, myAddress)

	fmt.Printf("USDR: %d (%.2f with decimals)\n", usdrBalance, float64(usdrBalance)/1000000)
	fmt.Printf("GAME: %d (%.6f with decimals)\n", gameBalance, float64(gameBalance)/1000000000000000000)
	fmt.Printf("RNRDAO: %d (%.6f with decimals)\n\n", daoBalance, float64(daoBalance)/1000000000000000000)

	// 8. Transfer some tokens
	var friendAddress [32]byte
	copy(friendAddress[:], []byte("friend-wallet-address-32bytes!"))

	fmt.Println("Transferring 1000 USDR to friend...")
	err = manager.Transfer(usdrToken.Address, myAddress, friendAddress, 1000000000) // 1000 USDR
	if err != nil {
		log.Fatal("Transfer failed:", err)
	}
	fmt.Println("✅ Transfer successful!\n")

	// 9. Check updated balances
	myNewBalance := manager.GetBalance(usdrToken.Address, myAddress)
	friendBalance := manager.GetBalance(usdrToken.Address, friendAddress)

	fmt.Println("=== Updated USDR Balances ===")
	fmt.Printf("My balance: %.2f USDR\n", float64(myNewBalance)/1000000)
	fmt.Printf("Friend's balance: %.2f USDR\n\n", float64(friendBalance)/1000000)

	// 10. List all tokens
	fmt.Println("=== All Registered Tokens ===")
	allTokens := registry.List()
	for i, t := range allTokens {
		fmt.Printf("%d. %s (%s) - Supply: %d\n", i+1, t.Name, t.Symbol, t.TotalSupply)
	}

	fmt.Println("\n🎉 RNR-20 Token System Working!")
	fmt.Println("You now have a multi-token platform like Ethereum!")
}
