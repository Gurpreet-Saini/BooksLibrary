package books

import (
	"ThreeLayer/delivery"
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"ThreeLayer/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"

	"net/http"
	"strconv"
	"strings"
)

type BookHandler struct {
	serviceBook service.Book
}

//dependency injection
func New(book service.Book) BookHandler {
	return BookHandler{serviceBook: book}
}

// GetBook function is to perform Handler Requests to get a book instance from the database
func (a BookHandler) GetBook(response http.ResponseWriter, request *http.Request) {
	title := request.URL.Query().Get("title")
	includeAuthor := request.URL.Query().Get("includeAuthor")

	ctx := context.WithValue(request.Context(), entities.Title, strings.TrimSpace(title))

	ctx = context.WithValue(ctx, entities.IncludeAuthor, includeAuthor)

	books, err := a.serviceBook.GetBook(ctx)
	delivery.SetStatusCode(response, request.Method, books, err)
}

// GetBookByID function is to perform Handler Requests to get an author instance using its ID from the database
func (a BookHandler) GetBookByID(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		delivery.SetStatusCode(response, request.Method, nil, errors.InValidDetails{Details: "id"})
		return
	}

	book, err := a.serviceBook.GetBookByID(context.Background(), id)
	delivery.SetStatusCode(response, request.Method, book, err)
}

// PostBook function is to perform Handler Requests to add a new book instance to the database
func (a BookHandler) PostBook(response http.ResponseWriter, request *http.Request) {
	book, err := getBook(request)
	if err != nil {
		delivery.SetStatusCode(response, request.Method, book, err)
	}

	book, err = a.serviceBook.PostBook(request.Context(), book)
	delivery.SetStatusCode(response, request.Method, book, err)
}
func (a BookHandler) PutBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		delivery.SetStatusCode(w, r.Method, nil, errors.InValidDetails{Details: "id"})
		return
	}

	book, err := getBook(r)
	if err != nil {
		delivery.SetStatusCode(w, r.Method, book, err)
		return
	}

	book, err = a.serviceBook.PutBook(context.Background(), id, book)
	delivery.SetStatusCode(w, r.Method, book, err)
}

// DeleteBook function is to performs Handler Requests to remove a book instance from the database
func (a BookHandler) DeleteBook(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		delivery.SetStatusCode(response, request.Method, nil, errors.InValidDetails{Details: "id"})
		return
	}

	err = a.serviceBook.DeleteBook(context.Background(), id)
	delivery.SetStatusCode(response, request.Method, nil, err)
}

func getBook(r *http.Request) (entities.Book, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return entities.Book{}, errors.InValidDetails{Details: "body"}
	}

	var book entities.Book

	err = json.Unmarshal(body, &book)
	if err != nil {
		return entities.Book{}, errors.InValidDetails{Details: "body"}
	}

	return book, nil
}
