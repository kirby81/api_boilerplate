package main

import (
	"errors"
	"flag"
	"os"
)

type config struct {
	TokenSecret string
	Hostname    string
	Port        string
}

func newConfig() (*config, error) {
	c := &config{}

	// Env
	if c.TokenSecret = os.Getenv("APP_TOKEN_SECRET"); c.TokenSecret == "" {
		return nil, errors.New("env APP_TOKEN_SECRET is not set")
	}

	// Flags
	flag.StringVar(&c.Hostname, "host", "localhost", "defines the server's hostname")
	flag.StringVar(&c.Port, "port", "8000", "defines the server's port")
	flag.Parse()

	return c, nil
}
