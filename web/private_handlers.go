package web

import (
	"net/http"
)

func messagesHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("messages"))

	return nil
}
