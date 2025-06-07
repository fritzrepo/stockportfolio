package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestStore(t *testing.T) MemoryDatabase {
	// In-Memory SQLite-Datenbank
	store := GetMemoryDatabase(uuid.New)
	store.Open()
	err := store.CreateDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	// Cleanup registrieren. Wird nach jedem Test ausgeführt.
	t.Cleanup(func() {
		store.Close()
	})
	return store
}

func TestGeneralDatabaseFunctions(t *testing.T) {
	store := setupTestStore(t)

	transaction := &Transaction{
		Date:            time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		TransactionType: "buy",
		IsClosed:        false,
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"}

	err := store.InsertTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to insert transaction: %v", err)
	}

	transactions, err := store.LoadAllTransactions()
	if err != nil {
		t.Errorf("Failed to load transactions: %v", err)
	}

	if transactions[0].Date != transaction.Date ||
		transactions[0].TransactionType != transaction.TransactionType ||
		transactions[0].IsClosed != transaction.IsClosed ||
		transactions[0].AssetType != transaction.AssetType ||
		transactions[0].Asset != transaction.Asset ||
		transactions[0].TickerSymbol != transaction.TickerSymbol ||
		transactions[0].Quantity != transaction.Quantity ||
		transactions[0].Price != transaction.Price ||
		transactions[0].Fees != transaction.Fees ||
		transactions[0].Currency != transaction.Currency {
		t.Errorf("Expected %+v, but got %+v", transaction, transactions[0])
	}
}
