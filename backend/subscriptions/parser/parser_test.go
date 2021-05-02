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

func Test_Parser(t *testing.T) {
	t.Parallel()

	//nolint: paralleltest
	for _, ext := range []string{"rss", "rdf"} {
		prefix := "./testdata/" + ext
		files, _ := filepath.Glob(prefix + "/*.xml")
		for _, f := range files {
			f := f
			t.Run(ext+f, func(t *testing.T) {
				t.Parallel()

				base := filepath.Base(f)
				name := strings.TrimSuffix(base, filepath.Ext(base))

				// Get actual source feed
				ff := fmt.Sprintf(prefix+"/%s.xml", name)
				f, _ := ioutil.ReadFile(ff)

				// Parse actual feed
				actual, err := Parse(f, &testLogger{})
				if err != nil {
					t.Error(err)
				}

				// Get json encoded expected feed result
				ef := fmt.Sprintf(prefix+"/%s.json", name)

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
}

type testLogger struct{}

func (tl *testLogger) Error(string, ...interface{}) {
}
