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
	router.HandleFunc("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks",
		ctrl.getTracksHandler)
	router.HandleFunc("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks/{"+TRACK_ID_PARAM+"}",
		ctrl.convertHandler).Methods("POST")
}

func (ctrl *TrackController) getTracksHandler(w http.ResponseWriter, r *http.Request) {
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

	user := ctrl.currentUser.Get(r)
	if user == nil {
		log.Printf("convert: No current user found from context")
		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
		return
	}

	client := user.GetClient(models.SPOTIFY_SERVICE)
	if client == nil {
		log.Printf("convert: No spotify client found for user %s", user.GetState())
		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
		return
	}

	playlist, err := client.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tracks, err := client.GetTracks(playlist)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(tracks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ctrl *TrackController) convertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerId, present := vars[OWNER_ID_PARAM]
	if !present {
		// TODO: Actually handle these errors lmao
		utils.RenderErrorTemplate(w, "URL is missing an owner ID", http.StatusUnprocessableEntity)
		return
	}

	playlistId, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		// TODO: Actually handle these errors lmao
		utils.RenderErrorTemplate(w, "URL is missing a playlist ID", http.StatusUnprocessableEntity)
		return
	}

	trackId, present := vars[TRACK_ID_PARAM]
	if !present {
		// TODO: Actually handle these errors lmao
		utils.RenderErrorTemplate(w, "URL is missing a track ID", http.StatusUnprocessableEntity)
		return
	}

	user := ctrl.currentUser.Get(r)
	if user == nil {
		log.Printf("convert: No current user found from context")
		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
		return
	}

	spotify := user.GetClient(models.SPOTIFY_SERVICE)
	if spotify == nil {
		log.Printf("convert: No spotify client found for user %s", user.GetState())
		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
		return
	}

	youtube := user.GetClient(models.YOUTUBE_SERVICE)
	if youtube == nil {
		log.Printf("convert: No youtube client found for user %s", user.GetState())
		http.Error(w, "User not logged into YouTube", http.StatusUnauthorized)
		return
	}

	spotifyPlaylist, err := spotify.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		http.Error(w, "error during conversion, could not retrieve playlist info", http.StatusInternalServerError)
		return
	}

	track, err := spotify.GetTrackByID(trackId)
	if err != nil {
		log.Printf("convert: error getting track to convert")
		log.Print(err)
		http.Error(w, "error during conversion, could not retrieve track", http.StatusInternalServerError)
		return
	}

	newPlaylistId := r.FormValue(NEW_PLAYLIST_ID_QUERY_PARAM)
	newPlaylist, err := youtube.GetOwnPlaylistInfo(newPlaylistId)
	if err != nil {
		log.Print(err)
		return
	}

	if newPlaylist == nil {
		newPlaylist, err = youtube.CreatePlaylist(spotifyPlaylist.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, "convert: unable to create YouTube playlist", http.StatusInternalServerError)
			return
		}
	}

	found, err := youtube.InsertTrack(newPlaylist, track)
	if err != nil {
		log.Print(err)
		http.Error(w, "convert: unable to insert track into YouTube playlist", http.StatusInternalServerError)
		return
	} else if !found {
		http.Error(w, "convert: unable to find video matching search results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newPlaylist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Write(js)
}
