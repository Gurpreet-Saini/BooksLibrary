package author

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

// testAuthorStore_CreateAuthor contains test cases for function to perform DB execution to add a new author instance in database
func TestStorer_CreateAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		reqBody        entities.Author
		expRes         entities.Author
		lastInsertedID int64
		expErr         error
	}{
		{
			"Error Case",
			entities.Author{},
			entities.Author{},
			0,
			fmt.Errorf("query error"),
		},
		{
			"Success Case",
			entities.Author{FirstName: "MG", LastName: "Verma", Dob: "13/07/2000", PenName: "Verma"},
			entities.Author{ID: 1, FirstName: "MG", LastName: "Verma", Dob: "13/07/2000", PenName: "Verma"},
			1,
			nil,
		},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)
		mock.ExpectExec(datastore.InsertAuthor).
			WithArgs(v.reqBody.FirstName, v.reqBody.LastName, v.reqBody.Dob, v.reqBody.PenName).
			WillReturnResult(sqlmock.NewResult(v.lastInsertedID, 0)).
			WillReturnError(v.expErr)

		resp, err := a.CreateAuthor(context.Background(), v.reqBody)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if resp != v.expRes {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.expRes)
		}
	}
}

// testAuthorStorer_PutAuthor contains test cases for function
// to perform required DB Queries to edit an author instance in database
func TestStorer_PutAuthor(t *testing.T) {
	testcases := []struct {
		desc         string
		reqID        int
		reqData      entities.Author
		expData      entities.Author
		rowsAffected int64
		expErr       error
	}{
		{
			"Valid case update firstname.", 1,
			entities.Author{ID: 1, FirstName: "Rahul", LastName: "Saini", Dob: "22/07/2000", PenName: "ABC"},
			entities.Author{ID: 1, FirstName: "Rahul", LastName: "Saini", Dob: "22/07/2000", PenName: "ABC"},
			1, nil,
		},
		{
			"Error Case", 2,
			entities.Author{},
			entities.Author{},
			0, fmt.Errorf("query error"),
		},
	}

	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)
		mock.ExpectExec(datastore.UpdateAuthor).
			WithArgs(v.reqData.FirstName, v.reqData.LastName, v.reqData.Dob, v.reqData.PenName, v.reqID).
			WillReturnResult(sqlmock.NewResult(0, v.rowsAffected)).
			WillReturnError(v.expErr)

		res, err := a.PutAuthor(context.Background(), v.reqID, v.reqData)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if res != v.expData {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, res, v.expData)
		}
	}
}

// testAuthorStorer_DeleteAuthor contains test cases for function to perform required DB Queries to
// remove an author instance from the database

func TestStorer_DeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc   string
		reqID  int
		rowAff int64
		expErr error
	}{
		{desc: "Success Case", reqID: 1, rowAff: 1},
		{desc: "Error Case", reqID: 99,
			expErr: errors.EntityNotFound{Entity: "Author", ID: 99}},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)
		mock.ExpectExec(datastore.DeleteAuthor).
			WithArgs(v.reqID).
			WillReturnResult(sqlmock.NewResult(0, v.rowAff))

		err := a.DeleteAuthor(context.Background(), v.reqID)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}
	}
}

//testAuthorStorer_GetAuthor contains test cases for function to perform required DB Queries to
// remove an author instance from the database
func TestAuthorStorer_GetAuthor(t *testing.T) {
	testcases := []struct {
		desc    string
		expRows *sqlmock.Rows
		expRes  []entities.Author
		expErr  error
	}{
		{
			desc: "get all books",
			expRows: sqlmock.NewRows([]string{"id", "first_name", "last_name", "dob", "pen_name"}).
				AddRow(1, "MG", "Verma", "13/07/2000", "Verma"),
			expRes: []entities.Author{{ID: 1, FirstName: "MG", LastName: "Verma", Dob: "13/07/2000",
				PenName: "Verma"}},
		},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(datastore.GetAuthor).WillReturnRows(v.expRows).WillReturnError(v.expErr)
		resp, err := a.GetAuthor(context.Background())

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if !reflect.DeepEqual(resp, v.expRes) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.expRes)
		}
	}
}

func TestAuthorStore_GetAuthorByID(t *testing.T) {
	testcases := []struct {
		desc   string
		reqID  int
		expRes entities.Author
		expRow *sqlmock.Rows
		expErr error
	}{

		{"get book", 1, entities.Author{ID: 1, FirstName: "MG", LastName: "Verma", Dob: "13/07/2000",
			PenName: "Verma"},
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "dob", "pen_name"}).
				AddRow(1, "MG", "Verma", "13/07/2000", "Verma"),
			nil},
		{"Id NotFOUND", 999, entities.Author{},
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "dob",
				"pen_name"}), errors.EntityNotFound{Entity: "Author"}},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(datastore.GetByIDAuthor).WithArgs(v.reqID).WillReturnRows(v.expRow)
		mock.ExpectQuery(datastore.GetByIDAuthor).WithArgs(v.reqID).WillReturnError(v.expErr)
		resp, err := a.GetAuthorByID(context.Background(), v.reqID)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, v.expErr)
		}

		if resp != v.expRes {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.expRes)
		}
	}
}
