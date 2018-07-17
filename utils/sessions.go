package utils

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SESSION_KEY = "spotitube_session"
)

type SessionManager struct {
	store *sessions.CookieStore
}

func NewSessionManager(authKey []byte) *SessionManager {
	store := sessions.NewCookieStore(authKey)
	return &SessionManager{store}
}

func (manager *SessionManager) GetClient(r *http.Request, clientType string) interface{} {
	return manager.get(r, clientType)
}

func (manager *SessionManager) GetState(r *http.Request) string {
	return manager.getString(r, USER_STATE_KEY)
}

func (manager *SessionManager) GetRedirect(r *http.Request) string {
	return manager.getString(r, "RedirectAfterLogin")
}

func (manager *SessionManager) get(r *http.Request, key string) interface{} {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Get: Error getting session from store: %s", err.Error())
		return ""
	}

	log.Printf("GETTING %v", session.Values)
	return session.Values[key]
}

func (manager *SessionManager) getString(r *http.Request, key string) string {
	val := manager.get(r, key)
	if str, ok := val.(string); !ok {
		log.Printf("SessionManager getString: Key %s not found in session", key)
		return ""
	} else {
		return str
	}
}

func (manager *SessionManager) SetState(r *http.Request, w http.ResponseWriter, state string) error {
	return manager.set(r, w, USER_STATE_KEY, state)
}

func (manager *SessionManager) SetClient(r *http.Request, w http.ResponseWriter, clientType string, client interface{}) error {
	log.Printf("Setting client %s", clientType)
	return manager.set(r, w, clientType, client)
}

func (manager *SessionManager) SetRedirect(r *http.Request, w http.ResponseWriter, path string) error {
	return manager.set(r, w, "RedirectAfterLogin", path)
}

func (manager *SessionManager) set(r *http.Request, w http.ResponseWriter, key string, value interface{}) (err error) {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Set: Error getting session from store: %s", err.Error())
		return err
	}

	session.Values[key] = value
	log.Printf("SETTING %v", session.Values["spotify"])
	err = session.Save(r, w)
	return
}

func (manager *SessionManager) DeleteRedirect(r *http.Request, w http.ResponseWriter) error {
	return manager.delete(r, w, "RedirectAfterLogin")
}

func (manager *SessionManager) delete(r *http.Request, w http.ResponseWriter, key string) (err error) {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Delete: Error getting session from store: %s", err.Error())
		return
	}

	delete(session.Values, key)
	err = session.Save(r, w)
	return
}
