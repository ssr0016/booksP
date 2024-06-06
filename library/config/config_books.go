package config

import (
	"practice/util/env"
	"time"
)

const (
	defaultBooksInterval = 60
)

type BooksConfig struct {
	BooksInterval time.Duration
}

func (c *Config) booksConfig() {
	booksInterval, _ := env.GetEnvAsInt("BOOKS_INTERVAL", defaultBooksInterval)
	c.Books.BooksInterval = time.Duration(booksInterval) * time.Second
}
