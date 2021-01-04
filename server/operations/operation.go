package operations

import "github.com/google/uuid"

// Error contains operations error result.
type Error struct {
	Message string `json:"message"`
}

// Result defined operation's result.
type Result struct {
	Error    *Error      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// Operation represents a longrunning operation.
type Operation struct {
	ID     string  `json:"id"`
	UserID string  `json:"-"`
	Done   bool    `json:"done"`
	Result *Result `json:"result,omitempty"`
}

// New creates new operation.
func New(userID string) *Operation {
	return &Operation{
		ID:     uuid.New().String(),
		UserID: userID,
	}
}

// Success marks operation as successful.
func (o *Operation) Success(response interface{}) {
	o.Result = &Result{
		Response: new(interface{}),
	}
	o.Result.Response = response
	o.Done = true
}

// Error marks operation as failed.
func (o *Operation) Error(err string) {
	o.Result = &Result{
		Error: &Error{
			Message: err,
		},
	}
	o.Done = true
}
