package controllers

import (
	"fmt"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const PLAYLIST_ID_PARAM = "playlistId"

type ConvertController struct {
	sessionManager *utils.SessionManager
}

func NewConvertController(sessionManager *utils.SessionManager) *ConvertController {
	return &ConvertController{sessionManager: sessionManager}
}

func (ctrl *ConvertController) Register(router *mux.Router) {
	router.HandleFunc("/convert-spotify/{"+PLAYLIST_ID_PARAM+"}", ctrl.convertSpotifyHandler)
}

func (ctrl *ConvertController) convertSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if playlistId, present := vars[PLAYLIST_ID_PARAM]; present {
		log.Printf("Got playlist ID: %s", playlistId)
	}
	fmt.Fprint(w, "Converting playlist...")
}
