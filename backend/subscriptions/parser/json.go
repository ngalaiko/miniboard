package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func parseJSON(data []byte) (*Feed, error) {
	feed := &jsonFeed{}
	if err := json.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse JSON feed: %s", err)
	}

	return feed.Convert()
}

type jsonFeed struct {
	Title string     `json:"title"`
	Link  string     `json:"home_page_url"`
	Icon  string     `json:"icon"`
	Items []jsonItem `json:"items"`
}

type jsonItem struct {
	Link          string `json:"url"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Text          string `json:"content_text"`
	DatePublished string `json:"date_published"`
	DateModified  string `json:"date_modified"`
}

func (f *jsonFeed) Convert() (*Feed, error) {
	feed := &Feed{}
	feed.Link = f.Link
	feed.Title = strings.TrimSpace(f.Title)
	if feed.Title == "" {
		feed.Title = feed.Link
	}

	if f.Icon != "" {
		feed.Image = &Image{
			URL: f.Icon,
		}
	}

	for _, i := range f.Items {
		item, err := i.Convert()
		if err != nil {
			return nil, err
		}
		itemLink, err := absoluteURL(feed.Link, item.Link)
		if err == nil {
			item.Link = itemLink
		}
		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

func (i *jsonItem) Convert() (*Item, error) {
	item := &Item{
		Link:  i.Link,
		Title: i.title(),
		Date:  i.date(),
	}
	if item.Title == "" {
		item.Title = item.Link
	}
	return item, nil
}

func (i *jsonItem) date() time.Time {
	for _, value := range []string{i.DatePublished, i.DateModified} {
		if value != "" {
			d, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return time.Now()
			}

			return d
		}
	}

	return time.Now()
}

func (i *jsonItem) title() string {
	for _, value := range []string{i.Title, i.Summary, i.Text, i.Link} {
		if value != "" {
			return truncate(value)
		}
	}

	return i.Link
}

func truncate(str string) string {
	max := 100
	str = strings.TrimSpace(str)
	if len(str) > max {
		return str[:max] + "..."
	}

	return str
}
