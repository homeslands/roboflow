package config

type NatsConfig struct {
	// URL is the URL of the NATS server. If not provided, the default URL will be used.
	// That means the NATS server will be embedded in the application.
	URL *string `env:"URL"`
	// EnableLog enables logging for the NATS server.
	EnableLog bool `env:"ENABLE_LOG" envDefault:"false"`
}
