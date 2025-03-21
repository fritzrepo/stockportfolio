package depot

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

type Transaction struct {
	id              uuid.UUID
	date            time.Time
	transactionType string //besser enum buy, sell
	IsClosed        bool
	assetType       string //stock, crypto, forex
	asset           string
	tickerSymbol    string
	quantity        float32
	price           float32
	fees            float32
	totalPrice      float32
	currency        currency.Unit
}

type DepotEntry struct {
	assetType    string
	asset        string
	tickerSymbol string
	quantity     float32
	price        float32
	totalPrice   float32
	currency     currency.Unit
}

type RealizedGain struct {
	SellTransactionID uuid.UUID   // ID der Verkaufstransaktion
	BuyTransactions   []uuid.UUID // Liste von Kauf-IDs, die diesem Verkauf zugeordnet sind
	Amount            float64     // Der Gewinn/Verlust-Betrag
	IsProfit          bool        // true für Gewinn, false für Verlust
	TaxRate           float64     // Anwendbarer Steuersatz
	Quantity          int
	BuyPrice          float64
	SellPrice         float64
}

// func calculateProfitLoss(sellTrans Transaction, buyTransactions []Transaction) RealizedGain {
// 	//ToDo => Implementieren
// 	//FiFo-Prinzip (First in, first out)

// 	//Durchschnittskostenmethode (average cost)
// }

func ComputeTransactions() map[string]DepotEntry {
	transactions := []Transaction{} // Ein Slice für Transaktionen
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 10, price: 100.5, fees: 4, totalPrice: 1009, currency: currency.EUR, id: uuid.New()})
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 20, price: 100.5, fees: 4, totalPrice: 2014, currency: currency.EUR, id: uuid.New()})
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "BASF", tickerSymbol: "BAS1",
		quantity: 100, price: 45.5, fees: 5, totalPrice: 4555, currency: currency.EUR, id: uuid.New()})

	unclosedTransactions := make(map[string][]Transaction)
	depotEntries := make(map[string]DepotEntry)
	//realizedGains := make([]RealizedGain, 10)

	for _, value := range transactions {
		fmt.Println(value.asset)

		//Handle die Asset-Transaktionen
		//Buy transaction
		if value.transactionType == "buy" {
			transaction, exists := unclosedTransactions[value.asset]
			//Gibt es schon Transaktionen für das Asset? Dann füge die neue Transaktion hinzu.
			if exists {
				transaction = append(transaction, value)
				unclosedTransactions[value.asset] = transaction
			} else {
				//Gibt es noch keine Transaktionen für das Asset? Dann erstelle einen neuen Slice mit der Transaktion.
				unclosedTransaction := make([]Transaction, 1)
				unclosedTransaction = append(unclosedTransaction, value)
				unclosedTransactions[value.asset] = unclosedTransaction
			}
		}

		//Sell transaction
		if value.transactionType == "sell" {
			transactions, exists := unclosedTransactions[value.asset]
			if exists {
				//Ziehe die Anzahl der verkauften Assets von der ersten buy Transaktion ab
				//Sollten mehr Assets verkauft werden, als gekauft wurden, dann wird
				//die nächste buy Transaktion noch verwendet.
				//Code wurde erstellt
				for i, transaction := range transactions {
					if transaction.transactionType == "buy" {
						//Buy und sell Transaktionen sind gleich
						if transaction.quantity == value.quantity {
							//Entferne die buy Transaktion aus dem Slice
							transactions = append(transactions[:i], transactions[i+1:]...)
							//Berechne den Gewinn / Verlust

							break
						}
						//Buy Transaktion ist größer als die Sell Transaktion
						if transaction.quantity > value.quantity {
							//Buy Transaktion verkleinern um die Anzahl der verkauften Assets
							transaction.quantity -= value.quantity
							transaction.totalPrice -= value.totalPrice
							transactions[i] = transaction
							//Berechne den Gewinn / Verlust

							break
						}
						//Buy Transaktion ist kleiner als die Sell Transaktion
						if transaction.quantity < value.quantity {
							value.quantity -= transaction.quantity
							value.totalPrice -= transaction.totalPrice
							//Entferne die Transaktion aus dem Slice
							transactions = append(transactions[:i], transactions[i+1:]...)
							//Berechne den Gewinn / Verlust

							//Sell transaction muss auf die nächste buy transaction angewendet werden
							//Daher kein break
						}
					}
				}
				unclosedTransactions[value.asset] = transactions
			} else {
				//ToDo => Richtige Fehlermeldung implementieren
				fmt.Println("Asset not found in	depot")
			}
		}
	}

	//Handle depot
	for _, transactions := range unclosedTransactions {
		for _, transaction := range transactions {
			//Wenn das Asset noch nicht im Depot ist, dann füge es hinzu
			entry, exists := depotEntries[transaction.asset]
			if !exists {

				depotEntries[transaction.asset] = DepotEntry{assetType: transaction.assetType, asset: transaction.asset,
					tickerSymbol: transaction.tickerSymbol, quantity: transaction.quantity, price: transaction.price,
					totalPrice: transaction.totalPrice, currency: transaction.currency}
			} else {
				//Wenn das Asset schon im Depot ist, dann aktualisiere die Anzahl und den Preis
				entry.quantity += transaction.quantity
				entry.totalPrice += transaction.totalPrice
				depotEntries[transaction.asset] = entry
			}
		}
	}

	//fmt.Println(depot)
	fmt.Println(depotEntries)
	fmt.Println(transactions[0].date.Format("2006-01-02"))
	return depotEntries
}
