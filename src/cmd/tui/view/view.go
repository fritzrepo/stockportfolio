package view

import (
	"fmt"

	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Package-sichtbare Variablen für die TUI-Komponenten
var app = tview.NewApplication()
var pages = tview.NewPages()
var formAddTransaction = tview.NewForm()
var mainFlex = tview.NewFlex()
var mainTable = tview.NewTable()

var depot *portfolio.Depot

func GetView(dep *portfolio.Depot) *tview.Application {

	depot = dep

	var formShowGains = GetGainsForm()
	var formDebugSession = GetSessionDebugForm()

	displayDepotEntries()

	mainMenu := tview.NewList().
		AddItem("Add transaction", "", 'a', func() {
			formAddTransaction.Clear(true)
			addTransactionForm()
			pages.SwitchToPage("AddTransaction")
		}).
		AddItem("Show Gains", "", 'g', func() {
			loadAndDisplayGains()
			pages.SwitchToPage("ShowGains")
		}).
		AddItem("Session debug", "", 's', func() {
			pages.SwitchToPage("DebugSession")
		}).
		AddItem("Quit", "", 'q', func() {
			app.Stop()
		}).ShowSecondaryText(false)

	mainFlex.SetDirection(tview.FlexRow).
		AddItem(mainMenu, 4, 0, true).
		AddItem(mainTable, 0, 1, false)

	// 97 = 'a', 113 = 'q', 103 = 'g'
	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			formAddTransaction.Clear(true)
			addTransactionForm()
			pages.SwitchToPage("AddTransaction")
		} else if event.Rune() == 103 {
			loadAndDisplayGains()
			pages.SwitchToPage("ShowGains")
		}
		return event
	})

	pages.AddPage("Main", mainFlex, true, true)
	pages.AddPage("AddTransaction", formAddTransaction, true, false)
	pages.AddPage("ShowGains", formShowGains, true, false)
	pages.AddPage("DebugSession", formDebugSession, true, false)

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
	gainsTable.SetFixed(1, 0)

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
