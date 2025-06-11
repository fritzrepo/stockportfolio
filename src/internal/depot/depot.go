package depot

import (
	"errors"
	"fmt"

	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/google/uuid"
)

type DepotEntry struct {
	AssetType    string  `json:"assetType"` //stock, crypto, forex
	Asset        string  `json:"asset"`     //Name des Assets
	TickerSymbol string  `json:"tickerSymbol"`
	Quantity     float64 `json:"quantity"` //Anzahl der Assets
	Price        float64 `json:"price"`    //Preis des Assets
	Currency     string  `json:"currency"` //Währung des Assets currency. Unit ist zu speziell für json und db.
}

// TotalPrice berechnet und gibt den gesamt Ankaufspreis zurück
func (d *DepotEntry) TotalPrice() float64 {
	return d.Quantity * d.Price
}

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

type Depot struct {
	DepotEntries         map[string]DepotEntry
	RealizedGains        []RealizedGain
	unclosedTransactions map[string][]storage.Transaction
	uuidGenerator        func() uuid.UUID
	store                storage.Store
}

// Diese Funktion berechnet die "Realized Gains" und die "unclosed transactions" aller ihr
// zugänglichen Transaktionen. Ist für den Import historischer Transaktionen oder für einen Fehlerfall gedacht,
// bei dem man alles neu berechnen muss.
// Sie ist auch für die Units Tests nützlich, da man damit den Algorithmus für "Realized Gains"
// und "unclosed transactions" gut der testen kann.
func (d *Depot) ComputeAllTransactions() error {

	//ToDo => Prüfen ob hier noch die persistierten "RealizedGains" und
	// die "unclosed transactions" gelöscht werden müsssen.

	transactions, err := d.store.LoadAllTransactions()
	if err != nil {
		return err
	}

	for _, newTransaction := range transactions {
		d.AddTransaction(newTransaction)
	}

	return nil
}

// func saveUnclosedTransactions() {
// 	//ToDo => Implementieren
// }

func loadUnclosedTransactions() {
	//ToDo => Implementieren
}

// func saveRealizedGains() {
// 	//ToDo => Implementieren
// }

// func loadRealizedGains() {
// 	//ToDo => Implementieren
// }

func (d *Depot) AddTransaction(newTransaction storage.Transaction) error {
	if newTransaction.TransactionType == "buy" {
		d.addBuyTransaction(newTransaction)
	} else if newTransaction.TransactionType == "sell" {
		d.addSellTransaction(newTransaction)
	} else {
		return errors.New("transaction type not supported")
	}

	//saveRealizedGains()
	d.createDepotEntries()
	return nil
}

func (d *Depot) addBuyTransaction(newTransaction storage.Transaction) {

	availableBuyTrans, exists := d.unclosedTransactions[newTransaction.TickerSymbol]
	//Gibt es schon Transaktionen für das Asset? Dann füge die neue Transaktion hinzu.
	if exists {
		availableBuyTrans = append(availableBuyTrans, newTransaction)
		d.unclosedTransactions[newTransaction.TickerSymbol] = availableBuyTrans
	} else {
		//Gibt es noch keine Transaktionen für das Asset? Dann erstelle einen neuen Slice mit der Transaktion.
		unclosedTransaction := []storage.Transaction{}
		unclosedTransaction = append(unclosedTransaction, newTransaction)
		d.unclosedTransactions[newTransaction.TickerSymbol] = unclosedTransaction
	}
}

func (d *Depot) addSellTransaction(newTransaction storage.Transaction) {

	transactions, exists := d.unclosedTransactions[newTransaction.TickerSymbol]
	if exists {
		//FiFo-Prinzip (First in, first out)
		//Ziehe die Anzahl der verkauften Assets von der ersten buy Transaktion ab.
		//Sollten mehr Assets verkauft werden, als gekauft wurden, dann wird
		//die nächste buy Transaktion verwendet.
		modifyTransactions := make([]storage.Transaction, len(transactions))

		_ = copy(modifyTransactions, transactions)

		for _, availableBuyTrans := range transactions {
			if availableBuyTrans.TransactionType == "buy" {
				//Buy und sell Transaktionen sind gleich
				if availableBuyTrans.Quantity == newTransaction.Quantity {
					//Entferne die Transaktion aus der modifyTransactions
					filteredTransactions := []storage.Transaction{}
					for _, transaction := range modifyTransactions {
						if transaction.Id != availableBuyTrans.Id {
							filteredTransactions = append(filteredTransactions, transaction)
						}
					}
					modifyTransactions = filteredTransactions

					//Berechne den Gewinn / Verlust
					d.RealizedGains = append(d.RealizedGains, calculateProfitLoss(d.uuidGenerator, newTransaction, availableBuyTrans))
					break
				}
				//Buy Transaktion ist größer als die Sell Transaktion
				if availableBuyTrans.Quantity > newTransaction.Quantity {
					//Berechne den Gewinn / Verlust
					d.RealizedGains = append(d.RealizedGains, calculateProfitLoss(d.uuidGenerator, newTransaction, availableBuyTrans))
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
					d.RealizedGains = append(d.RealizedGains, calculateProfitLoss(d.uuidGenerator, newTransaction, availableBuyTrans))
					newTransaction.Quantity -= availableBuyTrans.Quantity

					//Entferne die Transaktion aus der modifyTransactions
					filteredTransactions := []storage.Transaction{}
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
				fmt.Printf("Transaction is not a buy transaction %s\n", newTransaction.TickerSymbol)
			}
		}
		//Wennn die tansactions leer sind, dann lösche den Eintrag
		if len(modifyTransactions) == 0 {
			delete(d.unclosedTransactions, newTransaction.TickerSymbol)
		} else {
			//Aktualisiere die Transaktionen
			d.unclosedTransactions[newTransaction.TickerSymbol] = modifyTransactions
		}

		//Durchschnittskostenmethode (average cost)
		//ToDo => Implementieren
	} else {
		//ToDo => Richtige Fehlermeldung / Fehlerbehandlung implementieren
		fmt.Printf("No buy transaction available for this sell transaction %s\n", newTransaction.TickerSymbol)
	}
}

func (d *Depot) createDepotEntries() {
	clear(d.DepotEntries)
	for _, transactions := range d.unclosedTransactions {
		for _, transaction := range transactions {
			//Wenn das Asset noch nicht im Depot ist, dann füge es hinzu
			entry, exists := d.DepotEntries[transaction.TickerSymbol]
			if !exists {
				d.DepotEntries[transaction.TickerSymbol] = DepotEntry{AssetType: transaction.AssetType, Asset: transaction.Asset,
					TickerSymbol: transaction.TickerSymbol, Quantity: transaction.Quantity, Price: transaction.Price,
					Currency: transaction.Currency}
			} else {
				//Wenn das Asset schon im Depot ist, dann aktualisiere den (durchschnitts) Preis und die Anzahl
				entry.Price = (entry.Price*entry.Quantity + transaction.Price*transaction.Quantity) / (entry.Quantity + transaction.Quantity)
				entry.Quantity += transaction.Quantity
				d.DepotEntries[transaction.TickerSymbol] = entry
			}
		}
	}
}

func (d *Depot) LoadStock() {
	loadUnclosedTransactions()
	d.createDepotEntries()
}

func GetDepot(uuidGen func() uuid.UUID, dataStore storage.Store) Depot {
	return Depot{
		DepotEntries:         make(map[string]DepotEntry),
		RealizedGains:        make([]RealizedGain, 0, 5),
		unclosedTransactions: make(map[string][]storage.Transaction),
		uuidGenerator:        uuidGen,
		store:                dataStore,
	}
}
