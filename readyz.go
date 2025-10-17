package main

import "net/http"

func readyz() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		if ctx.Err() != nil {
			return ErrServerShutdown
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return nil
	}
	return HandlerFuncErr(fn)
}
