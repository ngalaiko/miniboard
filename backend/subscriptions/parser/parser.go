package parser

import "fmt"

// Item contains relevant feed item information.
type Item struct {
	Title string
	Link  string
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

// Parse returns a parsed feed.
func Parse(data []byte) (*Feed, error) {
	switch detectType(data) {
	case feedTypeRSS:
		return parseRSS(data)
	case feedTypeRDF:
		return parseRDF(data)
	case feedTypeJSON:
		return parseJSON(data)
	case feedTypeAtom03:
		return parseAtom03(data)
	case feedTypeAtom10:
		return parseAtom10(data)
	default:
		return nil, fmt.Errorf("unkwown type")
	}
}
