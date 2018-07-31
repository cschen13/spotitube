package server

import (
	"log"
	"net/http"
	"os"

	"github.com/cschen13/spotitube/controllers"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Server struct {
	*negroni.Negroni
	host string
	port string
}

func NewServer(host string, port string, sessionSecret string, userManagerKey int, devPort string) *Server {
	isDev := port != devPort
	server := Server{negroni.Classic(), host, port}
	sessionManager := utils.NewSessionManager([]byte(sessionSecret))
	// currentUser := utils.NewCurrentUserManager(userManagerKey)
	// userContext := middleware.NewUserContext(currentUser, sessionManager)
	// server.Use(userContext.Middleware())

	router := mux.NewRouter()
	var spotifyAuth *models.SpotifyAuthenticator
	if isDev {
		spotifyAuth = models.NewSpotifyAuthenticator(host + port)
	} else {
		spotifyAuth = models.NewSpotifyAuthenticator(host)
	}

	configJson := os.Getenv("YOUTUBE_SECRET")
	if configJson == "" {
		log.Fatalf("Client secret for youtube not found")
	}

	youtubeAuth := models.NewYoutubeAuthenticator(configJson, host+port, isDev)

	auths := make(map[string]models.Authenticator)
	auths[spotifyAuth.GetType()] = spotifyAuth
	auths[youtubeAuth.GetType()] = youtubeAuth

	authCtrl := controllers.NewAuthController(sessionManager, auths)
	playlistCtrl := controllers.NewPlaylistController(sessionManager, auths)
	trackCtrl := controllers.NewTrackController(sessionManager, auths)
	authCtrl.Register(router)
	playlistCtrl.Register(router)
	trackCtrl.Register(router)

	if isDev {
		log.Printf("DEVELOPMENT: Dev server port %s", devPort)
		router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if acceptsJson := r.Header.Get("Accept") == "application/json"; acceptsJson {
				http.Error(w, "The requested resource does not exist.", http.StatusNotFound)
			} else {
				http.Redirect(w, r, host+devPort+r.URL.Path, http.StatusFound)
			}
		})
		// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	http.Redirect(w, r, host+devPort, http.StatusFound)
		// })
	} else {
		log.Print("PRODUCTION BUILD")
		// serve images, JS files, etc.
		router.Path("/").Handler(http.FileServer(http.Dir("client/build")))
		router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("client/build")))
	}

	server.UseHandler(router)
	return &server
}

func (server *Server) Start() {
	// go models.HandleUsers()
	log.Printf("Spinning up the server at %s%s...", server.host, server.port)
	err := http.ListenAndServe(server.port, server)
	log.Printf(err.Error())
}
