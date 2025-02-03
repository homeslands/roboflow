package config

import "fmt"

type PostgresConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	DB       string `env:"DB,required"`
	SSLMode  string `env:"SSLMODE" envDefault:"disable"`

	EnableLog bool  `env:"ENABLE_LOG" envDefault:"false"`
	MaxConns  int32 `env:"MAX_CONNS" envDefault:"20"`
	MinConns  int32 `env:"MIN_CONNS" envDefault:"10"`
}

// ConnectionString returns the PostgreSQL connection string
func (c PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DB, c.SSLMode)
}
