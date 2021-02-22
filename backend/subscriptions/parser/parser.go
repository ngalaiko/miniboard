package parser

import "fmt"

// Item contains relevant feed item information.
type Item struct {
	Title *string
	Link  string
}

// Feed contains relevant feed information.
type Feed struct {
	Title   string
	IconURL *string
	Items   []*Item
}

// Parse returns a parsed feed.
func Parse(data []byte) (*Feed, error) {
	return nil, fmt.Errorf("not implemented")
}
