package models

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

const SPOTIFY_PLAYLISTS_PAGE_LIMIT = 21

var (
	redirect = "http://localhost" + getPort() + "/callback"
	auth     = spotify.NewAuthenticator(redirect, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistReadCollaborative)
)

type spotifyClient struct {
	*spotify.Client
	token *oauth2.Token
}

func newSpotifyClient(state string, r *http.Request) (*spotifyClient, error) {
	// acquire access token (also checks state parameter)
	tok, err := auth.Token(state, r)
	if err != nil {
		return nil, err
	}

	client := auth.NewClient(tok)
	return &spotifyClient{&client, tok}, nil
}

type SpotifyPlaylist struct {
	ID           spotify.ID
	Name         string
	Images       []spotify.Image
	ExternalURLs map[string]string
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}

func BuildSpotifyAuthURL(state string) string {
	return auth.AuthURL(state)
}

func (client *spotifyClient) getPlaylists(pageNumber int) (playlists []SpotifyPlaylist, err error) {
	// limit is one more than PAGE_LIMIT so we can "peek ahead"
	// and determine if this is the last page of playlists
	limit := SPOTIFY_PLAYLISTS_PAGE_LIMIT + 1
	offset := (pageNumber - 1) * SPOTIFY_PLAYLISTS_PAGE_LIMIT
	options := spotify.Options{Limit: &limit, Offset: &offset}
	simplePlaylistPage, err := client.CurrentUsersPlaylistsOpt(&options)
	if err != nil {
		return
	}

	playlists = make([]SpotifyPlaylist, len(simplePlaylistPage.Playlists))
	for i, playlist := range simplePlaylistPage.Playlists {
		playlists[i] = SpotifyPlaylist{ID: playlist.ID, Name: playlist.Name, Images: playlist.Images, ExternalURLs: playlist.ExternalURLs}
	}
	return
}
