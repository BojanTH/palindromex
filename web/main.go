package web

import (
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "palindromex/web/docs"
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

// @title PalindromeX
// @version 1.0
// @description Discover hidden world of palindromes

// @license.name BSD 2-Clause License
// @license.url https://choosealicense.com/licenses/bsd-2-clause/

// @host palindromex.ml
// @BasePath /v1/users/{userID}
// @schemes https

// @securityDefinitions.apikey ApiToken
// @in header
// @name X-Auth-Token
func Make() {
	c := NewContainer()

	// Public paths
	c.Router.Handle("/doc/{any}", httpSwagger.Handler(
		httpSwagger.URL("/doc/doc.json"),
	))

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

	auth.Handle("/messages", Handler{c, getMessagesHandler}).
		Methods("GET").
		Name("messages")

	auth.Handle("/messages/{id}", Handler{c, getOneMessageHandler}).
		Methods("GET")

	auth.Handle("/messages", Handler{c, createMessageHandler}).
		Methods("POST")

	auth.Handle("/messages/{id}", Handler{c, updateMessageHandler}).
		Methods("PUT")

	auth.Handle("/messages/{id}", Handler{c, deleteMessageHandler}).
		Methods("DELETE")


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