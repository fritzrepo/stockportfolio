package depot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

type DepotEntryForTest struct {
	assetType    string
	asset        string
	tickerSymbol string
	quantity     float32
	price        float32
	currency     currency.Unit
	totalPrice   float32
}

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

	// expectedDepot := map[string]DepotEntryForTest{
	// 	"BASF": {
	// 		assetType:    "stock",
	// 		asset:        "BASF",
	// 		tickerSymbol: "BAS1",
	// 		quantity:     100,
	// 		price:        45.5,
	// 		totalPrice:   4550,
	// 		currency:     currency.EUR,
	// 	},
	// }
	expectedDepot := map[string]DepotEntry{
		"BASF": {
			assetType:    "stock",
			asset:        "BASF",
			tickerSymbol: "BAS1",
			quantity:     100,
			price:        45.5,
			currency:     currency.EUR,
		},
	}

	expectedGains := []RealizedGain{
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			BuyTransactionsId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Asset:             "Apple",
			Amount:            95,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          10,
			BuyPrice:          100.5,
			SellPrice:         110,
		},
	}

	var uuidGenerator = NewMockUUIDGenerator()
	dep := NewDepot(uuidGenerator.GetUUID)

	err := dep.ComputeTransactions("../data/RawTransactionsTest1.csv")
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
		// if resultEntry.assetType != expectedEntry.assetType ||
		// 	resultEntry.asset != expectedEntry.asset ||
		// 	resultEntry.tickerSymbol != expectedEntry.tickerSymbol ||
		// 	resultEntry.quantity != expectedEntry.quantity ||
		// 	resultEntry.price != expectedEntry.price ||
		// 	resultEntry.currency != expectedEntry.currency ||
		// 	resultEntry.TotalPrice() != expectedEntry.totalPrice {
		// 	t.Errorf("For asset %s, expected %+v, but got %+v", key, expectedEntry, resultEntry)
		// }

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

	// Also jedes Feld einzeln prüfen
	// for i, exp := range expectedGains {
	// 	res := realizedGains[i]
	// 	if exp.Amount != res.Amount || exp.Asset != res.Asset ||
	// 		exp.BuyPrice != res.BuyPrice || exp.IsProfit != res.IsProfit ||
	// 		exp.Quantity != res.Quantity || exp.SellPrice != res.SellPrice ||
	// 		exp.TaxRate != res.TaxRate {
	// 		t.Errorf("Element at index %d differs:\nExpected: %+v\nGot: %+v", i, exp, res)
	// 	}
	// }

}

//Checken, ob wenn ein Asset zwei mal gekauft wird, der Durchschnitts-Kaufpreis richtig berechnet wird

func Test2ComputeTransactions(t *testing.T) {

	expectedDepot := map[string]DepotEntryForTest{
		"Apple": {
			assetType:    "stock",
			asset:        "Apple",
			tickerSymbol: "AAPL",
			quantity:     15,
			price:        110.5,
			totalPrice:   1657.5, //Hier fester Wert, im Originalcode wird die totalPrice berechnet
			currency:     currency.EUR,
		},
		"BASF": {
			assetType:    "stock",
			asset:        "BASF",
			tickerSymbol: "BAS1",
			quantity:     100,
			price:        45.5,
			totalPrice:   4550,
			currency:     currency.EUR,
		},
	}

	var uuidGenerator = NewMockUUIDGenerator()
	dep := NewDepot(uuidGenerator.GetUUID)

	err := dep.ComputeTransactions("../data/RawTransactions.csv")
	if err != nil {
		t.Fatalf("Error computing transactions: %v", err)
	}
	result := dep.DepotEntries

	for key, expectedEntry := range expectedDepot {
		resultEntry, exists := result[key]
		if !exists {
			t.Errorf("Expected asset %s not found in result", key)
			continue
		}

		if resultEntry.assetType != expectedEntry.assetType ||
			resultEntry.asset != expectedEntry.asset ||
			resultEntry.tickerSymbol != expectedEntry.tickerSymbol ||
			resultEntry.quantity != expectedEntry.quantity ||
			resultEntry.price != expectedEntry.price ||
			resultEntry.currency != expectedEntry.currency ||
			resultEntry.TotalPrice() != expectedEntry.totalPrice {
			t.Errorf("For asset %s, expected %+v, but got %+v", key, expectedEntry, resultEntry)
		}
	}
}

func TestComputeTransactionsRealizedGains(t *testing.T) {
	expected := []RealizedGain{
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000007"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			BuyTransactionsId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Asset:             "Siemens",
			Amount:            550,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          100,
			BuyPrice:          90,
			SellPrice:         95.5,
		},
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000008"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			BuyTransactionsId: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Asset:             "Apple",
			Amount:            199,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          10,
			BuyPrice:          100.5,
			SellPrice:         120.4,
		},
		{
			Id:                uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			SellTransactionId: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			BuyTransactionsId: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Asset:             "Apple",
			Amount:            49.5,
			IsProfit:          true,
			TaxRate:           0.25,
			Quantity:          5,
			BuyPrice:          110.5,
			SellPrice:         120.4,
		},
	}

	var uuidGenerator = NewMockUUIDGenerator()
	dep := NewDepot(uuidGenerator.GetUUID)

	err := dep.ComputeTransactions("../data/RawTransactions.csv")
	if err != nil {
		t.Fatalf("Error computing transactions: %v", err)
	}
	realizedGains := dep.RealizedGains

	for _, expectedEntry := range expected {
		for _, realizedGain := range realizedGains {
			if realizedGain.Id == expectedEntry.Id {
				if realizedGain.Amount != expectedEntry.Amount ||
					realizedGain.Asset != expectedEntry.Asset ||
					realizedGain.SellTransactionId != expectedEntry.SellTransactionId ||
					realizedGain.BuyTransactionsId != expectedEntry.BuyTransactionsId ||
					realizedGain.IsProfit != expectedEntry.IsProfit ||
					realizedGain.TaxRate != expectedEntry.TaxRate ||
					realizedGain.Quantity != expectedEntry.Quantity ||
					realizedGain.BuyPrice != expectedEntry.BuyPrice ||
					realizedGain.SellPrice != expectedEntry.SellPrice {
					t.Errorf("For asset %s, expected %+v, but got %+v", expectedEntry.Asset, expectedEntry, realizedGain)
				}
				break
			}
		}
	}
}
