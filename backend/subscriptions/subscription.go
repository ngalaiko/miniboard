package subscriptions

import (
	"fmt"
	"strings"
	"time"
)

// Subscription is the subscription model.
type Subscription struct {
	ID      string       `json:"id"`
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
