package web

import (
	"log"
	"net/http"
	"os"

	"palindromex/web/container"
	"palindromex/web/controller"
	_ "palindromex/web/docs"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

var (
	// PORT is used by Google Cloud as well
	port          string = os.Getenv("PORT")
	jwtKey        string = os.Getenv("JWT_KEY")
	sessionSecret string = os.Getenv("SESSION_SECRET")
	dbConnection  string = os.Getenv("DB_CONNECTION")
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
// @name Authorization
func Make() {
	c := container.NewContainer(dbConnection, sessionSecret, jwtKey)

	// Public paths
	c.Router.Handle("/doc/{any}", httpSwagger.Handler(
		httpSwagger.URL("/doc/doc.json"),
	))

	c.Router.Handle("/signup", controller.Handler{c, controller.SignupHandler}).
		Methods("GET", "POST").
		Name("signup")

	c.Router.Handle("/signin", controller.Handler{c, controller.SigninHandler}).
		Methods("GET", "POST").
		Name("signin")


	// Secured paths
	auth := c.Router.PathPrefix("/v1/users/{userID}").Subrouter()
	auth.Use(controller.VerifyJwtCookie(c))

	auth.Handle("/api-credentials", controller.Handler{c, controller.ApiCredentialsHandler}).
		Methods("GET").
		Name("api_credentials")

	auth.Handle("/messages", controller.Handler{c, controller.GetMessagesHandler}).
		Methods("GET").
		Name("messages")

	auth.Handle("/messages/{id}", controller.Handler{c, controller.GetOneMessageHandler}).
		Methods("GET")

	auth.Handle("/messages", controller.Handler{c, controller.CreateMessageHandler}).
		Methods("POST")

	auth.Handle("/messages/{id}", controller.Handler{c, controller.UpdateMessageHandler}).
		Methods("PUT")

	auth.Handle("/messages/{id}", controller.Handler{c, controller.DeleteMessageHandler}).
		Methods("DELETE")


	// Static file paths
	c.Router.HandleFunc("/static/{file:[^/]+.(?:js|css)[?0-9]*$}", func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		http.ServeFile(response, request, "web/static/dist/"+vars["file"])
	})

	// Redirect everything else to 404
	c.Router.NotFoundHandler = controller.Handler{c, controller.NotFoundHandler}

    if port == "" {
        port = "8080"
        log.Printf("defaulting to port %s", port)
    }
	log.Fatal(http.ListenAndServe(":" + port, c.Router))
}
