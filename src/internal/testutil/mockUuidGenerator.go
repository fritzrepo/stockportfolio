package testutil

import (
	"fmt"

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

func (m *MockUUIDGenerator) GetUUID() uuid.UUID {
	if m.callCount >= len(m.uuids) {
		// Wenn alle UUIDs verwendet wurden, beginne von vorne
		m.callCount = 0
	}

	id := m.uuids[m.callCount]
	m.callCount++
	return id
}
