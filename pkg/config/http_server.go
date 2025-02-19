package config

type HTTPServerConfig struct {
	Port          int  `env:"PORT" envDefault:"8080"`
	EnableSwagger bool `env:"ENABLE_SWAGGER" envDefault:"true"`
}
