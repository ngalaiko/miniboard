package parser

import (
	"fmt"
	"strings"
	"time"
)

var replacements = map[string]string{
	"пн":    "Mon",
	"вт":    "Tue",
	"ср":    "Wed",
	"чт":    "Thu",
	"пт":    "Fri",
	"сб":    "Sat",
	"вс":    "Sun",
	"янв.":  "Jan",
	"февр.": "Feb",
	"мар.":  "Mar",
	"апр.":  "Apr",
	"мая":   "May",
	"июня":  "Jun",
	"июля":  "Jul",
	"авг.":  "Aug",
	"сент.": "Sep",
	"окт.":  "Oct",
	"нояб.": "Nov",
	"дек.":  "Dec",
}

func parseDateTime(raw string) (*time.Time, error) {
	for from, to := range replacements {
		raw = strings.Replace(raw, from, to, 1)
	}
	for _, layout := range []string{
		time.RFC3339,
		time.RFC822,
		time.RFC822Z,
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 2 Jan 2006 15:04:05 -0700", // time.RFC1123Z but single-digit date
		"Mon, 2 Jan 2006 15:04:05 MST",   // time.RFC1123 but single-digit date
		"2006-01-02T15:04",
	} {
		result, err := time.Parse(layout, raw)
		if err == nil {
			return &result, nil
		}
	}
	return nil, fmt.Errorf("invalid format: '%s'", raw)
}
