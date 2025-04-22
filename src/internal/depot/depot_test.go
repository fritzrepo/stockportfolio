package depot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/fritzrepo/stockportfolio/internal/storage"
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

func TestComputeTransactions(t *testing.T) {

	type TestCases = []struct {
		Name          string                `json:"name"`
		ExpectedGains []RealizedGain        `json:"expectedGains"`
		ExpectedDepot map[string]DepotEntry `json:"expectedDepot"`
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

			var uuidGenerator = NewMockUUIDGenerator()
			filenameTrans := fmt.Sprintf("../../testdata/depot/RawTransactionsTest%d.csv", i)
			i = i + 1

			store := storage.NewCsvStorage(filenameTrans, uuidGenerator.GetUUID)
			dep := NewDepot(uuidGenerator.GetUUID, &store)

			err := dep.ComputeTransactions()
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
