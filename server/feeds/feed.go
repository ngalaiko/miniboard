package feeds

import (
	"net/url"
	"time"
)

// Feed is the feed model.
type Feed struct {
	ID      string     `json:"id"`
	UserID  string     `json:"user_id"`
	URL     *url.URL   `json:"url"`
	Title   string     `json:"title"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
}
