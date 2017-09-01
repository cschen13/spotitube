package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const PLAYLIST_ID_PARAM = "playlistId"

func convertSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if playlistId, present := vars[PLAYLIST_ID_PARAM]; present {
		log.Printf("Got playlist ID: %s", playlistId)
	}
	fmt.Fprint(w, "Converting playlist...")
}
