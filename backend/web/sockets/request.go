package sockets

type event string

const (
	tagToggled    event = "tag:toggled"
	itemSelected  event = "item:selected"
	itemsLoad     event = "items:load"
	itemsLoadmore event = "items:loadmore"
)

type request struct {
	ID     uint              `json:"id"`
	Event  event             `json:"event"`
	Params map[string]string `json:"params"`
}
