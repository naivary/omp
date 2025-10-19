package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_pgCodeRelationDoesNotExist = "42P01"
)

const (
	_advisoryLockProvision = iota + 1
)

func Connect(ctx context.Context, host string, port int, username, password, database string) (*pgxpool.Pool, error) {
	if username == "" {
		return nil, errors.New("pg username not defined")
	}
	if password == "" {
		return nil, errors.New("pg password not defined")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return pool, provision(ctx, pool)
}

func acquireAdvisoryLock(ctx context.Context, tx pgx.Tx, lockFor int) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_lock($1)", lockFor)
	return err
}

func isProvisioned(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	var pgErr *pgconn.PgError
	var isAlreadyProvisioned string
	rows, err := conn.Query(ctx, `SELECT value FROM omp_metadata WHERE key = 'isProvisioned'`)
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == _pgCodeRelationDoesNotExist {
			return false, nil
		}
	}
	defer rows.Close()
	if !rows.Next() {
		return true, errors.New(
			"now row found for key 'isProvisioned'. This error should never occur. If it does then the key might have changed or the name of the table",
		)
	}
	err = rows.Scan(&isAlreadyProvisioned)
	return isAlreadyProvisioned == "true", err
}

func provision(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	if isAlreadyProvisioned, err := isProvisioned(ctx, conn); err != nil || isAlreadyProvisioned {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	err = acquireAdvisoryLock(ctx, tx, _advisoryLockProvision)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
	CREATE OR REPLACE FUNCTION pseudo_encrypt(value bigint) returns bigint AS $$
	DECLARE
	l1 int;
	l2 int;
	r1 int;
	r2 int;
	i int:=0;
	BEGIN
	 l1:= (value >> 16) & 65535;
	 r1:= value & 65535;
	 WHILE i < 3 LOOP
	   l2 := r1;
	   r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * 32767)::int;
	   l1 := l2;
	   r1 := r2;
	   i := i + 1;
	 END LOOP;
	 return ((r1 << 16) + l1);
	END;
	$$ LANGUAGE plpgsql strict immutable;

	-- metadata table
	CREATE TABLE omp_metadata(
		key text PRIMARY KEY,
		value text NOT NULL
	);

	CREATE SEQUENCE IF NOT EXISTS id_seq START 1;
	CREATE TABLE IF NOT EXISTS club(
		id int PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
		name text,
		timezone text DEFAULT 'Europe/Berlin'
	);

	CREATE SEQUENCE IF NOT EXISTS id_seq START 1;
	CREATE TABLE IF NOT EXISTS team(
		id int PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
		club_id int REFERENCES club(id),
		name text,
		league text
	);
	CREATE SEQUENCE IF NOT EXISTS id_seq START 1;
	CREATE TYPE scope AS ENUM('Club', 'Team', 'Private');

	CREATE TABLE IF NOT EXISTS metric_type(
		name text PRIMARY KEY
	);

	INSERT INTO metric_type VALUES('Counter');
	INSERT INTO metric_type VALUES('Gauge');
	CREATE TABLE IF NOT EXISTS metric_defintion(
		id int PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
		name text,
		type text REFERENCES metric_type(name),
		scope scope,
		description text,
		club_id int REFERENCES club(id),
		team_id int REFERENCES team(id),
		UNIQUE(name, scope, club_id),
		UNIQUE(name, scope, team_id)
	);

	-- make sure to update the omp_metadata table to not provision again.
	INSERT INTO omp_metadata VALUES('isProvisioned', 'true');
	`)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
