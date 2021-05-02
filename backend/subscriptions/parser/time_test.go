package parser

import (
	"testing"
	"time"
)

func Test_parseDateTime(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		In  string
		Out time.Time
	}{
		{
			In:  "Tue, 17 Jan 2017 21:19:47 +0000",
			Out: time.Date(2017, time.January, 17, 21, 19, 47, 0, time.UTC),
		},
		{
			In:  "Wed, 1 Feb 2017 10:18:28 +0000",
			Out: time.Date(2017, time.February, 1, 10, 18, 28, 0, time.UTC),
		},
		{
			In:  "Thu, 2 Feb 2017 10:40:19 +0000",
			Out: time.Date(2017, time.February, 2, 10, 40, 19, 0, time.UTC),
		},
		{
			In:  "чт, 26 дек. 2019 22:21:00 +0000",
			Out: time.Date(2019, time.December, 26, 22, 21, 0, 0, time.UTC),
		},
		{
			In:  "2019-06-05T00:00",
			Out: time.Date(2019, time.June, 5, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			t.Parallel()

			dt, err := parseDateTime(tc.In)
			if err != nil {
				t.Error(err)
				return
			}
			if !dt.Equal(tc.Out) {
				t.Errorf("expected '%s', got '%s'", tc.Out, *dt)
			}
		})
	}
}
