package profiler

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	playerv1 "github.com/naivary/omp/api/player/v1"
)

type PlayerProfiler interface {
	Create(profile *playerv1.Profile) (int64, error)
}

var _ PlayerProfiler = (*playerProfiler)(nil)

type playerProfiler struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewPlayerProfiler(ctx context.Context, pool *pgxpool.Pool) (PlayerProfiler, error) {
	pp := &playerProfiler{
		ctx:  ctx,
		pool: pool,
	}
	return pp, pool.Ping(ctx)
}

func (p *playerProfiler) Create(profile *playerv1.Profile) (int64, error) {
	var id int64
	tx, err := p.pool.Begin(p.ctx)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(p.ctx,
		`INSERT INTO player_profile(email, first_name, last_name, jersey_number, position, strong_foot, team_id) VALUES($1, $2, $3, $4, $5, $6)`,
		profile.Email, profile.FirstName, profile.LastName, profile.JerseyNumber, profile.Position, profile.StrongFoot, profile.TeamID,
	)
	if err != nil {
		return 0, err
	}
	err = tx.Commit(p.ctx)
	if err != nil {
		return 0, err
	}
	row := p.pool.QueryRow(p.ctx,
		`SELECT id FROM player_profile WHERE email = $1`,
		profile.Email,
	)
	return id, row.Scan(&id)
}
