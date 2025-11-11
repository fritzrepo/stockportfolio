package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var sessionTable = tview.NewTable()

// GetSessionForm baut das Formular zur Anzeige der Gewinne und verwendet
// package-sichtbare Variablen wie formShowGains, depot, pages und showModal.
func GetSessionDebugForm() *tview.Flex {

	var sessionFlex = tview.NewFlex()

	sessionMenu := tview.NewList().
		AddItem("Next step", "", 'n', nil).
		AddItem("Back", "", 'b', func() {
			pages.SwitchToPage("Main")
		}).ShowSecondaryText(false)

	sessionFlex.SetDirection(tview.FlexRow).
		AddItem(sessionMenu, 4, 0, true).
		AddItem(sessionTable, 0, 1, false)

	// 98 = 'b'
	sessionFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 98 {
			pages.SwitchToPage("Main")
		}
		return event
	})

	return sessionFlex
}
