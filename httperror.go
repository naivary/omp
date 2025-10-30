package main

import "fmt"

// +openapi:schema:title="http error response"
type HTTPError struct {
	StatusCode int
	Msg        string
	RequestID  string
	SpanID     string
}

func NewHTTPError(status int, format string, args ...any) *HTTPError {
	return &HTTPError{
		Msg:        fmt.Sprintf(format, args...),
		StatusCode: status,
	}
}

func (h *HTTPError) Error() string {
	return h.Msg
}
