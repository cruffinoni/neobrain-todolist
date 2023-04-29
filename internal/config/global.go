package config

type Global struct {
	Environment string   `env:"ENVIRONMENT" envDefault:"local"`
	APIPort     int      `env:"PORT" envDefault:"8080"`
	Database    Database `envPrefix:"DATABASE_"`
}
