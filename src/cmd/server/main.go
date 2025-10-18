package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fritzrepo/stockportfolio/cmd/server/handlers"
	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/gin-gonic/gin"
)

var appConfig *config.Config

// Interface-Variablen sind bereits "Referenzen", deshalb kein *storage.Store
var store storage.Store
var depot *portfolio.Depot

func main() {
	router := gin.Default()

	router.GET("/ping", handlers.PingHandler(appConfig))
	router.GET("/api/depot/getentries", handlers.GetEntries(depot))
	router.GET("/api/depot/getperformance", handlers.GetPerformanceHandler(depot))
	router.GET("/api/depot/getrealizedgains", handlers.GetRealizedGains(depot))
	router.POST("/api/depot/addTransaction", handlers.AddTransactionHandler(depot))
	router.GET("/api/depot/getalltransactions", handlers.GetAllTransactionsHandler(depot))

	router.Run()
}

func init() {
	log.Println("Initializing server...")
	loadAppConfig()
	initializingStore()
	err := initializingDepot()
	if err != nil {
		log.Fatalf("Failed to initialize depot: %v", err)
		//Programm beenden
		os.Exit(1)
	}
	log.Println("Server initialized successfully.")
}

func loadAppConfig() {
	var err error
	appConfig, err = config.LoadConfigFromJSON("../../configs/appConfig.json")
	if err != nil {
		fmt.Println("Error loading config")
		panic(err)
	}
}

func initializingStore() {
	log.Println("Initializing store...")

	_, err := os.Stat(appConfig.DatabaseFilePath)
	dbNotExists := os.IsNotExist(err)

	store = storage.GetFileDatabase(appConfig.DatabaseFilePath)

	if dbNotExists {
		log.Println("Database file does not exist, creating a new one...")
		store.CreateDatabase()
	}

	err = store.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	} else {
		log.Println("Store successful initialized and connected to the database.")
	}
}

func initializingDepot() error {
	log.Println("Initializing depot...")
	depot = portfolio.GetDepot(store)
	err := depot.CalculateSecuritiesAccountBalance()
	if err != nil {
		log.Fatalf("Failed to calculate securities account balance: %v", err)
		return errors.New("failed to initialize depot")
	} else {
		log.Println("Depot successful initialized.")
	}
	return nil
}
