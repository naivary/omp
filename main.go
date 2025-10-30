//go:generate codemark gen -o openapi:fs ./... -- --fs.path=api/openapi/schemas
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/naivary/omp/keycloak"
	"github.com/naivary/omp/logger"
	"github.com/naivary/omp/postgres"
	"github.com/naivary/omp/profiler"
	"github.com/naivary/omp/team"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}
	logger := logger.New(&slog.HandlerOptions{AddSource: true})
	pool, err := postgres.Connect(ctx, cfg.pgHost, cfg.pgPort, cfg.pgUsername, cfg.pgPassword, cfg.pgDatabaseName)
	if err != nil {
		return err
	}
	if cfg.pgInsertTestData {
		err = postgres.InsertRandTestdata(ctx, pool)
		if err != nil {
			return err
		}
	}
	kc, err := keycloak.New(ctx, cfg.oidcURL, cfg.oidcRealm, cfg.oidcClientID, cfg.oidcClientSecret, cfg.oidcInsecure)
	if err != nil {
		return err
	}
	playerProfiler, err := profiler.NewPlayerProfiler(ctx, pool)
	if err != nil {
		return err
	}
	clubProfiler, err := profiler.NewClubProfiler(ctx, pool)
	if err != nil {
		return err
	}
	teamer, err := team.NewTeamer(pool)
	if err != nil {
		return err
	}

	// start the server with graceful handling
	interuptCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	host, port := cfg.host, strconv.Itoa(cfg.port)
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: newHandler(pool, kc, playerProfiler, clubProfiler, teamer),
		BaseContext: func(net.Listener) context.Context {
			return interuptCtx
		},
	}
	go func() {
		logger.Info("Server started!", "host", host, "port", port, "oidc.url", cfg.oidcURL, "oidc.realm", cfg.oidcRealm)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	// wait for interrupt signal for graceful shutdown procedure
	<-interuptCtx.Done()
	logger.Info("Interrupt signal received. Gracefully shutting down server")
	cancel() // instantly stop the application on further interrupt signals

	// new context to have a finite shutdown time
	shutdownCtx, shutdown := context.WithTimeout(ctx, 15*time.Second)
	defer shutdown()
	return srv.Shutdown(shutdownCtx)
}

func newHandler(
	pgPool *pgxpool.Pool,
	kc keycloak.Keycloak,
	playerProfiler profiler.PlayerProfiler,
	clubProfiler profiler.ClubProfiler,
	teamer team.TeamManager,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, pgPool, kc, playerProfiler, clubProfiler, teamer)
	return mux
}
