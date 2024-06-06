package booksimpl

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"practice/infra/storage/postgres"
	"practice/library/books"
	"practice/util/pointer"
	"strings"

	"go.uber.org/zap"
)

type store struct {
	db     postgres.DB
	logger *zap.Logger
}

func NewStore(db postgres.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("books.store"),
	}
}

func (s *store) create(ctx context.Context, entity *books.Books) error {
	return s.db.WithTransaction(ctx, func(tx postgres.Tx) error {
		rawSQL := `
			INSERT INTO books (
				name,
				created_at,
				created_by
			)
			VALUES (
				:name,
				:created_at,
				:created_by
			)
		`
		_, err := tx.NamedExec(ctx, rawSQL, entity)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) isBookTaken(ctx context.Context, id int64, name string) ([]*books.Books, error) {
	var result []*books.Books

	rawSQL := `
		SELECT 
			id,
			name
		FROM books
		WHERE
			id = ? OR
			name = ?
	`

	err := s.db.Select(ctx, &result, rawSQL, id, name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) update(ctx context.Context, entity *books.Books) error {
	return s.db.WithTransaction(ctx, func(tx postgres.Tx) error {
		rawSQL := `
			UPDATE books
			SET
				name = :name,
				updated_at = :updated_at,
				updated_by = :updated_by	
			WHERE
				id = :id
		`

		_, err := tx.NamedExec(ctx, rawSQL, entity)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) getByID(ctx context.Context, id int64) (*books.Books, error) {
	var result books.Books

	rawSQL := `
		SELECT
			*
		FROM
			books
		WHERE
			id = ? 
	`

	err := s.db.Get(ctx, &result, rawSQL, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) search(ctx context.Context, query *books.BooksQuery) (*books.BooksResult, error) {
	var (
		result = &books.BooksResult{
			Books: make([]books.Books, 0),
		}

		sql             bytes.Buffer
		whereConditions = make([]string, 0)
		whereParams     = make([]interface{}, 0)
	)

	sql.WriteString(`
		SELECT
			*
		FROM books
	
	`)

	if len(query.Name) > 0 {
		whereConditions = append(whereConditions, "name = ILIKE ?")
		whereParams = append(whereParams, pointer.LikeString(query.Name))
	}

	if len(whereConditions) > 0 {
		sql.WriteString(" WHERE " + strings.Join(whereConditions, " AND "))
	}

	sql.WriteString(" ORDER BY created_at DESC")

	count, err := s.getCount(ctx, sql, whereParams)
	if err != nil {
		return nil, err
	}

	if query.PerPage != 0 {
		offset := query.PerPage * (query.Page - 1)
		sql.WriteString(" LIMIT ? OFFSET ?")
		whereParams = append(whereParams, query.PerPage, offset)
	}

	err = s.db.Select(ctx, &result.Books, sql.String(), whereParams...)
	if err != nil {
		return nil, err
	}

	result.TotalCount = count

	return result, nil
}
func (s *store) getCount(ctx context.Context, sql bytes.Buffer, whereParams []interface{}) (int64, error) {
	var count int64

	err := s.db.Get(ctx, &count, "SELECT COUNT(id) AS count FROM ("+sql.String()+") AS t1", whereParams...)
	if err != nil {
		return count, err
	}

	return count, nil
}
