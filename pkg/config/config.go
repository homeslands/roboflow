package config

import (
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

const (
	ProdEnvVal = "production"
)

type Config struct {
	Env         string `env:"ROBOFLOW_ENV" envDefault:"development"`
	Port        int    `env:"ROBOFLOW_PORT" envDefault:"8000"`
	PostgresDsn string `env:"ROBOFLOW_POSTGRES_DSN"`
}

func (c *Config) IsProd() bool {
	return c.Env == ProdEnvVal
}

func MustLoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}
