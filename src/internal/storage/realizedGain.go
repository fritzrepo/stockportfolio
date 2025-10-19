package storage

import "github.com/google/uuid"

// RealizedGain repräsentiert einen realisierten Gewinn oder Verlust aus dem Verkauf von Assets.
// Die Gebühren sind in der Berechnung des Amounts bereits berücksichtigt.
type RealizedGain struct {
	Id                uuid.UUID `json:"id"`                // ID der Realisierung
	SellTransactionId uuid.UUID `json:"sellTransactionId"` // ID der Verkaufstransaktion
	BuyTransactionId  uuid.UUID `json:"buytransactionId"`  // ID der Kauftransaktion
	Asset             string    // Asset-Name
	Amount            float64   // Der Gewinn/Verlust-Betrag
	IsProfit          bool      // true für Gewinn, false für Verlust
	TaxRate           float64   // Anwendbarer Steuersatz
	Quantity          float64
	BuyPrice          float64
	SellPrice         float64
	Currency          string
}
