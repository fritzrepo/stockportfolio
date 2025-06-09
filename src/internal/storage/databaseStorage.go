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
	sqlStmt = "CREATE TABLE transactions (id TEXT not null primary key, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
		"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity INTEGER, price INTEGER, fees REAL, currency TEXT);"
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
	sqlStmt = "CREATE TABLE unclosed_assets (asset_id INTEGER PRIMARY KEY AUTOINCREMENT, asset_name TEXT UNIQUE NOT NULL);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create table unclosed_assets. %w", err)
	}

	// Create the unclosed_transactions table
	sqlStmt = "CREATE TABLE unclosed_trans (unclosed_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"asset_id INTEGER NOT NULL, " +
		"transaction_id TEXT, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
		"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity INTEGER, price INTEGER, fees REAL, currency TEXT, " +
		"FOREIGN KEY (asset_id) REFERENCES unclosed_assets(asset_id) ON DELETE CASCADE);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error at create table unclosed_trans. %w", err)
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
	sqlStmt := "INSERT OR IGNORE INTO unclosed_assets (asset_name) VALUES (?);"
	// Modern
	// sqlStmt := "INSERT INTO unclosed_assets (asset_name) VALUES (?) ON CONFLICT(asset_name) DO NOTHING;"
	// Beides wird von SQLite unterstützt
	_, err := db.Exec(sqlStmt, trans.Asset)

	if err != nil {

		return err
	}

	// Get the asset_id from the asset_name table
	var assetId int
	sqlStmt = "SELECT asset_id FROM unclosed_assets WHERE asset_name = ?;"
	err = db.QueryRow(sqlStmt, trans.Asset).Scan(&assetId)

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

func (s *DatabaseStorage) readUnclosedAssetNames(db *sql.DB) ([]string, error) {
	assetNames := make([]string, 0)

	sqlStmt := "SELECT asset_name FROM unclosed_assets;"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error at read unclosed asset names. %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var assetName string
		err = rows.Scan(
			&assetName)
		if err != nil {
			return nil, err
		}
		assetNames = append(assetNames, assetName)
	}
	return assetNames, nil
}
