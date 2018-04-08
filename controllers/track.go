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
	PLAYLIST_ID_PARAM           = "playlistId"
	OWNER_ID_PARAM              = "ownerId"
	TRACK_ID_PARAM              = "trackId"
	NEW_PLAYLIST_ID_QUERY_PARAM = "newPlaylistId"
)

type TrackController struct {
	sessionManager *utils.SessionManager
	currentUser    *utils.CurrentUserManager
}

func NewTrackController(sessionManager *utils.SessionManager, currentUser *utils.CurrentUserManager) *TrackController {
	return &TrackController{sessionManager: sessionManager, currentUser: currentUser}
}

func (ctrl *TrackController) Register(router *mux.Router) {
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks",
		utils.Handler(ctrl.getTracks))
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks/{"+TRACK_ID_PARAM+"}",
		utils.Handler(ctrl.convert)).Methods("POST")
}

func (ctrl *TrackController) getTracks(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ownerId, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("getTracks: URL is missing an owner ID")}
	}

	playlistId, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("getTracks: URL is missing a playlist ID")}
	}

	user := ctrl.currentUser.Get(r)
	if user == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New("getTracks: user not logged in")}
	}

	client := user.GetClient(models.SPOTIFY_SERVICE)
	if client == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("getTracks: no %s client found for user %s", models.SPOTIFY_SERVICE, user.GetState()))}
	}

	playlist, err := client.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	tracks, err := client.GetTracks(playlist)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	js, err := json.Marshal(tracks)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func (ctrl *TrackController) convert(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ownerId, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("convert: URL is missing an owner ID")}
	}

	playlistId, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("convert: URL is missing a playlist ID")}
	}

	trackId, present := vars[TRACK_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("convert: URL is missing a track ID")}
	}

	user := ctrl.currentUser.Get(r)
	if user == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New("getTracks: user not logged in")}
	}

	spotify := user.GetClient(models.SPOTIFY_SERVICE)
	if spotify == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("convert: no %s client found for user %s", models.SPOTIFY_SERVICE, user.GetState()))}
	}

	youtube := user.GetClient(models.YOUTUBE_SERVICE)
	if youtube == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("convert: no %s client found for user %s", models.YOUTUBE_SERVICE, user.GetState()))}
	}

	spotifyPlaylist, err := spotify.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	track, err := spotify.GetTrackByID(trackId)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	newPlaylistId := r.FormValue(NEW_PLAYLIST_ID_QUERY_PARAM)
	newPlaylist, err := youtube.GetOwnPlaylistInfo(newPlaylistId)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	if newPlaylist == nil {
		newPlaylist, err = youtube.CreatePlaylist(spotifyPlaylist.Name)
		if err != nil {
			return utils.StatusError{http.StatusInternalServerError, err}
		}
	}

	found, err := youtube.InsertTrack(newPlaylist, track)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	} else if !found {
		return utils.StatusError{http.StatusInternalServerError, errors.New("convert: unable to find video matching search results")}
	}

	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newPlaylist)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	w.Write(js)
	return nil
}
