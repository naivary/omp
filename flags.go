package main

import "flag"

type config struct {
	port int
	host string
}

func parseFlags(args []string) (*config, error) {
	cfg := config{}
	fs := flag.NewFlagSet("omp", flag.ExitOnError)
	fs.IntVar(&cfg.port, "port", 9443, "port of http server")
	fs.StringVar(&cfg.host, "host", "127.0.0.1", "host of http server")
	return &cfg, fs.Parse(args)
}
