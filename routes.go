package main

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func addRoutes(mux *http.ServeMux, pg *pgx.Conn) {
	// system
	mux.Handle("GET /livez", livez())
	mux.Handle("GET /readyz", readyz())
	// metrics management
	mux.Handle("POST /metrics/definition", CreateMetricDefinition(pg))
}
