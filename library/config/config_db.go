package config

import (
	"fmt"
	"practice/util/env"
)

const (
	defaultDBHost        = "localhost"
	defaultDBPort        = "5432"
	defaultDBUser        = "postgres"
	defaultDBPassword    = "secret"
	defaultDBName        = "libraryDB"
	defaultDBApplication = "library"
)

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DB              string
	ApplicationName string
}

func (p *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable application_name=%s binary_parameters=yes",
		p.Host, p.Port, p.User, p.Password, p.DB, p.ApplicationName)
}

func (cfg *Config) postgresConfig() {
	cfg.Postgres.Host = env.GetEnvAsString("POSTGRES_HOST", defaultDBHost)
	cfg.Postgres.Port = env.GetEnvAsString("POSTGRES_PORT", defaultDBPort)
	cfg.Postgres.User = env.GetEnvAsString("POSTGRES_USER", defaultDBUser)
	cfg.Postgres.Password = env.GetEnvAsString("POSTGRES_PASSWORD", defaultDBPassword)
	cfg.Postgres.DB = env.GetEnvAsString("POSTGRES_DB_NAME", defaultDBName)
	cfg.Postgres.ApplicationName = env.GetEnvAsString("POSTGRES_APPLICATION_NAME", defaultDBApplication)
}
