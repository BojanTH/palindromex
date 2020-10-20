package web

import (
	"net/http"
)

func signupHandler(w http.ResponseWriter, r *http.Request) (err error) {
	w.Write([]byte("test"))

	return nil
}

func notFoundHandler(response http.ResponseWriter, request *http.Request) error {
	response.WriteHeader(http.StatusNotFound)

	return nil
}
