package storage

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CsvStorage struct {
	filePath      string
	uuidGenerator func() uuid.UUID
}

func (s *CsvStorage) CreateDatabase() error {
	// Not implemented for CSV storage
	return errors.New("CreateDatabase not implemented for CSV storage")
}

func (s *CsvStorage) Ping() error {
	_, err := os.Open(s.filePath)
	return err
}

func (s *CsvStorage) AddTransaction(transaction *Transaction) error {
	//Not implemented for CSV storage
	return errors.New("AddTransaction not implemented for CSV storage")
}

func (s *CsvStorage) ReadAllTransactions() ([]Transaction, error) {
	lines, err := loadFile(s.filePath)
	if err != nil {
		return nil, err
	}
	transactions := make([]Transaction, 0, len(lines))
	for _, line := range lines {
		// Parse line
		values := strings.Split(line, ";")
		transaction := Transaction{}
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
		transaction.Quantity = float64(quantity)
		price, err := strconv.ParseFloat(values[6], 32)
		if err != nil {
			return nil, err
		}
		transaction.Price = float64(price)
		fees, err := strconv.ParseFloat(values[7], 32)
		if err != nil {
			return nil, err
		}
		transaction.Fees = float64(fees)
		// currency, err := currency.ParseISO(values[8])
		// if err != nil {
		// 	return nil, err
		// }
		transaction.Currency = values[8]
		transaction.Id = s.uuidGenerator()
		transaction.IsClosed = false
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *CsvStorage) AddUnclosedTransaction(asset Transaction) error {
	// Not implemented for CSV storage
	return errors.New("AddUnclosedTransaction not implemented for CSV storage")
}

func (s *CsvStorage) ReadAllUnclosedTickerSymbols() ([]string, error) {
	// Not implemented for CSV storage
	return nil, errors.New("ReadAllUnclosedTickerSymbols not implemented for CSV storage")
}

func (s *CsvStorage) ReadAllUnclosedTransactions() (map[string][]Transaction, error) {
	// Not implemented for CSV storage
	return nil, errors.New("ReadAllUnclosedTransactions not implemented for CSV storage")
}

func (s *CsvStorage) AddRealizedGain(realizedGain RealizedGain) error {
	// Not implemented for CSV storage
	return errors.New("AddRealizedGain not implemented for CSV storage")
}

func (s *CsvStorage) ReadAllRealizedGains() ([]RealizedGain, error) {
	// Not implemented for CSV storage
	return nil, errors.New("ReadAllRealizedGains not implemented for CSV storage")
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

func GetCsvStorage(pathToFile string, uuidGen func() uuid.UUID) CsvStorage {
	return CsvStorage{
		uuidGenerator: uuidGen,
		filePath:      pathToFile,
	}
}
