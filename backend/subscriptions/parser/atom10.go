package parser

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
	"time"
)

func parseAtom10(data []byte, logger logger) (*Feed, error) {
	feed := &atom10Feed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse atom10 feed: %s", err)
	}
	return feed.Convert(logger)
}

type atom10Feed struct {
	XMLName xml.Name      `xml:"http://www.w3.org/2005/Atom feed"`
	Title   atom10Text    `xml:"title"`
	Links   atomLinks     `xml:"link"`
	IconURL atom10Text    `xml:"icon"`
	Items   []*atom10Item `xml:"entry"`
}

func (f *atom10Feed) Convert(logger logger) (*Feed, error) {
	feed := &Feed{}
	feed.Link = f.Links.originalLink()
	feed.Title = f.Title.String()
	if feed.Title == "" {
		feed.Title = feed.Link
	}
	if iconURL := f.IconURL.String(); iconURL != "" {
		feedIconURL, err := absoluteURL(feed.Link, iconURL)
		if err != nil {
			return nil, err
		}
		feed.Image = &Image{
			URL: feedIconURL,
		}
	}
	for _, i := range f.Items {
		item := i.Convert(logger)
		itemLink, err := absoluteURL(feed.Link, item.Link)
		if err == nil {
			item.Link = itemLink
		}
		if item.Title == "" {
			item.Title = item.Link
		}
		feed.Items = append(feed.Items, item)
	}
	return feed, nil
}

type atom10Item struct {
	Title     atom10Text `xml:"title"`
	Links     atomLinks  `xml:"link"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
	Summary   atom10Text `xml:"summary"`
	Content   atom10Text `xml:"http://www.w3.org/2005/Atom content"`
}

func (i *atom10Item) Convert(logger logger) *Item {
	return &Item{
		Title:   i.Title.String(),
		Link:    i.Links.originalLink(),
		Date:    i.date(logger),
		Content: i.content(),
	}
}

func (i *atom10Item) content() *string {
	if content := i.Content.String(); content != "" {
		return &content
	}

	if summary := i.Summary.String(); summary != "" {
		return &summary
	}

	return nil
}

func (i *atom10Item) date(logger logger) *time.Time {
	dateText := i.Published
	if dateText == "" {
		dateText = i.Updated
	}

	if dateText != "" {
		t, err := parseDateTime(dateText)
		if err != nil {
			logger.Error("atom10: failed to parse date '%s': %w", dateText, err)
			return nil
		}
		return t
	}

	return nil
}

type atomLinks []*atomLink

type atomLink struct {
	URL    string `xml:"href,attr"`
	Type   string `xml:"type,attr"`
	Rel    string `xml:"rel,attr"`
	Length string `xml:"length,attr"`
}

func (a atomLinks) originalLink() string {
	for _, link := range a {
		if strings.ToLower(link.Rel) == "alternate" {
			return strings.TrimSpace(link.URL)
		}

		if link.Rel == "" && link.Type == "" {
			return strings.TrimSpace(link.URL)
		}
	}

	return ""
}

type atom10Text struct {
	Type string `xml:"type,attr"`
	Data string `xml:",chardata"`
	XML  string `xml:",innerxml"`
}

func (a *atom10Text) String() string {
	content := ""

	switch {
	case a.Type == "xhtml":
		content = a.XML
	default:
		content = a.Data
	}

	return html.UnescapeString(strings.TrimSpace(content))
}
