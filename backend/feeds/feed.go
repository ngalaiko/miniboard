package feeds

import (
	"time"
)

// Feed is the feed model.
type Feed struct {
	ID      string     `json:"id"`
	UserID  string     `json:"-"`
	URL     string     `json:"url"`
	Title   string     `json:"title"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
}
