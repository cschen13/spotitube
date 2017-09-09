package models

import (
	"errors"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

const (
	SPOTIFY_PLAYLISTS_PAGE_LIMIT = 21
	SPOTIFY_SERVICE              = "spotify"
)

var spotifyPermissions = []string{
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopePlaylistReadCollaborative,
}

type SpotifyAuthenticator struct {
	auth spotify.Authenticator
}

func NewSpotifyAuthenticator(addr string) *SpotifyAuthenticator {
	return &SpotifyAuthenticator{
		auth: spotify.NewAuthenticator(addr+"/callback/"+SPOTIFY_SERVICE, spotifyPermissions...),
	}
}

func (sa *SpotifyAuthenticator) BuildAuthURL(state string) string {
	return sa.auth.AuthURL(state)
}

func (sa *SpotifyAuthenticator) newClient(state string, r *http.Request) (Client, error) {
	// acquire access token (also checks state parameter)
	tok, err := sa.auth.Token(state, r)
	if err != nil {
		return nil, err
	}

	client := sa.auth.NewClient(tok)
	return &spotifyClient{&client, tok}, nil
}

func (sa *SpotifyAuthenticator) GetType() string {
	return SPOTIFY_SERVICE
}

type spotifyClient struct {
	*spotify.Client
	token *oauth2.Token
}

func (client *spotifyClient) GetPlaylists(page string) (playlistPage *PlaylistsPage, err error) {
	pageNumber := 1
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil {
			return
		}
	}

	limit := SPOTIFY_PLAYLISTS_PAGE_LIMIT
	offset := (pageNumber - 1) * SPOTIFY_PLAYLISTS_PAGE_LIMIT
	options := spotify.Options{Limit: &limit, Offset: &offset}
	simplePlaylistPage, err := client.CurrentUsersPlaylistsOpt(&options)
	if err != nil {
		return
	}

	playlists := make([]Playlist, len(simplePlaylistPage.Playlists))
	for i, playlist := range simplePlaylistPage.Playlists {
		playlists[i] = &spotifyPlaylist{playlist}
	}

	playlistPage = &PlaylistsPage{Playlists: playlists, PageNumber: pageNumber}
	if simplePlaylistPage.Previous != "" {
		playlistPage.PreviousPageParam = strconv.Itoa(pageNumber - 1)
	}

	if simplePlaylistPage.Next != "" {
		playlistPage.NextPageParam = strconv.Itoa(pageNumber + 1)
	}

	return
}

func (client *spotifyClient) CreatePlaylist(name string) (Playlist, error) {
	return nil, errors.New("Unimplemented")
}

func (client *spotifyClient) GetPlaylistTracks(playlist Playlist) []PlaylistTrack {
	return nil
}

func (client *spotifyClient) InsertTrack(playlist Playlist, track PlaylistTrack) (bool, error) {
	return false, errors.New("Unimplemented")
}

type spotifyPlaylist struct {
	obj spotify.SimplePlaylist
}

func (playlist *spotifyPlaylist) GetID() string {
	return playlist.obj.ID.String()
}

func (playlist *spotifyPlaylist) GetName() string {
	return playlist.obj.Name
}

func (playlist *spotifyPlaylist) GetURL() string {
	if url, present := playlist.obj.ExternalURLs["spotify"]; present {
		return url
	}
	return ""
}

func (playlist *spotifyPlaylist) GetCoverURL() string {
	if len(playlist.obj.Images) > 0 {
		return playlist.obj.Images[0].URL
	}
	return ""
}
