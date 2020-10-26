package service

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

const sessionName = "flash-session"
const typeError = "danger";
const typeWarning = "warning";
const typeSuccess = "success";

type Flasher interface {
	SetFlashes(map[string][]string)
}

type Flash struct {
	cookieStore *sessions.CookieStore
}

func NewFlash(cookieStore *sessions.CookieStore) *Flash {
	return &Flash{cookieStore}
}

func (f *Flash) AddError(w http.ResponseWriter, r *http.Request, messageValue string) {
	f.addFlash(w, r, typeError, messageValue)
}

func (f *Flash) AddWarning(w http.ResponseWriter, r *http.Request, messageValue string) {
	f.addFlash(w, r, typeWarning, messageValue)
}

func (f *Flash) AddSuccess(w http.ResponseWriter, r *http.Request, messageValue string) {
	f.addFlash(w, r, typeSuccess, messageValue)
}

func (f *Flash) addFlash(w http.ResponseWriter, r *http.Request, messageType string, messageValue string) {
	session, err := f.cookieStore.Get(r, sessionName)
	if err != nil {
		f.RemoveSessionCookie(w)
	}
	session.AddFlash(messageValue, messageType)
	session.Save(r, w)
}

func (f *Flash) GetFlashes(w http.ResponseWriter, r *http.Request) map[string][]string {
	session, err := f.cookieStore.Get(r, sessionName)
	if err != nil {
		f.RemoveSessionCookie(w)
	}

	messages := make(map[string][]string)
	for _, messageType := range []string{typeError, typeWarning, typeSuccess} {
	    for _, value := range session.Flashes(messageType) {
	        if messageValue, ok := value.(string); ok {
				messages[messageType] = append(messages[messageType], messageValue)
	        }
	    }
	}
	session.Save(r, w)

	return messages
}

func (f *Flash) RemoveSessionCookie(w http.ResponseWriter) {
	newCookie := http.Cookie {
		Name: sessionName,
		Value: "",
		Path: "/",
		Expires: time.Unix(0,0),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)
}