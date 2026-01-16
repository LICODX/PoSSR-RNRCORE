package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/LICODX/PoSSR-RNRCORE/internal/blockchain"
	"github.com/LICODX/PoSSR-RNRCORE/internal/metrics"
)

// CreateDashboard creates the main dashboard view
func CreateDashboard(chain *blockchain.Blockchain) fyne.CanvasObject {
	// Metrics cards
	blockHeight := widget.NewLabel("0")
	tpsCounter := widget.NewLabel("0")
	totalTx := widget.NewLabel("0")
	peerCount := widget.NewLabel("0")

	// Update metrics periodically
	go func() {
		for {
			time.Sleep(1 * time.Second)
			stats := metrics.GetStats()

			tip := chain.GetTip()
			blockHeight.SetText(fmt.Sprintf("%d", tip.Height))
			tpsCounter.SetText(fmt.Sprintf("%.2f", stats["blocks_per_second"]))
			totalTx.SetText(fmt.Sprintf("%d", stats["transactions_total"]))
			peerCount.SetText(fmt.Sprintf("%d", stats["peer_count"]))
		}
	}()

	// Create metric cards
	blockCard := createMetricCard("Block Height", blockHeight)
	tpsCard := createMetricCard("TPS", tpsCounter)
	txCard := createMetricCard("Total Transactions", totalTx)
	peerCard := createMetricCard("Peers", peerCount)

	// Layout
	grid := container.NewGridWithColumns(2,
		blockCard,
		tpsCard,
		txCard,
		peerCard,
	)

	// Recent blocks list
	recentBlocks := widget.NewList(
		func() int { return 10 },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Block #"),
				widget.NewLabel("Hash"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("Block #%d", id))
			obj.(*fyne.Container).Objects[1].(*widget.Label).SetText("0x1234...abcd")
		},
	)

	return container.NewBorder(
		widget.NewLabelWithStyle("Dashboard", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		recentBlocks,
		nil,
		nil,
		grid,
	)
}

func createMetricCard(title string, value *widget.Label) fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	card := container.NewVBox(
		titleLabel,
		value,
	)

	return card
}
