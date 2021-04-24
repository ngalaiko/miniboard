package sockets

type position string

const (
	Afterbegin position = "afterbegin"
	Beforeend  position = "beforeend"
)

type Response struct {
	ID     uint     `json:"id"`
	Error  string   `json:"error,omitempty"`
	Target string   `json:"target"`
	HTML   string   `json:"html"`
	Insert position `json:"insert"`
	Reset  bool     `json:"reset"`
}

func Error(req *Request, err error) *Response {
	return &Response{
		ID:    req.ID,
		Error: err.Error(),
	}
}
