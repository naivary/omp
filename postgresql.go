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
