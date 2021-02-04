package feeds

import (
	"fmt"
	"strings"
	"time"
)

// Feed is the feed model.
type Feed struct {
	ID      string       `json:"id"`
	UserID  string       `json:"-"`
	URL     string       `json:"url"`
	Title   string       `json:"title"`
	Created time.Time    `json:"created"`
	Updated *time.Time   `json:"updated,omitempty"`
	IconURL *string      `json:"icon_url,omitempty"`
	TagIDs  stringsArray `json:"tag_ids"`
}

type stringsArray []string

// Scan implements database Scanner.
func (sa *stringsArray) Scan(value interface{}) error {
	if value == nil {
		*sa = []string{}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("unexpected type: %T", value)
	}

	*sa = strings.Split(str, ",")

	return nil
}
