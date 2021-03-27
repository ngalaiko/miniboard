package parser

import (
	"fmt"
	"time"
)

func parseDateTime(raw string) (*time.Time, error) {
	for _, layout := range []string{
		time.RFC3339,
		time.RFC822,
		time.RFC822Z,
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 2 Jan 2006 15:04:05 -0700", // time.RFC1123Z but single-digit date
		"Mon, 2 Jan 2006 15:04:05 MST",   // time.RFC1123 but single-digit date
	} {
		result, err := time.Parse(layout, raw)
		if err == nil {
			return &result, nil
		}
	}
	return nil, fmt.Errorf("invalid format: '%s'", raw)
}
