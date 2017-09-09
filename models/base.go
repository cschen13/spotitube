package models

import (
	"net/http"
)

type Authenticator interface {
	BuildAuthURL(string) string
	GetType() string
	newClient(string, *http.Request) (Client, error)
}

type Client interface {
	GetPlaylists(string) (*PlaylistsPage, error)
	GetPlaylistInfo(string) (Playlist, error)
	CreatePlaylist(string) (Playlist, error)
	GetPlaylistTracks(Playlist, string) ([]PlaylistTrack, bool, error)
	InsertTrack(Playlist, PlaylistTrack) error
}

type PlaylistsPage struct {
	Playlists         []Playlist
	PageNumber        int
	NextPageParam     string
	PreviousPageParam string
}

type Playlist interface {
	GetID() string
	GetName() string
	GetURL() string
	GetCoverURL() string
}

type PlaylistTrack interface {
	GetTitle() string
	GetArtist() string //Main artist
}
