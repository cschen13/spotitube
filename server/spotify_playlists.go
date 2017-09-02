package server

import (
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"strconv"
)

const (
	PAGE_PARAM = "page"
	PAGE_LIMIT = 21
)

type playlistPage struct {
	SimplePlaylistPage *spotify.SimplePlaylistPage
	PageNumber         int
	LastPage           bool
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
	client := getSpotify(state)
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
		renderErrorTemplate(w, "An error occurred while generating the playlists.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// for _, playlist := range page.SimplePlaylistPage.Playlists {
	// 	log.Printf("%v", playlist.Images)
	// }

	renderTemplate(w, "playlists", page)
}

func generatePlaylistPage(client *spotify.Client, pageNumberParam string) (page playlistPage, err error) {
	pageNumber, err := strconv.Atoi(pageNumberParam)
	if err != nil {
		return
	}

	// limit is one more than PAGE_LIMIT so we can "peek ahead"
	// and determine if this is the last page of playlists
	limit := PAGE_LIMIT + 1
	offset := (pageNumber - 1) * PAGE_LIMIT
	options := spotify.Options{Limit: &limit, Offset: &offset}
	simplePlaylistPage, err := client.CurrentUsersPlaylistsOpt(&options)
	if err != nil {
		return
	}

	page = playlistPage{PageNumber: pageNumber}
	numPlaylists := len(simplePlaylistPage.Playlists)
	if numPlaylists < PAGE_LIMIT+1 {
		page.LastPage = true
	} else {
		simplePlaylistPage.Playlists = simplePlaylistPage.Playlists[:numPlaylists-1]
	}

	page.SimplePlaylistPage = simplePlaylistPage
	return
}
