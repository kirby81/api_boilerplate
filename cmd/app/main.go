package main

import (
	"os"

	"github.com/kirby81/api-boilerplate/internal/api"
	"github.com/kirby81/api-boilerplate/internal/auth"
	"github.com/kirby81/api-boilerplate/internal/auth/memory"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func initConfig() (*config, error) {
	conf, err := newConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func main() {
	conf, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init configuration")
	}

	initLogger()

	authRepo := memory.NewRepository()
	authService, err := auth.NewService(authRepo, conf.TokenSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init auth service")
	}

	authHandler, err := auth.NewHandler(authService)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init auth handler")
	}

	s, err := api.NewServer(routes(authHandler))
	if err != nil {
		log.Fatal().Err(err).Msg("Server init error")
	}

	if err = s.Run(conf.Hostname, conf.Port); err != nil {
		log.Fatal().Err(err).Msg("Server runtime error")
	}

}
