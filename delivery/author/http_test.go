package author

import (
	"ThreeLayer/entities"
	"ThreeLayer/service"
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"strconv"

	"ThreeLayer/errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerAuthor_post function contains test cases to check the function which performs Handler Requests
// to add a new author instance from the database
func TestHandlerAuthor_Post(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockAuthor(ctrl)
	mock := New(mockService)

	testcases := []struct {
		desc      string
		reqBody   entities.Author
		expRes    entities.Author
		expStatus int
		expError  error
	}{
		{"Valid details",
			entities.Author{FirstName: "HC", LastName: "Verma", Dob: "2/11/1989",
				PenName: "Verma"},
			entities.Author{ID: 1, FirstName: "HC", LastName: "Verma", Dob: "2/11/1989",
				PenName: "Verma"},
			http.StatusCreated, nil},
		{"InValid details",
			entities.Author{FirstName: "", LastName: "Verma", Dob: "2/11/1989",
				PenName: "Verma"},
			entities.Author{},
			http.StatusBadRequest, errors.InValidDetails{Details: "FirstName"}},
	}

	for i, v := range testcases {
		mockService.EXPECT().PostAuthor(context.TODO(), v.reqBody).Return(v.expRes, v.expError)

		body, _ := json.Marshal(v.reqBody)
		req := httptest.NewRequest(http.MethodPost, "/author", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mock.PostAuthor(w, req)
		res, err := io.ReadAll(w.Result().Body)
		if err != nil {
			log.Print(err)
		}
		resID := entities.Author{}
		err = json.Unmarshal(res, &resID)
		if err != nil {
			log.Print(err)
		}
		if resID != v.expRes {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expStatus, w.Code)
		}
		if w.Code != v.expStatus {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expStatus, w.Code)
		}

	}
}

// TestAuthorHandler_putAuthor function contains test cases to check the function which performs Handler Requests
// to make changes to an existing author instance in the database
func TestAuthorHandler_PutAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockAuthor(ctrl)
	mock := New(mockService)

	testcases := []struct {
		desc          string
		reqID         string
		reqData       entities.Author
		expData       entities.Author
		expStatusCode int
		expError      error
	}{
		{"success case update firstname.", "1",
			entities.Author{ID: 1, FirstName: "Rahul", LastName: "Saini", Dob: "22/07/2000",
				PenName: "ABC"},
			entities.Author{ID: 1, FirstName: "Rahul", LastName: "Saini", Dob: "22/07/2000",
				PenName: "ABC"},
			http.StatusOK, nil},
		{"success case id not present.", "1000",
			entities.Author{ID: 1000, FirstName: "Ram", LastName: "lal", Dob: "22/07/2000", PenName: "ABC"},
			entities.Author{},
			http.StatusNotFound, errors.EntityNotFound{Entity: "Author", ID: 1000}},
		//{"invalid id", "id",
		//	entities.Author{ID: 1000, FirstName: "Ram", LastName: "lal", Dob: "22/07/2000", PenName: "ABC"},
		//	entities.Author{},
		//	http.StatusBadRequest, errors.EntityNotFound{"Author", 1000}},
		{"exist Already", "2",
			entities.Author{ID: 2, FirstName: "Ram", LastName: "lal", Dob: "22/07/2000", PenName: "ABC"},
			entities.Author{},
			http.StatusConflict, errors.ExistAlready{Entity: "Author"}},
	}

	for i, v := range testcases {
		mockService.EXPECT().PutAuthor(context.Background(), v.reqData.ID, v.reqData).Return(v.expData, v.expError)

		body, _ := json.Marshal(v.reqData)
		req := httptest.NewRequest(http.MethodPut, "/author/id", bytes.NewReader(body))
		resAuthor := entities.Author{}
		req = mux.SetURLVars(req, map[string]string{"id": v.reqID})
		w := httptest.NewRecorder()
		mock.PutAuthor(w, req)
		res, err := io.ReadAll(w.Result().Body)

		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal(res, &resAuthor)
		if err != nil {
			log.Print(err)
		}
		if w.Code != v.expStatusCode {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expStatusCode, w.Code)
		}
		if resAuthor != v.expData {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expData, resAuthor)
		}
	}

}

// TestAuthorHandler_delete function contains test cases to check the function which performs Handler Requests
// to remove an author instance from the database
func TestAuthorHandler_DeleteAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockAuthor(ctrl)
	mock := New(mockService)

	testcases := []struct {
		desc          string
		reqID         string
		expStatusCode int
		expError      error
	}{
		{"Valid Details", "1", http.StatusNoContent, nil},
		{"Author does not exists", "100", http.StatusNotFound, errors.EntityNotFound{Entity: "Author", ID: 100}},
	}
	for i, v := range testcases {
		id, err := strconv.Atoi(v.reqID)
		if err != nil {
			log.Print(err)
		}

		mockService.EXPECT().DeleteAuthor(context.Background(), id).Return(v.expError)

		req := httptest.NewRequest(http.MethodDelete, "/author/{id}", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.reqID})
		w := httptest.NewRecorder()

		mock.DeleteAuthor(w, req)

		if w.Code != v.expStatusCode {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, v.expStatusCode, w.Code)
		}

	}

}

//for checking the invalid id and handeling error
func TestHandler_DeleteAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockAuthor(ctrl)

	defer ctrl.Finish()
	mock := New(mockService)
	testcases := []struct {
		desc          string
		reqID         string
		expStatusCode int
	}{
		{"Invalid ID", "id", http.StatusBadRequest},
		{"Invalid ID", "true", http.StatusBadRequest},
	}
	for i, tc := range testcases {
		req := httptest.NewRequest("PUT", "/author/{id}", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.reqID})
		w := httptest.NewRecorder()

		mock.DeleteAuthor(w, req)

		if w.Code != tc.expStatusCode {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, tc.expStatusCode, w.Code)
		}
	}
}

// for handeling invalid and unmarshling error
func TestHandler_PutAuthor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := service.NewMockAuthor(mockCtrl)
	mock := New(mockService)

	testcases := []struct {
		desc      string
		reqID     string
		reqBody   []byte
		expBody   entities.Author
		expStatus int
	}{
		{"unmarshalling error", "1",
			[]byte(`{"id": 1,"FirstName":RD, "LastName": "Sharma", "DOB": "2/12/1990", "PenName": "Sharma"}`),
			entities.Author{}, http.StatusBadRequest},
		{"invalid id", "id",
			[]byte(``),
			entities.Author{}, http.StatusBadRequest},
	}

	for i, tc := range testcases {
		req := httptest.NewRequest(http.MethodPut, "/author/{id}", bytes.NewReader(tc.reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": tc.reqID})

		w := httptest.NewRecorder()

		mock.PutAuthor(w, req)

		if tc.expStatus != w.Code {
			t.Errorf("[TEST%d]Failed. Expected %v\tGot %v", i, tc.expStatus, w.Code)
		}
	}
}
