package config

import (
	"fmt"
	//nolint:revive

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Log        LogConfig        `envPrefix:"LOG_"`
	HTTPServer HTTPServerConfig `envPrefix:"HTTP_SERVER_"`
	Postgres   PostgresConfig   `envPrefix:"PG_"`
	Nats       NatsConfig       `envPrefix:"NATS_"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing env: %w", err)
	}

	return cfg, nil
}
