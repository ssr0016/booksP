package books

import (
	"practice/api/errors"
)

var (
	ErrBookTaken     = errors.New("books.already_taken", "Book is already taken")
	ErrBookNotFound  = errors.New("books.not_found", "Book not found")
	ErrNameInvalid   = errors.New("books.name_invalid", "Name is invalid")
	ErrBookIDInvalid = errors.New("books.id_invalid", "Book ID is invalid")
)

type Books struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	CreatedBy string `db:"created_by" json:"created_by"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
}

type CreateBooksCommand struct {
	Name string `json:"name"`
}

type UpdateBooksCommand struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UpdatedBy string
}

type BooksQuery struct {
	Name    string `schema:"name"`
	Page    int    `schema:"page"`
	PerPage int    `schema:"per_page"`
}

type BooksResult struct {
	TotalCount int64   `json:"total_count"`
	Books      []Books `json:"result"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
}

func (c *CreateBooksCommand) Validate() error {
	if c.Name == "" {
		return ErrNameInvalid
	}
	return nil
}

func (c *UpdateBooksCommand) Validate() error {
	if c.ID <= 0 {
		return ErrBookIDInvalid
	}

	if c.Name == "" {
		return ErrNameInvalid
	}

	return nil
}
