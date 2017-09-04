package controllers

import (
	"fmt"
	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	PAGE_PARAM = "page"
)

type playlistPage struct {
	Playlists  []models.SpotifyPlaylist
	PageNumber int
	LastPage   bool
}

func RegisterPlaylistController(router *mux.Router) {
	router.HandleFunc("/playlists", getPlaylistsHandler)
}

func getPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	// check the request for a state cookie
	cookie, err := r.Cookie(STATE_KEY)
	if err != nil {
		log.Print("No cookie for spotify auth state found")
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}
	state := cookie.Value

	// retrieve the Spotify client
	client := models.GetUser(state)
	if client == nil {
		log.Print("No associated user found for state %s", state)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// get page number from route
	pageNumberParam := r.URL.Query().Get(PAGE_PARAM)
	if pageNumberParam == "" {
		pageNumberParam = "1"
	}

	page, err := generatePlaylistPage(client, pageNumberParam)
	if err != nil {
		utils.RenderErrorTemplate(w, "An error occurred while generating the playlists.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	utils.RenderTemplate(w, "playlists", page)
}

func generatePlaylistPage(client *models.User, pageNumberParam string) (page playlistPage, err error) {
	pageNumber, err := strconv.Atoi(pageNumberParam)
	if err != nil {
		return
	}

	playlists, err := client.GetSpotifyPlaylists(pageNumber)
	if err != nil {
		return
	}

	// for _, playlist := range playlists {
	// 	log.Printf("%s", playlist.Name)
	// }

	page = playlistPage{PageNumber: pageNumber}
	numPlaylists := len(playlists)
	if numPlaylists < models.SPOTIFY_PLAYLISTS_PAGE_LIMIT+1 {
		page.LastPage = true
	} else {
		playlists = playlists[:numPlaylists-1]
	}

	page.Playlists = playlists
	return
}
