package models

type ClientPlaylist interface {
	GetID() string
	GetOwnerID() string
	GetName() string
	GetURL() string
	GetCoverURL() string
}

type Playlist struct {
	ID       string
	OwnerID  string
	Name     string
	URL      string
	CoverURL string
}

func NewPlaylist(playlist ClientPlaylist) *Playlist {
	return &Playlist{
		playlist.GetID(),
		playlist.GetOwnerID(),
		playlist.GetName(),
		playlist.GetURL(),
		playlist.GetCoverURL(),
	}
}

type Playlists []*Playlist
