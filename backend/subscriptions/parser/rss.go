package parser

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"
)

func parseRSS(data []byte) (*Feed, error) {
	feed := &rssFeed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse RSS feed: %s", err)
	}

	return feed.Convert(), nil
}

func (f *rssFeed) Convert() *Feed {
	feed := &Feed{
		Title: f.Title,
		Link:  f.link(),
		Image: f.Image.Convert(),
	}

	if feed.Title == "" {
		feed.Title = feed.Link
	}

	for _, item := range f.Items {
		i := item.Convert()

		if i.Link == "" {
			i.Link = feed.Link
		} else {
			itemURL, err := absoluteURL(feed.Link, i.Link)
			if err == nil {
				i.Link = itemURL
			}

			if i.Title == "" {
				i.Title = i.Link
			}
		}

		feed.Items = append(feed.Items, i)
	}

	return feed
}

func absoluteURL(baseURL, input string) (string, error) {
	if strings.HasPrefix(input, "//") {
		input = "https://" + input[2:]
	}

	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("unable to parse input URL: %v", err)
	}

	if u.IsAbs() {
		return u.String(), nil
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("unable to parse base URL: %v", err)
	}

	return base.ResolveReference(u).String(), nil
}

func (i *rssImage) Convert() *Image {
	if i == nil {
		return nil
	}

	return &Image{
		URL: i.url(),
	}
}

func (i *rssImage) url() string {
	for _, element := range i.URL {
		if element.XMLName.Space == "" {
			return strings.TrimSpace(element.Data)
		}
	}

	return ""
}

func (f *rssFeed) link() string {
	for _, element := range f.Links {
		if element.XMLName.Space == "" {
			return strings.TrimSpace(element.Data)
		}
	}

	return ""
}

func (i *rssItem) Convert() *Item {
	return &Item{
		Title: i.title(),
		Link:  i.link(),
		Date:  i.date(),
	}
}

func (i *rssItem) link() string {
	if i.FeedBurnerLink != "" {
		return i.FeedBurnerLink
	}
	for _, link := range i.Links {
		if link.XMLName.Space == "http://www.w3.org/2005/Atom" && link.Href != "" && isValidLinkRelation(link.Rel) {
			return strings.TrimSpace(link.Href)
		}

		if link.Data != "" {
			return strings.TrimSpace(link.Data)
		}
	}

	return ""
}

func isValidLinkRelation(rel string) bool {
	switch rel {
	case "", "alternate", "enclosure", "related", "self", "via":
		return true
	default:
		if strings.HasPrefix(rel, "http") {
			return true
		}
		return false
	}
}

func (i *rssItem) title() string {
	var title string

	for _, rssTitle := range i.Title {
		switch rssTitle.XMLName.Space {
		case "http://search.yahoo.com/mrss/":
			// Ignore title in media namespace
		case "http://purl.org/dc/elements/1.1/":
			title = rssTitle.Data
		default:
			title = rssTitle.Data
		}

		if title != "" {
			break
		}
	}

	return html.UnescapeString(strings.TrimSpace(title))
}

func (i *rssItem) date() time.Time {
	value := i.PubDate
	if i.DublinCoreDate != "" {
		value = i.DublinCoreDate
	}

	if value == "" {
		return time.Now()
	}

	for _, layout := range []string{
		time.RFC3339,
		time.RFC822,
		time.RFC822Z,
		time.RFC1123,
		time.RFC1123Z,
	} {
		result, err := time.Parse(layout, value)
		if err == nil {
			return result
		}
	}

	return time.Now()
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Title   string     `xml:"channel>title"`
	Links   []rssLink  `xml:"channel>link"`
	Image   *rssImage  `xml:"channel>image"`
	Items   []*rssItem `xml:"channel>item"`
}

type rssImage struct {
	URL []rssLink `xml:"url"`
}

type rssItem struct {
	GUID           string         `xml:"guid"`
	Title          []rssItemtitle `xml:"title"`
	Links          []rssLink      `xml:"link"`
	PubDate        string         `xml:"pubDate"`
	DublinCoreDate string         `xml:"http://purl.org/dc/elements/1.1/ date"`
	feedBurnerRssItem
}

type feedBurnerRssItem struct {
	FeedBurnerLink string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origLink"`
}

type rssLink struct {
	XMLName xml.Name
	Data    string `xml:",chardata"`
	Href    string `xml:"href,attr"`
	Rel     string `xml:"rel,attr"`
}

type rssItemtitle struct {
	XMLName xml.Name
	Data    string `xml:",chardata"`
	Inner   string `xml:",innerxml"`
}
