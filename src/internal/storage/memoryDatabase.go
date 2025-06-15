package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type MemoryDatabase struct {
	baseDb DatabaseStorage
	db     *sql.DB
}

func (s *MemoryDatabase) CreateDatabase() error {
	err := s.baseDb.createDatabase(s.db)
	if err != nil {
		return err
	}
	return nil
}

func GetMemoryDatabase(uuidGen func() uuid.UUID) MemoryDatabase {
	var memoryDB = MemoryDatabase{
		db: nil,
	}
	memoryDB.baseDb.uuidGenerator = uuidGen
	return memoryDB
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
	return nil
}

func (s *MemoryDatabase) InsertTransaction(transaction *Transaction) error {
	return s.baseDb.insertTransaction(s.db, transaction)
}

func (s *MemoryDatabase) LoadAllTransactions() ([]Transaction, error) {
	transactions, err := s.baseDb.loadAllTransactions(s.db)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *MemoryDatabase) InsertUnclosedTransaction(asset Transaction) error {
	return s.baseDb.insertUnclosedTransaction(s.db, asset)
}

func (s *MemoryDatabase) LoadAllUnclosedTransactions() (map[string][]Transaction, error) {
	unclosedTransactions, err := s.baseDb.readUnclosedTransactions(s.db)
	if err != nil {
		return nil, err
	}
	return unclosedTransactions, nil
}

// Wird eigentlich nicht benötigt.
func (s *MemoryDatabase) LoadAllUnclosedTickerSymbols() ([]string, error) {
	tickerSymbols, err := s.baseDb.readUnclosedTickerSymbol(s.db)
	if err != nil {
		return nil, err
	}
	return tickerSymbols, nil
}

func (s *MemoryDatabase) InsertRealizedGain(realizedGain RealizedGain) error {
	return s.baseDb.insertRealizedGain(s.db, &realizedGain)
}
