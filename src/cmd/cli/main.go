package main

import (
	"fmt"
	"os"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/google/uuid"
)

func main() {
	var buildDb = false
	var fillDb = false
	var compute = false
	var readTransaktions = false

	// the first argument is always program name
	argLength := len(os.Args[1:])
	fmt.Printf("Arg length is %d\n", argLength)

	for i, a := range os.Args[1:] {
		fmt.Printf("Arg %d is %s\n", i+1, a)
		if a == "buildDb" {
			buildDb = true
		}
		if a == "fillDb" {
			fillDb = true
		}
		if a == "readTransactions" {
			readTransaktions = true
		}
		if a == "compute" {
			compute = true
		}
	}

	config, err := config.LoadConfigFromJSON("../../configs/appConfig.json")
	if err != nil {
		// Fehlerbehandlung
		fmt.Println("Error loading config")
		panic(err)
	}

	if buildDb {
		fmt.Println("Building database")
		store := storage.GetFileDatabase(config.DatabaseFilePath, uuid.New)
		err := store.CreateDatabase()
		if err != nil {
			fmt.Println("Database not created or already exists")
			panic(err)
		}
		fmt.Println("Database created")
	}

	if fillDb {
		fmt.Println("Fill up database")
		store := storage.GetCsvStorage(config.TransactionFilePath, uuid.New)
		transactions, err := store.ReadAllTransactions()
		if err != nil {
			// Fehlerbehandlung
			fmt.Println("Error loading transactions")
			panic(err)
		}

		dbStore := storage.GetFileDatabase(config.DatabaseFilePath, uuid.New)

		for _, transaction := range transactions {
			fmt.Println(transaction)
			err := dbStore.AddTransaction(&transaction)
			if err != nil {
				fmt.Println("Database not created or already exists")
				panic(err)
			}
			fmt.Println("Inserted transaction")
		}
	}

	if readTransaktions {
		fmt.Println("Reading transactions from database")
		store := storage.GetFileDatabase(config.DatabaseFilePath, uuid.New)

		transactions, err := store.ReadAllTransactions()
		if err != nil {
			// Fehlerbehandlung
			fmt.Println("Error reading transactions")
			panic(err)
		}

		for _, transaction := range transactions {
			fmt.Println(transaction)
		}
		fmt.Println("End of transactions")
		return
	}

	if compute {
		fmt.Println("Computing transactions")
		store := storage.GetCsvStorage(config.TransactionFilePath, uuid.New)

		dep := portfolio.GetDepot(uuid.New, &store)

		err = dep.ComputeAllTransactions()
		if err != nil {
			// Fehlerbehandlung
			fmt.Println("Error computing transactions")
			panic(err)
		}
		fmt.Println("Depot:")
		fmt.Println(dep.GetEntries())
		fmt.Println("Realized Gains:")
		fmt.Println(dep.RealizedGains)
		fmt.Println("End")
	}

}
