package parser

import (
	"testing"
)

func Test_parseRSS__Rss2Sample(t *testing.T) {
	data := `
		<?xml version="1.0"?>
		<rss version="2.0">
		<channel>
			<title>Liftoff News</title>
			<link>http://liftoff.msfc.nasa.gov/</link>
			<description>Liftoff to Space Exploration.</description>
			<language>en-us</language>
			<pubDate>Tue, 10 Jun 2003 04:00:00 GMT</pubDate>
			<lastBuildDate>Tue, 10 Jun 2003 09:41:01 GMT</lastBuildDate>
			<docs>http://blogs.law.harvard.edu/tech/rss</docs>
			<generator>Weblog Editor 2.0</generator>
			<managingEditor>editor@example.com</managingEditor>
			<webMaster>webmaster@example.com</webMaster>
			<image>
				<url>http://liftoff.msfc.nasa.gov/image.png</url>
			</image>
			<item>
				<title>Star City</title>
				<link>http://liftoff.msfc.nasa.gov/news/2003/news-starcity.asp</link>
				<description>How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, language and protocol at Russia's &lt;a href="http://howe.iki.rssi.ru/GCTC/gctc_e.htm"&gt;Star City&lt;/a&gt;.</description>
				<pubDate>Tue, 03 Jun 2003 09:39:21 GMT</pubDate>
				<guid>http://liftoff.msfc.nasa.gov/2003/06/03.html#item573</guid>
			</item>
			<item>
				<description>Sky watchers in Europe, Asia, and parts of Alaska and Canada will experience a &lt;a href="http://science.nasa.gov/headlines/y2003/30may_solareclipse.htm"&gt;partial eclipse of the Sun&lt;/a&gt; on Saturday, May 31st.</description>
				<pubDate>Fri, 30 May 2003 11:06:42 GMT</pubDate>
				<guid>http://liftoff.msfc.nasa.gov/2003/05/30.html#item572</guid>
			</item>
			<item>
				<title>The Engine That Does More</title>
				<link>http://liftoff.msfc.nasa.gov/news/2003/news-VASIMR.asp</link>
				<description>Before man travels to Mars, NASA hopes to design new engines that will let us fly through the Solar System more quickly.  The proposed VASIMR engine would do that.</description>
				<pubDate>Tue, 27 May 2003 08:37:32 GMT</pubDate>
				<guid>http://liftoff.msfc.nasa.gov/2003/05/27.html#item571</guid>
			</item>
			<item>
				<title>Astronauts' Dirty Laundry</title>
				<link>http://liftoff.msfc.nasa.gov/news/2003/news-laundry.asp</link>
				<description>Compared to earlier spacecraft, the International Space Station has many luxuries, but laundry facilities are not one of them.  Instead, astronauts have other options.</description>
				<pubDate>Tue, 20 May 2003 08:56:02 GMT</pubDate>
				<guid>http://liftoff.msfc.nasa.gov/2003/05/20.html#item570</guid>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "Liftoff News" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Link != "http://liftoff.msfc.nasa.gov/" {
		t.Errorf("Incorrect site URL, got: %s", feed.Link)
	}

	if feed.Image == nil {
		t.Errorf("Image is missing")
	}

	if feed.Image.URL != "http://liftoff.msfc.nasa.gov/image.png" {
		t.Errorf("Incorrect image url, got: %s", feed.Image.URL)
	}

	if len(feed.Items) != 4 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://liftoff.msfc.nasa.gov/news/2003/news-starcity.asp" {
		t.Errorf("Incorrect entry Link, got: %s", feed.Items[0].Link)
	}

	if feed.Items[0].Title != "Star City" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_parseRSS__FeedWithoutTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0">
		<channel>
			<link>https://example.org/</link>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "https://example.org/" {
		t.Errorf("Incorrect feed title, got: %s", feed.Title)
	}
}

func Test_parseRSS__ItemWithoutTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0">
		<channel>
			<link>https://example.org/</link>
			<item>
				<link>https://example.org/item</link>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "https://example.org/item" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_parseRSS__ItemWithMediaTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:media="http://search.yahoo.com/mrss/">
		<channel>
			<link>https://example.org/</link>
			<item>
				<title>Item Title</title>
				<link>https://example.org/item</link>
				<media:title>Media Title</media:title>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Item Title" {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_parseRSS__ItemWithDCTitleOnly(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:media="http://search.yahoo.com/mrss/" xmlns:dc="http://purl.org/dc/elements/1.1/">
		<channel>
			<link>https://example.org/</link>
			<item>
				<dc:title>Item Title</dc:title>
				<link>https://example.org/item</link>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Item Title" {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_parseRSS__ItemWithoutLink(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0">
		<channel>
			<link>https://example.org/</link>
			<item>
				<guid isPermaLink="false">1234</guid>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "https://example.org/" {
		t.Errorf("Incorrect entry link, got: %s", feed.Items[0].Link)
	}
}

func Test_parseRSS__ItemWithAtomLink(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
		<channel>
			<link>https://example.org/</link>
			<item>
				<title>Test</title>
				<atom:link href="https://example.org/item" />
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Link != "https://example.org/" {
		t.Errorf("Incorrect site URL, got: %s", feed.Link)
	}

	if feed.Items[0].Link != "https://example.org/item" {
		t.Errorf("Incorrect entry link, got: %s", feed.Items[0].Link)
	}
}

func Test_parseRSS__ItemWithMultipleAtomLinks(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
		<channel>
			<link>https://example.org/</link>
			<item>
				<title>Test</title>
				<atom:link rel="payment" href="https://example.org/a" />
				<atom:link rel="http://foobar.tld" href="https://example.org/b" />
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "https://example.org/b" {
		t.Errorf("Incorrect entry link, got: %s", feed.Items[0].Link)
	}
}

func Test_parseRSS__ItemWithFeedBurnerLink(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:feedburner="http://rssnamespace.org/feedburner/ext/1.0">
		<channel>
			<title>Example</title>
			<link>http://example.org/</link>
			<item>
				<title>Item 1</title>
				<link>http://example.org/item1</link>
				<feedburner:origLink>http://example.org/original</feedburner:origLink>
			</item>
		</channel>
	</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "http://example.org/original" {
		t.Errorf("Incorrect entry content, got: %s", feed.Items[0].Link)
	}
}

func Test_parseRSS__ItemTitleWithWhitespaces(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rss version="2.0">
	<channel>
		<title>Example</title>
		<link>http://example.org</link>
		<item>
			<title>
				Some Title
			</title>
			<link>http://www.example.org/entries/1</link>
			<pubDate>Fri, 15 Jul 2005 00:00:00 -0500</pubDate>
		</item>
	</channel>
	</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Some Title" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_parseRSS__ItemWithRelativeLink(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0">
		<channel>
			<link>https://example.org/</link>
			<item>
				<link>item.html</link>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "https://example.org/item.html" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_parseRSS__InvalidXml(t *testing.T) {
	data := `garbage`
	_, err := parseRSS([]byte(data))
	if err == nil {
		t.Error("Parse should returns an error")
	}
}

func Test_parseRSS__ItemTitleWithHTMLEntity(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:slash="http://purl.org/rss/1.0/modules/slash/">
		<channel>
			<link>https://example.org/</link>
			<title>Example</title>
			<item>
				<title>&lt;/example&gt;</title>
				<link>http://www.example.org/entries/1</link>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "</example>" {
		t.Errorf(`Incorrect title, got: %q`, feed.Items[0].Title)
	}
}

func Test_parseRSS__ItemTitleWithNumericCharacterReference(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:slash="http://purl.org/rss/1.0/modules/slash/">
		<channel>
			<link>https://example.org/</link>
			<title>Example</title>
			<item>
				<title>&#931; &#xDF;</title>
				<link>http://www.example.org/article.html</link>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Σ ß" {
		t.Errorf(`Incorrect title, got: %q`, feed.Items[0].Title)
	}
}

func Test_parseRSS__ItemTitleWithDoubleEncodedEntities(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:slash="http://purl.org/rss/1.0/modules/slash/">
		<channel>
			<link>https://example.org/</link>
			<title>Example</title>
			<item>
				<title>&amp;#39;Text&amp;#39;</title>
				<link>http://www.example.org/article.html</link>
			</item>
		</channel>
		</rss>`

	feed, err := parseRSS([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "'Text'" {
		t.Errorf(`Incorrect title, got: %q`, feed.Items[0].Title)
	}
}
