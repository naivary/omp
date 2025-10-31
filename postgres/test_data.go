package postgres

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	kcv1 "github.com/naivary/omp/api/keycloak/v1"
	"github.com/naivary/omp/keycloak"
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

func InsertTestdata(ctx context.Context, pool *pgxpool.Pool, kc keycloak.Keycloak) error {
	sql := `
	--  Insert 100 Clubs
	INSERT INTO club_profile (name, email, location, timezone)
	SELECT
		initcap(md5(random()::text)) || '_Club_' || g AS name,
		lower(md5(random()::text)) || '@omptest.de' as email,
		initcap(md5(random()::text)) || '_City' AS location,
		(ARRAY[
			'America/New_York','America/Los_Angeles','Europe/London','Europe/Berlin',
			'Asia/Tokyo','Asia/Dubai','Africa/Lagos','America/Sao_Paulo'
		 ])[1 + floor(random() * 8)::int] AS timezone
	FROM generate_series(1, 10) AS g;

	-- Insert 10 Teams per Club (1,000 total)
	INSERT INTO team (name, league, club_id)
	SELECT
		initcap(md5(random()::text)) || '_Team_' || t AS name,
		'League_' || ceil(random() * 5)::int AS league,
		c.id AS club_id
	FROM club_profile c
	CROSS JOIN LATERAL generate_series(1, 5) AS t;

	-- Insert 20 Players per Team (20,000 total)
	INSERT INTO player_profile (
		email, first_name, last_name, jersey_number, position, strong_foot, team_id
	)
	SELECT
		lower(md5(random()::text)) || '@omptest.de' as email,
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
	CROSS JOIN generate_series(1, 5) AS gs;
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
	_, err = tx.Exec(ctx, sql)
	if err != nil {
		return err
	}
	clubUsers, err := getClubUsers(ctx, tx)
	if err != nil {
		return err
	}
	playerUsers, err := getPlayerUsers(ctx, tx)
	if err != nil {
		return err
	}
	users := slices.Concat(clubUsers, playerUsers)
	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		go func(u *kcv1.User, w *sync.WaitGroup) {
			err = kc.CreateUser(ctx, user)
			if err != nil {
				panic(err)
			}
			w.Done()
		}(user, &wg)
	}
	wg.Wait()
	err = AddMetadata(ctx, tx, _ompMetadataKeyTestdataAlreadyInserted, "true")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func getClubUsers(ctx context.Context, tx pgx.Tx) ([]*kcv1.User, error) {
	profiles, err := tx.Query(ctx, "SELECT id, email FROM club_profile")
	if err != nil {
		return nil, err
	}
	users := make([]*kcv1.User, 0, profiles.CommandTag().RowsAffected())
	for profiles.Next() {
		var id int64
		var email string
		err = profiles.Scan(&id, &email)
		if err != nil {
			return nil, err
		}
		user := keycloak.NewUser(email, "omppasswordtest", nil, &keycloak.Attributes{ProfileID: id})
		user.Enabled = true
		users = append(users, user)
	}
	return users, nil
}

func getPlayerUsers(ctx context.Context, tx pgx.Tx) ([]*kcv1.User, error) {
	profiles, err := tx.Query(ctx, "SELECT id, email FROM player_profile")
	if err != nil {
		return nil, err
	}
	users := make([]*kcv1.User, 0, profiles.CommandTag().RowsAffected())
	for profiles.Next() {
		var id int64
		var email string
		err = profiles.Scan(&id, &email)
		if err != nil {
			return nil, err
		}
		user := keycloak.NewUser(email, "omppasswordtest", nil, &keycloak.Attributes{ProfileID: id})
		users = append(users, user)
	}
	return users, nil
}
