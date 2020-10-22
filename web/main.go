package web

import (
	"net/http"
	"log"

	"github.com/gorilla/mux"
)

var (
	AppPort       string
	SessionSecret string
)

func Make() {
	c := NewContainer()

	c.Router.Handle("/signup", Handler{c, signupHandler}).
		Methods("GET", "POST").
		Name("signup")

	c.Router.Handle("/signin", Handler{c, signinHandler}).
		Methods("GET", "POST").
		Name("signin")

	// for files in dist dir
	c.Router.HandleFunc("/static/{file:[^/]+.(?:js|css)[?0-9]*$}", func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		http.ServeFile(response, request, "web/static/dist/"+vars["file"])
	})

	// Redirect everything else to 404
	c.Router.NotFoundHandler = Handler{c, notFoundHandler}

	// Serve
	log.Fatal(http.ListenAndServe(AppPort, c.Router))
}