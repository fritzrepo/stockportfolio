package view

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Tview
var app = tview.NewApplication()
var pages = tview.NewPages()
var formAddTransaction = tview.NewForm()
var mainFlex = tview.NewFlex()
var mainTable = tview.NewTable()

var depot *portfolio.Depot

func GetApp(dep *portfolio.Depot) *tview.Application {

	depot = dep

	displayDepotEntries()

	mainMenu := tview.NewList().
		AddItem("Add transaction", "", 'a', func() {
			formAddTransaction.Clear(true)
			addTransactionForm()
			pages.SwitchToPage("Add Contact")
		}).
		AddItem("Option 2", "", '2', func() {
			// Action für Option 2
		}).
		AddItem("Beenden", "", 'q', func() {
			app.Stop()
		}).ShowSecondaryText(false)

	mainFlex.SetDirection(tview.FlexRow).
		AddItem(mainMenu, 4, 0, true).
		AddItem(mainTable, 0, 1, false)

	pages.AddPage("Main", mainFlex, true, true)
	pages.AddPage("Add Contact", formAddTransaction, true, false)

	return app.SetRoot(pages, true).EnableMouse(true)
}

func displayDepotEntries() {
	var entries = depot.GetEntries()

	mainTable.Clear()

	//create table headers
	mainTable.SetCell(0, 0, tview.NewTableCell("Portfolio").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 1, tview.NewTableCell("Ticker symbol").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 2, tview.NewTableCell("Name").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 3, tview.NewTableCell("Asset class").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 4, tview.NewTableCell("Quantity").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 5, tview.NewTableCell("Price").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 6, tview.NewTableCell("Total Price").SetAlign(tview.AlignCenter))
	mainTable.SetCell(0, 7, tview.NewTableCell("Currency").SetAlign(tview.AlignCenter))

	formatRow(mainTable, 0, tcell.ColorYellow, tcell.ColorBlack, tcell.AttrBold)

	//populate table with entries
	i := 1
	totalVolume := 0.0
	for key, value := range entries {
		mainTable.SetCell(i, 0, tview.NewTableCell(value.DepotName))
		mainTable.SetCell(i, 1, tview.NewTableCell(key))
		mainTable.SetCell(i, 2, tview.NewTableCell(value.Asset))
		mainTable.SetCell(i, 3, tview.NewTableCell(value.AssetType))
		mainTable.SetCell(i, 4, tview.NewTableCell(fmt.Sprintf("%.2f", value.Quantity)))
		mainTable.SetCell(i, 5, tview.NewTableCell(fmt.Sprintf("%.2f", value.Price)))
		mainTable.SetCell(i, 6, tview.NewTableCell(fmt.Sprintf("%.2f", value.TotalPrice())))
		mainTable.SetCell(i, 7, tview.NewTableCell(value.Currency))
		totalVolume += value.TotalPrice()
		i++
	}

	//add total row
	mainTable.SetCell(i, 0, tview.NewTableCell("Total").SetAlign(tview.AlignRight).SetAttributes(tcell.AttrBold))
	mainTable.SetCell(i, 1, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	mainTable.SetCell(i, 2, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	mainTable.SetCell(i, 3, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	mainTable.SetCell(i, 4, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	mainTable.SetCell(i, 5, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	mainTable.SetCell(i, 6, tview.NewTableCell(fmt.Sprintf("%.2f", totalVolume)).SetAttributes(tcell.AttrBold))
	mainTable.SetCell(i, 7, tview.NewTableCell("").SetAlign(tview.AlignCenter))

	mainTable.SetBorder(true).SetTitle("Portfolio Entries")
	mainTable.SetBorders(true)

}

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
			showModal(fmt.Sprintf("Ungültiges Datum:\n%s", err.Error()), []string{"OK"}, func(index int, label string) {
				pages.RemovePage("modal")
			})
			return
		}

		transaction.Date = parsedDate

		err = depot.AddTransaction(transaction, true)
		if err != nil {
			// Zeige Modal mit OK-Button und entferne es beim Schließen
			showModal(fmt.Sprintf("Fehler beim Speichern:\n%s", err.Error()), []string{"OK"}, func(index int, label string) {
				// nur Modal schließen. Alternative Aktionen können hier hinzugefügt werden.
				// Dafür kann man die Parameter index und label verwenden.
				pages.RemovePage("modal")
			})
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

// formatRow setzt Farben/Attribute für alle vorhandenen Zellen in einer Zeile.
func formatRow(table *tview.Table, row int, fg, bg tcell.Color, attrs tcell.AttrMask) {
	count := table.GetColumnCount()
	for col := 0; col < count; col++ {
		cell := table.GetCell(row, col)
		if cell == nil {
			// keine weitere Zelle in dieser Zeile
			break
		}
		cell.SetTextColor(fg).
			SetBackgroundColor(bg).
			SetAttributes(attrs)
	}
}

// Zeigt ein Modal an. Buttons und Done-Funktion können übergeben werden.
func showModal(message string, buttons []string, done func(buttonIndex int, buttonLabel string)) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons(buttons).
		SetDoneFunc(func(index int, label string) {
			// Standard: Modal entfernen
			pages.RemovePage("modal")
			// optionales Callback
			if done != nil {
				done(index, label)
			}
		})

	// Modal zur Seiten-Struktur hinzufügen und anzeigen
	pages.AddPage("modal", modal, true, true)
}
