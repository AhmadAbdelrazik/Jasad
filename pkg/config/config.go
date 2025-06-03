package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DSN                string `env:"JASAD_DB_DSN"`
	Origin             string `env:"ORIGIN"`
	Port               int    `env:"PORT"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

func Load(fileNames ...string) (*Config, error) {
	if err := godotenv.Load(fileNames...); err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment")
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
