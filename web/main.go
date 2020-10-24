package web

import (
	"net/http"
	"log"

	"github.com/gorilla/mux"
)

var (
	AppPort       string
	JwtKey        string
	SessionSecret string
	DbHost        string
	DbName        string
	DbUser        string
	DbPassword    string
	DbPort        string
	DbSslMode     string
)

func Make() {
	c := NewContainer()

	// Public paths
	c.Router.Handle("/signup", Handler{c, signupHandler}).
		Methods("GET", "POST").
		Name("signup")

	c.Router.Handle("/signin", Handler{c, signinHandler}).
		Methods("GET", "POST").
		Name("signin")


	// Secured paths
	auth := c.Router.PathPrefix("/v1/users/{userID}").Subrouter()
	auth.Use(VerifyJwtCookie(c))

	auth.Handle("/api-credentials", Handler{c, apiCredentialsHandler}).
		Methods("GET").
		Name("api_credentials")

	auth.Handle("/messages", Handler{c, messagesHandler}).
		Methods("GET").
		Name("messages")


	// Static file paths
	c.Router.HandleFunc("/static/{file:[^/]+.(?:js|css)[?0-9]*$}", func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		http.ServeFile(response, request, "web/static/dist/"+vars["file"])
	})

	// Redirect everything else to 404
	c.Router.NotFoundHandler = Handler{c, notFoundHandler}

	// Serve
	log.Fatal(http.ListenAndServe(AppPort, c.Router))
}