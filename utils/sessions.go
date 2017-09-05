package utils

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	SPOTIFY_STATE_KEY = "spotify_auth_state"
)

type SessionManager struct {
	store *sessions.CookieStore
}

func NewSessionManager(authKey []byte) *SessionManager {
	store := sessions.NewCookieStore(authKey)
	return &SessionManager{store}
}

func (manager *SessionManager) Get(r *http.Request, key string) string {
	session, err := manager.store.Get(r, SPOTIFY_STATE_KEY)
	if err != nil {
		log.Printf("SessionManager Get: Error getting session from store: %s", err.Error())
		return ""
	}

	val := session.Values[key]
	if str, ok := val.(string); !ok {
		log.Printf("SessionManager Get: Key %s not found in session", key)
		return ""
	} else {
		return str
	}
}

func (manager *SessionManager) Set(r *http.Request, w http.ResponseWriter, key string, value string) (err error) {
	session, err := manager.store.Get(r, SPOTIFY_STATE_KEY)
	if err != nil {
		log.Printf("SessionManager Set: Error getting session from store: %s", err.Error())
		return err
	}

	session.Values[key] = value
	session.Save(r, w)
	return
}

func (manager *SessionManager) Delete(r *http.Request, w http.ResponseWriter, key string) (err error) {
	session, err := manager.store.Get(r, SPOTIFY_STATE_KEY)
	if err != nil {
		log.Printf("SessionManager Delete: Error getting session from store: %s", err.Error())
		return
	}

	delete(session.Values, key)
	session.Save(r, w)
	return
}
