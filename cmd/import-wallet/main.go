package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/wallet"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🔓 RNR Wallet Import Tool")
	fmt.Println("═══════════════════════════════════════════════")

	// Get mnemonic from user
	fmt.Println("\n📝 Enter your 12-word mnemonic phrase:")
	fmt.Println("   (separate words with spaces)")
	fmt.Print("\n> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\n❌ Error reading input: %v\n", err)
		return
	}

	mnemonic := strings.TrimSpace(input)

	// Validate mnemonic
	if !wallet.ValidateMnemonic(mnemonic) {
		fmt.Println("\n❌ Invalid mnemonic phrase!")
		fmt.Println("   Please check that:")
		fmt.Println("   - You have exactly 12 words")
		fmt.Println("   - All words are spelled correctly")
		fmt.Println("   - Words are separated by spaces")
		return
	}

	// Import wallet
	w, err := wallet.CreateWalletFromMnemonic(mnemonic)
	if err != nil {
		fmt.Printf("\n❌ Error importing wallet: %v\n", err)
		return
	}

	fmt.Println("\n✅ Wallet Imported Successfully!")

	// Display wallet info
	fmt.Println("\n📬 ADDRESS:")
	fmt.Printf("   %s\n", w.Address)

	fmt.Println("\n🛤️  DERIVATION PATH:")
	fmt.Printf("   %s\n", w.Path)

	fmt.Println("\n🔑 PUBLIC KEY:")
	fmt.Printf("   %s\n", hex.EncodeToString(w.PublicKey))

	// Ask if user wants to save
	fmt.Print("\n💾 Save wallet information to file? (y/n): ")
	saveInput, _ := reader.ReadString('\n')
	saveInput = strings.TrimSpace(strings.ToLower(saveInput))

	if saveInput == "y" || saveInput == "yes" {
		// Save imported wallet
		data := map[string]string{
			"address":         w.Address,
			"derivation_path": w.Path,
			"public_key":      hex.EncodeToString(w.PublicKey),
		}

		jsonData, _ := json.MarshalIndent(data, "", "  ")
		filename := "imported_wallet.json"
		err = os.WriteFile(filename, jsonData, 0600)
		if err != nil {
			fmt.Printf("\n❌ Error saving wallet: %v\n", err)
			return
		}

		fmt.Printf("\n✅ Wallet information saved to: %s\n", filename)
		fmt.Println("\n⚠️  NOTE: Mnemonic is NOT saved for security reasons")
	}

	fmt.Println("\n✨ Import complete!")
}
