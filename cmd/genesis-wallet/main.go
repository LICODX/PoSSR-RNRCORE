package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/wallet"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🔑 RNR Genesis Wallet Generator")
	fmt.Println("═══════════════════════════════════════════════")

	// Generate Genesis Wallet
	w, err := wallet.CreateWallet()
	if err != nil {
		fmt.Printf("\n❌ Error creating wallet: %v\n", err)
		return
	}

	fmt.Println("\n✅ Genesis Wallet Created Successfully!")

	// Display Mnemonic (MOST IMPORTANT)
	fmt.Println("\n🔐 MNEMONIC PHRASE (12 Words):")
	fmt.Println("───────────────────────────────────────────────")
	fmt.Printf("   %s\n", w.Mnemonic)
	fmt.Println("───────────────────────────────────────────────")

	// Display Address
	fmt.Println("\n📬 ADDRESS (Bech32 format):")
	fmt.Printf("   %s\n", w.Address)

	// Display Derivation Path
	fmt.Println("\n🛤️  DERIVATION PATH:")
	fmt.Printf("   %s\n", w.Path)

	// Save to file (mnemonic + address, NO private key for security)
	data := map[string]string{
		"mnemonic":        w.Mnemonic,
		"address":         w.Address,
		"derivation_path": w.Path,
		"public_key":      hex.EncodeToString(w.PublicKey),
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	err = os.WriteFile("genesis_wallet.json", jsonData, 0600)
	if err != nil {
		fmt.Printf("\n❌ Error saving wallet: %v\n", err)
		return
	}

	// Dump raw keys for script reading
	rawContent := fmt.Sprintf("ADDRESS=%s\nMNEMONIC=%s", w.Address, w.Mnemonic)
	os.WriteFile("genesis_keys.txt", []byte(rawContent), 0644)

	fmt.Println("\n💾 Wallet saved to: genesis_wallet.json")

	// Critical warnings
	fmt.Println("\n⚠️  CRITICAL SECURITY WARNINGS:")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("1. 📝 WRITE DOWN the 12 words above IN ORDER")
	fmt.Println("2. 🔒 Store them in a SAFE, OFFLINE location")
	fmt.Println("3. ❌ NEVER share your mnemonic with anyone")
	fmt.Println("4. 💰 This wallet receives 5 BILLION RNR")
	fmt.Println("5. 🚨 Anyone with these words controls the funds")
	fmt.Println("6. ⚠️  Use this address in mainnet genesis config")
	fmt.Println("═══════════════════════════════════════════════")

	fmt.Println("\n✨ Genesis wallet generation complete!")
}
