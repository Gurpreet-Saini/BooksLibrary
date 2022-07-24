package author

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
)

type Handler struct {
	service service.Author
}

//dependency injection
func New(author service.Author) Handler {
	return Handler{service: author}
}

// PutAuthor function is to perform Handler Requests to make changes to an existing author instance in the database
func (a Handler) PutAuthor(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		delivery.SetStatusCode(w, r.Method, nil, errors.InValidDetails{Details: "id"})
		return
	}

	author, err := getAuthor(r)
	if err != nil {
		delivery.SetStatusCode(w, r.Method, author, err)
		return
	}

	author, err = a.service.PutAuthor(context.Background(), id, author)

	delivery.SetStatusCode(w, r.Method, author, err)
}

// PostAuthor function is to perform Handler Requests add new author instance in the database
func (a Handler) PostAuthor(w http.ResponseWriter, r *http.Request) {
	author, err := getAuthor(r)
	if err != nil {
		delivery.SetStatusCode(w, r.Method, author, err)
	}

	author, err = a.service.PostAuthor(r.Context(), author)
	delivery.SetStatusCode(w, r.Method, author, err)
}

// DeleteAuthor function is to perform Handler Requests remove an author instance from the database

func (a Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		delivery.SetStatusCode(w, r.Method, nil, errors.InValidDetails{Details: "id"})
		return
	}

	err = a.service.DeleteAuthor(context.Background(), id)
	delivery.SetStatusCode(w, r.Method, nil, err)
}
func getAuthor(r *http.Request) (entities.Author, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return entities.Author{}, errors.InValidDetails{Details: "body"}
	}

	var author entities.Author

	err = json.Unmarshal(body, &author)
	if err != nil {
		return entities.Author{}, errors.InValidDetails{Details: "body"}
	}

	return author, nil
}
