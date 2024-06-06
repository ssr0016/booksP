package dbimpl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type sqltx struct {
	sqlxtx *sqlx.Tx
}

func (gtx *sqltx) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return gtx.sqlxtx.NamedExecContext(ctx, gtx.sqlxtx.Rebind(query), arg)
}

func (gtx *sqltx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return gtx.sqlxtx.ExecContext(ctx, gtx.sqlxtx.Rebind(query), args...)
}

func (gtx *sqltx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return gtx.sqlxtx.QueryContext(ctx, gtx.sqlxtx.Rebind(query), args...)
}

func (gtx *sqltx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return gtx.sqlxtx.GetContext(ctx, dest, gtx.sqlxtx.Rebind(query), args...)
}

func (gtx *sqltx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return gtx.sqlxtx.SelectContext(ctx, dest, gtx.sqlxtx.Rebind(query), args...)
}

func (gtx *sqltx) Rollback() error {
	return gtx.sqlxtx.Rollback()
}

func (gtx *sqltx) Commit() error {
	return gtx.sqlxtx.Commit()
}

func (gtx *sqltx) driverName() string {
	return gtx.sqlxtx.DriverName()
}

func (gtx *sqltx) ExecWithReturningId(ctx context.Context, query string, args ...interface{}) (int64, error) {
	return execWithReturningId(ctx, gtx.driverName(), query, *gtx, args...)
}

func (gtx *sqltx) ExecWithReturningIds(ctx context.Context, query string, args ...interface{}) ([]int64, error) {
	return execWithReturningIds(ctx, gtx.driverName(), query, *gtx, args...)
}

func (gtx *sqltx) NamedExecReturningId(ctx context.Context, query string, arg interface{}) (int64, error) {

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return 0, err
	}
	query = gtx.sqlxtx.Rebind(query)
	var id int64
	err = gtx.sqlxtx.Get(&id, query, args...)
	if err != nil {
		return id, err
	}
	return id, nil
}

func execWithReturningId(ctx context.Context, driverName string, query string, sess sqltx, args ...interface{}) (int64, error) {
	supported := false
	var id int64
	if driverName == "postgres" {
		query = fmt.Sprintf("%s RETURNING id", query)
		supported = true
	}
	if supported {
		err := sess.Get(ctx, &id, query, args...)
		if err != nil {
			return id, err
		}
		return id, nil
	} else {
		res, err := sess.Exec(ctx, query, args...)
		if err != nil {
			return id, err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func execWithReturningIds(ctx context.Context, driverName string, query string, sess sqltx, args ...interface{}) ([]int64, error) {
	supported := false
	ids := make([]int64, 0)
	if driverName == "postgres" {
		query = fmt.Sprintf("%s RETURNING id", query)
		supported = true
	}
	if supported {
		err := sess.Select(ctx, &ids, query, args...)
		if err != nil {
			return ids, err
		}
		return ids, nil
	} else {
		res, err := sess.Exec(ctx, query, args...)
		if err != nil {
			return ids, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
