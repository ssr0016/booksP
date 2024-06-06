package booksimpl

import (
	"context"
	"practice/infra/storage/postgres"
	"practice/library/books"
	"practice/library/config"
	"time"

	"go.uber.org/zap"
)

const (
	TimestampFormat = time.RFC3339Nano
)

type service struct {
	store *store
	db    postgres.DB
	cfg   *config.Config
	log   *zap.Logger
}

func NewService(
	db postgres.DB,
	cfg *config.Config,
) books.Service {
	return &service{
		db:    db,
		store: NewStore(db),
		cfg:   cfg,
		log:   zap.L().Named("books.service"),
	}
}

func (s *service) Create(ctx context.Context, cmd *books.CreateBooksCommand) error {
	return s.db.InTransaction(ctx, func(ctx context.Context) error {
		bt, err := s.store.isBookTaken(ctx, 0, cmd.Name)
		if err != nil {
			return err
		}

		if len(bt) != 0 {
			return books.ErrBookTaken
		}

		err = s.store.create(ctx, &books.Books{
			Name:      cmd.Name,
			CreatedAt: time.Now().UTC().Format(TimestampFormat),
			UpdatedAt: time.Now().UTC().Format(TimestampFormat),
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) Update(ctx context.Context, cmd *books.UpdateBooksCommand) error {
	return s.db.InTransaction(ctx, func(ctx context.Context) error {
		now := time.Now().Format(TimestampFormat)

		bt, err := s.store.isBookTaken(ctx, cmd.ID, cmd.Name)
		if err != nil {
			return err
		}

		if len(bt) == 0 {
			return books.ErrBookNotFound
		}

		if len(bt) > 1 || (len(bt) == 1 && bt[0].ID != cmd.ID) {
			return books.ErrBookTaken
		}

		err = s.store.update(ctx, &books.Books{
			ID:        cmd.ID,
			Name:      cmd.Name,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetByID(ctx context.Context, id int64) (*books.Books, error) {
	result, err := s.store.getByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, err
	}

	return result, nil
}

func (s *service) Search(ctx context.Context, query *books.BooksQuery) (*books.BooksResult, error) {
	if query.Page <= 0 {
		query.Page = s.cfg.Pagination.Page
	}

	if query.PerPage <= 0 {
		query.PerPage = s.cfg.Pagination.PageLimit
	}

	result, err := s.store.search(ctx, query)
	if err != nil {
		return nil, err
	}

	result.Page = query.Page
	result.PerPage = query.PerPage

	return result, nil
}
