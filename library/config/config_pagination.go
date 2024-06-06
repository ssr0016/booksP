package config

import "practice/util/env"

const (
	defaultPage     = 1
	defaulPageLimit = 10000
)

type PaginationConfig struct {
	Page      int
	PageLimit int
}

func (cfg *Config) PaginationConfig() {
	cfg.Pagination.Page, _ = env.GetEnvAsInt("PAGE", defaultPage)
	cfg.Pagination.PageLimit, _ = env.GetEnvAsInt("PER_PAGE", defaulPageLimit)
}
