package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func parseRDF(data []byte, logger logger) (*Feed, error) {
	feed := &rdfFeed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse RDF feed: %s", err)
	}

	return feed.Convert(logger), nil
}

type rdfFeed struct {
	XMLName xml.Name   `xml:"RDF"`
	Title   string     `xml:"channel>title"`
	Link    string     `xml:"channel>link"`
	Image   *rdfImage  `xml:"image"`
	Items   []*rdfItem `xml:"item"`
}

func (f *rdfFeed) Convert(logger logger) *Feed {
	feed := &Feed{
		Title: strings.TrimSpace(f.Title),
		Link:  strings.TrimSpace(f.Link),
		Image: f.Image.Convert(),
	}
	for _, i := range f.Items {
		item := i.Convert(logger)

		if item.Link == "" {
			item.Link = feed.Link
		} else {
			linkURL, err := absoluteURL(feed.Link, item.Link)
			if err == nil {
				item.Link = linkURL
			}
		}

		feed.Items = append(feed.Items, item)
	}
	return feed
}

type rdfImage struct {
	URL string `xml:"url"`
}

func (i *rdfImage) Convert() *Image {
	if i == nil {
		return nil
	}
	return &Image{
		URL: i.URL,
	}
}

type rdfItem struct {
	Title             string `xml:"title"`
	Link              string `xml:"link"`
	Description       string `xml:"description"`
	DublinCoreDate    string `xml:"http://purl.org/dc/elements/1.1/ date"`
	DublinCoreContent string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

func (i *rdfItem) Convert(logger logger) *Item {
	return &Item{
		Title:   strings.TrimSpace(i.Title),
		Link:    strings.TrimSpace(i.Link),
		Date:    i.date(logger),
		Content: i.content(),
	}
}

func (i *rdfItem) content() string {
	switch {
	case i.DublinCoreContent != "":
		return i.DublinCoreContent
	default:
		return i.Description
	}
}

func (i *rdfItem) date(logger logger) time.Time {
	if i.DublinCoreDate != "" {
		t, err := parseDateTime(i.DublinCoreDate)
		if err != nil {
			logger.Error("rdf: failed to parse date '%s': %s", i.DublinCoreDate, err)
			return time.Now()
		}

		return *t
	}

	return time.Now()
}
