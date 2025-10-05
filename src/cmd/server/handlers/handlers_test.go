package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/portfolio"
	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/gin-gonic/gin"
)

// mockDepot implements the AddTransaction method for testing
type mockDepot struct {
	addTransaction func(storage.Transaction) error
	getEntries     func() map[string]portfolio.DepotEntry
}

func (m *mockDepot) AddTransaction(t storage.Transaction) error {
	return m.addTransaction(t)
}

func (m *mockDepot) GetEntries() map[string]portfolio.DepotEntry {
	return m.getEntries()
}

func TestAddTransactionHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockDepot{
		addTransaction: func(tr storage.Transaction) error {
			return nil
		},
		getEntries: func() map[string]portfolio.DepotEntry {
			return make(map[string]portfolio.DepotEntry)
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
