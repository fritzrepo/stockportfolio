package comdirect

import (
	"fmt"
	"strings"
	"time"
)

// FlexTime is a custom time type that can handle multiple time formats during JSON unmarshalling.
type FlexTime struct {
	time.Time
}

func (ft *FlexTime) UnmarshalJSON(data []byte) error {
	// Anführungszeichen entfernen
	s := strings.Trim(string(data), `"`)

	layouts := []string{
		"2006-01-02T15:04:05Z07:00", // Standard RFC3339
		"2006-01-02T15:04:05Z07",    // Kurzes Offset (z.B. +01)
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			ft.Time = t
			return nil
		}
	}

	return fmt.Errorf("cannot parse time: %q", s)
}
