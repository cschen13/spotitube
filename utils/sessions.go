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
	REDIRECT_KEY           = "RedirectAfterLogin"
	STATE_KEY              = "state"
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

func (manager *SessionManager) GetToken(r *http.Request, clientType string) (*oauth2.Token, error) {
	val, err := manager.get(r, clientType)
	if err != nil {
		return nil, err
	}

	if tok, ok := val.(*oauth2.Token); !ok {
		log.Printf("SessionManager Get: Token %s not found in session", clientType)
		return nil, nil
	} else {
		return tok, nil
	}
}

func (manager *SessionManager) GetState(r *http.Request) (string, error) {
	return manager.getString(r, STATE_KEY)
}

// GetRedirect retrieves the URL that a user last visited before logging in
// so that they can be redirected to a page other than home.
func (manager *SessionManager) GetRedirect(r *http.Request) (path string, err error) {
	path, err = manager.getString(r, REDIRECT_KEY)
	if err != nil {
		log.Printf("SessionManager GetRedirect: Error getting redirect URL from session: %s", err.Error())
	}

	return
}

func (manager *SessionManager) getString(r *http.Request, key string) (string, error) {
	val, err := manager.get(r, key)
	if err != nil {
		return "", err
	}

	if str, ok := val.(string); !ok {
		log.Printf("SessionManager Get: Key %s not found in session", key)
		return "", nil
	} else {
		return str, nil
	}
}

func (manager *SessionManager) get(r *http.Request, key string) (interface{}, error) {
	session, err := manager.store.Get(r, SESSION_KEY)
	if err != nil {
		log.Printf("SessionManager Get: Error getting session from store: %s", err.Error())
		return nil, err
	}

	return session.Values[key], nil
}

func (manager *SessionManager) SetToken(r *http.Request, w http.ResponseWriter, clientType string, tok *oauth2.Token) error {
	return manager.set(r, w, clientType, tok)
}

// SetState associates a state string with a particular session.
// It should be used as part of the authentication flow to verify login callback requests.
func (manager *SessionManager) SetState(r *http.Request, w http.ResponseWriter, state string) error {
	return manager.set(r, w, STATE_KEY, state)
}

func (manager *SessionManager) SetRedirect(r *http.Request, w http.ResponseWriter, path string) error {
	return manager.set(r, w, REDIRECT_KEY, path)
}

func (manager *SessionManager) set(r *http.Request, w http.ResponseWriter, key string, value interface{}) (err error) {
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

func (manager *SessionManager) DeleteRedirect(r *http.Request, w http.ResponseWriter) error {
	return manager.delete(r, w, REDIRECT_KEY)
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
