package web

import (
	"palindromex/web/db"
	"palindromex/web/service"
	"palindromex/web/repository"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Container struct {
	Connection    *db.Connection
	Router        *mux.Router
	Templates     map[string]*service.Template
	Flash         *service.Flash
	UserService   *service.User
	ApiKeyService *service.APIKey
}

func NewContainer() *Container {
	// DB
	connection := db.NewConnection(DbHost, DbUser, DbPassword, DbName, DbPort, DbSslMode)

	// Router
	router := mux.NewRouter()

	// Flash
	cookieStore := sessions.NewCookieStore([]byte(SessionSecret))
	cookieStore.Options = &sessions.Options{SameSite: http.SameSiteStrictMode, Path: "/"}
	flash := service.NewFlash(cookieStore)

	// Templates
	templates, err := service.GetTemplates(flash)
	if err != nil {
		panic(err)
	}

	// Services
	userRepository := repository.NewUser(connection)
	userService := service.NewUser(connection, userRepository)
	apiKeyRepository := repository.NewAPIKey(connection)
	apiKeyService := service.NewAPIKey(apiKeyRepository)

	return &Container {
		Connection: connection,
		Router: router,
		Templates: templates,
		Flash:  flash,
		UserService: userService,
		ApiKeyService: apiKeyService,
	}
}
