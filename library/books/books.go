package books

import "context"

type Service interface {
	Search(ctx context.Context, query *BooksQuery) (*BooksResult, error)
	Create(ctx context.Context, cmd *CreateBooksCommand) error
	Update(ctx context.Context, cmd *UpdateBooksCommand) error
	GetByID(ctx context.Context, id int64) (*Books, error)
}
