package handlers

import (
	"log"
	"net/http"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/depot"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApiResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	ErrorMessage string      `json:"error_message"`
	ErrorDetails string      `json:"error_details"`
	Data         interface{} `json:"data,omitempty"`
}

// PingHandler returns a simple JSON response with a message and the database file path
// With closure to access appConfig
func PingHandler(appConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"path":    appConfig.DatabaseFilePath,
		})
	}
}

func AddTransactionHandler(portfolio *depot.Depot) gin.HandlerFunc {
	return func(c *gin.Context) {

		response := &ApiResponse{
			Status:       "success",
			Message:      "Transaction added successfully",
			ErrorMessage: "",
			ErrorDetails: "",
			Data:         nil,
		}

		var transaction storage.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			response.Status = "error"
			response.Message = "Failed to add transaction"
			response.ErrorMessage = "Invalid request body"
			response.ErrorDetails = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}

		transaction.Id = uuid.New()
		log.Printf("Received transaction: %+v\n", transaction)
		err := portfolio.AddTransaction(transaction)
		if err != nil {
			log.Printf("Error adding transaction: %v\n", err)
			response.Status = "error"
			response.Message = "Failed to add transaction"
			response.ErrorDetails = err.Error()
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		log.Printf("Transaction added successfully: %+v\n", transaction)
		response.Data = transaction

		c.JSON(http.StatusCreated, response)
	}
}
