package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
)

const (
	SERVICE_PARAM = "service"
)

type AuthController struct {
	sessionManager *utils.SessionManager
	auths          map[string]models.Authenticator
	// currentUser    *utils.CurrentUserManager
}

func NewAuthController(sessionManager *utils.SessionManager, auths map[string]models.Authenticator) *AuthController {
	return &AuthController{sessionManager: sessionManager, auths: auths}
}

func (ctrl *AuthController) Register(router *mux.Router) {
	router.Handle("/login/{"+SERVICE_PARAM+"}", utils.Handler(ctrl.initiateAuth))
	router.Handle("/callback/{"+SERVICE_PARAM+"}", utils.Handler(ctrl.completeAuth))
	// router.Handle("/logout", utils.Handler(ctrl.logout))
}

func (ctrl *AuthController) initiateAuth(w http.ResponseWriter, r *http.Request) error {
	if returnURL := r.FormValue("returnURL"); returnURL != "" {
		ctrl.sessionManager.SetRedirect(r, w, returnURL)
	}

	service := mux.Vars(r)[SERVICE_PARAM]
	auth, present := ctrl.auths[service]
	if !present {
		return utils.PageError{
			http.StatusInternalServerError,
			fmt.Errorf("initiateAuth: unrecognized service %s", service),
			"An error occurred while logging in. Please try again.",
		}
	}

	state := utils.GenerateRandStr(128)
	err := ctrl.sessionManager.SetState(r, w, state)
	if err != nil {
		return utils.PageError{
			http.StatusInternalServerError,
			err,
			"An error occurred while logging in. Please clear your cookies and try again.",
		}
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
			fmt.Errorf("completeAuth: unrecognized service %s", service),
			"An error occurred while logging in. Please try again.",
		}
	}

	state, err := ctrl.sessionManager.GetState(r)
	if err != nil {
		return utils.PageError{
			http.StatusInternalServerError,
			err,
			"An error occurred while logging in. Please try again.",
		}
	}

	tok, err := auth.GenerateToken(state, r)
	if err != nil {
		return utils.PageError{
			http.StatusInternalServerError,
			err,
			"An error occurred while logging in. Please try again.",
		}
	}

	err = ctrl.sessionManager.SetToken(r, w, auth.GetType(), tok)
	if err != nil {
		return utils.PageError{
			http.StatusInternalServerError,
			err,
			"An error occurred while logging in. Please try again.",
		}
	}

	redirectTo := "/"
	if returnURL, err := ctrl.sessionManager.GetRedirect(r); returnURL != "" && err == nil {
		ctrl.sessionManager.DeleteRedirect(r, w)
		redirectTo = returnURL
	}

	http.Redirect(w, r, redirectTo, http.StatusFound)
	return nil
}

// func (ctrl *AuthController) logout(w http.ResponseWriter, r *http.Request) error {
// 	models.DeleteUser(ctrl.sessionManager.Get(r, utils.USER_STATE_KEY))
// 	err := ctrl.sessionManager.Delete(r, w, utils.USER_STATE_KEY)
// 	if err != nil {
// 		return utils.PageError{
// 			http.StatusInternalServerError,
// 			err,
// 			"An error occurred while logging out.",
// 		}
// 	}

// 	http.Redirect(w, r, "/", http.StatusFound)
// 	return nil
// }
