package main

import (
	"net/http"

	"github.com/naivary/omp/api/system"
	"github.com/naivary/omp/openapi"
)

func Readyz() *Endpoint {
	return &Endpoint{
		Handler:     readyz(),
		Error:       defaultErrorHandler(),
		Summary:     "Reports wheter the API is ready to accept requests",
		Tags:        openapi.Tags("system"),
		Pattern:     "GET /readyz",
		OperationID: "readyz",
		RequestBody: nil,
		Parameters: []*openapi.Parameter{
			openapi.NewQueryParam[int](
				"verbose",
				"vebrosity of the response",
				false,
			),
		},
		Responses: map[string]*openapi.Response{
			"200": openapi.NewResponse[system.ReadyzResponse]("API is ready and can accept requests"),
			"503": openapi.NewResponse[HTTPError]("Service is unavailable"),
		},
	}
}

func readyz() HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		if ctx.Err() != nil {
			return ErrServerShutdown
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return nil
	})
}
