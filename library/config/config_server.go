package config

import "practice/util/env"

const (
	defaultHTTPPort = "3000"
	defaultHost     = "localhost"
	// defaultAllowedDomains = "localhost"
	// defaultGrpcPort       = ""
)

type ServerConfig struct {
	Host           string
	HTTPPort       string
	AllowedDomains string
	RPCPort        string
}

func (cfg *Config) serverConfig() {
	cfg.Server.HTTPPort = env.GetEnvAsString("BOOKS_HTTP_PORT", defaultHTTPPort)
	cfg.Server.Host = env.GetEnvAsString("BOOKS_HOST", defaultHost)
	// cfg.Server.AllowedDomains = env.GetEnvAsString("ALLOWED_DOMAINS", defaultAllowedDomains)
	// cfg.Server.RPCPort = env.GetEnvAsString("BOOKS_RPC_PORT", defaultGrpcPort)
}
