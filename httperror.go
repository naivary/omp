package main

// +openapi:schema:title="http error response"
type HTTPError struct {
	StatusCode int
	Msg        string
	RequestID  string
	SpanID     string
}

func NewHTTPError(msg string, status int) *HTTPError {
	return &HTTPError{
		Msg:        msg,
		StatusCode: status,
	}
}

func (h *HTTPError) Error() string {
	return h.Msg
}
