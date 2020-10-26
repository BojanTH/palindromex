package controller

import (
	"palindromex/web/container"

	"log"
	"net/http"
)

// StatusError is a custom error type with a status code
type StatusError struct {
	Err error
	Code int
}

// NewStatusError creates a new StatusError
func NewStatusError(err error, code int) error {
	return StatusError{err, code}
}

// Error is used so that StatuError satisfies the error interface
func (se StatusError) Error() string {
	return se.Err.Error()
}

// StatusCode will return the status code
func (se StatusError) StatusCode() int {
	return se.Code
}

// Handler is a wrapper around http.Handler, all requests will go trough this handler
type Handler struct {
	Container *container.Container
	Handle func(container *container.Container, response http.ResponseWriter, request *http.Request) error
}

// ServeHTTP is calling the handle method and allows us to handle all error responses in one place based on the error type
func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := h.Handle(h.Container, response, request)

	// The error handeling below is added for demonstration purposes
	if err != nil {
		switch e := err.(type) {
		case StatusError:
			log.Printf("Error code: %d Error message: %s", e.StatusCode(), e.Error())

			http.Error(response, e.Error(), e.StatusCode())
		default:
			log.Printf("Error: %s", e.Error())

			http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
