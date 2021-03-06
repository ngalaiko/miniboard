package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func parseJSON(data []byte, logger logger) (*Feed, error) {
	feed := &jsonFeed{}
	if err := json.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("unable to parse JSON feed: %w", err)
	}

	return feed.Convert(logger), nil
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
	HTML          string `json:"content_html"`
	DatePublished string `json:"date_published"`
	DateModified  string `json:"date_modified"`
}

func (f *jsonFeed) Convert(logger logger) *Feed {
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
		item := i.Convert(logger)
		itemLink, err := absoluteURL(feed.Link, item.Link)
		if err == nil {
			item.Link = itemLink
		}
		feed.Items = append(feed.Items, item)
	}

	return feed
}

func (i *jsonItem) Convert(logger logger) *Item {
	item := &Item{
		Link:    i.Link,
		Title:   i.title(),
		Date:    i.date(logger),
		Content: i.content(),
	}
	if item.Title == "" {
		item.Title = item.Link
	}
	return item
}

func (i *jsonItem) content() *string {
	for _, value := range []string{i.HTML, i.Text, i.Summary} {
		if value != "" {
			return &value
		}
	}

	return nil
}

func (i *jsonItem) date(logger logger) *time.Time {
	for _, value := range []string{i.DatePublished, i.DateModified} {
		if value != "" {
			t, err := parseDateTime(value)
			if err != nil {
				logger.Error("json: failed to parse date '%s': %s", value, err)
				return nil
			}
			return t
		}
	}

	return nil
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
