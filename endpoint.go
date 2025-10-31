package main

import (
	"net/http"

	"github.com/naivary/omp/openapi"
)

var _ http.Handler = (*Endpoint)(nil)

// Endpoint implements http.Handler and allows to customize the used ErroHandler
// for a given HandlerFuncErr. Furhtermore it allows you to define OpenAPI
// documentation metadata to automatically generate OpenAPI specs.
//
// TODO: Callbacks
// TODO: Servers
type Endpoint struct {
	// Handler to handle the incoming request.
	Handler HandlerFuncErr
	// Error is the handler used if the returned error of the Handler is
	// non-nil.
	Error ErrorHandler

	Pattern string
	// Summary included in the OpenAPI spec.
	Summary string
	// Description included in the OpenAPI spec. By default it will use the doc
	// string of your `Handler`.
	Description string
	// Tags included in the OpenAPI spec.
	Tags []string
	// OperationID is used to generate the OpenAPI spec. By default it will use
	// the name of your `Handler` function.
	OperationID string
	// Deprecated flags an endpoint as deprecated in the OpenAPI documentation.
	Deprecated bool

	// Parameters defines teh required and optional parameters for this
	// endpoint. Parameters might be in Cookies, Headers, Queries and Path. Path
	// parameters cannot be generated automatically without breaking the
	// standard patter notation of the http package. Thats why you have to
	// define them by yourself using the helper functions.
	Parameters []*openapi.Parameter

	// RequestBody of the Endpoint.
	RequestBody *openapi.RequestBody
	Responses   map[string]*openapi.Response
	Security    openapi.SecurityRequirement
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.Handler(w, r)
	if err == nil {
		return
	}
	e.Error.ServeError(w, r, err)
}

// ErrorHandler is a handler used to handle errors which might occur in a
// request handler.
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
		if !isHTTPErr {
			httpErr = NewHTTPError(http.StatusInternalServerError, "internal error: %v", err.Error())
		}
		encode(w, r, httpErr.StatusCode, &httpErr)
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
