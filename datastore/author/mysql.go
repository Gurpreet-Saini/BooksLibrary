package author

import (
	"ThreeLayer/datastore"
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
	"database/sql"
)

type Storer struct {
	db *sql.DB
}

func New(db *sql.DB) Storer {
	return Storer{db: db}
}

// get the list of authors
func (a Storer) GetAuthor(ctx context.Context) ([]entities.Author, error) {

	rows, err := a.db.QueryContext(ctx, datastore.GetAuthor)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	authors := make([]entities.Author, 0)

	for rows.Next() {
		var author entities.Author

		err = rows.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

//getAuthorByID  is used in book for checking the exixting author
func (a Storer) GetAuthorByID(ctx context.Context, id int) (entities.Author, error) {

	var author entities.Author

	err := a.db.QueryRowContext(ctx, datastore.GetByIDAuthor, id).Scan(&author.ID, &author.FirstName, &author.LastName,
		&author.Dob, &author.PenName)
	if err != nil {
		return entities.Author{}, errors.EntityNotFound{Entity: "Author"}
	}

	return author, nil
}

// PostAuthor function is to perform DB execution to add a new author instance in database
func (a Storer) CreateAuthor(ctx context.Context, author entities.Author) (entities.Author, error) {

	res, err := a.db.ExecContext(ctx, datastore.InsertAuthor,
		author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil {
		return entities.Author{}, err
	}

	id, _ := res.LastInsertId()

	author.ID = int(id)

	return author, nil
}

// PutAuthor function is to perform required DB Queries to edit an author instance in database.
func (a Storer) PutAuthor(ctx context.Context, id int, author entities.Author) (entities.Author, error) {

	_, err := a.db.ExecContext(ctx, datastore.UpdateAuthor, author.FirstName, author.LastName, author.Dob, author.PenName, id)
	if err != nil {
		return entities.Author{}, err
	}

	return author, nil
}

// DeleteAuthor function is to perform required DB Queries to remove an author instance from database.
func (a Storer) DeleteAuthor(ctx context.Context, id int) error {
	res, err := a.db.ExecContext(ctx, datastore.DeleteAuthor, id)
	r, _ := res.RowsAffected()
	if int(r) == 0 || err != nil {
		return errors.EntityNotFound{Entity: "Author", ID: id}
	}

	return nil

}
