package controllers

import (
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	PAGE_PARAM   = "page"
	CLIENT_PARAM = "client"
)

type PlaylistController struct {
	sessionManager *utils.SessionManager
}

func NewPlaylistController(sessionManager *utils.SessionManager) *PlaylistController {
	return &PlaylistController{sessionManager: sessionManager}
}

func (ctrl *PlaylistController) Register(router *mux.Router) {
	router.HandleFunc("/playlists", ctrl.getPlaylistsHandler)
}

func (ctrl *PlaylistController) getPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	// check the request for a state cookie
	state := ctrl.sessionManager.Get(r, USER_STATE_KEY)
	if state == "" {
		log.Print("No cookie for user found")
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}

	user := models.GetUser(state)
	if user == nil {
		log.Print("No associated user found for state %s", state)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	clientParam := r.URL.Query().Get(CLIENT_PARAM)
	if clientParam == "" {
		clientParam = models.SPOTIFY_SERVICE
	}

	client := user.GetClient(clientParam)
	if client == nil {
		log.Printf("No %s client found for user %s", clientParam, state)
		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
		return
	}

	playlistPage, err := client.GetPlaylists(r.URL.Query().Get(PAGE_PARAM))
	if err != nil {
		utils.RenderErrorTemplate(w, "An error occurred while generating the playlists.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	utils.RenderTemplate(w, "playlists", playlistPage)
}
