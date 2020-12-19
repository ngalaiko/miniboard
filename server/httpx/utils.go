package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON marshals objet as json and writes it to the response body.
func JSON(w http.ResponseWriter, response interface{}) {
	body, _ := json.Marshal(response)
	size, _ := w.Write(body)

	w.Header().Add("Content-Length", fmt.Sprint(size))
	w.Header().Add("Content-Type", "application/json")
}

type errorMessage struct {
	Error string `json:"error"`
}

// Error writes error response.
func Error(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	JSON(w, &errorMessage{Error: err.Error()})
}

// InternalError responsds with unknown internal error.
func InternalError(w http.ResponseWriter) {
	Error(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
}
