package config

import "practice/util/env"

const (
	Dev         = "development"
	Prod        = "production"
	Test        = "test"
	ServiceName = "library"
)

var (
	Env = Dev
)

type Config struct {
	App          string
	BuildVersion string

	Server     ServerConfig
	Postgres   PostgresConfig
	Pagination PaginationConfig
	Books      BooksConfig
}

func FromEnv() (*Config, error) {
	cfg := &Config{}

	Env = env.GetEnvAsString("APP_ENV", Env)
	cfg.BuildVersion = "2.0"

	cfg.serverConfig()
	cfg.postgresConfig()
	cfg.PaginationConfig()
	cfg.booksConfig()

	return cfg, nil
}
