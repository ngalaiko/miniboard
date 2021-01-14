package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorLogger interface {
	Error(string, ...interface{})
}

// JSON marshals objet as json and writes it to the response body.
func JSON(w http.ResponseWriter, logger errorLogger, response interface{}) {
	body, err := json.Marshal(response)
	if err != nil {
		logger.Error("failed to marshal response: %s", err)
		InternalError(w, logger)
		return
	}

	size, err := w.Write(body)
	if err != nil {
		logger.Error("failed to write response: %s", err)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprint(size))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

type errorMessage struct {
	Message string `json:"message"`
}

// Error writes error response.
func Error(w http.ResponseWriter, logger errorLogger, err error, code int) {
	w.WriteHeader(code)
	JSON(w, logger, &errorMessage{Message: err.Error()})
}

// InternalError responsds with unknown internal error.
func InternalError(w http.ResponseWriter, logger errorLogger) {
	Error(w, logger, fmt.Errorf("internal server error"), http.StatusInternalServerError)
}
