package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id              uuid.UUID
	Date            time.Time
	TransactionType string // besser enum buy, sell
	IsClosed        bool
	AssetType       string //stock, crypto, forex
	Asset           string
	TickerSymbol    string
	Quantity        float32
	Price           float32
	Fees            float32
	Currency        string
}

// TotalPrice berechnet und gibt den Gesamtpreis zurück
func (d *Transaction) TotalPrice() float32 {
	return d.Quantity * d.Price
}
