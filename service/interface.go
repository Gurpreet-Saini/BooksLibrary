package service

import (
	"ThreeLayer/entities"
	"context"
)

type Book interface {
	GetBook(ctx context.Context) ([]entities.Book, error)
	GetBookByID(ctx context.Context, id int) (entities.Book, error)
	PostBook(ctx context.Context, book entities.Book) (entities.Book, error)
	DeleteBook(ctx context.Context, id int) error
	PutBook(ctx context.Context, id int, book entities.Book) (entities.Book, error)
}

type Author interface {
	PostAuthor(ctx context.Context, author entities.Author) (entities.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	PutAuthor(ctx context.Context, id int, author entities.Author) (entities.Author, error)
}
