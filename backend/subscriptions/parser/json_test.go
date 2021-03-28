package parser

import (
	"strings"
	"testing"
)

func Test_Parse_json__JsonFeed(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "My Example Feed",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"icon": "https://example.org/icon.png",
		"items": [
			{
				"id": "2",
				"content_text": "This is a second item.",
				"url": "https://example.org/second-item"
			},
			{
				"id": "1",
				"content_html": "<p>Hello, world!</p>",
				"url": "https://example.org/initial-post"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "My Example Feed" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Image.URL != "https://example.org/icon.png" {
		t.Errorf("Incorrect image url, got: %s", feed.Image.URL)
	}

	if feed.Link != "https://example.org/" {
		t.Errorf("Incorrect site URL, got: %s", feed.Link)
	}

	if len(feed.Items) != 2 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "https://example.org/second-item" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}

	if feed.Items[0].Title != "This is a second item." {
		t.Errorf(`Incorrect entry title, got: "%s"`, feed.Items[0].Title)
	}

	if feed.Items[0].Content != "This is a second item." {
		t.Errorf("Incorrect entry content, got: %s", feed.Items[0].Content)
	}

	if feed.Items[1].Link != "https://example.org/initial-post" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[1].Link)
	}

	if feed.Items[1].Title != "https://example.org/initial-post" {
		t.Errorf(`Incorrect entry title, got: "%s"`, feed.Items[1].Title)
	}

	if feed.Items[1].Content != "<p>Hello, world!</p>" {
		t.Errorf("Incorrect entry content, got: %s", feed.Items[0].Content)
	}
}

func Test_Parse_json__Podcast(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"user_comment": "This is a podcast feed. You can add this feed to your podcast client using the following URL: http://therecord.co/feed.json",
		"title": "The Record",
		"home_page_url": "http://therecord.co/",
		"feed_url": "http://therecord.co/feed.json",
		"items": [
			{
				"id": "http://therecord.co/chris-parrish",
				"title": "Special #1 - Chris Parrish",
				"url": "http://therecord.co/chris-parrish",
				"content_text": "Chris has worked at Adobe and as a founder of Rogue Sheep, which won an Apple Design Award for Postage. Chris’s new company is Aged & Distilled with Guy English — which shipped Napkin, a Mac app for visual collaboration. Chris is also the co-host of The Record. He lives on Bainbridge Island, a quick ferry ride from Seattle.",
				"content_html": "Chris has worked at <a href=\"http://adobe.com/\">Adobe</a> and as a founder of Rogue Sheep, which won an Apple Design Award for Postage. Chris’s new company is Aged & Distilled with Guy English — which shipped <a href=\"http://aged-and-distilled.com/napkin/\">Napkin</a>, a Mac app for visual collaboration. Chris is also the co-host of The Record. He lives on <a href=\"http://www.ci.bainbridge-isl.wa.us/\">Bainbridge Island</a>, a quick ferry ride from Seattle.",
				"summary": "Brent interviews Chris Parrish, co-host of The Record and one-half of Aged & Distilled.",
				"date_published": "2014-05-09T14:04:00-07:00",
				"attachments": [
					{
						"url": "http://therecord.co/downloads/The-Record-sp1e1-ChrisParrish.m4a",
						"mime_type": "audio/x-m4a",
						"size_in_bytes": 89970236,
						"duration_in_seconds": 6629
					}
				]
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "The Record" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Link != "http://therecord.co/" {
		t.Errorf("Incorrect site URL, got: %s", feed.Link)
	}

	if len(feed.Items) != 1 {
		t.Fatalf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://therecord.co/chris-parrish" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}

	if feed.Items[0].Title != "Special #1 - Chris Parrish" {
		t.Errorf(`Incorrect entry title, got: "%s"`, feed.Items[0].Title)
	}
}

func Test_Parse_json__FeedWithRelativeURL(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "Example",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"id": "2347259",
				"url": "something.html",
				"date_published": "2016-02-09T14:22:00-07:00"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "https://example.org/something.html" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}
}

func Test_Parse_json__FeedWithoutTitle(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"id": "2347259",
				"url": "https://example.org/2347259",
				"content_text": "Cats are neat. \n\nhttps://example.org/cats",
				"date_published": "2016-02-09T14:22:00-07:00"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "https://example.org/" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}
}

func Test_Parse_json__FeedItemWithoutTitle(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "My Example Feed",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"url": "https://example.org/item"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Title != "https://example.org/item" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_json__TruncateItemTitle(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "My Example Feed",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"title": "` + strings.Repeat("a", 200) + `"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if len(feed.Items[0].Title) != 103 {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_json__ItemTitleWithXMLTags(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "My Example Feed",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"title": "</example>"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Title != "</example>" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_json__InvalidJSON(t *testing.T) {
	data := `garbage`
	_, err := Parse([]byte(data), &testLogger{})
	if err == nil {
		t.Error("Parse should returns an error")
	}
}

func TestParseFeedItemWithInvalidDate(t *testing.T) {
	data := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "My Example Feed",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"id": "2347259",
				"url": "https://example.org/2347259",
				"content_text": "Cats are neat. \n\nhttps://example.org/cats",
				"date_published": "Tomorrow"
			}
		]
	}`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}

	if feed.Items[0].Date != nil {
		t.Errorf("Incorrect entry date, got: %v", feed.Items[0].Date)
	}
}
