package items

import "time"

// Item is the item model.
type Item struct {
	ID             string    `json:"id"`
	UserID         string    `json:"-"`
	URL            string    `json:"url"`
	Title          string    `json:"title"`
	SubscriptionID string    `json:"subscription_id"`
	Created        time.Time `json:"created"`
}
