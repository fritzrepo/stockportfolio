package storage

type Store interface {
	Ping() error
	CreateDatabase() error
	AddTransaction(transaction *Transaction) error
	ExistsTransaction(transaction *Transaction) (*Transaction, error)
	ReadAllTransactions() ([]Transaction, error)
	RemoveAllTransactions() error
	AddUnclosedTransaction(asset Transaction) error
	RemoveAllUnclosedTransactions() error
	ReadAllUnclosedTickerSymbols() ([]string, error)
	ReadAllUnclosedTransactions() (map[string][]Transaction, error)
	AddRealizedGain(realizedGain RealizedGain) error
	ReadAllRealizedGains() ([]RealizedGain, error)
	RemoveAllRealizedGains() error
}
