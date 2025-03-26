package main

import (
	"fmt"

	"github.com/fritzrepo/stockportfolio/config"
	"github.com/fritzrepo/stockportfolio/depot"
)

func main() {
	config, err := config.LoadConfigFromJSON("appConfig.json")
	if err != nil {
		// Fehlerbehandlung
		fmt.Println("Error loading config")
		panic(err)
	}

	depot, err := depot.ComputeTransactions(config.TransactionFilePath)
	if err != nil {
		// Fehlerbehandlung
		fmt.Println("Error computing transactions")
		panic(err)
	}

	fmt.Println("Depot:")
	fmt.Println(depot)
	fmt.Println("End")
}
