package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Parser__rss(t *testing.T) {
	files, _ := filepath.Glob("./testdata/rss/*.xml")
	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			base := filepath.Base(f)
			name := strings.TrimSuffix(base, filepath.Ext(base))

			// Get actual source feed
			ff := fmt.Sprintf("./testdata/rss/%s.xml", name)
			f, _ := ioutil.ReadFile(ff)

			// Parse actual feed
			actual, err := Parse(f)
			if err != nil {
				t.Error(err)
			}

			// Get json encoded expected feed result
			ef := fmt.Sprintf("./testdata/rss/%s.json", name)
			e, _ := ioutil.ReadFile(ef)

			// Unmarshal expected feed
			expected := &Feed{}

			if err := json.Unmarshal(e, &expected); err != nil {
				t.Error(err)
			}

			if !cmp.Equal(expected, actual) {
				t.Error(cmp.Diff(expected, actual))
			}
		})
	}
}
