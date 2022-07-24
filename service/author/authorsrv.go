package author

import (
	"ThreeLayer/datastore"
	_ "ThreeLayer/datastore/author"
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
)

type authorService struct {
	authorstore datastore.Author
	bookstore   datastore.Book
}

//dependency injection factory function
func New(author datastore.Author, book datastore.Book) authorService {
	return authorService{authorstore: author, bookstore: book}
}

func (s authorService) PostAuthor(ctx context.Context, author entities.Author) (entities.Author, error) {
	if err := checkDetails(author); err != nil {
		return entities.Author{}, err
	}
	///getting details of all authors
	authors, err := s.authorstore.GetAuthor(context.Background())
	if err != nil {
		return entities.Author{}, err
	}

	for i := range authors {
		if checkDuplicate(authors[i], author) { //checking duplicacy
			return entities.Author{}, errors.ExistAlready{Entity: "Author"}
		}
	}

	return s.authorstore.CreateAuthor(ctx, author)
}
func (s authorService) PutAuthor(ctx context.Context, id int, author entities.Author) (entities.Author, error) {
	_, err := s.authorstore.GetAuthorByID(ctx, id)
	if err != nil {
		return entities.Author{}, err
	}

	return s.authorstore.PutAuthor(ctx, id, author)
}
func (s authorService) DeleteAuthor(ctx context.Context, id int) error {
	_, err := s.authorstore.GetAuthorByID(ctx, id)
	if err != nil {
		return err
	}
	books, err := s.bookstore.GetAllBook(ctx)
	if err != nil {
		return err
	}

	books = matchBook(books, func(book entities.Book) bool {
		return book.Author.ID == id
	})

	for i := range books {
		err = s.bookstore.DeleteBook(ctx, (books[i].ID))
		if err != nil {
			return err
		}

	}
	return s.authorstore.DeleteAuthor(ctx, id)
}

//<--------------functions----------------->

// checking duplicacy
func checkDuplicate(a1, a2 entities.Author) bool {
	return a1.FirstName == a2.FirstName && a1.LastName == a2.LastName && a1.Dob == a2.Dob && a1.PenName == a2.PenName
}

// CheckDetails function is to check all the validations for an author
func checkDetails(author entities.Author) error {
	switch {
	case author.FirstName == "":
		return errors.InValidDetails{Details: "FirstName"}
	case author.LastName == "":
		return errors.InValidDetails{Details: "LastName"}
	case author.PenName == "":
		return errors.InValidDetails{Details: "PenName"}
	case author.Dob == "":
		return errors.InValidDetails{Details: "Dob"}
	default:
		return nil
	}
}

//match function is used of matching the all same book with same author id
func matchBook(books []entities.Book, fn func(book entities.Book) bool) []entities.Book {
	count := 0

	for i := range books {

		if fn(books[i]) {
			books[count] = books[i]
			count++
		}
	}
	return books[:count]
}
