package server

import (
	"github.com/zmb3/spotify"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	STATE_KEY = "spotify_auth_state"
)

var (
	redirect = "http://localhost" + getPort() + "/callback"
	letters  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	auth     = spotify.NewAuthenticator(redirect, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistReadCollaborative)
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func generateRandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func initiateAuthHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandStr(128)
	cookie := http.Cookie{Name: STATE_KEY, Value: state}
	http.SetCookie(w, &cookie)
	url := auth.AuthURL(state)
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
	tok, err := auth.Token(storedState, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Print("Couldn't get token:")
		log.Print(err)
		return
	}

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	addSpotifyChan <- spotifySession{state: storedState, client: &client}
	http.Redirect(w, r, "playlists", http.StatusFound)
}
