package models

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
)

const SPOTIFY_PLAYLISTS_PAGE_LIMIT = 21
const SPOTIFY_CALLBACK_PATH = "/callback"

var spotifyPermissions = []string{
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopePlaylistReadCollaborative,
}

type SpotifyAuthenticator struct {
	auth spotify.Authenticator
}

func NewSpotifyAuthenticator(addr string) *SpotifyAuthenticator {
	return &SpotifyAuthenticator{
		auth: spotify.NewAuthenticator(addr+SPOTIFY_CALLBACK_PATH, spotifyPermissions...),
	}
}

type spotifyClient struct {
	*spotify.Client
	token *oauth2.Token
}

func (sa *SpotifyAuthenticator) newSpotifyClient(state string, r *http.Request) (*spotifyClient, error) {
	// acquire access token (also checks state parameter)
	tok, err := sa.auth.Token(state, r)
	if err != nil {
		return nil, err
	}

	client := sa.auth.NewClient(tok)
	return &spotifyClient{&client, tok}, nil
}

type SpotifyPlaylist struct {
	ID           spotify.ID
	Name         string
	Images       []spotify.Image
	ExternalURLs map[string]string
}

func (sa *SpotifyAuthenticator) BuildSpotifyAuthURL(state string) string {
	return sa.auth.AuthURL(state)
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
