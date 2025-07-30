package portfolio

import "github.com/fritzrepo/stockportfolio/internal/storage"

type Portfolio interface {
	AddTransaction(transaction storage.Transaction) error
}
