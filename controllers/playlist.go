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

// PlaylistController holds all API endpoints related to high-level playlist information.
type PlaylistController struct {
	sessionManager *utils.SessionManager
	// currentUser    *utils.CurrentUserManager
	auths map[string]models.Authenticator
}

// NewPlaylistController instantiates a new PlaylistController.
func NewPlaylistController(sessionManager *utils.SessionManager, auths map[string]models.Authenticator) *PlaylistController {
	return &PlaylistController{sessionManager: sessionManager, auths: auths}
}

// Register takes all endpoints in the controller and registers them with a router.
func (ctrl *PlaylistController) Register(router *mux.Router) {
	router.Handle("/playlists", utils.Handler(ctrl.getPlaylists))
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}", utils.Handler(ctrl.getPlaylistInfo))
}

type playlistsClient interface {
	GetPlaylists() (models.Playlists, error)
}

func (ctrl *PlaylistController) getPlaylists(w http.ResponseWriter, r *http.Request) error {
	clientType := r.URL.Query().Get(CLIENT_PARAM)
	if clientType == "" {
		clientType = models.SPOTIFY_SERVICE
	}

	tok, err := ctrl.sessionManager.GetToken(r, clientType)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("getPlaylists: no %s client found", clientType),
		}
	}

	c, err := ctrl.auths[clientType].NewClient(tok)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("getPlaylists: no %s client found", clientType),
		}
	}

	client, ok := c.(playlistsClient)
	if !ok {
		return utils.StatusError{
			Code: http.StatusMethodNotAllowed,
			Err:  fmt.Errorf("getPlaylists: %s client does not satisfy interface", clientType),
		}
	}

	playlists, err := client.GetPlaylists()
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	js, err := json.Marshal(playlists)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

type playlistInfoClient interface {
	GetPlaylistInfo(string, string) (*models.Playlist, error)
}

func (ctrl *PlaylistController) getPlaylistInfo(w http.ResponseWriter, r *http.Request) error {
	clientType := r.URL.Query().Get(CLIENT_PARAM)
	if clientType == "" {
		clientType = models.SPOTIFY_SERVICE
	}

	tok, err := ctrl.sessionManager.GetToken(r, clientType)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("getPlaylistInfo: no %s client found", clientType),
		}
	}

	c, err := ctrl.auths[clientType].NewClient(tok)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("getPlaylistInfo: no %s client found", clientType),
		}
	}

	client, ok := c.(playlistInfoClient)
	if !ok {
		return utils.StatusError{
			Code: http.StatusMethodNotAllowed,
			Err:  fmt.Errorf("getPlaylistInfo: %s client cannot be used to get playlist info", clientType),
		}
	}

	vars := mux.Vars(r)
	ownerID, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("getPlaylistInfo: URL is missing an owner ID"),
		}
	}

	playlistID, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("getPlaylistInfo: URL is missing a playlist ID"),
		}
	}

	playlist, err := client.GetPlaylistInfo(ownerID, playlistID)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	js, err := json.Marshal(playlist)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}
