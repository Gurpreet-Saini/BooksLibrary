package book

import (
	"ThreeLayer/datastore"
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

// testGetAllBook contains test cases for function to perform DB Executions to get a book instance from the database
func TestStorer_GetAllBook(t *testing.T) {
	testcases := []struct {
		desc    string
		expRows *sqlmock.Rows
		expRes  []entities.Book
		expErr  error
	}{
		{
			desc: "get all books",
			expRows: sqlmock.NewRows([]string{"id", "title", "publication", "publication_date",
				"author_id"}).AddRow(1, "Rahul", "Penguin", "22/07/2000", 1),
			expRes: []entities.Book{{ID: 1, Title: "Rahul", Publication: "Penguin",
				PublishedDate: "22/07/2000",
				Author:        entities.Author{ID: 1}}},
		},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(datastore.GetBook).WillReturnRows(v.expRows)
		resp, err := a.GetAllBook(context.Background())

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if !reflect.DeepEqual(resp, v.expRes) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.expRes)
		}
	}
}

// testGetByBookID contains test cases for function to perform DB Executions to get a book instance using its ID
// from the database
func TestStorer_GetBookByID(t *testing.T) {
	testcases := []struct {
		desc   string
		reqID  int
		expRes entities.Book
		expRow *sqlmock.Rows
		expErr error
	}{
		{desc: "get book", reqID: 1, expRes: entities.Book{ID: 1, Title: "Rahul",
			Publication: "Penguin", PublishedDate: "22/07/2000", Author: entities.Author{ID: 1}},
			expRow: sqlmock.NewRows([]string{"id", "title", "publication", "publication_date",
				"author_id"}).AddRow(1, "Rahul", "Penguin", "22/07/2000", 1)},
		{desc: "Id doesn't exist", reqID: 1000, expRow: sqlmock.NewRows([]string{"id", "title",
			"publication", "publication_date",
			"author_id"}), expErr: errors.EntityNotFound{Entity: "Book"}},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(datastore.GetByIDBook).WithArgs(v.reqID).WillReturnRows(v.expRow)
		mock.ExpectQuery(datastore.GetByIDBook).WithArgs(v.reqID).WillReturnError(v.expErr)
		resp, err := a.GetBookByID(context.Background(), v.reqID)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if resp != v.expRes {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.expRes)
		}
	}
}

func TestStorer_CreateBook(t *testing.T) {
	testcases := []struct {
		desc         string
		reqBody      entities.Book
		expRes       entities.Book
		lastInsertID int64
		expErr       error
	}{
		{
			"Valid Details",
			entities.Book{Title: "Rahul", Author: entities.Author{ID: 1},
				Publication: "Penguin", PublishedDate: "22/07/2000"},
			entities.Book{ID: 1, Title: "Rahul", Author: entities.Author{ID: 1},
				Publication: "Penguin", PublishedDate: "22/07/2000"},
			1, nil,
		},
		{
			"Error Case",
			entities.Book{},
			entities.Book{},
			0, fmt.Errorf("query error"),
		},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)
		mock.ExpectExec(datastore.InsertBook).
			WithArgs(v.reqBody.Title, v.reqBody.Publication, v.reqBody.PublishedDate, v.reqBody.Author.ID).
			WillReturnResult(sqlmock.NewResult(v.lastInsertID, 0)).
			WillReturnError(v.expErr)

		resp, err := a.CreateBook(context.Background(), v.reqBody)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if resp != v.expRes {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.expRes)
		}
	}
}

// testUpdateBook contains test cases for function to perform DB Executions to make changes to
// a book instance in the database
func TestStorer_UpdateBook(t *testing.T) {
	testcases := []struct {
		desc        string
		reqID       int
		reqBody     entities.Book
		expBody     entities.Book
		affectedRow int64
		expErr      error
	}{
		{
			desc: "valid case id exist", reqID: 1,
			reqBody: entities.Book{ID: 1, Title: "title", Author: entities.Author{ID: 1},
				Publication: "Arihanth", PublishedDate: "22/08/1999"},
			expBody: entities.Book{ID: 1, Title: "title", Author: entities.Author{ID: 1},
				Publication: "Arihanth", PublishedDate: "22/08/1999"},
			affectedRow: 1,
		},
		{
			desc: "error case", reqID: 1,
			affectedRow: 1, expErr: fmt.Errorf("query error"),
		},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(datastore.UpdateBook).
			WithArgs(v.reqBody.Title, v.reqBody.Publication, v.reqBody.PublishedDate, v.reqBody.ID, v.reqID).
			WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(v.expErr)

		res, err := a.UpdateBook(context.Background(), testcases[i].reqID, testcases[i].reqBody)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, testcases[i].expErr)
		}

		if res != v.expBody {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, res, testcases[i].expBody)
		}
	}
}

// testDeleteBook contains test cases for function to perform DB Executions to remove a
// book instance from the database
func TestStorer_DeleteBook(t *testing.T) {
	testcases := []struct {
		desc        string
		reqID       int
		affectedRow int64
		expErr      error
	}{
		{"Valid Details", 1, 1, nil},
		{"Book does not exists", 10, 0,
			errors.EntityNotFound{Entity: "Book", ID: 10}},
	}
	for i, tc := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(datastore.DeleteBook).
			WithArgs(tc.reqID).
			WillReturnResult(sqlmock.NewResult(0, tc.affectedRow)).
			WillReturnError(nil)

		err := a.DeleteBook(context.Background(), tc.reqID)

		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, tc.expErr)
		}
	}
}
