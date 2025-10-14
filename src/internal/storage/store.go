package storage

import "time"

type Store interface {
	Ping() error
	CreateDatabase() error
	AddTransaction(transaction *Transaction) error
	LoadTransactionByParams(date time.Time, transType string, tickSymbol string) (*Transaction, error)
	ReadAllTransactions() ([]Transaction, error)
	AddUnclosedTransaction(asset Transaction) error
	RemoveAllUnclosedTransactions() error
	ReadAllUnclosedTickerSymbols() ([]string, error)
	ReadAllUnclosedTransactions() (map[string][]Transaction, error)
	AddRealizedGain(realizedGain RealizedGain) error
	ReadAllRealizedGains() ([]RealizedGain, error)
	RemoveAllRealizedGains() error
}
