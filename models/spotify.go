package models

import (
	"errors"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

const (
	SPOTIFY_PLAYLISTS_PAGE_LIMIT       = 21
	SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT = 100
	SPOTIFY_SERVICE                    = "spotify"
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
		playlists[i] = &spotifyPlaylist{&playlist}
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

func (client *spotifyClient) GetPlaylistInfo(playlistId string) (Playlist, error) {
	user, err := client.CurrentUser()
	if err != nil {
		return nil, err
	}

	fullPlaylist, err := client.GetPlaylist(user.ID, spotify.ID(playlistId))
	if err != nil {
		return nil, err
	}

	return &spotifyPlaylist{&fullPlaylist.SimplePlaylist}, nil
}

// GetPlaylistTracks gets page number "page" of PlaylistTracks from playlist.
// Returns the slice of playlist tracks, along with a boolean indicating whether it is the last page
func (client *spotifyClient) GetPlaylistTracks(playlist Playlist, page string) (tracks []PlaylistTrack, lastPage bool, err error) {
	pageNumber := 1
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil {
			return
		}
	}

	user, err := client.CurrentUser()
	if err != nil {
		return
	}

	limit := SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT
	offset := (pageNumber - 1) * SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT
	options := spotify.Options{Limit: &limit, Offset: &offset}
	trackPage, err := client.GetPlaylistTracksOpt(user.ID, spotify.ID(playlist.GetID()), &options, "")
	if err != nil {
		return
	}

	tracks = make([]PlaylistTrack, len(trackPage.Tracks))
	for i, track := range trackPage.Tracks {
		tracks[i] = &spotifyTrack{&track.Track.SimpleTrack}
	}

	lastPage = len(tracks) < SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT
	return
}

func (client *spotifyClient) InsertTrack(playlist Playlist, track PlaylistTrack) error {
	return errors.New("Unimplemented")
}

type spotifyPlaylist struct {
	obj *spotify.SimplePlaylist
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

type spotifyTrack struct {
	obj *spotify.SimpleTrack
}

func (track *spotifyTrack) GetTitle() string {
	return track.obj.Name
}

func (track *spotifyTrack) GetArtist() string {
	if len(track.obj.Artists) > 0 {
		return track.obj.Artists[0].Name
	}
	return ""
}
