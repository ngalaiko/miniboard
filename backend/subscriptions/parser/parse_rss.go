package parser

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func parseRSS(data []byte) (*Feed, error) {
	var feed *Feed

	z := html.NewTokenizer(bytes.NewReader(data))
	z.AllowCDATA(true)
	for {
		switch z.Next() {
		case html.StartTagToken:
			tn, _ := z.TagName()
			if string(tn) == "channel" {
				f, err := parseRSSChannel(z)
				if err != nil {
					return nil, err
				}
				feed = f
			}
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return feed, nil
			}
			return nil, z.Err()
		}
	}
}

func parseRSSChannel(z *html.Tokenizer) (*Feed, error) {
	feed := &Feed{
		Items: []*Item{},
	}
	for tt := z.Next(); tt != html.EndTagToken; tt = z.Next() {
		switch tt {
		case html.StartTagToken:
			tn, _ := z.TagName()
			switch string(tn) {
			case tagTitle:
				title, err := parseText(z)
				if err != nil {
					return nil, err
				}
				feed.Title = title
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
			return nil, z.Err()
		}
	}
	return feed, nil
}

func parseItem(z *html.Tokenizer) (*Item, error) {
	item := &Item{}
	for tt := z.Next(); tt != html.EndTagToken; tt = z.Next() {
		switch tt {
		case html.ErrorToken:
			return nil, z.Err()
		case html.StartTagToken:
			tn, _ := z.TagName()
			switch string(tn) {
			case "title":
				text, err := parseText(z)
				if err != nil {
					return nil, err
				}
				item.Title = text
			case "link":
				text, err := parseText(z)
				if err != nil {
					return nil, err
				}
				item.Link = text
			}
		}
	}
	return item, nil
}

func parseImage(z *html.Tokenizer) (*Image, error) {
	var image *Image
	for tt := z.Next(); tt != html.EndTagToken; tt = z.Next() {
		switch tt {
		case html.ErrorToken:
			return nil, z.Err()
		case html.StartTagToken:
			tn, _ := z.TagName()
			if string(tn) == "url" {
				url, err := parseText(z)
				if err != nil {
					return nil, err
				}
				image = &Image{
					URL: url,
				}
			}
		}
	}
	return image, nil
}

func parseText(z *html.Tokenizer) (string, error) {
	var text string
	for tt := z.Next(); tt != html.EndTagToken; tt = z.Next() {
		switch tt {
		case html.ErrorToken:
			return "", z.Err()
		case html.TextToken:
			text = string(z.Text())
		}
	}
	if strings.HasPrefix(text, "<![CDATA[") {
		return text[9 : len(text)-3], nil
	}
	return text, nil
}
