package server

import (
	"html/template"
	"log"
	"net/http"
)

var FUNC_MAP = template.FuncMap{
	"add": func(i int, j int) int {
		return i + j
	},
}

type errorPage struct {
	Message string
	Code    int
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	t, err := template.New(tmpl+".tmpl").Funcs(FUNC_MAP).ParseFiles("templates/_head.tmpl", "templates/"+tmpl+".tmpl")
	if err != nil {
		log.Print(err)
		http.Error(w, "Error generating the HTML template.", http.StatusInternalServerError)
		return
	}
	t.Execute(w, &p)
}

func renderErrorTemplate(w http.ResponseWriter, message string, code int) {
	renderTemplate(w, "error", errorPage{message, code})
}
