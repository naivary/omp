package profiler

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	playerv1 "github.com/naivary/omp/api/player/v1"
)

type PlayerProfiler interface {
	Create(ctx context.Context, profile *playerv1.Profile) (int64, error)
}

var _ PlayerProfiler = (*playerProfiler)(nil)

type playerProfiler struct {
	pool *pgxpool.Pool
}

func NewPlayerProfiler(ctx context.Context, pool *pgxpool.Pool) (PlayerProfiler, error) {
	pp := &playerProfiler{
		pool: pool,
	}
	return pp, pool.Ping(ctx)
}

func (p *playerProfiler) Create(ctx context.Context, profile *playerv1.Profile) (int64, error) {
	var id int64
	if err := ctx.Err(); err != nil {
		return id, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return id, err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO player_profile(email, first_name, last_name, jersey_number, position, strong_foot, team_id) VALUES($1, $2, $3, $4, $5, $6, $7)`,
		profile.Email, profile.FirstName, profile.LastName, profile.JerseyNumber, profile.Position, profile.StrongFoot, profile.TeamID,
	)
	if err != nil {
		return id, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return id, err
	}
	row := p.pool.QueryRow(ctx,
		`SELECT id FROM player_profile WHERE email = $1`,
		profile.Email,
	)
	return id, row.Scan(&id)
}
