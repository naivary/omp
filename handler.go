package main

import (
	"net/http"
)

var _ http.Handler = (*Endpoint)(nil)

// Endpoint implements http.Handler and allows to customize the used ErroHandler
// for a given HandlerFuncErr.
type Endpoint struct {
	Handler HandlerFuncErr
	Error   ErrorHandler
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.Handler(w, r)
	if err == nil {
		return
	}
	e.Error.ServeError(w, r, err)
}

// ErrorHandler is a handler used to handle errors which might occur in a
// http.Handler.
type ErrorHandler interface {
	ServeError(w http.ResponseWriter, r *http.Request, err error)
}

var _ ErrorHandler = (ErrorHandlerFunc)(nil)

// ErrorHandlerFunc is implementing ErrorHandler
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

func (e ErrorHandlerFunc) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	e(w, r, err)
}

func defaultErrorHandler() ErrorHandler {
	fn := func(w http.ResponseWriter, r *http.Request, err error) {
		httpErr, isHTTPErr := err.(*HTTPError)
		msg := err.Error()
		code := http.StatusInternalServerError
		if isHTTPErr {
			msg = httpErr.Msg
			code = httpErr.StatusCode
		}
		http.Error(w, msg, code)
	}
	return ErrorHandlerFunc(fn)
}

var _ http.Handler = (HandlerFuncErr)(nil)

// HandlerFuncErr is an http.Handler allowing to return an error for idiomatic
// error handling. If a non-nil error is returned it will be handled using the
// `defaultErrorHandler`. If a custom ErrorHandler is needed you should return
// an `Endpoint` with your custom ErrorHandlerFunc.
type HandlerFuncErr func(w http.ResponseWriter, r *http.Request) error

func (h HandlerFuncErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}
	defaultErrorHandler().ServeError(w, r, err)
}
