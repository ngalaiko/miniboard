package parser

import (
	"fmt"
	"time"
)

// Item contains relevant feed item information.
type Item struct {
	Title   string
	Link    string
	Date    *time.Time
	Content string
}

// Image is a feed's image.
type Image struct {
	URL string
}

// Feed contains relevant feed information.
type Feed struct {
	Title string
	Link  string
	Image *Image
	Items []*Item
}

type logger interface {
	Error(string, ...interface{})
}

// Parse returns a parsed feed.
func Parse(data []byte, logger logger) (*Feed, error) {
	switch detectType(data) {
	case feedTypeRSS:
		return parseRSS(data, logger)
	case feedTypeRDF:
		return parseRDF(data, logger)
	case feedTypeJSON:
		return parseJSON(data, logger)
	case feedTypeAtom03:
		return parseAtom03(data, logger)
	case feedTypeAtom10:
		return parseAtom10(data, logger)
	default:
		return nil, fmt.Errorf("unkwown type")
	}
}
