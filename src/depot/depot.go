package depot

import (
	"fmt"
	"slices"
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
	currency        currency.Unit
}

// TotalPrice berechnet und gibt den Gesamtpreis zurück
func (d *Transaction) TotalPrice() float32 {
	return d.quantity * d.price
}

type DepotEntry struct {
	assetType    string
	asset        string
	tickerSymbol string
	quantity     float32
	price        float32 //Kann ein Durchschnittspreis sein, wenn mehrere Transaktionen vorhanden sind
	currency     currency.Unit
}

// TotalPrice berechnet und gibt den Gesamtpreis zurück
func (d *DepotEntry) TotalPrice() float32 {
	return d.quantity * d.price
}

type RealizedGain struct {
	SellTransactionID uuid.UUID // ID der Verkaufstransaktion
	BuyTransactions   uuid.UUID // ID der Kauftransaktion
	Asset             string    // Asset-Name
	Amount            float32   // Der Gewinn/Verlust-Betrag
	IsProfit          bool      // true für Gewinn, false für Verlust
	TaxRate           float32   // Anwendbarer Steuersatz
	Quantity          float32
	BuyPrice          float32
	SellPrice         float32
}

func ComputeTransactions() map[string]DepotEntry {

	transactions := []Transaction{} // Ein Slice für alle Transaktionen
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 10, price: 100.5, fees: 4, currency: currency.EUR, id: uuid.New()})
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 20, price: 100.5, fees: 4, currency: currency.EUR, id: uuid.New()})
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "BASF", tickerSymbol: "BAS1",
		quantity: 100, price: 45.5, fees: 5, currency: currency.EUR, id: uuid.New()})
	transactions = append(transactions, Transaction{date: time.Now(), transactionType: "sell",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 10, price: 100.5, fees: 4, currency: currency.EUR, id: uuid.New()})

	unclosedTransactions := make(map[string][]Transaction)
	depotEntries := make(map[string]DepotEntry)
	realizedGains := make([]RealizedGain, 0, 5)

	for _, newTransaction := range transactions {
		fmt.Println(newTransaction.asset)

		//Handle die Asset-Transaktionen
		//Add buy transaction to unclosed transactionsavailableBuytrans
		if newTransaction.transactionType == "buy" {
			availableBuyTrans, exists := unclosedTransactions[newTransaction.asset]
			//Gibt es schon Transaktionen für das Asset? Dann füge die neue Transaktion hinzu.
			if exists {
				availableBuyTrans = append(availableBuyTrans, newTransaction)
				unclosedTransactions[newTransaction.asset] = availableBuyTrans
			} else {
				//Gibt es noch keine Transaktionen für das Asset? Dann erstelle einen neuen Slice mit der Transaktion.
				unclosedTransaction := []Transaction{}
				unclosedTransaction = append(unclosedTransaction, newTransaction)
				unclosedTransactions[newTransaction.asset] = unclosedTransaction
			}
		}

		//Sell transaction
		if newTransaction.transactionType == "sell" {
			transactions, exists := unclosedTransactions[newTransaction.asset]
			if exists {
				//FiFo-Prinzip (First in, first out)
				//Ziehe die Anzahl der verkauften Assets von der ersten buy Transaktion ab
				//Sollten mehr Assets verkauft werden, als gekauft wurden, dann wird
				//die nächste buy Transaktion noch verwendet.
				//Code wurde erstellt
				for i, availableBuyTrans := range transactions {
					if availableBuyTrans.transactionType == "buy" {
						//Buy und sell Transaktionen sind gleich
						if availableBuyTrans.quantity == newTransaction.quantity {
							//Entferne die buy Transaktion aus dem Slice
							transactions = slices.Delete(transactions, i, i+1)
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(newTransaction, availableBuyTrans))
							break
						}
						//Buy Transaktion ist größer als die Sell Transaktion
						if availableBuyTrans.quantity > newTransaction.quantity {
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(newTransaction, availableBuyTrans))
							//Buy Transaktion verkleinern um die Anzahl der verkauften Assets
							availableBuyTrans.quantity -= newTransaction.quantity
							//availableBuyTrans.totalPrice -= newTransaction.totalPrice //Total price könnte sich auch selbst berechnen
							transactions[i] = availableBuyTrans
							break
						}
						//Buy Transaktion ist kleiner als die Sell Transaktion
						if availableBuyTrans.quantity < newTransaction.quantity {
							newTransaction.quantity -= availableBuyTrans.quantity
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(newTransaction, availableBuyTrans))

							//newTransaction.totalPrice -= availableBuyTrans.totalPrice
							//Entferne die Transaktion aus dem Slice
							transactions = slices.Delete(transactions, i, i+1)

							//Der Rest der Sell transaction muss auf die nächste buy transaction angewendet werden
							//Daher kein break
						}
					} else {
						//ToDo => Richtige Fehlermeldung implementieren
						fmt.Println("Transaction is not a buy transaction")
					}
				}

				//Wennn die tansactions leer sind, dann lösche den Eintrag
				if len(transactions) == 0 {
					delete(unclosedTransactions, newTransaction.asset)
				} else {
					//Aktualisiere die Transaktionen
					unclosedTransactions[newTransaction.asset] = transactions
				}

				//Durchschnittskostenmethode (average cost)
				//ToDo => Implementieren
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

				// depotEntries[transaction.asset] = DepotEntry{assetType: transaction.assetType, asset: transaction.asset,
				// 	tickerSymbol: transaction.tickerSymbol, quantity: transaction.quantity, price: transaction.price,
				// 	totalPrice: transaction.totalPrice, currency: transaction.currency}
				depotEntries[transaction.asset] = DepotEntry{assetType: transaction.assetType, asset: transaction.asset,
					tickerSymbol: transaction.tickerSymbol, quantity: transaction.quantity, price: transaction.price,
					currency: transaction.currency}

			} else {
				//Wenn das Asset schon im Depot ist, dann aktualisiere die Anzahl und den Preis
				entry.quantity += transaction.quantity
				//entry.totalPrice += transaction.totalPrice
				depotEntries[transaction.asset] = entry
			}
		}
	}

	//fmt.Println(depot)
	fmt.Println(("Abgeschlossene Transaktionen"))
	fmt.Println(realizedGains)
	fmt.Println("Depot:")
	fmt.Println(depotEntries)
	fmt.Println(transactions[0].date.Format("2006-01-02"))
	return depotEntries
}
