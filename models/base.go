package models

import (
	"net/http"
)

type Authenticator interface {
	BuildAuthURL(string) string
	newClient(string, *http.Request) (Client, error)
}

type Client interface {
	GetPlaylists(string) (*PlaylistPage, error)
}

type PlaylistPage struct {
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
