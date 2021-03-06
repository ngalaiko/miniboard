package sockets

type position string

const (
	beforebegin position = "beforebegin"
	afterbegin  position = "afterbegin"
	beforeend   position = "beforeend"
	afterend    position = "afterend"
)

type response struct {
	ID     uint     `json:"id"`
	Error  string   `json:"error,omitempty"`
	Target string   `json:"target"`
	HTML   string   `json:"html"`
	Insert position `json:"insert"`
	Reset  bool     `json:"reset"`
}

func errResponse(req *request, err error) *response {
	return &response{
		ID:    req.ID,
		Error: err.Error(),
	}
}
