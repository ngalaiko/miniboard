package parser

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"html"
	"strings"
	"time"
)

func parseAtom03(data []byte) (*Feed, error) {
	feed := &atom03Feed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse atom10 feed: %s", err)
	}
	return feed.Convert()
}

type atom03Feed struct {
	Title atom03Text   `xml:"title"`
	Links atomLinks    `xml:"link"`
	Items []atom03Item `xml:"entry"`
}

func (f *atom03Feed) Convert() (*Feed, error) {
	feed := &Feed{
		Title: f.Title.String(),
		Link:  f.Links.originalLink(),
	}
	if feed.Title == "" {
		feed.Title = feed.Link
	}
	for _, i := range f.Items {
		item, err := i.Convert()
		if err != nil {
			return nil, err
		}
		if item.Title == "" {
			item.Title = item.Link
		}
		feed.Items = append(feed.Items, item)
	}
	return feed, nil
}

type atom03Text struct {
	Type string `xml:"type,attr"`
	Mode string `xml:"mode,attr"`
	Data string `xml:",chardata"`
	XML  string `xml:",innerxml"`
}

func (a *atom03Text) String() string {
	content := ""

	switch {
	case a.Mode == "xml":
		content = a.XML
	case a.Mode == "escaped":
		content = a.Data
	case a.Mode == "base64":
		b, err := base64.StdEncoding.DecodeString(a.Data)
		if err == nil {
			content = string(b)
		}
	default:
		content = a.Data
	}

	if a.Type != "text/html" {
		content = html.EscapeString(content)
	}

	return strings.TrimSpace(content)
}

type atom03Item struct {
	Title    atom03Text `xml:"title"`
	Links    atomLinks  `xml:"link"`
	Modified string     `xml:"modified"`
	Issued   string     `xml:"issued"`
	Created  string     `xml:"created"`
}

func (i *atom03Item) Convert() (*Item, error) {
	return &Item{
		Title: i.Title.String(),
		Link:  i.Links.originalLink(),
		Date:  i.date(),
	}, nil
}

func (i *atom03Item) date() time.Time {
	dateText := ""
	for _, value := range []string{i.Issued, i.Modified, i.Created} {
		if value != "" {
			dateText = value
			break
		}
	}

	if dateText != "" {
		result, err := time.Parse(time.RFC3339, dateText)
		if err != nil {
			return time.Now()
		}

		return result
	}

	return time.Now()
}
