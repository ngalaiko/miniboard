package opml

import (
	"encoding/xml"
	"fmt"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Tags    []*Tag   `xml:"body>outline"`
}

type Tag struct {
	Title string  `xml:"title,attr"`
	Feeds []*Feed `xml:"outline"`
}

type Feed struct {
	Title string `xml:"title,attr"`
	URL   string `xml:"xmlUrl,attr"`
}

func Parse(data []byte) (*OPML, error) {
	parsed := &OPML{}
	if err := xml.Unmarshal(data, parsed); err != nil {
		return nil, fmt.Errorf("failed to parse opml: %w", err)
	}
	return parsed, nil
}
