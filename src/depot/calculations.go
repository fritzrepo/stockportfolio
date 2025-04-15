package depot

import (
	"github.com/fritzrepo/stockportfolio/models"
	"github.com/google/uuid"
)

func calculateProfitLoss(uuidGenerator func() uuid.UUID, sellTrans models.Transaction, buyTransaction models.Transaction) RealizedGain {
	result := RealizedGain{}
	result.Id = uuidGenerator()
	result.SellTransactionId = sellTrans.Id
	result.BuyTransactionId = buyTransaction.Id
	result.Asset = sellTrans.Asset
	//So könnte die Steuerberechnung aussehen, von Copilot vorgeschlagen
	// if sellTrans.date.Before(buyTransaction.date) {
	// 	result.TaxRate = 0.25
	// } else {
	// 	result.TaxRate = 0.15
	// }
	result.TaxRate = 0.25
	if sellTrans.Quantity > buyTransaction.Quantity {
		result.Quantity = buyTransaction.Quantity
	} else {
		result.Quantity = sellTrans.Quantity
	}
	result.BuyPrice = buyTransaction.Price
	result.SellPrice = sellTrans.Price
	result.Amount = calculateAmount(result.Quantity, buyTransaction.Price, sellTrans.Price)
	result.IsProfit = result.Amount > 0
	return result
}

// calculateAmount berechnet den Gewinn/Verlust-Betrag
// unter Verwendung von Ganzzahlen, um Rundungsfehler zu vermeiden.
func calculateAmount(quantity, buyPrice, sellPrice float32) float32 {
	// Skalieren auf ganze Zahlen (z.B. Cent)
	const scale = 100

	q := int64(quantity * scale)
	bp := int64(buyPrice * scale)
	sp := int64(sellPrice * scale)

	// Integer-Berechnung
	diff := sp - bp
	amount := (q * diff) / scale

	return float32(amount) / scale
}
