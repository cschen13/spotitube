package controllers

import (
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	USER_STATE_KEY = "state"
	SERVICE_PARAM  = "service"
)

type AuthController struct {
	sessionManager *utils.SessionManager
	auths          *map[string]models.Authenticator
}

func NewAuthController(sessionManager *utils.SessionManager, auths *map[string]models.Authenticator) *AuthController {
	return &AuthController{sessionManager: sessionManager, auths: auths}
}

func (ctrl *AuthController) Register(router *mux.Router) {
	router.HandleFunc("/login/{"+SERVICE_PARAM+"}", ctrl.initiateAuthHandler)
	router.HandleFunc("/callback/{"+SERVICE_PARAM+"}", ctrl.completeAuthHandler)
}

func (ctrl *AuthController) initiateAuthHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)[SERVICE_PARAM]
	state := ctrl.sessionManager.Get(r, USER_STATE_KEY)
	user := models.GetUser(state)
	if state != "" && user != nil {
		if user.GetClient(service) != nil {
			log.Printf("User already logged in to %s, redirecting to playlists", service)
			http.Redirect(w, r, "/playlists", http.StatusFound)
			return
		}
	} else {
		state = utils.GenerateRandStr(128)
		err := ctrl.sessionManager.Set(r, w, USER_STATE_KEY, state)
		if err != nil {
			utils.RenderErrorTemplate(w, "An error occurred while logging in. Please clear your cookies and try again.", http.StatusInternalServerError)
			return
		}
	}

	if auth, present := (*ctrl.auths)[service]; !present {
		log.Printf("Unrecognized service %s", service)
		utils.RenderErrorTemplate(w, "An error occurred while logging in.", http.StatusInternalServerError)
	} else {
		url := auth.BuildAuthURL(state)
		log.Printf("Redirecting user to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (ctrl *AuthController) completeAuthHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)[SERVICE_PARAM]
	storedState := ctrl.sessionManager.Get(r, USER_STATE_KEY)
	if storedState == "" {
		log.Print("No cookie for user found")
		http.Redirect(w, r, "/login/"+service, http.StatusFound)
		return
	}

	if auth, present := (*ctrl.auths)[service]; !present {
		log.Printf("Unrecognized service %s", service)
		utils.RenderErrorTemplate(w, "An error occurred while logging in.", http.StatusInternalServerError)
	} else {
		user := models.GetUser(storedState)
		if user != nil {
			err := user.AddClient(storedState, r, auth)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				log.Print("Couldn't add new client to user:")
				log.Print(err)
				return
			}
			log.Printf("New %s client added to user %s", auth.GetType(), storedState)
		} else {
			user, err := models.NewUser(storedState, r, auth)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				log.Print("Couldn't create user:")
				log.Print(err)
				return
			}
			user.Add()
		}

		http.Redirect(w, r, "/playlists", http.StatusFound)
	}
}
