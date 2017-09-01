package server

import (
	"log"
	"net/http"
)

type playlistPage struct {
}

func getPlaylists(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(STATE_KEY)
	if err != nil {
		http.Error(w, "Access denied. Please try logging in.", http.StatusForbidden)
		log.Print("No stored state found")
		return
	}

	state := cookie.Value
	client := getSpotify(state)
	if client == nil {
		http.Error(w, "Access denied. Please try logging in.", http.StatusForbidden)
		log.Print("No associated user found for state %s", state)
		return
	}

	page, err := client.CurrentUsersPlaylists()
	if err != nil {
		http.Error(w, "Error while getting playlists", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// for _, playlist := range page.Playlists {
	// 	log.Printf("%v", playlist.Images)
	// }

	renderTemplate(w, "playlists", page)
}
