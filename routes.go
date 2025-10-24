package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/profiler"
	"github.com/naivary/omp/team"
)

func addRoutes(
	mux *http.ServeMux,
	pgPool *pgxpool.Pool,
	kc keycloak.Keycloak,
	playerProfiler profiler.PlayerProfiler,
	clubProfiler profiler.ClubProfiler,
	teamer team.Teamer,
) {
	// system
	mux.Handle("GET /livez", Livez())
	mux.Handle("GET /readyz", Readyz())

	// clubs
	mux.Handle("POST /clubs", CreateClub(kc, clubProfiler))

	// teams
	mux.Handle("POST /teams", CreateTeam(teamer))

	// players
	mux.Handle("POST /players", CreatePlayer(kc, playerProfiler))
}
