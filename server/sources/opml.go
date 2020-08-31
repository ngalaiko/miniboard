package sources

import (
	"encoding/xml"
	"fmt"
)

type opml struct {
	XMLName    xml.Name    `xml:"opml"`
	Categories []*category `xml:"body>outline"`
}

type category struct {
	Title string  `xml:"title,attr"`
	Feeds []*feed `xml:"outline"`
}

type feed struct {
	Title string `xml:"title,attr"`
	URL   string `xml:"xmlUrl,attr"`
}

func parseOPML(data []byte) ([]*Source, error) {
	parsed := &opml{}
	if err := xml.Unmarshal(data, parsed); err != nil {
		return nil, fmt.Errorf("failed to parse opml: %w", err)
	}

	ss := make([]*Source, 0, 32)
	for _, c := range parsed.Categories {
		for _, f := range c.Feeds {
			ss = append(ss, &Source{
				Url: f.URL,
			})
		}
	}

	return ss, nil
}
