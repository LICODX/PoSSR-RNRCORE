package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/wallet"
)

var currentWallet *wallet.Wallet

// CreateWallet creates the wallet management view
func CreateWallet() fyne.CanvasObject {
	addressLabel := widget.NewLabel("No wallet loaded")
	balanceLabel := widget.NewLabel("Balance: 0 RNR")

	// Create wallet button
	createBtn := widget.NewButton("Create New Wallet", func() {
		w, err := wallet.CreateWallet()
		if err != nil {
			dialog.ShowError(err, nil)
			return
		}
		currentWallet = w
		addressLabel.SetText(fmt.Sprintf("Address: %s", w.Address))
		dialog.ShowInformation("Wallet Created",
			fmt.Sprintf("Your new wallet address:\n%s\n\nPrivate Key (SAVE THIS!):\n%s",
				w.Address, w.ExportPrivateKey()), nil)
	})

	// Import wallet button
	importBtn := widget.NewButton("Import Wallet", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter private key...")

		dialog.ShowForm("Import Wallet", "Import", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Private Key", entry),
			}, func(ok bool) {
				if ok && entry.Text != "" {
					w, err := wallet.ImportPrivateKey(entry.Text)
					if err != nil {
						dialog.ShowError(err, nil)
						return
					}
					currentWallet = w
					addressLabel.SetText(fmt.Sprintf("Address: %s", w.Address))
				}
			}, nil)
	})

	// Export private key button
	exportBtn := widget.NewButton("Export Private Key", func() {
		if currentWallet == nil {
			dialog.ShowInformation("No Wallet", "Please create or import a wallet first", nil)
			return
		}
		dialog.ShowInformation("Private Key",
			fmt.Sprintf("⚠️ KEEP THIS SECRET!\n\n%s", currentWallet.ExportPrivateKey()), nil)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Wallet Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		addressLabel,
		balanceLabel,
		widget.NewSeparator(),
		createBtn,
		importBtn,
		exportBtn,
	)
}
