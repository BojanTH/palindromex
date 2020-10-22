package web

import (
	"palindromex/web/dto"

	"net/http"
)

func signupHandler(c *Container, response http.ResponseWriter, request *http.Request) (err error) {
	if request.Method == "GET" {
		submitURL, err := c.Router.Get("signup").URL()
		redirectURL, err := c.Router.Get("signin").URL()
		pageData := dto.PageData{"submitURL": submitURL.Path, "redirectURL": redirectURL.Path}

		err = c.Templates["signup.html"].Execute(response, request, pageData)
		if err != nil {
			return StatusError{err, http.StatusInternalServerError}
		}

		return nil
	}

	return nil
}

func signinHandler(container *Container, response http.ResponseWriter, request *http.Request) error {
	response.Write([]byte("Signin"))

	return nil
}

func notFoundHandler(container *Container, response http.ResponseWriter, request *http.Request) error {
	response.WriteHeader(http.StatusNotFound)

	return nil
}
