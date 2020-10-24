package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func apiCredentialsHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])

	user := c.UserService.GetUserByID(userID)
	if user.ID == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	apiKey, tokenString, err := GetAPICredentials(c, w, user)
	if err != nil {
		return NewStatusError(err, http.StatusInternalServerError)
	}

	message := fmt.Sprintf(
		"A new API token has been successfully created: '%s'. " +
		"This is a permanent token, keep it safe. In case you would like to disable this token, please disable the associated API key: '%s'.",
		tokenString,
		apiKey,
	)

	w.Write([]byte(message))

	return nil
}


// @TODO add pagination and limit
func getMessagesHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])

	messages, err := c.MessageService.FindAllByUserID(userID)
	if err != nil {
		return err
	}

	value, _ := json.Marshal(messages)
	w.Write(value)

	return nil
}


func getOneMessageHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])
	messageID, _ := strconv.Atoi(vars["id"])

	message := c.MessageService.FindByUserIDAndID(userID, messageID)
	if message.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	value, _ := json.Marshal(message)
	w.Write(value)

	return nil
}


func createMessageHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{err, http.StatusBadRequest}
	}
	if len(content) == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	user := c.UserService.GetUserByID(userID)
	if user.ID == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	return c.MessageService.CreateNewMessage(user, string(content))
}
