package server

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}, funcMap template.FuncMap) {
	t, _ := template.New(tmpl + ".tmpl").Funcs(funcMap).ParseFiles("templates/" + tmpl + ".tmpl")
	t.Execute(w, &p)
}

// Routes
func init() {
	router := mux.NewRouter()
	router.HandleFunc("/", landingHandler)
	router.HandleFunc("/login", initiateAuthHandler)
	router.HandleFunc("/callback", completeAuthHandler)
	router.HandleFunc("/playlists/{"+PAGE_PARAM+"}", getPlaylistsHandler)
	router.HandleFunc("/convert-spotify/{"+PLAYLIST_ID_PARAM+"}", convertSpotifyHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(".")))
	http.Handle("/", router)
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil, nil)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "404", nil, nil)
}

func Start() {
	go handleUsers()
	log.Println("Spinning up the server...")
	http.ListenAndServe(getPort(), nil)
}
