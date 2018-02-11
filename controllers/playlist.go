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
	router.HandleFunc("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}", ctrl.getPlaylistInfoHandler)
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

	playlists, err := client.GetPlaylists()
	if err != nil {
		http.Error(w, "An error occurred while generating the playlists.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	js, err := json.Marshal(playlists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ctrl *PlaylistController) getPlaylistInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	ownerId, present := vars[OWNER_ID_PARAM]
	if !present {
		utils.RenderErrorTemplate(w, "URL is missing an owner ID", http.StatusUnprocessableEntity)
		return
	}

	playlistId, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		utils.RenderErrorTemplate(w, "URL is missing a playlist ID", http.StatusUnprocessableEntity)
		return
	}

	playlist, err := client.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	js, err := json.Marshal(playlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
