package controllers

import (
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	STATE_KEY = "spotify_auth_state"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RegisterAuthController(router *mux.Router) {
	router.HandleFunc("/login", initiateAuthHandler)
	router.HandleFunc("/callback", completeAuthHandler)
}

func initiateAuthHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandStr(128)
	cookie := http.Cookie{Name: STATE_KEY, Value: state}
	http.SetCookie(w, &cookie)
	url := models.BuildSpotifyAuthURL(state)
	// url := auth.AuthURL(state)
	log.Printf("Redirecting user to %s", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func completeAuthHandler(w http.ResponseWriter, r *http.Request) {
	// check the request for a state cookie
	cookie, err := r.Cookie(STATE_KEY)
	if err != nil {
		log.Print("No cookie for spotify auth state found")
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}
	storedState := cookie.Value

	// acquire access token (also checks state parameter)
	user, err := models.NewUser(storedState, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Print("Couldn't create user:")
		log.Print(err)
		return
	}

	user.Add()
	http.Redirect(w, r, "playlists", http.StatusFound)
}
