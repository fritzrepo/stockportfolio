package storage

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id              uuid.UUID
	Date            time.Time
	TransactionType string // buy, sell
	IsClosed        bool
	AssetType       string //stock, crypto, forex
	Asset           string
	TickerSymbol    string
	Quantity        float64 //float64, um kombatibel mit der SQLite Datenbank zu sein.
	Price           float64
	Fees            float64
	Currency        string
}

// TotalPrice berechnet und gibt den Gesamtpreis zurück
func (d *Transaction) TotalPrice() float64 {
	return d.Quantity * d.Price
}
