package service

import (
	"net/http"

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

func (f *Flash) AddError(response http.ResponseWriter, request *http.Request, messageValue string) {
	f.addFlash(response, request, typeError, messageValue)
}

func (f *Flash) AddWarning(response http.ResponseWriter, request *http.Request, messageValue string) {
	f.addFlash(response, request, typeWarning, messageValue)
}

func (f *Flash) AddSuccess(response http.ResponseWriter, request *http.Request, messageValue string) {
	f.addFlash(response, request, typeSuccess, messageValue)
}

func (f *Flash) addFlash(response http.ResponseWriter, request *http.Request, messageType string, messageValue string) {
	session, err := f.cookieStore.Get(request, sessionName)
	if err != nil {
	  http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	session.AddFlash(messageValue, messageType)
	session.Save(request, response)
}

func (f *Flash) GetFlashes(w http.ResponseWriter, r *http.Request) map[string][]string {
	session, err := f.cookieStore.Get(r, sessionName)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
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