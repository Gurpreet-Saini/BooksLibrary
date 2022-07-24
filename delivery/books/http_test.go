package books

import (
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"ThreeLayer/service"
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

// TestBookHandler_GetAll function contains test cases for function to perform Handler Requests to get a
// book instance from the database

func TestBookHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockBook(ctrl)
	mock := New(mockService)
	defer ctrl.Finish()

	testcases := []struct {
		desc          string
		title         string
		includeAuthor string
		expError      error
		expRes        []entities.Book
		expStatusCode int
	}{
		{desc: "get all books", title: "", includeAuthor: "", expRes: []entities.Book{
			{ID: 1, Title: "Rahul", Publication: "Penguin", PublishedDate: "22/07/2000"},
		}, expStatusCode: http.StatusOK, expError: nil},
		{desc: "get all books with query param", title: "Rahul", includeAuthor: "", expRes: []entities.Book{
			{ID: 1, Title: "Rahul", Publication: "Penguin", PublishedDate: "22/07/2000"}},
			expStatusCode: http.StatusOK, expError: nil},
		{desc: "get all books with query param", title: "", includeAuthor: "true",
			expRes: []entities.Book{
				{ID: 1, Title: "Rahul",
					Author: entities.Author{ID: 1, FirstName: "HC", LastName: "Verma",
						Dob: "2/12/1999", PenName: "Verma"},
					Publication: "Penguin", PublishedDate: "22/07/2000"}}, expStatusCode: http.StatusOK, expError: nil},
	}
	for i, tc := range testcases {
		//ctx := context.Background()
		//ctx = context.WithValue(ctx, entities.Title, tc.title)
		//ctx = context.WithValue(ctx, entities.IncludeAuthor, tc.includeAuthor == "true")

		mockService.EXPECT().GetBook(context.TODO()).Return(tc.expRes, tc.expError)
		req := httptest.NewRequest(http.MethodGet, "/book?title="+tc.title+"&includeAuthor="+tc.includeAuthor,
			nil)
		w := httptest.NewRecorder()

		mock.GetBook(w, req)

		res, err := io.ReadAll(w.Result().Body)

		if err != nil {
			log.Print(err)
		}

		resBooks := make([]entities.Book, 0)

		err = json.Unmarshal(res, &resBooks)
		if err != nil {
			log.Print("expected error to be nil got ", err)
		}

		if w.Code != tc.expStatusCode {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, tc.expStatusCode, w.Code)
		}

		if !reflect.DeepEqual(resBooks, tc.expRes) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, tc.expRes, resBooks)
		}
	}
}

// TestBookDeliveryGetBookByID function contains test cases for function to perform Handler Requests to get a
// book instance using its ID from the database
//error
func TestBookHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockBook(ctrl)
	mock := New(mockService)
	defer ctrl.Finish()

	testcases := []struct {
		desc          string
		req           string
		expRes        entities.Book
		expStatusCode int
		expError      error
	}{
		{desc: "get book", req: "1", expRes: entities.Book{ID: 1, Title: "Rahul",
			Author:      entities.Author{ID: 1, FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"},
			Publication: "Penguin", PublishedDate: "22/07/2000"}, expStatusCode: http.StatusOK},
		{"Id doesn't exist", "1000", entities.Book{}, http.StatusNotFound, errors.EntityNotFound{Entity: "Book", ID: 1000}},
		//	{"invalid id", "id", entities.Book{}, http.StatusBadRequest, errors.EntityNotFound{"Book", 1000}},
	}
	for i, tc := range testcases {
		id, _ := strconv.Atoi(tc.req)
		mockService.EXPECT().GetBookByID(context.Background(), id).Return(tc.expRes, tc.expError)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/book/{id}", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.req})

		mock.GetBookByID(w, req)
		res, err := io.ReadAll(w.Result().Body)
		if err != nil {
			log.Printf("expected err to be nil got %v", err)
		}

		resBook := entities.Book{}

		err = json.Unmarshal(res, &resBook)
		if err != nil {
			log.Printf("expected error to be nil got %v", err)
		}

		if w.Result().StatusCode != tc.expStatusCode {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, tc.expStatusCode, w.Code)
		}

		if resBook != tc.expRes {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, tc.expRes, resBook)
		}
	}
}

// TestBookHandler_POST function contains test cases for function to perform Handler Requests to add a
// book instance to the database
func TestBookHandler_Post(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockBook(ctrl)
	mock := New(mockService)
	defer ctrl.Finish()

	testcases := []struct {
		desc      string
		reqBody   entities.Book
		expRes    entities.Book
		expStatus int
		expError  error
	}{
		{"Publication should be Scholastic/Penguin/Arihanth", entities.Book{Title: "Rahul",
			Author:      entities.Author{ID: 1},
			Publication: "Rahul", PublishedDate: "22/07/2000"}, entities.Book{},
			http.StatusBadRequest, errors.InValidDetails{Details: "Publication"}},
		{"Publication date should be between 1880 and 2022", entities.Book{Title: "",
			Author:      entities.Author{ID: 1},
			Publication: "", PublishedDate: "1/1/1600"}, entities.Book{},
			http.StatusBadRequest, errors.InValidDetails{Details: "PublishedDate"}},
		{"Author should exist", entities.Book{Title: "Rahul",
			Author:      entities.Author{ID: 2},
			Publication: "Penguin", PublishedDate: "22/07/2000"}, entities.Book{},
			http.StatusBadRequest, errors.InValidDetails{Details: "Author ID"}},
		{"Title can't be empty", entities.Book{Title: "",
			Author:      entities.Author{ID: 1},
			Publication: "", PublishedDate: ""}, entities.Book{},
			http.StatusBadRequest, errors.InValidDetails{Details: "Title"}},
	}
	for i, tc := range testcases {
		mockService.EXPECT().PostBook(context.Background(), tc.reqBody).Return(tc.expRes, tc.expError)

		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader(body))

		mock.PostBook(w, req)

		res, err := io.ReadAll(w.Result().Body)
		if err != nil {
			log.Printf("expected error to be nil got %v", err)
		}
		resBook := entities.Book{}

		err = json.Unmarshal(res, &resBook)
		if err != nil {
			log.Printf("expected error to be nil got %v", err)
		}

		if w.Code != tc.expStatus {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, tc.expStatus)
		}

		if !reflect.DeepEqual(resBook, tc.expRes) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resBook, tc.expRes)
		}
	}
}

// TestBookHandler_Put function contains test cases for function to perform Handler Requests to make changes
// to an existing book instance in the database
//error
func TestBookHandler_Put(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockBook(ctrl)
	mock := New(mockService)
	defer ctrl.Finish()

	testcases := []struct {
		desc      string
		reqID     string
		reqBody   entities.Book
		expStatus int
		expError  error
	}{
		{desc: "invalid case id not exist", reqID: "1000", reqBody: entities.Book{ID: 1000, Title: "title1", Author: entities.Author{ID: 9},
			Publication: "Arihanth", PublishedDate: "18/08/2018"}, expStatus: http.StatusNotFound, expError: errors.EntityNotFound{Entity: "Book", ID: 1000}},
		{"Invalid book name.", "1", entities.Book{ID: 1, Author: entities.Author{ID: 1},
			Publication: "Oxford", PublishedDate: "21/04/1985"}, http.StatusBadRequest, errors.InValidDetails{Details: "Title"}},
		//{"invalid id", "id", entities.Book{}, http.StatusBadRequest, errors.InValidDetails{"id"}},
	}
	for i, tc := range testcases {
		id, err := strconv.Atoi(tc.reqID)
		if err != nil {
			log.Print(err)
		}
		mockService.EXPECT().PutBook(context.Background(), id, tc.reqBody).Return(entities.Book{}, tc.expError)
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPut, "/book/{id}", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": tc.reqID})

		mock.PutBook(w, req)

		if w.Code != tc.expStatus {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, tc.expStatus)
		}
	}
}

// TestBookHandler_DeleteBook function contains test cases for function to perform Handler Requests to remove a
// book instance from the database
//error
func TestBookHandler_DeleteBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockBook(ctrl)
	mock := New(mockService)
	defer ctrl.Finish()
	testcases := []struct {
		desc      string
		reqID     string
		expStatus int
		expError  error
	}{
		{"Valid Details", "1", http.StatusNoContent, nil},
		{"Book does not exists", "100", http.StatusNotFound, errors.EntityNotFound{Entity: "Book", ID: 100}},
		//{"Invalid id", "id", http.StatusBadRequest}, //passing as a character
	}
	for i, tc := range testcases {
		id, err := strconv.Atoi(tc.reqID)
		if err != nil {
			log.Print(err)
		}
		mockService.EXPECT().DeleteBook(context.Background(), id).Return(tc.expError)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/book/{id}", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.reqID})
		mock.DeleteBook(w, req)

		if w.Code != tc.expStatus {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, tc.expStatus)
		}
	}
}
