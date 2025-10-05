package storage

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type FileDatabase struct {
	baseDb   DatabaseStorage
	filePath string
}

func GetFileDatabase(pathToFile string, uuidGen func() uuid.UUID) *FileDatabase {

	var fileDB = &FileDatabase{
		filePath: pathToFile,
	}
	fileDB.baseDb.uuidGenerator = uuidGen

	return fileDB
}

func (s *FileDatabase) CreateDatabase() error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.createDatabase(db)
	})
}

func (s *FileDatabase) Ping() error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.ping(db)
	})
}

func (s *FileDatabase) AddTransaction(transaction *Transaction) error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.insertTransaction(db, transaction)
	})
}

func (s *FileDatabase) ReadAllTransactions() ([]Transaction, error) {
	var transactions []Transaction

	err := s.withDatabase(func(db *sql.DB) error {
		var errorSql error
		transactions, errorSql = s.baseDb.loadAllTransactions(db)
		return errorSql
	})

	if err != nil {
		return nil, err
	}
	return transactions, err
}

func (s *FileDatabase) AddUnclosedTransaction(asset Transaction) error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.insertUnclosedTransaction(db, asset)
	})
}

func (s *FileDatabase) RemoveAllUnclosedTransactions() error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.deleteAllUnclosedTransaction(db)
	})
}

func (s *FileDatabase) ReadAllUnclosedTickerSymbols() ([]string, error) {
	var tickerSymbols []string

	err := s.withDatabase(func(db *sql.DB) error {
		var errorSql error
		tickerSymbols, errorSql = s.baseDb.loadUnclosedTickerSymbols(db)
		return errorSql
	})

	if err != nil {
		return nil, err
	}
	return tickerSymbols, nil
}

func (s *FileDatabase) ReadAllUnclosedTransactions() (map[string][]Transaction, error) {
	var unclosedTransactions map[string][]Transaction

	err := s.withDatabase(func(db *sql.DB) error {
		var errorSql error
		unclosedTransactions, errorSql = s.baseDb.loadUnclosedTransactions(db)
		return errorSql
	})

	if err != nil {
		return nil, err
	}
	return unclosedTransactions, nil
}

func (s *FileDatabase) AddRealizedGain(realizedGain RealizedGain) error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.insertRealizedGain(db, &realizedGain)
	})
}

func (s *FileDatabase) ReadAllRealizedGains() ([]RealizedGain, error) {
	var realizedGains []RealizedGain

	err := s.withDatabase(func(db *sql.DB) error {
		var errorSql error
		realizedGains, errorSql = s.baseDb.loadAllRealizedGains(db)
		return errorSql
	})

	if err != nil {
		return nil, err
	}
	return realizedGains, nil
}

func (s *FileDatabase) withDatabase(action func(db *sql.DB) error) error {
	dbPath := s.filePath

	// Datenbankverbindung öffnen
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close() // Verbindung sicher schließen

	// Aktion mit der Datenbank ausführen
	return action(db)
}
