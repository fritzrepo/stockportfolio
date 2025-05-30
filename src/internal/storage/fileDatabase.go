package storage

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type FileDatabase struct {
	baseDb        DatabaseStorage
	filePath      string
	uuidGenerator func() uuid.UUID
}

func (s *FileDatabase) CreateDatabase() error {
	return s.withDatabase(func(db *sql.DB) error {
		err := s.baseDb.createDatabase(db)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *FileDatabase) InsertTransaction(transaction *Transaction) error {
	return s.withDatabase(func(db *sql.DB) error {
		return s.baseDb.insertTransaction(db, transaction)
	})
}

func (s *FileDatabase) LoadAllTransactions() ([]Transaction, error) {
	var transactions []Transaction
	var err error

	_ = s.withDatabase(func(db *sql.DB) error {
		transactions, err = s.baseDb.loadAllTransactions(db)
		return err
	})

	if err != nil {
		return nil, err
	}
	return transactions, err
}

func GetFileDatabase(pathToFile string, uuidGen func() uuid.UUID) FileDatabase {
	return FileDatabase{
		uuidGenerator: uuidGen,
		filePath:      pathToFile,
	}
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
