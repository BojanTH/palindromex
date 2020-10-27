package web

import (
	"log"
	"net/http"

	"palindromex/web/container"
	"palindromex/web/controller"
	_ "palindromex/web/docs"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/swaggo/http-swagger"
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
	// Load config/env variables
	port := viper.GetString("port")
	jwtKey := viper.GetString("jwt_key")
	sessionSecret := viper.GetString("session_secret")
	dbConnection := viper.GetString("db_connection")

	c := container.NewContainer(jwtKey, dbConnection, sessionSecret)

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


	/** Secured API paths **/
	auth := c.Router.PathPrefix("/v1/users/{userID}").Subrouter()
	auth.Use(controller.VerifyJwtCookie(c))

	auth.Handle("/messages", controller.Handler{c, controller.GetMessagesHandler}).
		Methods("GET").
		Name("messages")

	auth.Handle("/messages/{id}", controller.Handler{c, controller.GetOneMessageHandler}).
		Methods("GET").
		Name("one_messages")

	auth.Handle("/messages", controller.Handler{c, controller.CreateMessageHandler}).
		Methods("POST")

	auth.Handle("/messages/{id}", controller.Handler{c, controller.UpdateMessageHandler}).
		Methods("PUT")

	auth.Handle("/messages/{id}", controller.Handler{c, controller.DeleteMessageHandler}).
		Methods("DELETE")

	/** Secured UI paths **/
	ui := c.Router.PathPrefix("/users/{userID}").Subrouter()
	ui.Use(controller.VerifyJwtCookie(c))


	ui.Handle("/credentials", controller.Handler{c, controller.CredentialsHandler}).
		Methods("GET").
		Name("ui_credentials")

	ui.Handle("/show-messages", controller.Handler{c, controller.UIShowMessagesHandler}).
		Methods("GET").
		Name("ui_show_messages")

	ui.Handle("/create-message", controller.Handler{c, controller.UICreateMessageHandler}).
		Methods("GET").
		Name("ui_create_message")

	ui.Handle("/edit-message/{id}", controller.Handler{c, controller.UIEditMessageHandler}).
		Methods("GET").
		Name("ui_edit_message")


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
