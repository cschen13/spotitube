package server

import (
	"github.com/cschen13/spotitube/controllers"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/server/middleware"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
)

type Server struct {
	*negroni.Negroni
	host string
	port string
}

func NewServer(host string, port string, sessionSecret string, userManagerKey int, isDev bool) *Server {
	server := Server{negroni.Classic(), host, port}
	sessionManager := utils.NewSessionManager([]byte(sessionSecret))
	currentUser := utils.NewCurrentUserManager(userManagerKey)
	userContext := middleware.NewUserContext(currentUser, sessionManager)
	server.Use(userContext.Middleware())

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, "index", nil)
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RenderErrorTemplate(w, "This page doesn't exist.", http.StatusNotFound)
	})

	var spotifyAuth *models.SpotifyAuthenticator
	if isDev {
		spotifyAuth = models.NewSpotifyAuthenticator(host + port)
	} else {
		spotifyAuth = models.NewSpotifyAuthenticator(host)
	}

	json := os.Getenv("YOUTUBE_SECRET")
	if json == "" {
		log.Fatalf("Client secret for youtube not found")
	}

	youtubeAuth := models.NewYoutubeAuthenticator(json, host+port, isDev)

	auths := make(map[string]models.Authenticator)
	auths[spotifyAuth.GetType()] = spotifyAuth
	auths[youtubeAuth.GetType()] = youtubeAuth

	authCtrl := controllers.NewAuthController(sessionManager, auths, currentUser)
	playlistCtrl := controllers.NewPlaylistController(sessionManager, currentUser)
	convertCtrl := controllers.NewConvertController(sessionManager, currentUser)
	authCtrl.Register(router)
	playlistCtrl.Register(router)
	convertCtrl.Register(router)

	// serve images, JS files, etc.
	router.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(".")))

	server.UseHandler(router)
	return &server
}

func (server *Server) Start() {
	go models.HandleUsers()
	log.Printf("Spinning up the server at %s%s...", server.host, server.port)
	err := http.ListenAndServe(server.port, server)
	log.Printf(err.Error())
}
