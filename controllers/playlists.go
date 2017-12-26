package controllers

import (
	"encoding/json"
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
	currentUser    *utils.CurrentUserManager
}

func NewPlaylistController(sessionManager *utils.SessionManager, currentUser *utils.CurrentUserManager) *PlaylistController {
	return &PlaylistController{sessionManager: sessionManager, currentUser: currentUser}
}

func (ctrl *PlaylistController) Register(router *mux.Router) {
	router.HandleFunc("/playlists", ctrl.getPlaylistsHandler)
}

func (ctrl *PlaylistController) getPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	user := ctrl.currentUser.Get(r)
	if user == nil {
		log.Println("Cannot retrieve playlists; user not logged in")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	clientParam := r.URL.Query().Get(CLIENT_PARAM)
	if clientParam == "" {
		clientParam = models.SPOTIFY_SERVICE
	}

	client := user.GetClient(clientParam)
	if client == nil {
		log.Printf("No %s client found for user %s", clientParam, user.GetState())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistPage, err := client.GetPlaylists(r.URL.Query().Get(PAGE_PARAM))
	if err != nil {
		http.Error(w, "An error occurred while generating the playlists.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	js, err := json.Marshal(playlistPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
