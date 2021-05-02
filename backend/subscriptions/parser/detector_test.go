package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_detectType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		file string
		out  feedType
	}{
		{"atom03_feed.xml", feedTypeAtom03},
		{"atom10_feed.xml", feedTypeAtom10},
		{"rss_feed.xml", feedTypeRSS},
		{"rss_feed_bom.xml", feedTypeRSS},
		{"rss_feed_leading_spaces.xml", feedTypeRSS},
		{"rdf_feed.xml", feedTypeRDF},
		{"unknown_feed.xml", feedTypeUnknown},
		{"empty_feed.xml", feedTypeUnknown},
		{"json_feed.json", feedTypeJSON},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.file, func(t *testing.T) {
			t.Parallel()

			path := fmt.Sprintf("testdata/%s", testCase.file)
			f, _ := ioutil.ReadFile(path)

			actual := detectType(f)
			if testCase.out != actual {
				t.Errorf("expected %v, got %v", actual, testCase.out)
			}
		})
	}
}
