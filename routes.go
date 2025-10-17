package main

import "net/http"

func addRoutes(mux *http.ServeMux) {
	mux.Handle("GET /livez", livez())
	mux.Handle("GET /readyz", readyz())
}
