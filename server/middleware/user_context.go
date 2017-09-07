package middleware

import (
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

type UserContext struct {
	currentUser    *utils.CurrentUserManager
	sessionManager *utils.SessionManager
}

func NewUserContext(currentUser *utils.CurrentUserManager, sessionManager *utils.SessionManager) *UserContext {
	return &UserContext{currentUser, sessionManager}
}

func (uc *UserContext) Middleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		state := uc.sessionManager.Get(r, utils.USER_STATE_KEY)
		if user := models.GetUser(state); state != "" && user != nil {
			log.Printf("Current user with state %s", state)
			uc.currentUser.Set(r, user)
		} else {
			log.Printf("No current user found")
		}
		next(w, r)
	}
}
