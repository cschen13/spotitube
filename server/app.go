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
		renderTemplate(w, "index", nil)
	})

	http.HandleFunc("/login", initiateAuth)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/playlists", getPlaylists)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	t, _ := template.ParseFiles("templates/" + tmpl + ".tmpl")
	t.Execute(w, &p)
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
