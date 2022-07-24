package datastore

import (
	"ThreeLayer/entities"
	"context"
)

type Author interface {
	GetAuthor(context.Context) ([]entities.Author, error)
	GetAuthorByID(ctx context.Context, id int) (entities.Author, error)
	CreateAuthor(ctx context.Context, author entities.Author) (entities.Author, error) //post
	PutAuthor(ctx context.Context, id int, author entities.Author) (entities.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
}

type Book interface {
	GetAllBook(ctx context.Context) ([]entities.Book, error)
	GetBookByID(ctx context.Context, id int) (entities.Book, error)
	CreateBook(ctx context.Context, book entities.Book) (entities.Book, error)
	UpdateBook(ctx context.Context, id int, book entities.Book) (entities.Book, error)
	DeleteBook(ctx context.Context, id int) error
}
