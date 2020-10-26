package controller

import (
	"palindromex/web/container"
	"palindromex/web/dto"

	"encoding/json"
	"errors"
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

	pageData := dto.PageData{"tokenString": tokenString, "apiKey": apiKey}
	err = c.Templates["credentials.html"].Execute(w, r, pageData)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

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
	if messages == nil {
		return nil
	}

	value, _ := json.Marshal(messages)
	w.WriteHeader(http.StatusOK)
	w.Write(value)

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
	w.WriteHeader(http.StatusOK)
	w.Write(value)

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


func UIShowMessagesHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID := vars["userID"]
	getMessagesURL, err := c.Router.Get("messages").URL("userID", userID)
	createMessageURL, err := c.Router.Get("ui_create_message").URL("userID", userID)
	pageData := dto.PageData{"getMessagesURL": getMessagesURL.Path, "createMessageURL": createMessageURL}

	err = c.Templates["messages.html"].Execute(w, r, pageData)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

	return nil
}

func UICreateMessageHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	submitURL, err := c.Router.Get("signin").URL()
	pageData := dto.PageData{"submitURL": submitURL.Path}

	err = c.Templates["create_message.html"].Execute(w, r, pageData)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

	return nil
}


func UIEditMessageHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	submitURL, err := c.Router.Get("signin").URL()
	pageData := dto.PageData{"submitURL": submitURL.Path}

	err = c.Templates["edit_message.html"].Execute(w, r, pageData)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

	return nil
}