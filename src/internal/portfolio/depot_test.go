package portfolio

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/storage"
)

func setupTestStore(t *testing.T) storage.Store {
	store := storage.GetMemoryDatabase()
	store.Open()
	err := store.CreateDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	// Cleanup registrieren. Wird nach jedem Test ausgeführt.
	t.Cleanup(func() {
		store.Close()
	})
	return store
}

// TestComputeTransactions
func TestComputeTransactions(t *testing.T) {

	type TestCases = []struct {
		Name          string                 `json:"name"`
		ExpectedGains []storage.RealizedGain `json:"expectedGains"`
		ExpectedDepot map[string]DepotEntry  `json:"expectedDepot"`
	}

	testcount := 5
	testCases := make(TestCases, testcount)
	for i := range testcount {

		filenameDepot := fmt.Sprintf("../../testdata/depot/expectedDepot%d.json", i+1)
		jsonFileDepot, err := os.Open(filenameDepot)
		if err != nil {
			t.Fatalf("Failed to open test data: %v", err)
		}
		defer jsonFileDepot.Close()

		filenameGains := fmt.Sprintf("../../testdata/depot/expectedGains%d.json", i+1)
		jsonFileGains, err := os.Open(filenameGains)
		if err != nil {
			t.Fatalf("Failed to open test data: %v", err)
		}
		defer jsonFileGains.Close()

		testCases[i].Name = fmt.Sprintf("Test%d", i+1)

		byteValueDepot, _ := io.ReadAll(jsonFileDepot)
		if err := json.Unmarshal(byteValueDepot, &testCases[i].ExpectedDepot); err != nil {
			t.Fatalf("Failed to unmarshal test data: %v", err)
		}

		byteValueGains, _ := io.ReadAll(jsonFileGains)
		if err := json.Unmarshal(byteValueGains, &testCases[i].ExpectedGains); err != nil {
			t.Fatalf("Failed to unmarshal test data: %v", err)
		}
	}

	i := 1
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {

			filenameTrans := fmt.Sprintf("../../testdata/depot/RawTransactionsTest%d.csv", i)
			i = i + 1

			store := storage.GetCsvStorage(filenameTrans)
			dep := GetDepot(&store)

			err := dep.ComputeAllTransactions()
			if err != nil {
				t.Fatalf("Error computing transactions: %v", err)
			}

			//Check the depot entries
			result := dep.depotEntries

			for key, expectedEntry := range tt.ExpectedDepot {
				resultEntry, exists := result[key]
				if !exists {
					t.Errorf("Expected asset %s not found in result", key)
					continue
				}

				if !reflect.DeepEqual(resultEntry, expectedEntry) {
					t.Errorf("For asset %s, expected %+v, but got %+v", key, expectedEntry, resultEntry)
				}
			}

			//Check the realized gains
			realizedGains, err := dep.GetAllRealizedGains()
			if err != nil {
				t.Fatalf("Error getting realized gains: %v", err)
			}

			for idx, expectedEntry := range tt.ExpectedGains {

				// Suche den entsprechenden RealizedGain-Eintrag
				var matchedGain *storage.RealizedGain
				if idx < len(realizedGains) {
					matchedGain = &realizedGains[idx]
				}

				if matchedGain == nil {
					t.Errorf("No matching realized gain found for SellTransactionId: %s and BuyTransactionId: %s", expectedEntry.SellTransactionId, expectedEntry.BuyTransactionId)
					continue
				}

				const epsilon = 1e-3
				if matchedGain.Asset != expectedEntry.Asset ||
					math.Abs(expectedEntry.Amount-matchedGain.Amount) > epsilon ||
					matchedGain.IsProfit != expectedEntry.IsProfit ||
					math.Abs(expectedEntry.TaxRate-matchedGain.TaxRate) > epsilon ||
					math.Abs(expectedEntry.Quantity-matchedGain.Quantity) > epsilon ||
					math.Abs(expectedEntry.BuyPrice-matchedGain.BuyPrice) > epsilon ||
					math.Abs(expectedEntry.SellPrice-matchedGain.SellPrice) > epsilon ||
					matchedGain.Currency != expectedEntry.Currency {
					t.Errorf("Realized gain for asset %s does not match expected values. Expected: %+v, Got: %+v", expectedEntry.Asset, expectedEntry, matchedGain)
				}
			}
		})
	}
}

// Test addTransaction
func TestAddTransactions(t *testing.T) {
	store := setupTestStore(t)
	dep := GetDepot(store)

	err := dep.AddTransaction(storage.Transaction{
		Date:            time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		TransactionType: "buy",
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"})

	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	err = dep.AddTransaction(storage.Transaction{
		Date:            time.Date(2023, 11, 1, 14, 0, 0, 0, time.UTC),
		TransactionType: "sell",
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        5,
		Price:           200,
		Fees:            1.5,
		Currency:        "USD"})

	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	realizedGains, _ := dep.GetAllRealizedGains()

	if len(realizedGains) != 1 {
		t.Errorf("Expected 1 realized gain, but got %d", len(realizedGains))
	} else {
		gain := realizedGains[0]
		if gain.Asset != "Apple" {
			t.Errorf("Realized gain values do not match expected values: %+v", gain)
		}
	}

	if len(dep.depotEntries) != 1 {
		t.Errorf("Expected depot to have 1 entry after selling all shares, but got %d entries", len(dep.depotEntries))
	}

}

func TestDoNotAddAnExistingTransaction(t *testing.T) {
	store := setupTestStore(t)
	dep := GetDepot(store)

	transaction := storage.Transaction{
		Date:            time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		TransactionType: "buy",
		AssetType:       "stock",
		Asset:           "Apple",
		TickerSymbol:    "AAPL",
		Quantity:        10,
		Price:           150,
		Fees:            1.5,
		Currency:        "USD"}

	err := dep.AddTransaction(transaction)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	//Versuche die gleiche Transaktion erneut hinzuzufügen
	err = dep.AddTransaction(transaction)
	if err == nil {
		t.Fatalf("Expected error when adding an existing transaction, but got none")
	}
}
