package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const PLAYLIST_ID_PARAM = "playlistId"

func RegisterConvertController(router *mux.Router) {
	router.HandleFunc("/convert-spotify/{"+PLAYLIST_ID_PARAM+"}", convertSpotifyHandler)
}

func convertSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if playlistId, present := vars[PLAYLIST_ID_PARAM]; present {
		log.Printf("Got playlist ID: %s", playlistId)
	}
	fmt.Fprint(w, "Converting playlist...")
}
