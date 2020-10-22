package web

import (
	"palindromex/web/service"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Container struct {
	Router    *mux.Router
	Templates map[string]*service.Template
	Flash     *service.Flash
}

func NewContainer() *Container {
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

	return &Container {
		Router: router,
		Templates: templates,
		Flash:  flash,

	}
}
