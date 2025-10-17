package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func connectToPostgresDB(ctx context.Context, host string, port int, username, password, database string) (*pgx.Conn, error) {
	if username == "" {
		return nil, errors.New("psql username not defined")
	}
	if password == "" {
		return nil, errors.New("psql password not defined")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)
	return pgx.Connect(ctx, connString)
}

func provisionPostgresDB(ctx context.Context, con *pgx.Conn) error {
	_, err := con.Exec(ctx, `
	CREATE TYPE scope AS ENUM('Club', 'Team', 'Private');
	CREATE TYPE metric_type AS ENUM('Counter', 'Gauge');
	CREATE TABLE IF NOT EXISTS metric_defintion(
		id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
		name text,
		type metric_type,
		scope scope,
		owner bigint,
		description text
	);`)
	if err != nil {
		return err
	}
	return nil
}
