package parser

import (
	"strings"
	"testing"
	"time"
)

func Test_Parse_rdf__RDFSample(t *testing.T) {
	data := `
	<?xml version="1.0"?>
	<rdf:RDF
	  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	  xmlns="http://purl.org/rss/1.0/"
	>
	  <channel rdf:about="http://www.xml.com/xml/news.rss">
		<title>XML.com</title>
		<link>http://xml.com/pub</link>
		<description>
		  XML.com features a rich mix of information and services
		  for the XML community.
		</description>
		<image rdf:resource="http://xml.com/universal/images/xml_tiny.gif" />
		<items>
		  <rdf:Seq>
			<rdf:li resource="http://xml.com/pub/2000/08/09/xslt/xslt.html" />
			<rdf:li resource="http://xml.com/pub/2000/08/09/rdfdb/index.html" />
		  </rdf:Seq>
		</items>
		<textinput rdf:resource="http://search.xml.com" />
	  </channel>
	  <image rdf:about="http://xml.com/universal/images/xml_tiny.gif">
		<title>XML.com</title>
		<link>http://www.xml.com</link>
		<url>http://xml.com/universal/images/xml_tiny.gif</url>
	  </image>
	  <item rdf:about="http://xml.com/pub/2000/08/09/xslt/xslt.html">
		<title>Processing Inclusions with XSLT</title>
		<link>http://xml.com/pub/2000/08/09/xslt/xslt.html</link>
		<description>
		 Processing document inclusions with general XML tools can be
		 problematic. This article proposes a way of preserving inclusion
		 information through SAX-based processing.
		</description>
	  </item>
	  <item rdf:about="http://xml.com/pub/2000/08/09/rdfdb/index.html">
		<title>Putting RDF to Work</title>
		<link>http://xml.com/pub/2000/08/09/rdfdb/index.html</link>
		<description>
		 Tool and API support for the Resource Description Framework
		 is slowly coming of age. Edd Dumbill takes a look at RDFDB,
		 one of the most exciting new RDF toolkits.
		</description>
	  </item>
	  <textinput rdf:about="http://search.xml.com">
		<title>Search XML.com</title>
		<description>Search XML.com's XML collection</description>
		<name>s</name>
		<link>http://search.xml.com</link>
	  </textinput>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "XML.com" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Link != "http://xml.com/pub" {
		t.Errorf("Incorrect feed URL, got: %s", feed.Link)
	}

	if feed.Image.URL != "http://xml.com/universal/images/xml_tiny.gif" {
		t.Errorf("Inconnect feed image, got: %s", feed.Image.URL)
	}

	if len(feed.Items) != 2 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}

	if strings.HasSuffix(feed.Items[1].Content, "Tool and API support") {
		t.Errorf("Incorrect entry content, got: %s", feed.Items[0].Content)
	}

	if feed.Items[1].Link != "http://xml.com/pub/2000/08/09/rdfdb/index.html" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}

	if feed.Items[1].Title != "Putting RDF to Work" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}

	if feed.Items[1].Date != nil {
		t.Errorf("Entry date should not empty")
	}
}

func Test_Parse_rdf__RDFSampleWithDublinCore(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF
	  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	  xmlns:dc="http://purl.org/dc/elements/1.1/"
	  xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
	  xmlns:co="http://purl.org/rss/1.0/modules/company/"
	  xmlns:ti="http://purl.org/rss/1.0/modules/textinput/"
	  xmlns="http://purl.org/rss/1.0/"
	>
	  <channel rdf:about="http://meerkat.oreillynet.com/?_fl=rss1.0">
		<title>Meerkat</title>
		<link>http://meerkat.oreillynet.com</link>
		<description>Meerkat: An Open Wire Service</description>
		<dc:publisher>The O'Reilly Network</dc:publisher>
		<dc:creator>Rael Dornfest (mailto:rael@oreilly.com)</dc:creator>
		<dc:rights>Copyright &#169; 2000 O'Reilly &amp; Associates, Inc.</dc:rights>
		<dc:date>2000-01-01T12:00+00:00</dc:date>
		<sy:updatePeriod>hourly</sy:updatePeriod>
		<sy:updateFrequency>2</sy:updateFrequency>
		<sy:updateBase>2000-01-01T12:00+00:00</sy:updateBase>
		<image rdf:resource="http://meerkat.oreillynet.com/icons/meerkat-powered.jpg" />
		<items>
		  <rdf:Seq>
			<rdf:li resource="http://c.moreover.com/click/here.pl?r123" />
		  </rdf:Seq>
		</items>
		<textinput rdf:resource="http://meerkat.oreillynet.com" />
	  </channel>
	  <image rdf:about="http://meerkat.oreillynet.com/icons/meerkat-powered.jpg">
		<title>Meerkat Powered!</title>
		<url>http://meerkat.oreillynet.com/icons/meerkat-powered.jpg</url>
		<link>http://meerkat.oreillynet.com</link>
	  </image>
	  <item rdf:about="http://c.moreover.com/click/here.pl?r123">
		<title>XML: A Disruptive Technology</title>
		<link>http://c.moreover.com/click/here.pl?r123</link>
		<dc:description>
		  XML is placing increasingly heavy loads on the existing technical
		  infrastructure of the Internet.
		</dc:description>
		<dc:publisher>The O'Reilly Network</dc:publisher>
		<dc:creator>Simon St.Laurent (mailto:simonstl@simonstl.com)</dc:creator>
		<dc:rights>Copyright &#169; 2000 O'Reilly &amp; Associates, Inc.</dc:rights>
		<dc:subject>XML</dc:subject>
		<co:name>XML.com</co:name>
		<co:market>NASDAQ</co:market>
		<co:symbol>XML</co:symbol>
	  </item>
	  <textinput rdf:about="http://meerkat.oreillynet.com">
		<title>Search Meerkat</title>
		<description>Search Meerkat's RSS Database...</description>
		<name>s</name>
		<link>http://meerkat.oreillynet.com/</link>
		<ti:function>search</ti:function>
		<ti:inputType>regex</ti:inputType>
	  </textinput>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "Meerkat" {
		t.Errorf("Incorrect title, got: %s", feed.Title)
	}

	if feed.Link != "http://meerkat.oreillynet.com" {
		t.Errorf("Incorrect feed URL, got: %s", feed.Link)
	}

	if len(feed.Items) != 1 {
		t.Errorf("Incorrect number of entries, got: %d", len(feed.Items))
	}

	if feed.Items[0].Link != "http://c.moreover.com/click/here.pl?r123" {
		t.Errorf("Incorrect entry URL, got: %s", feed.Items[0].Link)
	}

	if feed.Items[0].Title != "XML: A Disruptive Technology" {
		t.Errorf("Incorrect entry title, got: %s", feed.Items[0].Title)
	}

	if strings.HasSuffix(feed.Items[0].Content, "XML is placing increasingly") {
		t.Errorf("Incorrect entry content, got: %s", feed.Items[0].Content)
	}
}

func Test_Parse_rdf__ItemRelativeURL(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/">
	  <channel>
			<title>Example</title>
			<link>http://example.org</link>
	  </channel>
	  <item>
			<title>Title</title>
			<description>Test</description>
			<link>something.html</link>
	  </item>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "http://example.org/something.html" {
		t.Errorf("Incorrect entry url, got: %s", feed.Items[0].Link)
	}
}

func Test_Parse_rdf__ItemWithoutLink(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF
	  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	  xmlns="http://purl.org/rss/1.0/"
	>
	  <channel rdf:about="http://meerkat.oreillynet.com/?_fl=rss1.0">
		<title>Meerkat</title>
		<link>http://meerkat.oreillynet.com</link>
	  </channel>
	  <item rdf:about="http://c.moreover.com/click/here.pl?r123">
		<title>Title</title>
		<description>Test</description>
	  </item>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Link != "http://meerkat.oreillynet.com" {
		t.Errorf("Incorrect entry url, got: %s", feed.Items[0].Link)
	}
}

func Test_Parse_rdf__InvalidXml(t *testing.T) {
	data := `garbage`
	_, err := Parse([]byte(data), &testLogger{})
	if err == nil {
		t.Fatal("Parse should returns an error")
	}
}

func Test_Parse_rdf__FeedWithURLWrappedInSpaces(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF xmlns:admin="http://webns.net/mvcb/" xmlns="http://purl.org/rss/1.0/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:prism="http://purl.org/rss/1.0/modules/prism/" xmlns:taxo="http://purl.org/rss/1.0/modules/taxonomy/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:syn="http://purl.org/rss/1.0/modules/syndication/">
	<channel rdf:about="http://biorxiv.org">
		<title>bioRxiv Subject Collection: Bioengineering</title>
		<link>http://biorxiv.org</link>
		<description>
			This feed contains articles for bioRxiv Subject Collection "Bioengineering"
		</description>
		<items>
			<rdf:Seq>
				<rdf:li rdf:resource="http://biorxiv.org/cgi/content/short/857789v1?rss=1"/>
			</rdf:Seq>
		</items>
		<prism:eIssn/>
		<prism:publicationName>bioRxiv</prism:publicationName>
		<prism:issn/>
		<image rdf:resource=""/>
	</channel>
	<image rdf:about="">
		<title>bioRxiv</title>
		<url/>
		<link>http://biorxiv.org</link>
	</image>
	<item rdf:about="http://biorxiv.org/cgi/content/short/857789v1?rss=1">
		<title>
			<![CDATA[
			Microscale Collagen and Fibroblast Interactions Enhance Primary Human Hepatocyte Functions in 3-Dimensional Models
			]]>
		</title>
		<link>
			http://biorxiv.org/cgi/content/short/857789v1?rss=1
		</link>
		<description><![CDATA[
		Human liver models that are 3-dimensional (3D) in architecture are proving to be indispensable for diverse applications, including compound metabolism and toxicity screening during preclinical drug development, to model human liver diseases for the discovery of novel therapeutics, and for cell-based therapies in the clinic; however, further development of such models is needed to maintain high levels of primary human hepatocyte (PHH) functions for weeks to months in vitro. Therefore, here we determined how microscale 3D collagen-I presentation and fibroblast interaction could affect the long-term functions of PHHs. High-throughput droplet microfluidics was utilized to rapidly generate reproducibly-sized (~300 micron diameter) microtissues containing PHHs encapsulated in collagen-I +/- supportive fibroblasts, namely 3T3-J2 murine embryonic fibroblasts or primary human hepatic stellate cells (HSCs); self-assembled spheroids and bulk collagen gels (macrogels) containing PHHs served as gold-standard controls. Hepatic functions (e.g. albumin and cytochrome-P450 or CYP activities) and gene expression were subsequently measured for up to 6 weeks. We found that collagen-based 3D microtissues rescued PHH functions within static multi-well plates at 2- to 30-fold higher levels than self-assembled spheroids or macrogels. Further coating of PHH microtissues with 3T3-J2s led to higher hepatic functions than when the two cell types were either coencapsulated together or when HSCs were used for the coating instead. Additionally, the 3T3-J2-coated PHH microtissues displayed 6+ weeks of relatively stable hepatic gene expression and function at levels similar to freshly thawed PHHs. Lastly, microtissues responded in a clinically-relevant manner to drug-mediated CYP induction or hepatotoxicity. In conclusion, fibroblast-coated collagen microtissues containing PHHs display hepatic functions for 6+ weeks without any fluid perfusion at higher levels than spheroids and macrogels, and such microtissues can be used to assess drug-mediated CYP induction and hepatotoxicity. Ultimately, microtissues may find broader utility for modeling liver diseases and as building blocks for cell-based therapies.
		]]></description>
		<dc:creator><![CDATA[ Kukla, D., Crampton, A., Wood, D., Khetani, S. ]]></dc:creator>
		<dc:date>2019-11-29</dc:date>
		<dc:identifier>doi:10.1101/857789</dc:identifier>
		<dc:title><![CDATA[Microscale Collagen and Fibroblast Interactions Enhance Primary Human Hepatocyte Functions in 3-Dimensional Models]]></dc:title>
		<dc:publisher>Cold Spring Harbor Laboratory</dc:publisher>
		<prism:publicationDate>2019-11-29</prism:publicationDate>
		<prism:section></prism:section>
	</item>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Fatalf(`Unexpected number of entries, got %d`, len(feed.Items))
	}

	if feed.Items[0].Link != `http://biorxiv.org/cgi/content/short/857789v1?rss=1` {
		t.Errorf(`Unexpected entry URL, got %q`, feed.Items[0].Link)
	}
}

func Test_Parse_rdf_ItemWithDublicCoreDate(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:slash="http://purl.org/rss/1.0/modules/slash/">
	  <channel>
			<title>Example</title>
			<link>http://example.org</link>
	  </channel>
	  <item>
			<title>Title</title>
			<description>Test</description>
			<link>http://example.org/test.html</link>
			<dc:creator>Tester</dc:creator>
			<dc:date>2018-04-10T05:00:00+00:00</dc:date>
	  </item>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	expectedDate := time.Date(2018, time.April, 10, 5, 0, 0, 0, time.UTC)
	if !feed.Items[0].Date.Equal(expectedDate) {
		t.Errorf("Incorrect entry date, got: %v, want: %v", feed.Items[0].Date, expectedDate)
	}
}

func Test_Parse_rdf_ItemWithoutDate(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/">
	  <channel>
			<title>Example</title>
			<link>http://example.org</link>
	  </channel>
	  <item>
			<title>Title</title>
			<description>Test</description>
			<link>http://example.org/test.html</link>
	  </item>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if feed.Items[0].Date != nil {
		t.Errorf("Incorrect entry date, got: %v", feed.Items[0].Date)
	}
}

func Test_Parse_rdf_WithContentEncoded(t *testing.T) {
	data := `<?xml version="1.0" encoding="utf-8"?>
	<rdf:RDF
		xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
		xmlns="http://purl.org/rss/1.0/"
		xmlns:content="http://purl.org/rss/1.0/modules/content/">
		<channel>
			<title>Example Feed</title>
			<link>http://example.org/</link>
		</channel>
		<item>
			<title>Item Title</title>
			<link>http://example.org/</link>
			<content:encoded><![CDATA[<p>Test</p>]]></content:encoded>
		</item>
	</rdf:RDF>`

	feed, err := Parse([]byte(data), &testLogger{})
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Items) != 1 {
		t.Fatalf(`Unexpected number of entries, got %d`, len(feed.Items))
	}

	expected := `<p>Test</p>`
	result := feed.Items[0].Content
	if result != expected {
		t.Errorf(`Unexpected entry URL, got %q instead of %q`, result, expected)
	}
}
