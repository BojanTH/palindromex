package web

import (
	"fmt"
	"net/http"
)

type Error interface {
	error
	StatusCode() int
}

type StatusError struct {
	Err error
	Code int
}

// Error is used so that StatuError satisfies the error interface
func (se StatusError) Error() string {
	return se.Err.Error()
}

// StatusCode will return the status code
func (se StatusError) StatusCode() int {
	return se.Code
}

type Handler struct {
	Container *Container
	Handle func(response http.ResponseWriter, request *http.Request) error
}

func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := h.Handle(response, request)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// @TODO: add logging
			http.Error(response, e.Error(), e.StatusCode())
		default:
			fmt.Println(e)
			// @TODO: template errors end up here, they should be better handled
			http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
