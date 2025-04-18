package depot

import (
	"fmt"

	"github.com/fritzrepo/stockportfolio/depot/importer"
	"github.com/fritzrepo/stockportfolio/models"
	"github.com/google/uuid"
)

type DepotEntry struct {
	AssetType    string  `json:"assetType"` //stock, crypto, forex
	Asset        string  `json:"asset"`     //Name des Assets
	TickerSymbol string  `json:"tickerSymbol"`
	Quantity     float32 `json:"quantity"` //Anzahl der Assets
	Price        float32 `json:"price"`    //Preis des Assets
	Currency     string  `json:"currency"` //Währung des Assets currency.Unit ist zu speziell für json und db.
}

// TotalPrice berechnet und gibt den gesamt Ankaufspreis zurück
func (d *DepotEntry) TotalPrice() float32 {
	return d.Quantity * d.Price
}

type RealizedGain struct {
	Id                uuid.UUID `json:"id"`                // ID der Realisierung
	SellTransactionId uuid.UUID `json:"sellTransactionId"` // ID der Verkaufstransaktion
	BuyTransactionId  uuid.UUID `json:"buytransactionId"`  // ID der Kauftransaktion
	Asset             string    // Asset-Name
	Amount            float32   // Der Gewinn/Verlust-Betrag
	IsProfit          bool      // true für Gewinn, false für Verlust
	TaxRate           float32   // Anwendbarer Steuersatz
	Quantity          float32
	BuyPrice          float32
	SellPrice         float32
	Currency          string
}

type Depot struct {
	DepotEntries         map[string]DepotEntry
	RealizedGains        []RealizedGain
	unclosedTransactions map[string][]models.Transaction
	uuidGenerator        func() uuid.UUID
}

func (d *Depot) ComputeTransactions(filePath string) error {

	transactions, err := importer.LoadTransactions(filePath, d.uuidGenerator)
	if err != nil {
		return err
	}

	unclosedTransactions := make(map[string][]models.Transaction)
	depotEntries := make(map[string]DepotEntry)
	realizedGains := make([]RealizedGain, 0, 5)

	for _, newTransaction := range transactions {

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
			continue
		}

		//Sell transaction
		if newTransaction.TransactionType == "sell" {
			transactions, exists := unclosedTransactions[newTransaction.Asset]
			if exists {
				//FiFo-Prinzip (First in, first out)
				//Ziehe die Anzahl der verkauften Assets von der ersten buy Transaktion ab.
				//Sollten mehr Assets verkauft werden, als gekauft wurden, dann wird
				//die nächste buy Transaktion verwendet.
				modifyTransactions := make([]models.Transaction, len(transactions))

				_ = copy(modifyTransactions, transactions)

				for _, availableBuyTrans := range transactions {
					if availableBuyTrans.TransactionType == "buy" {
						//Buy und sell Transaktionen sind gleich
						if availableBuyTrans.Quantity == newTransaction.Quantity {
							//Entferne die Transaktion aus der modifyTransactions
							filteredTransactions := []models.Transaction{}
							for _, transaction := range modifyTransactions {
								if transaction.Id != availableBuyTrans.Id {
									filteredTransactions = append(filteredTransactions, transaction)
								}
							}
							modifyTransactions = filteredTransactions

							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(d.uuidGenerator, newTransaction, availableBuyTrans))
							break
						}
						//Buy Transaktion ist größer als die Sell Transaktion
						if availableBuyTrans.Quantity > newTransaction.Quantity {
							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(d.uuidGenerator, newTransaction, availableBuyTrans))
							//Buy Transaktion verkleinern um die Anzahl der verkauften Assets
							availableBuyTrans.Quantity -= newTransaction.Quantity
							//Suche in den modifyTransactions die Transaktion
							for i, transaction := range modifyTransactions {
								if transaction.Id == availableBuyTrans.Id {
									modifyTransactions[i] = availableBuyTrans
									break
								}
							}
							break
						}
						//Buy Transaktion ist kleiner als die Sell Transaktion
						if availableBuyTrans.Quantity < newTransaction.Quantity {

							//Berechne den Gewinn / Verlust
							realizedGains = append(realizedGains, calculateProfitLoss(d.uuidGenerator, newTransaction, availableBuyTrans))
							newTransaction.Quantity -= availableBuyTrans.Quantity

							//Entferne die Transaktion aus der modifyTransactions
							filteredTransactions := []models.Transaction{}
							for _, transaction := range modifyTransactions {
								if transaction.Id != availableBuyTrans.Id {
									filteredTransactions = append(filteredTransactions, transaction)
								}
							}
							modifyTransactions = filteredTransactions
							//Der Rest der Sell transaction muss auf die nächste buy transaction angewendet werden
							//Daher kein break
						}
					} else {
						//ToDo => Richtige Fehlermeldung / Fehlerbehandlung implementieren
						fmt.Printf("Transaction is not a buy transaction %s\n", newTransaction.Asset)
					}
				}
				//Wennn die tansactions leer sind, dann lösche den Eintrag
				if len(modifyTransactions) == 0 {
					delete(unclosedTransactions, newTransaction.Asset)
				} else {
					//Aktualisiere die Transaktionen
					unclosedTransactions[newTransaction.Asset] = modifyTransactions
				}

				//Durchschnittskostenmethode (average cost)
				//ToDo => Implementieren
			} else {
				//ToDo => Richtige Fehlermeldung / Fehlerbehandlung implementieren
				fmt.Printf("No buy transaction available for this sell transaction %s\n", newTransaction.Asset)
			}
		}
	}

	//Handle depot
	for _, transactions := range unclosedTransactions {
		for _, transaction := range transactions {
			//Wenn das Asset noch nicht im Depot ist, dann füge es hinzu
			entry, exists := depotEntries[transaction.Asset]
			if !exists {
				depotEntries[transaction.Asset] = DepotEntry{AssetType: transaction.AssetType, Asset: transaction.Asset,
					TickerSymbol: transaction.TickerSymbol, Quantity: transaction.Quantity, Price: transaction.Price,
					Currency: transaction.Currency}
			} else {
				//Wenn das Asset schon im Depot ist, dann aktualisiere den (durchschnitts) Preis und die Anzahl
				entry.Price = (entry.Price*entry.Quantity + transaction.Price*transaction.Quantity) / (entry.Quantity + transaction.Quantity)
				entry.Quantity += transaction.Quantity
				depotEntries[transaction.Asset] = entry
			}
		}
	}

	d.unclosedTransactions = unclosedTransactions
	d.DepotEntries = depotEntries
	d.RealizedGains = realizedGains

	return nil
}

func NewDepot(uuidGen func() uuid.UUID) Depot {
	return Depot{
		DepotEntries:         make(map[string]DepotEntry),
		RealizedGains:        make([]RealizedGain, 0, 5),
		unclosedTransactions: make(map[string][]models.Transaction),
		uuidGenerator:        uuidGen,
	}
}
