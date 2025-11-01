package main

import "net/http"

func Livez() *Endpoint {
	return &Endpoint{
		Handler: livez(),
		Error:   defaultErrorHandler(),
	}
}

func livez() HandlerFuncErr {
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
