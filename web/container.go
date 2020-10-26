package web

import (
	"palindromex/web/db"
	"palindromex/web/service"
	"palindromex/web/repository"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Container is a collection on singleton objects
// All of the objects must be stateless because they are shared between requests
type Container struct {
	Connection     *db.Connection
	Router         *mux.Router
	Templates      map[string]*service.Template
	Flash          *service.Flash
	UserService    *service.User
	ApiKeyService  *service.APIKey
	MessageService *service.Message
}

func NewContainer() *Container {
	// DB
	connection := db.NewConnection(DbConnection)

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
	userService := service.NewUser(userRepository)
	apiKeyRepository := repository.NewAPIKey(connection)
	apiKeyService := service.NewAPIKey(apiKeyRepository)
	messageRepository := repository.NewMessage(connection)
	messageService := service.NewMessage(messageRepository)

	return &Container {
		Connection: connection,
		Router: router,
		Templates: templates,
		Flash:  flash,
		UserService: userService,
		ApiKeyService: apiKeyService,
		MessageService: messageService,
	}
}
