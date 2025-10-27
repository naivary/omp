package team

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	teamv1 "github.com/naivary/omp/api/team/v1"
)

type TeamManager interface {
	Create(ctx context.Context, team *teamv1.Team) (int64, error)
}

var _ TeamManager = (*teamManager)(nil)

type teamManager struct {
	pool *pgxpool.Pool
}

func NewTeamer(pool *pgxpool.Pool) (TeamManager, error) {
	return &teamManager{pool: pool}, nil
}

func (t *teamManager) Create(ctx context.Context, team *teamv1.Team) (int64, error) {
	var id int64
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO team(name, league, club_id) VALUES($1, $2, $3)`,
		team.Name, team.League, team.ClubID,
	)
	if err != nil {
		return 0, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	row := t.pool.QueryRow(ctx,
		"SELECT id FROM team WHERE club_id = $1 AND name = $2",
		team.ClubID, team.Name,
	)
	return id, row.Scan(&id)
}
