package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var gainsTable = tview.NewTable()

// GetGainsForm baut das Formular zur Anzeige der Gewinne und verwendet
// package-sichtbare Variablen wie formShowGains, depot, pages und showModal.
func GetGainsForm() *tview.Flex {

	var gainFlex = tview.NewFlex()

	gainMenu := tview.NewList().
		AddItem("Place holder", "", 'p', func() {
			// Placeholder action
		}).
		AddItem("Back", "", 'b', func() {
			pages.SwitchToPage("Main")
		}).ShowSecondaryText(false)

	gainFlex.SetDirection(tview.FlexRow).
		AddItem(gainMenu, 4, 0, true).
		AddItem(gainsTable, 0, 1, false)

	return gainFlex
}

func loadAndDisplayGains() {
	// Placeholder function to load and display gains
	gains, err := depot.GetAllRealizedGains()
	if err != nil {
		showModal(fmt.Sprintf("Fehler beim Laden der realisierten Gewinne:\n%s", err.Error()), []string{"OK"}, nil)
		return
	}

	performance, err := depot.GetPerformance()
	if err != nil {
		showModal(fmt.Sprintf("Fehler beim Laden der Performance:\n%s", err.Error()), []string{"OK"}, nil)
		return
	}

	gainsTable.Clear()

	//Create table headers
	gainsTable.SetCell(0, 0, tview.NewTableCell("Portfolio").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(0, 1, tview.NewTableCell("Asset").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(0, 2, tview.NewTableCell("Quantity").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(0, 3, tview.NewTableCell("Buy price").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(0, 4, tview.NewTableCell("Sell price").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(0, 5, tview.NewTableCell("Gain Amount").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(0, 6, tview.NewTableCell("Currency").SetAlign(tview.AlignCenter))
	formatRow(gainsTable, 0, tcell.ColorYellow, tcell.ColorBlack, tcell.AttrBold)
	gainsTable.SetFixed(1, 0)

	//Populate table with gains data

	for i, gain := range gains {
		gainsTable.SetCell(i+1, 0, tview.NewTableCell(gain.DepotName))
		gainsTable.SetCell(i+1, 1, tview.NewTableCell(gain.Asset))
		gainsTable.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("%.2f", gain.Quantity)))
		gainsTable.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%.2f", gain.BuyPrice)))
		gainsTable.SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf("%.2f", gain.SellPrice)))
		gainsTable.SetCell(i+1, 5, tview.NewTableCell(fmt.Sprintf("%.2f", gain.Amount)))
		gainsTable.SetCell(i+1, 6, tview.NewTableCell(gain.Currency))
		if gain.IsProfit {
			formatRow(gainsTable, i+1, tcell.ColorGreen, tcell.ColorBlack, tcell.AttrNone)
		} else {
			formatRow(gainsTable, i+1, tcell.ColorRed, tcell.ColorBlack, tcell.AttrNone)
		}
	}

	//add total row
	totalRow := len(gains) + 1
	gainsTable.SetCell(totalRow, 0, tview.NewTableCell("Count of gains").SetAlign(tview.AlignRight).SetAttributes(tcell.AttrBold))
	gainsTable.SetCell(totalRow, 1, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 2, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 3, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 4, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 5, tview.NewTableCell(fmt.Sprintf("%d", performance.CountOfRealizedGains)).SetAttributes(tcell.AttrBold))
	gainsTable.SetCell(totalRow, 6, tview.NewTableCell("").SetAlign(tview.AlignCenter))

	totalRow = len(gains) + 2
	gainsTable.SetCell(totalRow, 0, tview.NewTableCell("Total amount").SetAlign(tview.AlignRight).SetAttributes(tcell.AttrBold))
	gainsTable.SetCell(totalRow, 1, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 2, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 3, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 4, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 5, tview.NewTableCell(fmt.Sprintf("%.2f", performance.TotalGains)).SetAttributes(tcell.AttrBold))
	gainsTable.SetCell(totalRow, 6, tview.NewTableCell("").SetAlign(tview.AlignCenter))

	totalRow = len(gains) + 3
	gainsTable.SetCell(totalRow, 0, tview.NewTableCell("Total invested amount").SetAlign(tview.AlignRight).SetAttributes(tcell.AttrBold))
	gainsTable.SetCell(totalRow, 1, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 2, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 3, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 4, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 5, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	gainsTable.SetCell(totalRow, 5, tview.NewTableCell(fmt.Sprintf("%.2f", performance.TotalInvestedAmount)).SetAttributes(tcell.AttrBold))
	gainsTable.SetCell(totalRow, 6, tview.NewTableCell("").SetAlign(tview.AlignCenter))

	gainsTable.SetBorder(true).SetTitle("Realized gains")
	gainsTable.SetBorders(true)

}
