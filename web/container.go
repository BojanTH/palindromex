package web

import (
	"github.com/gorilla/mux"
)

type Container struct {
	Router *mux.Router
}

func NewContainer() *Container {
	router := mux.NewRouter()

	return &Container {
		Router: router,
	}
}
