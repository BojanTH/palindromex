package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"palindromex/web/container"
	"palindromex/web/dto"

	"errors"
	"net/http"
	"strconv"
)

func SignupHandler(c *container.Container, w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method == "GET" {
		submitURL, err := c.Router.Get("signup").URL()
		if err != nil {
			log.Printf("Cant generate URL: signup")
			return StatusError{err, http.StatusInternalServerError}
		}
		redirectURL, err := c.Router.Get("signin").URL()
		if err != nil {
			log.Printf("Cant generate URL: signin")
			return StatusError{err, http.StatusInternalServerError}
		}
		pageData := dto.PageData{"submitURL": submitURL.Path, "redirectURL": redirectURL.Path}

		err = c.Templates["signup.html"].Execute(w, r, pageData)
		if err != nil {
			return StatusError{err, http.StatusInternalServerError}
		}

		return nil
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{err, http.StatusBadRequest}
	}
	if len(content) == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	credentials := dto.Credentials{}
	json.Unmarshal(content, &credentials)
	// @TODO create more advanced validation
	if credentials.Email == "" || credentials.Name == "" || credentials.Password == "" {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = StatusError{errors.New(x), http.StatusBadRequest}
			case error:
				err = StatusError{x, http.StatusBadRequest}
			default:
				err = StatusError{errors.New("Unknown error"), http.StatusBadRequest}
			}
		}
	}()

	c.UserService.CreateNewUser(&credentials)
	c.Flash.AddSuccess(w, r, "Success! Your account has been created.")
	w.WriteHeader(http.StatusCreated)

	return nil
}

func SigninHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		submitURL, err := c.Router.Get("signin").URL()
		if err != nil {
			log.Printf("Cant generate URL: signup")
			return StatusError{err, http.StatusInternalServerError}
		}
		pageData := dto.PageData{"submitURL": submitURL.Path}

		err = c.Templates["signin.html"].Execute(w, r, pageData)
		if err != nil {
			return StatusError{err, http.StatusInternalServerError}
		}

		return nil
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{err, http.StatusBadRequest}
	}
	if len(content) == 0 {
		return StatusError{errors.New("Bad request"), http.StatusBadRequest}
	}
	credentials := dto.Credentials{}
	json.Unmarshal(content, &credentials)

	user, err := c.UserService.GetUserByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		return StatusError{err, http.StatusUnauthorized}
	}

	err = SetJwtCookie(c, w, user)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

	url, _ := c.Router.Get("ui_show_messages").URL("userID", strconv.Itoa(int(user.ID)))
	response := map[string]string{"url": url.Path}
	resp, _ := json.Marshal(&response)
	w.Write(resp)

	return nil
}

func NotFoundHandler(c *container.Container, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	err := c.Templates["404.html"].Execute(w, r, nil)
	if err != nil {
		return err
	}

	return nil
}
