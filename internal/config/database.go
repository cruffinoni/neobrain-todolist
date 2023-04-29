package config

type Database struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT" envDefault:"3306"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}
