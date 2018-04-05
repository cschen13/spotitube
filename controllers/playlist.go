package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	// "log"
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
	router.Handle("/playlists", utils.Handler(ctrl.getPlaylists))
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}", utils.Handler(ctrl.getPlaylistInfo))
}

func (ctrl *PlaylistController) getPlaylists(w http.ResponseWriter, r *http.Request) error {
	user := ctrl.currentUser.Get(r)
	if user == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New("getPlaylists: user not logged in")}
		// log.Println("Cannot retrieve playlists; user not logged in")
		// w.WriteHeader(http.StatusUnauthorized)
		// return
	}

	clientParam := r.URL.Query().Get(CLIENT_PARAM)
	if clientParam == "" {
		clientParam = models.SPOTIFY_SERVICE
	}

	client := user.GetClient(clientParam)
	if client == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("getPlaylists: no %s client found for user %s", clientParam, user.GetState()))}
		// log.Printf("No %s client found for user %s", clientParam, user.GetState())
		// w.WriteHeader(http.StatusUnauthorized)
		// return
	}

	playlists, err := client.GetPlaylists()
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
		// http.Error(w, "An error occurred while generating the playlists.", http.StatusInternalServerError)
		// log.Println(err)
		// return
	}

	js, err := json.Marshal(playlists)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// log.Print(err)
		// return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func (ctrl *PlaylistController) getPlaylistInfo(w http.ResponseWriter, r *http.Request) error {
	user := ctrl.currentUser.Get(r)
	if user == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New("getPlaylists: user not logged in")}
		// log.Println("Cannot retrieve playlists; user not logged in")
		// w.WriteHeader(http.StatusUnauthorized)
		// return
	}

	clientParam := r.URL.Query().Get(CLIENT_PARAM)
	if clientParam == "" {
		clientParam = models.SPOTIFY_SERVICE
	}

	client := user.GetClient(clientParam)
	if client == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("getPlaylists: no %s client found for user %s", clientParam, user.GetState()))}
		// log.Printf("No %s client found for user %s", clientParam, user.GetState())
		// w.WriteHeader(http.StatusUnauthorized)
		// return
	}

	vars := mux.Vars(r)
	ownerId, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("getPlaylists: URL is missing an owner ID")}
		// utils.RenderErrorTemplate(w, "URL is missing an owner ID", http.StatusUnprocessableEntity)
		// return
	}

	playlistId, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("getPlaylists: URL is missing a playlist ID")}
		// utils.RenderErrorTemplate(w, "URL is missing a playlist ID", http.StatusUnprocessableEntity)
		// return
	}

	playlist, err := client.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// log.Print(err)
		// return
	}

	js, err := json.Marshal(playlist)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// log.Print(err)
		// return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}
