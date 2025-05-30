package main

import (
	"fmt"
	"os"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/depot"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/google/uuid"
)

func main() {
	var buildDb = false
	var fillDb = false
	var compute = false

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
		store := storage.NewCsvStorage(config.TransactionFilePath, uuid.New)
		transactions, err := store.LoadAllTransactions()
		if err != nil {
			// Fehlerbehandlung
			fmt.Println("Error loading transactions")
			panic(err)
		}

		dbStore := storage.GetFileDatabase(config.DatabaseFilePath, uuid.New)

		for _, transaction := range transactions {
			fmt.Println(transaction)
			err := dbStore.InsertTransaction(&transaction)
			if err != nil {
				fmt.Println("Database not created or already exists")
				panic(err)
			}
			fmt.Println("Inserted transaction")
		}
	}

	if compute {
		fmt.Println("Computing transactions")
		store := storage.NewCsvStorage(config.TransactionFilePath, uuid.New)

		// db, err := sql.Open("sqlite3", "../../data/depot.sqlite")
		// if err != nil {
		// 	log.Panic(err)
		// }
		// defer db.Close()

		dep := depot.NewDepot(uuid.New, &store)

		err = dep.ComputeAllTransactions()
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

}
