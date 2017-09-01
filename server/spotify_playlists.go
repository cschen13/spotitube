package server

import (
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"html/template"
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
	cookie, err := r.Cookie(STATE_KEY)
	if err != nil {
		log.Print("No stored state found")
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}

	state := cookie.Value
	client := getSpotify(state)
	if client == nil {
		log.Print("No associated user found for state %s", state)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	params := mux.Vars(r)
	pageNumberParam, present := params[PAGE_PARAM]
	if !present {
		pageNumberParam = "1"
	}

	page, err := generatePlaylistPage(client, pageNumberParam)
	if err != nil {
		http.Error(w, "Error while preparing playlists", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// for _, playlist := range page.SimplePlaylistPage.Playlists {
	// 	log.Printf("%v", playlist.Images)
	// }
	funcMap := template.FuncMap{
		"add": func(i int, j int) int {
			return i + j
		},
	}
	renderTemplate(w, "playlists", page, funcMap)
}

func generatePlaylistPage(client *spotify.Client, pageNumberParam string) (page playlistPage, err error) {
	pageNumber, err := strconv.Atoi(pageNumberParam)
	if err != nil {
		return
	}

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
