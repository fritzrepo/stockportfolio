package main

import (
	"fmt"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/depot"
	"github.com/google/uuid"
)

func main() {
	config, err := config.LoadConfigFromJSON("../../configs/appConfig.json")
	if err != nil {
		// Fehlerbehandlung
		fmt.Println("Error loading config")
		panic(err)
	}

	// var generateUUID = func() uuid.UUID {
	// 	return uuid.New() // Echte, zufällige UUID
	// }
	//Könnte sein das ich das hier nicht brauche
	dep := depot.NewDepot(uuid.New)

	err = dep.ComputeTransactions(config.TransactionFilePath)
	if err != nil {
		// Fehlerbehandlung
		fmt.Println("Error computing transactions")
		panic(err)
	}

	fmt.Println("Depot:")
	fmt.Println(dep.DepotEntries)
	fmt.Println("Realized Gains:")
	fmt.Println(dep.RealizedGains)
	fmt.Println("End")
}
