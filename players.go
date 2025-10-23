package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	kcv1 "github.com/naivary/omp/api/keycloak/v1"
	playerv1 "github.com/naivary/omp/api/player/v1"
	"github.com/naivary/omp/keycloak"
)

func CreatePlayer(kc keycloak.Keycloak, pg *pgxpool.Pool) *Endpoint {
	return &Endpoint{
		Handler: createPlayer(kc, pg),
		Error:   defaultErrorHandler(),
	}
}

func createPlayer(kc keycloak.Keycloak, pg *pgxpool.Pool) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		p, err := decode[playerv1.CreatePlayerRequest](r)
		if err != nil {
			return err
		}
		kcUser := &kcv1.User{
			Email: p.Email,
			Credentials: []*kcv1.Credential{
				{Type: "password", Value: p.Password},
			},
		}
		err = kc.CreateUser(kcUser)
		if err != nil {
			return err
		}
		// TODO: Create the empty profile. Player will be enabled when he joins
		// a team/club
		return encode[any](w, r, http.StatusCreated, nil)
	})
}
