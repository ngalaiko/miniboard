package sockets

type Request struct {
	ID     uint              `json:"id"`
	Event  string            `json:"event"`
	Params map[string]string `json:"params"`
}
