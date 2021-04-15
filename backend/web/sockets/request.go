package sockets

type event string

const (
	tagToggled           event = "tag:toggled"
	tagSelected          event = "tag:selected"
	subscriptionSelected event = "subscription:selected"
	itemSelected         event = "item:selected"
)

type request struct {
	ID     uint              `json:"id"`
	Event  event             `json:"event"`
	Params map[string]string `json:"params"`
}
