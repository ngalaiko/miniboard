package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func parseRDF(data []byte) (*Feed, error) {
	feed := &rdfFeed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse RDF feed: %s", err)
	}

	return feed.Convert(), nil
}

type rdfFeed struct {
	XMLName xml.Name   `xml:"RDF"`
	Title   string     `xml:"channel>title"`
	Link    string     `xml:"channel>link"`
	Image   *rdfImage  `xml:"image"`
	Items   []*rdfItem `xml:"item"`
}

func (f *rdfFeed) Convert() *Feed {
	feed := &Feed{
		Title: strings.TrimSpace(f.Title),
		Link:  strings.TrimSpace(f.Link),
		Image: f.Image.Convert(),
	}
	for _, i := range f.Items {
		item := i.Convert()

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
	Title          string `xml:"title"`
	Link           string `xml:"link"`
	DublinCoreDate string `xml:"http://purl.org/dc/elements/1.1/ date"`
}

func (i *rdfItem) Convert() *Item {
	return &Item{
		Title: strings.TrimSpace(i.Title),
		Link:  strings.TrimSpace(i.Link),
		Date:  i.date(),
	}
}

func (i *rdfItem) date() time.Time {
	if i.DublinCoreDate != "" {
		result, err := time.Parse(time.RFC3339, i.DublinCoreDate)
		if err != nil {
			return time.Now()
		}

		return result
	}

	return time.Now()
}