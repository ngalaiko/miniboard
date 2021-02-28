package parser

import (
	"bytes"
	"encoding/xml"
)

type feedType int8

const (
	feedTypeUnknown feedType = iota
	feedTypeRSS
	feedTypeRDF
	feedTypeAtom03
	feedTypeAtom10
	feedTypeJSON
)

// cc: https://github.com/mmcdole/gofeed/blob/96998c2d6be59bc75c2d86863bedd2fe5a09a20d/detector.go

func detectType(data []byte) feedType {
	buffer := bytes.NewBuffer(data)
	var firstChar byte
loop:
	for {
		ch, err := buffer.ReadByte()
		if err != nil {
			return feedTypeUnknown
		}
		// ignore leading whitespace & byte order marks
		switch ch {
		case ' ', '\r', '\n', '\t':
		case 0xFE, 0xFF, 0x00, 0xEF, 0xBB, 0xBF: // utf 8-16-32 bom
		default:
			firstChar = ch
			if err := buffer.UnreadByte(); err != nil {
				return feedTypeUnknown
			}
			break loop
		}
	}

	switch firstChar {
	case '{':
		return feedTypeJSON
	case '<':
		var rss struct {
			XMLName xml.Name `xml:"rss"`
		}

		var rdf struct {
			XMLName xml.Name `xml:"RDF"`
		}

		var atom struct {
			XMLName xml.Name `xml:"feed"`
			Version string   `xml:"version,attr"`
		}

		switch {
		case xml.Unmarshal(data, &rss) == nil:
			return feedTypeRSS
		case xml.Unmarshal(data, &rdf) == nil:
			return feedTypeRDF
		case xml.Unmarshal(data, &atom) == nil:
			if atom.Version == "0.3" {
				return feedTypeAtom03
			}
			return feedTypeAtom10
		default:
			return feedTypeUnknown
		}
	default:
		return feedTypeUnknown
	}
}
