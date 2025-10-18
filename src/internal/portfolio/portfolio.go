package portfolio

import "github.com/fritzrepo/stockportfolio/internal/storage"

type Portfolio interface {
	GetEntries() map[string]DepotEntry
	AddTransaction(transaction storage.Transaction) error
	GetAllTransactions() ([]storage.Transaction, error)
	GetPerformance() (Performance, error)
	GetAllRealizedGains() ([]storage.RealizedGain, error)
}
