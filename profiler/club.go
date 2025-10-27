package profiler

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	clubv1 "github.com/naivary/omp/api/club/v1"
)

type ClubProfiler interface {
	Create(profile *clubv1.Profile) (int64, error)
	Remove(id int64) error
}

var _ ClubProfiler = (*clubProfiler)(nil)

func NewClubProfiler(ctx context.Context, pool *pgxpool.Pool) (ClubProfiler, error) {
	c := clubProfiler{
		ctx:  ctx,
		pool: pool,
	}
	return &c, pool.Ping(ctx)
}

type clubProfiler struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func (c *clubProfiler) Create(profile *clubv1.Profile) (int64, error) {
	var id int64
	tx, err := c.pool.Begin(c.ctx)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(c.ctx,
		`INSERT INTO club_profile(name, location, timezone) VALUES ($1, $2, $3)`,
		profile.Name, profile.Location, profile.Timezone,
	)
	if err != nil {
		return 0, err
	}
	err = tx.Commit(c.ctx)
	if err != nil {
		return 0, err
	}
	row := c.pool.QueryRow(c.ctx,
		`SELECT id FROM club_profile WHERE name = $1`,
		profile.Name,
	)
	return id, row.Scan(&id)
}

func (c *clubProfiler) Remove(id int64) error {
	return nil
}
