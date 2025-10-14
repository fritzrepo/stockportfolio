package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestStore(t *testing.T) Store {
	store := GetMemoryDatabase()
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
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"}

	err := store.AddTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to insert transaction: %v", err)
	}

	transactions, err := store.ReadAllTransactions()
	if err != nil {
		t.Errorf("Failed to load transactions: %v", err)
	}

	if transactions[0].Date != transaction.Date ||
		transactions[0].TransactionType != transaction.TransactionType ||
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
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"}

	err := store.AddUnclosedTransaction(*transaction)
	if err != nil {
		t.Errorf("Failed to insert unclosed asset name: %v", err)
	}
	//Füge es nochmal ein, um zu testen, das es bei der Tabelle "unclosed_assets" zu keinem Insert-Fehler kommt.
	//Bzw. dass der Insert dann nicht durchgeführt wird.
	err = store.AddUnclosedTransaction(*transaction)
	if err != nil {
		t.Errorf("Failed to insert unclosed asset name: %v", err)
	}

	tickerSymbols, err := store.ReadAllUnclosedTickerSymbols()
	if err != nil {
		t.Errorf("Failed to load unclosed asset names: %v", err)
	}

	if len(tickerSymbols) != 1 || tickerSymbols[0] != transaction.TickerSymbol {
		t.Errorf("Expected unclosed asset name 'Apple', but got %v", tickerSymbols)
	}

	unclosedTransactions, err := store.ReadAllUnclosedTransactions()
	if err != nil {
		t.Errorf("Failed to load unclosed transactions: %v", err)
	}

	if len(unclosedTransactions) != 1 || len(unclosedTransactions[transaction.TickerSymbol]) != 2 {
		t.Errorf("Expected 2 unclosed transaction for 'Apple', but got %v", len(unclosedTransactions[transaction.TickerSymbol]))
	}

	//Füge ein dritte Transaktion für ein neues asset hinzu.
	transaction = &Transaction{
		Id:              uuid.New(),
		Date:            time.Date(2022, 11, 7, 12, 0, 0, 0, time.UTC),
		TransactionType: "buy",
		AssetType:       "stock",
		Asset:           "BASF",
		TickerSymbol:    "BAS1",
		Quantity:        20,
		Price:           99,
		Fees:            1.5,
		Currency:        "USD"}

	err = store.AddUnclosedTransaction(*transaction)
	if err != nil {
		t.Errorf("Failed to insert unclosed asset name: %v", err)
	}

	unclosedTransactions, err = store.ReadAllUnclosedTransactions()
	if err != nil {
		t.Errorf("Failed to load unclosed transactions: %v", err)
	}

	if len(unclosedTransactions) != 2 || len(unclosedTransactions[transaction.TickerSymbol]) != 1 {
		t.Errorf("Expected 2 unclosed transaction for 'Apple', but got %v", len(unclosedTransactions[transaction.TickerSymbol]))
	}
}

func TestInsertRealizedGains(t *testing.T) {
	store := setupTestStore(t)

	Id1 := uuid.New() //Wird für die Foreign Key-Referenz benötigt.
	Id2 := uuid.New()

	//Wird für die Foreign Key-Referenz benötigt.
	transaction := &Transaction{
		Id:              Id1,
		Date:            time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		TransactionType: "buy",
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"}

	err := store.AddTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to insert transaction: %v", err)
	}

	transaction = &Transaction{
		Id:              Id2,
		Date:            time.Date(2023, 11, 1, 12, 0, 0, 0, time.UTC),
		TransactionType: "sell",
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           200,
		Fees:            1.5,
		Currency:        "USD"}

	err = store.AddTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to insert transaction: %v", err)
	}

	gain := &RealizedGain{
		Id:                uuid.New(),
		SellTransactionId: Id2,
		BuyTransactionId:  Id1,
		Asset:             "Apple",
		Amount:            500.0,
		IsProfit:          true,
		TaxRate:           0.25,
		Quantity:          10,
		BuyPrice:          150.0,
		SellPrice:         200.0,
		Currency:          "USD",
	}

	err = store.AddRealizedGain(*gain)
	if err != nil {
		t.Errorf("Failed to insert realized gain: %v", err)
	}

	realizedGains, err := store.ReadAllRealizedGains()
	if err != nil {
		t.Errorf("Failed to load realized gains: %v", err)
	}

	if len(realizedGains) != 1 {
		t.Errorf("Expected 1 realized gain, but got %d", len(realizedGains))
	}
	//Hier traten bisher keine Rundungsfehler auf. Wenn doch, dann epsilon verwenden. Siehe DepotTest.
	if realizedGains[0].Asset != gain.Asset ||
		realizedGains[0].Amount != gain.Amount ||
		realizedGains[0].IsProfit != gain.IsProfit ||
		realizedGains[0].TaxRate != gain.TaxRate ||
		realizedGains[0].Quantity != gain.Quantity ||
		realizedGains[0].BuyPrice != gain.BuyPrice ||
		realizedGains[0].SellPrice != gain.SellPrice ||
		realizedGains[0].Currency != gain.Currency {
		t.Errorf("Expected %+v, but got %+v", gain, realizedGains[0])
	}
}

func TestLoadTransactionByParams(t *testing.T) {
	store := setupTestStore(t)

	transaction := &Transaction{
		Id:              uuid.New(),
		Date:            time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		TransactionType: "buy",
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"}

	err := store.AddTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to insert transaction: %v", err)
	}

	loadedTransaction, err := store.LoadTransactionByParams(transaction.Date, transaction.TransactionType, transaction.TickerSymbol)
	if err != nil {
		t.Errorf("Failed to load transaction by params: %v", err)
	}
	if loadedTransaction == nil {
		t.Error("Loaded transaction is nil")
		return
	}
	if *loadedTransaction != *transaction {
		t.Errorf("Loaded transaction does not match original: %+v != %+v", loadedTransaction, transaction)
	}

	//Teste das Laden einer nicht existierenden Transaktion
	loadedTransaction, err = store.LoadTransactionByParams(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), "sell", "MSFT")
	if err != nil {
		t.Errorf("Failed to load transaction by params: %v", err)
	}
	if loadedTransaction != nil {
		t.Errorf("Expected nil for non-existing transaction, but got: %+v", loadedTransaction)
	}
}
