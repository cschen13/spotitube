package utils

import (
	// "github.com/cschen13/spotitube/models"
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	t, err := template.New(tmpl+".tmpl").ParseFiles("views/_head.tmpl", "views/"+tmpl+".tmpl")
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
