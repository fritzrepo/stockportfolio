package portfolio

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

type Depot struct {
	depotEntries         map[string]DepotEntry
	unclosedTransactions map[string][]storage.Transaction
	store                storage.Store
}

func GetDepot(dataStore storage.Store) *Depot {
	return &Depot{
		depotEntries:         make(map[string]DepotEntry),
		unclosedTransactions: make(map[string][]storage.Transaction),
		store:                dataStore,
	}
}

func (d *Depot) CalculateSecuritiesAccountBalance() error {
	err := d.loadUnclosedTransactions()
	if err != nil {
		return err
	} else {
		d.createDepotEntries()
	}
	return nil
}

func (d *Depot) GetEntries() map[string]DepotEntry {
	return d.depotEntries
}

func (d *Depot) GetAllRealizedGains() ([]storage.RealizedGain, error) {
	realizedGains, err := d.store.ReadAllRealizedGains()
	if err != nil {
		return nil, fmt.Errorf("failed to read realized gains from store: %w", err)
	}
	return realizedGains, nil
}

// Diese Funktion berechnet die "Realized Gains" und die "unclosed transactions" aller ihr
// zugänglichen Transaktionen. Ist für den Import historischer Transaktionen oder für einen Fehlerfall gedacht,
// bei dem man alles neu berechnen muss.
// Sie ist auch für die Units Tests nützlich, da man damit den Algorithmus für "Realized Gains"
// und "unclosed transactions" gut testen kann.
func (d *Depot) ComputeAllTransactions() error {

	err := d.store.RemoveAllRealizedGains()
	if err != nil {
		return fmt.Errorf("failed to remove all realized gains from store: %w", err)
	}

	transactions, err := d.store.ReadAllTransactions()
	if err != nil {
		return err
	}

	for _, newTransaction := range transactions {
		//Die neue Transaktion kann auch mehrere Realized Gains erzeugen (bei FiFo-Prinzip)
		//Das ist der Fehler. Hier wird nur ein Realized Gain erzeugt, wenn die Transaktion verkauft wird.
		areNewRealizedGains, newRealizedGains, err := d.processNewTransaction(newTransaction)
		if err != nil {
			return err
		}

		if areNewRealizedGains {
			for _, newRealizedGain := range newRealizedGains {
				newRealizedGain.Id = uuid.New()
				err = d.store.AddRealizedGain(newRealizedGain)
				if err != nil {
					return fmt.Errorf("failed to add realized gain to store: %w", err)
				}
			}
		}
	}

	err = d.store.RemoveAllUnclosedTransactions()
	if err != nil {
		return fmt.Errorf("failed to remove all unclosed transaction from store: %w", err)
	}

	//Schleife zum Speichern aller unclosed transactions
	for _, asset := range d.unclosedTransactions {
		for _, transaction := range asset {
			err = d.store.AddUnclosedTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to save unclosed transactions to store: %w", err)
			}
		}
	}

	d.createDepotEntries()

	return nil
}

func (d *Depot) AddTransaction(newTransaction storage.Transaction) error {

	//Überprüfen, ob die Transaction schon existiert
	transaction, err := d.store.LoadTransactionByParams(newTransaction.Date, newTransaction.TransactionType, newTransaction.TickerSymbol)
	if err != nil {
		return err
	}
	if transaction != nil {
		return errors.New("transaction already exists")
	}

	newTransaction.Id = uuid.New()

	areNewRealizedGains, newRealizedGains, err := d.processNewTransaction(newTransaction)
	if err != nil {
		return err
	}

	err = d.store.AddTransaction(&newTransaction)
	if err != nil {
		return fmt.Errorf("failed to add transaction to store: %w", err)
	}

	if areNewRealizedGains {
		for _, newRealizedGain := range newRealizedGains {
			newRealizedGain.Id = uuid.New()
			err = d.store.AddRealizedGain(newRealizedGain)
			if err != nil {
				return fmt.Errorf("failed to add realized gain to store: %w", err)
			}
		}
	}

	err = d.store.RemoveAllUnclosedTransactions()
	if err != nil {
		return fmt.Errorf("failed to remove all unclosed transaction from store: %w", err)
	}

	//Schleife zum Speichern aller unclosed transactions
	for _, asset := range d.unclosedTransactions {
		for _, transaction := range asset {
			err = d.store.AddUnclosedTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to save unclosed transactions to store: %w", err)
			}
		}
	}

	d.createDepotEntries()
	return nil
}

func (d *Depot) processNewTransaction(newTransaction storage.Transaction) (bool, []storage.RealizedGain, error) {
	isNewRealizedGain := false
	var newRealizedGains []storage.RealizedGain
	var err error

	switch newTransaction.TransactionType {
	case "buy":
		d.addBuyTransaction(newTransaction)
	case "sell":
		isNewRealizedGain, newRealizedGains, err = d.addSellTransaction(newTransaction)
		if err != nil {
			return false, nil, fmt.Errorf("failed to process sell transaction: %w", err)
		}
	default:
		return false, nil, errors.New("transaction type not supported")
	}
	return isNewRealizedGain, newRealizedGains, nil
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

func (d *Depot) addSellTransaction(newTransaction storage.Transaction) (bool, []storage.RealizedGain, error) {

	areNewRealizedGains := false
	var newRealizedGains []storage.RealizedGain
	var newRealizedGain storage.RealizedGain

	transactions, exists := d.unclosedTransactions[newTransaction.TickerSymbol]
	if !exists {
		return areNewRealizedGains, nil, fmt.Errorf("no buy transaction available for this sell transaction %s", newTransaction.TickerSymbol)
	}

	//FiFo-Prinzip (First in, first out)

	//Ziehe die Anzahl der verkauften Assets von der ersten buy Transaktion ab.
	//Sollten mehr Assets verkauft werden, als gekauft wurden, dann wird
	//die nächste buy Transaktion verwendet.

	modifyTransactions := make([]storage.Transaction, len(transactions))

	_ = copy(modifyTransactions, transactions)

	for _, availableBuyTrans := range transactions {

		if availableBuyTrans.TransactionType != "buy" {
			return areNewRealizedGains, nil, fmt.Errorf("transaction is not a buy transaction %s", newTransaction.TickerSymbol)
		}

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
			areNewRealizedGains = true
			newRealizedGain = calculateProfitLoss(newTransaction, availableBuyTrans)
			newRealizedGains = append(newRealizedGains, newRealizedGain)
			break
		}
		//Buy Transaktion ist größer als die Sell Transaktion
		if availableBuyTrans.Quantity > newTransaction.Quantity {
			//Berechne den Gewinn / Verlust
			areNewRealizedGains = true
			newRealizedGain = calculateProfitLoss(newTransaction, availableBuyTrans)
			newRealizedGains = append(newRealizedGains, newRealizedGain)
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
			areNewRealizedGains = true
			newRealizedGain = calculateProfitLoss(newTransaction, availableBuyTrans)
			newRealizedGains = append(newRealizedGains, newRealizedGain)
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
	}

	//Wennn die tansactions leer sind, dann lösche den Eintrag
	if len(modifyTransactions) == 0 {
		delete(d.unclosedTransactions, newTransaction.TickerSymbol)
	} else {
		//Aktualisiere die Transaktionen
		d.unclosedTransactions[newTransaction.TickerSymbol] = modifyTransactions
	}

	//Durchschnittskostenmethode (average cost)
	//ToDo => Implementieren wenn nötig.

	return areNewRealizedGains, newRealizedGains, nil
}

func (d *Depot) loadUnclosedTransactions() error {
	var err error
	d.unclosedTransactions, err = d.store.ReadAllUnclosedTransactions()
	if err != nil {
		return fmt.Errorf("failed to read unclosed transactions from store: %w", err)
	}
	return nil
}

func (d *Depot) createDepotEntries() {
	clear(d.depotEntries)
	for _, transactions := range d.unclosedTransactions {
		for _, transaction := range transactions {
			//Wenn das Asset noch nicht im Depot ist, dann füge es hinzu
			entry, exists := d.depotEntries[transaction.TickerSymbol]
			if !exists {
				d.depotEntries[transaction.TickerSymbol] = DepotEntry{AssetType: transaction.AssetType, Asset: transaction.Asset,
					TickerSymbol: transaction.TickerSymbol, Quantity: transaction.Quantity, Price: transaction.Price,
					Currency: transaction.Currency}
			} else {
				//Wenn das Asset schon im Depot ist, dann aktualisiere den (durchschnitts) Preis und die Anzahl
				entry.Price = (entry.Price*entry.Quantity + transaction.Price*transaction.Quantity) / (entry.Quantity + transaction.Quantity)
				entry.Quantity += transaction.Quantity
				d.depotEntries[transaction.TickerSymbol] = entry
			}
		}
	}
}
