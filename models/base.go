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
}

type PlaylistsPage struct {
	Playlists         []PlaylistInfo
	PageNumber        int
	NextPageParam     string
	PreviousPageParam string
}

type PlaylistInfo interface {
	GetID() string
	GetName() string
	GetURL() string
	GetCoverURL() string
}
