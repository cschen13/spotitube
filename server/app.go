package server

import (
	"github.com/cschen13/spotitube/controllers"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Routes
func init() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, "index", nil)
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RenderErrorTemplate(w, "This page doesn't exist.", http.StatusNotFound)
	})

	controllers.RegisterAuthController(router)
	controllers.RegisterPlaylistController(router)
	controllers.RegisterConvertController(router)

	// serve images, JS files, etc.
	router.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(".")))
	http.Handle("/", router)
}

func Start() {
	go models.HandleUsers()
	log.Println("Spinning up the server...")
	http.ListenAndServe(utils.GetPort(), nil)
}
