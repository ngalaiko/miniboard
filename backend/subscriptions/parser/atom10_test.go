package parser

import (
	"testing"
)

func Test_Parse_atom10__AtomSample(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <updated>2003-12-13T18:30:02Z</updated>
	  <icon>/icon.png</icon>
	  <author>
		<name>John Doe</name>
	  </author>
	  <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>
	  <entry>
		<title>Atom-Powered Robots Run Amok</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "Example Feed" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Link != "http://example.org/" {
		t.Errorf("Incorrect site URL, got: %s", feed.Link)
	}

	if feed.Image.URL != "http://example.org/icon.png" {
		t.Errorf("Incorrect image URL, got: %s", feed.Image.URL)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://example.org/2003/12/13/atom03" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}

	if feed.Items[0].Title != "Atom-Powered Robots Run Amok" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__FeedWithoutTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<link rel="alternate" type="text/html" href="https://example.org/"/>
			<link rel="self" type="application/atom+xml" href="https://example.org/feed"/>
			<updated>2003-12-13T18:30:02Z</updated>
		</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "https://example.org/" {
		t.Errorf("Incorrect feed title, got: %s", feed.Title)
	}
}

func Test_Parse_atom10__EntryWithoutTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <updated>2003-12-13T18:30:02Z</updated>
	  <author>
		<name>John Doe</name>
	  </author>
	  <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>
	  <entry>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "http://example.org/2003/12/13/atom03" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__FeedURL(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link rel="alternate" type="text/html" href="https://example.org/"/>
	  <link rel="self" type="application/atom+xml" href="https://example.org/feed"/>
	  <updated>2003-12-13T18:30:02Z</updated>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Link != "https://example.org/" {
		t.Errorf("Incorrect site URL, got: %s", feed.Link)
	}
}

func Test_Parse_atom10__EntryWithRelativeURL(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title>Test</title>
		<link href="something.html"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "http://example.org/something.html" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}
}

func Test_Parse_atom10__EntryTitleWithWhitespaces(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title>
			Some Title
		</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Some Title" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__EntryTitleWithHTMLAndCDATA(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title type="html"><![CDATA[Test &#8220;Test&#8221;]]></title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Test “Test”" {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__EntryTitleWithHTML(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title type="html">&lt;code&gt;Test&lt;/code&gt; Test</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "<code>Test</code> Test" {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__EntryTitleWithXHTML(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title type="xhtml"><code>Test</code> Test</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "<code>Test</code> Test" {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__EntryTitleWithNumericCharacterReference(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title>&#931; &#xDF;</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != "Σ ß" {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__EntryTitleWithDoubleEncodedEntities(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	  <title>Example Feed</title>
	  <link href="http://example.org/"/>
	  <entry>
		<title>&amp;#39;AT&amp;amp;T&amp;#39;</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	  </entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Title != `'AT&T'` {
		t.Errorf("Incorrect entry title, got: %q", feed.Items[0].Title)
	}
}

func Test_Parse_atom10__InvalidXml(t *testing.T) {
	data := `garbage`
	_, err := Parse([]byte(data))
	if err == nil {
		t.Error("Parse should returns an error")
	}
}

func Test_Parse_atom10__TitleWithSingleQuote(t *testing.T) {
	data := `
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title>' or ’</title>
			<link href="http://example.org/"/>
		</feed>
	`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "' or ’" {
		t.Errorf(`Incorrect title, got: %q`, feed.Title)
	}
}

func Test_Parse_atom10__TitleWithEncodedSingleQuote(t *testing.T) {
	data := `
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title type="html">Test&#39;s Blog</title>
			<link href="http://example.org/"/>
		</feed>
	`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "Test's Blog" {
		t.Errorf(`Incorrect title, got: %q`, feed.Title)
	}
}

func Test_Parse_atom10__TitleWithSingleQuoteAndHTMLType(t *testing.T) {
	data := `
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title type="html">O’Hara</title>
			<link href="http://example.org/"/>
		</feed>
	`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "O’Hara" {
		t.Errorf(`Incorrect title, got: %q`, feed.Title)
	}
}

func Test_Parse_atom10__RepliesLinkRelationWithHTMLType(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom"
			xmlns:thr="http://purl.org/syndication/thread/1.0">
		<id>http://www.example.org/myfeed</id>
		<title>My Example Feed</title>
		<updated>2005-07-28T12:00:00Z</updated>
		<link href="http://www.example.org/myfeed" />
		<author><name>James</name></author>
		<entry>
			<id>tag:items.com,2005:1</id>
			<title>My original entry</title>
			<updated>2006-03-01T12:12:12Z</updated>
			<link href="http://www.example.org/items/1" />
			<link rel="replies"
				type="application/atom+xml"
				href="http://www.example.org/mycommentsfeed.xml"
				thr:count="10" thr:updated="2005-07-28T12:10:00Z" />
			<link rel="replies"
				type="text/html"
				href="http://www.example.org/comments.html"
				thr:count="10" thr:updated="2005-07-28T12:10:00Z" />
			<summary>This is my original entry</summary>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://www.example.org/items/1" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}
}

func Test_Parse_atom10__RepliesLinkRelationWithXHTMLType(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom"
			xmlns:thr="http://purl.org/syndication/thread/1.0">
		<id>http://www.example.org/myfeed</id>
		<title>My Example Feed</title>
		<updated>2005-07-28T12:00:00Z</updated>
		<link href="http://www.example.org/myfeed" />
		<author><name>James</name></author>
		<entry>
			<id>tag:items.com,2005:1</id>
			<title>My original entry</title>
			<updated>2006-03-01T12:12:12Z</updated>
			<link href="http://www.example.org/items/1" />
			<link rel="replies"
				type="application/atom+xml"
				href="http://www.example.org/mycommentsfeed.xml"
				thr:count="10" thr:updated="2005-07-28T12:10:00Z" />
			<link rel="replies"
				type="application/xhtml+xml"
				href="http://www.example.org/comments.xhtml"
				thr:count="10" thr:updated="2005-07-28T12:10:00Z" />
			<summary>This is my original entry</summary>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://www.example.org/items/1" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}
}

func Test_Parse_atom10__RepliesLinkRelationWithNoType(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom"
			xmlns:thr="http://purl.org/syndication/thread/1.0">
		<id>http://www.example.org/myfeed</id>
		<title>My Example Feed</title>
		<updated>2005-07-28T12:00:00Z</updated>
		<link href="http://www.example.org/myfeed" />
		<author><name>James</name></author>
		<entry>
			<id>tag:items.com,2005:1</id>
			<title>My original entry</title>
			<updated>2006-03-01T12:12:12Z</updated>
			<link href="http://www.example.org/items/1" />
			<link rel="replies"
				href="http://www.example.org/mycommentsfeed.xml"
				thr:count="10" thr:updated="2005-07-28T12:10:00Z" />
			<summary>This is my original entry</summary>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of items, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://www.example.org/items/1" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}
}
