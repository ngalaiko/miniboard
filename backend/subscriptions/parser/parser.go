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
	default:
		return nil, fmt.Errorf("unkwown type")
	}
}
