package main

import (
	"fmt"
	"os"

	"github.com/fritzrepo/stockportfolio/cmd/tui/view"
	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}
}

func run() error {

	appConfig, err := config.LoadConfigFromJSON("./appConfig.json")
	if err != nil {
		return err
	}

	_, err = os.Stat(appConfig.DatabaseFilePath)
	dbNotExists := os.IsNotExist(err)
	if dbNotExists {
		return fmt.Errorf("database not exists: %s", appConfig.DatabaseFilePath)
	}

	store := storage.GetFileDatabase(appConfig.DatabaseFilePath)
	depot := portfolio.GetDepot(store)

	err = depot.CalculateSecuritiesAccountBalance()
	if err != nil {
		return err
	} else {
		fmt.Println("Depot successful initialized.")
	}

	var app = view.GetView(depot)
	return app.Run()
}
