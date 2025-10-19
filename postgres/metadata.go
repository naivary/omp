package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddMetadata(ctx context.Context, pool *pgxpool.Pool, key, value string) error {
	return nil
}

func GetMetadata(ctx context.Context, pool *pgxpool.Pool, key string) (string, error) {
	return "", nil
}
