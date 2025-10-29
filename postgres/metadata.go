package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func AddMetadata(ctx context.Context, e Execer, key, value string) error {
	_, err := e.Exec(ctx, `INSERT INTO omp_metadata VALUES($1, $2);`, key, value)
	return err
}

func GetMetadata(ctx context.Context, q Querier, key string) (string, error) {
	var value string
	row := q.QueryRow(ctx, `SELECT value FROM omp_metadata WHERE key = $1`, key)
	return value, row.Scan(&value)
}
