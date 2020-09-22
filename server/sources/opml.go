package sources

import (
	"encoding/xml"
	"fmt"

	sources "github.com/ngalaiko/miniboard/server/genproto/sources/v1"
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

func parseOPML(data []byte) ([]*sources.Source, error) {
	parsed := &opml{}
	if err := xml.Unmarshal(data, parsed); err != nil {
		return nil, fmt.Errorf("failed to parse opml: %w", err)
	}

	ss := make([]*sources.Source, 0, 32)
	for _, c := range parsed.Categories {
		for _, f := range c.Feeds {
			ss = append(ss, &sources.Source{
				Url: f.URL,
			})
		}
	}

	return ss, nil
}
