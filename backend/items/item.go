package items

import "time"

// UserItem is the item model with user id.
type UserItem struct {
	Item
	UserID string `json:"-"`
}

// Item is the item model.
type Item struct {
	ID             string     `json:"id"`
	URL            string     `json:"url"`
	Title          string     `json:"title"`
	SubscriptionID string     `json:"subscription_id"`
	Created        *time.Time `json:"created,omitempty"`
	Summary        *string    `json:"summary,omitempty"`
}
