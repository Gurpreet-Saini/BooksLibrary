package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"ThreeLayer/driver"

	datastoreAuthor "ThreeLayer/datastore/author"
	datastoreBook "ThreeLayer/datastore/books"
	handlerAuthor "ThreeLayer/delivery/author"
	handlerBook "ThreeLayer/delivery/books"
	serviceAuthor "ThreeLayer/service/author"
	serviceBook "ThreeLayer/service/books"
)

func main() {

	var err error

	db, err := driver.ConnectToSQL()
	if err != nil {
		log.Println("could not connect to sql, Connection Fail, err:", err)
		return
	}
	bookStore := datastoreBook.New(db)
	authorStore := datastoreAuthor.New(db)

	svcBook := serviceBook.New(bookStore, authorStore)
	svcAuthor := serviceAuthor.New(authorStore, bookStore)

	book := handlerBook.New(svcBook)
	author := handlerAuthor.New(svcAuthor)

	r := mux.NewRouter()
	r.HandleFunc("/book", book.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", book.PostBook).Methods(http.MethodPost)
	r.HandleFunc("/book/{id}", book.GetBookByID).Methods(http.MethodGet)
	r.HandleFunc("/book/{id}", book.PutBook).Methods(http.MethodPut)
	r.HandleFunc("/book/{id}", book.DeleteBook).Methods(http.MethodDelete)

	r.HandleFunc("/author", author.PostAuthor).Methods(http.MethodPost)
	r.HandleFunc("/author/{id}", author.PutAuthor).Methods(http.MethodPut)
	r.HandleFunc("/author/{id}", author.DeleteAuthor).Methods(http.MethodDelete)

	server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	fmt.Println("server stared At localhost:8000")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
