package main

import (
	"github.com/ahmadabdelrazik/jasad/internal/application"
	"github.com/ahmadabdelrazik/jasad/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	// initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	app := application.New(*cfg)

	if err := app.Serve(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

