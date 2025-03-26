package importer

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fritzrepo/stockportfolio/models"
	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

func LoadTransactions(filename string) ([]models.Transaction, error) {
	lines, err := loadFile(filename)
	if err != nil {
		return nil, err
	}
	transactions := make([]models.Transaction, 0, len(lines))
	for _, line := range lines {
		// Parse line
		values := strings.Split(line, ";")
		transaction := models.Transaction{}
		parsedDate, err := time.Parse("02.01.2006", values[0]) // Adjust the format as per your date format
		if err != nil {
			return nil, err
		}
		transaction.Date = parsedDate
		transaction.TransactionType = values[1]
		transaction.AssetType = values[2]
		transaction.Asset = values[3]
		transaction.TickerSymbol = values[4]
		quantity, err := strconv.ParseFloat(values[5], 32)
		if err != nil {
			return nil, err
		}
		transaction.Quantity = float32(quantity)
		price, err := strconv.ParseFloat(values[6], 32)
		if err != nil {
			return nil, err
		}
		transaction.Price = float32(price)
		fees, err := strconv.ParseFloat(values[7], 32)
		if err != nil {
			return nil, err
		}
		transaction.Fees = float32(fees)
		currency, err := currency.ParseISO(values[8])
		if err != nil {
			return nil, err
		}
		transaction.Currency = currency
		transaction.Id = uuid.New()
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func loadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
