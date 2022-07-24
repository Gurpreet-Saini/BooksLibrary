package books

import (
	"ThreeLayer/datastore"
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	book   datastore.Book
	author datastore.Author
}

func New(b datastore.Book, a datastore.Author) Service {
	return Service{book: b, author: a}
}

const (
	LowestPubYear = 1880
	Publisher1    = "Arihanth"
	Publisher2    = "Scholastic"
	Publisher3    = "Penguin"
)

func (s Service) PostBook(ctx context.Context, book entities.Book) (entities.Book, error) {

	err := checkDetails(book)
	if err != nil {
		return entities.Book{}, err
	}

	_, err = s.author.GetAuthorByID(ctx, book.Author.ID)
	if err != nil {
		return entities.Book{}, errors.InValidDetails{Details: "Author ID "}
	}

	books, err := s.GetBook(ctx)
	if err != nil {
		return entities.Book{}, err
	}
	for i := range books {
		if book.Author.ID == books[i].Author.ID {
			return entities.Book{}, errors.ExistAlready{Entity: "Book"}
		}
	}

	return s.book.CreateBook(ctx, book)
}

func (s Service) GetBook(ctx context.Context) ([]entities.Book, error) {
	books, err := s.book.GetAllBook(ctx)
	if err != nil {
		return nil, err
	}
	title := ctx.Value(entities.Title)
	includeAuthor, _ := ctx.Value(entities.IncludeAuthor).(bool)
	if title != "" {
		books = matchDetails(books, func(book entities.Book) bool {
			return book.Title == title
		})
	}

	if includeAuthor {
		for i := range books {
			auth, err := s.author.GetAuthorByID(ctx, books[i].Author.ID)
			if err != nil {
				log.Print(err)
			}
			books[i].Author = auth
		}
	}

	return books, nil
}

func (s Service) GetBookByID(ctx context.Context, id int) (entities.Book, error) {
	book, err := s.book.GetBookByID(ctx, id)
	if err != nil {
		return entities.Book{}, err
	}
	auth, err := s.author.GetAuthorByID(ctx, book.Author.ID)
	if err != nil {
		return entities.Book{}, err
	}
	book.Author = auth
	return book, nil

}

func (s Service) PutBook(ctx context.Context, id int, book entities.Book) (entities.Book, error) {
	err := checkDetails(book)
	if err != nil {
		return entities.Book{}, err
	}
	_, err = s.author.GetAuthorByID(ctx, book.Author.ID)
	if err != nil {
		return entities.Book{}, errors.InValidDetails{Details: "Author ID"}
	}
	_, err = s.book.GetBookByID(ctx, id)
	if err != nil {
		return entities.Book{}, err
	}
	return s.book.UpdateBook(ctx, id, book)
}

func (s Service) DeleteBook(ctx context.Context, id int) error {
	_, err := s.book.GetBookByID(ctx, id)
	if err != nil {
		return err
	}
	return s.book.DeleteBook(ctx, id)
}

//<-------------functions----------->
func matchDetails(books []entities.Book, fn func(book entities.Book) bool) []entities.Book {
	count := 0
	for i := range books {
		if fn(books[i]) {
			books[count] = books[i]
			count++
		}
	}
	return books[:count]
}
func publicationCheck(p string) bool {
	return !(p == Publisher1 || p == Publisher2 || p == Publisher3)
}

func publishedDateCheck(date string) bool {
	p := strings.Split(date, "/")

	year, err := strconv.Atoi(p[2])
	if err != nil {
		return false
	}

	if year > time.Now().Year() || year < LowestPubYear {
		return false
	}

	return true
}
func checkDetails(book entities.Book) error {
	switch {
	case book.Title == "":
		return errors.InValidDetails{Details: "Title"}
	case publicationCheck(book.Publication):
		return errors.InValidDetails{Details: "Publication"}
	case !publishedDateCheck(book.PublishedDate):
		return errors.InValidDetails{Details: "PublishedDate"}
	case book.Author.ID <= 0:
		return errors.InValidDetails{Details: "Author ID"}
	default:
		return nil
	}
}
