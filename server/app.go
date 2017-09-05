package server

import (
	"github.com/cschen13/spotitube/controllers"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	host   string
	port   string
	router http.Handler
}

func NewServer(host string, port string, sessionSecret string) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, "index", nil)
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RenderErrorTemplate(w, "This page doesn't exist.", http.StatusNotFound)
	})

	sessionManager := utils.NewSessionManager([]byte(sessionSecret))
	spotifyAuth := models.NewSpotifyAuthenticator(host + port)

	authCtrl := controllers.NewAuthController(sessionManager, spotifyAuth)
	playlistCtrl := controllers.NewPlaylistController(sessionManager)
	convertCtrl := controllers.NewConvertController(sessionManager)
	authCtrl.Register(router)
	playlistCtrl.Register(router)
	convertCtrl.Register(router)

	// serve images, JS files, etc.
	router.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(".")))
	// http.Handle("/", router)
	return &Server{host: host, port: port, router: router}
}

func (server *Server) Start() {
	go models.HandleUsers()
	log.Printf("Spinning up the server at %s...", server.host)
	// http.ListenAndServe(utils.GetPort(), nil)
	err := http.ListenAndServe(server.port, server.router)
	log.Printf(err.Error())
}
