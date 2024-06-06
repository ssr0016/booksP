package dbimpl

import (
	"context"
	"database/sql"
	"fmt"
	"practice/infra/storage/db"

	"github.com/jmoiron/sqlx"
)

type ContextSessionKey struct{}

type sqlxdb struct {
	db *sqlx.DB
}

func NewSqlx(db *sqlx.DB) db.DB {
	return &sqlxdb{db: db}
}

func (gs *sqlxdb) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return gs.db.GetContext(ctx, dest, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return gs.db.SelectContext(ctx, dest, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return gs.db.QueryContext(ctx, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return gs.db.ExecContext(ctx, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return gs.db.NamedExecContext(ctx, gs.db.Rebind(query), arg)
}

func (gs *sqlxdb) Beginx() (db.Tx, error) {
	tx, err := gs.db.Beginx()
	return &sqltx{sqlxtx: tx}, err
}

func (gs *sqlxdb) checkSession(ctx context.Context) (db.Tx, bool, error) {
	value := ctx.Value(ContextSessionKey{})

	var tx db.Tx
	sess, ok := value.(db.Tx)
	if ok {
		return sess, false, nil
	}

	tx, err := gs.Beginx()
	if err != nil {
		return tx, false, nil
	}

	return tx, true, nil
}

func (gs *sqlxdb) inTransaction(ctx context.Context, callback func(db.Tx) error) error {
	tx, isNew, err := gs.checkSession(ctx)
	if err != nil {
		return err
	}

	err = callback(tx)

	if !isNew {
		return err
	}

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (gs *sqlxdb) WithTransaction(ctx context.Context, fn func(sess db.Tx) error) error {
	return gs.inTransaction(ctx, fn)
}

func (gs *sqlxdb) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return gs.inTransaction(ctx, func(sess db.Tx) error {
		withValue := context.WithValue(ctx, ContextSessionKey{}, sess)
		return fn(withValue)
	})
}

func (gs *sqlxdb) Ping() error {
	return gs.db.Ping()
}

func (gs *sqlxdb) Close() error {
	if err := gs.db.Close(); err != nil {
		return err
	}

	return nil
}
