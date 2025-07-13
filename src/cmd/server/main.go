package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var appConfig *config.Config

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"path":    appConfig.DatabaseFilePath,
		})
	})

	router.POST("/api/depot/addTransaction", func(c *gin.Context) {
		var transaction storage.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		transaction.Id = uuid.New()
		log.Printf("Received transaction: %+v\n", transaction)

		response := &ApiResponse{
			Status:       "success",
			Message:      "Transaction added successfully",
			ErrorMessage: "",
			ErrorDetails: ""}

		response.Data = transaction
		c.JSON(http.StatusCreated, response)

	})

	router.Run()
}

func init() {
	log.Println("Initializing server...")
	loadAppConfig()
	initializingStore()
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

	store := storage.GetFileDatabase(appConfig.DatabaseFilePath, uuid.New)

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
