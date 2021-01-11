package operations

import (
	"context"

	"github.com/google/uuid"
)

// Error contains operations error result.
type Error struct {
	Message string `json:"message"`
}

// Result defined operation's result.
type Result struct {
	Error    *Error      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// Task is a single long running operation.
type Task func(context.Context, *Operation, chan<- *Operation) error

// Operation represents a longrunning operation.
type Operation struct {
	ID     string  `json:"id"`
	Done   bool    `json:"done"`
	UserID string  `json:"-"`
	Result *Result `json:"result,omitempty"`

	task Task
}

// New creates new operation.
func New(userID string) *Operation {
	return &Operation{
		ID:     uuid.New().String(),
		UserID: userID,
	}
}

func (o *Operation) copy() *Operation {
	return &Operation{
		ID:     o.ID,
		Done:   o.Done,
		UserID: o.UserID,
		Result: o.Result,
		task:   o.task,
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
