package sockets

type position string

const (
	afterbegin position = "afterbegin"
)

type response struct {
	ID     uint     `json:"id"`
	Error  string   `json:"error,omitempty"`
	Target string   `json:"target"`
	HTML   string   `json:"html"`
	Insert position `json:"insert"`
	Reset  bool     `json:"reset"`
}
