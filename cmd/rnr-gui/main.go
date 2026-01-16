package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/gui"
	"github.com/LICODX/PoSSR-RNRCORE/internal/storage"
)

func main() {
	// Create Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("rnr-core - PoSSR Blockchain Node")
	myWindow.Resize(fyne.NewSize(1024, 768))

	// Initialize blockchain backend
	db, err := storage.NewLevelDB("./data/chaindata")
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		return
	}
	defer db.GetDB().Close()

	chain := blockchain.NewBlockchain(db)

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Dashboard", gui.CreateDashboard(chain)),
		container.NewTabItem("Wallet", gui.CreateWallet()),
		container.NewTabItem("Send", gui.CreateSendTab()),
		container.NewTabItem("Peers", gui.CreatePeersTab()),
		container.NewTabItem("Settings", gui.CreateSettingsTab()),
	)

	// Status bar
	statusLabel := widget.NewLabel("🟢 Node Running | Height: 0 | Peers: 0")

	// Update status periodically
	go func() {
		for {
			time.Sleep(1 * time.Second)
			tip := chain.GetTip()
			statusLabel.SetText(fmt.Sprintf("🟢 Node Running | Height: %d | Peers: 0", tip.Height))
		}
	}()

	// Layout
	content := container.NewBorder(
		nil,         // top
		statusLabel, // bottom
		nil,         // left
		nil,         // right
		tabs,        // center
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
