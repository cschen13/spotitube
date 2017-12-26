package models

import (
	"encoding/json"
	"log"
	"net/http"
)

type Authenticator interface {
	BuildAuthURL(string) string
	GetType() string
	newClient(string, *http.Request) (Client, error)
}

type Client interface {
	GetPlaylists(string) (*PlaylistsPage, error)
	GetPlaylistInfo(string, string) (Playlist, error)
	CreatePlaylist(string) (Playlist, error)
	GetPlaylistTracks(Playlist, string) ([]PlaylistTrack, bool, error)
	InsertTrack(Playlist, PlaylistTrack) (bool, error)
}

type PlaylistsPage struct {
	Playlists         []Playlist
	PageNumber        int
	NextPageParam     string
	PreviousPageParam string
}

type ConvertPage struct {
	Playlist Playlist
	Tracks   []PlaylistTrack
}

type Playlist interface {
	GetID() string
	GetOwnerID() string
	GetName() string
	GetURL() string
	GetCoverURL() string
}

func (p PlaylistsPage) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["pageNumber"] = p.PageNumber
	m["nextPageParam"] = p.NextPageParam
	m["previousPageParam"] = p.PreviousPageParam

	m["playlists"] = make([]map[string]string, len(p.Playlists))
	mPlaylists := m["playlists"].([]map[string]string)
	for i, playlist := range p.Playlists {
		mPlaylists[i] = make(map[string]string)
		mPlaylists[i]["id"] = playlist.GetID()
		mPlaylists[i]["ownerId"] = playlist.GetOwnerID()
		mPlaylists[i]["name"] = playlist.GetName()
		mPlaylists[i]["url"] = playlist.GetURL()
		mPlaylists[i]["coverUrl"] = playlist.GetCoverURL()
	}

	return json.Marshal(m)
}

type PlaylistTrack interface {
	GetTitle() string
	GetArtist() string //Main artist
}
