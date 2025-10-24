package main

import (
	"net/http"

	playerv1 "github.com/naivary/omp/api/player/v1"
	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/profiler"
)

func CreatePlayer(kc keycloak.Keycloak, playerProfiler profiler.PlayerProfiler) *Endpoint {
	return &Endpoint{
		Handler: createPlayer(kc, playerProfiler),
		Error:   defaultErrorHandler(),
	}
}

func createPlayer(kc keycloak.Keycloak, playerProfiler profiler.PlayerProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		p, err := decode[playerv1.CreatePlayerRequest](r)
		if err != nil {
			return err
		}
		profile := playerv1.Profile{
			Email:      p.Email,
			TeamID:     p.TeamID,
			FirstName:  p.FirstName,
			LastName:   p.LastName,
			StrongFoot: p.StrongFoot,
			Position:   p.Position,
		}
		profileID, err := playerProfiler.Create(&profile)
		if err != nil {
			return err
		}
		kcUser := keycloak.NewUser(
			p.Email,
			p.Password,
			nil,
			&keycloak.Attributes{
				ProfileID: profileID,
			},
		)
		err = kc.CreateUser(kcUser)
		if err != nil {
			return err
		}
		return encode(w, r, http.StatusCreated, playerv1.CreatePlayerResponse{ID: profileID})
	})
}
