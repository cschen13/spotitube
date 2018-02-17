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
	PLAYLIST_ID_PARAM = "playlistId"
	OWNER_ID_PARAM    = "ownerId"
)

type TrackController struct {
	sessionManager *utils.SessionManager
	currentUser    *utils.CurrentUserManager
}

func NewTrackController(sessionManager *utils.SessionManager, currentUser *utils.CurrentUserManager) *TrackController {
	return &TrackController{sessionManager: sessionManager, currentUser: currentUser}
}

func (ctrl *TrackController) Register(router *mux.Router) {
	router.HandleFunc("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}/tracks", ctrl.getTracksHandler)
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

// func (ctrl *TrackController) convertSpotifyHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	ownerId, present := vars[OWNER_ID_PARAM]
// 	if !present {
// 		utils.RenderErrorTemplate(w, "URL is missing an owner ID", http.StatusUnprocessableEntity)
// 		return
// 	}

// 	playlistId, present := vars[PLAYLIST_ID_PARAM]
// 	if !present {
// 		utils.RenderErrorTemplate(w, "URL is missing a playlist ID", http.StatusUnprocessableEntity)
// 		return
// 	}

// 	user := ctrl.currentUser.Get(r)
// 	if user == nil {
// 		log.Printf("convert: No current user found from context")
// 		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
// 		return
// 	}

// 	spotify := user.GetClient(models.SPOTIFY_SERVICE)
// 	if spotify == nil {
// 		log.Printf("convert: No spotify client found for user %s", user.GetState())
// 		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
// 		return
// 	}

// 	youtube := user.GetClient(models.YOUTUBE_SERVICE)
// 	if youtube == nil {
// 		log.Printf("convert: No youtube client found for user %s", user.GetState())
// 		//TODO: RETURN PAGE PARAMETER
// 		http.Redirect(w, r, "/login/"+models.YOUTUBE_SERVICE, http.StatusFound)
// 		return
// 	}

// 	playlist, err := spotify.GetPlaylistInfo(ownerId, playlistId)
// 	if err != nil {
// 		log.Printf("convert: error getting playlist info")
// 		log.Print(err)
// 		utils.RenderErrorTemplate(w, "Error occurred while retrieving playlist.", http.StatusInternalServerError)
// 		return
// 	}

// 	tracks, err := ctrl.getAllTracks(spotify, playlist)
// 	if err != nil {
// 		log.Printf("convert: Error occurred while retrieving playlist tracks")
// 		log.Print(err)
// 		utils.RenderErrorTemplate(w, "Error occurred while retrieving playlist tracks.", http.StatusInternalServerError)
// 		return
// 	}
// 	page := &models.ConvertPage{playlist, tracks}
// 	utils.RenderTemplate(w, "convert", page)
// }

// // TODO: Make getting tracks its own endpoint; call from front end
// func (ctrl *TrackController) getAllTracks(spotify models.Client, playlist models.Playlist) ([]models.PlaylistTrack, error) {
// 	pageNum := 1
// 	tracks, lastPage, err := spotify.GetPlaylistTracks(playlist, strconv.Itoa(pageNum))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var nextPageTracks []models.PlaylistTrack
// 	for !lastPage {
// 		pageNum += 1
// 		nextPageTracks, lastPage, err = spotify.GetPlaylistTracks(playlist, strconv.Itoa(pageNum))
// 		if err != nil {
// 			return nil, err
// 		}
// 		tracks = append(tracks, nextPageTracks...)
// 	}

// 	log.Printf("Found %d tracks in playlist %s", len(tracks), playlist.GetName())

// 	return tracks, nil
// }
