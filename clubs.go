package main

import (
	"net/http"

	clubv1 "github.com/naivary/omp/api/club/v1"
	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/profiler"
)

func CreateClub(kc keycloak.Keycloak, p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler: createClub(kc, p),
		Error:   defaultErrorHandler(),
	}
}

func createClub(kc keycloak.Keycloak, p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		c, err := decode[clubv1.CreateClubRequest](r)
		if err != nil {
			return err
		}
		profile := clubv1.Profile{
			Name:     c.Name,
			Location: c.Location,
			Timezone: c.Timezone,
		}
		profileID, err := p.Create(ctx, &profile)
		if err != nil {
			return err
		}
		user := keycloak.NewUser(
			c.Email,
			c.Password,
			nil,
			&keycloak.Attributes{
				ProfileID: profileID,
			},
		)
		// club root user is always enabled by default
		user.Enabled = true
		err = kc.CreateUser(ctx, user)
		if err != nil {
			return err
		}
		return encode(w, r, http.StatusCreated, clubv1.CreateClubResponse{ID: profileID})
	})
}

func RemoveClub(kc keycloak.Keycloak, p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler: removeClub(kc, p),
		Error:   defaultErrorHandler(),
	}
}

func removeClub(kc keycloak.Keycloak, p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		c, err := decode[clubv1.DeleteClubRequest](r)
		if err != nil {
			return err
		}
		err = kc.RemoveUser(ctx, c.Email)
		if err != nil {
			return err
		}
		return p.Remove(ctx, c.ClubID)
	})
}
