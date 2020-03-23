package main

import (
	"os"

	"github.com/kirby81/api-boilerplate/internal/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	s, err := api.NewServer()
	if err != nil {
		log.Fatal().Err(err).Msg("Server init error")
	}

	if err = s.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server runtime error")
	}

}
