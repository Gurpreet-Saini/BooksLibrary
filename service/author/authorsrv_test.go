package author

import (
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
	"fmt"
	"reflect"
	"testing"
)

type mockAuthorStore struct {
}

func (m mockAuthorStore) GetAuthor(ctx context.Context) ([]entities.Author, error) {
	return []entities.Author{{ID: 2, FirstName: "MG", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"}}, nil
}
func (m mockAuthorStore) GetAuthorByID(ctx context.Context, id int) (entities.Author, error) {
	if id == 10 {
		return entities.Author{}, errors.EntityNotFound{Entity: "Author", ID: id}
	}
	return entities.Author{ID: id, FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"}, nil
}
func (m mockAuthorStore) PutAuthor(ctx context.Context, id int, author entities.Author) (entities.Author, error) {
	if author.FirstName == "" {
		return entities.Author{}, errors.InValidDetails{Details: "FirstName"}
	}
	if author.FirstName == "HC" {
		return entities.Author{FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"}, nil
	}
	return entities.Author{}, errors.EntityNotFound{Entity: "Author", ID: 100}
}

func (m mockAuthorStore) DeleteAuthor(ctx context.Context, id int) error {
	if id == 1 {
		return nil
	}
	if id == 10 {
		return nil
	}
	return errors.EntityNotFound{Entity: "Author", ID: id}
}
func (m mockAuthorStore) CreateAuthor(ctx context.Context, author entities.Author) (entities.Author, error) {
	if author.FirstName != "" {
		return entities.Author{ID: 1, FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"}, nil
	}
	return entities.Author{}, errors.InValidDetails{Details: "FirstName"}
}

//<-------------------BookStruct----------------------->
//<-------------------BookStruct----------------------->

type mockBookStore struct {
}

func (m mockBookStore) GetAllBook(ctx context.Context) ([]entities.Book, error) {
	return []entities.Book{{ID: 1, Title: "Rahul", Author: entities.Author{ID: 1},
		Publication: "Rahul", PublishedDate: "11/03/2002"}, {ID: 2, Title: "Rahul", Author: entities.Author{ID: 1},
		Publication: "Rahul", PublishedDate: "11/03/2002"}}, nil
}

func (m mockBookStore) GetBookByID(ctx context.Context, id int) (entities.Book, error) {
	return entities.Book{}, nil
}

func (m mockBookStore) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {
	return entities.Book{}, nil
}

func (m mockBookStore) UpdateBook(ctx context.Context, id int, book entities.Book) (entities.Book, error) {
	//TODO implement me
	return entities.Book{}, nil
}

func (m mockBookStore) DeleteBook(ctx context.Context, id int) error {
	if id == 3 {
		return fmt.Errorf("temp err")
	}
	return nil
}

//<------------------------------main functions--------------------------------------->
func TestServiceAuthor_PostAuthor(t *testing.T) {

	testcases := []struct {
		desc      string
		reqResult entities.Author
		expResult entities.Author
		expErr    error
	}{
		{
			"Valid details",
			entities.Author{FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"},
			entities.Author{ID: 1, FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"},
			nil,
		},

		{
			"InValid details Firstname",
			entities.Author{FirstName: "", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"},
			entities.Author{},
			errors.InValidDetails{Details: "FirstName"},
		},
		{
			"InValid details lastname",
			entities.Author{FirstName: "HC", LastName: "", Dob: "2/12/1999", PenName: "Verma"},
			entities.Author{},
			errors.InValidDetails{Details: "LastName"},
		},
		{
			"InValid details dob",
			entities.Author{FirstName: "HC", LastName: "Verma", Dob: "", PenName: "Verma"},
			entities.Author{},
			errors.InValidDetails{Details: "Dob"},
		},
		{
			"InValid details penname",
			entities.Author{FirstName: "HC", LastName: "Verma", Dob: "2/12/1999", PenName: ""},
			entities.Author{},
			errors.InValidDetails{Details: "PenName"},
		},
		{
			"exist Already",
			entities.Author{FirstName: "MG", LastName: "Verma", Dob: "2/12/1999", PenName: "Verma"},
			entities.Author{},
			errors.ExistAlready{Entity: "Author"},
		},
	}

	//ctrl := gomock.NewController(t)
	//mockAuthorStore := datastore.NewMockAuthor(ctrl)
	//mockBookStore := datastore.NewMockBook(ctrl)
	//mock := New(mockAuthorStore, mockBookStore)
	//defer ctrl.Finish()

	//ctrl := gomock.NewController(t)
	//mockAuthorStore := datastore.NewMockAuthor(ctrl)
	//mockBookStore := datastore.NewMockBook(ctrl)
	//mock := New(mockAuthorStore, mockBookStore)
	//
	//for i, v := range testcases {
	//	ctx := context.Background()
	//	ctx = context.WithValue(ctx, entities.FirstName, v.reqResult.FirstName)
	//	mockAuthorStore.EXPECT().CreateAuthor(ctx, v.reqResult).Return(v.expResult, v.expErr)
	//	resp, err := mock.PostAuthor(context.Background(), v.reqResult)
	//	if err != nil {
	//		return
	//	}
	//	assert.Equalf(t, v.expErr, err, "Actual  %v and expected error %v not equal", v.expErr, err)
	//	assert.Equalf(t, v.reqResult, resp, "Actual %v and expected body %v not equal, Test Case %d failed", v.expResult, resp, i)

	for i, v := range testcases {
		a := New(mockAuthorStore{}, mockBookStore{})

		ctx := context.Background()
		ctx = context.WithValue(ctx, entities.FirstName, v.reqResult.FirstName)
		res, err := a.PostAuthor(ctx, v.reqResult)
		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. expected error %v got %v", i, v.expErr, err)
		}

		if res != v.expResult {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expResult, res)
		}
	}
}

func TestServiceAuthor_PutAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqID     int
		reqResult entities.Author
		expResult entities.Author
		expErr    error
	}{
		{
			desc: "id does not exists", reqID: 10,
			reqResult: entities.Author{ID: 10, FirstName: "Rahul", LastName: "Saini",
				Dob: "22/07/2000", PenName: "ABC"},
			expErr: errors.EntityNotFound{Entity: "Author", ID: 10},
		},
		{
			desc: "InValid details first name", reqID: 1,
			reqResult: entities.Author{FirstName: "", LastName: "Verma", Dob: "2/12/1999",
				PenName: "Verma"},
			expErr: errors.InValidDetails{Details: "FirstName"},
		},
	}
	for i, v := range testcases {
		a := New(mockAuthorStore{}, mockBookStore{})

		resBook, err := a.PutAuthor(context.Background(), v.reqID, v.reqResult)
		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expErr, err)
		}

		if !reflect.DeepEqual(resBook, v.expResult) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expResult, resBook)
		}
	}
}

func TestServiceAuthor_DeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc   string
		reqID  int
		expErr error
	}{
		{"Valid Details", 1, nil},
		{desc: "Author does not exists", reqID: 100, expErr: errors.EntityNotFound{Entity: "Author", ID: 100}},
	}

	for i, v := range testcases {

		a := New(mockAuthorStore{}, mockBookStore{})
		ctx := context.Background()
		err := a.DeleteAuthor(ctx, v.reqID)
		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expErr, err)
		}
	}
}
