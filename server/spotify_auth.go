package server

import (
	"github.com/zmb3/spotify"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	REDIRECT  = "http://localhost:8080/callback"
	STATE_KEY = "spotify_auth_state"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	auth    = spotify.NewAuthenticator(REDIRECT, spotify.ScopePlaylistReadPrivate)
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

func initiateAuth(w http.ResponseWriter, r *http.Request) {
	state := generateRandStr(128)
	cookie := http.Cookie{Name: STATE_KEY, Value: state}
	http.SetCookie(w, &cookie)
	url := auth.AuthURL(state)
	log.Printf("Redirecting user to %s", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(STATE_KEY)
	if err != nil {
		log.Print("No cookie for spotify auth state found")
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}

	storedState := cookie.Value
	tok, err := auth.Token(storedState, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Print(err)
		return
	}

	if st := r.FormValue("state"); st != storedState {
		http.NotFound(w, r)
		log.Print("State mismatch: %s != %s\n", st, storedState)
		return
	}

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	addSpotifyChan <- spotifySession{state: storedState, client: &client}
	http.Redirect(w, r, "playlists", http.StatusFound)
}
