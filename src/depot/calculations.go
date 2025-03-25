package depot

import (
	"github.com/fritzrepo/stockportfolio/models"
)

func calculateProfitLoss(sellTrans models.Transaction, buyTransaction models.Transaction) RealizedGain {
	result := RealizedGain{}
	result.SellTransactionID = sellTrans.Id
	result.BuyTransactions = buyTransaction.Id
	result.Asset = sellTrans.Asset
	//So könnte die Steuerberechnung aussehen, von Copilot vorgeschlagen
	// if sellTrans.date.Before(buyTransaction.date) {
	// 	result.TaxRate = 0.25
	// } else {
	// 	result.TaxRate = 0.15
	// }
	if sellTrans.Quantity > buyTransaction.Quantity {
		result.Quantity = buyTransaction.Quantity
	} else {
		result.Quantity = sellTrans.Quantity
	}
	result.BuyPrice = buyTransaction.Price
	result.SellPrice = sellTrans.Price
	result.Amount = sellTrans.Quantity * (sellTrans.Price - buyTransaction.Price)
	result.IsProfit = result.Amount > 0
	return result
}
