package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type MemoryDatabase struct {
	baseDb DatabaseStorage
	db     *sql.DB
}

func GetMemoryDatabase(uuidGen func() uuid.UUID) MemoryDatabase {
	var memoryDB = &MemoryDatabase{
		db: nil,
	}
	memoryDB.baseDb.uuidGenerator = uuidGen
	return *memoryDB
}

func (s *MemoryDatabase) CreateDatabase() error {
	return s.baseDb.createDatabase(s.db)
}

func (s *MemoryDatabase) Ping() error {
	return s.baseDb.ping(s.db)
}

func (s *MemoryDatabase) Open() error {
	var err error
	s.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	return nil
}

func (s *MemoryDatabase) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return errors.New("database instance is nil, cannot close")
}

func (s *MemoryDatabase) AddTransaction(transaction *Transaction) error {
	return s.baseDb.insertTransaction(s.db, transaction)
}

func (s *MemoryDatabase) ReadAllTransactions() ([]Transaction, error) {
	return s.baseDb.loadAllTransactions(s.db)
}

func (s *MemoryDatabase) AddUnclosedTransaction(asset Transaction) error {
	return s.baseDb.insertUnclosedTransaction(s.db, asset)
}

func (s *MemoryDatabase) ReadAllUnclosedTransactions() (map[string][]Transaction, error) {
	return s.baseDb.loadUnclosedTransactions(s.db)
}

// Wird eigentlich nicht benötigt.
func (s *MemoryDatabase) ReadAllUnclosedTickerSymbols() ([]string, error) {
	return s.baseDb.loadUnclosedTickerSymbols(s.db)
}

func (s *MemoryDatabase) AddRealizedGain(realizedGain RealizedGain) error {
	return s.baseDb.insertRealizedGain(s.db, &realizedGain)
}

func (s *MemoryDatabase) ReadAllRealizedGains() ([]RealizedGain, error) {
	return s.baseDb.loadAllRealizedGains(s.db)
}
