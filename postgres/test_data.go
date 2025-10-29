package postgres

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func isTestdataAlreadyInserted(ctx context.Context, tx pgx.Tx) (bool, error) {
	value, err := GetMetadata(ctx, tx, _ompMetadataKeyTestdataAlreadyInserted)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return strconv.ParseBool(value)
}

func InsertRandTestdata(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	--  Insert 100 Clubs
	INSERT INTO club_profile (name, location, timezone)
	SELECT
		initcap(md5(random()::text)) || '_Club_' || g AS name,
		initcap(md5(random()::text)) || '_City' AS location,
		(ARRAY[
			'America/New_York','America/Los_Angeles','Europe/London','Europe/Berlin',
			'Asia/Tokyo','Asia/Dubai','Africa/Lagos','America/Sao_Paulo'
		 ])[1 + floor(random() * 8)::int] AS timezone
	FROM generate_series(1, 100) AS g;

	-- Insert 10 Teams per Club (1,000 total)
	INSERT INTO team (name, league, club_id)
	SELECT
		initcap(md5(random()::text)) || '_Team_' || t AS name,
		'League_' || ceil(random() * 5)::int AS league,
		c.id AS club_id
	FROM club_profile c
	CROSS JOIN LATERAL generate_series(1, 10) AS t;

	-- Insert 20 Players per Team (20,000 total)
	INSERT INTO player_profile (
		email, first_name, last_name, jersey_number, position, strong_foot, team_id
	)
	SELECT
		lower(
			(ARRAY['alex','jordan','taylor','morgan','riley','casey','jamie','avery','skyler','quinn'])
			[1 + floor(random() * 10)::int]
			|| '.' ||
			(ARRAY['smith','johnson','williams','brown','jones','garcia','miller','davis','martinez','lopez'])
			[1 + floor(random() * 10)::int]
			|| '_' || t.id || '_' || gs || '@example.com'
		) AS email,
		(ARRAY['Alex','Jordan','Taylor','Morgan','Riley','Casey','Jamie','Avery','Skyler','Quinn'])
			[1 + floor(random() * 10)::int] AS first_name,
		(ARRAY['Smith','Johnson','Williams','Brown','Jones','Garcia','Miller','Davis','Martinez','Lopez'])
			[1 + floor(random() * 10)::int] AS last_name,
		1 + floor(random() * 99)::int AS jersey_number,
		(ARRAY[
			-- Goalkeepers
			'Goalkeeper',
			-- Defenders
			'Right Back','Left Back','Centre Back','Sweeper','Wing Back',
			-- Midfielders
			'Defensive Midfielder','Central Midfielder','Attacking Midfielder',
			'Right Midfielder','Left Midfielder','Wide Midfielder',
			-- Forwards
			'Striker','Centre Forward','Second Striker','Right Winger','Left Winger'
		])[1 + floor(random() * 17)::int] AS position,
		(ARRAY['Right','Left','Both'])[1 + floor(random() * 3)::int]::foot AS strong_foot,
		t.id AS team_id
	FROM team t
	CROSS JOIN generate_series(1, 20) AS gs;
	`
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	isAlreadyInserted, err := isTestdataAlreadyInserted(ctx, tx)
	if err != nil {
		return err
	}
	if isAlreadyInserted {
		return nil
	}
	err = acquireAdvisoryLock(ctx, tx, _advisoryLockInsertTestdata)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, sql)
	if err != nil {
		return err
	}
	err = AddMetadata(ctx, tx, _ompMetadataKeyTestdataAlreadyInserted, "true")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
