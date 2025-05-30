package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func TestInsertTransaction(t *testing.T) {
	// In-Memory SQLite-Datenbank
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Tabelle erstellen
	CreateDatabase := func(db *sql.DB) {
		sqlStmt := "CREATE TABLE transactions (id TEXT not null primary key, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
			"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity INTEGER, price INTEGER, fees REAL, currency TEXT);"
		_, err := db.Exec(sqlStmt)
		if err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
	}
	CreateDatabase(db)

	transaction := &Transaction{
		Id:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Date:            time.Now(),
		TransactionType: "buy",
		IsClosed:        false,
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD",
	}

	// Insert-Funktion testen
	sqlStmt := "INSERT INTO transactions (id, date, transactionType, isClosed, assetType, asset, tickerSymbol, quantity, price, fees, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = db.Exec(sqlStmt,
		transaction.Id,
		transaction.Date,
		transaction.TransactionType,
		transaction.IsClosed,
		transaction.AssetType,
		transaction.Asset,
		transaction.TickerSymbol,
		transaction.Quantity,
		transaction.Price,
		transaction.Fees,
		transaction.Currency)
	if err != nil {
		t.Fatalf("Failed to insert transaction: %v", err)
	}

	// Überprüfen, ob die Daten korrekt eingefügt wurden
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", transaction.Id).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query transaction: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 transaction, got %d", count)
	}
}

func TestGeneralDatabaseFunctions(t *testing.T) {
	store := NewMemoryDatabase(uuid.New)

	err := store.Open()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	defer store.Close()

	err = store.CreateDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

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

	err = store.InsertTransaction(transaction)
	if err != nil {
		t.Fatalf("Failed to insert transaction: %v", err)
	}

	transactions, err := store.LoadAllTransactions()
	if err != nil {
		t.Fatalf("Failed to load transactions: %v", err)
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
