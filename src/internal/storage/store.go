package storage

type Store interface {
	Ping() error
	CreateDatabase() error
	AddTransaction(transaction *Transaction) error
	ReadAllTransactions() ([]Transaction, error)
	AddUnclosedTransaction(asset Transaction) error
	ReadAllUnclosedTickerSymbols() ([]string, error)
	ReadAllUnclosedTransactions() (map[string][]Transaction, error)
	AddRealizedGain(realizedGain RealizedGain) error
	ReadAllRealizedGains() ([]RealizedGain, error)
}
