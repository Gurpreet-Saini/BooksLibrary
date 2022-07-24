package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"ThreeLayer/errors"
)

// SetStatusCode writes the status code based on the error type
func SetStatusCode(w http.ResponseWriter, method string, data interface{}, err error) {
	switch err.(type) {
	case errors.ExistAlready:
		w.WriteHeader(http.StatusConflict)
	case errors.InValidDetails:
		w.WriteHeader(http.StatusBadRequest)
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		writeSuccessResponse(method, w, data)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// writeSuccessResponse based on the method type it calls function writeResponseBody
func writeSuccessResponse(method string, w http.ResponseWriter, data interface{}) {
	switch method {
	case http.MethodPost:
		writeResponseBody(w, http.StatusCreated, data)
	case http.MethodGet:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodPut:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodDelete:
		writeResponseBody(w, http.StatusNoContent, data)
	}
}

// writeResponseBody marshals the data and writes the body which is sent to client
func writeResponseBody(response http.ResponseWriter, statusCode int, data interface{}) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	if data == nil {
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(resp)
	if err != nil {
		log.Println("error in writing response")
		return
	}
}
