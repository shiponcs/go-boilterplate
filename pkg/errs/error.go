package errs

import (
	"fmt"
	"maps"
	"net/http"
)

// StatusCoder reports the HTTP status code for an error.
type StatusCoder interface {
	StatusCode() int
}

// HTTPFormatter renders an error into a JSON-serializable body.
type HTTPFormatter interface {
	HTTPFormat() map[string]any
}

// Error is a client-facing error carrying an internal log message, a client
// message, and an HTTP status code. The serve helpers inspect it to render a
// consistent error response.
type Error struct {
	LogMsg    string
	ClientMsg any
	Code      int
}

func (e *Error) Error() string {
	return fmt.Sprintf("log msg: %s, status code: %d", e.LogMsg, e.Code)
}

func (e *Error) HTTPFormat() map[string]any {
	m := map[string]any{"status": false}
	switch msg := e.ClientMsg.(type) {
	case map[string]any:
		maps.Copy(m, msg)
	default:
		m["error"] = msg
	}
	return m
}

func (e *Error) StatusCode() int {
	return e.Code
}

func NewClientError(logMsg string, clientMsg any, code int) error {
	return &Error{LogMsg: logMsg, ClientMsg: clientMsg, Code: code}
}

func NewBadReq(logMsg string, clientMsg any) error {
	return NewClientError(logMsg, clientMsg, http.StatusBadRequest)
}

func NewUnauthorized(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusUnauthorized)
}

func NewForbidden(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusForbidden)
}

func NewConflict(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusConflict)
}

func NewNotFound(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusNotFound)
}

func NewUnprocessable(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusUnprocessableEntity)
}

func NewInternalServer(logMsg, clientMsg string) error {
	if clientMsg == "" {
		clientMsg = "internal server error"
	}
	return NewClientError(logMsg, clientMsg, http.StatusInternalServerError)
}
