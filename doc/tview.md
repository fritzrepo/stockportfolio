## tview
Man kann Listeneinträge eine Funktion zuweisen, die beim selektieren eines Eintrags ausgeführt wird.
```
contactsList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		setConcatText(&contacts[index])
	})
```


Man kann Tastatureingaben bestimmten Elemente zuordnen. So das die Eingaben nur wenn das Element vorhanden ist gefangen werden.
```
mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			formAddTransaction.Clear(true)
			addContactForm()
			pages.SwitchToPage("Add Contact")
		}
		return event
	})
```
