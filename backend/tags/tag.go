package tags

import "time"

// Tag is the tag model.
type Tag struct {
	ID      string    `json:"id"`
	UserID  string    `json:"-"`
	Title   string    `json:"title"`
	Created time.Time `json:"created"`
}
