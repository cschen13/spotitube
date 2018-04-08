package controllers

import (
	"errors"
	"fmt"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	SERVICE_PARAM = "service"
)

type AuthController struct {
	sessionManager *utils.SessionManager
	auths          map[string]models.Authenticator
	currentUser    *utils.CurrentUserManager
}

func NewAuthController(sessionManager *utils.SessionManager, auths map[string]models.Authenticator, currentUser *utils.CurrentUserManager) *AuthController {
	return &AuthController{sessionManager: sessionManager, auths: auths, currentUser: currentUser}
}

func (ctrl *AuthController) Register(router *mux.Router) {
	router.Handle("/login/{"+SERVICE_PARAM+"}", utils.Handler(ctrl.initiateAuth))
	router.Handle("/callback/{"+SERVICE_PARAM+"}", utils.Handler(ctrl.completeAuth))
	router.Handle("/logout", utils.Handler(ctrl.logout))
}

func (ctrl *AuthController) initiateAuth(w http.ResponseWriter, r *http.Request) error {
	if returnURL := r.FormValue("returnURL"); returnURL != "" {
		ctrl.sessionManager.Set(r, w, "RedirectAfterLogin", returnURL)
	}

	service := mux.Vars(r)[SERVICE_PARAM]
	auth, present := ctrl.auths[service]
	if !present {
		return utils.PageError{
			http.StatusInternalServerError,
			errors.New(fmt.Sprintf("initiateAuth: unrecognized service %s", service)),
			"An error occurred while logging in. Please try again.",
		}
	}

	state := utils.GenerateRandStr(128)
	if user := ctrl.currentUser.Get(r); user == nil {
		err := ctrl.sessionManager.Set(r, w, utils.USER_STATE_KEY, state)
		if err != nil {
			return utils.PageError{
				http.StatusInternalServerError,
				err,
				"An error occurred while logging in. Please clear your cookies and try again.",
			}
		}

	} else {
		state = user.GetState()
	}

	url := auth.BuildAuthURL(state)
	log.Printf("Redirecting user to %s", url)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func (ctrl *AuthController) completeAuth(w http.ResponseWriter, r *http.Request) error {
	service := mux.Vars(r)[SERVICE_PARAM]
	auth, present := ctrl.auths[service]
	if !present {
		return utils.PageError{
			http.StatusInternalServerError,
			errors.New(fmt.Sprintf("initiateAuth: unrecognized service %s", service)),
			"An error occurred while logging in. Please try again.",
		}
	}

	clientType := auth.GetType()
	if user := ctrl.currentUser.Get(r); user != nil {
		if err := user.AddClient(r, auth); err != nil {
			return utils.PageError{
				http.StatusInternalServerError,
				err,
				"An error occurred while logging in. Please try again.",
			}
		}

		log.Printf("New %s client added to user %s", clientType, user.GetState())
	} else if storedState := ctrl.sessionManager.Get(r, utils.USER_STATE_KEY); storedState != "" {
		user, err := models.NewUser(storedState, r, auth)
		if err != nil {
			return utils.PageError{
				http.StatusInternalServerError,
				err,
				"An error occurred while logging in. Please try again.",
			}
		}

		user.Add()
		log.Printf("New %s client added to NEW user %s", clientType, storedState)
	} else {
		log.Print("No cookie for user found")
		http.Redirect(w, r, "/login/"+service, http.StatusFound)
		return nil
	}

	redirectTo := "/"
	if returnURL := ctrl.sessionManager.Get(r, "RedirectAfterLogin"); returnURL != "" {
		ctrl.sessionManager.Delete(r, w, "RedirectAfterLogin")
		redirectTo = returnURL
	}

	http.Redirect(w, r, redirectTo, http.StatusFound)
	return nil
}

func (ctrl *AuthController) logout(w http.ResponseWriter, r *http.Request) error {
	models.DeleteUser(ctrl.sessionManager.Get(r, utils.USER_STATE_KEY))
	err := ctrl.sessionManager.Delete(r, w, utils.USER_STATE_KEY)
	if err != nil {
		return utils.PageError{
			http.StatusInternalServerError,
			err,
			"An error occurred while logging out.",
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
