package storage

type Store interface {
	LoadAllTransactions() ([]Transaction, error)
}
