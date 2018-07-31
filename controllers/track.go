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

type TrackController struct {
	sessionManager *utils.SessionManager
	// currentUser    *utils.CurrentUserManager
	auths map[string]models.Authenticator
}

func NewTrackController(sessionManager *utils.SessionManager, auths map[string]models.Authenticator) *TrackController {
	return &TrackController{sessionManager: sessionManager, auths: auths}
}

func (ctrl *TrackController) Register(router *mux.Router) {
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks",
		utils.Handler(ctrl.getTracks))
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks/{"+TRACK_ID_PARAM+"}",
		utils.Handler(ctrl.convert)).Methods("POST")
}

type tracksClient interface {
	GetPlaylistInfo(string, string) (*models.Playlist, error)
	GetTracks(*models.Playlist) (models.Tracks, error)
}

func (ctrl *TrackController) getTracks(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ownerID, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("getTracks: URL is missing an owner ID"),
		}
	}

	playlistID, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("getTracks: URL is missing a playlist ID"),
		}
	}

	clientType := r.URL.Query().Get(CLIENT_PARAM)

	// TODO: getTracks functionality for other services besides Spotify
	// if clientType == "" {
	clientType = models.SPOTIFY_SERVICE
	// }

	tok, err := ctrl.sessionManager.GetToken(r, clientType)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("getTracks: no %s client found", clientType),
		}
	}

	c, err := ctrl.auths[clientType].NewClient(tok)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("getTracks: no %s client found", clientType),
		}
	}

	client, ok := c.(tracksClient)
	if !ok {
		return utils.StatusError{
			Code: http.StatusMethodNotAllowed,
			Err:  fmt.Errorf("getTracks: %s client does not satisfy interface", clientType),
		}
	}

	playlist, err := client.GetPlaylistInfo(ownerID, playlistID)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	tracks, err := client.GetTracks(playlist)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	js, err := json.Marshal(tracks)
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

type convertSrcClient interface {
	GetPlaylistInfo(string, string) (*models.Playlist, error)
	GetTrackByID(string) (*models.Track, error)
}

type convertDstClient interface {
	GetOwnPlaylistInfo(string) (*models.Playlist, error)
	CreatePlaylist(string) (*models.Playlist, error)
	InsertTrack(*models.Playlist, *models.Track) (bool, error)
}

func (ctrl *TrackController) convert(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ownerID, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("convert: URL is missing an owner ID"),
		}
	}

	playlistID, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("convert: URL is missing a playlist ID"),
		}
	}

	trackId, present := vars[TRACK_ID_PARAM]
	if !present {
		return utils.StatusError{
			Code: http.StatusUnprocessableEntity,
			Err:  errors.New("convert: URL is missing a track ID"),
		}
	}

	clientType := r.URL.Query().Get(CLIENT_PARAM)

	// TODO: getTracks functionality for other services besides Spotify
	// if clientType == "" {
	clientType = models.SPOTIFY_SERVICE
	// }

	tok, err := ctrl.sessionManager.GetToken(r, clientType)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("convert: no %s client found", clientType),
		}
	}

	c, err := ctrl.auths[clientType].NewClient(tok)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("convert: no %s client found", clientType),
		}
	}

	spotify, ok := c.(convertSrcClient)
	if !ok {
		return utils.StatusError{
			Code: http.StatusMethodNotAllowed,
			Err:  fmt.Errorf("convert: %s client does not satisfy source interface", clientType),
		}
	}

	clientType = models.YOUTUBE_SERVICE

	tok, err = ctrl.sessionManager.GetToken(r, clientType)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("convert: no %s client found", clientType),
		}
	}

	c, err = ctrl.auths[clientType].NewClient(tok)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("convert: no %s client found", clientType),
		}
	}

	youtube, ok := c.(convertDstClient)
	if !ok {
		return utils.StatusError{
			Code: http.StatusMethodNotAllowed,
			Err:  fmt.Errorf("convert: %s client does not satisfy destination interface"),
		}
	}

	spotifyPlaylist, err := spotify.GetPlaylistInfo(ownerID, playlistID)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	track, err := spotify.GetTrackByID(trackId)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	newPlaylistID := r.FormValue(NEW_PLAYLIST_ID_QUERY_PARAM)
	newPlaylist, err := youtube.GetOwnPlaylistInfo(newPlaylistID)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if newPlaylist == nil {
		newPlaylist, err = youtube.CreatePlaylist(spotifyPlaylist.Name)
		if err != nil {
			return utils.StatusError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
	}

	found, err := youtube.InsertTrack(newPlaylist, track)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	} else if !found {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("convert: unable to find video matching search results"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newPlaylist)
	if err != nil {
		return utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	w.Write(js)
	return nil
}
