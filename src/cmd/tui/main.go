package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fritzrepo/stockportfolio/cmd/tui/view"
	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
)

func main() {

	appConfig, err := config.LoadConfigFromJSON("../../configs/appConfig.json")
	if err != nil {
		fmt.Println("Error loading config")
		panic(err)
	}

	_, err = os.Stat(appConfig.DatabaseFilePath)
	dbNotExists := os.IsNotExist(err)
	if dbNotExists {
		fmt.Println("Database file does not exist.")
		panic("Database not exists")
	}

	store := storage.GetFileDatabase(appConfig.DatabaseFilePath)
	depot := portfolio.GetDepot(store)
	err = depot.CalculateSecuritiesAccountBalance()
	if err != nil {
		log.Fatalf("Failed to calculate securities account balance: %v", err)
		panic(err)
	} else {
		log.Println("Depot successful initialized.")
	}

	var app = view.GetView(depot)

	if err := app.Run(); err != nil {
		panic(err)
	}

}
