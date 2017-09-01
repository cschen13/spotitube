package server

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// Routes
func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
		t, _ := template.ParseFiles("templates/index.tmpl")
		t.Execute(w, nil)
	})

	http.HandleFunc("/login", initiateAuth)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/playlists", getPlaylists)
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func Start() {
	go handleUsers()
	log.Println("Spinning up the server...")
	http.ListenAndServe(getPort(), nil)
}
