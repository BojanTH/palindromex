package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func messagesHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("messages"))

	return nil
}


func apiCredentialsHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])

	user := c.UserService.GetUserByID(userID)
	if user.ID == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	apiKey, tokenString := GetAPICredentials(c, w, user)
	message := fmt.Sprintf(
		"A new API token has been successfully created: '%s'. " +
		"This is a permanent token, keep it safe. In case you would like to disable this token, please disable the associated API key: '%s'.",
		tokenString,
		apiKey,
	)

	w.Write([]byte(message))

	return nil
}
