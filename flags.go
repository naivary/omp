package main

import (
	"flag"
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
	fs.StringVar(&cfg.psqlDatabaseName, "psql.database", "omp", "database name")
	return &cfg, fs.Parse(args[1:])
}
