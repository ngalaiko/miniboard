package sockets

type event string

const (
	subscriptionsCreated event = "subscriptions:create"
	itemsSelect          event = "items:select"
	itemsLoad            event = "items:load"
	itemsLoadmore        event = "items:loadmore"
)

type request struct {
	ID     uint              `json:"id"`
	Event  event             `json:"event"`
	Params map[string]string `json:"params"`
}
