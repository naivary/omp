package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/openapi"
	"github.com/naivary/omp/profiler"
	"github.com/naivary/omp/team"
)

func addRoutes(
	mux *http.ServeMux,
	pgPool *pgxpool.Pool,
	kc keycloak.Keycloak,
	playerProfiler profiler.PlayerProfiler,
	clubProfiler profiler.ClubProfiler,
	teamer team.TeamManager,
) error {
	rootOpenAPI := openapi.New("3.0.0", "mustafa.hussaini", "info@omp.de", openapi.Apache)
	endpoints := []*Endpoint{
		// system
		Readyz(),
		// club profile
		// CreateClub(kc, clubProfiler),
		// ReadClub(clubProfiler),
		// UpdateClub(clubProfiler),
		// DeleteClub(kc, clubProfiler),
		//
		// // teams
		// CreateTeam(teamer),

		// player profile
		CreatePlayerProfile(kc, playerProfiler),
	}
	err := GenOpenAPISpecs(rootOpenAPI, endpoints...)
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		mux.Handle(endpoint.Pattern, endpoint)
	}
	return nil
}
