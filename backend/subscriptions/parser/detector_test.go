package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_detctType(t *testing.T) {
	testCases := []struct {
		file string
		out  feedType
	}{
		{"atom03_feed.xml", feedTypeAtom},
		{"atom10_feed.xml", feedTypeAtom},
		{"rss_feed.xml", feedTypeRSS},
		{"rss_feed_bom.xml", feedTypeRSS},
		{"rss_feed_leading_spaces.xml", feedTypeRSS},
		{"rdf_feed.xml", feedTypeRSS},
		{"unknown_feed.xml", feedTypeUnknown},
		{"empty_feed.xml", feedTypeUnknown},
		{"json_feed.json", feedTypeJSON},
	}

	for _, testCase := range testCases {
		t.Run(testCase.file, func(t *testing.T) {
			path := fmt.Sprintf("testdata/%s", testCase.file)
			f, _ := ioutil.ReadFile(path)

			actual := detectType(f)
			if testCase.out != actual {
				t.Errorf("expected %v, got %v", actual, testCase.out)
			}
		})
	}
}
