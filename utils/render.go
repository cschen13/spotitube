package utils

import (
	"github.com/cschen13/spotitube/models"
	"html/template"
	"log"
	"net/http"
)

var FUNC_MAP = template.FuncMap{
	"id": func(playlist models.Playlist) string {
		return playlist.GetID()
	},
	"name": func(playlist models.Playlist) string {
		return playlist.GetName()
	},
	"url": func(playlist models.Playlist) string {
		return playlist.GetURL()
	},
	"coverUrl": func(playlist models.Playlist) string {
		return playlist.GetCoverURL()
	},
	"title": func(track models.PlaylistTrack) string {
		return track.GetTitle()
	},
	"artist": func(track models.PlaylistTrack) string {
		return track.GetArtist()
	},
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	t, err := template.New(tmpl+".tmpl").Funcs(FUNC_MAP).ParseFiles("views/_head.tmpl", "views/"+tmpl+".tmpl")
	if err != nil {
		log.Print(err)
		http.Error(w, "Error generating the HTML template.", http.StatusInternalServerError)
		return
	}
	t.Execute(w, &p)
}

type errorPage struct {
	Message string
	Code    int
}

func RenderErrorTemplate(w http.ResponseWriter, message string, code int) {
	RenderTemplate(w, "error", errorPage{message, code})
}
