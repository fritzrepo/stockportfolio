package depot

import (
	"testing"

	"golang.org/x/text/currency"
)

func TestComputeTransactions(t *testing.T) {
	expected := map[string]DepotEntry{
		"Apple": {
			assetType:    "stock",
			asset:        "Apple",
			tickerSymbol: "AAPL",
			quantity:     20,
			price:        100.5,
			//totalPrice:   2014,
			currency: currency.EUR,
		},
		"BASF": {
			assetType:    "stock",
			asset:        "BASF",
			tickerSymbol: "BAS1",
			quantity:     100,
			price:        45.5,
			//totalPrice:   4555,
			currency: currency.EUR,
		},
	}

	result := ComputeTransactions()

	for key, expectedEntry := range expected {
		resultEntry, exists := result[key]
		if !exists {
			t.Errorf("Expected asset %s not found in result", key)
			continue
		}
		if resultEntry != expectedEntry {
			t.Errorf("For asset %s, expected %+v, but got %+v", key, expectedEntry, resultEntry)
		}
	}
}
