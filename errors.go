package main

import "net/http"

var ErrServerShutdown = &HTTPError{
	StatusCode: http.StatusServiceUnavailable,
	Msg:        `Shutting Down. No new connections accepted`,
}
