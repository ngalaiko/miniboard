package sockets

type position string

const (
	afterbegin position = "afterbegin"
	beforeend  position = "beforeend"
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

func resetResponse(req *request, target string) *response {
	return &response{
		ID:     req.ID,
		Target: target,
		Reset:  true,
	}
}
