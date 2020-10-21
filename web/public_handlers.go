package web

import (
	"net/http"
)

func signupHandler(container *Container, response http.ResponseWriter, request *http.Request) (err error) {
	response.Write([]byte("signup"))

	return nil
}

func notFoundHandler(container *Container, response http.ResponseWriter, request *http.Request) error {
	response.WriteHeader(http.StatusNotFound)

	return nil
}
