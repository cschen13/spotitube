package controllers

import (
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const USER_STATE_KEY = "state"

type AuthController struct {
	sessionManager *utils.SessionManager
}

func NewAuthController(sessionManager *utils.SessionManager) *AuthController {
	return &AuthController{sessionManager: sessionManager}
}

func (ctrl *AuthController) Register(router *mux.Router) {
	router.HandleFunc("/login", ctrl.initiateAuthHandler)
	router.HandleFunc("/callback", ctrl.completeAuthHandler)
}

func (ctrl *AuthController) initiateAuthHandler(w http.ResponseWriter, r *http.Request) {
	if state := ctrl.sessionManager.Get(r, USER_STATE_KEY); state == "" {
		state = utils.GenerateRandStr(128)
		err := ctrl.sessionManager.Set(r, w, USER_STATE_KEY, state)
		if err != nil {
			utils.RenderErrorTemplate(w, "An error occurred while logging in. Please clear your cookies and try again.", http.StatusInternalServerError)
			return
		}

		url := models.BuildSpotifyAuthURL(state)
		log.Printf("Redirecting user to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		log.Printf("Found existing session, redirecting to playlists")
		http.Redirect(w, r, "playlists", http.StatusFound)
	}
}

func (ctrl *AuthController) completeAuthHandler(w http.ResponseWriter, r *http.Request) {
	storedState := ctrl.sessionManager.Get(r, USER_STATE_KEY)
	if storedState == "" {
		log.Print("No cookie for spotify auth state found")
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}

	// acquire access token (also checks state parameter)
	user, err := models.NewUser(storedState, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Print("Couldn't create user:")
		log.Print(err)
		return
	}

	user.Add()
	http.Redirect(w, r, "playlists", http.StatusFound)
}
