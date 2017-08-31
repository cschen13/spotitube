package server

import (
	"log"
	"net/http"
)

// Routes
func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	http.HandleFunc("/login", initiateAuth)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/playlists", getPlaylists)
}

func Start() {
	go handleUsers()
	log.Println("Spinning up the server...")
	http.ListenAndServe(":8080", nil)
}
