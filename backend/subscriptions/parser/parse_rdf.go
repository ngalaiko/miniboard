package parser

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

const tagTitle = "title"

func parseRDF(data []byte) (*Feed, error) {
	var feed *Feed

	z := html.NewTokenizer(bytes.NewReader(data))
	z.AllowCDATA(true)
	for {
		switch z.Next() {
		case html.StartTagToken:
			tn, _ := z.TagName()
			switch string(tn) {
			case "channel":
				f, err := parseRDFChannel(z)
				if err != nil {
					return nil, err
				}
				feed = f
			case "image":
				image, err := parseImage(z)
				if err != nil {
					return nil, err
				}
				feed.Image = image
			case "item":
				item, err := parseItem(z)
				if err != nil {
					return nil, err
				}
				feed.Items = append(feed.Items, item)
			}
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return feed, nil
			}
			return nil, z.Err()
		}
	}
}

func parseRDFChannel(z *html.Tokenizer) (*Feed, error) {
	feed := &Feed{
		Items: []*Item{},
	}
	for tt := z.Next(); tt != html.EndTagToken; tt = z.Next() {
		switch tt {
		case html.StartTagToken:
			tn, _ := z.TagName()
			if string(tn) == tagTitle {
				title, err := parseText(z)
				if err != nil {
					return nil, err
				}
				feed.Title = title
			}
		case html.ErrorToken:
			return nil, z.Err()
		}
	}
	return feed, nil
}
