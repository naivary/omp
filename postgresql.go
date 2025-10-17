package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func connectToPostgresDB(ctx context.Context, cfg *config) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.psqlUsername, cfg.psqlPassword, cfg.psqlHost, cfg.psqlPort, cfg.psqlDatabaseName)
	return pgx.Connect(ctx, connString)
}
