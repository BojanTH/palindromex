package web

import (
	"palindromex/web/dto"

	"strconv"
	"errors"
	"net/http"
)

func signupHandler(c *Container, w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method == "GET" {
		submitURL, err := c.Router.Get("signup").URL()
		redirectURL, err := c.Router.Get("signin").URL()
		pageData := dto.PageData{"submitURL": submitURL.Path, "redirectURL": redirectURL.Path}

		err = c.Templates["signup.html"].Execute(w, r, pageData)
		if err != nil {
			return StatusError{err, http.StatusInternalServerError}
		}

		return nil
	}

	r.ParseForm()
	credentials := dto.Credentials{Name: r.FormValue("name"), Email: r.FormValue("email"), Password: r.FormValue("password")}

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
	if err != nil {
		return StatusError{err, http.StatusBadRequest}
	}

	c.Flash.AddSuccess(w, r, "Success! Your account has been created.")
	url, _ := c.Router.Get("signin").URL()
	http.Redirect(w, r, url.String(), http.StatusFound)

	return nil
}

func signinHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		submitURL, err := c.Router.Get("signin").URL()
		pageData := dto.PageData{"submitURL": submitURL.Path}

		err = c.Templates["signin.html"].Execute(w, r, pageData)
		if err != nil {
			return StatusError{err, http.StatusInternalServerError}
		}

		return nil
	}

	r.ParseForm()
	credentials := dto.Credentials{Email: r.FormValue("email"), Password: r.FormValue("password")}

	user, err := c.UserService.GetUserByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		return StatusError{err, http.StatusUnauthorized}
	}

	credentials.UserID = user.ID
	err = SetJwtCookie(c, w, credentials)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

	url, _ := c.Router.Get("messages").URL("userID", strconv.Itoa(int(user.ID)))
	http.Redirect(w, r, url.String(), http.StatusFound)

	return nil
}

func notFoundHandler(c *Container, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	err := c.Templates["404.html"].Execute(w, r, nil)
	if err != nil {
		return err
	}

	return nil
}
