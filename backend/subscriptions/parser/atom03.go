package parser

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"html"
	"strings"
	"time"
)

func parseAtom03(data []byte, logger logger) (*Feed, error) {
	feed := &atom03Feed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse atom10 feed: %s", err)
	}
	return feed.Convert(logger), nil
}

type atom03Feed struct {
	Title atom03Text   `xml:"title"`
	Links atomLinks    `xml:"link"`
	Items []atom03Item `xml:"entry"`
}

func (f *atom03Feed) Convert(logger logger) *Feed {
	feed := &Feed{
		Title: f.Title.String(),
		Link:  f.Links.originalLink(),
	}
	if feed.Title == "" {
		feed.Title = feed.Link
	}
	for _, i := range f.Items {
		item := i.Convert(logger)
		if item.Title == "" {
			item.Title = item.Link
		}
		feed.Items = append(feed.Items, item)
	}
	return feed
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
	Content  atom03Text `xml:"content"`
	Summary  atom03Text `xml:"summary"`
}

func (i *atom03Item) Convert(logger logger) *Item {
	return &Item{
		Title:   i.Title.String(),
		Link:    i.Links.originalLink(),
		Date:    i.date(logger),
		Content: i.content(),
	}
}

func (i *atom03Item) content() string {
	content := i.Content.String()
	if content != "" {
		return content
	}

	summary := i.Summary.String()
	if summary != "" {
		return summary
	}

	return ""
}

func (i *atom03Item) date(logger logger) time.Time {
	dateText := ""
	for _, value := range []string{i.Issued, i.Modified, i.Created} {
		if value != "" {
			dateText = value
			break
		}
	}

	if dateText != "" {
		t, err := parseDateTime(dateText)
		if err != nil {
			logger.Error("atom03: failed to parse date '%s': %w", dateText, err)
			return time.Now()
		}
		return *t
	}

	return time.Now()
}
