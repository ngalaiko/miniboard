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

	return feed.Convert()
}

type rdfFeed struct {
	XMLName xml.Name   `xml:"RDF"`
	Title   string     `xml:"channel>title"`
	Link    string     `xml:"channel>link"`
	Image   *rdfImage  `xml:"image"`
	Items   []*rdfItem `xml:"item"`
}

func (f *rdfFeed) Convert() (*Feed, error) {
	feed := &Feed{
		Title: strings.TrimSpace(f.Title),
		Link:  strings.TrimSpace(f.Link),
	}
	image, err := f.Image.Convert()
	if err != nil {
		return nil, err
	}
	feed.Image = image
	for _, i := range f.Items {
		item, err := i.Convert()
		if err != nil {
			return nil, err
		}

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
	return feed, nil
}

type rdfImage struct {
	URL string `xml:"url"`
}

func (i *rdfImage) Convert() (*Image, error) {
	if i == nil {
		return nil, nil
	}
	return &Image{
		URL: i.URL,
	}, nil
}

type rdfItem struct {
	Title          string `xml:"title"`
	Link           string `xml:"link"`
	DublinCoreDate string `xml:"http://purl.org/dc/elements/1.1/ date"`
}

func (i *rdfItem) Convert() (*Item, error) {
	return &Item{
		Title: strings.TrimSpace(i.Title),
		Link:  strings.TrimSpace(i.Link),
		Date:  i.date(),
	}, nil
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
