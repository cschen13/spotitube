package utils

import (
	"log"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const (
	SESSION_KEY            = "spotitube_session"
	SESSION_DIRECTORY_NAME = ".spotitube"
)

type SessionManager struct {
	store *sessions.FilesystemStore
}

func NewSessionManager(authKey []byte) *SessionManager {
	usr, err := user.Current()
	if err != nil {
		log.Printf("SessionManager: Error getting current user during session manager instantiation")
		return nil
	}
	store := sessions.NewFilesystemStore(filepath.Join(usr.HomeDir, SESSION_DIRECTORY_NAME), authKey)
	return &SessionManager{store}
}

func (manager *SessionManager) Get(r *http.Request, key string) string {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Get: Error getting session from store: %s", err.Error())
		return ""
	}

	val := session.Values[key]
	log.Printf("%v", session.Values)
	if str, ok := val.(string); !ok {
		log.Printf("SessionManager Get: Key %s not found in session", key)
		return ""
	} else {
		return str
	}
}

func (manager *SessionManager) GetToken(r *http.Request, key string) *oauth2.Token {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Get: Error getting session from store: %s", err.Error())
		return nil
	}

	val := session.Values[key]
	log.Printf("%v", session.Values)
	if tok, ok := val.(*oauth2.Token); !ok {
		log.Printf("SessionManager Get: Key %s not found in session", key)
		return nil
	} else {
		return tok
	}
}

func (manager *SessionManager) Set(r *http.Request, w http.ResponseWriter, key string, value interface{}) (err error) {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Set: Error getting session from store: %s", err.Error())
		return err
	}

	session.Values[key] = value
	log.Printf("Set: %v", session.Values)
	err = session.Save(r, w)
	return
}

func (manager *SessionManager) Delete(r *http.Request, w http.ResponseWriter, key string) (err error) {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Delete: Error getting session from store: %s", err.Error())
		return
	}

	delete(session.Values, key)
	err = session.Save(r, w)
	return
}
