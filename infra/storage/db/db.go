package db

import (
	"context"
	"database/sql"
)

type DB interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	WithTransaction(ctx context.Context, fn func(sess Tx) error) error
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	Ping() error
	Close() error
}

type Tx interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecWithReturningId(ctx context.Context, query string, args ...interface{}) (int64, error)
	NamedExecReturningId(ctx context.Context, query string, arg interface{}) (int64, error)
	ExecWithReturningIds(ctx context.Context, query string, args ...interface{}) ([]int64, error)
	Commit() error
	Rollback() error
}
