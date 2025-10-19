package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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
