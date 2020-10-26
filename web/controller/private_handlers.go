package controller

import (
	"palindromex/web/container"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ApiCredentialsHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
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
//
// @Summary Retrieves messages that belong to a specified user
// @Produce json
// @Param userID path integer true "userID"
// @Security ApiToken
// @Success 200 {object} []dto.Message
// @Router /messages [get]
func GetMessagesHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])

	messages, err := c.MessageService.FindAllByUserID(userID)
	if err != nil {
		return err
	}

	value, _ := json.Marshal(messages)
	w.Write(value)
	w.WriteHeader(http.StatusOK)

	return nil
}


// @Summary Retrieves one message
// @Produce json
// @Param userID path integer true "userID"
// @Param messageID path integer true "messageID"
// @Security ApiToken
// @Success 200 {object} dto.Message
// @Router /messages/{mesageID} [get]
func GetOneMessageHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])
	messageID, _ := strconv.Atoi(vars["id"])

	message, err := c.MessageService.FindMessage(userID, messageID)
	if err != nil {
		return NewStatusError(err, http.StatusNotFound)
	}

	value, _ := json.Marshal(message)
	w.Write(value)
	w.WriteHeader(http.StatusOK)

	return nil
}


// @Summary Creates a new message
// @Param userID path integer true "userID"
// @Param message body string true "Message (palindrome text)"
// @Accept plain
// @Security ApiToken
// @Success 201
// @Failure 400
// @Router /messages [post]
func CreateMessageHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
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

	err = c.MessageService.CreateNewMessage(user, string(content))
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}


// @Summary Updates existing message
// @Param userID path integer true "userID"
// @Param messageID path integer true "messageID"
// @Param message body string true "Message (palindrome text)"
// @Accept plain
// @Security ApiToken
// @Success 200
// @Failure 400
// @Router /messages [put]
func UpdateMessageHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])
	messageID, _ := strconv.Atoi(vars["id"])
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{err, http.StatusBadRequest}
	}
	if len(content) == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	err = c.MessageService.UpdateMessage(userID, messageID, string(content))
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}


// @Summary Deletes existing message
// @Param userID path integer true "userID"
// @Param messageID path integer true "messageID"
// @Security ApiToken
// @Success 201
// @Failure 404
// @Router /messages [delete]
func DeleteMessageHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])
	messageID, _ := strconv.Atoi(vars["id"])

	err := c.MessageService.DeleteMessage(userID, messageID)
	if err != nil {
		return StatusError{err, http.StatusNotFound}
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
