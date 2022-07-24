package book

import (
	"ThreeLayer/datastore"
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
	"database/sql"
	"log"
)

type Storer struct {
	db *sql.DB
}

func New(db *sql.DB) Storer {
	return Storer{db: db}
}

// GetALLBook function is to perform DB Queries to get one or multiple book instances from database
func (a Storer) GetAllBook(ctx context.Context) ([]entities.Book, error) {
	rows, err := a.db.QueryContext(ctx, datastore.GetBook)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("error in closing: %v", err)
		}
	}(rows)

	books := make([]entities.Book, 0)

	for rows.Next() {
		var book entities.Book

		err = rows.Scan(&book.ID, &book.Title, &book.Publication, &book.PublishedDate, &book.Author.ID)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

// GetBookByID function is to perform DB Queries to get a particular book instance using its ID number from database
func (a Storer) GetBookByID(ctx context.Context, id int) (entities.Book, error) {

	var book entities.Book

	err := a.db.QueryRowContext(ctx, datastore.GetByIDBook, id).Scan(&book.ID, &book.Title, &book.Publication, &book.PublishedDate,
		&book.Author.ID)
	if err != nil {
		return entities.Book{}, errors.EntityNotFound{Entity: "Book"}
	}

	return book, err
}

// CreateBook function is to perform DB Executions to add new book instance in the database
func (a Storer) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {

	res, err := a.db.ExecContext(ctx, datastore.InsertBook, book.Title, book.Publication, book.PublishedDate, book.Author.ID)
	if err != nil {
		return entities.Book{}, err
	}

	id, _ := res.LastInsertId()

	book.ID = int(id)

	return book, nil
}

// Updatebook function is to perform required DB Queries to make changes to a book instance in database
func (a Storer) UpdateBook(ctx context.Context, id int, book entities.Book) (entities.Book, error) {
	_, err := a.db.ExecContext(ctx, datastore.UpdateBook, book.Title, book.Publication, book.PublishedDate, book.Author.ID, id)
	if err != nil {
		return entities.Book{}, err
	}

	book.ID = id

	return book, nil
}

// DeleteBook function is to perform DB Queries to get a particular book instance using its ID number from database
func (a Storer) DeleteBook(ctx context.Context, id int) error {
	res, err := a.db.ExecContext(ctx, datastore.DeleteBook, id)
	r, _ := res.RowsAffected()
	if r == 0 || err != nil {
		return errors.EntityNotFound{Entity: "Book", ID: id}
	}

	return nil
}
