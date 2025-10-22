package storage

import "github.com/google/uuid"

// RealizedGain repräsentiert einen realisierten Gewinn oder Verlust aus dem Verkauf von Assets.
// Die Gebühren sind in der Berechnung des Amounts bereits berücksichtigt.
type RealizedGain struct {
	Id                uuid.UUID `json:"id"`                // ID der Realisierung
	SellTransactionId uuid.UUID `json:"sellTransactionId"` // ID der Verkaufstransaktion
	BuyTransactionId  uuid.UUID `json:"buytransactionId"`  // ID der Kauftransaktion
	Asset             string    `json:"asset"`             // Asset-Name
	DepotName         string    `json:"depotName"`         // Name des Depots
	Amount            float64   `json:"amount"`            // Der Gewinn/Verlust-Betrag
	IsProfit          bool      `json:"isProfit"`          // true für Gewinn, false für Verlust
	TaxRate           float64   `json:"taxRate"`           // Anwendbarer Steuersatz
	Quantity          float64   `json:"quantity"`          // Verkaufte Menge
	BuyPrice          float64   `json:"buyPrice"`          // Kaufpreis
	SellPrice         float64   `json:"sellPrice"`         // Verkaufspreis
	Currency          string    `json:"currency"`          // Währung
}
