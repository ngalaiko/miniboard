package items

import "time"

// Item is the item model.
type Item struct {
	ID                string
	URL               string
	Title             string
	SubscriptionID    string
	Created           *time.Time
	Summary           *string
	SubscriptionTitle *string
	SubscriptionIcon  *string
}
