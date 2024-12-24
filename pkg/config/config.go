package config

import (
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

const (
	prodEnvVal = "production"
)

type Config struct {
	Env         string `env:"ENV" envDefault:"development"`
	Port        int    `env:"PORT" envDefault:"8000"`
	PostgresDsn string `env:"POSTGRES_DSN"`
}

func (c *Config) IsProd() bool {
	return c.Env == prodEnvVal
}

func MustLoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}
