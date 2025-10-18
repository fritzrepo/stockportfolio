package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/portfolio"
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

func GetEntries(depot portfolio.Portfolio) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := &ApiResponse{
			Status:       "success",
			Message:      "Depot entries loaded",
			ErrorMessage: "",
			ErrorDetails: "",
			Data:         depot.GetEntries(),
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetRealizedGains(depot portfolio.Portfolio) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := depot.GetAllRealizedGains()
		if err != nil {
			response := &ApiResponse{
				Status:       "error",
				Message:      "",
				ErrorMessage: "Could not retrieve realized gains",
				ErrorDetails: err.Error(),
				Data:         nil,
			}
			c.JSON(http.StatusOK, response)
			return
		}
		response := &ApiResponse{
			Status:       "success",
			Message:      "Realized gains loaded",
			ErrorMessage: "",
			ErrorDetails: "",
			Data:         data,
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetPerformanceHandler(depot portfolio.Portfolio) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := depot.GetPerformance()
		if err != nil {
			response := &ApiResponse{
				Status:       "error",
				Message:      "",
				ErrorMessage: "Could not retrieve performance data",
				ErrorDetails: err.Error(),
				Data:         nil,
			}
			c.JSON(http.StatusOK, response)
			return
		}
		response := &ApiResponse{
			Status:       "success",
			Message:      "Performance data loaded",
			ErrorMessage: "",
			ErrorDetails: "",
			Data:         data,
		}
		c.JSON(http.StatusOK, response)
	}
}

func AddTransactionHandler(depot portfolio.Portfolio) gin.HandlerFunc {
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
			c.JSON(http.StatusOK, response)
			return
		}

		transaction.Id = uuid.New()
		log.Printf("Received transaction: %+v\n", transaction)
		err := depot.AddTransaction(transaction)
		if err != nil {
			log.Printf("Error adding transaction: %v\n", err)
			response.Status = "error"
			response.Message = "Failed to add transaction"
			response.ErrorDetails = err.Error()
			c.JSON(http.StatusOK, response)
			return
		}

		log.Printf("Transaction added successfully: %+v\n", transaction)
		response.Data = transaction

		c.JSON(http.StatusOK, response)
	}
}

func GetAllTransactionsHandler(depot portfolio.Portfolio) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := depot.GetAllTransactions()
		if err != nil {
			response := &ApiResponse{
				Status:       "error",
				Message:      "",
				ErrorMessage: "Could not retrieve transactions",
				ErrorDetails: err.Error(),
				Data:         nil,
			}
			c.JSON(http.StatusOK, response)
			return
		}

		accept := c.GetHeader("Accept")
		// JSON by default or when Accept contains application/json
		if accept == "" || strings.Contains(accept, "application/json") {
			response := &ApiResponse{
				Status:       "success",
				Message:      "Transactions loaded",
				ErrorMessage: "",
				ErrorDetails: "",
				Data:         data,
			}
			c.JSON(http.StatusOK, response)
			return
		}

		// Fallback: plain text output
		var b strings.Builder
		for i, t := range data {
			if i > 0 {
				b.WriteString("\n")
			}
			// Format nach Bedarf anpassen: hier einige Standardfelder
			//b.WriteString(fmt.Sprintf("Date: %s | Type: %s | AssetType: %s | Asset: %s | Ticker: %s | Qty: %v | Price: %v | Fees: %v | Currency: %s",
			b.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s;%v;%v;%v;%s",
				// t.Date.Format(time.RFC3339),
				t.Date.Format(time.DateOnly),
				t.TransactionType,
				t.AssetType,
				t.Asset,
				t.TickerSymbol,
				t.Quantity,
				t.Price,
				t.Fees,
				t.Currency))
		}
		// Optional: Setze Content-Disposition Header für Dateidownload
		//c.Header("Content-Disposition", "attachment; filename=\"datei.txt\"")
		c.String(http.StatusOK, b.String())
	}
}
