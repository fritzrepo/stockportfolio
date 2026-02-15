package banking

import (
	"time"
)

type Banking interface {
	StartSession() (string, error)
	EndSession() error
	RefreshTokenPeriodically(interval time.Duration) error
	GetAccountBalances() ([]Account, error)
}
