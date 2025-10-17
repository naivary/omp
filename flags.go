package main

import (
	"flag"
)

type config struct {
	port int
	host string

	// postgres database
	pgHost         string
	pgPort         int
	pgUsername     string
	pgPassword     string
	pgDatabaseName string
}

func parseFlags(args []string) (*config, error) {
	const required = "<None>"
	cfg := config{}
	fs := flag.NewFlagSet("omp", flag.ExitOnError)
	// http server
	fs.IntVar(&cfg.port, "port", 9443, "port of http server")
	fs.StringVar(&cfg.host, "host", "127.0.0.1", "host of http server")
	// psql database
	fs.StringVar(&cfg.pgHost, "pg.host", "127.0.0.1", "host of postgresql server")
	fs.IntVar(&cfg.pgPort, "pg.port", 5432, "port of postgresql server")
	fs.StringVar(&cfg.pgUsername, "pg.username", required, "username of postgresql server to use for authentication")
	fs.StringVar(&cfg.pgPassword, "pg.password", required, "password of postregsql server to use for authentication")
	fs.StringVar(&cfg.pgDatabaseName, "pg.database", "omp", "database name")
	return &cfg, fs.Parse(args[1:])
}
