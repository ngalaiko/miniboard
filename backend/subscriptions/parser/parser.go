package parser

import "fmt"

// Item contains relevant feed item information.
type Item struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// Image is a feed's image.
type Image struct {
	URL string `json:"url"`
}

// Feed contains relevant feed information.
type Feed struct {
	Title string  `json:"title"`
	Image *Image  `json:"image"`
	Items []*Item `json:"items"`
}

// Parse returns a parsed feed.
func Parse(data []byte) (*Feed, error) {
	switch detectType(data) {
	case feedTypeRSS:
		return parseRSS(data)
	case feedTypeAtom:
		return parseAtom(data)
	case feedTypeJSON:
		return parseJSON(data)
	default:
		return nil, fmt.Errorf("unkwown type")
	}
}
