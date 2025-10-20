package imports

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/storage"
	"github.com/google/uuid"
)

func ReadRealizedGainsCsv(fileName string) []storage.Transaction {

	var transactions []storage.Transaction

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Setze das Trennzeichen auf Semikolon
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {

		// Erstelle buy transaction
		newTransaction := storage.Transaction{}
		newTransaction.Id = uuid.New()
		newTransaction.Asset = record[0]
		newTransaction.Date, err = time.Parse("02.01.06", record[7])
		if err != nil {
			panic(err)
		}
		newTransaction.TransactionType = "buy"
		newTransaction.AssetType = record[1]
		newTransaction.TickerSymbol = ""
		newTransaction.Quantity, _ = strconv.ParseFloat(record[2], 64)
		str := record[3]
		str = strings.ReplaceAll(str, ".", "")  // Punkte entfernen
		str = strings.Replace(str, ",", ".", 1) // Komma zu Punkt
		str = strings.ReplaceAll(str, "EUR", "")
		str = strings.TrimSpace(str)
		f, err := strconv.ParseFloat(str, 64)
		newTransaction.Price = f
		newTransaction.Fees = 0.0
		newTransaction.Currency = "EUR"

		if err != nil {
			panic(err)
		}

		transactions = append(transactions, newTransaction)

		// Erstelle sell transaction
		newTransaction = storage.Transaction{}
		newTransaction.Id = uuid.New()
		newTransaction.Asset = record[0]
		newTransaction.Date, err = time.Parse("02.01.06", record[8])
		if err != nil {
			panic(err)
		}
		newTransaction.TransactionType = "sell"
		newTransaction.AssetType = record[1]
		newTransaction.TickerSymbol = ""
		newTransaction.Quantity, _ = strconv.ParseFloat(record[2], 64)
		str2 := record[4]
		str2 = strings.ReplaceAll(str2, ".", "")  // Punkte entfernen
		str2 = strings.Replace(str2, ",", ".", 1) // Komma zu Punkt
		str2 = strings.ReplaceAll(str2, "EUR", "")
		str2 = strings.TrimSpace(str2)
		f2, err := strconv.ParseFloat(str2, 64)
		newTransaction.Price = f2
		newTransaction.Fees = 0.0
		newTransaction.Currency = "EUR"

		if err != nil {
			panic(err)
		}

		transactions = append(transactions, newTransaction)
	}

	return transactions
}
