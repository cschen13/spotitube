package controllers

import (
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const PLAYLIST_ID_PARAM = "playlistId"

type ConvertController struct {
	sessionManager *utils.SessionManager
	currentUser    *utils.CurrentUserManager
}

func NewConvertController(sessionManager *utils.SessionManager, currentUser *utils.CurrentUserManager) *ConvertController {
	return &ConvertController{sessionManager: sessionManager, currentUser: currentUser}
}

func (ctrl *ConvertController) Register(router *mux.Router) {
	router.HandleFunc("/convert-"+models.SPOTIFY_SERVICE+"/{"+PLAYLIST_ID_PARAM+"}", ctrl.convertSpotifyHandler)
}

func (ctrl *ConvertController) convertSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
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

	spotify := user.GetClient(models.SPOTIFY_SERVICE)
	if spotify == nil {
		log.Printf("convert: No spotify client found for user %s", user.GetState())
		http.Redirect(w, r, "/login/"+models.SPOTIFY_SERVICE, http.StatusFound)
		return
	}

	youtube := user.GetClient(models.YOUTUBE_SERVICE)
	if youtube == nil {
		log.Printf("convert: No youtube client found for user %s", user.GetState())
		//TODO: RETURN PAGE PARAMETER
		http.Redirect(w, r, "/login/"+models.YOUTUBE_SERVICE, http.StatusFound)
		return
	}

	playlist, err := spotify.GetPlaylistInfo(playlistId)
	if err != nil {
		log.Printf("convert: error getting playlist info")
		log.Print(err)
		utils.RenderErrorTemplate(w, "Error occurred while retrieving playlist.", http.StatusInternalServerError)
		return
	}

	tracks, err := ctrl.getAllTracks(spotify, playlist)
	if err != nil {
		log.Printf("convert: Error occurred while retrieving playlist tracks")
		log.Print(err)
		utils.RenderErrorTemplate(w, "Error occurred while retrieving playlist tracks.", http.StatusInternalServerError)
		return
	}
	page := &models.ConvertPage{playlist, tracks}
	utils.RenderTemplate(w, "convert", page)
}

func (ctrl *ConvertController) getAllTracks(spotify models.Client, playlist models.Playlist) ([]models.PlaylistTrack, error) {
	pageNum := 1
	tracks, lastPage, err := spotify.GetPlaylistTracks(playlist, strconv.Itoa(pageNum))
	if err != nil {
		return nil, err
	}

	var nextPageTracks []models.PlaylistTrack
	for !lastPage {
		pageNum += 1
		nextPageTracks, lastPage, err = spotify.GetPlaylistTracks(playlist, strconv.Itoa(pageNum))
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, nextPageTracks...)
	}

	return tracks, nil
}
