package depot

import (
	"fmt"
	"time"

	"golang.org/x/text/currency"
)

type transaction struct {
	date            time.Time
	transactionType string //besser enum buy, sell
	assetType       string //stock, crypto, forex
	asset           string
	tickerSymbol    string
	quantity        float32
	price           float32
	fees            float32
	totalPrice      float32
	currency        currency.Unit
}

type depotEntry struct {
	assetType    string
	asset        string
	tickerSymbol string
	quantity     float32
	price        float32
	totalPrice   float32
	currency     currency.Unit
}

func calculateProfitLoss() {
	//ToDo => Implementieren
	//FiFo-Prinzip (First in, first out)

	//Durchschnittskostenmethode (average cost)
}

func ComputeTransactions() map[string]depotEntry {
	transactions := []transaction{} // Ein Slice für Transaktionen
	transactions = append(transactions, transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 10, price: 100.5, fees: 4, totalPrice: 1009, currency: currency.EUR})
	transactions = append(transactions, transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "Apple", tickerSymbol: "AAPL",
		quantity: 20, price: 100.5, fees: 4, totalPrice: 2014, currency: currency.EUR})
	transactions = append(transactions, transaction{date: time.Now(), transactionType: "buy",
		assetType: "stock", asset: "BASF", tickerSymbol: "BAS1",
		quantity: 100, price: 45.5, fees: 5, totalPrice: 4555, currency: currency.EUR})

	//depot := []depotEntry{} // Ein Slice für Depot Einträge

	depotMap := make(map[string]depotEntry)

	for _, value := range transactions {
		fmt.Println(value.asset)

		//Wenn das Asset noch nicht im Depot ist, dann füge es hinzu
		// if !contains(depot, value.asset) {

		// 	depot = append(depot, depotEntry{assetType: value.assetType, asset: value.asset,
		// 		tickerSymbol: value.tickerSymbol, quantity: value.quantity, price: value.price,
		// 		totalPrice: value.totalPrice, currency: value.currency})

		// 	continue
		// }

		//Wenn das Asset schon im Depot ist, dann aktualisiere die Anzahl und den Preis
		// for i, v := range depot {
		// 	if v.asset == value.asset {
		// 		depot[i].quantity += value.quantity
		// 		depot[i].totalPrice += value.totalPrice
		// 	}
		// }

		//Buy transaction
		if value.transactionType == "buy" {
			entry, exists := depotMap[value.asset]
			if !exists {
				//if depotMap[value.asset] == (depotEntry{}) {
				depotMap[value.asset] = depotEntry{assetType: value.assetType, asset: value.asset,
					tickerSymbol: value.tickerSymbol, quantity: value.quantity, price: value.price,
					totalPrice: value.totalPrice, currency: value.currency}
			} else {
				//entry = depotMap[value.asset]
				entry.quantity += value.quantity
				entry.totalPrice += value.totalPrice
				depotMap[value.asset] = entry
			}
		}
		//Sell transaction
		if value.transactionType == "sell" {
			entry, exists := depotMap[value.asset]
			if exists {
				entry.quantity -= value.quantity
				entry.totalPrice -= value.totalPrice
				depotMap[value.asset] = entry
				//Stelle Gewinn/Verlust fest
				//ToDo => Implementieren

				//Wenn die Anzahl 0 ist, lösche den Eintrag aus dem Depot
				if entry.quantity == 0 {
					delete(depotMap, value.asset)
				}
			} else {
				//ToDo => Richtige Fehlermeldung implementieren
				fmt.Println("Asset not found in	depot")
			}
		}
	}

	//fmt.Println(depot)
	fmt.Println(depotMap)
	fmt.Println(transactions[0].date.Format("2006-01-02"))
	return depotMap
}
