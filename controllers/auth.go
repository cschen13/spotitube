package controllers

import (
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
	router.HandleFunc("/login/{"+SERVICE_PARAM+"}", ctrl.initiateAuthHandler)
	router.HandleFunc("/callback/{"+SERVICE_PARAM+"}", ctrl.completeAuthHandler)
	router.HandleFunc("/logout", ctrl.logoutHandler)
}

func (ctrl *AuthController) initiateAuthHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)[SERVICE_PARAM]
	auth, present := ctrl.auths[service]
	if !present {
		log.Printf("Unrecognized service %s", service)
		utils.RenderErrorTemplate(w, "An error occurred while logging in.", http.StatusInternalServerError)
		return
	}

	state := utils.GenerateRandStr(128)
	if user := ctrl.currentUser.Get(r); user == nil {
		log.Printf("auth: No current user found from context")
		err := ctrl.sessionManager.Set(r, w, utils.USER_STATE_KEY, state)
		if err != nil {
			utils.RenderErrorTemplate(w, "An error occurred while logging in. Please clear your cookies and try again.", http.StatusInternalServerError)
			return
		}
	} else {
		state = user.GetState()
	}

	url := auth.BuildAuthURL(state)
	log.Printf("Redirecting user to %s", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func (ctrl *AuthController) completeAuthHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)[SERVICE_PARAM]
	auth, present := ctrl.auths[service]
	if !present {
		log.Printf("Unrecognized service %s", service)
		utils.RenderErrorTemplate(w, "An error occurred while logging in.", http.StatusInternalServerError)
		return
	}

	if user := ctrl.currentUser.Get(r); user != nil {
		if err := user.AddClient(r, auth); err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			log.Print("Couldn't add new client to user:")
			log.Print(err)
			return
		}
		log.Printf("New %s client added to user %s", auth.GetType(), user.GetState())
	} else if storedState := ctrl.sessionManager.Get(r, utils.USER_STATE_KEY); storedState != "" {
		user, err := models.NewUser(storedState, r, auth)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			log.Print("Couldn't create user:")
			log.Print(err)
			return
		}
		user.Add()
		log.Printf("New %s client added to NEW user %s", auth.GetType(), storedState)
	} else {
		log.Print("No cookie for user found")
		http.Redirect(w, r, "/login/"+service, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/playlists", http.StatusFound)
}

func (ctrl *AuthController) logoutHandler(w http.ResponseWriter, r *http.Request) {
	models.DeleteUser(ctrl.sessionManager.Get(r, utils.USER_STATE_KEY))
	err := ctrl.sessionManager.Delete(r, w, utils.USER_STATE_KEY)
	if err != nil {
		log.Printf("Error logging out:")
		log.Print(err)
		utils.RenderErrorTemplate(w, "An error occurred while logging out.", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
