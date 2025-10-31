package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/openapi"
	"github.com/naivary/omp/profiler"
	"github.com/naivary/omp/team"
)

var rootOpenAPI = openapi.New("3.0.0", "mustafa.hussaini", "info@omp.de", openapi.Apache)

func addRoutes(
	mux *http.ServeMux,
	pgPool *pgxpool.Pool,
	kc keycloak.Keycloak,
	playerProfiler profiler.PlayerProfiler,
	clubProfiler profiler.ClubProfiler,
	teamer team.TeamManager,
) error {
	endpoints := []*Endpoint{
		// CreateClub(kc, clubProfiler),
		// ReadClub(clubProfiler),
		// UpdateClub(clubProfiler),
		// DeleteClub(kc, clubProfiler),
		//
		// // teams
		// CreateTeam(teamer),

		// players
		CreatePlayerProfile(kc, playerProfiler),
	}
	err := GenOpenAPISpecs(rootOpenAPI, endpoints...)
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		mux.Handle(endpoint.Pattern, endpoint)
	}
	// system
	mux.Handle("GET /livez", Livez())
	mux.Handle("GET /readyz", Readyz())
	return nil
}
