package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/naivary/omp/keycloak"
)

func addRoutes(mux *http.ServeMux, pgPool *pgxpool.Pool, kc keycloak.Keycloak) {
	// system
	mux.Handle("GET /livez", Livez())
	mux.Handle("GET /readyz", Readyz())

	// players
	mux.Handle("POST /players", CreatePlayer(kc, pgPool))
}
