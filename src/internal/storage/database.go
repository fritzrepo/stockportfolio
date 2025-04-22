package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() {
	db, err := sql.Open("sqlite3", "../../data/depot.sqlite")
	if err != nil {
		log.Panic(err)
	}

	sqlStmt := "CREATE TABLE transactions (id TEXT not null primary key, date DATETIME, transactionType TEXT, isClosed INTEGER, " +
		"assetType TEXT, asset TEXT, tickerSymbol TEXT, quantity INTEGER, price INTEGER, fees REAL, currency TEXT);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Panic(err)
	}

	sqlStmt = "CREATE UNIQUE INDEX idx_transactions_id ON transactions(id);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Panic(err)
	}
}

func InsertTransaction(Transaction *Transaction) {

	db, err := sql.Open("sqlite3", "../../data/depot.sqlite")
	if err != nil {
		log.Panic(err)
	}

	sqlStmt := "INSERT INTO transactions (id, date, transactionType, isClosed, assetType, asset, tickerSymbol, quantity, price, fees, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = db.Exec(sqlStmt,
		Transaction.Id,
		Transaction.Date,
		Transaction.TransactionType,
		Transaction.IsClosed,
		Transaction.AssetType,
		Transaction.Asset,
		Transaction.TickerSymbol,
		Transaction.Quantity,
		Transaction.Price,
		Transaction.Fees,
		Transaction.Currency)
	if err != nil {
		log.Panic(err)
	}

}
