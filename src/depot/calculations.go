package depot

func calculateProfitLoss(sellTrans Transaction, buyTransaction Transaction) RealizedGain {
	result := RealizedGain{}
	result.SellTransactionID = sellTrans.id
	result.BuyTransactions = buyTransaction.id
	result.Asset = sellTrans.asset
	//So könnte die Steuerberechnung aussehen, von Copilot vorgeschlagen
	// if sellTrans.date.Before(buyTransaction.date) {
	// 	result.TaxRate = 0.25
	// } else {
	// 	result.TaxRate = 0.15
	// }
	if sellTrans.quantity > buyTransaction.quantity {
		result.Quantity = buyTransaction.quantity
	} else {
		result.Quantity = sellTrans.quantity
	}
	result.BuyPrice = buyTransaction.price
	result.SellPrice = sellTrans.price
	result.Amount = sellTrans.quantity * (sellTrans.price - buyTransaction.price)
	result.IsProfit = result.Amount > 0
	return result
}
