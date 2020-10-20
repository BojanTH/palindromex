package web

import (
	"net/http"
	"log"
)

var (
	AppPort string
)

func Make() {
	c := NewContainer()

	c.Router.Handle("/signup", Handler{c, signupHandler}).
		Methods("GET", "POST").
		Name("signup")

	// Redirect everything else to 404
	c.Router.NotFoundHandler = Handler{c, notFoundHandler}

	// Serve
	log.Fatal(http.ListenAndServe(AppPort, c.Router))
}