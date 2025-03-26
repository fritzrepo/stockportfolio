package depot

import (
	"fmt"
	"slices"

	"github.com/fritzrepo/stockportfolio/depot/importer"
	"github.com/fritzrepo/stockportfolio/models"
	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

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

func ComputeTransactions(filePath string) (map[string]DepotEntry, error) {

	// transactions := []models.Transaction{} // Ein Slice für alle Transaktionen
	// transactions = append(transactions, models.Transaction{Date: time.Now(), TransactionType: "buy",
	// 	AssetType: "stock", Asset: "Apple", TickerSymbol: "AAPL",
	// 	Quantity: 10, Price: 100.5, Fees: 4, Currency: currency.EUR, Id: uuid.New()})
	// transactions = append(transactions, models.Transaction{Date: time.Now(), TransactionType: "buy",
	// 	AssetType: "stock", Asset: "Apple", TickerSymbol: "AAPL",
	// 	Quantity: 20, Price: 100.5, Fees: 4, Currency: currency.EUR, Id: uuid.New()})
	// transactions = append(transactions, models.Transaction{Date: time.Now(), TransactionType: "buy",
	// 	AssetType: "stock", Asset: "BASF", TickerSymbol: "BAS1",
	// 	Quantity: 100, Price: 45.5, Fees: 5, Currency: currency.EUR, Id: uuid.New()})
	// transactions = append(transactions, models.Transaction{Date: time.Now(), TransactionType: "sell",
	// 	AssetType: "stock", Asset: "Apple", TickerSymbol: "AAPL",
	// 	Quantity: 10, Price: 100.5, Fees: 4, Currency: currency.EUR, Id: uuid.New()})

	transactions, err := importer.LoadTransactions(filePath)
	if err != nil {
		return nil, err
	}

	unclosedTransactions := make(map[string][]models.Transaction)
	depotEntries := make(map[string]DepotEntry)
	realizedGains := make([]RealizedGain, 0, 5)

	for _, newTransaction := range transactions {
		fmt.Println(newTransaction.Asset)

		//Handle die Asset-Transaktionen
		//Add buy transaction to unclosed transactionsavailableBuytrans
		if newTransaction.TransactionType == "buy" {
			availableBuyTrans, exists := unclosedTransactions[newTransaction.Asset]
			//Gibt es schon Transaktionen für das Asset? Dann füge die neue Transaktion hinzu.
			if exists {
				availableBuyTrans = append(availableBuyTrans, newTransaction)
				unclosedTransactions[newTransaction.Asset] = availableBuyTrans
			} else {
				//Gibt es noch keine Transaktionen für das Asset? Dann erstelle einen neuen Slice mit der Transaktion.
				unclosedTransaction := []models.Transaction{}
				unclosedTransaction = append(unclosedTransaction, newTransaction)
				unclosedTransactions[newTransaction.Asset] = unclosedTransaction
			}
		}

		//Sell transaction
		if newTransaction.TransactionType == "sell" {
			transactions, exists := unclosedTransactions[newTransaction.Asset]
			if exists {
				//FiFo-Prinzip (First in, first out)
				//Ziehe die Anzahl der verkauften Assets von der ersten buy Transaktion ab
				//Sollten mehr Assets verkauft werden, als gekauft wurden, dann wird
				//die nächste buy Transaktion noch verwendet.
				//Code wurde erstellt
				for i, availableBuyTrans := range transactions {
					if availableBuyTrans.TransactionType == "buy" {
						//Buy und sell Transaktionen sind gleich
						if availableBuyTrans.Quantity == newTransaction.Quantity {
							//Entferne die buy Transaktion aus dem Slice
							transactions = slices.Delete(transactions, i, i+1)
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(newTransaction, availableBuyTrans))
							break
						}
						//Buy Transaktion ist größer als die Sell Transaktion
						if availableBuyTrans.Quantity > newTransaction.Quantity {
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(newTransaction, availableBuyTrans))
							//Buy Transaktion verkleinern um die Anzahl der verkauften Assets
							availableBuyTrans.Quantity -= newTransaction.Quantity
							transactions[i] = availableBuyTrans
							break
						}
						//Buy Transaktion ist kleiner als die Sell Transaktion
						if availableBuyTrans.Quantity < newTransaction.Quantity {
							newTransaction.Quantity -= availableBuyTrans.Quantity
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(newTransaction, availableBuyTrans))

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
					delete(unclosedTransactions, newTransaction.Asset)
				} else {
					//Aktualisiere die Transaktionen
					unclosedTransactions[newTransaction.Asset] = transactions
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
			entry, exists := depotEntries[transaction.Asset]
			if !exists {
				depotEntries[transaction.Asset] = DepotEntry{assetType: transaction.AssetType, asset: transaction.Asset,
					tickerSymbol: transaction.TickerSymbol, quantity: transaction.Quantity, price: transaction.Price,
					currency: transaction.Currency}
			} else {
				//Wenn das Asset schon im Depot ist, dann aktualisiere die Anzahl und den Preis
				entry.quantity += transaction.Quantity
				depotEntries[transaction.Asset] = entry
			}
		}
	}

	//fmt.Println(depot)
	fmt.Println(("Abgeschlossene Transaktionen"))
	fmt.Println(realizedGains)
	//fmt.Println(transactions[0].Date.Format("2006-01-02"))
	return depotEntries, nil
}
