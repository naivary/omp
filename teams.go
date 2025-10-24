package main

import (
	"net/http"

	teamv1 "github.com/naivary/omp/api/team/v1"
	"github.com/naivary/omp/team"
)

func CreateTeam(teamer team.Teamer) *Endpoint {
	return &Endpoint{
		Handler: createTeam(teamer),
		Error:   defaultErrorHandler(),
	}
}

func createTeam(teamer team.Teamer) HandlerFuncErr {
	return HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		c, err := decode[teamv1.CreateTeamRequest](r)
		if err != nil {
			return err
		}
		team := teamv1.Team{
			Name:   c.Name,
			ClubID: c.ClubID,
			League: c.League,
		}
		teamID, err := teamer.Create(ctx, &team)
		if err != nil {
			return err
		}
		return encode(w, r, http.StatusCreated, teamv1.CreateTeamResponse{ID: teamID})
	})
}
