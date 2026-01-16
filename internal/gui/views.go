package gui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// CreateSendTab creates the send transaction view
func CreateSendTab() fyne.CanvasObject {
	toEntry := widget.NewEntry()
	toEntry.SetPlaceHolder("Recipient address (0x...)")

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Amount in RNR")

	sendBtn := widget.NewButton("Send Transaction", func() {
		if currentWallet == nil {
			dialog.ShowInformation("No Wallet", "Please create or import a wallet first", nil)
			return
		}

		amount, err := strconv.ParseUint(amountEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(err, nil)
			return
		}

		// Create and sign transaction
		tx, err := currentWallet.CreateTransaction(toEntry.Text, amount, 1)
		if err != nil {
			dialog.ShowError(err, nil)
			return
		}

		// TODO: Broadcast to network
		_ = tx
		dialog.ShowInformation("Success", "Transaction sent!", nil)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "To", Widget: toEntry},
			{Text: "Amount", Widget: amountEntry},
		},
		OnSubmit: func() {
			sendBtn.OnTapped()
		},
	}

	return container.NewVBox(
		widget.NewLabelWithStyle("Send Transaction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		form,
		sendBtn,
	)
}

// CreatePeersTab creates the peers view
func CreatePeersTab() fyne.CanvasObject {
	peerList := widget.NewList(
		func() int { return 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Peer")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText("Peer " + strconv.Itoa(id+1) + ": 192.168.1.100:3000")
		},
	)

	return container.NewBorder(
		widget.NewLabelWithStyle("Connected Peers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		peerList,
	)
}

// CreateSettingsTab creates the settings view
func CreateSettingsTab() fyne.CanvasObject {
	dataDirEntry := widget.NewEntry()
	dataDirEntry.SetText("./data/chaindata")

	p2pPortEntry := widget.NewEntry()
	p2pPortEntry.SetText("3000")

	rpcPortEntry := widget.NewEntry()
	rpcPortEntry.SetText("9545")

	saveBtn := widget.NewButton("Save Settings", func() {
		dialog.ShowInformation("Settings", "Settings saved!", nil)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data Directory", Widget: dataDirEntry},
			{Text: "P2P Port", Widget: p2pPortEntry},
			{Text: "RPC Port", Widget: rpcPortEntry},
		},
	}

	return container.NewVBox(
		widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		form,
		saveBtn,
	)
}
