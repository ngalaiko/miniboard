package parser

import (
	"testing"
)

func Test_Parse_atom03__Atom03(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed version="0.3" xmlns="http://purl.org/atom/ns#">
		<title>dive into mark</title>
		<link rel="alternate" type="text/html" href="http://diveintomark.org/"/>
		<modified>2003-12-13T18:30:02Z</modified>
		<author><name>Mark Pilgrim</name></author>
		<entry>
			<title>Atom 0.3 snapshot</title>
			<link rel="alternate" type="text/html" href="http://diveintomark.org/2003/12/13/atom03"/>
			<id>tag:diveintomark.org,2003:3.2397</id>
			<issued>2003-12-13T08:29:29-04:00</issued>
			<modified>2003-12-13T18:30:02Z</modified>
			<summary type="text/plain">It&apos;s a test</summary>
			<content type="text/html" mode="escaped"><![CDATA[<p>HTML content</p>]]></content>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "dive into mark" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Link != "http://diveintomark.org/" {
		t.Errorf("Incorrect feed URL, got: %s", feed.Link)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://diveintomark.org/2003/12/13/atom03" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}

	if feed.Items[0].Title != "Atom 0.3 snapshot" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_atom03__Atom03WithoutFeedTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed version="0.3" xmlns="http://purl.org/atom/ns#">
		<link rel="alternate" type="text/html" href="http://diveintomark.org/"/>
		<modified>2003-12-13T18:30:02Z</modified>
		<author><name>Mark Pilgrim</name></author>
		<entry>
			<title>Atom 0.3 snapshot</title>
			<link rel="alternate" type="text/html" href="http://diveintomark.org/2003/12/13/atom03"/>
			<id>tag:diveintomark.org,2003:3.2397</id>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "http://diveintomark.org/" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}
}

func Test_Parse_atom03__Atom03WithoutEntryTitle(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed version="0.3" xmlns="http://purl.org/atom/ns#">
		<title>dive into mark</title>
		<link rel="alternate" type="text/html" href="http://diveintomark.org/"/>
		<modified>2003-12-13T18:30:02Z</modified>
		<author><name>Mark Pilgrim</name></author>
		<entry>
			<link rel="alternate" type="text/html" href="http://diveintomark.org/2003/12/13/atom03"/>
			<id>tag:diveintomark.org,2003:3.2397</id>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}

	if feed.Items[0].Title != "http://diveintomark.org/2003/12/13/atom03" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}
}

func Test_Parse_atom03__Atom03WithSummaryOnly(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed version="0.3" xmlns="http://purl.org/atom/ns#">
		<title>dive into mark</title>
		<link rel="alternate" type="text/html" href="http://diveintomark.org/"/>
		<modified>2003-12-13T18:30:02Z</modified>
		<author><name>Mark Pilgrim</name></author>
		<entry>
			<title>Atom 0.3 snapshot</title>
			<link rel="alternate" type="text/html" href="http://diveintomark.org/2003/12/13/atom03"/>
			<id>tag:diveintomark.org,2003:3.2397</id>
			<issued>2003-12-13T08:29:29-04:00</issued>
			<modified>2003-12-13T18:30:02Z</modified>
			<summary type="text/plain">It&apos;s a test</summary>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}
}

func Test_Parse_atom03__Atom03WithXMLContent(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed version="0.3" xmlns="http://purl.org/atom/ns#">
		<title>dive into mark</title>
		<link rel="alternate" type="text/html" href="http://diveintomark.org/"/>
		<modified>2003-12-13T18:30:02Z</modified>
		<author><name>Mark Pilgrim</name></author>
		<entry>
			<title>Atom 0.3 snapshot</title>
			<link rel="alternate" type="text/html" href="http://diveintomark.org/2003/12/13/atom03"/>
			<id>tag:diveintomark.org,2003:3.2397</id>
			<issued>2003-12-13T08:29:29-04:00</issued>
			<modified>2003-12-13T18:30:02Z</modified>
			<content mode="xml" type="text/html"><p>Some text.</p></content>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}
}

func Test_Parse_atom03__Atom03WithBase64Content(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<feed version="0.3" xmlns="http://purl.org/atom/ns#">
		<title>dive into mark</title>
		<link rel="alternate" type="text/html" href="http://diveintomark.org/"/>
		<modified>2003-12-13T18:30:02Z</modified>
		<author><name>Mark Pilgrim</name></author>
		<entry>
			<title>Atom 0.3 snapshot</title>
			<link rel="alternate" type="text/html" href="http://diveintomark.org/2003/12/13/atom03"/>
			<id>tag:diveintomark.org,2003:3.2397</id>
			<issued>2003-12-13T08:29:29-04:00</issued>
			<modified>2003-12-13T18:30:02Z</modified>
			<content mode="base64" type="text/html">PHA+U29tZSB0ZXh0LjwvcD4=</content>
		</entry>
	</feed>`

	feed, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}
}
