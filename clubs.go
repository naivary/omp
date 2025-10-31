package main

import (
	"fmt"
	"net/http"
	"strconv"

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
		isUsed, err := kc.IsEmailUsed(ctx, c.Email)
		if err != nil {
			return err
		}
		if isUsed {
			return NewHTTPError(http.StatusBadRequest, "user with the email already exists: %s", c.Email)
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

func ReadClub(p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler: readClub(p),
		Error:   defaultErrorHandler(),
	}
}

func readClub(p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return err
		}
		profile, err := p.Read(ctx, id)
		if err != nil {
			return err
		}
		return encode(w, r, http.StatusOK, profile)
	})
}

func UpdateClub(p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler: updateClub(p),
		Error:   defaultErrorHandler(),
	}
}

func updateClub(p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return err
		}
		update, err := decode[clubv1.UpdateClubRequest](r)
		if err != nil {
			return err
		}
		profile := clubv1.Profile{
			ID:       id,
			Name:     update.Name,
			Location: update.Location,
			Timezone: update.Timezone,
		}
		err = p.Update(ctx, &profile)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func DeleteClub(kc keycloak.Keycloak, p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler: deleteClub(kc, p),
		Error:   defaultErrorHandler(),
	}
}

func deleteClub(kc keycloak.Keycloak, p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		c, err := decode[clubv1.DeleteClubRequest](r)
		if err != nil {
			return err
		}
		isExisting := p.IsExisting(ctx, c.ClubID)
		if !isExisting {
			fmt.Println("###################### IM HERE")
			return NewHTTPError(http.StatusBadRequest, "club does not exist: %d", c.ClubID)
		}
		err = p.Remove(ctx, c.ClubID)
		if err != nil {
			return err
		}
		return kc.RemoveUser(ctx, c.Email)
	})
}

func ReadAllClubs(p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler: readAllClubs(p),
		Error:   defaultErrorHandler(),
	}
}

func readAllClubs(p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		profiles, err := p.All(ctx)
		if err != nil {
			return err
		}
		return encode(w, r, http.StatusOK, profiles)
	})
}
