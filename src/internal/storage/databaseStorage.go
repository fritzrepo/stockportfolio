package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type DatabaseStorage struct {
	uuidGenerator func() uuid.UUID
}

func (s *DatabaseStorage) createDatabase(db *sql.DB) error {

	sqlStmt := "CREATE TABLE transactions (id TEXT not null primary key, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
		"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity INTEGER, price INTEGER, fees REAL, currency TEXT);"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = "CREATE UNIQUE INDEX idx_transactions_id ON transactions(id);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
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
