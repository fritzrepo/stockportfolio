package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type DatabaseStorage struct {
	uuidGenerator func() uuid.UUID
}

func (s *DatabaseStorage) createDatabase(db *sql.DB) error {

	// Fremdschlüssel-Unterstützung aktivieren
	sqlStmt := "PRAGMA foreign_keys = ON;"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at enable foreign key support. %w", err)
	}

	// Create the transactions table
	sqlStmt = "CREATE TABLE transactions (id TEXT(36) not null primary key, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
		"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity REAL, price REAL, fees REAL, currency TEXT);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create table transactions. %w", err)
	}

	sqlStmt = "CREATE UNIQUE INDEX idx_transactions_id ON transactions(id);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create index on table transactions. %w", err)
	}

	// Unclosed_transactions
	// 1:n asset -> unclosed_transactions

	// Create the unclosed_assets table
	sqlStmt = "CREATE TABLE unclosed_assets (asset_id INTEGER PRIMARY KEY AUTOINCREMENT, ticker_symbol TEXT UNIQUE NOT NULL);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create table unclosed_assets. %w", err)
	}

	// Create the unclosed_transactions table
	sqlStmt = "CREATE TABLE unclosed_trans (unclosed_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"asset_id INTEGER NOT NULL, " +
		"transaction_id TEXT, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
		"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity REAL, price REAL, fees REAL, currency TEXT, " +
		"FOREIGN KEY (asset_id) REFERENCES unclosed_assets(asset_id) ON DELETE CASCADE);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create table unclosed_trans. %w", err)
	}

	sqlStmt = "CREATE INDEX IF NOT EXISTS idx_nclosed_id ON unclosed_trans(unclosed_id)"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create index on table unclosed. %w", err)
	}

	// Create the RealizedGains table
	sqlStmt = "CREATE TABLE realized_gains (id TEXT(36) not null primary key, sellTransactionId TEXT(36), buyTransactionId TEXT(36), " +
		"asset TEXT, amount REAL, isProfit INTEGER, taxRate REAL, quantity REAL, buyPrice REAL, sellPrice REAL, currency TEXT, " +
		"FOREIGN KEY (sellTransactionId) REFERENCES transactions(id) ON DELETE CASCADE, " +
		"FOREIGN KEY (buyTransactionId) REFERENCES transactions(id) ON DELETE CASCADE);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create table realized_gains. %w", err)
	}

	return nil
}

func (s *DatabaseStorage) insertTransaction(db *sql.DB, transaction *Transaction) error {
	transaction.Id = s.uuidGenerator()
	sqlStmt := "INSERT INTO transactions (id, date, transactionType, isClosed, assetType, asset, tickerSymbol, quantity, price, fees, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := db.Exec(sqlStmt,
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
		return err
	}
	return nil
}

func (s *DatabaseStorage) loadAllTransactions(db *sql.DB) ([]Transaction, error) {

	transactions := make([]Transaction, 0)

	rows, err := db.Query("SELECT id, date, transactionType, isClosed, assetType, asset, tickerSymbol, quantity, price, fees, currency FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(
			&transaction.Id,
			&transaction.Date,
			&transaction.TransactionType,
			&transaction.IsClosed,
			&transaction.AssetType,
			&transaction.Asset,
			&transaction.TickerSymbol,
			&transaction.Quantity,
			&transaction.Price,
			&transaction.Fees,
			&transaction.Currency)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *DatabaseStorage) insertUnclosedTransaction(db *sql.DB, trans Transaction) error {

	// Save Asset-Name in unclosed_assets table
	// SQLite spezifisch
	sqlStmt := "INSERT OR IGNORE INTO unclosed_assets (ticker_symbol) VALUES (?);"
	// Modern
	// sqlStmt := "INSERT INTO unclosed_assets (ticker_symbol) VALUES (?) ON CONFLICT(ticker_symbol) DO NOTHING;"
	// Beides wird von SQLite unterstützt
	_, err := db.Exec(sqlStmt, trans.TickerSymbol)

	if err != nil {
		return err
	}

	// Get the asset_id from the unclosed_assets table
	var assetId int
	sqlStmt = "SELECT asset_id FROM unclosed_assets WHERE ticker_symbol = ?;"
	err = db.QueryRow(sqlStmt, trans.TickerSymbol).Scan(&assetId)

	if err != nil {
		return err
	}

	// Insert the transaction into unclosed
	sqlStmt = "INSERT INTO unclosed_trans (asset_id, transaction_id, date, transactionType, isClosed, assetType, asset, tickerSymbol, quantity, price, fees, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = db.Exec(sqlStmt,
		assetId,
		trans.Id,
		trans.Date,
		trans.TransactionType,
		trans.IsClosed,
		trans.AssetType,
		trans.Asset,
		trans.TickerSymbol,
		trans.Quantity,
		trans.Price,
		trans.Fees,
		trans.Currency)

	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseStorage) deleteAllUnclosedTransaction(db *sql.DB) error {
	sqlStmt := "DELETE FROM unclosed_trans;"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at delete all unclosed transactions. %w", err)
	}

	sqlStmt = "DELETE FROM unclosed_assets;"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at delete all unclosed assets. %w", err)
	}

	return nil
}

func (s *DatabaseStorage) loadUnclosedTickerSymbols(db *sql.DB) ([]string, error) {
	tickerSymbols := make([]string, 0)

	sqlStmt := "SELECT ticker_symbol FROM unclosed_assets;"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error at read unclosed ticker symbol. %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tickerSymbol string
		err = rows.Scan(
			&tickerSymbol)
		if err != nil {
			return nil, err
		}
		tickerSymbols = append(tickerSymbols, tickerSymbol)
	}
	return tickerSymbols, nil
}

func (s *DatabaseStorage) loadUnclosedTransactions(db *sql.DB) (map[string][]Transaction, error) {
	unclosedTransactions := make(map[string][]Transaction)
	var tickerSymbols []string

	tickerSymbols, err := s.loadUnclosedTickerSymbols(db)
	if err != nil {
		return nil, fmt.Errorf("error at read unclosed transactions. %w", err)
	}

	for _, tickerSymbol := range tickerSymbols {
		sqlStmt := `SELECT transaction_id, date, transactionType, isClosed, assetType, asset, tickerSymbol, 
		quantity, price, fees, currency FROM unclosed_trans 
		WHERE asset_id = (SELECT asset_id FROM unclosed_assets WHERE ticker_symbol = ?);`

		rows, err := db.Query(sqlStmt, tickerSymbol)
		if err != nil {
			return nil, fmt.Errorf("error at read unclosed transactions. %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var transaction Transaction
			err = rows.Scan(
				&transaction.Id,
				&transaction.Date,
				&transaction.TransactionType,
				&transaction.IsClosed,
				&transaction.AssetType,
				&transaction.Asset,
				&transaction.TickerSymbol,
				&transaction.Quantity,
				&transaction.Price,
				&transaction.Fees,
				&transaction.Currency)
			if err != nil {
				return nil, err
			}
			unclosedTransactions[tickerSymbol] = append(unclosedTransactions[tickerSymbol], transaction)
		}
	}
	return unclosedTransactions, nil
}

func (s *DatabaseStorage) insertRealizedGain(db *sql.DB, realizedGain *RealizedGain) error {
	realizedGain.Id = s.uuidGenerator()
	sqlStmt := "INSERT INTO realized_gains (id, sellTransactionId, buyTransactionId, asset, amount, isProfit, taxRate, quantity, buyPrice, sellPrice, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := db.Exec(sqlStmt,
		realizedGain.Id,
		realizedGain.SellTransactionId,
		realizedGain.BuyTransactionId,
		realizedGain.Asset,
		realizedGain.Amount,
		realizedGain.IsProfit,
		realizedGain.TaxRate,
		realizedGain.Quantity,
		realizedGain.BuyPrice,
		realizedGain.SellPrice,
		realizedGain.Currency)
	if err != nil {
		return err
	}
	return nil
}

func (s *DatabaseStorage) loadAllRealizedGains(db *sql.DB) ([]RealizedGain, error) {
	realizedGains := make([]RealizedGain, 0)

	rows, err := db.Query("SELECT id, sellTransactionId, buyTransactionId, asset, amount, isProfit, taxRate, quantity, buyPrice, sellPrice, currency FROM realized_gains")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var realizedGain RealizedGain
		err = rows.Scan(
			&realizedGain.Id,
			&realizedGain.SellTransactionId,
			&realizedGain.BuyTransactionId,
			&realizedGain.Asset,
			&realizedGain.Amount,
			&realizedGain.IsProfit,
			&realizedGain.TaxRate,
			&realizedGain.Quantity,
			&realizedGain.BuyPrice,
			&realizedGain.SellPrice,
			&realizedGain.Currency)
		if err != nil {
			return nil, err
		}
		realizedGains = append(realizedGains, realizedGain)
	}
	return realizedGains, nil
}

func (s *DatabaseStorage) ping(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("error at ping database. %w", err)
	}
	return nil
}
