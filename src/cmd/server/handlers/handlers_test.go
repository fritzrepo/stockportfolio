package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/config"
	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// mockDepot implements the AddTransaction method for testing
type mockDepot struct {
	addTransaction      func(storage.Transaction) error
	getEntries          func() map[string]portfolio.DepotEntry
	getAllRealizedGains func() ([]storage.RealizedGain, error)
	getPerformance      func() (portfolio.Performance, error)
	getAllTransactions  func() ([]storage.Transaction, error)
}

func (m *mockDepot) AddTransaction(t storage.Transaction) error {
	return m.addTransaction(t)
}

func (m *mockDepot) GetEntries() map[string]portfolio.DepotEntry {
	return m.getEntries()
}

func (m *mockDepot) GetAllRealizedGains() ([]storage.RealizedGain, error) {
	return m.getAllRealizedGains()
}

func (m *mockDepot) GetPerformance() (portfolio.Performance, error) {
	return m.getPerformance()
}

func (m *mockDepot) GetAllTransactions() ([]storage.Transaction, error) {
	return m.getAllTransactions()
}

func TestPingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	appConfig := &config.Config{
		DatabaseFilePath:    "test_db_path",
		TransactionFilePath: "",
	}

	router := gin.New()
	router.GET("/ping", PingHandler(appConfig))

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "pong" {
		t.Errorf("Expected message 'pong', got '%s'", response["message"])
	}

	if response["path"] != "test_db_path" {
		t.Errorf("Expected path 'test_db_path', got '%s'", response["path"])
	}
}

func TestAddTransactionHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		addTransaction: func(tr storage.Transaction) error {
			return nil
		},
	}

	router := gin.New()
	router.POST("/transaction", AddTransactionHandler(mock))

	// Prepare valid transaction JSON (without Id, as it is set in handler)
	tx := storage.Transaction{
		Date:            time.Now(),
		TransactionType: "buy",
		Asset:           "Apple Inc.",
		Currency:        "USD",
		Fees:            1.0,
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150.0,
		AssetType:       "buy",
	}
	body, _ := json.Marshal(tx)

	req, _ := http.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected success status, got %s", resp.Status)
	}

	if resp.Message != "Transaction added successfully" {
		t.Errorf("Expected success message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "" {
		t.Errorf("Expected no error message, got %s", resp.ErrorMessage)
	}

	if resp.Data == nil {
		t.Error("Expected data in response, got nil")
	}

	// Check that Id was set
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected Data to be a map[string]interface{}")
	}
	_, hasId := dataMap["Id"]
	if !hasId {
		t.Error("Expected Id to be present in response data")
	}
}

func TestAddTransactionHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		addTransaction: func(tr storage.Transaction) error {
			return nil
		},
	}

	router := gin.New()
	router.POST("/transaction", AddTransactionHandler(mock))

	// Invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/transaction", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "error" {
		t.Errorf("Expected error status, got %s", resp.Status)
	}

	if resp.Message != "Failed to add transaction" {
		t.Errorf("Expected error message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "Invalid request body" {
		t.Errorf("Expected error message 'Invalid request body', got %s", resp.ErrorMessage)
	}
}

func TestAddTransactionHandler_AddTransactionError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		addTransaction: func(tr storage.Transaction) error {
			return errors.New("db error")
		},
	}

	router := gin.New()
	router.POST("/transaction", AddTransactionHandler(mock))

	tx := storage.Transaction{
		Date:            time.Now(),
		TransactionType: "buy",
		Asset:           "Apple Inc.",
		Currency:        "USD",
		Fees:            1.0,
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150.0,
		AssetType:       "buy",
	}
	body, _ := json.Marshal(tx)

	req, _ := http.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "error" {
		t.Errorf("Expected error status, got %s", resp.Status)
	}

	if resp.Message != "Failed to add transaction" {
		t.Errorf("Expected error message 'Failed to add transaction', got %s", resp.Message)
	}

	if resp.ErrorDetails != "db error" {
		t.Error("Expected error details, got empty string")
	}

}

func TestGetEntriesHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getEntries: func() map[string]portfolio.DepotEntry {
			entries := make(map[string]portfolio.DepotEntry)
			entries["AAPL"] = portfolio.DepotEntry{
				TickerSymbol: "AAPL",
				Asset:        "Apple Inc.",
				AssetType:    "stock",
				Quantity:     50,
				Price:        145.00,
				Currency:     "USD",
			}
			return entries
		},
	}

	router := gin.New()
	router.GET("/getentries", GetEntries(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getentries", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected success status, got %s", resp.Status)
	}

	if resp.Message != "Depot entries loaded" {
		t.Errorf("Expected success message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "" {
		t.Errorf("Expected no error message, got %s", resp.ErrorMessage)
	}

	if resp.Data == nil {
		t.Error("Expected data in response, got nil")
	}

}

func TestGetRealizedGainsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getAllRealizedGains: func() ([]storage.RealizedGain, error) {
			gains := []storage.RealizedGain{
				{
					Id:                uuid.New(),
					BuyTransactionId:  uuid.New(),
					SellTransactionId: uuid.New(),
					Asset:             "Apple Inc.",
					Amount:            100.00,
					IsProfit:          true,
					TaxRate:           10,
					Quantity:          10,
					BuyPrice:          140.00,
					SellPrice:         150.00,
					Currency:          "USD",
				},
			}
			return gains, nil
		},
	}

	router := gin.New()
	router.GET("/getrealizedgains", GetRealizedGains(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getrealizedgains", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected success status, got %s", resp.Status)
	}

	if resp.Message != "Realized gains loaded" {
		t.Errorf("Expected success message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "" {
		t.Errorf("Expected no error message, got %s", resp.ErrorMessage)
	}

	if resp.Data == nil {
		t.Error("Expected data in response, got nil")
	}
}

func TestGetRealizedGainsHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getAllRealizedGains: func() ([]storage.RealizedGain, error) {
			return nil, errors.New("db error")
		},
	}

	router := gin.New()
	router.GET("/getrealizedgains", GetRealizedGains(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getrealizedgains", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "error" {
		t.Errorf("Expected error status, got %s", resp.Status)
	}

	if resp.Message != "" {
		t.Errorf("Expected empty message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "Could not retrieve realized gains" {
		t.Errorf("Expected error message 'Could not retrieve realized gains', got %s", resp.ErrorMessage)
	}

	if resp.ErrorDetails != "db error" {
		t.Errorf("Expected error details 'db error', got %s", resp.ErrorDetails)
	}
}

func TestGetPerformanceHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getPerformance: func() (portfolio.Performance, error) {
			performance := portfolio.Performance{
				TotalInvestedAmount:  10000.00,
				CountOfRealizedGains: 5,
				TotalGains:           1500.00,
				RealizedGains:        []storage.RealizedGain{},
			}
			return performance, nil
		},
	}

	router := gin.New()
	router.GET("/getperformance", GetPerformanceHandler(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getperformance", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected success status, got %s", resp.Status)
	}

	if resp.Message != "Performance data loaded" {
		t.Errorf("Expected success message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "" {
		t.Errorf("Expected no error message, got %s", resp.ErrorMessage)
	}

	if resp.Data == nil {
		t.Error("Expected data in response, got nil")
	}
}

func TestGetPerformance_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getPerformance: func() (portfolio.Performance, error) {
			return portfolio.Performance{}, errors.New("db error")
		},
	}

	router := gin.New()
	router.GET("/getperformance", GetPerformanceHandler(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getperformance", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "error" {
		t.Errorf("Expected error status, got %s", resp.Status)
	}

	if resp.Message != "" {
		t.Errorf("Expected empty message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "Could not retrieve performance data" {
		t.Errorf("Expected error message 'Could not retrieve performance data', got %s", resp.ErrorMessage)
	}

	if resp.ErrorDetails != "db error" {
		t.Errorf("Expected error details 'db error', got %s", resp.ErrorDetails)
	}
}

func TestGetAllTransactionsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getAllTransactions: func() ([]storage.Transaction, error) {
			transactions := []storage.Transaction{
				{
					Id:              uuid.New(),
					Date:            time.Now(),
					TransactionType: "buy",
					Asset:           "Apple Inc.",
					Currency:        "USD",
					Fees:            1.0,
					TickerSymbol:    "AAPL",
					Quantity:        10,
					Price:           150.0,
					AssetType:       "stock",
				},
			}
			return transactions, nil
		},
	}
	router := gin.New()
	router.GET("/getalltransactions", GetAllTransactionsHandler(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getalltransactions", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected success status, got %s", resp.Status)
	}

	if resp.Message != "Transactions loaded" {
		t.Errorf("Expected success message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "" {
		t.Errorf("Expected no error message, got %s", resp.ErrorMessage)
	}

	if resp.Data == nil {
		t.Error("Expected data in response, got nil")
	}
}

func TestGetAllTransactionsAsTextHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getAllTransactions: func() ([]storage.Transaction, error) {
			transactions := []storage.Transaction{
				{
					Id:              uuid.New(),
					Date:            time.Now(),
					TransactionType: "buy",
					Asset:           "Apple Inc.",
					Currency:        "USD",
					Fees:            1.0,
					TickerSymbol:    "AAPL",
					Quantity:        10,
					Price:           150.0,
					AssetType:       "stock",
				},
				{
					Id:              uuid.New(),
					Date:            time.Now(),
					TransactionType: "buy",
					Asset:           "BASF SE",
					Currency:        "EUR",
					Fees:            3.0,
					TickerSymbol:    "BAS",
					Quantity:        100,
					Price:           50.0,
					AssetType:       "stock",
				},
			}
			return transactions, nil
		},
	}
	router := gin.New()
	router.GET("/getalltransactions", GetAllTransactionsHandler(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getalltransactions", nil)
	req.Header.Set("Accept", "text/plain")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expectedContentType := "text/plain; charset=utf-8"
	if w.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type '%s', got '%s'", expectedContentType, w.Header().Get("Content-Type"))
	}

	if len(w.Body.Bytes()) == 0 {
		t.Error("Expected non-empty response body")
	}

	//Body should contain transaction details
	if !bytes.Contains(w.Body.Bytes(), []byte("Apple Inc.")) {
		t.Error("Expected response body to contain transaction details")
	}
}

func TestGetAllTransactionsHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		getAllTransactions: func() ([]storage.Transaction, error) {
			return nil, errors.New("db error")
		},
	}

	router := gin.New()
	router.GET("/getalltransactions", GetAllTransactionsHandler(mock))

	req, _ := http.NewRequest(http.MethodGet, "/getalltransactions", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "error" {
		t.Errorf("Expected error status, got %s", resp.Status)
	}

	if resp.Message != "" {
		t.Errorf("Expected empty message, got %s", resp.Message)
	}

	if resp.ErrorMessage != "Could not retrieve transactions" {
		t.Errorf("Expected error message 'Could not retrieve transactions', got %s", resp.ErrorMessage)
	}

	if resp.ErrorDetails != "db error" {
		t.Errorf("Expected error details 'db error', got %s", resp.ErrorDetails)
	}
}
