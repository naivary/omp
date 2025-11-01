package main

import (
	"net/http"
	"strconv"

	clubv1 "github.com/naivary/omp/api/club/v1"
	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/openapi"
	"github.com/naivary/omp/profiler"
)

func CreateClubProfile(kc keycloak.Keycloak, p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler:     createClubProfile(kc, p),
		Error:       defaultErrorHandler(),
		Pattern:     "POST /clubs",
		Summary:     "Create a new club profile",
		Tags:        openapi.Tags("ClubProfile"),
		OperationID: "createClubProfile",
		RequestBody: openapi.NewReqBody[clubv1.CreateClubRequest]("create a new club profile", true),
		Responses: map[string]*openapi.Response{
			"201": openapi.NewResponse[clubv1.CreateClubResponse]("successfull response for a created club profile"),
			"400": openapi.NewResponse[HTTPError]("user with the email already exists"),
		},
	}
}

func createClubProfile(kc keycloak.Keycloak, p profiler.ClubProfiler) HandlerFuncErr {
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
			Email:    c.Email,
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

func ReadClubProfile(p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler:     readClubProfile(p),
		Error:       defaultErrorHandler(),
		Pattern:     "GET /clubs/{id}",
		Summary:     "read a single club profile based on the id",
		Tags:        openapi.Tags("ClubProfile"),
		OperationID: "readClubProfile",
		Responses: map[string]*openapi.Response{
			"201": openapi.NewResponse[clubv1.ReadClubProfileResponse]("successfull response"),
		},
		Parameters: []*openapi.Parameter{
			openapi.NewPathParam[int]("id"),
		},
	}
}

func readClubProfile(p profiler.ClubProfiler) HandlerFuncErr {
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

func UpdateClubProfile(p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler:     updateClubProfile(p),
		Error:       defaultErrorHandler(),
		Pattern:     "PATCH /clubs",
		Summary:     "update a club profile partially or completly",
		Tags:        openapi.Tags("ClubProfile"),
		OperationID: "updateClubProfile",
		RequestBody: openapi.NewReqBody[clubv1.UpdateClubRequest]("update a club profile", true),
		Responses: map[string]*openapi.Response{
			"204": {Description: "successfull response"},
		},
		Parameters: []*openapi.Parameter{
			openapi.NewPathParam[int]("id"),
		},
	}
}

func updateClubProfile(p profiler.ClubProfiler) HandlerFuncErr {
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

func DeleteClubProfile(kc keycloak.Keycloak, p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler:     deleteClubProfile(kc, p),
		Error:       defaultErrorHandler(),
		Pattern:     "DELETE /clubs/{id}",
		Summary:     "delete a club profile. This operation will be cascading and delete all entities related to this club profile.",
		Tags:        openapi.Tags("ClubProfile"),
		OperationID: "deleteClubProfile",
		RequestBody: openapi.NewReqBody[clubv1.DeleteClubRequest]("delete club profile and root user", true),
		Responses: map[string]*openapi.Response{
			"204": {Description: "Club Profile deleted successfully"},
		},
	}
}

func deleteClubProfile(kc keycloak.Keycloak, p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return err
		}
		c, err := decode[clubv1.DeleteClubRequest](r)
		if err != nil {
			return err
		}
		if !p.IsExisting(ctx, id) {
			return NewHTTPError(http.StatusBadRequest, "club does not exist: %d", id)
		}
		err = kc.RemoveUser(ctx, c.Email)
		if err != nil {
			return err
		}
		err = p.Remove(ctx, id)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func ReadAllClubProfiles(p profiler.ClubProfiler) *Endpoint {
	return &Endpoint{
		Handler:     readAllClubProfiles(p),
		Error:       defaultErrorHandler(),
		Pattern:     "GET /clubs",
		Summary:     "Get all available club profiles",
		Tags:        openapi.Tags("ClubProfile"),
		OperationID: "readAllClubProfiles",
		Responses: map[string]*openapi.Response{
			"200": openapi.NewResponse[](),
		},
	}
}

func readAllClubProfiles(p profiler.ClubProfiler) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		profiles, err := p.All(ctx)
		if err != nil {
			return err
		}
		return encode(w, r, http.StatusOK, profiles)
	})
}
