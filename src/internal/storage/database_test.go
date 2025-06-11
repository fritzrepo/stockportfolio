package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestStore(t *testing.T) MemoryDatabase {
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

func TestInsertTransaction(t *testing.T) {
	store := setupTestStore(t)

	transaction := &Transaction{
		Id:              uuid.New(),
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

func TestInsertUclosedTransaction(t *testing.T) {
	store := setupTestStore(t)

	transaction := &Transaction{
		Id:              uuid.New(),
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

	err := store.InsertUnclosedTransaction(*transaction)
	if err != nil {
		t.Errorf("Failed to insert unclosed asset name: %v", err)
	}
	//Füge es nochmal ein, um zu testen, das es bei der Tabelle "unclosed_assets" zu keinem Insert-Fehler kommt.
	//Bzw. dass der Insert dann nicht durchgeführt wird.
	err = store.InsertUnclosedTransaction(*transaction)
	if err != nil {
		t.Errorf("Failed to insert unclosed asset name: %v", err)
	}

	tickerSymbols, err := store.LoadAllUnclosedTickerSymbols()
	if err != nil {
		t.Errorf("Failed to load unclosed asset names: %v", err)
	}

	if len(tickerSymbols) != 1 || tickerSymbols[0] != transaction.TickerSymbol {
		t.Errorf("Expected unclosed asset name 'Apple', but got %v", tickerSymbols)
	}

	unclosedTransactions, err := store.LoadAllUnclosedTransactions()
	if err != nil {
		t.Errorf("Failed to load unclosed transactions: %v", err)
	}

	if len(unclosedTransactions) != 1 || len(unclosedTransactions[transaction.TickerSymbol]) != 2 {
		t.Errorf("Expected 2 unclosed transaction for 'Apple', but got %v", len(unclosedTransactions[transaction.TickerSymbol]))
	}
}
