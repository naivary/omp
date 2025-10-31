package profiler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	clubv1 "github.com/naivary/omp/api/club/v1"
)

type ClubProfiler interface {
	// Create a new profile
	Create(ctx context.Context, profile *clubv1.Profile) (int64, error)
	// Remove the profile associatted with the id
	Remove(ctx context.Context, id int64) error
	// Update a profile partially or completly.
	Update(ctx context.Context, profile *clubv1.Profile) error
	// Read a profile associatted with the id
	Read(ctx context.Context, id int64) (*clubv1.Profile, error)
	// Retrieve all profiles
	All(ctx context.Context) ([]*clubv1.Profile, error)

	// Non-nil error means the user does not exist
	IsExisting(ctx context.Context, id int64) bool
}

var _ ClubProfiler = (*clubProfiler)(nil)

func NewClubProfiler(ctx context.Context, pool *pgxpool.Pool) (ClubProfiler, error) {
	c := clubProfiler{
		pool: pool,
	}
	return &c, pool.Ping(ctx)
}

type clubProfiler struct {
	pool *pgxpool.Pool
}

func (c *clubProfiler) Create(ctx context.Context, profile *clubv1.Profile) (int64, error) {
	var id int64
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO club_profile(name, email, location, timezone) VALUES ($1, $2, $3)`,
		profile.Name, profile.Email, profile.Location, profile.Timezone,
	)
	if err != nil {
		return 0, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	row := c.pool.QueryRow(ctx,
		`SELECT id FROM club_profile WHERE name = $1`,
		profile.Name,
	)
	return id, row.Scan(&id)
}

func (c *clubProfiler) Read(ctx context.Context, id int64) (*clubv1.Profile, error) {
	p := clubv1.Profile{}
	row := c.pool.QueryRow(ctx,
		`SELECT * FROM club_profile WHERE id = $1`,
		id,
	)
	return &p, row.Scan(&p.ID, &p.Name, &p.Location, &p.Timezone)
}

func (c *clubProfiler) Update(ctx context.Context, profile *clubv1.Profile) error {
	if profile.ID == 0 {
		return errors.New("update club profile: missing id")
	}
	_, err := c.pool.Exec(ctx,
		`
		UPDATE club_profile
		SET
		  name     = CASE WHEN $2 <> '' THEN $2 ELSE name END,
		  location = CASE WHEN $3 <> '' THEN $3 ELSE location END,
		  timezone = CASE WHEN $4 <> '' THEN $4 ELSE timezone END
		WHERE id = $1;
		`,
		profile.ID, profile.Name, profile.Location, profile.Timezone,
	)
	return err
}

func (c *clubProfiler) Remove(ctx context.Context, id int64) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		`DELETE FROM club_profile WHERE id = $1`,
		id,
	)
	return err
}

func (c *clubProfiler) All(ctx context.Context) ([]*clubv1.Profile, error) {
	rows, err := c.pool.Query(ctx, `SELECT * FROM club_profile`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	affectedRows := rows.CommandTag().RowsAffected()
	profiles := make([]*clubv1.Profile, 0, affectedRows)
	for rows.Next() {
		profile := clubv1.Profile{}
		err = rows.Scan(&profile.ID, &profile.Name, &profile.Location, &profile.Timezone)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, &profile)
	}
	return profiles, nil
}

func (c *clubProfiler) IsExisting(ctx context.Context, id int64) bool {
	_, err := c.Read(ctx, id)
	return err == nil
}
