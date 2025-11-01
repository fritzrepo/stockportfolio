package view

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/rivo/tview"
)

// addTransactionForm baut das Formular für neue Transaktionen und verwendet
// package-sichtbare Variablen wie formAddTransaction, depot, pages und showModal.
func addTransactionForm() *tview.Form {

	transaction := storage.Transaction{}

	defaultDate := time.Now().Format("02.01.2006")
	dateField := tview.NewInputField().
		SetLabel("Date").
		SetText(defaultDate).
		SetFieldWidth(20).
		SetAcceptanceFunc(tview.InputFieldMaxLength(10))

	formAddTransaction.AddFormItem(dateField)

	formAddTransaction.AddDropDown("Transaction type", []string{"buy", "sell"}, 0, func(option string, index int) {
		transaction.TransactionType = option
	})

	formAddTransaction.AddDropDown("Asset class", []string{"stock", "crypto", "forex", "warrent"}, 0, func(option string, index int) {
		transaction.AssetType = option
	})

	formAddTransaction.AddInputField("Asset", "", 20, nil, func(asset string) {
		transaction.Asset = asset
	})

	formAddTransaction.AddInputField("Portfolio", "", 20, nil, func(portfolio string) {
		transaction.DepotName = portfolio
	})

	formAddTransaction.AddInputField("Ticker symbol", "", 20, nil, func(tickerSymbol string) {
		transaction.TickerSymbol = tickerSymbol
	})

	formAddTransaction.AddInputField("Quantity", "", 20, nil, func(quantity string) {
		parsedQuantity, err := strconv.ParseFloat(quantity, 64)
		if err != nil {
			// Handle error
			return
		}
		transaction.Quantity = parsedQuantity
	})

	formAddTransaction.AddInputField("Price", "", 20, nil, func(price string) {
		parsedPrice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			// Handle error
			return
		}
		transaction.Price = parsedPrice
	})

	formAddTransaction.AddInputField("Fees", "", 20, nil, func(fees string) {
		parsedFees, err := strconv.ParseFloat(fees, 64)
		if err != nil {
			// Handle error
			return
		}
		transaction.Fees = parsedFees
	})

	formAddTransaction.AddDropDown("Currency", []string{"EUR", "USD", "GBP"}, 0, func(option string, index int) {
		transaction.Currency = option
	})

	formAddTransaction.AddButton("Save", func() {
		// Datum parsen
		dateText := dateField.GetText()
		parsedDate, err := time.Parse("02.01.2006", dateText)
		if err != nil {
			showModal(fmt.Sprintf("Ungültiges Datum:\n%s", err.Error()), []string{"OK"}, nil)
			return
		}

		transaction.Date = parsedDate

		err = depot.AddTransaction(transaction, true)
		if err != nil {
			showModal(fmt.Sprintf("Fehler beim Speichern:\n%s", err.Error()), []string{"OK"}, nil)
			return
		}
		pages.SwitchToPage("Main")
		displayDepotEntries()
	})

	formAddTransaction.AddButton("Cancel", func() {
		pages.SwitchToPage("Main")
	})

	return formAddTransaction
}
