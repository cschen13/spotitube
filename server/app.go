package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// Routes
func init() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", nil)
	})

	router.HandleFunc("/login", initiateAuthHandler)
	router.HandleFunc("/callback", completeAuthHandler)
	router.HandleFunc("/playlists", getPlaylistsHandler)
	router.HandleFunc("/convert-spotify/{"+PLAYLIST_ID_PARAM+"}", convertSpotifyHandler)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderErrorTemplate(w, "This page doesn't exist.", http.StatusNotFound)
	})

	// serve images, JS files, etc.
	router.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(".")))
	http.Handle("/", router)
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}

func Start() {
	go handleUsers()
	log.Println("Spinning up the server...")
	http.ListenAndServe(getPort(), nil)
}
