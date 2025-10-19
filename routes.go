package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func addRoutes(mux *http.ServeMux, pgPool *pgxpool.Pool) {
	// system
	mux.Handle("GET /livez", Livez())
	mux.Handle("GET /readyz", Readyz())
	// metrics management
	mux.Handle("POST /metrics/definition", CreateMetricDefinition(pgPool))
}
