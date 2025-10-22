package storage

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id              uuid.UUID
	Date            time.Time `json:"date" xml:"dat" binding:"required"`
	TransactionType string    `json:"transactionType" xml:"transactionType" binding:"required"` // buy, sell
	AssetType       string    `json:"assetType" xml:"assetType" binding:"required"`             //stock, crypto, forex
	Asset           string    `json:"asset" xml:"asset" binding:"required"`
	DepotName       string    `json:"depotName" xml:"depotName" binding:"required"`
	TickerSymbol    string    `json:"tickerSymbol" xml:"tickerSymbol" binding:"required"`
	Quantity        float64   `json:"quantity" xml:"quantity" binding:"required"` //float64, um kombatibel mit der SQLite Datenbank zu sein.
	Price           float64   `json:"price" xml:"price" binding:"required"`
	Fees            float64   `json:"fees" xml:"fees" binding:"required"`
	Currency        string    `json:"currency" xml:"currency" binding:"required"`
}
