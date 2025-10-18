package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, host string, port int, username, password, database string) (*pgx.Conn, error) {
	if username == "" {
		return nil, errors.New("pg username not defined")
	}
	if password == "" {
		return nil, errors.New("pg password not defined")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)
	return pgx.Connect(ctx, connString)
}

func provision(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
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
	$$ LANGUAGE plpgsql strict immutable;`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `
	CREATE SEQUENCE IF NOT EXISTS id_seq START 1;
	CREATE TABLE IF NOT EXISTS club(
		id int PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
		name text,
		timezone text DEFAULT 'Europe/Berlin'
	);`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `
	CREATE SEQUENCE IF NOT EXISTS id_seq START 1;
	CREATE TABLE IF NOT EXISTS team(
		id int PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
		club_id int REFERENCES club(id),
		name text,
		league text
	);`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `
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
	);`)
	if err != nil {
		return err
	}
	return nil
}
