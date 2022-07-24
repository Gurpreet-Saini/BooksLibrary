package books

import (
	"ThreeLayer/entities"
	"ThreeLayer/errors"
	"context"
	"reflect"
	"testing"
)

//<---------------------AUTHOR STORE--------------------------->
type mockAuthorStore struct {
}

func (m mockAuthorStore) GetAuthor(ctx context.Context) ([]entities.Author, error) {
	return []entities.Author{}, nil
}

func (m mockAuthorStore) GetAuthorByID(ctx context.Context, id int) (entities.Author, error) {
	if id == 2 {
		return entities.Author{}, errors.EntityNotFound{Entity: "Author", ID: 2}
	}
	if id == 9 {
		return entities.Author{}, errors.InValidDetails{Details: "Author ID"}
	}
	return entities.Author{ID: id, FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, nil
}

func (m mockAuthorStore) CreateAuthor(ctx context.Context, author entities.Author) (entities.Author, error) {
	return entities.Author{}, nil
}

func (m mockAuthorStore) PutAuthor(ctx context.Context, id int, author entities.Author) (entities.Author, error) {
	return entities.Author{}, nil
}

func (m mockAuthorStore) DeleteAuthor(ctx context.Context, id int) error {
	return nil
}
func TestServiceBook_GetBook(t *testing.T) {
	testcases := []struct {
		desc          string
		title         string
		includeAuthor string
		expResult     []entities.Book
		expErr        error
	}{
		{desc: "get all books", expResult: []entities.Book{
			{ID: 1, Title: "Rahul", Publication: "Penguin", PublishedDate: "22/07/2000", Author: entities.Author{ID: 3}},
		}},

		{desc: "get all books with query param", title: "Rahul", includeAuthor: "false", expResult: []entities.Book{
			{ID: 1, Title: "Rahul", Publication: "Penguin", PublishedDate: "22/07/2000", Author: entities.Author{ID: 3}},
		}},
		{desc: "get all books with query param", includeAuthor: "true", expResult: []entities.Book{
			{ID: 1, Title: "Rahul", Author: entities.Author{ID: 3, FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989",
				PenName: "Sharma"}, Publication: "Penguin", PublishedDate: "22/07/2000"},
		}},
	}
	for i, v := range testcases {
		a := New(mockBookStore{}, mockAuthorStore{})

		ctx := context.Background()
		ctx = context.WithValue(ctx, entities.Title, v.title)
		ctx = context.WithValue(ctx, entities.IncludeAuthor, v.includeAuthor == "true")
		output, err := a.GetBook(ctx)

		if !reflect.DeepEqual(err, v.expErr) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expErr, err)
		}

		if !reflect.DeepEqual(output, v.expResult) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expResult, output)
		}
	}
}

func TestServiceBook_GetBookByID(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		expResult entities.Book
		expErr    error
	}{

		{desc: "Book ID doesn't exist", id: 999, expResult: entities.Book{}, expErr: errors.EntityNotFound{Entity: "Book", ID: 999}},
		//{desc: "Book ID doesn't exist", id: 2, expResult: entities.Book{}, expErr: errors.EntityNotFound{Entity: "Book", ID: 2}},
	}
	for i, v := range testcases {
		a := New(mockBookStore{}, mockAuthorStore{})

		output, err := a.GetBookByID(context.Background(), v.id)
		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expErr, err)
		}

		if !reflect.DeepEqual(output, v.expResult) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expResult, output)
		}
	}
}

func TestServiceBook_PostBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqResult entities.Book
		expResult entities.Book
		expErr    error
	}{
		{desc: "Valid case", reqResult: entities.Book{Title: "Rahul", Author: entities.Author{ID: 1},
			Publication: "Arihanth", PublishedDate: "22/07/2000"},
			expResult: entities.Book{ID: 1, Title: "Rahul",
				Author: entities.Author{ID: 1}, Publication: "Arihanth",
				PublishedDate: "22/07/2000"}},
		{desc: "Already Exists", reqResult: entities.Book{Title: "Rahul",
			Author:      entities.Author{ID: 3},
			Publication: "Arihanth", PublishedDate: "22/07/2000"},
			expErr: errors.ExistAlready{Entity: "Book"}},
		{desc: "Publication should be Scholastic/Penguin/Arihanth", reqResult: entities.Book{Title: "Rahul",
			Author: entities.Author{ID: 1}, Publication: "Rahul",
			PublishedDate: "22/07/2000"},
			expErr: errors.InValidDetails{Details: "Publication"}},

		{desc: "Published date should be in between 1880 and 2022", reqResult: entities.Book{Title: "Rahul",
			Author: entities.Author{ID: 1}, Publication: "",
			PublishedDate: "1/1/1600"},
			expErr: errors.InValidDetails{Details: "PublishedDate"}},

		{desc: "Author id invalid", reqResult: entities.Book{Title: "Rahul",
			Author:        entities.Author{ID: 2},
			Publication:   "Penguin",
			PublishedDate: "22/07/2000"},
			expErr: errors.InValidDetails{Details: "Author ID "}},
		{desc: "Author id invalid", reqResult: entities.Book{Title: "Rahul",
			Author:        entities.Author{ID: 0},
			Publication:   "Penguin",
			PublishedDate: "22/07/2000"},
			expErr: errors.InValidDetails{Details: "Author ID"}},
		{desc: "Title empty", reqResult: entities.Book{Title: "",
			Author:        entities.Author{ID: 1},
			Publication:   "",
			PublishedDate: ""},
			expErr: errors.InValidDetails{Details: "Title"}},
	}
	for i, v := range testcases {
		a := New(mockBookStore{}, mockAuthorStore{})

		ctx := context.Background()
		ctx = context.WithValue(ctx, entities.Title, v.reqResult.Title)

		res, err := a.PostBook(ctx, v.reqResult)

		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expErr, err)
		}

		if res != v.expResult {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expResult, res)
		}
	}
}

func TestServiceBook_PutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqID     int
		reqResult entities.Book
		expResult entities.Book
		expErr    error
	}{

		{desc: "invalid case id not exist", reqID: 999, reqResult: entities.Book{ID: 999, Title: "title1",
			Author: entities.Author{ID: 9}, Publication: "Arihanth", PublishedDate: "18/08/2018"},
			expErr: errors.EntityNotFound{Entity: "Book", ID: 999}},
		{desc: "Invalid name.", reqID: 1, reqResult: entities.Book{ID: 1,
			Author: entities.Author{ID: 1}, Publication: "Oxford", PublishedDate: "22/07/2000"},
			expErr: errors.InValidDetails{Details: "Title"}},
		{desc: "invalid case id not found", reqID: 1, reqResult: entities.Book{ID: 1, Title: "title1",
			Author: entities.Author{ID: 9}, Publication: "Arihanth", PublishedDate: "22/07/2000"},
			expErr: errors.EntityNotFound{Entity: "Book", ID: 1}},
	}
	for i, v := range testcases {
		a := New(mockBookStore{}, mockAuthorStore{})

		resBook, err := a.PutBook(context.Background(), v.reqID, v.reqResult)
		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expErr, err)
		}

		if !reflect.DeepEqual(resBook, v.expResult) {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expResult, resBook)
		}
	}
}

func TestServiceBook_DeleteBook(t *testing.T) {
	testcases := []struct {
		desc   string
		reqID  int
		expErr error
	}{

		{
			desc: "id does not exists", reqID: 99, expErr: errors.EntityNotFound{Entity: "Book", ID: 99},
		},
	}
	for i, v := range testcases {
		a := New(mockBookStore{}, mockAuthorStore{})
		ctx := context.Background()
		err := a.DeleteBook(ctx, v.reqID)

		if !reflect.DeepEqual(v.expErr, err) {
			t.Errorf("[TEST%d]Failed. Expected %v\ttGot %v", i, v.expErr, err)
		}

	}

}

//<--------------------BookSTORE-------------------------->
type mockBookStore struct {
}

func (m mockBookStore) GetAllBook(ctx context.Context) ([]entities.Book, error) {
	return []entities.Book{{ID: 1, Title: "Rahul", Publication: "Penguin", PublishedDate: "22/07/2000",
		Author: entities.Author{ID: 3}}}, nil
}

func (m mockBookStore) GetBookByID(ctx context.Context, id int) (entities.Book, error) {
	return entities.Book{}, errors.EntityNotFound{Entity: "Book", ID: id}
}

func (m mockBookStore) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {
	if book.Publication == "Rahul" || book.Title == "" {
		return entities.Book{}, errors.InValidDetails{Details: "Title"}
	}
	if book.Publication == "Arihanth" {
		return entities.Book{ID: 1, Title: "Rahul",
			Author: entities.Author{ID: 1}, Publication: "Arihanth", PublishedDate: "22/07/2000"}, nil
	}
	return entities.Book{}, errors.InValidDetails{Details: "Author ID"}
}

func (m mockBookStore) UpdateBook(ctx context.Context, id int, book entities.Book) (entities.Book, error) {
	if id == 1 {
		return book, nil
	}
	if id == 999 {
		return entities.Book{}, errors.EntityNotFound{Entity: "Book", ID: 999}
	}

	return entities.Book{}, errors.InValidDetails{Details: "Title"}
}

func (m mockBookStore) DeleteBook(ctx context.Context, id int) error {
	if id == 1 {
		return nil
	}

	return errors.EntityNotFound{Entity: "Book", ID: 100}
}
