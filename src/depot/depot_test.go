package depot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// Mock-Generator für Tests
type MockUUIDGenerator struct {
	uuids     []uuid.UUID
	callCount int
}

// Erstellt einen neuen Mock-Generator mit 20 vordefinierten UUIDs
func NewMockUUIDGenerator() *MockUUIDGenerator {
	mockGenerator := &MockUUIDGenerator{
		uuids:     make([]uuid.UUID, 20),
		callCount: 0,
	}

	// 20 feste UUIDs generieren
	for i := 0; i < 20; i++ {
		// Erzeuge vorhersehbare UUIDs, die sich nur in einer Ziffer unterscheiden
		mockGenerator.uuids[i] = uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-0000000000%02d", i+1))
	}

	return mockGenerator
}

// GetUUID liefert bei jedem Aufruf die nächste UUID aus der Liste
func (m *MockUUIDGenerator) GetUUID() uuid.UUID {
	if m.callCount >= len(m.uuids) {
		// Wenn alle UUIDs verwendet wurden, beginne von vorne
		m.callCount = 0
	}

	id := m.uuids[m.callCount]
	m.callCount++
	return id
}

// TestComputeTransactions

// One realized gain and one stock in depot
func Test1ComputeTransactions(t *testing.T) {

	expectedDepot := map[string]DepotEntry{
		"BASF": {
			AssetType:    "stock",
			Asset:        "BASF",
			TickerSymbol: "BAS1",
			Quantity:     100,
			Price:        45.5,
			Currency:     "EUR",
		},
	}

	expectedGains := []RealizedGain{
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			BuyTransactionId:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Asset:             "Apple",
			Amount:            95,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          10,
			BuyPrice:          100.5,
			SellPrice:         110,
			Currency:          "EUR",
		},
	}

	var uuidGenerator = NewMockUUIDGenerator()
	dep := NewDepot(uuidGenerator.GetUUID)

	err := dep.ComputeTransactions("./test_data/RawTransactionsTest1.csv")
	if err != nil {
		t.Fatalf("Error computing transactions: %v", err)
	}
	depotEntry := dep.DepotEntries
	realizedGains := dep.RealizedGains

	// Check the depot entries
	for key, expectedEntry := range expectedDepot {
		resultEntry, exists := depotEntry[key]
		if !exists {
			t.Errorf("Expected asset %s not found in result", key)
			continue
		}

		if !reflect.DeepEqual(resultEntry, expectedEntry) {
			t.Errorf("For asset %s, expected %+v, but got %+v", key, expectedEntry, resultEntry)
		}
	}

	// Check count of realized gains
	if len(realizedGains) != len(expectedGains) {
		t.Errorf("Slice lengths differ: expected %d, got %d", len(expectedGains), len(realizedGains))
		return
	}

	// Check each element in the realized gains
	for i, exp := range expectedGains {
		if !reflect.DeepEqual(realizedGains[i], exp) {
			t.Errorf("Element at index %d differs:\nExpected: %+v\nGot: %+v", i, exp, realizedGains[i])
		}
	}

}

//Checken, ob wenn ein Asset zwei mal gekauft wird, der Durchschnitts-Kaufpreis richtig berechnet wird

func Test2ComputeTransactions(t *testing.T) {

	expectedDepot := map[string]DepotEntry{
		"Apple": {
			AssetType:    "stock",
			Asset:        "Apple",
			TickerSymbol: "AAPL",
			Quantity:     15,
			Price:        110.5,
			Currency:     "EUR",
		},
		"BASF": {
			AssetType:    "stock",
			Asset:        "BASF",
			TickerSymbol: "BAS1",
			Quantity:     100,
			Price:        45.5,
			Currency:     "EUR",
		},
	}

	expected := []RealizedGain{
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000007"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			BuyTransactionId:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Asset:             "Siemens",
			Amount:            550,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          100,
			BuyPrice:          90,
			SellPrice:         95.5,
			Currency:          "EUR",
		},
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000008"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			BuyTransactionId:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Asset:             "Apple",
			Amount:            199,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          10,
			BuyPrice:          100.5,
			SellPrice:         120.4,
			Currency:          "EUR",
		},
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			BuyTransactionId:  uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Asset:             "Apple",
			Amount:            49.5,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          5,
			BuyPrice:          110.5,
			SellPrice:         120.4,
			Currency:          "EUR",
		},
	}

	var uuidGenerator = NewMockUUIDGenerator()
	dep := NewDepot(uuidGenerator.GetUUID)

	err := dep.ComputeTransactions("./test_data/RawTransactionsTest2.csv")
	if err != nil {
		t.Fatalf("Error computing transactions: %v", err)
	}

	//Check the depot entries
	result := dep.DepotEntries

	for key, expectedEntry := range expectedDepot {
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
	realizedGains := dep.RealizedGains

	for _, expectedEntry := range expected {
		for _, realizedGain := range realizedGains {
			if realizedGain.Id == expectedEntry.Id {
				if reflect.DeepEqual(realizedGain, expectedEntry) {
					t.Logf("Realized gain for asset %s matches expected values: %+v", expectedEntry.Asset, realizedGain)
				}
			}
		}
	}
}

func TestComputeTransactions(t *testing.T) {

	type TestCases = []struct {
		Name          string                `json:"name"`
		ExpectedGains []RealizedGain        `json:"expectedGains"`
		ExpectedDepot map[string]DepotEntry `json:"expectedDepot"`
	}

	testcount := 1
	testCases := make(TestCases, testcount)

	for i := range testcount {

		filenameDepot := fmt.Sprintf("./test_data/expectedDepot%d.json", i+1)
		jsonFileDepot, err := os.Open(filenameDepot)
		if err != nil {
			t.Fatalf("Failed to open test data: %v", err)
		}
		defer jsonFileDepot.Close()

		filenameGains := fmt.Sprintf("./test_data/expectedGains%d.json", i+1)
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

	fmt.Println(testCases[0].ExpectedDepot)
	fmt.Println(testCases[0].ExpectedGains)

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {

			var uuidGenerator = NewMockUUIDGenerator()
			dep := NewDepot(uuidGenerator.GetUUID)

			err := dep.ComputeTransactions("./test_data/RawTransactionsTest2.csv")
			if err != nil {
				t.Fatalf("Error computing transactions: %v", err)
			}

			//Check the depot entries
			result := dep.DepotEntries

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
			realizedGains := dep.RealizedGains

			for _, expectedEntry := range tt.ExpectedGains {
				for _, realizedGain := range realizedGains {
					if realizedGain.Id == expectedEntry.Id {
						if reflect.DeepEqual(realizedGain, expectedEntry) {
							t.Logf("Realized gain for asset %s matches expected values: %+v", expectedEntry.Asset, realizedGain)
						}
					}
				}
			}

		})
	}

}
