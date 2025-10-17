package main

import (
	"context"
	"flag"
	"fmt"
)

type config struct {
	port int
	host string

	// postgres database
	psqlHost         string
	psqlPort         int
	psqlUsername     string
	psqlPassword     string
	psqlDatabaseName string
}

func parseFlags(args []string) (*config, error) {
	cfg := config{}
	fs := flag.NewFlagSet("omp", flag.ExitOnError)
	// api
	fs.IntVar(&cfg.port, "port", 9443, "port of http server")
	fs.StringVar(&cfg.host, "host", "127.0.0.1", "host of http server")
	// psql database
	fs.StringVar(&cfg.psqlHost, "psql.host", "127.0.0.1", "host of postgresql server")
	fs.IntVar(&cfg.psqlPort, "psql.port", 5432, "port of postgresql server")
	fs.StringVar(&cfg.psqlUsername, "psql.username", "", "username of postgresql server to use for authentication")
	fs.StringVar(&cfg.psqlPassword, "psql.password", "", "password of postregsql server to use for authentication")
	fs.StringVar(&cfg.psqlDatabaseName, "psql.db.name", "omp", "database name")
	fmt.Println(args)
	return &cfg, fs.Parse(args[1:])
}

func (c *config) Validate(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if c.psqlUsername == "" {
		problems["psql username empty"] = "psql username cannot be empty"
	}
	if c.psqlPassword == "" {
		problems["psql password empty"] = "psql password cannot be empty"
	}
	return problems
}
