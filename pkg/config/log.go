package config

import (
	"log/slog"

	"github.com/tuanvumaihuynh/roboflow/pkg/log"
)

type LogConfig struct {
	Format    log.Format `env:"FORMAT" envDefault:"json"`
	Level     slog.Level `env:"LEVEL" envDefault:"info"`
	AddSource bool       `env:"ADD_SOURCE" envDefault:"false"`
}
